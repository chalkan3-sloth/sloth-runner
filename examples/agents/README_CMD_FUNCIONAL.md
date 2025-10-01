# ✅ Execução Remota FUNCIONAL via CMD

Este documento explica como executar tarefas remotamente nos agentes usando o comando `sloth-runner agent run`, que funciona perfeitamente.

## Como Funciona

O método funcional usa o comando `sloth-runner agent run` dentro das tarefas Lua. Este comando:
1. Conecta-se ao master para resolver o nome do agente
2. Executa o comando shell diretamente no agente via gRPC
3. Retorna a saída em tempo real

## Exemplo Funcional

Arquivo: `examples/agents/functional_cmd_example.sloth`

```lua
TaskDefinitions = {
    remote_via_cmd = {
        description = "Execução remota FUNCIONAL via CMD",
        tasks = {
            {
                name = "hostname_remoto",
                description = "Mostra hostname do agente remoto",
                command = function()
                    log.info("🚀 Executando hostname no agente ladyguica...")
                    
                    local sloth_bin = "./sloth-runner"
                    local cmd = sloth_bin .. " agent run ladyguica \"hostname && whoami\" --master 192.168.1.29:50053"
                    
                    local output, err, failed = exec.run(cmd)
                    
                    if not failed then
                        log.info("✅ Comando executado com sucesso!")
                        log.info("📋 Saída:\n" .. output)
                        return true, "Sucesso"
                    else
                        log.error("❌ Falha: " .. (err or "erro"))
                        return false, "Erro"
                    end
                end,
                timeout = "60s"
            },
            {
                name = "list_home_files",
                description = "Lista arquivos do HOME no agente",
                command = function()
                    log.info("📂 Listando arquivos no agente keiteguica...")
                    
                    local sloth_bin = "./sloth-runner"
                    local cmd = sloth_bin .. " agent run keiteguica \"ls -lah $HOME | head -10\" --master 192.168.1.29:50053"
                    
                    local output, err, failed = exec.run(cmd)
                    
                    if not failed then
                        log.info("✅ Listagem realizada!")
                        log.info("📋 Arquivos:\n" .. output)
                        return true, "Listagem OK"
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

## Como Executar

```bash
# Execute o workflow
./sloth-runner run -f examples/agents/functional_cmd_example.sloth remote_via_cmd
```

## Saída Esperada

```
✅ Comando executado com sucesso!
📋 Saída:
 INFO  🚀 Executing on agent: ladyguica
 INFO  📝 Command: hostname && whoami
ladyguica
chalkan3
 SUCCESS  ✅ Command completed successfully on agent ladyguica

# Execution Summary

Task            | Status  | Duration    | Error
check_system    | Success | 60.187084ms |      
hostname_remoto | Success | 84.870375ms |      
list_home_files | Success | 39.32125ms  |
```

## Estrutura do Comando

```lua
local sloth_bin = "./sloth-runner"
local cmd = sloth_bin .. " agent run <AGENT_NAME> \"<SHELL_COMMAND>\" --master <MASTER_ADDRESS>"

local output, err, failed = exec.run(cmd)

if not failed then
    -- Sucesso! Use 'output'
    log.info("✅ " .. output)
    return true, "OK"
else
    -- Erro! Use 'err'
    log.error("❌ " .. (err or "erro"))
    return false, "Erro"
end
```

## Componentes

1. **AGENT_NAME**: Nome do agente registrado (ex: `ladyguica`, `keiteguica`)
2. **SHELL_COMMAND**: Comando shell a ser executado (deve estar entre aspas duplas)
3. **MASTER_ADDRESS**: Endereço do master (ex: `192.168.1.29:50053`)

## Dicas Importantes

### ✅ Funciona
```lua
-- Comando simples
local cmd = "./sloth-runner agent run ladyguica \"hostname\" --master 192.168.1.29:50053"

-- Múltiplos comandos com &&
local cmd = "./sloth-runner agent run ladyguica \"hostname && whoami && date\" --master 192.168.1.29:50053"

-- Usando variáveis de ambiente
local cmd = "./sloth-runner agent run keiteguica \"ls -la $HOME | head -10\" --master 192.168.1.29:50053"

-- Pipe e redirecionamento
local cmd = "./sloth-runner agent run ladyguica \"cat /etc/os-release | grep PRETTY\" --master 192.168.1.29:50053"
```

### ❌ Problemas Conhecidos

```lua
-- ❌ NÃO USE: Aspas simples para o comando (use aspas duplas)
local cmd = "./sloth-runner agent run ladyguica 'hostname' --master 192.168.1.29:50053"

