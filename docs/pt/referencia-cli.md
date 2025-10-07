# 📚 Referência Completa de Comandos CLI

## Visão Geral

O Sloth Runner oferece uma interface de linha de comando (CLI) completa e poderosa para gerenciar workflows, agentes, módulos, hooks, eventos e muito mais. Esta documentação cobre **todos** os comandos disponíveis com exemplos práticos.

---

## 🎯 Comandos Principais

### `run` - Executar Workflows

Executa um workflow Sloth a partir de um arquivo.

```bash
# Sintaxe básica
sloth-runner run <workflow-name> --file <arquivo.sloth> [opções]

# Exemplos
sloth-runner run deploy --file deploy.sloth
sloth-runner run deploy --file deploy.sloth --yes                    # Modo não-interativo
sloth-runner run deploy --file deploy.sloth --group production       # Executa grupo específico
sloth-runner run deploy --file deploy.sloth --delegate-to agent1     # Delega para agente
sloth-runner run deploy --file deploy.sloth --delegate-to agent1 --delegate-to agent2  # Múltiplos agentes
sloth-runner run deploy --file deploy.sloth --values vars.yaml       # Passa variáveis
sloth-runner run deploy --file deploy.sloth --var "env=production"   # Variável inline
```

**Opções:**
- `--file, -f` - Caminho do arquivo Sloth
- `--yes, -y` - Modo não-interativo (não pede confirmação)
- `--group, -g` - Executa apenas um grupo específico
- `--delegate-to` - Delega execução para agente(s) remoto(s)
- `--values` - Arquivo YAML com variáveis
- `--var` - Define variável inline (pode usar múltiplas vezes)
- `--verbose, -v` - Modo verboso

---

## 🤖 Gerenciamento de Agentes

### `agent list` - Listar Agentes

Lista todos os agentes registrados no servidor master.

```bash
# Sintaxe
sloth-runner agent list [opções]

# Exemplos
sloth-runner agent list                    # Lista todos os agentes
sloth-runner agent list --format json      # Saída em JSON
sloth-runner agent list --format yaml      # Saída em YAML
sloth-runner agent list --status active    # Apenas agentes ativos
```

**Opções:**
- `--format` - Formato de saída: table (padrão), json, yaml
- `--status` - Filtrar por status: active, inactive, all

---

### `agent get` - Detalhes do Agente

Obtém informações detalhadas sobre um agente específico.

```bash
# Sintaxe
sloth-runner agent get <agent-name> [opções]

# Exemplos
sloth-runner agent get web-server-01
sloth-runner agent get web-server-01 --format json
sloth-runner agent get web-server-01 --show-metrics       # Inclui métricas
```

**Opções:**
- `--format` - Formato de saída: table, json, yaml
- `--show-metrics` - Mostra métricas do agente

---

### `agent install` - Instalar Agente Remoto

Instala o agente Sloth Runner em um servidor remoto via SSH.

```bash
# Sintaxe
sloth-runner agent install <agent-name> --ssh-host <host> --ssh-user <user> [opções]

# Exemplos
sloth-runner agent install web-01 --ssh-host 192.168.1.100 --ssh-user root
sloth-runner agent install web-01 --ssh-host 192.168.1.100 --ssh-user root --ssh-port 2222
sloth-runner agent install web-01 \
  --ssh-host 192.168.1.100 \
  --ssh-user root \
  --master 192.168.1.1:50053 \
  --bind-address 0.0.0.0 \
  --port 50060 \
  --report-address 192.168.1.100:50060
```

**Opções:**
- `--ssh-host` - Host SSH do servidor remoto (obrigatório)
- `--ssh-user` - Usuário SSH (obrigatório)
- `--ssh-port` - Porta SSH (padrão: 22)
- `--ssh-key` - Caminho da chave SSH privada
- `--master` - Endereço do servidor master (padrão: localhost:50053)
- `--bind-address` - Endereço de bind do agente (padrão: 0.0.0.0)
- `--port` - Porta do agente (padrão: 50060)
- `--report-address` - Endereço que o agente reporta ao master

