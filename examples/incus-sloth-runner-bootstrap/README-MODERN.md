# 🚀 Bootstrap Moderno de Agent com Fluent API

## 📖 Overview

Este é um exemplo **moderno** que demonstra as últimas features do sloth-runner:

- ✅ **Fluent API** - Encadeamento elegante de métodos
- ✅ **Padrão (result, error)** - Tratamento consistente de erros
- ✅ **Idempotência** - Seguro para re-executar
- ✅ **Delegação Remota** - Execução em agents distribuídos
- ✅ **Error Handling Robusto** - Verificação em cada passo

## 🆚 Diferenças do Exemplo Original

### Exemplo Original (`bootstrap-agent.sloth`)
```lua
-- Padrão antigo (sem fluent API)
local container = incus.instance(container_name)
container:image("images:archlinux")
container:create()
container:start()

-- Sem verificação de erros
incus.exec({instance = name, command = "..."})
```

### Exemplo Moderno (`create-agent-vm-modern.sloth`)
```lua
-- ✨ Fluent API com encadeamento
local result, err = incus.instance(config.vm_name)
    :image(config.image)
    :config({
        ["limits.memory"] = "2GB",
        ["limits.cpu"] = "2"
    })
    :profile("default")
    :launch()  -- Cria e inicia em uma operação

-- ✅ Sempre verificar (result, error)
if err then
    log.error("❌ Erro: " .. err)
    return false, "Falha: " .. err
end

-- ✅ Verificação de sucesso
local exec_result, exec_err = incus.exec({
    instance = config.vm_name,
    command = "pacman -Sy"
})

if exec_err then
    return false, "Erro ao executar: " .. exec_err
end
```

## 🎯 Features Demonstradas

### 1. Fluent API do Incus

```lua
-- Encadear configurações elegantemente
incus.instance("my-vm")
    :image("images:archlinux")
    :config({
        ["limits.memory"] = "2GB",
        ["limits.cpu"] = "2"
    })
    :profile("default")
    :ephemeral(false)
    :launch()
```

### 2. Padrão (result, error) Consistente

Todos os módulos retornam `(result, error)`:

```lua
-- Módulo Incus
local result, err = incus.instance("vm"):create()
if err then
    return false, "Erro: " .. err
end

-- Módulo Sloth
local result, err = sloth.agent.install({...})
if err then
    return false, "Erro: " .. err
end

-- Módulo Exec
local result, err = exec.run("command")
if err then
    return false, "Erro: " .. err
end
```

### 3. Idempotência

O script detecta quando recursos já existem:

```lua
local result, err = sloth.agent.install({...})

if not err and result then
    if result.changed then
        log.info("✅ Agent instalado!")
    else
        log.info("ℹ️  Agent já estava instalado")
    end
end
```

### 4. Retry Automático

```lua
-- Tentar conectar SSH com retry
local max_retries = 5
for attempt = 1, max_retries do
    local result, err = exec.run(ssh_command)

    if not err and result.success then
        return true, "Conectado!"
    end

    if attempt < max_retries then
        exec.run("sleep 3")
    end
end
```

### 5. Delegação Remota

```lua
-- Criar VM no host remoto
:delegate_to(config.delegate_to)

-- Executar comandos na VM recém-criada
:delegate_to(config.vm_name)
```

## 🚀 Uso

### Configuração Rápida

Edite as variáveis no início do arquivo:

```lua
local config = {
    vm_name = "agent-keite-01",      -- Nome da VM
    image = "images:archlinux",       -- Imagem a usar
    delegate_to = "keite-guica",      -- Agent host
    master_addr = "192.168.1.29:50053",
    host_ip = "192.168.1.17",
    ssh_port = 50060,
    agent_port = 50060,
    memory = "2GB",
    cpus = 2
}
```

### Executar

```bash
# Executar bootstrap completo
sloth-runner run \
    --file examples/incus-sloth-runner-bootstrap/create-agent-vm-modern.sloth \
    --yes
```

### Output Esperado

```
🚀 Criando VM com fluent API: agent-keite-01
✅ VM criada e iniciada com sucesso!
🔐 Configurando SSH com tratamento de erros
📋 Chave SSH obtida (568 bytes)
✅ Sistema atualizado
✅ Pacotes instalados
✅ Chave SSH configurada
✅ SSH daemon configurado
🌐 Criando proxy devices...
📍 VM IP: 10.xxx.xxx.xxx
✅ Proxy SSH: 0.0.0.0:50060 -> 10.xxx.xxx.xxx:22
✅ Proxy Agent: 0.0.0.0:50060 -> 10.xxx.xxx.xxx:50060
🧪 Testando conexão SSH...
✅ SSH conectado! Hostname: agent-keite-01
🦥 Instalando sloth-runner agent...
✅ Agent instalado com sucesso!
🔍 Verificando agent no master...
✅ Agent agent-keite-01 encontrado e ativo!
🧪 Testando execução remota...
✅ Hostname: agent-keite-01
✅ Kernel: 6.x.x-arch1-1
✅ Uptime: up 2 minutes
✅ Memória: 2.0Gi total
✅ Disco: /dev/sda1 20G

═══════════════════════════════════════════════
🎉 BOOTSTRAP CONCLUÍDO COM SUCESSO!
═══════════════════════════════════════════════

📋 Configuração:
   VM: agent-keite-01
   Imagem: images:archlinux
   Memória: 2GB
   CPUs: 2

🌐 Acesso:
   SSH: ssh -p 50060 root@192.168.1.17
   Agent: agent-keite-01 @ 192.168.1.29:50053

🚀 Uso:
   sloth-runner run <task> --delegate-to agent-keite-01

✨ Features demonstradas:
   ✅ Fluent API para Incus
   ✅ Padrão (result, error) consistente
   ✅ Tratamento de erros robusto
   ✅ Operações idempotentes
   ✅ Delegação remota
```

