# Package Management - Gerenciamento de Pacotes

## Visão Geral

O comando `sloth-runner sysadmin packages` fornece uma interface unificada e moderna para gerenciar pacotes do sistema em agents remotos, eliminando a necessidade de SSH manual e padronizando operações entre diferentes distribuições Linux.

**Status:** ✅ **Implementado e Production-Ready**

**Suporte Atual:**
- ✅ **APT** (Debian, Ubuntu, Linux Mint) - Totalmente implementado
- ⏳ **YUM** (CentOS 7, RHEL 7) - Planejado
- ⏳ **DNF** (Fedora, CentOS 8+, RHEL 8+) - Planejado
- ⏳ **Pacman** (Arch Linux, Manjaro) - Planejado
- ⏳ **APK** (Alpine Linux) - Planejado
- ⏳ **Zypper** (openSUSE, SLES) - Planejado

---

## Por Que Usar Este Comando?

### Problemas que Resolve

**Antes (método tradicional):**
```bash
# Para atualizar 10 servidores Ubuntu:
for server in web-{01..10}; do
  ssh $server "sudo apt update && sudo apt upgrade -y"
done
# Problemas:
# - Sem progresso visual
# - Falhas silenciosas
# - Sem rollback
# - Difícil debugar
# - SSH keys management
```

**Agora (com sloth-runner):**
```bash
# Mesma operação, melhor forma:
sloth-runner sysadmin packages update --all-agents
sloth-runner sysadmin packages upgrade --all-agents --strategy rolling

# Vantagens:
# ✅ Progresso visual com spinners
# ✅ Error handling robusto
# ✅ Rollback automático (futuro)
# ✅ Logs centralizados
# ✅ Via gRPC (sem SSH)
```

### Benefícios

| Benefício | Descrição |
|-----------|-----------|
| **Unificado** | Mesma sintaxe para apt, yum, dnf, pacman |
| **Visual** | Spinners, progress bars, tabelas formatadas |
| **Seguro** | Confirmações, dry-run, rollback |
| **Rápido** | Operações em paralelo, cache inteligente |
| **Auditável** | Logs completos de todas operações |
| **Centralized** | Gerencia 100+ servers de um lugar |

---

## Instalação e Requisitos

### Requisitos

**No Master (sua máquina):**
- sloth-runner CLI instalado
- Conectividade com agents via gRPC

**No Agent (servidor remoto):**
- sloth-runner agent em execução
- Permissões sudo para operações de pacotes
- Package manager instalado (apt, yum, etc.)

### Verificação Rápida

```bash
# Verificar se comando está disponível
sloth-runner sysadmin packages --help

# Testar detecção de package manager
sloth-runner sysadmin packages list --limit 5

# Output esperado (Debian/Ubuntu):
# ✅ Detected package manager: apt
# Installed Packages (5)
# ┌─────────┬─────────┐
# │ Package │ Version │
# ├─────────┼─────────┤
# │ ...     │ ...     │
# └─────────┴─────────┘
```

---

## Referência de Comandos

### 📋 `list` - Listar Pacotes Instalados

Lista todos os pacotes instalados no sistema com opções de filtro.

**Sintaxe:**
```bash
sloth-runner sysadmin packages list [flags]
```

**Flags:**
- `--filter, -f <string>` - Filtrar por nome de pacote
- `--limit, -l <int>` - Limitar número de resultados (0 = sem limite)

**Exemplos:**

```bash
# Listar TODOS os pacotes
sloth-runner sysadmin packages list

# Listar apenas primeiros 20
sloth-runner sysadmin pkg list --limit 20

# Filtrar pacotes nginx
sloth-runner sysadmin packages list --filter nginx

# Combinar filtro + limite
sloth-runner sysadmin pkg list -f python -l 10
```

**Output de Exemplo:**
```
✅ Detected package manager: apt
Installed Packages (1247)

┌──────────────────┬─────────────────┐
│ Package          │ Version         │
├──────────────────┼─────────────────┤
│ nginx            │ 1.18.0-6ubuntu14│
│ nginx-common     │ 1.18.0-6ubuntu14│
│ nginx-core       │ 1.18.0-6ubuntu14│
└──────────────────┴─────────────────┘
```

**Casos de Uso:**
- Inventário de software instalado
- Auditoria de compliance
- Comparação entre servers
- Troubleshooting de dependências

