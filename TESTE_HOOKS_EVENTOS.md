# 🧪 Procedimento de Teste Completo - Hooks e Eventos

Este documento contém um procedimento passo a passo para testar todo o sistema de hooks, eventos e watchers do sloth-runner.

---

## 📋 Pré-requisitos

- ✅ sloth-runner instalado em `$HOME/.local/bin/`
- ✅ Agent rodando em `lady-guica` ou outro host
- ✅ Master server rodando localmente
- ✅ Web UI acessível em `http://localhost:8080`

---

## 🧹 Etapa 1: Limpeza dos Bancos de Dados

### 1.1 Parar o Master Server (se estiver rodando)

```bash
# Encontrar processo do master
ps aux | grep "sloth-runner.*master"

# Matar o processo (substitua <PID> pelo PID encontrado)
kill <PID>

# Ou usar pkill
pkill -f "sloth-runner.*master"
```

### 1.2 Limpar Bancos de Dados

```bash
# Navegar para o diretório do projeto
cd ~/.projects/task-runner

# Remover banco de eventos
rm -f ~/.sloth-runner/events.db
echo "✓ Banco de eventos removido"

# Remover banco de hooks
rm -f ~/.sloth-runner/hooks.db
echo "✓ Banco de hooks removido"

# Remover banco de métricas (opcional)
rm -f ~/.sloth-runner/metrics.db
echo "✓ Banco de métricas removido"

# Limpar logs de hooks
rm -f ~/.sloth-runner/hook-executions.log
echo "✓ Logs de hooks limpos"

# Verificar limpeza
ls -lh ~/.sloth-runner/
```

### 1.3 Recriar Diretório Limpo

```bash
# Garantir que o diretório existe
mkdir -p ~/.sloth-runner

# Verificar permissões
chmod 755 ~/.sloth-runner

echo "✓ Diretório limpo e pronto"
```

---

## 🚀 Etapa 2: Iniciar o Sistema

### 2.1 Iniciar Master Server com Web UI

```bash
# Iniciar master server em background
cd ~/.projects/task-runner
$HOME/.local/bin/sloth-runner master \
    --bind-address 0.0.0.0 \
    --port 50051 \
    --web-ui \
    --web-port 8080 > /tmp/sloth-master.log 2>&1 &

# Guardar PID
MASTER_PID=$!
echo "Master server iniciado com PID: $MASTER_PID"

# Aguardar inicialização
sleep 3

# Verificar se está rodando
curl -s http://localhost:8080 > /dev/null && echo "✓ Web UI acessível" || echo "✗ Web UI não respondeu"
```

### 2.2 Verificar Agent (deve já estar rodando)

```bash
# Verificar conectividade com agent
$HOME/.local/bin/sloth-runner agent list

# Deve mostrar lady-guica (ou seu agent)
```

### 2.3 Acessar Web UI

```bash
# Abrir Web UI no navegador
open http://localhost:8080

# Ou apenas mostrar URL
echo "📊 Web UI disponível em: http://localhost:8080"
```

---

## 🧪 Etapa 3: Testes de Hooks

### 3.1 Testar Hook Simples

```bash
# Criar arquivo de teste de hook
cat > /tmp/test_hook.sloth << 'SLOTH'
task("test_simple_hook", {
    on = function()
        shell("echo 'Hook simples executado'")
        shell("date")
        return true
    end
})
SLOTH

# Executar
$HOME/.local/bin/sloth-runner run test_simple_hook \
    --file /tmp/test_hook.sloth \
    --yes

echo "✓ Teste de hook simples concluído"
```

### 3.2 Testar Hook com Event Module

```bash
# Criar hook que emite eventos
cat > /tmp/test_hook_event.sloth << 'SLOTH'
task("test_hook_with_event", {
    on = function()
        local event = require("event")
        
        -- Emitir evento de teste
        event.emit({
            type = "test.hook.execution",
            severity = "info",
            message = "Hook de teste executado com sucesso",
            metadata = {
                timestamp = os.time(),
                test_id = "hook_test_001"
            }
        })
        
        shell("echo 'Hook com evento emitido'")
        return true
    end
})
SLOTH

# Executar
$HOME/.local/bin/sloth-runner run test_hook_with_event \
    --file /tmp/test_hook_event.sloth \
    --yes

echo "✓ Hook com evento executado"
```

### 3.3 Testar Hook com Condicional

