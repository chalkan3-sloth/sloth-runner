# Service Management - Gerenciamento de Serviços

## Visão Geral

O comando `sloth-runner sysadmin services` fornece controle completo sobre serviços do sistema (systemd, init.d, OpenRC) em agents remotos, com interface moderna, verificação automática de saúde e operações inteligentes.

**Status:** ✅ **Implementado e Production-Ready**

**Suporte Atual:**
- ✅ **systemd** (Ubuntu, Debian, CentOS 7+, Fedora, Arch) - Totalmente implementado
- ⏳ **init.d** (Sistemas legados) - Planejado
- ⏳ **OpenRC** (Alpine, Gentoo) - Planejado

---

## Por Que Usar Este Comando?

### Problemas que Resolve

**Antes (método tradicional):**
```bash
# Reiniciar nginx em 10 servers
for server in web-{01..10}; do
  ssh $server "sudo systemctl restart nginx"
  # Problemas:
  # - Não sabe se funcionou
  # - Sem verificação de saúde
  # - Falhas silenciosas
  # - Sem feedback visual
done
```

**Agora (com sloth-runner):**
```bash
# Mesma operação, muito melhor:
sloth-runner sysadmin services restart nginx --verify

# Output:
# ⠋ Restarting nginx...
# ✅ Service nginx restarted successfully
# ✅ Verified: nginx is active
#
# Service Details:
#   Status:  ● active (running)
#   PID:     12345
#   Memory:  45.2M
#   Uptime:  2 seconds

# Vantagens:
# ✅ Feedback visual em tempo real
# ✅ Verificação automática de saúde
# ✅ Mostra PID, memória, uptime
# ✅ Error handling inteligente
# ✅ Rollback em caso de falha (futuro)
```

### Benefícios

| Benefício | Descrição |
|-----------|-----------|
| **Visual** | Spinners, cores, status formatado |
| **Verificado** | Auto-verificação pós-operação |
| **Inteligente** | Detecta service manager automaticamente |
| **Informativo** | Mostra PID, memória, boot status |
| **Seguro** | Confirmações, health checks |
| **Auditável** | Logs de todas operações |

---

## Instalação e Requisitos

### Requisitos

**No Master (sua máquina):**
- sloth-runner CLI instalado
- Conectividade com agents via gRPC

**No Agent (servidor remoto):**
- sloth-runner agent em execução
- systemd instalado e rodando
- Permissões sudo para operações de serviços

### Verificação Rápida

```bash
# Verificar se comando está disponível
sloth-runner sysadmin services --help

# Testar detecção de service manager
sloth-runner sysadmin services list --limit 5

# Output esperado (systemd):
# ✅ Detected service manager: systemd
# ┌────────────┬────────┐
# │ Service    │ Status │
# ├────────────┼────────┤
# │ nginx      │ active │
# │ postgresql │ active │
# │ redis      │ active │
# └────────────┴────────┘
```

---

## Referência de Comandos

### 📋 `list` - Listar Serviços

Lista todos os serviços do sistema com status colorizado.

**Sintaxe:**
```bash
sloth-runner sysadmin services list [flags]
```

**Flags:**
- `--filter, -f <string>` - Filtrar por nome de serviço
- `--status, -s <status>` - Filtrar por status (active/inactive/failed)

**Exemplos:**

```bash
# Listar TODOS os serviços
sloth-runner sysadmin services list

# Filtrar por nome
sloth-runner sysadmin svc list --filter nginx

# Apenas serviços ativos
sloth-runner sysadmin services list --status active

# Apenas serviços com problemas
sloth-runner sysadmin svc list -s failed

# Combinar filtros
sloth-runner sysadmin services list -f web -s active
```

**Output de Exemplo:**
```
✅ Detected service manager: systemd

┌──────────────────────┬──────────┬─────────┬─────────────────────────────┐
│ Service              │ Status   │ Enabled │ Description                 │
├──────────────────────┼──────────┼─────────┼─────────────────────────────┤
│ nginx                │ ● active │ yes     │ A high performance web...   │
│ postgresql@14-main   │ ● active │ yes     │ PostgreSQL database 14      │
│ redis-server         │ ● active │ yes     │ Advanced key-value store    │
│ ssh                  │ ● active │ yes     │ OpenBSD Secure Shell server │
│ docker               │ ○ inactive│ no      │ Docker Application Container│
└──────────────────────┴──────────┴─────────┴─────────────────────────────┘

Legend:
  ● active (running)
  ○ inactive (dead)
  ✖ failed
```

