# 🎯 Exemplos de Execução Remota - Sloth Runner

Este diretório contém exemplos de execução de tarefas em agentes remotos usando o Sloth Runner.

## 🚀 Início Rápido

### Método Recomendado: ✅ agent run via CMD

**Este é o método que FUNCIONA!**

```bash
# 1. Execute o exemplo mínimo (Hello World)
./sloth-runner run -f examples/agents/hello_remote_cmd.sloth hello_remote

# 2. Execute o exemplo funcional completo
./sloth-runner run -f examples/agents/functional_cmd_example.sloth remote_via_cmd

# 3. Execute o pipeline de infraestrutura
./sloth-runner run -f examples/agents/complete_infrastructure_check.sloth distributed_infrastructure_check
```

## 📚 Documentação

### Para Iniciantes
- **[QUICK_START.md](./QUICK_START.md)** - Guia rápido para começar
  - Templates prontos para usar
  - Exemplos simples
  - Como executar

### Documentação Completa
- **[README_CMD_FUNCIONAL.md](./README_CMD_FUNCIONAL.md)** - Guia completo
  - Como funciona
  - Dicas e workarounds
  - Exemplos avançados
  
- **[DELEGATE_TO_SOLUTION.md](./DELEGATE_TO_SOLUTION.md)** - Descrição da solução
  - Problema original
  - Solução implementada
  - Comparação de métodos

## 📝 Exemplos Disponíveis

### 1. Hello World Remoto (Mínimo)
**Arquivo**: `hello_remote_cmd.sloth`

O exemplo mais simples possível que demonstra execução remota.

```bash
./sloth-runner run -f examples/agents/hello_remote_cmd.sloth hello_remote
```

**Saída**:
```
✅ Comando executado!
📋 Saída: Hello from remote agent!

# Execution Summary
hello | Success | 54ms |
```

### 2. Exemplo Funcional (Intermediário)
**Arquivo**: `functional_cmd_example.sloth`

Demonstra múltiplos comandos em diferentes agentes.

```bash
./sloth-runner run -f examples/agents/functional_cmd_example.sloth remote_via_cmd
```

**Executa**:
- ✅ Hostname remoto
- ✅ Listagem de arquivos
- ✅ Informações do sistema

### 3. Pipeline de Infraestrutura (Avançado)
**Arquivo**: `complete_infrastructure_check.sloth`

Pipeline completo de verificação de infraestrutura.

```bash
./sloth-runner run -f examples/agents/complete_infrastructure_check.sloth distributed_infrastructure_check
```

**Executa**:
- ✅ Verificação de conectividade
- ✅ Coleta de informações do sistema
- ✅ Verificação de recursos (CPU, memória, disco)
- ✅ Verificação de serviços
- ✅ Teste de performance básico
- ✅ Geração de relatório final

## 🔧 Pré-requisitos

### 1. Master Rodando
```bash
./sloth-runner master start --port 50053 --daemon
```

### 2. Agentes Registrados
```bash
# No agente 1
./sloth-runner agent start --name ladyguica --master 192.168.1.29:50053 --daemon

# No agente 2
./sloth-runner agent start --name keiteguica --master 192.168.1.29:50053 --daemon
```

### 3. Verificar Agentes
```bash
./sloth-runner agent list --master 192.168.1.29:50053
```

Saída esperada:
```
AGENT NAME     ADDRESS              STATUS   LAST HEARTBEAT
------------   ----------           ------   --------------
keiteguica     192.168.1.17:50051   Active   2025-10-01T12:27:35-03:00
ladyguica      192.168.1.16:50051   Active   2025-10-01T12:27:35-03:00
```

## 💡 Template Básico

Copie e adapte este template para seus casos de uso:

```lua
TaskDefinitions = {
    minha_task = {
        description = "Descrição da minha task",
        tasks = {
            {
                name = "nome_task",
                description = "O que esta task faz",
                command = function()
                    log.info("🚀 Executando comando remoto...")
                    
                    local cmd = "./sloth-runner agent run <AGENT_NAME> \"<COMANDO>\" --master <MASTER_ADDR>"
                    local output, err, failed = exec.run(cmd)
                    
                    if not failed then
                        log.info("✅ Sucesso!")
                        log.info("📋 Saída: " .. output)
                        return true, "Sucesso"
                    else
                        log.error("❌ Erro: " .. (err or "erro"))
                        return false, "Erro"
                    end
                end,
                timeout = "60s"
            }
        }
    }
}
```

