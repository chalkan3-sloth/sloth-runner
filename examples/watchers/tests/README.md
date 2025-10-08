# Watcher System Test Suite

Suite completa de testes para validar o sistema de watchers, eventos e hooks do Sloth Runner.

## 📋 Testes Incluídos

### Test 01: File Watcher
**Arquivo**: `01_file_watcher_test.sloth`
**Objetivo**: Validar detecção de eventos de arquivo (create, change, delete)
**Valida**:
- Registro de file watcher com `delegate_to`
- Detecção de criação de arquivo
- Detecção de modificação de arquivo
- Detecção de deleção de arquivo
- Hash checking funcional

**Uso**:
```bash
sloth-runner run file_watcher_test --file 01_file_watcher_test.sloth --yes
```

---

### Test 02: CPU Watcher
**Arquivo**: `02_cpu_watcher_test.sloth`
**Objetivo**: Validar detecção de threshold de CPU
**Valida**:
- Registro de CPU watcher
- Detecção quando CPU > threshold
- Intervalos de checagem funcionais
- Geração de eventos de CPU

**Uso**:
```bash
sloth-runner run cpu_watcher_test --file 02_cpu_watcher_test.sloth --yes
```

---

### Test 03: Memory Watcher
**Arquivo**: `03_memory_watcher_test.sloth`
**Objetivo**: Validar detecção de threshold de memória
**Valida**:
- Registro de memory watcher
- Detecção quando memory > threshold
- Monitoramento contínuo de memória
- Geração de eventos de memória

**Uso**:
```bash
sloth-runner run memory_watcher_test --file 03_memory_watcher_test.sloth --yes
```

---

### Test 04: Process Watcher
**Arquivo**: `04_process_watcher_test.sloth`
**Objetivo**: Validar detecção de início e fim de processos
**Valida**:
- Registro de process watcher
- Detecção de processo iniciado
- Detecção de processo finalizado
- Filtro por nome de processo

**Uso**:
```bash
sloth-runner run process_watcher_test --file 04_process_watcher_test.sloth --yes
```

---

### Test 05: Port Watcher
**Arquivo**: `05_port_watcher_test.sloth`
**Objetivo**: Validar detecção de portas em listening
**Valida**:
- Registro de port watcher
- Detecção de porta aberta
- Detecção de porta fechada
- Suporte a TCP/UDP

**Uso**:
```bash
sloth-runner run port_watcher_test --file 05_port_watcher_test.sloth --yes
```

---

### Test 06: Watcher + Hook Integration
**Arquivo**: `06_watcher_with_hook_test.sloth`
**Objetivo**: Validar integração completa watcher → event → hook
**Valida**:
- Watchers gerando eventos
- Eventos sendo enfileirados
- Hooks sendo disparados por eventos
- Execução de ações em hooks
- Persistência de execução de hooks

**Uso**:
```bash
sloth-runner run watcher_hook_integration_test --file 06_watcher_with_hook_test.sloth --yes
```

---

### Test 07: Complete End-to-End
**Arquivo**: `07_complete_e2e_test.sloth`
**Objetivo**: Teste completo de todo o sistema
**Valida**:
- Múltiplos watchers simultaneamente
- Múltiplos tipos de eventos
- Processamento de event queue
- Execução de múltiplos hooks
- Integridade do fluxo completo
- Logs e auditoria

**Uso**:
```bash
sloth-runner run complete_e2e_test --file 07_complete_e2e_test.sloth --yes
```

---

## 🚀 Executando os Testes

### Executar Todos os Testes
```bash
cd examples/watchers/tests
./run_all_tests.sh
```

Este script:
- Executa todos os 7 testes sequencialmente
- Salva logs individuais em `/tmp/watcher_tests/`
- Apresenta relatório final com PASSED/FAILED
- Aguarda 10 segundos entre testes para processamento

### Executar Teste Individual
```bash
sloth-runner run <workflow_name> --file <test_file> --yes
```

Exemplo:
```bash
sloth-runner run file_watcher_test --file 01_file_watcher_test.sloth --yes
```

---

## 📊 Verificando Resultados