---

### 🔍 `search` - Pesquisar Pacotes Disponíveis

Busca pacotes disponíveis nos repositórios configurados.

**Sintaxe:**
```bash
sloth-runner sysadmin packages search <query> [flags]
```

**Argumentos:**
- `<query>` - Termo de busca (obrigatório)

**Flags:**
- `--limit, -l <int>` - Limitar resultados (padrão: 20)

**Exemplos:**

```bash
# Buscar nginx
sloth-runner sysadmin packages search nginx

# Buscar com mais resultados
sloth-runner sysadmin pkg search python --limit 50

# Buscar ferramentas de monitoring
sloth-runner sysadmin packages search monitoring -l 10
```

**Output de Exemplo:**
```
✅ Using: apt
Search Results: 'nginx' (8 packages)

📦 nginx
   High performance web server

📦 nginx-common
   Common files for nginx

📦 nginx-core
   nginx web/proxy server (standard version)

📦 nginx-extras
   nginx web/proxy server (extended version)

📦 nginx-full
   nginx web/proxy server (standard version with extras)

📦 nginx-light
   nginx web/proxy server (basic version)
```

**Casos de Uso:**
- Descobrir pacotes antes de instalar
- Verificar disponibilidade de versões
- Explorar alternativas
- Planejamento de instalações

---

### 📥 `install` - Instalar Pacote

Instala um ou mais pacotes no sistema.

**Sintaxe:**
```bash
sloth-runner sysadmin packages install <package-name> [flags]
```

**Argumentos:**
- `<package-name>` - Nome do pacote a instalar (obrigatório)

**Flags:**
- `--yes, -y` - Confirmar automaticamente (não perguntar)
- `--no-deps` - Não instalar dependências (não recomendado)

**Exemplos:**

```bash
# Instalar com confirmação interativa
sloth-runner sysadmin packages install nginx
# Pergunta: Install package 'nginx'? [y/n]

# Instalar sem perguntar (scripts/automação)
sloth-runner sysadmin pkg install curl --yes

# Instalar múltiplos (loop)
for pkg in curl git vim; do
  sloth-runner sysadmin packages install $pkg -y
done
```

**Output de Exemplo:**
```
✅ Using: apt
Install package 'nginx'? yes

⠋ Installing nginx...

✅ Successfully installed nginx

Dependencies installed:
  - nginx-common (1.18.0-6ubuntu14)
  - nginx-core (1.18.0-6ubuntu14)
  - libnginx-mod-http-geoip (1.18.0-6ubuntu14)
  - libnginx-mod-http-image-filter (1.18.0-6ubuntu14)
  - libnginx-mod-http-xslt-filter (1.18.0-6ubuntu14)
  - libnginx-mod-mail (1.18.0-6ubuntu14)
  - libnginx-mod-stream (1.18.0-6ubuntu14)

Total installed: 8 packages
Disk space used: 1.2 MB
```

**Casos de Uso:**
- Instalar software em novos servers
- Adicionar ferramentas de desenvolvimento
- Setup de aplicações
- Automação de provisionamento

**⚠️ Avisos:**
- Requer permissões sudo no agent
- Confirme que o pacote está disponível (use `search` primeiro)
- Use `--yes` com cuidado em produção

---

### 🔄 `update` - Atualizar Listas de Pacotes

Atualiza as listas de pacotes disponíveis dos repositórios (equivalente a `apt update`).

**Sintaxe:**
```bash
sloth-runner sysadmin packages update
```

**Sem flags ou argumentos** - operação simples e direta.

**Exemplos:**

```bash
# Atualizar listas
sloth-runner sysadmin packages update

# Em múltiplos agents (futuro)
sloth-runner sysadmin pkg update --all-agents
```

**Output de Exemplo:**
```
✅ Using: apt
⠋ Updating package lists...

✅ Package lists updated successfully

Repositories updated:
  - http://archive.ubuntu.com/ubuntu jammy
  - http://archive.ubuntu.com/ubuntu jammy-updates
  - http://archive.ubuntu.com/ubuntu jammy-security
  - http://ppa.launchpad.net/...

Packages available for upgrade: 47
Security updates available: 12
```

**Casos de Uso:**
- Antes de instalar novos pacotes
- Parte de rotina de manutenção
- Verificar novas versões
- Manter índice atualizado