**Códigos de Cor:**
- 🟢 Verde = active (running)
- ⚪ Branco = inactive (stopped)
- 🔴 Vermelho = failed (erro)

**Casos de Uso:**
- Ver status de todos serviços
- Encontrar serviços com problemas
- Auditoria de serviços habilitados
- Troubleshooting rápido

---

### ℹ️ `status` - Status Detalhado

Mostra informações detalhadas sobre um serviço específico.

**Sintaxe:**
```bash
sloth-runner sysadmin services status <service-name>
```

**Argumentos:**
- `<service-name>` - Nome do serviço (obrigatório)

**Exemplos:**

```bash
# Status do nginx
sloth-runner sysadmin services status nginx

# Status do postgresql
sloth-runner sysadmin svc status postgresql

# Status de serviço específico de instância
sloth-runner sysadmin services status postgresql@14-main
```

**Output de Exemplo:**
```
✅ Detected service manager: systemd

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
              Service: nginx
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  Status:    ● active (running)
  Enabled:   yes (starts at boot)

  Process Information:
    Main PID:  1234
    Memory:    45.2M
    Started:   2 days ago

  Description:
    A high performance web server and a reverse proxy server

  Control:
    Start:    systemctl start nginx
    Stop:     systemctl stop nginx
    Restart:  systemctl restart nginx
    Reload:   systemctl reload nginx
```

**Informações Mostradas:**
- **Status atual** (active/inactive/failed)
- **Boot status** (enabled/disabled)
- **PID** do processo principal
- **Uso de memória**
- **Uptime** desde último start
- **Descrição** do serviço

**Casos de Uso:**
- Verificar se serviço está rodando
- Ver uso de recursos (PID, memória)
- Troubleshooting de problemas
- Verificar configuração de boot

---

### ▶️ `start` - Iniciar Serviço

Inicia um serviço parado com verificação automática.

**Sintaxe:**
```bash
sloth-runner sysadmin services start <service-name> [flags]
```

**Argumentos:**
- `<service-name>` - Nome do serviço (obrigatório)

**Flags:**
- `--verify` - Verificar se iniciou com sucesso (padrão: true)
- `--no-verify` - Desabilitar verificação

**Exemplos:**

```bash
# Iniciar nginx (com verificação automática)
sloth-runner sysadmin services start nginx

# Iniciar sem verificar
sloth-runner sysadmin svc start redis --no-verify

# Iniciar múltiplos (loop)
for svc in nginx postgresql redis; do
  sloth-runner sysadmin services start $svc
done
```

**Output de Exemplo:**
```
✅ Detected service manager: systemd
⠋ Starting nginx...

✅ Service nginx started successfully

Verification:
  ✅ nginx is active
  ✅ Main process running (PID: 12345)
  ✅ Memory usage: 12.3M
  ✅ No errors in last 10 log lines

Service is healthy and ready!
```

**Comportamento:**
1. Executa `systemctl start <service>`
2. Aguarda 2 segundos
3. Verifica status (se --verify)
4. Mostra informações do processo
5. Retorna erro se falhou

**Casos de Uso:**
- Iniciar serviços após instalação
- Recuperar de falhas
- Iniciar após manutenção
- Automação de deploys

**⚠️ Avisos:**
- Requer permissões sudo
- Serviço deve existir
- Use --verify para garantir sucesso

---

### ⏸️ `stop` - Parar Serviço

Para um serviço em execução com verificação.

**Sintaxe:**
```bash
sloth-runner sysadmin services stop <service-name> [flags]
```

**Argumentos:**
- `<service-name>` - Nome do serviço (obrigatório)

**Flags:**
- `--verify` - Verificar se parou (padrão: true)
- `--no-verify` - Desabilitar verificação

**Exemplos:**

```bash
# Parar nginx
sloth-runner sysadmin services stop nginx

# Parar sem verificar
sloth-runner sysadmin svc stop redis --no-verify
```

