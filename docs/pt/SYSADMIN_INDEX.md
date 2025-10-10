# Índice da Documentação Sysadmin

## 📚 Guia Completo de Documentação

Bem-vindo à documentação completa dos comandos sysadmin do sloth-runner. Esta página serve como índice navegável para toda a documentação disponível.

---

## 🎯 Começando

### Para Iniciantes
1. **[Visão Geral Sysadmin](sysadmin.md)** - Comece aqui para entender todos os comandos disponíveis
2. **[Roadmap de Features](../SYSADMIN_FEATURES_ROADMAP.md)** - Veja o que está implementado e o que vem por aí

### Para Usuários Avançados
1. **[Package Management - Guia Completo](sysadmin-packages.md)** - Documentação detalhada de packages
2. **[Service Management - Guia Completo](sysadmin-services.md)** - Documentação detalhada de services

---

## 📖 Documentação Por Comando

### ✅ Comandos Implementados (Production-Ready)

| Comando | Status | Documentação | Nível |
|---------|--------|--------------|-------|
| **logs** | ✅ Funcional | [Logs Command](logs-command.md) | Detalhada |
| **health** | ✅ Funcional | [Health Command](health-command.md) | Detalhada |
| **debug** | ✅ Funcional | [Debug Command](debug-command.md) | Detalhada |
| **packages** | ✅ Funcional | **[📦 Packages - Guia Completo](sysadmin-packages.md)** | **Muito Detalhada** |
| **services** | ✅ Funcional | **[🔧 Services - Guia Completo](sysadmin-services.md)** | **Muito Detalhada** |

### 🔨 Comandos com CLI Pronto (Implementação Pendente)

