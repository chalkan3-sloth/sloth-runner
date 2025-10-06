# 🦥 Incus Sloth Runner Bootstrap

Bootstrap automatizado de agent sloth-runner em container Incus.

## 📋 Overview

Este exemplo demonstra como:
- ✅ Criar container Arch Linux usando módulo `incus`
- ✅ Configurar SSH com chaves públicas
- ✅ Criar proxy device para acesso SSH externo
- ✅ Instalar sloth-runner agent usando módulo `sloth`
- ✅ Verificar agent conectado ao master

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────────────────────┐
│                    Master Server                             │
│              (192.168.1.29:50053)                           │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ gRPC
                         │
┌────────────────────────┴────────────────────────────────────┐
│              Host: keite-guica (192.168.1.17)               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │           Incus Container: agent-keite-01            │   │
│  │                                                      │   │
│  │  - OS: Arch Linux                                    │   │
│  │  - SSH: Port 22 (interno)                           │   │
│  │  - Agent: Port 50060                                │   │
│  │  - User: root (com SSH key)                         │   │
│  └──────────────────────────────────────────────────────┘   │
│         │                                  │                 │
│         │ Proxy Device                     │ Proxy Device    │
│         │ 0.0.0.0:50060 -> 22             │ (agent port)     │
│         ▼                                  ▼                 │
│    SSH: 192.168.1.17:50060           Agent: 50060           │
└─────────────────────────────────────────────────────────────┘
```

## 📦 Pré-requisitos

1. **Incus instalado** no host `keite-guica`
   ```bash
   # Verificar se incus está instalado
   incus version
   ```

2. **Imagem Arch Linux** disponível
   ```bash
   # Verificar imagens disponíveis
   incus image list images: | grep archlinux
   ```

3. **Chave SSH pública** em `~/.ssh/id_rsa.pub` no agent keite-guica
   ```bash
   # Gerar chave se não existir
   ssh-keygen -t rsa -b 4096 -N "" -f ~/.ssh/id_rsa
   ```

4. **Agent keite-guica** conectado ao master
   ```bash
   # Verificar agents
   sloth-runner agent list --master 192.168.1.29:50053
   ```

## 🚀 Uso

### Executar Bootstrap Completo

```bash
sloth-runner run --file examples/incus-sloth-runner-bootstrap/bootstrap-agent.sloth --yes
```

### Workflow Passo a Passo

O workflow executa as seguintes tarefas:

1. **create-container** - Cria container Arch Linux
   ```lua
   incus.instance({
       name = "agent-keite-01",
       image = "images:archlinux"
   }):create():start():wait_running()
   ```

2. **configure-ssh** - Configura SSH no container
   - Instala OpenSSH e sudo
   - Adiciona chave pública SSH para root
   - Habilita PermitRootLogin e PubkeyAuthentication
   - Inicia sshd

3. **create-ssh-proxy** - Cria proxy device SSH
   - Mapeia porta 50060 do host para porta 22 do container
   - Permite acesso SSH externo: `ssh -p 50060 root@192.168.1.17`

4. **test-ssh** - Testa conexão SSH
   - Verifica conectividade SSH
   - Confirma que SSH está funcionando

5. **install-agent** - Instala sloth-runner agent
   ```lua
   sloth.agent.install({
       name = "agent-keite-01",
       ssh_host = "192.168.1.17",
       ssh_port = "50060",
       master = "192.168.1.29:50053"
   })
   ```

6. **verify-agent** - Verifica agent no master
   - Lista agents ativos
   - Confirma que novo agent está conectado

7. **test-agent** - Testa execução remota
   - Executa comandos no novo agent via delegate_to
   - Verifica funcionalidade completa

## 🎯 Resultado Esperado

Após execução bem-sucedida:

```
🎉 Bootstrap concluído com sucesso!

📋 Resumo:
  ✅ Container criado: agent-keite-01
  ✅ SSH configurado na porta: 50060
  ✅ Agent instalado e conectado
  ✅ Master: 192.168.1.29:50053

🔧 Para acessar o container via SSH:
   ssh -p 50060 root@192.168.1.17

🦥 Para usar o agent:
   sloth-runner run <task> --delegate-to agent-keite-01
```

## 📝 Usar o Novo Agent

### Executar task no agent

```bash
# Criar arquivo de teste
cat > test-agent.sloth << 'EOF'
local test = task("test")
    :description("Testar novo agent")
    :delegate_to("agent-keite-01")
    :command(function(this, params)
        local result = exec.run("uname -a")
        log.info("Sistema: " .. result.stdout)
        return true, "Teste OK"
    end)
    :build()

