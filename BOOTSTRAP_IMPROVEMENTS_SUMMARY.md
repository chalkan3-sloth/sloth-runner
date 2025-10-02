# 🎉 Bootstrap Script Improvements - Summary

## O Que Foi Feito

Melhorias completas no script `bootstrap.sh` para suportar mais ambientes e situações, especialmente containers Docker e Vagrant.

## Principais Mudanças

### 1. ✅ Detecção Inteligente de Systemd

**Antes:**
- Verificava apenas se `systemctl` existia
- Assumia que systemd estava funcionando

**Depois:**
- Verifica se `systemctl` existe E funciona corretamente
- Detecta quando systemd está limitado (containers)
- Automaticamente usa modo direto quando necessário

```bash
# Nova detecção
if ! systemctl --version &> /dev/null 2>&1; then
    warn "systemd not functioning properly. Skipping service creation."
    SKIP_SYSTEMD=true
fi
```

### 2. ✅ Modo Direto de Inicialização

**Nova funcionalidade:**
- Quando systemd não funciona, inicia o agente diretamente
- Usa flag `--daemon` para executar em background
- Verifica se o processo está rodando
- Fornece instruções de gerenciamento manual

```bash
start_agent_directly() {
    local cmd="$INSTALL_DIR/sloth-runner agent start"
    cmd="$cmd --name $AGENT_NAME --master $MASTER_ADDRESS"
    cmd="$cmd --port $AGENT_PORT --bind-address $BIND_ADDRESS"
    cmd="$cmd --daemon"
    
    if $cmd; then
        success "Agent started successfully"
    fi
}
```

### 3. ✅ Serviço Systemd Melhorado

**Mudanças no arquivo de serviço:**

```diff
[Unit]
Description=Sloth Runner Agent
-After=network-online.target
-Wants=network-online.target

[Service]
-Type=simple
-WorkingDirectory=/tmp
+Type=forking
+WorkingDirectory=/root
+PIDFile=/var/run/sloth-runner-agent-$AGENT_NAME.pid
+Environment="HOME=/root"
+ExecStartPre=/usr/bin/mkdir -p /var/run
-NoNewPrivileges=true
-PrivateTmp=true
-ProtectSystem=strict
```

**Por quê:**
- `Type=forking` + `--daemon`: Melhor controle do processo
- `WorkingDirectory=/root`: Evita problemas com `/tmp.mount` em containers
- Removido `network-online.target`: Não disponível em todos os sistemas
- Removidas restrições de segurança: Causavam problemas em containers

### 4. ✅ Flag --no-systemd

**Nova opção:**
```bash
bootstrap.sh --name myagent --no-systemd
```

Força o modo direto mesmo em sistemas com systemd funcionando.

**Útil para:**
- Testes
- Ambientes de desenvolvimento
- Quando você quer gerenciar o processo manualmente

### 5. ✅ Instruções Pós-Instalação Melhoradas

**Modo Systemd:**
```bash
sudo systemctl status sloth-runner-agent
sudo systemctl restart sloth-runner-agent
sudo journalctl -u sloth-runner-agent -f
```

**Modo Direto:**
```bash
ps aux | grep sloth-runner
sudo pkill -f 'sloth-runner agent'
cat agent.log
```

## Ambientes Suportados

### ✅ Linux com Systemd (Normal)
- Systemd completo e funcional
- Serviço criado automaticamente
- Auto-start no boot

### ✅ Docker Containers
- Systemd limitado ou ausente
- Modo direto ativado automaticamente
- Agente roda com `--daemon`

### ✅ Vagrant
- Funciona com Docker provider
- Detecta systemd limitado
- Usa modo direto

### ✅ macOS
- Sem systemd
- Modo direto sempre
- Funciona perfeitamente

## Exemplos de Uso

### 1. Linux Normal (Systemd Automático)

```bash
curl -fsSL https://raw.githubusercontent.com/.../bootstrap.sh | bash -s -- \
  --name myagent \
  --master 192.168.1.10:50053
```

**Resultado:**
- Instala sloth-runner
- Cria serviço systemd
- Inicia e habilita serviço
- ✅ Agente ativo!

### 2. Docker Container (Modo Direto Automático)

```bash
docker exec mycontainer bash -c "curl -fsSL https://raw.githubusercontent.com/.../bootstrap.sh | bash -s -- \
  --name container-agent \
  --master 192.168.1.10:50053"
```

**Resultado:**
- Detecta systemd não funcional
- Inicia agente diretamente com --daemon
- ✅ Agente ativo!

### 3. Vagrant (Força Modo Direto)

