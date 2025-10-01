# 🤖 Sloth Runner Agent Setup Guide

Este guia explica como configurar e iniciar agentes do Sloth Runner para apontar para o master server.

## 📋 Pré-requisitos

1. **Sloth Runner instalado** em todas as máquinas (master e agents)
2. **Conectividade de rede** entre master e agents
3. **Portas abertas**:
   - Master: porta 50053 (gRPC)
   - Agents: porta 50051 (gRPC) ou outra de sua escolha

## 🏗️ Arquitetura

```
Master Server (192.168.1.29:50053)
         ↓
    ┌────┴────┐
    ↓         ↓
Agent 1    Agent 2
(50051)    (50051)
```

## 🚀 Iniciando o Master Server

Primeiro, você precisa iniciar o master server na máquina central:

```bash
# Inicie o master server em modo daemon
./sloth-runner master start --port 50053 --bind-address 192.168.1.29 --daemon

# Ou sem daemon (para debug)
./sloth-runner master start --port 50053 --bind-address 192.168.1.29
```

**Flags importantes:**
- `--port`: Porta que o master vai escutar (padrão: 50053)
- `--bind-address`: IP da máquina master
- `--daemon`: Executa em background

## 🤖 Iniciando um Agent

### Opção 1: Agent Local (mesma máquina do master)

```bash
# Agent local para teste
./sloth-runner agent start \
  --name "local-agent" \
  --port 50051 \
  --master "192.168.1.29:50053" \
  --bind-address "192.168.1.29" \
  --daemon
```

### Opção 2: Agent Remoto

**Na máquina remota**, execute:

```bash
# Exemplo: Agent na máquina ladyguica (192.168.1.16)
./sloth-runner agent start \
  --name "ladyguica" \
  --port 50051 \
  --master "192.168.1.29:50053" \
  --bind-address "192.168.1.16" \
  --daemon
```

**Flags importantes:**
- `--name`: Nome único para o agent (obrigatório)
- `--port`: Porta que o agent vai escutar (padrão: 50051)
- `--master`: Endereço do master server (IP:PORTA)
- `--bind-address`: IP da máquina agent
- `--daemon`: Executa em background

## 📡 Configuração para Esta Máquina (192.168.1.29)

### 1. Iniciar o Master Server

```bash
cd ~/.local/bin
./sloth-runner master start --port 50053 --bind-address 192.168.1.29 --daemon
```

### 2. Configurar Agents Remotos

#### Agent ladyguica (192.168.1.16)

```bash
# Conecte via SSH
ssh usuario@192.168.1.16

# Na máquina remota
cd ~/.local/bin
./sloth-runner agent start \
  --name "ladyguica" \
  --port 50051 \
  --master "192.168.1.29:50053" \
  --bind-address "192.168.1.16" \
  --daemon
```

#### Agent keiteguica (192.168.1.17)

```bash
# Conecte via SSH
ssh usuario@192.168.1.17

# Na máquina remota
cd ~/.local/bin
./sloth-runner agent start \
  --name "keiteguica" \
  --port 50051 \
  --master "192.168.1.29:50053" \
  --bind-address "192.168.1.17" \
  --daemon
```

### 3. Verificar Agents Registrados

```bash
./sloth-runner agent list --master 192.168.1.29:50053
```

## 🔧 Comandos Úteis

### Listar Agents Registrados

```bash
./sloth-runner agent list --master 192.168.1.29:50053
```

### Executar Comando em um Agent

```bash
./sloth-runner agent run \
  --master "192.168.1.29:50053" \
  --agent "ladyguica" \
  --command "hostname"
```

### Parar um Agent

```bash
./sloth-runner agent stop \
  --master "192.168.1.29:50053" \
  --agent "ladyguica"
```

## 📝 Usando Agents em Scripts Lua

### Exemplo com :delegate_to()

```lua
-- Tarefa que executa em um agent específico
local deploy_task = task("deploy_app")
    :description("Deploy application on remote server")
    :delegate_to("ladyguica")  -- Executa no agent 'ladyguica'
    :command(function(this, params)
        local exec = require("exec")
        log.info("🚀 Deploying on " .. this.agent.get())
        
        -- Este comando será executado no agent 'ladyguica'
        local result = exec.run("systemctl restart nginx")
        
        if result.exit_code == 0 then
            return true, "Deployment successful"
        else
            return false, "Deployment failed: " .. result.stderr
        end
    end)
    :build()

-- Executar a tarefa
deploy_task:run()
```

