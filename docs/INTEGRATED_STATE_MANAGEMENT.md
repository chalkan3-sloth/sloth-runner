# Integrated State Management System

## Visão Geral

O sloth-runner agora possui um sistema de gerenciamento de estado unificado, inspirado no Pulumi e Terraform, que integra o gerenciamento de stacks com o rastreamento de estados de workflow em um único backend SQLite.

## Arquitetura

### Componentes Principais

1. **StateBackend** (`internal/stack/state_backend.go`)
   - Backend unificado para gerenciamento de estado
   - Integra stacks, recursos e estados de workflow
   - Fornece versionamento, drift detection e rollback

2. **StackManager** (`internal/stack/stack.go`)
   - Gerencia stacks e recursos
   - Mantém metadados e histórico de execução
   - Suporta consultas e relacionamentos

3. **Migrator** (`internal/stack/migration.go`)
   - Ferramenta para migrar dados do sistema antigo
   - Consolida workflow_state em stacks
   - Preserva histórico e recursos

## Banco de Dados

### Tabelas Principais

#### Stacks
```sql
CREATE TABLE stacks (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT,
  version TEXT,
  status TEXT NOT NULL,
  created_at DATETIME,
  updated_at DATETIME,
  completed_at DATETIME,
  workflow_file TEXT,
  task_results TEXT,  -- JSON
  outputs TEXT,       -- JSON
  configuration TEXT, -- JSON
  metadata TEXT,      -- JSON
  execution_count INTEGER,
  last_duration INTEGER,
  last_error TEXT,
  resource_version TEXT
);
```

#### Resources
```sql
CREATE TABLE resources (
  id TEXT PRIMARY KEY,
  stack_id TEXT NOT NULL,
  type TEXT NOT NULL,
  name TEXT NOT NULL,
  module TEXT NOT NULL,
  properties TEXT,    -- JSON
  dependencies TEXT,  -- JSON array
  state TEXT NOT NULL,
  checksum TEXT,
  created_at DATETIME,
  updated_at DATETIME,
  last_applied DATETIME,
  error_message TEXT,
  metadata TEXT,      -- JSON
  FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE
);
```

#### State Versions (Snapshots)
```sql
CREATE TABLE state_versions (
  id TEXT PRIMARY KEY,
  stack_id TEXT NOT NULL,
  version INTEGER NOT NULL,
  state_snapshot TEXT NOT NULL, -- JSON snapshot completo
  created_at DATETIME,
  created_by TEXT,
  description TEXT,
  checksum TEXT,
  FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE,
  UNIQUE(stack_id, version)
);
```

#### Drift Detection
```sql
CREATE TABLE drift_detections (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  stack_id TEXT NOT NULL,
  resource_id TEXT NOT NULL,
  detected_at DATETIME,
  expected_state TEXT,     -- JSON
  actual_state TEXT,       -- JSON
  drift_fields TEXT,       -- JSON array
  is_drifted BOOLEAN,
  resolution_status TEXT,
  FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE,
  FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE
);
```

#### State Locks
```sql
CREATE TABLE state_locks (
  stack_id TEXT PRIMARY KEY,
  lock_id TEXT NOT NULL,
  operation TEXT NOT NULL,
  who TEXT NOT NULL,
  created_at DATETIME,
  expires_at DATETIME NOT NULL,
  info TEXT,
  FOREIGN KEY (stack_id) REFERENCES stacks(id) ON DELETE CASCADE
);
```

#### Resource Dependencies
```sql
CREATE TABLE resource_dependencies (
  resource_id TEXT NOT NULL,
  depends_on_id TEXT NOT NULL,
  dependency_type TEXT,  -- explicit, implicit
  created_at DATETIME,
  PRIMARY KEY (resource_id, depends_on_id),
  FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
  FOREIGN KEY (depends_on_id) REFERENCES resources(id) ON DELETE CASCADE
);
```

## Comandos CLI

### Gerenciamento de Stacks

```bash
# Listar todos os stacks
sloth-runner stack state list

# Mostrar detalhes de um stack
sloth-runner stack state show my-stack

# Listar versões de um stack
sloth-runner stack state versions my-stack

# Criar snapshot manual
sloth-runner stack state snapshot my-stack --description "Before deployment"
```

### Versionamento e Rollback

```bash
# Fazer rollback para versão anterior
sloth-runner stack state rollback my-stack --version 3

# Ver diferenças entre versões (futuro)
sloth-runner stack state diff my-stack --from 2 --to 3
```

### Drift Detection

```bash
# Detectar drift em recursos
sloth-runner stack state drift my-stack

# Formato JSON
sloth-runner stack state drift my-stack --output json
```

### Locking

```bash
# Adquirir lock
sloth-runner stack state lock my-stack --operation "deployment" --who "admin" --duration 30m

# Liberar lock
sloth-runner stack state unlock my-stack --lock-id <lock-id>
```

### Tags

```bash
# Adicionar tag
sloth-runner stack state tags add my-stack production

# Listar tags
sloth-runner stack state tags list my-stack
```

### Activity Log

```bash
# Ver log de atividades
sloth-runner stack state activity my-stack --limit 20
```

## Migração

### Migrar do Sistema Antigo

```bash
# Migrar dados do workflow_state para o novo sistema
sloth-runner stack migrate \
  --source ~/.sloth-runner/state.db \
  --target /etc/sloth-runner/stacks.db

# Dry run (sem modificar dados)
sloth-runner stack migrate --dry-run

# Gerar script SQL para migração manual
sloth-runner stack migrate --generate-script migration.sql
```