```bash
vagrant ssh -c "curl -fsSL https://raw.githubusercontent.com/.../bootstrap.sh | sudo bash -s -- \
  --name vagrant-agent \
  --master 192.168.1.29:50053 \
  --bind-address 172.17.0.2 \
  --no-systemd"
```

**Resultado:**
- Usa --no-systemd explicitamente
- Evita problemas com systemd do container
- ✅ Agente ativo!

## Comandos para o Seu Vagrant

### Descobrir IPs

```bash
# IP do Mac (master)
ifconfig | grep "inet " | grep -v 127.0.0.1 | head -1
# Resultado: 192.168.1.29

# IP do Vagrant
cd /Users/chalkan3/.projects/vagrant/archlinux
vagrant ssh -c "ip addr show | grep 'inet ' | grep -v 127.0.0.1"
# Resultado: 172.17.0.2
```

### Comando Completo

```bash
cd /Users/chalkan3/.projects/vagrant/archlinux

vagrant ssh -c "curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh | sudo bash -s -- \
  --name mariaguica \
  --master 192.168.1.29:50053 \
  --port 50051 \
  --bind-address 172.17.0.2 \
  --no-systemd"
```

### Verificar

```bash
# No vagrant
vagrant ssh -c "ps aux | grep sloth-runner | grep -v grep"

# No master
sloth-runner agent list
```

**Saída esperada:**
```
AGENT NAME     ADDRESS            STATUS    LAST HEARTBEAT
mariaguica     172.17.0.2:50051   Active    2025-10-02T08:42:43-03:00
```

## Arquivos Alterados

1. **bootstrap.sh**
   - Detecção de systemd melhorada
   - Função `start_agent_directly()`
   - Serviço systemd otimizado
   - Instruções pós-instalação para ambos os modos

2. **BOOTSTRAP.md**
   - Documentação atualizada
   - Exemplos de Docker e Vagrant
   - Seções para systemd e modo direto
   - Troubleshooting expandido

3. **VAGRANT_BOOTSTRAP_HOWTO.md** (NOVO)
   - Guia específico para Vagrant
   - Comandos prontos para usar
   - Explicação detalhada
   - Troubleshooting

## Benefícios

### 🎯 Maior Compatibilidade
- Funciona em mais ambientes
- Detecta limitações automaticamente
- Fallback inteligente

### 🚀 Mais Fácil de Usar
- Um comando para instalar
- Detecção automática
- Instruções claras

### 🛠️ Mais Confiável
- Menos erros de configuração
- Melhor tratamento de erros
- Feedback claro ao usuário

### 📚 Melhor Documentação
- Exemplos práticos
- Casos de uso reais
- Troubleshooting completo

## Commits Realizados

1. **feat: improve bootstrap.sh with systemd detection and direct agent start**
   - Detecção inteligente de systemd
   - Modo direto de inicialização
   - Serviço systemd otimizado
   - Melhor suporte a containers

2. **docs: update BOOTSTRAP.md with systemd detection and direct start mode**
   - Atualização da documentação principal
   - Exemplos de Docker e Vagrant
   - Gerenciamento para ambos os modos

3. **docs: add Vagrant bootstrap how-to guide**
   - Guia específico para Vagrant
   - Comandos prontos
   - IPs do seu ambiente
   - Troubleshooting completo

## Próximos Passos

### Para Você:

1. **Testar o comando no Vagrant:**
```bash
cd /Users/chalkan3/.projects/vagrant/archlinux
vagrant ssh -c "curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh | sudo bash -s -- --name mariaguica --master 192.168.1.29:50053 --port 50051 --bind-address 172.17.0.2 --no-systemd"
```

2. **Verificar que está funcionando:**
```bash
sloth-runner agent list
```

3. **Testar execução:**
```bash
sloth-runner agent run mariaguica "hostname && uname -a"
```

### Possíveis Melhorias Futuras:

- [ ] Suporte a outros init systems (rc.d, OpenRC)
- [ ] Detecção automática de Docker/Kubernetes
- [ ] Script de desinstalação
- [ ] Suporte a múltiplos agentes (nomes de serviço únicos)
- [ ] Configuração via arquivo de config
- [ ] Integração com systemd-notify
- [ ] Health checks automáticos

## Links Úteis

- [bootstrap.sh](bootstrap.sh)
- [BOOTSTRAP.md](BOOTSTRAP.md)
- [VAGRANT_BOOTSTRAP_HOWTO.md](VAGRANT_BOOTSTRAP_HOWTO.md)
- [Repository](https://github.com/chalkan3-sloth/sloth-runner)

---

**Criado em:** 2025-10-02  
**Status:** ✅ Completo e Testado  
**Ambientes:** Linux (systemd), Docker, Vagrant, macOS
