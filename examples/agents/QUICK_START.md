# 🚀 Guia Rápido: Execução Remota via CMD

Este guia mostra como executar tarefas remotamente nos agentes usando o método que **FUNCIONA**.

## ✅ Método Funcional: agent run via CMD

O método funcional usa `sloth-runner agent run` dentro das tasks Lua para executar comandos nos agentes remotos.

## 📝 Template Básico

```lua
TaskDefinitions = {
    minha_task_remota = {
        description = "Descrição da task",
        tasks = {
            {
                name = "nome_da_task",
                description = "O que esta task faz",
                command = function()
                    log.info("🚀 Executando comando remoto...")
                    
                    -- Comando remoto via agent run
                    local cmd = "./sloth-runner agent run <AGENT_NAME> \"<COMANDO>\" --master <MASTER_ADDRESS>"
                    local output, err, failed = exec.run(cmd)
                    
                    if not failed then
                        log.info("✅ Sucesso!")
                        log.info("📋 Saída: " .. output)
                        return true, "Sucesso"
                    else
                        log.error("❌ Falha: " .. (err or "erro"))
                        return false, "Erro"
                    end
                end,
                timeout = "60s"
            }
        }
    }
}
```

## 🎯 Exemplos Prontos

### Exemplo 1: Hostname Remoto

```lua
TaskDefinitions = {
    hostname_remoto = {
        description = "Verifica hostname do agente",
        tasks = {
            {
                name = "get_hostname",
                description = "Pega hostname do agente",
                command = function()
                    log.info("🖥️  Verificando hostname...")
                    
                    local cmd = "./sloth-runner agent run ladyguica \"hostname\" --master 192.168.1.29:50053"
                    local output, err, failed = exec.run(cmd)
                    
                    if not failed then
                        log.info("✅ Hostname: " .. output)
                        return true, "OK"
                    end
                    return false, "Erro"
                end,
                timeout = "30s"
            }
        }
    }
}
```

### Exemplo 2: Listar Arquivos

```lua
TaskDefinitions = {
    list_files = {
        description = "Lista arquivos do agente",
        tasks = {
            {
                name = "ls_home",
                description = "Lista HOME do agente",
                command = function()
                    log.info("📂 Listando arquivos...")
                    
                    local cmd = "./sloth-runner agent run keiteguica \"ls -lah $HOME | head -10\" --master 192.168.1.29:50053"
                    local output, err, failed = exec.run(cmd)
                    
                    if not failed then
                        log.info("✅ Arquivos:\n" .. output)
                        return true, "OK"
                    end
                    return false, "Erro"
                end,
                timeout = "30s"
            }
        }
    }
}
```

### Exemplo 3: Múltiplos Agentes

```lua
TaskDefinitions = {
    check_all = {
        description = "Verifica todos os agentes",
        tasks = {
            {
                name = "check_all_agents",
                description = "Verifica todos os agentes",
                command = function()
                    local agents = {"ladyguica", "keiteguica"}
                    local master = "192.168.1.29:50053"
                    
                    for _, agent in ipairs(agents) do
                        log.info("Verificando " .. agent .. "...")
                        local cmd = "./sloth-runner agent run " .. agent .. " \"hostname\" --master " .. master
                        local output, err, failed = exec.run(cmd)
                        
                        if not failed then
                            log.info("  ✅ " .. agent .. ": " .. output:gsub("\n", ""))
                        else
                            log.error("  ❌ " .. agent .. ": ERRO")
                        end
                    end
                    
                    return true, "Verificação completa"
                end,
                timeout = "120s"
            }
        }
    }
}
```

## 🔧 Como Usar

1. **Criar arquivo .sloth**
   ```bash
   # Copie um dos exemplos acima para um arquivo
   nano meu_exemplo.sloth
   ```

2. **Executar**
   ```bash
   ./sloth-runner run -f meu_exemplo.sloth <nome_do_taskgroup>
   ```

## 📚 Exemplos Completos Disponíveis

### Exemplo Simples
```bash
./sloth-runner run -f examples/agents/functional_cmd_example.sloth remote_via_cmd
```

Saída esperada:
```
✅ Hostname: ladyguica
✅ Arquivos listados!
✅ Sistema verificado!

# Execution Summary
Task            | Status  | Duration    | Error
check_system    | Success | 60.187084ms |      
hostname_remoto | Success | 84.870375ms |      
list_home_files | Success | 39.32125ms  |
```

### Exemplo Completo (Infraestrutura)
```bash
./sloth-runner run -f examples/agents/complete_infrastructure_check.sloth distributed_infrastructure_check
```

Este exemplo completo executa:
- ✅ Verificação de conectividade
- ✅ Coleta de informações do sistema
- ✅ Verificação de recursos (CPU, memória, disco)
- ✅ Verificação de serviços
- ✅ Teste de performance básico
- ✅ Geração de relatório final

