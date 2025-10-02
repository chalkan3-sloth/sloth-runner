---
title: Usando Values com delegate_to
description: Como usar values.yaml para controlar dinamicamente onde as tasks são executadas
---

# Usando Values com delegate_to

O `delegate_to` agora suporta valores dinâmicos através de variáveis locais que capturam dados de `values.yaml`. Isso permite que você execute as mesmas tasks em diferentes ambientes apenas mudando o arquivo de values.

## Conceito Básico

### 1. Crie um arquivo values.yaml

```yaml
host: mariaguica
environment: production
region: us-east-1
agents:
  web: mariaguica
  db: lady-guica
  cache: keite-guica
```

### 2. Use os valores no seu script

```lua
-- ✨ Capture o valor em uma variável local
local target_agent = values and values.host or "default-agent"

-- Use a variável no delegate_to
local my_task = task("my_task")
    :description("Task that runs on dynamic agent")
    :command(function(this, params)
        log.info("🚀 Running on remote agent")
        return true
    end)
    :delegate_to(target_agent)  -- ✨ Usa a variável
    :build()
```

### 3. Execute com o values

```bash
sloth-runner run -f script.sloth workflow_name --values values.yaml
```

## Exemplo Rápido

```lua
-- Captura o host do values
local agent = values and values.host or "mariaguica"

local test_task = task("test")
    :description("Test delegate_to with values")
    :command(function(this, params)
        local hostname, _, _ = exec.run("hostname")
        log.info("Running on: " .. hostname)
        return true
    end)
    :delegate_to(agent)  -- ✨ Usa valor de values.host
    :build()

workflow.define("test")
    :tasks({ test_task })
```

Execute com:
```bash
sloth-runner run -f test.sloth test --values values.yaml
```

## Resumo

O suporte a `values` no `delegate_to` torna o Sloth Runner extremamente flexível para:

- ✅ Deploy multi-ambiente (dev, staging, prod)
- ✅ Teste A/B em diferentes servers
- ✅ Blue/Green deployments
- ✅ Disaster recovery e failover
- ✅ Geographic distribution
- ✅ Load testing em diferentes clusters

Você escreve o script uma vez e reutiliza em todos os ambientes! 🎉