```bash
# Criar hook condicional
cat > /tmp/test_hook_conditional.sloth << 'SLOTH'
task("test_conditional_hook", {
    on = function()
        local hour = tonumber(os.date("%H"))
        
        if hour >= 9 and hour <= 18 then
            shell("echo 'Horário comercial - hook executado'")
            return true
        else
            shell("echo 'Fora do horário comercial - hook pulado'")
            return false
        end
    end
})
SLOTH

# Executar
$HOME/.local/bin/sloth-runner run test_conditional_hook \
    --file /tmp/test_hook_conditional.sloth \
    --yes

echo "✓ Hook condicional testado"
```

---

## 📡 Etapa 4: Testes de Eventos

### 4.1 Testar Emissão de Eventos

```bash
# Criar task que emite vários eventos
cat > /tmp/test_events.sloth << 'SLOTH'
task("test_event_emission", {
    on = function()
        local event = require("event")
        
        -- Evento de info
        event.emit({
            type = "test.info",
            severity = "info",
            message = "Evento de informação de teste"
        })
        
        -- Evento de warning
        event.emit({
            type = "test.warning",
            severity = "warning",
            message = "Evento de aviso de teste"
        })
        
        -- Evento de error
        event.emit({
            type = "test.error",
            severity = "error",
            message = "Evento de erro de teste"
        })
        
        -- Evento com metadata complexo
        event.emit({
            type = "test.complex",
            severity = "info",
            message = "Evento complexo de teste",
            metadata = {
                user = "test_user",
                action = "test_action",
                timestamp = os.time(),
                nested = {
                    key1 = "value1",
                    key2 = "value2"
                }
            }
        })
        
        shell("echo 'Eventos emitidos com sucesso'")
        return true
    end
})
SLOTH

# Executar
$HOME/.local/bin/sloth-runner run test_event_emission \
    --file /tmp/test_events.sloth \
    --yes

echo "✓ Eventos emitidos"
```

### 4.2 Verificar Eventos no CLI

```bash
# Listar todos os eventos
$HOME/.local/bin/sloth-runner events list

# Filtrar por tipo
$HOME/.local/bin/sloth-runner events list --type "test.info"

# Filtrar por severidade
$HOME/.local/bin/sloth-runner events list --severity "error"

# Últimos 5 eventos
$HOME/.local/bin/sloth-runner events list --limit 5

echo "✓ Eventos listados via CLI"
```

### 4.3 Verificar Eventos no Web UI

```bash
echo "📊 Verificação no Web UI:"
echo "1. Acesse: http://localhost:8080/events"
echo "2. Verifique se os eventos aparecem na lista"
echo "3. Teste os filtros (tipo, severidade, agent)"
echo "4. Verifique o stream em tempo real"
echo ""
echo "Pressione ENTER após verificar no navegador..."
read
```

---

## 👁️ Etapa 5: Testes de Watchers

### 5.1 Testar File Watcher

```bash
# Criar diretório de teste
mkdir -p /tmp/sloth-watcher-test

# Criar watcher de arquivo
cat > /tmp/test_file_watcher.sloth << 'SLOTH'
task("test_file_watcher", {
    on = function()
        local watcher = require("watcher")
        local event = require("event")
        
        -- Registrar watcher
        watcher.register({
            id = "test_file_watcher_001",
            type = "file",
            file_path = "/tmp/sloth-watcher-test/test.txt",
            when = {"file.created", "file.modified", "file.deleted"},
            interval = "5s",
            action = function(evt)
                event.emit({
                    type = "watcher.file.event",
                    severity = "info",
                    message = "File watcher detectou: " .. evt.type,
                    metadata = {
                        file = evt.file_path,
                        event_type = evt.type
                    }
                })
            end
        })
        
        shell("echo 'File watcher registrado'")
        return true
    end
})
SLOTH

# Executar no agent
$HOME/.local/bin/sloth-runner run test_file_watcher \
    --file /tmp/test_file_watcher.sloth \
    --delegate-to lady-guica \
    --yes

echo "✓ File watcher registrado"

# Aguardar um pouco
sleep 2

# Criar arquivo para triggerar watcher
ssh lady-guica "echo 'teste' > /tmp/sloth-watcher-test/test.txt"
echo "✓ Arquivo criado"

# Aguardar watcher detectar
sleep 6

# Modificar arquivo
ssh lady-guica "echo 'modificado' >> /tmp/sloth-watcher-test/test.txt"
echo "✓ Arquivo modificado"

# Aguardar
sleep 6

# Verificar eventos
$HOME/.local/bin/sloth-runner events list --type "watcher.file.event"
```

### 5.2 Testar Process Watcher

