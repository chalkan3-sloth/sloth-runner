# Configurando Agentes Sloth-Runner em Containers Incus/LXC

Este guia explica como configurar agentes sloth-runner dentro de containers Incus (ou LXC), incluindo configuração de port forwarding e endereços de reporte.

## Quick Start

Para uma instalação rápida em container Incus:

```bash
# 1. No HOST (192.168.1.17) - Configure port forwarding
sudo incus config device add main sloth-proxy proxy \
  listen=tcp:0.0.0.0:50052 \
  connect=tcp:127.0.0.1:50051

# 2. No CONTAINER - Instale com bootstrap script
sudo incus exec main -- bash -c "curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh | bash -s -- --name main --master 192.168.1.29:50053 --incus 192.168.1.17:50052"

# Pronto! O agente já está rodando e configurado.
```

## Problema

Quando você executa um agente dentro de um container Incus, o container recebe um IP interno (ex: `10.193.121.186`) que não é acessível diretamente do master server. Isso causa timeouts quando você tenta executar comandos no agente.

## Solução

A solução envolve três passos principais:

1. Configurar port forwarding no host Incus
2. Configurar o agente para escutar em todas as interfaces
3. Usar o flag `--report-address` para informar ao master como se conectar

## Passo a Passo

### 1. Configure Port Forwarding no Host

No host que está rodando o Incus, adicione um dispositivo proxy para fazer o forward da porta:

```bash
sudo incus config device add <nome_container> sloth-proxy proxy \
  listen=tcp:0.0.0.0:<porta_host> \
  connect=tcp:127.0.0.1:<porta_agente>
```

**Exemplo prático:**
```bash
# Forward da porta 50052 do host para porta 50051 do container "main"
sudo incus config device add main sloth-proxy proxy \
  listen=tcp:0.0.0.0:50052 \
  connect=tcp:127.0.0.1:50051
```

### 2. Instale e Configure o Agente no Container

Dentro do container, instale o sloth-runner e inicie o agente:

#### Opção 1: Usando Bootstrap Script (Recomendado)

O bootstrap script agora suporta a flag `--incus` que configura automaticamente tudo:

```bash
# Dentro do container
bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \
  --name main \
  --master 192.168.1.29:50053 \
  --incus 192.168.1.17:50052
```

Isso configura automaticamente:
- `--bind-address 0.0.0.0` (escuta em todas as interfaces)
- `--report-address 192.168.1.17:50052` (IP do host + porta forwarded)
- Cria e habilita o serviço systemd

#### Opção 2: Instalação Manual

```bash
# Instalar o sloth-runner (adapte conforme seu método de instalação)
# Por exemplo, copie o binário:
# sudo cp /caminho/do/sloth-runner /usr/local/bin/

# Iniciar o agente com as configurações corretas
sloth-runner agent start \
  --name <nome_agente> \
  --master <ip_master>:<porta_master> \
  --port <porta_agente> \
  --bind-address 0.0.0.0 \
  --report-address <ip_host>:<porta_host> \
  --daemon
```

**Exemplo prático:**
```bash
# Container "main" no host 192.168.1.17, conectando ao master em 192.168.1.29
sloth-runner agent start \
  --name main \
  --master 192.168.1.29:50053 \
  --port 50051 \
  --bind-address 0.0.0.0 \
  --report-address 192.168.1.17:50052 \
  --daemon
```

**Parâmetros importantes:**
- `--bind-address 0.0.0.0`: Faz o agente escutar em todas as interfaces de rede
- `--report-address <ip_host>:<porta_host>`: Informa ao master qual endereço usar para se conectar ao agente

### 3. Configure como Serviço Systemd (Recomendado)

Para garantir que o agente inicie automaticamente com o container:

```bash
# Criar diretório de trabalho
sudo mkdir -p /var/lib/sloth-runner

# Criar arquivo de serviço
sudo tee /etc/systemd/system/sloth-runner-agent.service > /dev/null <<'EOF'
[Unit]
Description=Sloth Runner Agent - main
Documentation=https://chalkan3.github.io/sloth-runner/
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
WorkingDirectory=/var/lib/sloth-runner
Restart=always
RestartSec=5s
StartLimitInterval=60s
StartLimitBurst=5

# Configuração do Agente
ExecStart=/usr/local/bin/sloth-runner agent start \
  --name main \
  --master 192.168.1.29:50053 \
  --port 50051 \
  --bind-address 0.0.0.0 \
  --report-address 192.168.1.17:50052

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=sloth-runner-agent

# Performance
LimitNOFILE=65536

# Security
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

# Habilitar e iniciar o serviço
sudo systemctl daemon-reload
sudo systemctl enable sloth-runner-agent
sudo systemctl start sloth-runner-agent

# Verificar status
sudo systemctl status sloth-runner-agent
```

## Múltiplos Containers no Mesmo Host

Se você tem vários containers no mesmo host, cada um precisa de uma porta diferente:

```bash
# Container 1: main -> porta 50052
sudo incus config device add main sloth-proxy proxy \
  listen=tcp:0.0.0.0:50052 connect=tcp:127.0.0.1:50051

# Container 2: dev -> porta 50053
sudo incus config device add dev sloth-proxy proxy \
  listen=tcp:0.0.0.0:50053 connect=tcp:127.0.0.1:50051

# Container 3: staging -> porta 50054
sudo incus config device add staging sloth-proxy proxy \
  listen=tcp:0.0.0.0:50054 connect=tcp:127.0.0.1:50051
```

Então configure cada agente com seu respectivo `--report-address`:

```bash
# No container "main"
sloth-runner agent start --name main --report-address 192.168.1.17:50052 ...

# No container "dev"
sloth-runner agent start --name dev --report-address 192.168.1.17:50053 ...

# No container "staging"
sloth-runner agent start --name staging --report-address 192.168.1.17:50054 ...
```

## Verificação

### 1. Verificar Port Forwarding

```bash
# Listar dispositivos do container
sudo incus config device list <nome_container>

# Ver detalhes do proxy
sudo incus config device show <nome_container>
```

### 2. Verificar Status do Agente

```bash
# No container
sudo systemctl status sloth-runner-agent

# Ver logs
sudo journalctl -u sloth-runner-agent -f
```

### 3. Testar do Master

```bash
# Listar agentes
sloth-runner agent list

# Executar comando de teste
sloth-runner agent run <nome_agente> 'hostname && whoami'
```

## Tabela de Referência Rápida

| Componente | IP Interno | IP:Porta Exposto | Master Enxerga |
|------------|------------|------------------|----------------|
| Agente no Container | 10.x.x.x:50051 | host_ip:50052 | host_ip:50052 |
| Agente no Host | host_ip:50051 | host_ip:50051 | host_ip:50051 |

## Troubleshooting

### Agente aparece como "Active" mas comandos dão timeout

**Causa:** O master não consegue alcançar o agente no endereço reportado.

**Soluções:**
1. Verifique se o port forwarding está configurado:
   ```bash
   sudo incus config device list <nome_container>
   ```

2. Verifique se o agente está usando `--report-address` correto:
   ```bash
   sudo incus exec <nome_container> -- systemctl status sloth-runner-agent
   ```

3. Teste conectividade do master:
   ```bash
   # Do master, teste se a porta está acessível
   nc -zv <host_ip> <porta_forwarded>
   telnet <host_ip> <porta_forwarded>
   ```

4. Verifique firewall do host:
   ```bash
   # No host
   sudo iptables -L -n | grep <porta>
   sudo ufw status | grep <porta>
   ```

### Agente não inicia

**Verifique logs:**
```bash
sudo journalctl -u sloth-runner-agent -n 50
```

**Problemas comuns:**
- Binário não encontrado: Verifique `/usr/local/bin/sloth-runner` existe
- Permissões: O binário precisa ser executável (`chmod +x`)
- Master inacessível: Verifique se o master está rodando e acessível

### Container reiniciado e agente não volta

**Solução:** Certifique-se que o serviço systemd está habilitado:
```bash
sudo systemctl enable sloth-runner-agent
```

## Exemplo Completo

Aqui está um exemplo completo de configuração de agente no container "main":

```bash
# 1. No HOST (192.168.1.17) - Configure port forwarding
sudo incus config device add main sloth-proxy proxy \
  listen=tcp:0.0.0.0:50052 \
  connect=tcp:127.0.0.1:50051

# 2. No CONTAINER - Use bootstrap script com flag --incus
sudo incus exec main -- bash -c "curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh | bash -s -- --name main --master 192.168.1.29:50053 --incus 192.168.1.17:50052"

# OU se preferir fazer dentro do container interativamente:
sudo incus exec main -- bash

# Dentro do container:
bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \
  --name main \
  --master 192.168.1.29:50053 \
  --incus 192.168.1.17:50052

# Verificar status
systemctl status sloth-runner-agent
exit

# 3. Do MASTER (192.168.1.29) - Testar o agente
sloth-runner agent list
sloth-runner agent run main 'uname -a'
```

### Exemplo com Instalação Manual de Binário

Se você já tem o binário compilado:

```bash
# 1. No HOST - Configure port forwarding
sudo incus config device add main sloth-proxy proxy \
  listen=tcp:0.0.0.0:50052 \
  connect=tcp:127.0.0.1:50051

# 2. Copiar binário para o container
sudo incus file push /caminho/sloth-runner main/usr/local/bin/sloth-runner
sudo incus exec main -- chmod +x /usr/local/bin/sloth-runner

# 3. No CONTAINER - Use bootstrap local
sudo incus exec main -- bash
mkdir -p /var/lib/sloth-runner

# Criar e iniciar o serviço
/usr/local/bin/sloth-runner agent start \
  --name main \
  --master 192.168.1.29:50053 \
  --port 50051 \
  --bind-address 0.0.0.0 \
  --report-address 192.168.1.17:50052 \
  --daemon

# Verificar
ps aux | grep sloth-runner
exit

# 4. Do MASTER - Verificar
sloth-runner agent list
sloth-runner agent run main 'hostname && whoami'
```

## Conclusão

Com essa configuração, você pode executar agentes sloth-runner em containers Incus de forma transparente, permitindo que o master execute comandos remotamente como se fossem máquinas físicas normais.