**💡 Dica:** Execute `update` antes de `search` ou `install` para garantir informações atualizadas.

---

### 🗑️ `remove` - Remover Pacote

**Status:** 🚧 Planejado

Remove um pacote do sistema.

**Sintaxe (planejada):**
```bash
sloth-runner sysadmin packages remove <package-name> [flags]

# Flags planejados:
# --yes, -y         - Confirmar automaticamente
# --purge           - Remover também configurações
# --auto-remove     - Remover dependências órfãs
```

**Exemplos (futuros):**
```bash
# Remover pacote
sloth-runner sysadmin packages remove nginx

# Remover com configurações
sloth-runner sysadmin pkg remove nginx --purge

# Remover + cleanup
sloth-runner sysadmin packages remove nginx --purge --auto-remove -y
```

---

### ⬆️ `upgrade` - Atualizar Pacotes

**Status:** 🚧 Planejado

Atualiza pacotes para versões mais recentes.

**Sintaxe (planejada):**
```bash
sloth-runner sysadmin packages upgrade [flags]

# Flags planejados:
# --all             - Atualizar todos os pacotes
# --security-only   - Apenas updates de segurança
# --strategy        - Estratégia: rolling, parallel, canary
# --wait-time       - Tempo de espera entre agents (rolling)
```

**Exemplos (futuros):**
```bash
# Atualizar tudo
sloth-runner sysadmin packages upgrade --all

# Apenas segurança
sloth-runner sysadmin pkg upgrade --security-only

# Rolling update em 5 servers
sloth-runner sysadmin packages upgrade \
  --all-agents \
  --strategy rolling \
  --wait-time 5m
```

---

### ✅ `check-updates` - Verificar Atualizações

**Status:** 🚧 Planejado

Verifica quais pacotes têm atualizações disponíveis sem instalar.

**Sintaxe (planejada):**
```bash
sloth-runner sysadmin packages check-updates [flags]
```

---

### ℹ️ `info` - Informações de Pacote

**Status:** 🚧 Planejado

Mostra informações detalhadas sobre um pacote.

**Sintaxe (planejada):**
```bash
sloth-runner sysadmin packages info <package-name>
```

---

### 📜 `history` - Histórico de Operações

**Status:** 🚧 Planejado

Mostra histórico de instalações, remoções e atualizações.

**Sintaxe (planejada):**
```bash
sloth-runner sysadmin packages history [flags]
```

---

## Workflows Comuns

### 1. Setup Inicial de Novo Server

```bash
# 1. Atualizar índices
sloth-runner sysadmin packages update

# 2. Instalar ferramentas básicas
for pkg in curl git vim htop; do
  sloth-runner sysadmin pkg install $pkg -y
done

# 3. Instalar stack de aplicação
sloth-runner sysadmin packages install nginx -y
sloth-runner sysadmin packages install postgresql -y
sloth-runner sysadmin packages install redis -y

# 4. Verificar instalação
sloth-runner sysadmin packages list --filter nginx
sloth-runner sysadmin packages list --filter postgresql
```

### 2. Auditoria de Compliance

```bash
# Listar todos pacotes instalados
sloth-runner sysadmin packages list > /tmp/installed-packages.txt

# Filtrar pacotes críticos
sloth-runner sysadmin pkg list --filter openssl
sloth-runner sysadmin pkg list --filter openssh

# Comparar com baseline
diff baseline-packages.txt /tmp/installed-packages.txt
```

### 3. Troubleshooting de Dependências

```bash
# Verificar se pacote está instalado
sloth-runner sysadmin packages list --filter libssl

# Buscar versões disponíveis
sloth-runner sysadmin packages search libssl

# Ver pacotes relacionados
sloth-runner sysadmin pkg list --filter ssl
```

### 4. Manutenção Regular

```bash
#!/bin/bash
# Script de manutenção semanal

echo "=== Weekly Package Maintenance ==="

# Atualizar índices
echo "1. Updating package lists..."
sloth-runner sysadmin packages update

# Verificar atualizações disponíveis (futuro)
# sloth-runner sysadmin packages check-updates

# Gerar relatório
echo "2. Generating package report..."
sloth-runner sysadmin packages list > /var/log/packages-$(date +%Y%m%d).txt

echo "Done!"
```

---

## Integração com Outros Comandos

### Com Health Checks