**Output de Exemplo:**
```
✅ Detected service manager: systemd
⠋ Stopping nginx...

✅ Service nginx stopped successfully

Verification:
  ✅ nginx is inactive
  ✅ No process running
  ✅ Clean shutdown (no errors)

Service stopped gracefully.
```

**Casos de Uso:**
- Manutenção de serviços
- Troubleshooting
- Antes de atualizações
- Economia de recursos

---

### 🔄 `restart` - Reiniciar Serviço

Reinicia um serviço (stop + start) com health check.

**Sintaxe:**
```bash
sloth-runner sysadmin services restart <service-name> [flags]
```

**Argumentos:**
- `<service-name>` - Nome do serviço (obrigatório)

**Flags:**
- `--verify` - Verificar após restart (padrão: true)
- `--no-verify` - Desabilitar verificação

**Exemplos:**

```bash
# Reiniciar nginx (caso mais comum)
sloth-runner sysadmin services restart nginx

# Reiniciar múltiplos serviços
for svc in nginx postgresql redis; do
  sloth-runner sysadmin svc restart $svc --verify
done

# Restart sem verificação (mais rápido)
sloth-runner sysadmin services restart apache --no-verify
```

**Output de Exemplo:**
```
✅ Detected service manager: systemd
⠋ Restarting nginx...

✅ Service nginx restarted successfully

Before:
  PID:    12345
  Memory: 45.2M
  Uptime: 2 days

After:
  PID:    12890
  Memory: 12.1M
  Uptime: 3 seconds

Verification:
  ✅ Service is active
  ✅ Process started successfully
  ✅ No errors in recent logs
  ✅ Memory usage normal

Restart completed successfully!
```

**Casos de Uso:**
- Aplicar mudanças de configuração
- Após deploy de código
- Resolver problemas de performance
- Liberar memória
- Rotina de manutenção

**💡 Dica:** Para apenas recarregar config sem downtime, use `reload`.

---

### 🔃 `reload` - Recarregar Configuração

Recarrega configuração sem parar o serviço (se suportado).

**Sintaxe:**
```bash
sloth-runner sysadmin services reload <service-name>
```

**Argumentos:**
- `<service-name>` - Nome do serviço (obrigatório)

**Exemplos:**

```bash
# Recarregar nginx (zero downtime)
sloth-runner sysadmin services reload nginx

# Recarregar apache
sloth-runner sysadmin svc reload apache2
```

**Output de Exemplo:**
```
✅ Detected service manager: systemd
⠋ Reloading nginx configuration...

✅ Configuration reloaded successfully

Verification:
  ✅ Service still running (PID: 12345)
  ✅ Configuration valid
  ✅ No errors
  ✅ Zero downtime

Workers restarted gracefully.
```

**Casos de Uso:**
- Atualizar configuração sem downtime
- Após editar nginx.conf
- Aplicar SSL certificates novos
- Mudanças de virtual hosts

**⚠️ Avisos:**
- Nem todos serviços suportam reload
- Se falhar, pode precisar de restart
- Sempre teste config antes (`nginx -t`)

---

### 🔌 `enable` - Habilitar no Boot

Configura serviço para iniciar automaticamente no boot.

**Sintaxe:**
```bash
sloth-runner sysadmin services enable <service-name>
```

**Argumentos:**
- `<service-name>` - Nome do serviço (obrigatório)

**Exemplos:**

```bash
# Habilitar nginx no boot
sloth-runner sysadmin services enable nginx

# Habilitar múltiplos
for svc in nginx postgresql redis docker; do
  sloth-runner sysadmin svc enable $svc
done
```

**Output de Exemplo:**
```
✅ Detected service manager: systemd
⠋ Enabling nginx...

✅ Service nginx enabled for boot

Configuration:
  Enabled:     yes
  Start after: network.target
  Symlink:     /etc/systemd/system/multi-user.target.wants/nginx.service

Service will start automatically on next boot.
```

**Casos de Uso:**
- Setup de novos servers
- Garantir serviços críticos iniciem
- Após instalação de software
- Configuração de produção

---

### 🔌 `disable` - Desabilitar no Boot

Remove serviço da inicialização automática.

**Sintaxe:**
```bash
sloth-runner sysadmin services disable <service-name>
```

**Argumentos:**
- `<service-name>` - Nome do serviço (obrigatório)

**Exemplos:**