workflow.define("test"):tasks({test})
EOF

# Executar
sloth-runner run --file test-agent.sloth --yes
```

### Acessar via SSH

```bash
# SSH direto
ssh -p 50060 root@192.168.1.17

# Verificar sloth-runner instalado
ssh -p 50060 root@192.168.1.17 "sloth-runner version"
```

### Verificar agent status

```bash
# Listar agents
sloth-runner agent list --master 192.168.1.29:50053

# Detalhes do agent
sloth-runner agent get agent-keite-01 --master 192.168.1.29:50053
```

## 🔧 Personalização

### Variáveis Configuráveis

Edite as variáveis no início do arquivo `bootstrap-agent.sloth`:

```lua
local container_name = "agent-keite-01"  -- Nome do container
local host_agent = "keite-guica"         -- Agent onde criar container
local master_addr = "192.168.1.29:50053" -- Endereço do master
local ssh_port = 50060                    -- Porta SSH no host
local agent_port = 50060                  -- Porta do agent
```

### Criar Múltiplos Agents

Para criar vários agents, modifique o workflow ou crie um loop:

```lua
local agents = {"agent-01", "agent-02", "agent-03"}
local base_port = 50060

for i, name in ipairs(agents) do
    local port = base_port + i
    -- Criar agent com porta diferente
end
```

## 🛠️ Troubleshooting

### Container não inicia
```bash
# Verificar status
incus list agent-keite-01

# Ver logs
incus console agent-keite-01 --show-log
```

### SSH não conecta
```bash
# Verificar proxy device
incus config device show agent-keite-01

# Verificar porta em uso
netstat -tuln | grep 50060

# Testar SSH verbose
ssh -vvv -p 50060 root@192.168.1.17
```

### Agent não conecta ao master
```bash
# Verificar logs do agent no container
ssh -p 50060 root@192.168.1.17 "journalctl -u sloth-runner-agent-* -f"

# Verificar serviço
ssh -p 50060 root@192.168.1.17 "systemctl status sloth-runner-agent-*"
```

### Recriar do zero
```bash
# Parar e deletar container
incus stop agent-keite-01 --force
incus delete agent-keite-01

# Executar bootstrap novamente
sloth-runner run --file examples/incus-sloth-runner-bootstrap/bootstrap-agent.sloth --yes
```

## 🧹 Cleanup

### Remover agent e container

```bash
# 1. Remover agent do master
sloth-runner agent delete agent-keite-01 --master 192.168.1.29:50053 --yes

# 2. Parar e deletar container
incus stop agent-keite-01 --force
incus delete agent-keite-01
```

### Remover workflow salvo (se aplicável)

```bash
sloth-runner sloth remove incus-bootstrap --yes
```

## 📚 Módulos Utilizados

Este exemplo demonstra o uso dos seguintes módulos:

- **incus** - Gerenciamento de containers/VMs
  - `incus.instance()` - Criar e gerenciar instâncias
  - `.create()`, `.start()`, `.wait_running()` - Lifecycle
  - `.exec()` - Executar comandos no container

- **sloth** - Automação do sloth-runner
  - `sloth.agent.install()` - Instalar agent remotamente
  - `sloth.agent.list()` - Listar agents
  - Idempotente e com estado

- **exec** - Execução de comandos
  - `exec.run()` - Executar comandos shell
  - Captura stdout, stderr, exit_code

- **facts** - Informações do sistema
  - `facts.os()` - Detectar sistema operacional
  - Útil para verificações

## 🌟 Features Demonstradas

- ✅ **Fluent API** - Encadeamento de métodos
- ✅ **Delegação Remota** - Execução em agents
- ✅ **Integração de Módulos** - incus + sloth + exec
- ✅ **Automação Completa** - Zero intervenção manual
- ✅ **Idempotência** - Safe para re-executar
- ✅ **Error Handling** - Tratamento de erros
- ✅ **Logging Rico** - Feedback detalhado

## 📖 Veja Também

- [Módulo Incus](../../docs/modules/incus.md)
- [Módulo Sloth](../../internal/modules/core/sloth.go)
- [Modern DSL](../../docs/modern-dsl/introduction.md)
- [Agent Management](../../docs/en/master-agent-architecture.md)
