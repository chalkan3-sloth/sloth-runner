# 🧪 Bootstrap Quick Test Guide

Este guia mostra como testar rapidamente o bootstrap.sh em diferentes cenários.

## Teste 1: Help e Validação

```bash
# Ver todas as opções
./bootstrap.sh --help

# Teste sem parâmetro obrigatório (deve falhar graciosamente)
./bootstrap.sh
# Esperado: erro pedindo --name
```

## Teste 2: Instalação Local (Sem Systemd)

```bash
# Instalar sem systemd para teste local
./bootstrap.sh \
  --name test-agent \
  --no-sudo \
  --no-systemd

# Verificar instalação
ls -la ~/.local/bin/sloth-runner
~/.local/bin/sloth-runner --version
```

## Teste 3: Instalação via Vagrant

```bash
# Assumindo que você tem um Vagrant rodando
cd /Users/chalkan3/.projects/vagrant/archlinux

# Instalar agent no Vagrant
vagrant ssh -c "bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \
  --name mariaguica \
  --master 192.168.1.29:50053 \
  --bind-address \$(hostname -I | awk '{print \$1}')"

# Verificar status
vagrant ssh -c "sudo systemctl status sloth-runner-agent"
```

## Teste 4: Verificar no Master

```bash
# No host do master
sloth-runner master start

# Em outro terminal, listar agents
sloth-runner agent list

# Testar comando no agent
sloth-runner agent run mariaguica "hostname"
sloth-runner agent run mariaguica "uname -a"
```

## Teste 5: Verificar Logs

```bash
# No host do agent
sudo journalctl -u sloth-runner-agent -f

# Ou via SSH
ssh user@agent-host "sudo journalctl -u sloth-runner-agent -n 50"
```

## Teste 6: Teste de Reconexão

```bash
# 1. Parar o master
# No terminal do master: Ctrl+C

# 2. Verificar que agent tenta reconectar
vagrant ssh -c "sudo journalctl -u sloth-runner-agent -f"
# Você deve ver tentativas de reconexão

# 3. Reiniciar master
sloth-runner master start

# 4. Verificar que agent reconectou
sloth-runner agent list
# Agent deve aparecer como online
```

## Teste 7: Múltiplos Agents

```bash
# Agent 1
./bootstrap.sh \
  --name agent-01 \
  --port 50051 \
  --no-sudo \
  --no-systemd

# Agent 2  
./bootstrap.sh \
  --name agent-02 \
  --port 50052 \
  --no-sudo \
  --no-systemd

# Verificar no master
sloth-runner agent list
```

## Teste 8: Limpeza

```bash
# Parar e remover serviço
sudo systemctl stop sloth-runner-agent
sudo systemctl disable sloth-runner-agent
sudo rm /etc/systemd/system/sloth-runner-agent.service
sudo systemctl daemon-reload

# Remover do master
sloth-runner agent delete test-agent --yes

# Remover binário
rm ~/.local/bin/sloth-runner
# ou
sudo rm /usr/local/bin/sloth-runner
```

## Troubleshooting Rápido

### Agent não registra no master

```bash
# 1. Verificar conectividade
telnet MASTER_IP 50053

# 2. Verificar logs do agent
sudo journalctl -u sloth-runner-agent -n 100

# 3. Verificar se agent está rodando
ps aux | grep sloth-runner

# 4. Verificar porta
sudo netstat -tulpn | grep 50051
```

### Erro de permissão

```bash
# Verificar usuário do serviço
sudo systemctl cat sloth-runner-agent | grep User

# Verificar permissões do binário
ls -la /usr/local/bin/sloth-runner

# Dar permissão de execução
sudo chmod +x /usr/local/bin/sloth-runner
```

### Service failed to start

```bash
# Ver erro detalhado
sudo systemctl status sloth-runner-agent -l

# Ver últimos logs
sudo journalctl -u sloth-runner-agent -n 50 --no-pager

# Testar comando manualmente
/usr/local/bin/sloth-runner agent start \
  --name test-agent \
  --master localhost:50053 \
  --port 50051
```

## Checklist de Validação

- [ ] Help funciona (`--help`)
- [ ] Detecta erro de parâmetro faltante
- [ ] Instala binário corretamente
- [ ] Cria serviço systemd
- [ ] Serviço inicia automaticamente
- [ ] Agent aparece no master
- [ ] Pode executar comandos no agent
- [ ] Reconecta após master restart
- [ ] Logs aparecem no journald
- [ ] Pode deletar agent

## Resultados Esperados

✅ **Sucesso**: Agent instalado, rodando e registrado no master  
✅ **Logs**: Limpos e informativos no journald  
✅ **Reconexão**: Automática após falhas  
✅ **Performance**: Responde rápido aos comandos  
✅ **Limpeza**: Fácil de remover completamente  

---

**Nota**: Este guia assume que você tem acesso root/sudo no host do agent.