### Exemplo com Múltiplos Agents

```lua
local check_task = task("check_status")
    :description("Check service status on all agents")
    :command(function(this, params)
        local agents = {"ladyguica", "keiteguica"}
        
        for _, agent_name in ipairs(agents) do
            log.info("Checking " .. agent_name)
            
            -- Criar subtarefa para cada agent
            local check = task("check_" .. agent_name)
                :delegate_to(agent_name)
                :command(function(t, p)
                    local exec = require("exec")
                    return exec.run("systemctl status nginx")
                end)
                :build()
            
            check:run()
        end
        
        return true
    end)
    :build()
```

## 🐛 Troubleshooting

### Agent não conecta ao Master

1. **Verificar se o master está rodando:**
   ```bash
   ps aux | grep sloth-runner
   netstat -an | grep 50053
   ```

2. **Verificar firewall:**
   ```bash
   # No master, liberar porta 50053
   sudo ufw allow 50053/tcp
   
   # No agent, liberar porta 50051
   sudo ufw allow 50051/tcp
   ```

3. **Verificar conectividade:**
   ```bash
   # Do agent, testar conexão com o master
   telnet 192.168.1.29 50053
   # ou
   nc -zv 192.168.1.29 50053
   ```

### Ver logs do Agent

```bash
# Se rodando em daemon, verificar logs
tail -f ~/.local/var/log/sloth-runner-agent.log

# Ou rodar sem daemon para ver logs em tempo real
./sloth-runner agent start \
  --name "test-agent" \
  --port 50051 \
  --master "192.168.1.29:50053" \
  --bind-address "192.168.1.29"
```

### Remover Agent desconectado

```bash
# Parar o agent
./sloth-runner agent stop --master "192.168.1.29:50053" --agent "nome-do-agent"

# Ou manualmente matar o processo
ps aux | grep sloth-runner
kill <PID>
```

## 🔐 Segurança

### Recomendações

1. **Use firewall** para restringir acesso às portas gRPC
2. **Configure SSH keys** para acesso remoto seguro
3. **Execute agents com usuário dedicado** (não root)
4. **Monitore logs** regularmente

### Exemplo de configuração de firewall

```bash
# No master (192.168.1.29)
sudo ufw allow from 192.168.1.0/24 to any port 50053 proto tcp

# Nos agents
sudo ufw allow from 192.168.1.29 to any port 50051 proto tcp
```

## 📊 Monitoramento

### Verificar status dos agents periodicamente

```bash
# Criar um script de monitoramento
watch -n 5 './sloth-runner agent list --master 192.168.1.29:50053'
```

### Health Check

```lua
-- Script Lua para health check
local check_agents = task("health_check")
    :description("Check all agents health")
    :command(function(this, params)
        local agents = {"ladyguica", "keiteguica"}
        local healthy = {}
        local unhealthy = {}
        
        for _, agent_name in ipairs(agents) do
            local health = task("ping_" .. agent_name)
                :delegate_to(agent_name)
                :command(function(t, p)
                    return true, "Agent is healthy"
                end)
                :timeout("5s")
                :build()
            
            local success = health:run()
            if success then
                table.insert(healthy, agent_name)
            else
                table.insert(unhealthy, agent_name)
            end
        end
        
        log.info("✅ Healthy agents: " .. table.concat(healthy, ", "))
        if #unhealthy > 0 then
            log.warn("⚠️ Unhealthy agents: " .. table.concat(unhealthy, ", "))
        end
        
        return true
    end)
    :build()
```

## 🎯 Próximos Passos

1. ✅ Instale o sloth-runner em todas as máquinas
2. ✅ Inicie o master server
3. ✅ Configure e inicie os agents
4. ✅ Verifique a conexão com `agent list`
5. ✅ Teste com um script Lua simples usando `:delegate_to()`

---

Para mais informações, consulte:
- [README principal](../README.md)
- [Documentação de módulos](./modules/)
- [Exemplos](../examples/)