---

### `agent update` - Atualizar Agente

Atualiza o binário do agente para a versão mais recente.

```bash
# Sintaxe
sloth-runner agent update <agent-name> [opções]

# Exemplos
sloth-runner agent update web-01
sloth-runner agent update web-01 --version v1.2.3
sloth-runner agent update web-01 --restart           # Reinicia após atualizar
```

**Opções:**
- `--version` - Versão específica (padrão: latest)
- `--restart` - Reinicia o agente após atualização
- `--force` - Força atualização mesmo se a versão for a mesma

---

### `agent modules` - Módulos do Agente

Lista ou verifica módulos disponíveis em um agente.

```bash
# Sintaxe
sloth-runner agent modules <agent-name> [opções]

# Exemplos
sloth-runner agent modules web-01                      # Lista todos os módulos
sloth-runner agent modules web-01 --check pkg          # Verifica se módulo 'pkg' está disponível
sloth-runner agent modules web-01 --check docker      # Verifica se Docker está instalado
sloth-runner agent modules web-01 --format json       # Saída em JSON
```

**Opções:**
- `--check` - Verifica módulo específico
- `--format` - Formato de saída: table, json, yaml

---

### `agent start` - Iniciar Agente

Inicia o serviço do agente localmente.

```bash
# Sintaxe
sloth-runner agent start [opções]

# Exemplos
sloth-runner agent start                                    # Inicia com configurações padrão
sloth-runner agent start --master 192.168.1.1:50053         # Conecta a master específico
sloth-runner agent start --port 50060                       # Usa porta específica
sloth-runner agent start --name my-agent                    # Define nome do agente
sloth-runner agent start --bind 0.0.0.0                     # Bind em todas as interfaces
sloth-runner agent start --foreground                       # Executa em primeiro plano
```

**Opções:**
- `--master` - Endereço do servidor master (padrão: localhost:50053)
- `--port` - Porta do agente (padrão: 50060)
- `--name` - Nome do agente (padrão: hostname)
- `--bind` - Endereço de bind (padrão: 0.0.0.0)
- `--report-address` - Endereço que o agente reporta
- `--foreground` - Executa em primeiro plano (não daemon)

---

### `agent stop` - Parar Agente

Para o serviço do agente.

```bash
# Sintaxe
sloth-runner agent stop [opções]

# Exemplos
sloth-runner agent stop                # Para agente local
sloth-runner agent stop --name web-01  # Para agente específico
```

---

### `agent restart` - Reiniciar Agente

Reinicia o serviço do agente.

```bash
# Sintaxe
sloth-runner agent restart [agent-name]

# Exemplos
sloth-runner agent restart               # Reinicia agente local
sloth-runner agent restart web-01        # Reinicia agente remoto
```

---

### `agent metrics` - Métricas do Agente

Visualiza métricas de performance e recursos do agente.

```bash
# Sintaxe
sloth-runner agent metrics <agent-name> [opções]

# Exemplos
sloth-runner agent metrics web-01
sloth-runner agent metrics web-01 --format json
sloth-runner agent metrics web-01 --watch              # Atualiza continuamente
sloth-runner agent metrics web-01 --interval 5         # Intervalo de 5 segundos
```

**Opções:**
- `--format` - Formato: table, json, yaml, prometheus
- `--watch` - Atualiza continuamente
- `--interval` - Intervalo de atualização em segundos (padrão: 2)

---

### `agent metrics grafana` - Dashboard Grafana

Gera e exibe dashboard do Grafana para um agente.

```bash
# Sintaxe
sloth-runner agent metrics grafana <agent-name> [opções]

# Exemplos
sloth-runner agent metrics grafana web-01
sloth-runner agent metrics grafana web-01 --export dashboard.json
```

**Opções:**
- `--export` - Exporta dashboard para arquivo JSON

---

## 📦 Gerenciamento de Sloths (Workflows Salvos)