```bash
# Criar watcher de processo
cat > /tmp/test_process_watcher.sloth << 'SLOTH'
task("test_process_watcher", {
    on = function()
        local watcher = require("watcher")
        local event = require("event")
        
        watcher.register({
            id = "test_process_watcher_001",
            type = "process",
            process_name = "sloth-runner",
            when = {"process.started", "process.stopped"},
            interval = "10s",
            action = function(evt)
                event.emit({
                    type = "watcher.process.event",
                    severity = "info",
                    message = "Process watcher detectou: " .. evt.type,
                    metadata = {
                        process = evt.process_name,
                        event_type = evt.type
                    }
                })
            end
        })
        
        shell("echo 'Process watcher registrado'")
        return true
    end
})
SLOTH

# Executar no agent
$HOME/.local/bin/sloth-runner run test_process_watcher \
    --file /tmp/test_process_watcher.sloth \
    --delegate-to lady-guica \
    --yes

echo "✓ Process watcher registrado"
```

### 5.3 Testar Resource Watcher

```bash
# Criar watcher de recursos
cat > /tmp/test_resource_watcher.sloth << 'SLOTH'
task("test_resource_watcher", {
    on = function()
        local watcher = require("watcher")
        local event = require("event")
        
        watcher.register({
            id = "test_resource_watcher_001",
            type = "resource",
            cpu_threshold = 50.0,
            memory_threshold = 70.0,
            when = {"resource.cpu_high", "resource.memory_high"},
            interval = "10s",
            action = function(evt)
                event.emit({
                    type = "watcher.resource.event",
                    severity = "warning",
                    message = "Resource watcher alertou: " .. evt.type,
                    metadata = {
                        cpu = evt.cpu_percent,
                        memory = evt.memory_percent,
                        event_type = evt.type
                    }
                })
            end
        })
        
        shell("echo 'Resource watcher registrado'")
        return true
    end
})
SLOTH

# Executar no agent
$HOME/.local/bin/sloth-runner run test_resource_watcher \
    --file /tmp/test_resource_watcher.sloth \
    --delegate-to lady-guica \
    --yes

echo "✓ Resource watcher registrado"
```

### 5.4 Listar Watchers Ativos

```bash
# Listar watchers no CLI
$HOME/.local/bin/sloth-runner agent watchers lady-guica

echo "✓ Watchers listados"
```

---

## 🔗 Etapa 6: Teste de Integração Completa

### 6.1 Cenário Completo: Deploy Automático

```bash
# Criar cenário de deploy com hooks e watchers
cat > /tmp/test_full_integration.sloth << 'SLOTH'
task("deploy_with_monitoring", {
    on = function()
        local event = require("event")
        local watcher = require("watcher")
        
        -- Emitir evento de início
        event.emit({
            type = "deploy.started",
            severity = "info",
            message = "Deploy iniciado",
            metadata = { version = "1.0.0" }
        })
        
        -- Registrar watcher para monitorar arquivo de deploy
        watcher.register({
            id = "deploy_monitor",
            type = "file",
            file_path = "/tmp/sloth-deploy/status.txt",
            when = {"file.modified"},
            interval = "5s",
            action = function(evt)
                event.emit({
                    type = "deploy.status_changed",
                    severity = "info",
                    message = "Status do deploy alterado"
                })
            end
        })
        
        -- Simular deploy
        shell("mkdir -p /tmp/sloth-deploy")
        shell("echo 'iniciando' > /tmp/sloth-deploy/status.txt")
        shell("sleep 2")
        shell("echo 'processando' > /tmp/sloth-deploy/status.txt")
        shell("sleep 2")
        shell("echo 'concluído' > /tmp/sloth-deploy/status.txt")
        
        -- Emitir evento de conclusão
        event.emit({
            type = "deploy.completed",
            severity = "info",
            message = "Deploy concluído com sucesso"
        })
        
        return true
    end
})
SLOTH

# Executar
$HOME/.local/bin/sloth-runner run deploy_with_monitoring \
    --file /tmp/test_full_integration.sloth \
    --delegate-to lady-guica \
    --yes

echo "✓ Integração completa testada"

# Aguardar eventos serem processados
sleep 5

# Verificar eventos gerados
$HOME/.local/bin/sloth-runner events list --type "deploy."
```

---

## 📊 Etapa 7: Verificação no Web UI

### 7.1 Dashboard Principal

```bash
echo "📊 Verificação no Dashboard:"
echo "1. Acesse: http://localhost:8080"
echo "2. Verifique cards de resumo (agents, hooks, eventos)"
echo "3. Verifique gráficos de métricas em tempo real"
echo ""
```