## 📚 Estrutura do Workflow

### Tarefas

1. **create-vm-fluent** - Criar VM usando fluent API
2. **setup-ssh-modern** - Configurar SSH com error checking
3. **create-proxy-modern** - Criar proxy devices
4. **test-ssh-retry** - Testar SSH com retry automático
5. **install-agent-modern** - Instalar agent com novo padrão
6. **verify-agent-modern** - Verificar agent no master
7. **test-remote-execution** - Testar execução remota

### Fluxo de Execução

```
┌─────────────────────┐
│   Create VM         │ ← Fluent API
│   (delegate_to)     │
└──────────┬──────────┘
           ↓
┌─────────────────────┐
│   Setup SSH         │ ← Error Checking
│   (delegate_to)     │
└──────────┬──────────┘
           ↓
┌─────────────────────┐
│   Create Proxies    │ ← Idempotency
│   (delegate_to)     │
└──────────┬──────────┘
           ↓
┌─────────────────────┐
│   Test SSH          │ ← Retry Logic
│   (delegate_to)     │
└──────────┬──────────┘
           ↓
┌─────────────────────┐
│   Install Agent     │ ← (result, error)
│   (local)           │
└──────────┬──────────┘
           ↓
┌─────────────────────┐
│   Verify Agent      │ ← Master Check
│   (local)           │
└──────────┬──────────┘
           ↓
┌─────────────────────┐
│   Test Execution    │ ← Remote Test
│   (vm agent)        │
└─────────────────────┘
```

## 🔧 Personalização

### Usar Ubuntu em vez de Arch

```lua
local config = {
    image = "images:ubuntu/22.04",
    ...
}

-- Ajustar comandos SSH
local install_cmd = "apt-get update && apt-get install -y openssh-server"
```

### Criar Múltiplas VMs

```lua
local agents = {
    {name = "agent-01", ssh_port = 50061, agent_port = 50061},
    {name = "agent-02", ssh_port = 50062, agent_port = 50062},
    {name = "agent-03", ssh_port = 50063, agent_port = 50063}
}

-- Loop para criar cada agent
for _, agent_config in ipairs(agents) do
    -- Criar task com config específica
end
```

### VM com Mais Recursos

```lua
local config = {
    memory = "4GB",
    cpus = 4,
    ...
}

-- Adicionar disco extra
:config({
    ["limits.memory"] = "4GB",
    ["limits.cpu"] = "4",
    ["devices.root.size"] = "50GB"
})
```

## 🐛 Troubleshooting

### Erro: "Fluent API não funciona"

Certifique-se que você está usando a versão atualizada dos módulos:

```bash
git pull origin master
go build ./cmd/sloth-runner
```

### Erro: "result is nil"

Sempre verifique ambos `result` e `err`:

```lua
local result, err = module.function(...)

if err then
    -- Tratar erro
    return false, err
end

if not result then
    -- Result também pode ser nil
    return false, "Resultado inesperado"
end
```

### SSH não conecta

```bash
# Verificar proxy devices
incus config device show agent-keite-01

# Testar SSH manualmente
ssh -v -p 50060 root@192.168.1.17

# Verificar logs sshd no container
incus exec agent-keite-01 -- journalctl -u sshd -f
```

## 📖 Documentação Relacionada

- [Padrão de Retorno de Módulos](../../docs/MODULE_RETURN_PATTERN.md)
- [Helper Functions](../../internal/modules/helpers.go)
- [Módulo Incus](../../internal/modules/infra/incus.go)
- [Módulo Sloth](../../internal/modules/core/sloth.go)

## 🎓 Aprendizado

Este exemplo ensina:

1. ✅ **Como usar Fluent API** - Encadeamento de métodos
2. ✅ **Padrão (result, error)** - Sempre verificar ambos
3. ✅ **Idempotência** - Checar `result.changed`
4. ✅ **Error Handling** - Propagar erros corretamente
5. ✅ **Retry Logic** - Implementar tentativas automáticas
6. ✅ **Delegação** - Distribuir tarefas entre agents

## 🌟 Próximos Passos

Após executar este exemplo com sucesso:

1. **Adaptar para seu caso de uso**
   - Mudar distro Linux
   - Adicionar configurações personalizadas
   - Instalar software adicional

2. **Integrar com CI/CD**
   - Criar VMs de teste automaticamente
   - Destruir após testes
   - Pipeline completo

3. **Escalar**
   - Criar múltiplos agents
   - Load balancing
   - Alta disponibilidade

4. **Monitorar**
   - Coletar métricas
   - Alertas
   - Dashboards

## 🙌 Contribuindo

Encontrou um bug ou tem uma sugestão? Abra uma issue ou PR!

---

**Desenvolvido com ❤️ usando sloth-runner fluent API** 🦥