### `sloth list` - Listar Sloths

Lista todos os workflows salvos no repositório local.

```bash
# Sintaxe
sloth-runner sloth list [opções]

# Exemplos
sloth-runner sloth list                   # Lista todos
sloth-runner sloth list --active          # Apenas sloths ativos
sloth-runner sloth list --inactive        # Apenas sloths inativos
sloth-runner sloth list --format json     # Saída em JSON
```

**Opções:**
- `--active` - Apenas sloths ativos
- `--inactive` - Apenas sloths inativos
- `--format` - Formato: table, json, yaml

---

### `sloth add` - Adicionar Sloth

Adiciona um novo workflow ao repositório.

```bash
# Sintaxe
sloth-runner sloth add <name> --file <caminho> [opções]

# Exemplos
sloth-runner sloth add deploy --file deploy.sloth
sloth-runner sloth add deploy --file deploy.sloth --description "Deploy production"
sloth-runner sloth add deploy --file deploy.sloth --tags "prod,deploy"
```

**Opções:**
- `--file` - Caminho do arquivo Sloth (obrigatório)
- `--description` - Descrição do sloth
- `--tags` - Tags separadas por vírgula

---

### `sloth get` - Obter Sloth

Exibe detalhes de um sloth específico.

```bash
# Sintaxe
sloth-runner sloth get <name> [opções]

# Exemplos
sloth-runner sloth get deploy
sloth-runner sloth get deploy --format json
sloth-runner sloth get deploy --show-content    # Mostra conteúdo do workflow
```

**Opções:**
- `--format` - Formato: table, json, yaml
- `--show-content` - Mostra conteúdo completo do workflow

---

### `sloth update` - Atualizar Sloth

Atualiza um sloth existente.

```bash
# Sintaxe
sloth-runner sloth update <name> [opções]

# Exemplos
sloth-runner sloth update deploy --file deploy-v2.sloth
sloth-runner sloth update deploy --description "New description"
sloth-runner sloth update deploy --tags "prod,deploy,updated"
```

**Opções:**
- `--file` - Novo arquivo Sloth
- `--description` - Nova descrição
- `--tags` - Novas tags

---

### `sloth remove` - Remover Sloth

Remove um sloth do repositório.

```bash
# Sintaxe
sloth-runner sloth remove <name>

# Exemplos
sloth-runner sloth remove deploy
sloth-runner sloth remove deploy --force    # Remove sem confirmação
```

**Opções:**
- `--force` - Remove sem pedir confirmação

---

### `sloth activate` - Ativar Sloth

Ativa um sloth desativado.

```bash
# Sintaxe
sloth-runner sloth activate <name>

# Exemplos
sloth-runner sloth activate deploy
```

---

### `sloth deactivate` - Desativar Sloth

Desativa um sloth (não remove, apenas marca como inativo).

```bash
# Sintaxe
sloth-runner sloth deactivate <name>

# Exemplos
sloth-runner sloth deactivate deploy
```

---

## 🎣 Gerenciamento de Hooks

### `hook list` - Listar Hooks

Lista todos os hooks registrados.

```bash
# Sintaxe
sloth-runner hook list [opções]

# Exemplos
sloth-runner hook list
sloth-runner hook list --format json
sloth-runner hook list --event workflow.started    # Filtra por evento
```

**Opções:**
- `--format` - Formato: table, json, yaml
- `--event` - Filtra por tipo de evento

---

### `hook add` - Adicionar Hook

Adiciona um novo hook.

```bash
# Sintaxe
sloth-runner hook add <name> --event <evento> --script <caminho> [opções]

# Exemplos
sloth-runner hook add notify-slack --event workflow.completed --script notify.sh
sloth-runner hook add backup --event task.completed --script backup.lua
sloth-runner hook add validate --event workflow.started --script validate.lua --priority 10
```

**Opções:**
- `--event` - Tipo de evento (obrigatório)
- `--script` - Caminho do script (obrigatório)
- `--priority` - Prioridade de execução (padrão: 0)
- `--enabled` - Hook habilitado (padrão: true)

