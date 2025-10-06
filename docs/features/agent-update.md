# 🔄 Agent Update via gRPC

## Overview

O comando `sloth-runner agent update` permite atualizar agentes remotos automaticamente via gRPC, sem necessidade de acesso SSH. O agente baixa e instala a atualização localmente, garantindo uma atualização segura e autônoma.

**Características principais:**
- ✅ Atualização via gRPC (sem SSH)
- ✅ Download direto do GitHub no agente
- ✅ Reinício automático após atualização
- ✅ Backup automático do binário
- ✅ Rollback em caso de falha

## Como Funciona

```
┌─────────┐                    ┌─────────────┐                    ┌────────────┐
│ Master  │                    │   Agent     │                    │  GitHub    │
└────┬────┘                    └──────┬──────┘                    └─────┬──────┘
     │                                │                                  │
     │ 1. UpdateAgent(gRPC) ─────────>│                                  │
     │                                │                                  │
     │                                │ 2. Fetch latest release ────────>│
     │                                │<─────────────────────────────────│
     │                                │                                  │
     │                                │ 3. Download binary ─────────────>│
     │                                │<─────────────────────────────────│
     │                                │                                  │
     │                                │ 4. Create update script          │
     │<───── 5. Success ──────────────│                                  │
     │                                │                                  │
     │                                │ 6. Shutdown (3s delay)           │
     │                                │                                  │
     │                                │ 7. Update script runs:           │
     │                                │    - Replace binary              │
     │                                │    - Restart agent               │
     │                                │                                  │
     │                                │ 8. Agent restarted               │
     │<───────── Heartbeat ───────────│                                  │
```

## Uso

### Atualizar um Agente

```bash
sloth-runner agent update <agent-name>
```

**Exemplo:**
```bash
sloth-runner agent update lady-arch
```

**Saída:**
```
Connecting to agent 'lady-arch'...
Sending update command to agent at 192.168.1.16:50052...
✓ Agent 'lady-arch' updated successfully

Old version: v4.18.3
New version: v4.18.4
Agent update prepared. Shutting down for binary replacement and restart.
```

### Atualizar para Versão Específica

```bash
sloth-runner agent update <agent-name> --version <version>
```

**Exemplo:**
```bash
sloth-runner agent update lady-arch --version v4.18.3
```

### Forçar Atualização

Força a atualização mesmo se já estiver na versão mais recente:

```bash
sloth-runner agent update <agent-name> --force
```

### Pular Reinício Automático

Atualiza mas não reinicia o agente automaticamente:

```bash
sloth-runner agent update <agent-name> --skip-restart
```

## Opções de Comando

| Flag | Descrição |
|------|-----------|
| `--force` | Força atualização mesmo se já estiver na última versão |
| `--version string` | Versão específica para instalar (padrão: latest) |
| `--skip-restart` | Pula reinício automático do serviço |

## Fluxo de Atualização Detalhado

### 1. Conexão gRPC

O master se conecta ao agente via gRPC obtendo o endereço do master registry:

```bash
Connecting to agent 'lady-arch'...
Sending update command to agent at 192.168.1.16:50052...
```

### 2. Verificação de Versão

O agente verifica a versão atual e compara com o GitHub:

```
Current version: v4.18.3
Latest version: v4.18.4
```

Se já estiver atualizado (sem `--force`), retorna imediatamente.

### 3. Download do Binário

O agente determina a arquitetura e baixa o binário apropriado:

- **Linux ARM64**: `sloth-runner_v4.18.4_linux_arm64.tar.gz`
- **Linux AMD64**: `sloth-runner_v4.18.4_linux_amd64.tar.gz`
- **Darwin ARM64**: `sloth-runner_v4.18.4_darwin_arm64.tar.gz`

### 4. Criação do Script de Atualização

O agente cria `/tmp/agent-update.sh`:

```bash
#!/bin/bash
# Wait for agent to stop
sleep 2

# Replace binary
cp /usr/local/bin/sloth-runner /usr/local/bin/sloth-runner.backup
mv /tmp/sloth-runner-new /usr/local/bin/sloth-runner
chmod +x /usr/local/bin/sloth-runner

# Restart service or agent
if systemctl is-active --quiet sloth-runner-agent 2>/dev/null; then
    systemctl restart sloth-runner-agent
elif systemctl is-active --quiet sloth-agent 2>/dev/null; then
    systemctl restart sloth-agent
else
    # Restart via nohup if no systemd service
    cd /home/igor && nohup /usr/local/bin/sloth-runner agent start \
        --name lady-arch \
        --master 192.168.1.29:50053 \
        --port 50051 \
        --bind-address 0.0.0.0 \
        --report-address 192.168.1.16:50052 \
        --telemetry \
        --metrics-port 9090 > agent.log 2>&1 &
fi

# Cleanup
rm -f /tmp/sloth-runner-new /usr/local/bin/sloth-runner.backup /tmp/agent-update.sh
```

### 5. Shutdown Programado

1. Agente responde sucesso ao master
2. Após 3 segundos, agente faz shutdown (`os.Exit(0)`)
3. Script de atualização executa em background
4. Binário é substituído
5. Agente reinicia automaticamente

## Exemplos Práticos

### Exemplo 1: Atualização Simples

```bash
$ sloth-runner agent update lady-arch
Connecting to agent 'lady-arch'...
Sending update command to agent at 192.168.1.16:50052...
✓ Agent 'lady-arch' updated successfully

Old version: v4.18.3
New version: v4.18.4
Agent update prepared. Shutting down for binary replacement and restart.
```

### Exemplo 2: Já Atualizado

```bash
$ sloth-runner agent update lady-arch
Connecting to agent 'lady-arch'...
Sending update command to agent at 192.168.1.16:50052...
✓ Agent 'lady-arch' updated successfully

Old version: v4.18.4
New version: v4.18.4
Already running the latest version
```

### Exemplo 3: Versão Específica

```bash
$ sloth-runner agent update lady-arch --version v4.18.2
Connecting to agent 'lady-arch'...
Sending update command to agent at 192.168.1.16:50052...
✓ Agent 'lady-arch' updated successfully

Old version: v4.18.4
New version: v4.18.2
Agent update prepared. Shutting down for binary replacement and restart.
```

## Integração com CI/CD

### GitHub Actions

O workflow de release atualiza automaticamente os agentes:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    name: Release
    runs-on: ubuntu-latest
    steps:
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

      - name: Update agent lady-arch
        run: |
          sloth-runner agent update lady-arch
```

### Script de Atualização em Lote

```bash
#!/bin/bash
# update-all-agents.sh

AGENTS=("lady-arch" "keite-guica" "prod-server-01")

for agent in "${AGENTS[@]}"; do
    echo "Updating agent: $agent"
    sloth-runner agent update "$agent" || echo "Failed to update $agent"

    # Aguarda 30 segundos antes do próximo
    sleep 30
done
```

### Atualização Programada com Cron

```bash
# crontab -e
# Atualiza agentes todas as segundas às 03:00
0 3 * * 1 /usr/local/bin/sloth-runner agent update lady-arch >> /var/log/agent-update.log 2>&1
```

## Requisitos

### No Master

- ✅ Conexão com o master registry
- ✅ Acesso ao agente via gRPC
- ✅ Rede acessível ao GitHub (para verificar versões)

### No Agente

- ✅ Conexão com GitHub para download
- ✅ Permissões para substituir `/usr/local/bin/sloth-runner`
- ✅ Espaço em disco em `/tmp` e `/usr/local/bin`
- ✅ Bash disponível para executar script de update

## Segurança

### Comunicação

- ✅ gRPC entre master e agente
- ✅ HTTPS para downloads do GitHub
- ✅ Apenas de `https://github.com/chalkan3-sloth/sloth-runner/releases`

### Backup Automático

O script sempre cria backup antes de substituir:

```bash
/usr/local/bin/sloth-runner.backup
```

### Rollback Manual

