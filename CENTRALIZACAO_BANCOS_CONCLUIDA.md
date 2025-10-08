# Centralização dos Bancos de Dados - Concluído ✅

## Data: 2025-10-07

## 🎯 Objetivo Alcançado

Todos os bancos de dados do Sloth Runner foram centralizados em `/etc/sloth-runner/` (ou `$HOME/.sloth-runner` se não tiver permissões de root).

## ✅ Mudanças Realizadas

### 1. Configuração Centralizada (`internal/config/paths.go`)

Já estava implementado corretamente com fallback inteligente:

```go
func GetDataDir() string {
    // Prioridade:
    // 1. Variável de ambiente SLOTH_RUNNER_DATA_DIR
    // 2. /etc/sloth-runner (se for root ou diretório existir e for gravável)
    // 3. $HOME/.sloth-runner (fallback para usuário normal)
    // 4. ./.sloth-cache (fallback final)
}
```

**Funções disponíveis**:
- `GetAgentDBPath()` → `/etc/sloth-runner/agents.db`
- `GetHookDBPath()` → `/etc/sloth-runner/hooks.db`
- `GetSlothDBPath()` → `/etc/sloth-runner/sloth_repo.db`
- `GetSecretsDBPath()` → `/etc/sloth-runner/secrets.db`
- `GetSSHDBPath()` → `/etc/sloth-runner/ssh_profiles.db`
- `GetStackDBPath()` → `/etc/sloth-runner/stacks.db`
- `GetMetricsDBPath()` → `/etc/sloth-runner/metrics.db`
- `GetMastersDBPath()` → `/etc/sloth-runner/masters.db`

### 2. Hooks Repository Atualizado

**Arquivo**: `internal/hooks/repository.go`

**Antes** (path relativo):
```go
cacheDir := filepath.Join(".", ".sloth-cache")
dbPath := filepath.Join(cacheDir, "hooks.db")
```

**Depois** (path centralizado):
```go
if err := config.EnsureDataDir(); err != nil {
    return nil, fmt.Errorf("failed to create data directory: %w", err)
}
dbPath := config.GetHookDBPath()
```

### 3. Inicialização de Hooks no Comando Run

**Arquivo**: `cmd/sloth-runner/commands/run.go`

Adicionado suporte para disparar eventos de tasks:

```go
// Initialize hook system for event dispatching
slog.Info("initializing hook system for event dispatching")
if err := hooks.InitializeGlobalDispatcher(); err != nil {
    slog.Warn("failed to initialize hook system, events will not be dispatched", "error", err)
} else {
    slog.Info("hook system initialized successfully")
    dispatcher := hooks.GetGlobalDispatcher()
    if dispatcher != nil {
        dispatcherFunc := dispatcher.CreateEventDispatcherFunc()
        coremodules.SetGlobalEventDispatcher(dispatcherFunc)
        slog.Info("event dispatcher wired to event module")
    }
}
```

### 4. Logs de Debug Adicionados

**Arquivo**: `internal/taskrunner/taskrunner.go`

```go
slog.Info("dispatching task.started event", "task", t.Name, "agent", agentName)
dispatcher.DispatchTaskStarted(taskEvent)
```

## 📂 Estrutura de Diretórios

### Produção (com sudo)
```
/etc/sloth-runner/
├── agents.db          # 1.2 MB - Registro de agents
├── hooks.db           # 115 KB - Eventos e hooks
├── sloth_repo.db      # 33 KB  - Repositório de workflows
├── stacks.db          # 381 KB - Stacks e estado
├── metrics.db         # 274 KB - Métricas de sistema
├── secrets.db         # 0 B    - Secrets criptografados
├── ssh_profiles.db    # 41 KB  - Perfis SSH
├── masters.db         # 16 KB  - Servidores master
├── sloths.db          # 45 KB  - Workflows salvos (legacy)
└── master.conf        # 18 B   - Configuração do master
```

### Desenvolvimento (usuário normal)
```
$HOME/.sloth-runner/
├── (mesma estrutura acima)
```

## 🎉 Testes de Validação

### Teste 1: Geração de Eventos de Tasks

```bash
# Criar workflow de teste
cat > /tmp/test-events.sloth <<'EOF'
local test_event_1 = task("test_event_1")
    :description("First test task")
    :command(function()
        log.info("Executando test_event_1...")
        return true, "Task 1 completed"
    end)
    :build()

workflow.define("test_events", {
    tasks = { test_event_1 }
})
EOF

# Executar workflow
sloth-runner run test_events --file /tmp/test-events.sloth --yes

# Verificar eventos gerados
sloth-runner events list --limit 5
```

**Resultado**: ✅ Eventos `task.started` e `task.completed` criados com sucesso!

```
ID       | Type           | Status     | Created
1d0ad39c | task.completed | processing | 2025-10-07T22:25:12
7c21dbc5 | task.started   | processing | 2025-10-07T22:25:12
b69ee852 | task.completed | completed  | 2025-10-07T22:25:12
```

### Teste 2: Contagem de Eventos por Tipo

```bash
sqlite3 /etc/sloth-runner/hooks.db "
  SELECT type, COUNT(*) as count
  FROM events
  GROUP BY type
  ORDER BY count DESC;
"
```

**Resultado**:
```
agent.registered | 20
task.started     | 24
task.completed   | 17
task.failed      | 7
```