### Logs do Agent
```bash
ssh chalkan3@192.168.1.16 "cat agent.log | tail -100"
```

### Eventos no Database
```bash
ssh chalkan3@192.168.1.16 "sqlite3 /etc/sloth-runner/events.db 'SELECT * FROM events ORDER BY created_at DESC LIMIT 20'"
```

### Hook Executions
```bash
ssh chalkan3@192.168.1.16 "sqlite3 /etc/sloth-runner/hooks.db 'SELECT * FROM hook_executions ORDER BY executed_at DESC LIMIT 20'"
```

### Logs de Teste
```bash
ls -lh /tmp/watcher_tests/
cat /tmp/watcher_tests/01_file_watcher_test.log
```

---

## 🎯 Critérios de Sucesso

Para cada teste passar, deve:

1. **Watchers**: Registrar com sucesso no agent remoto
2. **Eventos**: Serem gerados quando condições são atendidas
3. **Hooks**: Executarem quando eventos correspondentes ocorrem
4. **Logs**: Conter evidências de todas as etapas
5. **Database**: Ter registros de eventos e hook executions

---

## 🐛 Troubleshooting

### Teste Falhou?

1. **Verificar agent está rodando**:
   ```bash
   ssh chalkan3@192.168.1.16 "ps aux | grep sloth-runner"
   ```

2. **Verificar logs do agent**:
   ```bash
   ssh chalkan3@192.168.1.16 "tail -100 agent.log"
   ```

3. **Verificar conectividade**:
   ```bash
   nc -zv 192.168.1.16 50051
   ```

4. **Verificar watcher manager**:
   ```bash
   ssh chalkan3@192.168.1.16 "cat agent.log | grep 'Watcher registered'"
   ```

5. **Verificar event processor**:
   ```bash
   ssh chalkan3@192.168.1.16 "cat agent.log | grep 'event processor started'"
   ```

### Eventos Não Gerados?

- Verificar intervalo do watcher (pode levar tempo)
- Verificar condições são realmente atendidas
- Verificar logs do watcher manager
- Verificar se watcher foi registrado corretamente

### Hooks Não Executam?

- Verificar hook está em `/etc/sloth-runner/hooks/`
- Verificar sintaxe Lua do hook
- Verificar event_type no hook.on() corresponde
- Verificar dispatcher de hooks está ativo

---

## 📚 Estrutura de Arquivos

```
examples/watchers/tests/
├── README.md                          # Este arquivo
├── run_all_tests.sh                   # Runner para todos os testes
├── 01_file_watcher_test.sloth        # Teste de file watcher
├── 02_cpu_watcher_test.sloth         # Teste de CPU watcher
├── 03_memory_watcher_test.sloth      # Teste de memory watcher
├── 04_process_watcher_test.sloth     # Teste de process watcher
├── 05_port_watcher_test.sloth        # Teste de port watcher
├── 06_watcher_with_hook_test.sloth   # Teste de integração
└── 07_complete_e2e_test.sloth        # Teste end-to-end completo
```

---

## 🔧 Requisitos

- Sloth Runner instalado e configurado
- Agent `lady-guica` rodando em 192.168.1.16:50051
- SSH access ao agent
- Master server rodando (para registry)
- SQLite3 instalado (para queries manuais)

---

## ✅ Exemplo de Saída Esperada

```
╔═══════════════════════════════════════════════════════╗
║     🧪 Watcher System Test Suite                     ║
║     Testing: Watchers, Events, Hooks                 ║
╚═══════════════════════════════════════════════════════╝

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Test 1/7: File Watcher (Create/Change/Delete)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ PASSED

...

╔═══════════════════════════════════════════════════════╗
║              📊 Test Results Summary                  ║
╚═══════════════════════════════════════════════════════╝

Total tests: 7
Passed: 7
Failed: 0

✅ ALL TESTS PASSED!
   Watcher system is working correctly
```

---

## 📞 Suporte

Para problemas ou dúvidas:
1. Verificar logs detalhados
2. Revisar documentação do watcher system
3. Verificar exemplos em `examples/watchers/`