**Eventos disponíveis:**
- `workflow.started`
- `workflow.completed`
- `workflow.failed`
- `task.started`
- `task.completed`
- `task.failed`
- `agent.connected`
- `agent.disconnected`

---

### `hook remove` - Remover Hook

Remove um hook.

```bash
# Sintaxe
sloth-runner hook remove <name>

# Exemplos
sloth-runner hook remove notify-slack
sloth-runner hook remove notify-slack --force
```

---

### `hook enable` - Habilitar Hook

Habilita um hook desabilitado.

```bash
# Sintaxe
sloth-runner hook enable <name>

# Exemplos
sloth-runner hook enable notify-slack
```

---

### `hook disable` - Desabilitar Hook

Desabilita um hook.

```bash
# Sintaxe
sloth-runner hook disable <name>

# Exemplos
sloth-runner hook disable notify-slack
```

---

### `hook test` - Testar Hook

Testa a execução de um hook.

```bash
# Sintaxe
sloth-runner hook test <name> [opções]

# Exemplos
sloth-runner hook test notify-slack
sloth-runner hook test notify-slack --payload '{"message": "test"}'
```

**Opções:**
- `--payload` - JSON com dados de teste

---

## 📡 Gerenciamento de Eventos

### `events list` - Listar Eventos

Lista eventos recentes do sistema.

```bash
# Sintaxe
sloth-runner events list [opções]

# Exemplos
sloth-runner events list
sloth-runner events list --limit 50               # Últimos 50 eventos
sloth-runner events list --type workflow.started  # Filtra por tipo
sloth-runner events list --since 1h               # Eventos da última hora
sloth-runner events list --format json
```

**Opções:**
- `--limit` - Número máximo de eventos (padrão: 100)
- `--type` - Filtra por tipo de evento
- `--since` - Filtra por tempo (ex: 1h, 30m, 24h)
- `--format` - Formato: table, json, yaml

---

### `events watch` - Monitorar Eventos

Monitora eventos em tempo real.

```bash
# Sintaxe
sloth-runner events watch [opções]

# Exemplos
sloth-runner events watch
sloth-runner events watch --type workflow.completed    # Apenas eventos de workflow completado
sloth-runner events watch --filter "status=success"    # Com filtro
```

**Opções:**
- `--type` - Filtra por tipo de evento
- `--filter` - Expressão de filtro

---

## 🗄️ Gerenciamento de Database

### `db backup` - Backup do Database

Cria backup do database SQLite.

```bash
# Sintaxe
sloth-runner db backup [opções]

# Exemplos
sloth-runner db backup
sloth-runner db backup --output /backup/sloth-backup.db
sloth-runner db backup --compress                     # Comprime com gzip
```

**Opções:**
- `--output` - Caminho do arquivo de backup
- `--compress` - Comprime o backup

---

### `db restore` - Restaurar Database

Restaura database de um backup.

```bash
# Sintaxe
sloth-runner db restore <arquivo-backup> [opções]

# Exemplos
sloth-runner db restore /backup/sloth-backup.db
sloth-runner db restore /backup/sloth-backup.db.gz --decompress
```

**Opções:**
- `--decompress` - Descomprime backup gzip

---

### `db vacuum` - Otimizar Database

Otimiza e compacta o database SQLite.

```bash
# Sintaxe
sloth-runner db vacuum

# Exemplos
sloth-runner db vacuum
```

---

### `db stats` - Estatísticas do Database

Mostra estatísticas do database.

```bash
# Sintaxe
sloth-runner db stats [opções]

# Exemplos
sloth-runner db stats
sloth-runner db stats --format json
```

**Opções:**
- `--format` - Formato: table, json, yaml

---

## 🌐 SSH Management

### `ssh list` - Listar Conexões SSH

Lista conexões SSH salvas.

