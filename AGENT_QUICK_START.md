# 🚀 Quick Agent Setup Reference

## 📋 TL;DR - Como Iniciar Agents

### 1️⃣ Inicie o Master Server (uma vez)
```bash
./sloth-runner master start --port 50053 --bind-address <SEU_IP> --daemon
```

### 2️⃣ Inicie um Agent Local (para teste)
```bash
./sloth-runner agent start \
  --name "meu-agent" \
  --port 50051 \
  --master "<IP_MASTER>:50053" \
  --bind-address "<SEU_IP>" \
  --daemon
```

### 3️⃣ Inicie um Agent Remoto via SSH
```bash
# Na máquina remota
ssh user@192.168.1.X
cd ~/.local/bin
./sloth-runner agent start \
  --name "agent-remoto" \
  --port 50051 \
  --master "<IP_MASTER>:50053" \
  --bind-address "<IP_AGENT>" \
  --daemon
```

### 4️⃣ Verifique os Agents
```bash
./sloth-runner agent list --master <IP_MASTER>:50053
```

### 5️⃣ Use em Scripts Lua
```lua
local task_remota = task("minha_tarefa")
    :delegate_to("meu-agent")  -- Nome do agent
    :command(function(this, params)
        -- Este código roda no agent remoto
        local exec = require("exec")
        return exec.run("hostname")
    end)
    :build()

task_remota:run()
```

## 🛠️ Scripts Auxiliares

### Iniciar Master
```bash
./start_master.sh
```

### Iniciar Agent Local
```bash
./start_local_agent.sh <nome> [ip] [porta]
```

### Gerenciar Agents Remotos
```bash
# Iniciar
./manage_remote_agent.sh start user@host nome-agent ip-agent

# Verificar status
./manage_remote_agent.sh status user@host nome-agent

# Parar
./manage_remote_agent.sh stop user@host nome-agent

# Instalar sloth-runner
./manage_remote_agent.sh install user@host
```

## 📚 Documentação Completa

- **[Agent Setup Guide](./docs/agent-setup.md)** - Guia completo de configuração
- **[Distributed Example](./examples/distributed_execution.sloth)** - Exemplo completo de uso
- **[README Principal](./README.md)** - Documentação geral do projeto

## 🔧 Troubleshooting Rápido

### Agent não conecta?
1. ✅ Master está rodando? `ps aux | grep "sloth-runner master"`
2. ✅ Firewall liberado? Portas 50053 (master) e 50051 (agent)
3. ✅ Conectividade de rede? `ping <IP_MASTER>`

### Ver logs do agent
```bash
# Se rodando em foreground
./sloth-runner agent start --name test (sem --daemon)

# Verificar processos
ps aux | grep sloth-runner
```

### Parar agents
```bash
# Parar agent específico
pkill -f "sloth-runner agent start.*--name <nome>"

# Parar master
pkill -f "sloth-runner master"
```

## 🎯 Exemplo Completo

```bash
# 1. Nesta máquina (master)
./start_master.sh

# 2. Em outra máquina ou terminal
./start_local_agent.sh test-agent

# 3. Verificar
./sloth-runner agent list --master 192.168.1.29:50053

# 4. Testar com script
cat > test.sloth << 'EOF'
local t = task("test")
    :delegate_to("test-agent")
    :command(function()
        local exec = require("exec")
        return exec.run("hostname && uptime")
    end)
    :build()
t:run()
EOF

./sloth-runner run -f test.sloth
```

---

Para mais informações, consulte a [documentação completa](./docs/agent-setup.md).