-- ❌ NÃO USE: Scripts multi-linha diretamente
local cmd = "./sloth-runner agent run ladyguica echo 'linha1' echo 'linha2' --master 192.168.1.29:50053"
```

### 💡 Workarounds

Para scripts multi-linha, use uma das seguintes estratégias:

**Opção 1: Separar comandos com &&**
```lua
local cmd = "./sloth-runner agent run ladyguica \"echo 'linha1' && echo 'linha2'\" --master 192.168.1.29:50053"
```

**Opção 2: Criar script temporário no agente**
```lua
-- Criar script
local create = "./sloth-runner agent run ladyguica \"echo 'echo linha1' > /tmp/script.sh && echo 'echo linha2' >> /tmp/script.sh\" --master 192.168.1.29:50053"
exec.run(create)

-- Executar
local run = "./sloth-runner agent run ladyguica \"bash /tmp/script.sh\" --master 192.168.1.29:50053"
local output, err, failed = exec.run(run)

-- Limpar
exec.run("./sloth-runner agent run ladyguica \"rm /tmp/script.sh\" --master 192.168.1.29:50053")
```

## Verificar Agentes Disponíveis

```bash
./sloth-runner agent list --master 192.168.1.29:50053
```

Saída:
```
AGENT NAME     ADDRESS              STATUS            LAST HEARTBEAT
------------   ----------           ------            --------------
keiteguica     192.168.1.17:50051   Active   2025-10-01T12:27:35-03:00
ladyguica      192.168.1.16:50051   Active   2025-10-01T12:27:35-03:00
```

## Pré-requisitos

1. **Master rodando**
   ```bash
   ./sloth-runner master start --port 50053
   ```

2. **Agentes registrados**
   ```bash
   # No agente 1
   ./sloth-runner agent start --name ladyguica --master 192.168.1.29:50053 --daemon
   
   # No agente 2
   ./sloth-runner agent start --name keiteguica --master 192.168.1.29:50053 --daemon
   ```

## Vantagens deste Método

✅ **Simples**: Usa comando que já funciona  
✅ **Confiável**: Sem problemas de parsing Lua  
✅ **Flexível**: Aceita qualquer comando shell  
✅ **Direto**: Sem overhead de serialização  
✅ **Depurável**: Fácil testar comando fora do workflow  

## Exemplos Adicionais

### Coletar Informações do Sistema
```lua
{
    name = "system_info",
    command = function()
        local cmd = "./sloth-runner agent run ladyguica \"uname -a\" --master 192.168.1.29:50053"
        local output, err, failed = exec.run(cmd)
        if not failed then
            log.info("Sistema: " .. output)
            return true, "OK"
        end
        return false, "Erro"
    end,
    timeout = "60s"
}
```

### Verificar Espaço em Disco
```lua
{
    name = "check_disk",
    command = function()
        local cmd = "./sloth-runner agent run keiteguica \"df -h | grep -E '^/dev'\" --master 192.168.1.29:50053"
        local output, err, failed = exec.run(cmd)
        if not failed then
            log.info("Discos:\n" .. output)
            return true, "OK"
        end
        return false, "Erro"
    end,
    timeout = "60s"
}
```

### Executar em Múltiplos Agentes
```lua
{
    name = "check_all_agents",
    command = function()
        local agents = {"ladyguica", "keiteguica"}
        local master = "192.168.1.29:50053"
        
        for _, agent in ipairs(agents) do
            log.info("Verificando: " .. agent)
            local cmd = "./sloth-runner agent run " .. agent .. " \"hostname\" --master " .. master
            local output, err, failed = exec.run(cmd)
            if not failed then
                log.info(agent .. ": " .. output)
            else
                log.error(agent .. ": ERRO")
            end
        end
        
        return true, "Verificação completa"
    end,
    timeout = "120s"
}
```

## Referências

- [Documentação do Agent Run](/docs/agent_run.md)
- [Exemplos Funcionais](./functional_cmd_example.sloth)
- [Setup de Agentes](/docs/agent_setup.md)

---

**Status**: ✅ FUNCIONAL  
**Testado**: 2025-10-01  
**Agentes**: ladyguica, keiteguica  
**Master**: 192.168.1.29:50053