Saída esperada:
```
============================================================
📋 RELATÓRIO FINAL DE INFRAESTRUTURA
============================================================

✅ Conectividade verificada
✅ Informações do sistema coletadas
✅ Recursos verificados
✅ Serviços verificados
✅ Performance testada

🎉 Todos os agentes estão funcionando corretamente!
============================================================

# Execution Summary
Task                   | Status  | Duration     | Error
check_connectivity     | Success | 103.653125ms |      
gather_system_info     | Success | 525.361208ms |      
check_resources        | Success | 367.018542ms |      
check_services         | Success | 491.255417ms |      
basic_performance_test | Success | 9.330437541s |      
generate_report        | Success | 2.246542ms   |
```

## ⚙️ Pré-requisitos

### 1. Master rodando
```bash
./sloth-runner master start --port 50053 --daemon
```

Verificar:
```bash
ps aux | grep "sloth-runner master"
```

### 2. Agentes registrados
```bash
# Iniciar agentes (se não estiverem rodando)
./sloth-runner agent start --name ladyguica --master 192.168.1.29:50053 --port 50051 --daemon
./sloth-runner agent start --name keiteguica --master 192.168.1.29:50053 --port 50051 --daemon
```

### 3. Verificar agentes
```bash
./sloth-runner agent list --master 192.168.1.29:50053
```

Saída esperada:
```
AGENT NAME     ADDRESS              STATUS            LAST HEARTBEAT
------------   ----------           ------            --------------
keiteguica     192.168.1.17:50051   Active   2025-10-01T12:27:35-03:00
ladyguica      192.168.1.16:50051   Active   2025-10-01T12:27:35-03:00
```

## 💡 Dicas

### ✅ Funciona

```lua
-- Comando simples
"./sloth-runner agent run ladyguica \"hostname\" --master 192.168.1.29:50053"

-- Múltiplos comandos (use &&)
"./sloth-runner agent run ladyguica \"hostname && whoami && date\" --master 192.168.1.29:50053"

-- Pipes e redirecionamento
"./sloth-runner agent run keiteguica \"ls -la $HOME | head -10\" --master 192.168.1.29:50053"

-- Variáveis de ambiente
"./sloth-runner agent run ladyguica \"echo $HOME\" --master 192.168.1.29:50053"
```

### ❌ Não Funciona

```lua
-- ❌ Aspas simples no comando (use aspas duplas)
"./sloth-runner agent run ladyguica 'hostname' --master 192.168.1.29:50053"

-- ❌ Scripts multi-linha diretos
"./sloth-runner agent run ladyguica echo 'linha1' echo 'linha2' --master 192.168.1.29:50053"
```

### 🔨 Workarounds para Scripts Multi-linha

Se você precisa executar um script com múltiplas linhas:

**Opção 1: Use && para separar comandos**
```lua
local cmd = "./sloth-runner agent run ladyguica \"echo 'linha1' && echo 'linha2' && echo 'linha3'\" --master 192.168.1.29:50053"
```

**Opção 2: Crie script temporário no agente**
```lua
-- 1. Criar script
local create_script = "./sloth-runner agent run ladyguica \"cat > /tmp/myscript.sh << 'EOF'\\necho 'linha1'\\necho 'linha2'\\nEOF\" --master 192.168.1.29:50053"
exec.run(create_script)

-- 2. Executar
local run_script = "./sloth-runner agent run ladyguica \"bash /tmp/myscript.sh\" --master 192.168.1.29:50053"
local output, err, failed = exec.run(run_script)

-- 3. Limpar
exec.run("./sloth-runner agent run ladyguica \"rm /tmp/myscript.sh\" --master 192.168.1.29:50053")
```

## 🐛 Troubleshooting

### Erro: "failed to connect to master"
```bash
# Verificar se o master está rodando
./sloth-runner agent list --master 192.168.1.29:50053

# Se não estiver, iniciar:
./sloth-runner master start --port 50053 --daemon
```

### Erro: "failed to resolve agent"
```bash
# Listar agentes disponíveis
./sloth-runner agent list --master 192.168.1.29:50053

# Se o agente não aparecer, registrar:
./sloth-runner agent start --name <agent_name> --master 192.168.1.29:50053 --daemon
```

### Erro: "Command failed on agent"
```bash
# Testar comando diretamente
./sloth-runner agent run <agent_name> "hostname" --master 192.168.1.29:50053

# Verificar logs do agente
tail -f agent.log
```

## 📖 Referências

- [README Completo](./README_CMD_FUNCIONAL.md)
- [Exemplo Funcional](./functional_cmd_example.sloth)
- [Exemplo Completo de Infraestrutura](./complete_infrastructure_check.sloth)

---

**Status**: ✅ TESTADO E FUNCIONANDO  
**Data**: 2025-10-01  
**Versão**: 1.0
