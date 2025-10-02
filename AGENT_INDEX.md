# 📚 Agent Setup - Índice de Documentação

Este é o índice central para toda a documentação sobre configuração e uso de agents no Sloth Runner.

## 🗂️ Estrutura da Documentação

### 📖 Por Nível de Detalhe

1. **Cheat Sheet** (1-2 min) → `AGENT_CHEATSHEET.md`
   - Comandos essenciais
   - Sintaxe rápida
   - Referência de bolso

2. **Quick Start** (5 min) → `AGENT_QUICK_START.md`
   - Guia rápido TL;DR
   - Passo a passo básico
   - Troubleshooting essencial

3. **Guia Completo** (15-20 min) → `docs/agent-setup.md`
   - Documentação detalhada
   - Segurança e monitoramento
   - Casos de uso avançados

4. **Exemplos Práticos** → `examples/distributed_execution.sloth`
   - Código executável
   - 6 exemplos completos
   - Base para seus scripts

### 🔧 Por Tipo de Recurso

#### Documentação
- `docs/agent-setup.md` - Guia completo
- `AGENT_QUICK_START.md` - Referência rápida
- `AGENT_CHEATSHEET.md` - Cheat sheet
- `AGENT_SETUP_SUMMARY.txt` - Resumo visual
- Este arquivo (`AGENT_INDEX.md`) - Índice

#### Scripts Executáveis
- `start_master.sh` - Iniciar master server
- `start_local_agent.sh` - Iniciar agent local
- `manage_remote_agent.sh` - Gerenciar agents remotos via SSH

#### Exemplos
- `examples/distributed_execution.sloth` - 6 exemplos práticos

#### Atualizações
- `README.md` - Seção "Distributed Task Execution"
- `mkdocs.yml` - Integração na documentação web

## 🎯 Guia de Uso Rápido

### Sua Primeira Vez?

```bash
# 1. Leia o Quick Start
cat AGENT_QUICK_START.md

# 2. Inicie o master
./start_master.sh

# 3. Inicie um agent local
./start_local_agent.sh test-agent

# 4. Verifique
./sloth-runner agent list --master 192.168.1.29:50053

# 5. Teste com exemplo
./sloth-runner run -f examples/distributed_execution.sloth
```

### Já Conhece o Básico?

```bash
# Ver comandos rápidos
cat AGENT_CHEATSHEET.md

# Gerenciar agents remotos
./manage_remote_agent.sh start user@host agent-name ip-address
```

### Quer Aprofundar?

```bash
# Ler documentação completa
cat docs/agent-setup.md

# Estudar exemplos
cat examples/distributed_execution.sloth
```

## 📊 Matriz de Escolha

| Situação | Leia | Tempo |
|----------|------|-------|
| Nunca usei agents | `AGENT_QUICK_START.md` | 5 min |
| Preciso de um comando específico | `AGENT_CHEATSHEET.md` | 1 min |
| Quero entender tudo | `docs/agent-setup.md` | 20 min |
| Quero ver código real | `examples/distributed_execution.sloth` | 10 min |
| Quero uma visão geral | `AGENT_SETUP_SUMMARY.txt` | 5 min |

## 🔗 Links Diretos

### Documentação
- [Guia Completo](./docs/agent-setup.md)
- [Quick Start](./AGENT_QUICK_START.md)
- [Cheat Sheet](./AGENT_CHEATSHEET.md)
- [Resumo Visual](./AGENT_SETUP_SUMMARY.txt)

### Scripts
- [start_master.sh](./start_master.sh)
- [start_local_agent.sh](./start_local_agent.sh)
- [manage_remote_agent.sh](./manage_remote_agent.sh)

### Exemplos
- [Distributed Execution](./examples/distributed_execution.sloth)

### Integração
- [README - Distributed](./README.md#distributed-task-execution)
- [MkDocs Config](./mkdocs.yml)

## 🎨 Fluxo de Aprendizado Recomendado

```
1. AGENT_QUICK_START.md
   ↓
2. ./start_master.sh
   ↓
3. ./start_local_agent.sh test
   ↓
4. examples/distributed_execution.sloth
   ↓
5. docs/agent-setup.md
   ↓
6. AGENT_CHEATSHEET.md (para referência futura)
```

## 💡 Casos de Uso Comuns

### Deploy em Múltiplos Servidores
**Leia:** `examples/distributed_execution.sloth` (Exemplo 3)  
**Script:** Use `:delegate_to()` em loop

### Monitoramento Distribuído
**Leia:** `examples/distributed_execution.sloth` (Exemplo 4)  
**Script:** Health check em múltiplos agents

### Backup Remoto
**Leia:** `examples/distributed_execution.sloth` (Exemplo 5)  
**Script:** Backup de configurações

### Gerenciamento de Serviços
**Leia:** `docs/agent-setup.md` (Seção "Using Agents")  
**Módulos:** `systemd`, `pkg`

## 🆘 Troubleshooting Rápido

### Agent não conecta?
**Ver:** `AGENT_CHEATSHEET.md` - Seção Troubleshooting  
**Ver:** `docs/agent-setup.md` - Seção Troubleshooting completa

### Como ver logs?
```bash
# Rodar agent sem daemon
./sloth-runner agent start --name test \
  --master 192.168.1.29:50053 \
  --bind-address <IP>
```

### Como parar tudo?
```bash
# Parar todos os processos
pkill -f sloth-runner
```

## 📈 Evolução do Aprendizado

### Nível 1: Iniciante
- ✅ Ler `AGENT_QUICK_START.md`
- ✅ Executar `start_master.sh`
- ✅ Executar `start_local_agent.sh`
- ✅ Testar conectividade

### Nível 2: Intermediário
- ✅ Configurar agents remotos
- ✅ Executar exemplos do `distributed_execution.sloth`
- ✅ Criar scripts simples com `:delegate_to()`
- ✅ Usar módulos `systemd` e `pkg`

### Nível 3: Avançado
- ✅ Ler `docs/agent-setup.md` completo
- ✅ Implementar workflows complexos
- ✅ Configurar segurança (firewall, SSH keys)
- ✅ Monitoramento e health checks
- ✅ Deploy em produção

## 🎯 Objetivos de Aprendizado

### Após ler `AGENT_QUICK_START.md`, você saberá:
- ✅ Como iniciar master e agents
- ✅ Comandos básicos
- ✅ Como verificar conectividade

### Após ler `docs/agent-setup.md`, você saberá:
- ✅ Arquitetura master-agent
- ✅ Configuração de segurança
- ✅ Troubleshooting avançado
- ✅ Monitoramento e saúde

### Após estudar `examples/distributed_execution.sloth`, você saberá:
- ✅ Como usar `:delegate_to()` na prática
- ✅ Deploy distribuído
- ✅ Health checks
- ✅ Workflows complexos

## 📞 Suporte e Recursos

### Documentação Online
- Site: https://chalkan3-sloth.github.io/sloth-runner/
- Agent Setup: https://chalkan3-sloth.github.io/sloth-runner/agent-setup/

### Repositório
- GitHub: https://github.com/chalkan3-sloth/sloth-runner
- Issues: https://github.com/chalkan3-sloth/sloth-runner/issues

### Comunidade
- Discord: (se disponível)
- Discussions: GitHub Discussions

## 🔄 Atualizações

Este índice foi criado em: Outubro 2024  
Versão: 1.0.0  
Compatível com: Sloth Runner v3.23.1+

---

**Comece agora:** `cat AGENT_QUICK_START.md` ou `./start_master.sh`

**Dúvidas?** Consulte `docs/agent-setup.md` para documentação completa.