Se a atualização falhar, restaure o backup:

```bash
sudo mv /usr/local/bin/sloth-runner.backup /usr/local/bin/sloth-runner
sudo systemctl restart sloth-runner-agent
```

## Troubleshooting

### Agente Não Reconecta

**Problema**: Agente não aparece após atualização.

**Solução**: Verifique logs e reinicie manualmente:

```bash
# Via systemd
sudo systemctl status sloth-runner-agent
sudo systemctl start sloth-runner-agent

# Logs
sudo journalctl -u sloth-runner-agent -f

# Manualmente (se sem systemd)
cd /home/igor
nohup /usr/local/bin/sloth-runner agent start \
    --name lady-arch \
    --master 192.168.1.29:50053 \
    --port 50051 \
    --bind-address 0.0.0.0 \
    --report-address 192.168.1.16:50052 \
    > agent.log 2>&1 &
```

### Erro de Conexão gRPC

**Problema**: `ERROR Failed to connect to agent: context deadline exceeded`

**Solução**: Verifique conectividade:

```bash
# Verificar porta gRPC
nc -zv 192.168.1.16 50052

# Verificar se agente está rodando
sloth-runner agent list
```

### Timeout no Download

**Problema**: Download do GitHub demora muito ou falha.

**Solução**:

```bash
# Verificar acesso ao GitHub do agente
ssh agent-host "curl -I https://github.com"

# Download manual e upload
wget https://github.com/chalkan3-sloth/sloth-runner/releases/download/v4.18.4/sloth-runner_v4.18.4_linux_arm64.tar.gz
scp sloth-runner_v4.18.4_linux_arm64.tar.gz agent-host:/tmp/
```

### Permissões Insuficientes

**Problema**: `permission denied` ao substituir binário

**Solução**: Verifique permissões:

```bash
# Verificar proprietário e permissões
ls -la /usr/local/bin/sloth-runner

# Corrigir se necessário
sudo chown root:root /usr/local/bin/sloth-runner
sudo chmod 755 /usr/local/bin/sloth-runner
```

### Binário Ocupado

**Problema**: `text file busy` ao substituir

**Solução**: O script aguarda o shutdown. Se persistir:

```bash
# Matar processo manualmente
sudo pkill -9 sloth-runner

# Substituir binário
sudo mv /tmp/sloth-runner-new /usr/local/bin/sloth-runner
sudo chmod +x /usr/local/bin/sloth-runner

# Reiniciar
sudo systemctl start sloth-runner-agent
```

## Melhores Práticas

### 1. Teste Primeiro

- Atualize agentes de desenvolvimento/staging primeiro
- Verifique funcionamento antes de produção
- Teste com `--skip-restart` se necessário

### 2. Monitore Após Update

```bash
# Verificar status
sloth-runner agent list

# Ver logs do agente
ssh agent-host "tail -f /path/to/agent.log"

# Testar execução de task
sloth-runner run test-task --delegate-to agent-name
```

### 3. Mantenha Backups

- Backups são criados automaticamente em `.backup`
- Mantenha cópias dos binários importantes
- Documente configurações customizadas

### 4. Atualizações Programadas

- Use janelas de manutenção
- Atualize em lotes para grandes deploys
- Configure alertas para falhas

## Diferenças da Versão Anterior

### Antes (SSH)

❌ Requeria acesso SSH ao host
❌ Executava comandos remotamente via SSH
❌ Dependia de chaves SSH configuradas
❌ Mais complexo para troubleshooting

### Agora (gRPC)

✅ Usa canal gRPC existente
✅ Agente executa update localmente
✅ Sem dependência de SSH
✅ Mais simples e confiável

## Referências

- [Agent Setup](../agent-setup.md) - Como configurar agentes
- [Distributed Agents](../distributed-agents.md) - Arquitetura distribuída
- [GitHub Releases](https://github.com/chalkan3-sloth/sloth-runner/releases) - Versões disponíveis
- [gRPC Protocol](https://grpc.io/) - Sobre gRPC
