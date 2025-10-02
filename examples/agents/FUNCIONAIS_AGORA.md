# 🚀 EXEMPLOS FUNCIONAIS DE DELEGATE_TO

## ✅ O que funciona AGORA:

### 1. Comando direto com agent run
```bash
# Executa comando simples no agent remoto
sloth-runner agent run ladyguica "ls -la $HOME" --master 192.168.1.29:50053
sloth-runner agent run keiteguica "hostname && date" --master 192.168.1.29:50053
```

### 2. Script para executar LS em múltiplos agents
```bash
#!/bin/bash
# Executa LS em todos os agents

echo "🔍 Executando LS em todos os agents..."

echo "📍 ladyguica:"
sloth-runner agent run ladyguica "ls -la \$HOME" --master 192.168.1.29:50053

echo ""
echo "📍 keiteguica:"
sloth-runner agent run keiteguica "ls -la \$HOME" --master 192.168.1.29:50053

echo ""
echo "✅ LS executado em todos os agents!"
```

## 🚧 O que NÃO funciona ainda:

### Arquivo .sloth com delegate_to
Os arquivos `.sloth` com `delegate_to` estão com problema no agent. O agent recebe o script mas falha ao executar as tasks remotamente.

**Error típico:**
```
task failed on agent: one or more task groups failed
```

**Causa:** O agent tenta executar o workflow inteiro em vez de apenas a task específica.

## 📝 Exemplo de uso funcional:

### Para executar LS em múltiplos hosts:
```bash
# Criar script executável
cat > run_ls_agents.sh << 'EOF'
#!/bin/bash
echo "🔍 Listando arquivos nos agents remotos..."

echo "📍 Agent: ladyguica"
sloth-runner agent run ladyguica "hostname && ls -la \$HOME | head -10" --master 192.168.1.29:50053

echo ""
echo "📍 Agent: keiteguica"  
sloth-runner agent run keiteguica "hostname && ls -la \$HOME | head -10" --master 192.168.1.29:50053

echo ""
echo "✅ Concluído!"
EOF

chmod +x run_ls_agents.sh
./run_ls_agents.sh
```

### Para comandos mais complexos:
```bash
# Info do sistema em todos os agents
sloth-runner agent run ladyguica "uname -a && uptime && df -h /" --master 192.168.1.29:50053
sloth-runner agent run keiteguica "uname -a && uptime && df -h /" --master 192.168.1.29:50053
```

## 🔧 Status do desenvolvimento:

- ✅ **Agent run**: Funciona perfeitamente 
- ✅ **Master/Agent connection**: Funciona
- ✅ **Agent name resolution**: Funciona
- 🚧 **Delegate_to em .sloth**: Em desenvolvimento
- 🚧 **Remote workflow execution**: Em desenvolvimento

## 🎯 Recomendação atual:

Use `sloth-runner agent run` para execução remota até que o `delegate_to` seja corrigido.