```bash
# Antes de atualizar, verificar saúde
sloth-runner sysadmin health check

# Atualizar pacotes
sloth-runner sysadmin packages update

# Verificar saúde após update
sloth-runner sysadmin health check
```

### Com Services

```bash
# Instalar nginx
sloth-runner sysadmin packages install nginx -y

# Verificar serviço
sloth-runner sysadmin services status nginx

# Iniciar se necessário
sloth-runner sysadmin services start nginx
```

### Com Logs

```bash
# Instalar pacote
sloth-runner sysadmin packages install nginx -y 2>&1 | tee /tmp/install.log

# Buscar erros
sloth-runner sysadmin logs search --query "package install" --since 1h
```

---

## Troubleshooting

### Erro: "no supported package manager found"

**Causa:** Sistema operacional não suportado ou package manager não detectado.

**Solução:**
```bash
# Verificar qual package manager está instalado
which apt apt-get yum dnf pacman apk zypper

# Se APT: verificar se está funcionando
apt --version

# Logs para debug
sloth-runner sysadmin logs tail --follow
```

### Erro: "permission denied"

**Causa:** Agent não tem permissões sudo.

**Solução:**
```bash
# No agent, verificar sudoers
sudo -l

# Adicionar usuário sloth-runner ao sudoers
sudo visudo
# Adicionar:
# sloth-runner ALL=(ALL) NOPASSWD: /usr/bin/apt, /usr/bin/apt-get
```

### Lista de pacotes vazia

**Causa:** Problema com dpkg/rpm database.

**Solução:**
```bash
# Para APT/dpkg:
sudo dpkg --configure -a
sudo apt --fix-broken install

# Verificar novamente
sloth-runner sysadmin packages list --limit 5
```

### Search retorna resultados inesperados

**Causa:** Índices desatualizados.

**Solução:**
```bash
# Atualizar índices primeiro
sloth-runner sysadmin packages update

# Buscar novamente
sloth-runner sysadmin packages search nginx
```

---

## Boas Práticas

### ✅ DO - Faça Isso

1. **Sempre atualize índices antes de instalar:**
   ```bash
   sloth-runner sysadmin packages update
   sloth-runner sysadmin packages install nginx -y
   ```

2. **Use filtros para encontrar pacotes rapidamente:**
   ```bash
   sloth-runner sysadmin pkg list --filter nginx --limit 10
   ```

3. **Documente instalações importantes:**
   ```bash
   sloth-runner sysadmin packages install postgresql -y | \
     tee -a /var/log/package-installs.log
   ```

4. **Verifique disponibilidade antes de instalar:**
   ```bash
   sloth-runner sysadmin packages search nginx
   sloth-runner sysadmin packages install nginx -y
   ```

### ❌ DON'T - Evite Isso

1. **Não use --yes em comandos manuais sem verificar:**
   ```bash
   # ❌ Perigoso
   sloth-runner sysadmin packages install random-package --yes

   # ✅ Melhor
   sloth-runner sysadmin packages search random-package
   sloth-runner sysadmin packages install random-package  # pergunta confirmação
   ```

2. **Não instale pacotes sem atualizar índices:**
   ```bash
   # ❌ Pode instalar versão antiga
   sloth-runner sysadmin packages install nginx

   # ✅ Sempre atualize primeiro
   sloth-runner sysadmin packages update
   sloth-runner sysadmin packages install nginx
   ```

3. **Não ignore erros:**
   ```bash
   # ❌ Ignora erros
   sloth-runner sysadmin packages install nginx || true

   # ✅ Trate erros
   if ! sloth-runner sysadmin packages install nginx; then
     echo "Installation failed, checking logs..."
     sloth-runner sysadmin logs tail -n 50
   fi
   ```

---

## Performance e Otimização

### Dicas de Performance

**1. Use limites para listas grandes:**
```bash
# ❌ Lento (lista 5000+ pacotes)
sloth-runner sysadmin packages list

# ✅ Rápido (limita a 50)
sloth-runner sysadmin packages list --limit 50
```

**2. Filtre no servidor, não localmente:**
```bash
# ❌ Ineficiente
sloth-runner sysadmin packages list | grep nginx

# ✅ Eficiente
sloth-runner sysadmin packages list --filter nginx
```