### Processo de Migração

1. **Backup**: Sempre faça backup dos bancos de dados antes
2. **Dry Run**: Execute com `--dry-run` primeiro
3. **Migração**: Execute a migração real
4. **Verificação**: Verifique os dados migrados
5. **Report**: Consulte o relatório em `migration_report.json`

## API Programática

### Usando o StateBackend em Go

```go
package main

import (
    "fmt"
    "github.com/chalkan3-sloth/sloth-runner/internal/stack"
)

func main() {
    // Criar backend
    backend, err := stack.NewStateBackend("/path/to/state.db")
    if err != nil {
        panic(err)
    }
    defer backend.Close()

    // Criar stack
    stackManager := backend.GetStackManager()
    myStack := &stack.StackState{
        ID:          "stack-123",
        Name:        "my-infrastructure",
        Description: "Production infrastructure",
        Version:     "1.0.0",
        Status:      "created",
        // ... outros campos
    }

    if err := stackManager.CreateStack(myStack); err != nil {
        panic(err)
    }

    // Criar snapshot
    version, err := backend.CreateSnapshot("stack-123", "admin", "Initial state")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Created snapshot version: %d\n", version)

    // Detectar drift
    expectedState := map[string]interface{}{
        "instances": 3,
        "region":    "us-east-1",
    }
    actualState := map[string]interface{}{
        "instances": 2,  // drift!
        "region":    "us-east-1",
    }

    if err := backend.DetectDrift("stack-123", "resource-456", expectedState, actualState); err != nil {
        panic(err)
    }

    // Obter informações de drift
    drifts, err := backend.GetDriftInfo("stack-123")
    if err != nil {
        panic(err)
    }

    for _, drift := range drifts {
        if drift.IsDrifted {
            fmt.Printf("Drift detected in resource %s\n", drift.ResourceID)
            fmt.Printf("Drifted fields: %v\n", drift.DriftedFields)
        }
    }

    // Rollback
    if err := backend.RollbackToSnapshot("stack-123", 1, "admin"); err != nil {
        panic(err)
    }

    // Tags
    backend.AddTag("stack-123", "production")
    backend.AddTag("stack-123", "critical")

    tags, _ := backend.GetTags("stack-123")
    fmt.Printf("Tags: %v\n", tags)
}
```

## Features Principais

### 1. Versionamento Automático

Cada mudança no stack cria automaticamente uma nova versão, permitindo:
- Histórico completo de mudanças
- Rollback para qualquer versão anterior
- Auditoria de alterações

### 2. Drift Detection

Detecta automaticamente quando o estado real difere do estado esperado:
- Compara expected vs actual state
- Identifica campos específicos que mudaram
- Status de resolução (pending, resolved, ignored)

### 3. State Locking

Previne modificações concorrentes:
- Locks com timeout automático
- Informações sobre quem/o quê está bloqueando
- Limpeza automática de locks expirados

### 4. Resource Dependencies

Rastreia dependências entre recursos:
- Dependências explícitas (definidas pelo usuário)
- Dependências implícitas (inferidas pelo sistema)
- Suporte para grafos de dependências

### 5. Activity Logging

Mantém log completo de todas as atividades:
- Quem fez o quê e quando
- Operações de create, update, delete, lock, unlock, rollback
- Útil para auditoria e debugging

### 6. Tags

Organização e categorização de stacks:
- Tags arbitrárias
- Busca por tags
- Útil para organizar ambientes (prod, dev, staging)

## Comparação com Terraform/Pulumi

| Feature | Terraform | Pulumi | sloth-runner |
|---------|-----------|--------|--------------|
| State Storage | Local/Remote | Local/Cloud | SQLite |
| Versioning | ✓ | ✓ | ✓ |
| Drift Detection | ✓ | ✓ | ✓ |
| State Locking | ✓ | ✓ | ✓ |
| Rollback | ✗ | ✗ | ✓ |
| Activity Log | ✗ | ✗ | ✓ |
| Resource Dependencies | ✓ | ✓ | ✓ |
| Tags | ✓ | ✓ | ✓ |

## Benefícios

1. **Simplicidade**: Tudo em um único banco SQLite, sem necessidade de servidores externos
2. **Performance**: SQLite é rápido e eficiente para a maioria dos casos de uso
3. **Portabilidade**: Fácil de fazer backup e mover entre ambientes
4. **Integração**: Totalmente integrado com o sistema de workflows do sloth-runner
5. **Auditoria**: Log completo de todas as atividades
6. **Flexibilidade**: Suporta qualquer tipo de recurso ou stack

## Limitações e Considerações

1. **Concorrência**: SQLite tem limitações de concorrência, mas o state locking ajuda
2. **Tamanho**: Para stacks muito grandes (>1GB), considere outras opções
3. **Distribuição**: SQLite é local, não distribuído (mas pode ser replicado)
4. **Backup**: Importante fazer backups regulares do banco de dados

## Próximos Passos

1. ✓ Sistema de migração
2. ✓ Comandos CLI
3. ✓ Testes automatizados
4. ✓ Documentação
5. ⚠ Integração com workflows existentes
6. ⚠ Backup automático
7. ⚠ Compressão de snapshots antigos
8. ⚠ Dashboard web para visualização

## Suporte

Para questões ou problemas:
- Issues: https://github.com/chalkan3-sloth/sloth-runner/issues
- Documentação: docs/
- Exemplos: examples/

## Referências

- [Terraform State](https://www.terraform.io/docs/language/state/index.html)
- [Pulumi State and Backends](https://www.pulumi.com/docs/intro/concepts/state/)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