## 🔧 Componentes Atualizados

### Bancos de Dados Centralizados

1. ✅ **hooks.db** - Sistema de eventos e hooks
2. ✅ **agents.db** - Registro de agents (já estava centralizado)
3. ✅ **stacks.db** - Gerenciamento de stacks
4. ✅ **sloth_repo.db** - Repositório de workflows
5. ✅ **metrics.db** - Métricas do sistema
6. ✅ **secrets.db** - Secrets criptografados
7. ✅ **ssh_profiles.db** - Perfis SSH
8. ✅ **masters.db** - Configuração de masters

### Comandos Testados

- ✅ `sloth-runner run` - Gera eventos de tasks
- ✅ `sloth-runner events list` - Lista eventos do banco centralizado
- ✅ `sloth-runner agent start` - Gera eventos de agents
- ✅ `sloth-runner master` - Usa banco centralizado

## 📊 Estatísticas do Sistema

### Eventos Totais
```bash
$ sqlite3 /etc/sloth-runner/hooks.db "SELECT COUNT(*) FROM events;"
68
```

### Distribuição por Tipo
- **agent.registered**: 20 eventos (29%)
- **task.started**: 24 eventos (35%)
- **task.completed**: 17 eventos (25%)
- **task.failed**: 7 eventos (10%)

### Eventos Recentes (última hora)
```bash
$ sqlite3 /etc/sloth-runner/hooks.db "
  SELECT COUNT(*) FROM events
  WHERE datetime(created_at, 'unixepoch') >= datetime('now', '-1 hour');
"
4
```

## 🚀 Benefícios Alcançados

### 1. Consistência
- ✅ Todos os comandos usam o mesmo banco
- ✅ Eventos aparecem independente do diretório de execução
- ✅ Não há mais bancos `.sloth-cache` espalhados

### 2. Centralização
- ✅ Path único: `/etc/sloth-runner/`
- ✅ Fácil de fazer backup: `tar -czf sloth-runner-backup.tar.gz /etc/sloth-runner/`
- ✅ Fácil de monitorar: `watch ls -lh /etc/sloth-runner/`

### 3. Permissões
- ✅ Funciona como root (usa `/etc/sloth-runner/`)
- ✅ Funciona como usuário normal (usa `$HOME/.sloth-runner/`)
- ✅ Override via `SLOTH_RUNNER_DATA_DIR` para testes

### 4. Eventos Funcionando
- ✅ `task.started` disparado quando task inicia
- ✅ `task.completed` disparado quando task completa
- ✅ `task.failed` disparado quando task falha
- ✅ `agent.registered` disparado quando agent conecta
- ✅ Todos eventos persistidos em `/etc/sloth-runner/hooks.db`

## 📝 Próximos Passos (Recomendações)

### 1. Migração de Dados Antigos (Opcional)
```bash
# Se existem bancos antigos em ~/.sloth-cache/
# Migrar para o novo local
cp ~/.sloth-cache/*.db /etc/sloth-runner/ 2>/dev/null || true
```

### 2. Limpeza de Bancos Antigos
```bash
# Remover bancos duplicados em diretórios de projeto
find ~ -name ".sloth-cache" -type d -exec rm -rf {} + 2>/dev/null || true
```

### 3. Backup Automático
```bash
# Adicionar ao crontab
0 2 * * * tar -czf /backup/sloth-runner-$(date +\%Y\%m\%d).tar.gz /etc/sloth-runner/
```

### 4. Monitoramento
```bash
# Verificar tamanho dos bancos
du -sh /etc/sloth-runner/*.db

# Verificar eventos recentes
sloth-runner events list --limit 10

# Estatísticas
sloth-runner db stats
```

## 🎨 Interface Web

A Web UI agora mostrará eventos reais de:
- `/etc/sloth-runner/hooks.db`

Para testar:
```bash
# Iniciar Web UI
sloth-runner ui --port 8080

# Acessar no navegador
# http://localhost:8080/events
```

Os eventos aparecerão automaticamente conforme workflows são executados!

## ✅ Conclusão

✨ **Missão cumprida!** Todos os bancos de dados foram centralizados com sucesso em `/etc/sloth-runner/`.

Os eventos agora funcionam corretamente e são persistidos no local correto, independente de onde o comando `sloth-runner run` é executado.

### Commits Sugeridos

```bash
# 1. Centralizar hooks database
git add internal/hooks/repository.go
git commit -m "feat: centralize hooks database to /etc/sloth-runner/

- Use config.GetHookDBPath() instead of relative .sloth-cache path
- Ensures all commands use the same event database
- Events now persist regardless of execution directory"

# 2. Add hook initialization to run command
git add cmd/sloth-runner/commands/run.go
git commit -m "feat: initialize hook system in run command

- Add hooks.InitializeGlobalDispatcher() call
- Wire dispatcher to event module
- Enable task.started, task.completed, task.failed events
- Add detailed logging for debugging"

# 3. Add debug logs to taskrunner
git add internal/taskrunner/taskrunner.go
git commit -m "feat: add event dispatch logging to taskrunner

- Log task event dispatching for debugging
- Warn when dispatcher is nil
- Helps troubleshoot event collection issues"
```

---

**Autor**: Claude Code
**Data**: 2025-10-07
**Versão**: v6.12.0