### 7.2 Página de Eventos

```bash
echo "📊 Verificação de Eventos:"
echo "1. Acesse: http://localhost:8080/events"
echo "2. Verifique lista de eventos"
echo "3. Teste filtros (tipo, severidade, agent)"
echo "4. Verifique timeline de eventos"
echo "5. Teste busca por texto"
echo ""
```

### 7.3 Página de Hooks

```bash
echo "📊 Verificação de Hooks:"
echo "1. Acesse: http://localhost:8080/hooks"
echo "2. Verifique lista de hooks executados"
echo "3. Verifique estatísticas (sucessos, falhas, tempo médio)"
echo "4. Clique em detalhes de um hook"
echo "5. Verifique logs de execução"
echo ""
```

### 7.4 Página de Watchers

```bash
echo "📊 Verificação de Watchers:"
echo "1. Acesse: http://localhost:8080/watchers"
echo "2. Verifique lista de watchers ativos"
echo "3. Verifique status de cada watcher"
echo "4. Clique em detalhes de um watcher"
echo ""
```

### 7.5 Dashboard de Agent

```bash
echo "📊 Verificação de Agent Dashboard:"
echo "1. Acesse: http://localhost:8080/agents"
echo "2. Clique em 'Dashboard' de lady-guica"
echo "3. Verifique métricas em tempo real (CPU, Memory, Disk)"
echo "4. Verifique aba de Network"
echo "5. Verifique gráficos históricos"
echo ""
```

### 7.6 Network Dashboard

```bash
echo "📊 Verificação de Network Dashboard:"
echo "1. Acesse: http://localhost:8080/network"
echo "2. Verifique cards de resumo (agents, interfaces, download, upload)"
echo "3. Verifique gráfico de tráfego por agent"
echo "4. Acesse: http://localhost:8080/network/topology"
echo "5. Verifique visualização de topologia D3.js"
echo "6. Teste zoom e drag dos nodes"
echo ""
```

---

## 📝 Etapa 8: Verificação de Logs

### 8.1 Logs do Master

```bash
# Ver logs do master
tail -50 /tmp/sloth-master.log

# Procurar por erros
grep -i error /tmp/sloth-master.log

# Procurar por eventos
grep -i event /tmp/sloth-master.log
```

### 8.2 Logs de Hooks

```bash
# Ver logs de execução de hooks
cat ~/.sloth-runner/hook-executions.log | tail -20

# Ver últimos hooks
$HOME/.local/bin/sloth-runner hook logs --limit 10
```

### 8.3 Logs do Agent

```bash
# Ver logs do agent (ajustar comando conforme seu setup)
ssh lady-guica "tail -50 /tmp/sloth-agent.log"

# Ou via journalctl se estiver como serviço
ssh lady-guica "journalctl -u sloth-runner-agent -n 50"
```

---

## 🧹 Etapa 9: Limpeza e Verificação Final

### 9.1 Remover Watchers de Teste

```bash
# Remover watchers via API (se implementado)
# Ou aguardar restart do agent

echo "Watchers ativos serão mantidos. Para remover:"
echo "1. Reinicie o agent"
echo "2. Ou aguarde timeout dos watchers"
```

### 9.2 Limpar Arquivos de Teste

```bash
# Remover arquivos de teste
rm -f /tmp/test_*.sloth
rm -rf /tmp/sloth-watcher-test
rm -rf /tmp/sloth-deploy

# No agent
ssh lady-guica "rm -rf /tmp/sloth-watcher-test /tmp/sloth-deploy"

echo "✓ Arquivos de teste limpos"
```

### 9.3 Estatísticas Finais

```bash
echo ""
echo "==================================="
echo "📊 ESTATÍSTICAS FINAIS"
echo "==================================="

# Contar eventos
EVENT_COUNT=$($HOME/.local/bin/sloth-runner events list | wc -l)
echo "Total de eventos: $EVENT_COUNT"

# Contar hooks
HOOK_COUNT=$(cat ~/.sloth-runner/hook-executions.log 2>/dev/null | wc -l)
echo "Total de hooks executados: $HOOK_COUNT"

# Listar watchers
echo ""
echo "Watchers ativos:"
$HOME/.local/bin/sloth-runner agent watchers lady-guica

echo ""
echo "==================================="
```

---

## ✅ Checklist de Validação

Marque cada item conforme completa:

### Sistema Básico
- [ ] Master server iniciou sem erros
- [ ] Web UI acessível em http://localhost:8080
- [ ] Agent conectado e listado
- [ ] Bancos de dados criados automaticamente