```bash
# Sintaxe
sloth-runner ssh list [opções]

# Exemplos
sloth-runner ssh list
sloth-runner ssh list --format json
```

**Opções:**
- `--format` - Formato: table, json, yaml

---

### `ssh add` - Adicionar Conexão SSH

Adiciona uma nova conexão SSH.

```bash
# Sintaxe
sloth-runner ssh add <name> --host <host> --user <user> [opções]

# Exemplos
sloth-runner ssh add web-server --host 192.168.1.100 --user ubuntu
sloth-runner ssh add web-server --host 192.168.1.100 --user ubuntu --port 2222
sloth-runner ssh add web-server --host 192.168.1.100 --user ubuntu --key ~/.ssh/id_rsa
```

**Opções:**
- `--host` - Host SSH (obrigatório)
- `--user` - Usuário SSH (obrigatório)
- `--port` - Porta SSH (padrão: 22)
- `--key` - Caminho da chave SSH privada

---

### `ssh remove` - Remover Conexão SSH

Remove uma conexão SSH salva.

```bash
# Sintaxe
sloth-runner ssh remove <name>

# Exemplos
sloth-runner ssh remove web-server
```

---

### `ssh test` - Testar Conexão SSH

Testa uma conexão SSH.

```bash
# Sintaxe
sloth-runner ssh test <name>

# Exemplos
sloth-runner ssh test web-server
```

---

## 📋 Módulos

### `modules list` - Listar Módulos

Lista todos os módulos disponíveis.

```bash
# Sintaxe
sloth-runner modules list [opções]

# Exemplos
sloth-runner modules list
sloth-runner modules list --format json
sloth-runner modules list --category cloud         # Filtra por categoria
```

**Opções:**
- `--format` - Formato: table, json, yaml
- `--category` - Filtra por categoria

---

### `modules info` - Informações do Módulo

Mostra informações detalhadas sobre um módulo.

```bash
# Sintaxe
sloth-runner modules info <module-name>

# Exemplos
sloth-runner modules info pkg
sloth-runner modules info docker
sloth-runner modules info terraform
```

---

## 🖥️ Servidor e UI

### `server` - Iniciar Servidor Master

Inicia o servidor master (gRPC).

```bash
# Sintaxe
sloth-runner server [opções]

# Exemplos
sloth-runner server                          # Inicia na porta padrão (50053)
sloth-runner server --port 50053             # Define porta específica
sloth-runner server --bind 0.0.0.0           # Bind em todas as interfaces
sloth-runner server --tls-cert cert.pem --tls-key key.pem  # Com TLS
```

**Opções:**
- `--port` - Porta do servidor (padrão: 50053)
- `--bind` - Endereço de bind (padrão: 0.0.0.0)
- `--tls-cert` - Certificado TLS
- `--tls-key` - Chave privada TLS

---

### `ui` - Iniciar Web UI

Inicia a interface web.

```bash
# Sintaxe
sloth-runner ui [opções]

# Exemplos
sloth-runner ui                              # Inicia na porta padrão (8080)
sloth-runner ui --port 8080                  # Define porta específica
sloth-runner ui --bind 0.0.0.0               # Bind em todas as interfaces
```

**Opções:**
- `--port` - Porta da web UI (padrão: 8080)
- `--bind` - Endereço de bind (padrão: 0.0.0.0)

---

### `terminal` - Terminal Interativo

Abre terminal interativo para um agente remoto.

```bash
# Sintaxe
sloth-runner terminal <agent-name>

# Exemplos
sloth-runner terminal web-01
```

---

## 🔧 Utilitários

### `version` - Versão

Mostra a versão do Sloth Runner.

```bash
# Sintaxe
sloth-runner version

# Exemplos
sloth-runner version
sloth-runner version --format json
```

---

### `completion` - Auto-completar

Gera scripts de auto-completar para o shell.

```bash
# Sintaxe
sloth-runner completion <shell>

# Exemplos
sloth-runner completion bash > /etc/bash_completion.d/sloth-runner
sloth-runner completion zsh > ~/.zsh/completion/_sloth-runner
sloth-runner completion fish > ~/.config/fish/completions/sloth-runner.fish
```

