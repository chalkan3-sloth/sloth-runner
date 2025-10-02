# Como Iniciar Agente no Vagrant com Bootstrap

## TL;DR - Comando Rápido

```bash
cd /Users/chalkan3/.projects/vagrant/archlinux

vagrant ssh -c "curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh | sudo bash -s -- \
  --name mariaguica \
  --master 192.168.1.29:50053 \
  --port 50051 \
  --bind-address 172.17.0.2 \
  --no-systemd"
```

## Explicação dos Parâmetros

| Parâmetro | Valor | Descrição |
|-----------|-------|-----------|
| `--name` | `mariaguica` | Nome do agente (deve ser único) |
| `--master` | `192.168.1.29:50053` | Endereço do master (IP do Mac host) |
| `--port` | `50051` | Porta que o agente vai escutar |
| `--bind-address` | `172.17.0.2` | IP do Vagrant para bind |
| `--no-systemd` | - | Força modo direto (systemd não funciona bem no container) |

## Como Descobrir os IPs

### IP do Master (Mac Host)

```bash
ifconfig | grep "inet " | grep -v 127.0.0.1 | head -1
# Resultado: inet 192.168.1.29
```

### IP do Vagrant

```bash
cd /Users/chalkan3/.projects/vagrant/archlinux
vagrant ssh -c "ip addr show | grep 'inet ' | grep -v 127.0.0.1"
# Resultado: inet 172.17.0.2/16
```

## Verificar Se Está Funcionando

### 1. No Vagrant - Ver se o processo está rodando

```bash
vagrant ssh -c "ps aux | grep sloth-runner | grep -v grep"
```

Saída esperada:
```
root  4887  3.2  0.3 1721556 27876 ?  Sl  11:42  0:00 /usr/local/bin/sloth-runner agent start...
```

### 2. No Master - Listar agentes

```bash
sloth-runner agent list
```

Saída esperada:
```
AGENT NAME     ADDRESS            STATUS    LAST HEARTBEAT
mariaguica     172.17.0.2:50051   Active    2025-10-02T08:42:43-03:00
```

### 3. Testar Execução de Comando

```bash
sloth-runner agent run mariaguica "hostname && uname -a"
```

## Gerenciar o Agente

### Parar o Agente

```bash
vagrant ssh -c "sudo pkill -f 'sloth-runner agent'"
```

### Ver Logs

```bash
vagrant ssh -c "cat /tmp/agent.log"
```

### Reiniciar o Agente

```bash
vagrant ssh -c "sudo pkill -f 'sloth-runner agent' && sleep 2 && \
  sudo /usr/local/bin/sloth-runner agent start \
  --name mariaguica \
  --master 192.168.1.29:50053 \
  --port 50051 \
  --bind-address 172.17.0.2 \
  --daemon"
```

## Troubleshooting

### Agente não aparece na lista

1. Verifique se o master está rodando:
```bash
ps aux | grep "sloth-runner master" | grep -v grep
```

2. Se não estiver, inicie o master:
```bash
sloth-runner master start &
```

3. Reinicie o agente no Vagrant

### Conectividade

Teste se o Vagrant consegue alcançar o master:

```bash
vagrant ssh -c "nc -zv 192.168.1.29 50053"
```

Saída esperada:
```
Connection to 192.168.1.29 50053 port [tcp/*] succeeded!
```

### Porta já em uso

Se a porta 50051 já estiver em uso, use outra porta:

```bash
vagrant ssh -c "curl -fsSL https://raw.githubusercontent.com/.../bootstrap.sh | sudo bash -s -- \
  --name mariaguica \
  --master 192.168.1.29:50053 \
  --port 50052 \
  --bind-address 172.17.0.2 \
  --no-systemd"
```

## Por Que --no-systemd?

O Vagrant com Docker provider roda um container que tem systemd limitado. Quando tentamos usar o systemd normalmente, ele falha com:

```
A dependency job for sloth-runner-agent.service failed
```

Usando `--no-systemd`, o bootstrap.sh:
1. Detecta que systemd não está funcionando
2. Inicia o agente diretamente com `--daemon`
3. O agente roda em background sem precisar de systemd

## Múltiplos Agentes

Para rodar vários agentes no mesmo Vagrant, use portas diferentes:

```bash
# Agente 1
vagrant ssh -c "curl ... | sudo bash -s -- --name agent1 --port 50051 ..."

# Agente 2
vagrant ssh -c "curl ... | sudo bash -s -- --name agent2 --port 50052 ..."
```

## Ver Todos os Agentes Ativos

```bash
$ sloth-runner agent list

AGENT NAME     ADDRESS              STATUS    LAST HEARTBEAT
------------   ----------           ------    --------------
ladyguica      192.168.1.16:50051   Active    2025-10-02T08:45:05-03:00
mariaguica     172.17.0.2:50051     Active    2025-10-02T08:42:43-03:00
```

## Deletar Agente do Master

```bash
# Com confirmação
sloth-runner agent delete mariaguica

# Sem confirmação
sloth-runner agent delete mariaguica --yes
```

## Links

- [Bootstrap Documentation](BOOTSTRAP.md)
- [Agent Documentation](docs/modules/agent.md)
- [Main README](README.md)
