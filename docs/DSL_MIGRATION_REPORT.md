# 🔄 Relatório de Migração do DSL

## 📊 Situação Atual

**Total de arquivos com DSL antigo encontrados:** 30 arquivos

### 🎯 Padrões a Migrar

**DSL Antigo (Deprecated):**
```lua
local my_task = task("name", function(this, params)
    -- código
end)

local my_workflow = workflow.create({
    tasks = {my_task}
})
```

**DSL Moderno (Atual):**
```lua
local my_task = task("name")
    :description("Descrição")
    :command(function(this, params)
        -- código
    end)
    :timeout("5m")
    :on_success(function(this, params, output)
        -- código
    end)
    :on_fail(function(this, params, output)
        -- código
    end)
    :build()

workflow.define("name")
    :description("Descrição")
    :version("1.0.0")
    :tasks({ my_task })
    :config({
        timeout = "10m"
    })
    :on_complete(function(success, results)
        -- código
    end)
```

---

## 📋 Arquivos a Atualizar

### Prioridade ALTA (Documentação Principal)

| # | Arquivo | Status | Prioridade |
|---|---------|--------|------------|
| 1 | `docs/README.md` | 🔄 Em progresso | 🔥🔥🔥 |
| 2 | `docs/index.md` | ⏳ Pendente | 🔥🔥🔥 |
| 3 | `docs/modern-dsl/introduction.md` | ⏳ Pendente | 🔥🔥🔥 |
| 4 | `docs/modern-dsl/index.md` | ⏳ Pendente | 🔥🔥🔥 |
| 5 | `docs/pt/quick-start.md` | ⏳ Pendente | 🔥🔥 |
| 6 | `docs/en/core-concepts.md` | ⏳ Pendente | 🔥🔥 |
| 7 | `docs/pt/core-concepts.md` | ⏳ Pendente | 🔥🔥 |

### Prioridade MÉDIA (Tutoriais e Guias)

| # | Arquivo | Status | Prioridade |
|---|---------|--------|------------|
| 8 | `docs/pt/contributing.md` | ⏳ Pendente | 🔥 |
| 9 | `docs/en/plugin-development.md` | ⏳ Pendente | 🔥 |
| 10 | `docs/pt/plugin-development.md` | ⏳ Pendente | 🔥 |
| 11 | `docs/commands/run.md` | ⏳ Pendente | 🔥 |
| 12 | `docs/en/security.md` | ⏳ Pendente | 🔥 |
| 13 | `docs/en/stack-management.md` | ⏳ Pendente | 🔥 |

### Prioridade BAIXA (Documentação de Módulos)

| # | Arquivo | Status |
|---|---------|--------|
| 14-30 | Vários `modules/*.md` | ⏳ Pendente |

---

## 🚀 Plano de Migração

### Fase 1 - Documentação Principal (HOJE)
- [ ] `docs/README.md` - Entrada principal do projeto
- [ ] `docs/modern-dsl/introduction.md` - Introdução ao DSL moderno
- [ ] `docs/pt/quick-start.md` - Quick start em português
- [ ] `docs/en/core-concepts.md` - Conceitos principais

**Meta:** 4 arquivos principais atualizados

### Fase 2 - Módulos Críticos (PRÓXIMO)
- [ ] `docs/modules/goroutine.md` - Exemplo de paralelismo
- [ ] `docs/modules/gitops.md` - GitOps workflows
- [ ] `docs/modules/pkg.md` - Package management
- [ ] `docs/modules/systemd.md` - Service management

**Meta:** 4 módulos com exemplos práticos

### Fase 3 - Resto da Documentação (DEPOIS)
- [ ] Restantes 22 arquivos
- [ ] Verificação final
- [ ] Build e teste

---

## 📝 Template de Migração

### Padrão de Transformação

**ANTES:**
```lua
local deploy_task = task("deploy", function(this, params)
    log.info("Deploying...")
    return true, "Success"
end)

local my_workflow = workflow.create({
    name = "deploy_app",
    tasks = {deploy_task}
})
```

**DEPOIS:**
```lua
local deploy_task = task("deploy")
    :description("Deploy application")
    :command(function(this, params)
        log.info("Deploying...")
        return true, "Success"
    end)
    :timeout("5m")
    :build()

workflow.define("deploy_app")
    :description("Deploy application workflow")
    :version("1.0.0")
    :tasks({ deploy_task })
    :config({
        timeout = "10m"
    })
```

---

## ✅ Checklist de Migração

Para cada arquivo:

- [ ] Ler arquivo completo
- [ ] Identificar todos os blocos com DSL antigo
- [ ] Transformar tasks:
  - [ ] Adicionar `:description()`
  - [ ] Converter function para `:command()`
  - [ ] Adicionar `:timeout()`
  - [ ] Adicionar `:on_success()` se relevante
  - [ ] Adicionar `:on_fail()` se relevante
  - [ ] Adicionar `:build()`
- [ ] Transformar workflows:
  - [ ] Mudar de `workflow.create()` para `workflow.define()`
  - [ ] Adicionar `:description()`
  - [ ] Adicionar `:version()`
  - [ ] Adicionar `:tasks()`
  - [ ] Adicionar `:config()`
  - [ ] Adicionar `:on_complete()`
- [ ] Verificar indentação
- [ ] Testar se código faz sentido
- [ ] Commit mudanças

---

## 🎯 Progresso

```
╔══════════════════════════════════════════════════════╗
║                                                      ║
║  MIGRAÇÃO DSL - PROGRESSO                          ║
║                                                      ║
║  Total de arquivos: 30                              ║
║  Concluídos:        0  ░░░░░░░░░░░░░░░░░░░  0%     ║
║  Em progresso:      1  █░░░░░░░░░░░░░░░░░░  3%     ║
║  Pendentes:        29  ░░░░░░░░░░░░░░░░░░░ 97%     ║
║                                                      ║
╚══════════════════════════════════════════════════════╝
```

---

## 📞 Notas

**Criado em:** 2025-10-10
**Status:** 🔄 Em Progresso
**Próxima revisão:** Após Fase 1

---

## 🔧 Script de Ajuda

Para buscar padrões antigos:
```bash
# Encontrar tasks no formato antigo
grep -r "local.*= task(" docs/ --include="*.md"

# Encontrar workflows no formato antigo
grep -r "workflow.create" docs/ --include="*.md"

# Listar arquivos com DSL
find docs -name "*.md" -exec grep -l "task(" {} \;
```

Para verificar após migração:
```bash
# Verificar se ainda tem padrões antigos
grep -r "task(\".*\", function" docs/ --include="*.md"

# Verificar se tem novo padrão
grep -r ":command(function" docs/ --include="*.md"
```