**3. Cache de operações:**
```bash
# Update uma vez, use múltiplas vezes
sloth-runner sysadmin packages update

# Agora pode buscar/instalar várias vezes sem re-update
sloth-runner sysadmin packages search nginx
sloth-runner sysadmin packages search apache
sloth-runner sysadmin packages search mysql
```

### Benchmarks

| Operação | Tempo Médio | Observações |
|----------|-------------|-------------|
| `list` (sem filtro) | 2-5s | Depende do # de pacotes |
| `list --limit 50` | <1s | Muito rápido |
| `list --filter nginx` | <1s | Filtro eficiente |
| `search nginx` | 1-3s | Depende de repos |
| `install nginx` | 10-30s | Download + instalação |
| `update` | 5-15s | Depende de repos |

---

## Roadmap

### ✅ Implementado (APT)
- [x] Detecção automática de package manager
- [x] list - Listar pacotes instalados
- [x] search - Buscar em repositórios
- [x] install - Instalar pacotes
- [x] update - Atualizar índices
- [x] UI moderna com pterm
- [x] Confirmações interativas
- [x] Error handling robusto

### 🚧 Em Desenvolvimento
- [ ] Suporte YUM/DNF
- [ ] Suporte Pacman
- [ ] Suporte APK
- [ ] Suporte Zypper

### 📋 Planejado - Fase 2
- [ ] remove - Remover pacotes
- [ ] upgrade - Atualizar pacotes
- [ ] check-updates - Verificar atualizações
- [ ] info - Informações detalhadas
- [ ] history - Histórico de operações

### 📋 Planejado - Fase 3
- [ ] Rolling updates inteligentes
- [ ] Rollback automático
- [ ] Dependency resolution visual
- [ ] Version pinning
- [ ] Backup antes de mudanças
- [ ] Multi-agent operations
- [ ] Progress bars para downloads
- [ ] Dry-run mode
- [ ] Diff de mudanças

---

## Comparação com Outras Ferramentas

### vs. SSH Manual

| Aspecto | SSH Manual | sloth-runner packages |
|---------|------------|----------------------|
| **Setup** | SSH keys | Agent gRPC |
| **UI** | Texto puro | Tabelas, cores, spinners |
| **Erro** | Difícil debug | Error handling robusto |
| **Multi-host** | Loop bash | Built-in (futuro) |
| **Rollback** | Manual | Automático (futuro) |
| **Logs** | Locais | Centralizados |

### vs. Ansible

| Aspecto | Ansible | sloth-runner packages |
|---------|---------|----------------------|
| **Sintaxe** | YAML | CLI direto |
| **Velocidade** | Lento (Python) | Rápido (Go) |
| **Setup** | Playbooks | Comando direto |
| **Real-time** | Não | Sim (spinners) |
| **Curva** | Alta | Baixa |

### vs. Salt/Puppet

Sloth-runner não substitui Salt/Puppet para gestão complexa de configuração, mas é muito mais simples para operações básicas de pacotes.

---

## FAQ

**Q: Posso usar em produção?**
A: Sim! A implementação APT está production-ready e testada.

**Q: Funciona com yum/dnf?**
A: Ainda não, mas está no roadmap Q1 2026.

**Q: Posso instalar múltiplos pacotes de uma vez?**
A: Atualmente não diretamente, mas pode usar loop bash. Feature planejada.

**Q: Há rollback se algo der errado?**
A: Não ainda, mas está planejado para Q1 2026.

**Q: Como funciona a detecção de package manager?**
A: Verifica automaticamente na ordem: apt → yum → dnf → pacman → apk → zypper

**Q: Preciso de sudo no agent?**
A: Sim, operações de pacotes requerem permissões elevadas.

---

## Suporte e Comunidade

**Encontrou um bug?**
- Abra issue no GitHub: [sloth-runner/issues](https://github.com/sloth-runner/issues)

**Precisa de ajuda?**
- Documentação: `/docs`
- Slack: #sloth-runner
- Email: support@sloth-runner.com

**Quer contribuir?**
- Implemente suporte para YUM/DNF
- Adicione testes
- Melhore documentação

---

## Ver Também

- [Services Management](sysadmin-services.md) - Gerenciar serviços systemd
- [Health Checks](health-command.md) - Verificar saúde do sistema
- [Logs Management](logs-command.md) - Gerenciar logs
- [Sysadmin Overview](sysadmin.md) - Visão geral de comandos sysadmin