**Shells suportados:** bash, zsh, fish, powershell

---

### `doctor` - Diagnóstico

Executa diagnóstico do sistema e configuração.

```bash
# Sintaxe
sloth-runner doctor [opções]

# Exemplos
sloth-runner doctor
sloth-runner doctor --format json
sloth-runner doctor --verbose             # Saída detalhada
```

**Opções:**
- `--format` - Formato: text, json
- `--verbose` - Saída detalhada

---

## 🔐 Variáveis de Ambiente

O Sloth Runner usa as seguintes variáveis de ambiente:

```bash
# Endereço do servidor master
export SLOTH_RUNNER_MASTER_ADDR="192.168.1.1:50053"

# Porta do agente
export SLOTH_RUNNER_AGENT_PORT="50060"

# Porta da Web UI
export SLOTH_RUNNER_UI_PORT="8080"

# Caminho do database
export SLOTH_RUNNER_DB_PATH="~/.sloth-runner/sloth.db"

# Nível de log
export SLOTH_RUNNER_LOG_LEVEL="info"  # debug, info, warn, error

# Habilitar modo debug
export SLOTH_RUNNER_DEBUG="true"
```

---

## 📊 Exemplos de Uso Comum

### 1. Deploy em Produção com Delegação

```bash
sloth-runner run production-deploy \
  --file deployments/prod.sloth \
  --delegate-to web-01 \
  --delegate-to web-02 \
  --values prod-vars.yaml \
  --yes
```

### 2. Monitorar Métricas de Todos os Agentes

```bash
# Em um terminal
sloth-runner agent metrics web-01 --watch

# Em outro terminal
sloth-runner agent metrics web-02 --watch
```

### 3. Backup Automatizado

```bash
# Criar backup comprimido com timestamp
sloth-runner db backup \
  --output /backup/sloth-$(date +%Y%m%d-%H%M%S).db \
  --compress
```

### 4. Workflow com Hook de Notificação

```bash
# Adicionar hook de notificação
sloth-runner hook add slack-notify \
  --event workflow.completed \
  --script /scripts/notify-slack.lua

# Executar workflow (hook será disparado automaticamente)
sloth-runner run deploy --file deploy.sloth --yes
```

### 5. Instalação de Agente em Múltiplos Servidores

```bash
# Loop para instalar em múltiplos hosts
for host in 192.168.1.{10..20}; do
  sloth-runner agent install "agent-$host" \
    --ssh-host "$host" \
    --ssh-user ubuntu \
    --master 192.168.1.1:50053
done
```

---

## 🎓 Próximos Passos

- [📖 Guia de Módulos](modulos-completos.md) - Documentação completa de todos os módulos
- [🎨 Web UI](web-ui-completo.md) - Guia completo da interface web
- [🎯 Exemplos Avançados](../en/advanced-examples.md) - Exemplos práticos de workflows
- [🏗️ Arquitetura](../architecture/sloth-runner-architecture.md) - Arquitetura do sistema

---

## 💡 Dicas e Truques

### Alias Úteis

Adicione ao seu `.bashrc` ou `.zshrc`:

```bash
alias sr='sloth-runner'
alias sra='sloth-runner agent'
alias srr='sloth-runner run'
alias srl='sloth-runner sloth list'
alias srui='sloth-runner ui --port 8080'
```

### Auto-completar

```bash
# Bash
sloth-runner completion bash > /etc/bash_completion.d/sloth-runner
source /etc/bash_completion.d/sloth-runner

# Zsh
sloth-runner completion zsh > ~/.zsh/completion/_sloth-runner
```

### Modo Debug

```bash
export SLOTH_RUNNER_DEBUG=true
export SLOTH_RUNNER_LOG_LEVEL=debug
sloth-runner run deploy --file deploy.sloth --verbose
```

---

**Última atualização:** 2025-10-07