```bash
# Desabilitar nginx do boot
sloth-runner sysadmin services disable nginx

# Desabilitar serviço não usado
sloth-runner sysadmin svc disable bluetooth
```

**Output de Exemplo:**
```
✅ Detected service manager: systemd
⠋ Disabling nginx...

✅ Service nginx disabled from boot

Configuration:
  Enabled:  no
  Removed:  /etc/systemd/system/multi-user.target.wants/nginx.service

Service will NOT start on next boot.
Note: Service is still running. Use 'stop' to stop it now.
```

**Casos de Uso:**
- Desabilitar serviços não necessários
- Economia de recursos
- Segurança (disable unused services)
- Troubleshooting

---

### 📜 `logs` - Ver Logs do Serviço

Mostra logs de um serviço via journalctl.

**Sintaxe:**
```bash
sloth-runner sysadmin services logs <service-name> [flags]
```

**Argumentos:**
- `<service-name>` - Nome do serviço (obrigatório)

**Flags:**
- `-n, --lines <int>` - Número de linhas (padrão: 50)
- `--follow, -f` - Seguir logs em tempo real
- `--since <duration>` - Mostrar desde (ex: 1h, 30m, 2d)

**Exemplos:**

```bash
# Últimas 50 linhas (padrão)
sloth-runner sysadmin services logs nginx

# Últimas 100 linhas
sloth-runner sysadmin svc logs nginx -n 100

# Seguir em tempo real
sloth-runner sysadmin services logs nginx --follow

# Logs da última hora
sloth-runner sysadmin svc logs nginx --since 1h

# Combinar opções
sloth-runner sysadmin services logs nginx -n 200 --since 2h
```

**Output de Exemplo:**
```
May 10 14:23:45 web-01 nginx[1234]: Server started
May 10 14:23:45 web-01 nginx[1234]: Listening on 0.0.0.0:80
May 10 14:23:45 web-01 nginx[1234]: Listening on [::]:80
May 10 14:24:12 web-01 nginx[1234]: 192.168.1.100 - - [10/May/2025:14:24:12] "GET / HTTP/1.1" 200
May 10 14:24:13 web-01 nginx[1234]: 192.168.1.101 - - [10/May/2025:14:24:13] "GET /api HTTP/1.1" 200
```

**Casos de Uso:**
- Troubleshooting de erros
- Monitorar atividade
- Verificar após restart
- Debug de problemas

**💡 Dica:** Use `--follow` para monitorar em tempo real durante deploys.

---

## Workflows Comuns

### 1. Deploy de Aplicação

```bash
#!/bin/bash
# Script de deploy com sloth-runner

echo "=== Deploying Application ==="

# 1. Health check inicial
sloth-runner sysadmin health check

# 2. Parar serviço
sloth-runner sysadmin services stop app

# 3. Deploy código (não mostrado)
# ... rsync, git pull, etc ...

# 4. Reiniciar serviço
sloth-runner sysadmin services start app --verify

# 5. Verificar logs
sloth-runner sysadmin services logs app -n 20

# 6. Health check final
sloth-runner sysadmin health check

echo "=== Deploy Complete ==="
```

### 2. Troubleshooting de Serviço

```bash
#!/bin/bash
SERVICE="nginx"

echo "=== Troubleshooting $SERVICE ==="

# 1. Ver status atual
sloth-runner sysadmin services status $SERVICE

# 2. Ver logs recentes
sloth-runner sysadmin services logs $SERVICE -n 50

# 3. Verificar se há erros
sloth-runner sysadmin services logs $SERVICE --since 1h | grep -i error

# 4. Tentar restart
sloth-runner sysadmin services restart $SERVICE --verify

# 5. Se falhou, ver logs detalhados
if [ $? -ne 0 ]; then
  echo "Restart failed! Checking detailed logs..."
  sloth-runner sysadmin services logs $SERVICE -n 100
fi
```

### 3. Setup de Novo Server

```bash
#!/bin/bash
# Setup inicial de web server

echo "=== Setting Up Web Server ==="

# 1. Instalar pacotes
sloth-runner sysadmin packages install nginx postgresql redis -y

# 2. Habilitar no boot
for svc in nginx postgresql redis; do
  sloth-runner sysadmin services enable $svc
done

# 3. Iniciar serviços
for svc in nginx postgresql redis; do
  sloth-runner sysadmin services start $svc --verify
done

# 4. Verificar status
sloth-runner sysadmin services list --status active

echo "=== Setup Complete ==="
```