| Comando | Status | Documentação | Timeline |
|---------|--------|--------------|----------|
| **backup** | 🔨 CLI Pronto | [Sysadmin Overview](sysadmin.md#backup) | Q1 2026 |
| **config** | 🔨 CLI Pronto | [Sysadmin Overview](sysadmin.md#config) | Q1 2026 |
| **deployment** | 🔨 CLI Pronto | [Sysadmin Overview](sysadmin.md#deployment) | Q1 2026 |
| **maintenance** | 🔨 CLI Pronto | [Sysadmin Overview](sysadmin.md#maintenance) | Q1 2026 |
| **network** | 🔨 CLI Pronto | [Sysadmin Overview](sysadmin.md#network) | Q4 2025 |
| **performance** | 🔨 CLI Pronto | [Sysadmin Overview](sysadmin.md#performance) | Q4 2025 |
| **resources** | 🔨 CLI Pronto | [Sysadmin Overview](sysadmin.md#resources) | Q4 2025 |
| **security** | 🔨 CLI Pronto | [Sysadmin Overview](sysadmin.md#security) | Q1 2026 |

---

## 🎓 Tutoriais e Guias

### Guias de Início Rápido

**5 Minutos para Começar:**
1. [Quick Start - Packages](sysadmin-packages.md#instalação-e-requisitos)
2. [Quick Start - Services](sysadmin-services.md#instalação-e-requisitos)

### Workflows Comuns

**Package Management:**
- [Setup Inicial de Novo Server](sysadmin-packages.md#1-setup-inicial-de-novo-server)
- [Auditoria de Compliance](sysadmin-packages.md#2-auditoria-de-compliance)
- [Troubleshooting de Dependências](sysadmin-packages.md#3-troubleshooting-de-dependências)
- [Manutenção Regular](sysadmin-packages.md#4-manutenção-regular)

**Service Management:**
- [Deploy de Aplicação](sysadmin-services.md#1-deploy-de-aplicação)
- [Troubleshooting de Serviço](sysadmin-services.md#2-troubleshooting-de-serviço)
- [Setup de Novo Server](sysadmin-services.md#3-setup-de-novo-server)
- [Manutenção Programada](sysadmin-services.md#4-manutenção-programada)

### Integrações

- [Packages + Services](sysadmin-services.md#com-packages)
- [Services + Health Checks](sysadmin-services.md#com-health-checks)
- [Services + Logs](sysadmin-services.md#com-logs)

---

## 🔍 Referência Rápida

### Comandos Mais Usados

**Packages:**
```bash
# Top 5 comandos packages
sloth-runner sysadmin packages list --filter nginx
sloth-runner sysadmin packages search nginx
sloth-runner sysadmin packages install nginx -y
sloth-runner sysadmin packages update
sloth-runner sysadmin pkg list -l 20
```

**Services:**
```bash
# Top 5 comandos services
sloth-runner sysadmin services list --status active
sloth-runner sysadmin services status nginx
sloth-runner sysadmin services restart nginx --verify
sloth-runner sysadmin services enable nginx
sloth-runner sysadmin svc logs nginx -n 50
```

### Atalhos e Aliases

```bash
# Packages aliases
pkg = packages
-y = --yes
-f = --filter
-l = --limit

# Services aliases
svc = services
-f = --filter
-s = --status
-n = --lines
```

---

## 🛠️ Troubleshooting

### Problemas Comuns

**Packages:**
- [Erro: "no supported package manager found"](sysadmin-packages.md#erro-no-supported-package-manager-found)
- [Erro: "permission denied"](sysadmin-packages.md#erro-permission-denied)
- [Lista de pacotes vazia](sysadmin-packages.md#lista-de-pacotes-vazia)
- [Search retorna resultados inesperados](sysadmin-packages.md#search-retorna-resultados-inesperados)

**Services:**
- [Erro: "no supported service manager found"](sysadmin-services.md#erro-no-supported-service-manager-found)
- [Serviço não inicia](sysadmin-services.md#serviço-não-inicia)
- [Erro: "Unit not found"](sysadmin-services.md#erro-unit-not-found)
- [Restart não aplica mudanças](sysadmin-services.md#restart-não-aplica-mudanças)

### Debug Avançado

1. **[Health Checks](health-command.md)** - Verificar saúde do sistema
2. **[Debug Command](debug-command.md)** - Troubleshooting avançado
3. **[Logs Command](logs-command.md)** - Análise de logs

---

## 📊 Comparações

### Packages vs Outras Ferramentas
- [vs SSH Manual](sysadmin-packages.md#vs-ssh-manual)
- [vs Ansible](sysadmin-packages.md#vs-ansible)
- [vs Salt/Puppet](sysadmin-packages.md#vs-saltpuppet)

### Services vs Outras Ferramentas
- [vs systemctl direto](sysadmin-services.md#vs-systemctl-direto)
- [vs Ansible service module](sysadmin-services.md#vs-ansible-service-module)

---

## 🎯 Por Caso de Uso

### Operações Diárias (SRE/DevOps)
1. [Monitoramento Diário](sysadmin.md#1-monitoramento-diário)
2. [Troubleshooting de Problema](sysadmin.md#2-troubleshooting-de-problema)
3. [Manutenção e Arquivamento](sysadmin.md#3-manutenção-e-arquivamento)

### Automação
1. [Scripts de Manutenção](sysadmin-packages.md#4-manutenção-regular)
2. [Deploy Automático](sysadmin-services.md#1-deploy-de-aplicação)
3. [Workflows de Automação](sysadmin.md#workflows-de-automação)

### Compliance e Segurança
1. [Auditoria de Compliance](sysadmin-packages.md#2-auditoria-de-compliance)
2. [Security Auditing](sysadmin.md#security)
3. [Inventory Management](sysadmin.md#ver-também)

---

## 📈 Performance e Otimização

### Dicas de Performance
- [Packages - Performance](sysadmin-packages.md#performance-e-otimização)
- [Services - Performance](sysadmin-services.md#performance-e-otimização)

### Benchmarks
- [Packages Benchmarks](sysadmin-packages.md#benchmarks)
- [Services Benchmarks](sysadmin-services.md#benchmarks)

---

## 🗺️ Roadmap

### Implementado (Atual)
- ✅ **5 comandos funcionais:** logs, health, debug, packages, services
- ✅ **Suporte APT completo** para packages
- ✅ **Suporte systemd completo** para services
- ✅ **85% test coverage** nos novos comandos
- ✅ **Documentação completa** em PT-BR e EN

### Próximos Passos (Q4 2025)
- 🚧 **Packages:** YUM, DNF, Pacman support
- 🚧 **Services:** init.d, OpenRC support
- 🚧 **Resources:** CPU, memory, disk monitoring
- 🚧 **Network:** Diagnostics tools
- 🚧 **Performance:** Real-time monitoring

### Futuro (2026)
- 📋 **Config management**
- 📋 **Backup & restore**
- 📋 **Deployment automation**
- 📋 **Security auditing**
- 📋 **Container management**

Ver **[Roadmap Completo](../SYSADMIN_FEATURES_ROADMAP.md)** para detalhes.

---

## 🌍 Idiomas Disponíveis

| Idioma | Disponibilidade | Link |
|--------|----------------|------|
| **Português (PT-BR)** | ✅ Completo | Você está aqui |
| **English (EN)** | ✅ Completo | [English Docs](../en/sysadmin-new-tools.md) |
| **中文 (ZH)** | 📋 Planejado | - |

---

## 💡 Boas Práticas

### Packages
- [DO's and DON'Ts](sysadmin-packages.md#boas-práticas)
- [Performance Tips](sysadmin-packages.md#performance-e-otimização)

### Services
- [DO's and DON'Ts](sysadmin-services.md#boas-práticas)
- [Performance Tips](sysadmin-services.md#performance-e-otimização)

---

## 🤝 Contribuindo

**Quer ajudar a melhorar a documentação?**

1. **Reportar erros:** Abra issue no GitHub
2. **Sugerir melhorias:** Pull requests bem-vindos
3. **Adicionar exemplos:** Compartilhe seus workflows
4. **Traduzir:** Ajude em outras línguas

**Áreas que precisam de ajuda:**
- [ ] Mais exemplos práticos
- [ ] Casos de uso enterprise
- [ ] Troubleshooting scenarios
- [ ] Performance benchmarks
- [ ] Integration guides

---

## 📞 Suporte

**Precisa de ajuda?**
- 📖 **Docs:** Você está nelas!
- 💬 **Slack:** #sloth-runner
- 🐛 **Issues:** [GitHub Issues](https://github.com/sloth-runner/issues)
- 📧 **Email:** support@sloth-runner.com

**Antes de abrir issue:**
1. ✅ Leia a documentação relevante
2. ✅ Verifique troubleshooting
3. ✅ Pesquise issues existentes
4. ✅ Prepare logs e contexto

---

## 📚 Recursos Adicionais

### Documentação Externa
- [systemd Documentation](https://www.freedesktop.org/wiki/Software/systemd/)
- [APT User Manual](https://www.debian.org/doc/manuals/apt-guide/)
- [Best Practices for SRE](https://sre.google/books/)

### Comunidade
- [Sloth-Runner Blog](https://blog.sloth-runner.com)
- [YouTube Tutorials](https://youtube.com/@sloth-runner)
- [Twitter @slothrunner](https://twitter.com/slothrunner)

### Exemplos e Templates
- [GitHub Examples](https://github.com/sloth-runner/examples)
- [Workflow Templates](https://github.com/sloth-runner/templates)
- [Community Scripts](https://github.com/sloth-runner/community)

---

## 🎉 Começe Agora!

**Novos usuários? Comece aqui:**

1. 📖 Leia a **[Visão Geral](sysadmin.md)**
2. 🚀 Siga o **[Quick Start - Packages](sysadmin-packages.md#verificação-rápida)**
3. 🔧 Experimente **[Quick Start - Services](sysadmin-services.md#verificação-rápida)**
4. 💪 Explore os **[Workflows Comuns](#workflows-comuns)**

**Usuários experientes? Vá direto para:**

- 📦 **[Package Management Completo](sysadmin-packages.md)**
- 🔧 **[Service Management Completo](sysadmin-services.md)**
- 🗺️ **[Roadmap de Features](../SYSADMIN_FEATURES_ROADMAP.md)**

---

**Última atualização:** Outubro 2025
**Versão da documentação:** 2.0
**Status:** ✅ Produção