### Hooks
- [ ] Hook simples executou com sucesso
- [ ] Hook com evento emitiu evento corretamente
- [ ] Hook condicional funcionou
- [ ] Hooks aparecem na página `/hooks` do Web UI
- [ ] Estatísticas de hooks aparecem corretamente

### Eventos
- [ ] Eventos foram emitidos com sucesso
- [ ] Eventos aparecem no CLI com `events list`
- [ ] Eventos aparecem na página `/events` do Web UI
- [ ] Filtros de eventos funcionam (tipo, severidade, agent)
- [ ] Timeline de eventos aparece corretamente
- [ ] Stream em tempo real funciona

### Watchers
- [ ] File watcher registrou e detectou mudanças
- [ ] Process watcher registrou corretamente
- [ ] Resource watcher registrou corretamente
- [ ] Watchers aparecem listados no CLI
- [ ] Eventos de watchers foram emitidos
- [ ] Watchers aparecem no Web UI

### Web UI
- [ ] Dashboard principal mostra resumo
- [ ] Gráficos de métricas funcionam
- [ ] Página de agents lista todos os agents
- [ ] Botões Dashboard/Details/Logs funcionam
- [ ] Agent dashboard mostra métricas em tempo real
- [ ] Network dashboard mostra estatísticas
- [ ] Network topology visualização funciona
- [ ] Todas as páginas carregam sem erros 404

### Integração
- [ ] Cenário completo (deploy) executou com sucesso
- [ ] Todos os eventos do cenário foram registrados
- [ ] Watchers do cenário detectaram mudanças
- [ ] Não há erros nos logs

---

## 🐛 Troubleshooting

### Problema: Web UI não carrega

```bash
# Verificar se master está rodando
ps aux | grep "sloth-runner.*master"

# Verificar logs
tail -100 /tmp/sloth-master.log | grep -i error

# Testar porta
lsof -i :8080

# Reiniciar master
pkill -f "sloth-runner.*master"
$HOME/.local/bin/sloth-runner master --bind-address 0.0.0.0 --port 50051 --web-ui --web-port 8080 &
```

### Problema: Eventos não aparecem

```bash
# Verificar banco de eventos
sqlite3 ~/.sloth-runner/events.db "SELECT COUNT(*) FROM events;"

# Verificar permissões
ls -lh ~/.sloth-runner/events.db

# Emitir evento de teste
cat > /tmp/test_event.sloth << 'SLOTH'
task("emit_test", {
    on = function()
        local event = require("event")
        event.emit({type = "test", severity = "info", message = "teste"})
        return true
    end
})
SLOTH

$HOME/.local/bin/sloth-runner run emit_test --file /tmp/test_event.sloth --yes
```

### Problema: Watchers não detectam mudanças

```bash
# Verificar se watcher está registrado
$HOME/.local/bin/sloth-runner agent watchers lady-guica

# Verificar logs do agent
ssh lady-guica "tail -100 /var/log/sloth-runner-agent.log | grep -i watcher"

# Verificar intervalo do watcher (pode estar muito longo)
```

### Problema: Agent não conecta

```bash
# Testar conectividade
nc -zv lady-guica 50060

# Verificar agent está rodando
ssh lady-guica "ps aux | grep sloth-runner"

# Verificar logs do agent
ssh lady-guica "journalctl -u sloth-runner-agent -n 50"

# Reiniciar agent
ssh lady-guica "systemctl restart sloth-runner-agent"
```

---

## 🎯 Resultado Esperado

Ao final deste procedimento, você deve ter:

✅ **Sistema completo funcionando**
- Master server rodando
- Web UI acessível
- Agent conectado

✅ **Hooks testados**
- Hooks simples, com eventos, condicionais
- Logs de hooks visíveis

✅ **Eventos testados**
- Eventos emitidos e armazenados
- Visíveis no CLI e Web UI
- Filtros funcionando

✅ **Watchers testados**
- File, process, resource watchers
- Detectando mudanças corretamente
- Emitindo eventos

✅ **Web UI completo**
- Todas as páginas funcionando
- Gráficos e visualizações renderizando
- Dados em tempo real atualizando

✅ **Zero erros**
- Logs limpos de erros críticos
- Todos os comandos executando com sucesso

---

## 📞 Suporte

Se encontrar problemas:

1. Verifique logs: `/tmp/sloth-master.log`
2. Verifique bancos: `~/.sloth-runner/*.db`
3. Execute troubleshooting acima
4. Verifique issues no GitHub

---

**Última atualização**: 2025-10-08
**Versão**: 1.0.0