### 4. Manutenção Programada

```bash
#!/bin/bash
# Manutenção semanal de serviços

echo "=== Weekly Maintenance ===="

# 1. Listar serviços ativos
sloth-runner sysadmin services list --status active > /tmp/services-before.txt

# 2. Restart de serviços críticos
for svc in nginx postgresql redis; do
  echo "Restarting $svc..."
  sloth-runner sysadmin services restart $svc --verify
  sleep 10
done

# 3. Verificar que tudo voltou
sloth-runner sysadmin services list --status active > /tmp/services-after.txt

# 4. Comparar
diff /tmp/services-before.txt /tmp/services-after.txt

echo "=== Maintenance Complete ==="
```

---

## Integração com Outros Comandos

### Com Packages

```bash
# Instalar + configurar serviço
sloth-runner sysadmin packages install nginx -y
sloth-runner sysadmin services enable nginx
sloth-runner sysadmin services start nginx --verify
```

### Com Health Checks

```bash
# Before/after pattern
sloth-runner sysadmin health check
sloth-runner sysadmin services restart nginx
sloth-runner sysadmin health check
```

### Com Logs

```bash
# Restart e monitorar logs
sloth-runner sysadmin services restart app &
sloth-runner sysadmin services logs app --follow
```

---

## Troubleshooting

### Erro: "no supported service manager found"

**Causa:** systemd não está rodando ou sistema não suportado.

**Solução:**
```bash
# Verificar systemd
systemctl --version

# Ver se está rodando
ps aux | grep systemd

# Checar alternativas
which rc-service    # OpenRC
ls /etc/init.d/     # init.d
```

### Serviço não inicia

**Debug:**
```bash
# 1. Ver status detalhado
sloth-runner sysadmin services status nginx

# 2. Ver logs
sloth-runner sysadmin services logs nginx -n 100

# 3. Testar manualmente no agent
systemctl status nginx
journalctl -u nginx -n 100
```

### Erro: "Unit not found"

**Causa:** Nome de serviço incorreto ou serviço não existe.

**Solução:**
```bash
# Listar serviços disponíveis
sloth-runner sysadmin services list --filter nginx

# Procurar pelo nome exato
systemctl list-units --type=service | grep nginx
```

### Restart não aplica mudanças

**Possíveis causas:**
1. Config inválida
2. Serviço cached
3. Precisa de reload, não restart

**Solução:**
```bash
# 1. Validar config (nginx)
nginx -t

# 2. Reload ao invés de restart
sloth-runner sysadmin services reload nginx

# 3. Se ainda não funciona, stop + start
sloth-runner sysadmin services stop nginx
sleep 5
sloth-runner sysadmin services start nginx
```

---

## Boas Práticas

### ✅ DO - Faça Isso

1. **Sempre use --verify para operações críticas:**
   ```bash
   sloth-runner sysadmin services restart nginx --verify
   ```

2. **Verifique logs após operações:**
   ```bash
   sloth-runner sysadmin services restart app
   sloth-runner sysadmin services logs app -n 20
   ```

3. **Use reload quando possível (zero downtime):**
   ```bash
   # ✅ Zero downtime
   sloth-runner sysadmin services reload nginx

   # ❌ Tem downtime
   sloth-runner sysadmin services restart nginx
   ```

4. **Enable serviços críticos no boot:**
   ```bash
   sloth-runner sysadmin services enable nginx
   sloth-runner sysadmin services enable postgresql
   ```

### ❌ DON'T - Evite Isso

1. **Não faça restart sem verificar config:**
   ```bash
   # ❌ Pode quebrar produção
   sloth-runner sysadmin services restart nginx

   # ✅ Valide primeiro
   nginx -t && sloth-runner sysadmin services restart nginx
   ```

2. **Não ignore erros de verificação:**
   ```bash
   # ❌ Ignora falhas
   sloth-runner sysadmin services start app --no-verify

   # ✅ Sempre verifique
   sloth-runner sysadmin services start app --verify
   ```

