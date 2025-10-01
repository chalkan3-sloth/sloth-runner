# 📦 Instalação do Sloth Runner

## ✅ Instalação Concluída

O executável `sloth-runner` foi instalado em:
```
$HOME/.local/bin/sloth-runner
```

### Verificação

O binário está disponível globalmente:
```bash
$ which sloth-runner
/Users/chalkan3/.local/bin/sloth-runner

$ sloth-runner agent list --master 192.168.1.29:50053
AGENT NAME     ADDRESS              STATUS   LAST HEARTBEAT
------------   ----------           ------   --------------
keiteguica     192.168.1.17:50051   Active   2025-10-01T12:41:13-03:00
ladyguica      192.168.1.16:50051   Active   2025-10-01T12:41:13-03:00
```

### PATH

O diretório `$HOME/.local/bin` já está no seu PATH, então você pode usar `sloth-runner` de qualquer lugar.

## 🚀 Como Usar

### Comandos Básicos

```bash
# Listar agentes
sloth-runner agent list --master 192.168.1.29:50053

# Executar comando em agente
sloth-runner agent run ladyguica "hostname" --master 192.168.1.29:50053

# Executar workflow
sloth-runner run -f examples/agents/hello_remote_cmd.sloth hello_remote

# Iniciar master
sloth-runner master start --port 50053 --daemon

# Iniciar agente
sloth-runner agent start --name myagent --master 192.168.1.29:50053 --daemon
```

### Exemplos Prontos

Os exemplos foram atualizados para usar o caminho completo do binário:

```bash
# De qualquer diretório
cd ~
sloth-runner run -f ~/.projects/task-runner/examples/agents/hello_remote_cmd.sloth hello_remote

# Ou do diretório do projeto
cd ~/.projects/task-runner
sloth-runner run -f examples/agents/functional_cmd_example.sloth remote_via_cmd
```

## 📝 Atualização dos Exemplos

Os exemplos foram atualizados para usar o caminho completo do binário:

```lua
-- Antes
local cmd = "./sloth-runner agent run ladyguica \"hostname\" --master 192.168.1.29:50053"

-- Depois
local cmd = "/Users/chalkan3/.local/bin/sloth-runner agent run ladyguica \"hostname\" --master 192.168.1.29:50053"
```

Isso garante que os exemplos funcionem mesmo quando executados de qualquer diretório.

## 🔧 Para Usar em Seus Scripts

Você pode usar tanto o caminho relativo quanto o absoluto:

```lua
-- Opção 1: Caminho completo (mais confiável)
local cmd = "/Users/chalkan3/.local/bin/sloth-runner agent run <agent> \"<comando>\" --master <master>"

-- Opção 2: Apenas sloth-runner (se PATH estiver configurado)
-- Nota: Pode não funcionar em alguns contextos Lua
local cmd = "sloth-runner agent run <agent> \"<comando>\" --master <master>"
```

**Recomendação**: Use o caminho completo nos scripts Lua para garantir compatibilidade.

## 📂 Localização dos Arquivos

```
$HOME/.local/bin/sloth-runner              → Executável instalado
$HOME/.projects/task-runner/sloth-runner   → Executável original
$HOME/.projects/task-runner/examples/      → Exemplos atualizados
```

## ✅ Teste de Verificação

Execute este comando para verificar a instalação:

```bash
sloth-runner agent list --master 192.168.1.29:50053
```

Se ver a lista de agentes, está tudo funcionando! ✅

## 🔄 Atualização

Para atualizar o binário:

```bash
cd ~/.projects/task-runner
# Após compilar/baixar nova versão
cp sloth-runner $HOME/.local/bin/
```

---

**Data de Instalação**: 2025-10-01  
**Localização**: $HOME/.local/bin/sloth-runner  
**Status**: ✅ INSTALADO E FUNCIONANDO