## 📖 Exemplos de Comandos

### Comando Simples
```lua
local cmd = "./sloth-runner agent run ladyguica \"hostname\" --master 192.168.1.29:50053"
```

### Múltiplos Comandos
```lua
local cmd = "./sloth-runner agent run ladyguica \"hostname && whoami && date\" --master 192.168.1.29:50053"
```

### Com Pipes
```lua
local cmd = "./sloth-runner agent run keiteguica \"ls -la $HOME | head -10\" --master 192.168.1.29:50053"
```

### Loop em Múltiplos Agentes
```lua
local agents = {"ladyguica", "keiteguica"}
local master = "192.168.1.29:50053"

for _, agent in ipairs(agents) do
    local cmd = "./sloth-runner agent run " .. agent .. " \"hostname\" --master " .. master
    local output, err, failed = exec.run(cmd)
    if not failed then
        log.info(agent .. ": " .. output:gsub("\n", ""))
    end
end
```

## 🎯 Status dos Exemplos

| Arquivo | Status | Descrição |
|---------|--------|-----------|
| `hello_remote_cmd.sloth` | ✅ FUNCIONA | Hello World remoto |
| `functional_cmd_example.sloth` | ✅ FUNCIONA | Exemplo funcional completo |
| `complete_infrastructure_check.sloth` | ✅ FUNCIONA | Pipeline de infraestrutura |
| `delegate_to_working.sloth` | ❌ NÃO FUNCIONA | Método antigo (não use) |
| `simple_delegate_ls.sloth` | ❌ NÃO FUNCIONA | Método antigo (não use) |

## 🐛 Troubleshooting

### Erro: "failed to connect to master"
```bash
# Verificar se o master está rodando
./sloth-runner agent list --master 192.168.1.29:50053
```

### Erro: "failed to resolve agent"
```bash
# Listar agentes disponíveis
./sloth-runner agent list --master 192.168.1.29:50053
```

### Comando falha no agente
```bash
# Testar comando diretamente
./sloth-runner agent run <agent_name> "hostname" --master 192.168.1.29:50053
```

## 🎓 Próximos Passos

1. **Leia o guia rápido**: [QUICK_START.md](./QUICK_START.md)
2. **Execute o exemplo mínimo**: `hello_remote_cmd.sloth`
3. **Teste com seus comandos**: Adapte o template básico
4. **Explore o pipeline completo**: `complete_infrastructure_check.sloth`
5. **Leia a documentação completa**: [README_CMD_FUNCIONAL.md](./README_CMD_FUNCIONAL.md)

## 🔗 Links Úteis

- [Guia Rápido](./QUICK_START.md) - Como começar
- [README Completo](./README_CMD_FUNCIONAL.md) - Documentação detalhada
- [Solução delegate_to](./DELEGATE_TO_SOLUTION.md) - Detalhes técnicos
- [Exemplo Hello World](./hello_remote_cmd.sloth) - Código mínimo
- [Exemplo Funcional](./functional_cmd_example.sloth) - Código intermediário
- [Exemplo Completo](./complete_infrastructure_check.sloth) - Pipeline avançado

## ⚡ Comandos Úteis

```bash
# Listar agentes
./sloth-runner agent list --master 192.168.1.29:50053

# Testar comando remoto
./sloth-runner agent run ladyguica "hostname" --master 192.168.1.29:50053

# Executar workflow
./sloth-runner run -f <arquivo.sloth> <task_group>

# Iniciar master
./sloth-runner master start --port 50053 --daemon

# Iniciar agente
./sloth-runner agent start --name <nome> --master <master_addr> --daemon

# Parar agente
./sloth-runner agent stop <nome> --master <master_addr>
```

---

**Status**: ✅ SOLUÇÃO IMPLEMENTADA E TESTADA  
**Última atualização**: 2025-10-01  
**Versão**: 1.0