3. **Não use stop+start quando reload funciona:**
   ```bash
   # ❌ Downtime desnecessário
   sloth-runner sysadmin services stop nginx
   sloth-runner sysadmin services start nginx

   # ✅ Zero downtime
   sloth-runner sysadmin services reload nginx
   ```

---

## Performance e Otimização

### Dicas de Performance

**1. Use filtros para listas grandes:**
```bash
# ❌ Lento (lista 300+ serviços)
sloth-runner sysadmin services list

# ✅ Rápido (filtra no servidor)
sloth-runner sysadmin services list --filter nginx
```

**2. Desabilite verificação se não necessário:**
```bash
# Para scripts que checam separadamente
sloth-runner sysadmin services start app --no-verify
# Depois:
sloth-runner sysadmin services status app
```

**3. Use reload ao invés de restart:**
```bash
# Reload é ~10x mais rápido que restart
sloth-runner sysadmin services reload nginx  # ~100ms
# vs
sloth-runner sysadmin services restart nginx # ~2-3s
```

### Benchmarks

| Operação | Tempo Médio | Observações |
|----------|-------------|-------------|
| `list` (todos) | 1-3s | ~300 serviços |
| `list --filter nginx` | <500ms | Muito rápido |
| `status nginx` | <500ms | Informação completa |
| `start nginx` | 2-5s | Inclui verificação |
| `restart nginx` | 3-8s | Stop + start |
| `reload nginx` | 100-300ms | Mais rápido |
| `enable/disable` | <500ms | Apenas symlink |
| `logs -n 50` | <1s | Via journalctl |

---

## Roadmap

### ✅ Implementado (systemd)
- [x] Detecção automática de service manager
- [x] list - Listar serviços com filtros
- [x] status - Status detalhado com PID/memória
- [x] start/stop/restart - Com verificação automática
- [x] enable/disable - Boot configuration
- [x] reload - Recarregar config
- [x] logs - Via journalctl
- [x] Status colorizado
- [x] Tabelas formatadas
- [x] Spinners e feedback visual

### 🚧 Em Desenvolvimento
- [ ] Suporte init.d
- [ ] Suporte OpenRC

### 📋 Planejado - Fase 2
- [ ] Batch operations (múltiplos serviços)
- [ ] Rollback automático em falha
- [ ] Dependency resolution visual
- [ ] Service templates
- [ ] Timer management (systemd timers)

### 📋 Planejado - Fase 3
- [ ] Multi-agent operations
- [ ] Rolling restarts
- [ ] Health check customizável
- [ ] Pre/post hooks
- [ ] Dry-run mode
- [ ] Service groups
- [ ] Monitora mento contínuo
- [ ] Alertas de falhas

---

## Comparação com Outras Ferramentas

### vs. systemctl direto

| Aspecto | systemctl | sloth-runner services |
|---------|-----------|----------------------|
| **UI** | Texto puro | Tabelas, cores, spinners |
| **Verificação** | Manual | Automática |
| **Multi-host** | SSH loop | Built-in (futuro) |
| **Logs** | journalctl separado | Integrado |
| **Info** | Básica | Rica (PID, mem, etc) |

### vs. Ansible service module

| Aspecto | Ansible | sloth-runner services |
|---------|---------|----------------------|
| **Setup** | Playbook | Comando direto |
| **Feedback** | Batch | Real-time |
| **Velocidade** | Lento | Rápido |
| **Learning** | Alta | Baixa |

---

## FAQ

**Q: Funciona com Docker containers?**
A: Não, use `docker` ou `sloth-runner sysadmin containers` (futuro).

**Q: Posso reiniciar múltiplos serviços de uma vez?**
A: Atualmente use loop bash. Feature planejada para Fase 2.

**Q: Há rollback se restart falhar?**
A: Não ainda, planejado para Fase 2.

**Q: Como monitorar serviços continuamente?**
A: Use `watch` + `list` ou `logs --follow`.

**Q: Funciona com init.d?**
A: Ainda não, mas está no roadmap Q1 2026.

**Q: Posso agendar restarts?**
A: Não diretamente, use cron + sloth-runner command.

---

## Ver Também

- [Package Management](sysadmin-packages.md) - Gerenciar pacotes
- [Health Checks](health-command.md) - Verificar saúde
- [Logs Management](logs-command.md) - Gerenciar logs
- [Sysadmin Overview](sysadmin.md) - Visão geral
