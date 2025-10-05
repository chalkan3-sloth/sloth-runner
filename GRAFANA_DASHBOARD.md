# 📊 Grafana Dashboard - Terminal-Based Metrics Visualization

## Visão Geral

O comando `sloth-runner agent metrics grafana` fornece um dashboard completo e detalhado no terminal, com visualizações ricas de todas as métricas do agente.

## Funcionalidades Implementadas

### ✅ Componentes Criados

1. **internal/telemetry/visualizer.go** (442 linhas)
   - `FetchMetrics()`: Busca e parseia métricas do endpoint Prometheus
   - `parseMetrics()`: Parseia formato de texto Prometheus
   - `DisplayDashboard()`: Exibe dashboard completo com pterm
   - `DisplayHistoricalTrends()`: Função preparada para tendências futuras

2. **Comando CLI: `agent metrics grafana`**
   - Localização: `cmd/sloth-runner/main.go` (linhas 2060-2151)
   - Flags:
     - `--watch`: Atualização contínua do dashboard
     - `--interval`: Intervalo de atualização em segundos (padrão: 5s)

## 📊 Seções do Dashboard

### 1. 🔧 Agent Information
Tabela com informações do agente:
- Version (versão do build)
- OS (sistema operacional)
- Architecture (arquitetura: arm64, amd64, etc.)
- Uptime (tempo de execução formatado: Xd Yh Zm)
- Last Updated (timestamp da última atualização)

### 2. 💻 System Resources
Barras de progresso visuais com cores dinâmicas:
- **Goroutines**: Contador de goroutines
  - Verde: < 60% do máximo (1000)
  - Amarelo: 60-80%
  - Vermelho: > 80%
- **Memory (MB)**: Memória alocada em megabytes
  - Verde: < 60% do máximo (512MB)
  - Amarelo: 60-80%
  - Vermelho: > 80%

### 3. 📋 Task Metrics
Tabela de resumo de tarefas por status:
- ✓ Success (verde)
- ✗ Failed (vermelho)
- ⊘ Skipped (amarelo)
- Contador de tarefas em execução com barra de progresso

### 4. ⏱️ Task Performance
Tabela de performance com latências:
- Nome da tarefa
- P50 latency (50º percentil em ms)
- P99 latency (99º percentil em ms)
- Status de performance:
  - 🟢 Fast: P99 < 1000ms
  - 🟡 Normal: P99 < 5000ms
  - 🔴 Slow: P99 >= 5000ms

### 5. 🌐 gRPC Metrics
Tabela de métricas gRPC:
- Method (nome do método)
- Requests (total de requisições)
- Avg Latency (latência média P50 em ms)

### 6. ⚠️ Errors
Tabela de erros por tipo (quando houver):
- Error Type (tipo do erro)
- Count (contador em vermelho)

### 7. 📈 Summary
Box de resumo com totais:
- Total Tasks (total de tarefas executadas)
- Running (tarefas em execução)
- Memory (memória atual em MB)
- Goroutines (goroutines ativas)

## 🚀 Uso

### Comando Básico
```bash
./sloth-runner-telemetry agent metrics grafana lady-arch
```

Exibe o dashboard uma vez com snapshot das métricas atuais.

### Watch Mode
```bash
# Atualiza a cada 5 segundos (padrão)
./sloth-runner-telemetry agent metrics grafana lady-arch --watch

# Atualiza a cada 10 segundos
./sloth-runner-telemetry agent metrics grafana lady-arch --watch --interval 10

# Atualiza a cada 1 segundo (monitoramento em tempo real)
./sloth-runner-telemetry agent metrics grafana lady-arch --watch --interval 1
```

**Recursos do Watch Mode**:
- Limpa a tela entre atualizações
- Atualiza métricas automaticamente
- Pressione Ctrl+C para parar
- Ideal para monitorar execução de tarefas em tempo real

### Help
```bash
./sloth-runner-telemetry agent metrics grafana --help
```

## 📋 Exemplo de Output

```
════════════════════════════════════════════════════════════════════════
             📊 Sloth Runner Metrics Dashboard - Agent: lady-arch
════════════════════════════════════════════════════════════════════════

🔧 Agent Information
┌────────────────┬─────────────────────────────────┐
│ Version        │ dev                             │
│ OS             │ linux                           │
│ Architecture   │ arm64                           │
│ Uptime         │ 2h 34m                          │
│ Last Updated   │ 2025-10-05 15:42:30             │
└────────────────┴─────────────────────────────────┘

💻 System Resources

Goroutines: [████████████░░░░░░░░░░░░░░░░░░░░░░░░░░] 342/1000 (34.2%)
Memory (MB): [██████░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 78/512 (15.2%)

📋 Task Metrics
┌───────────┬──────────┐
│  Status   │  Count   │
├───────────┼──────────┤
│ ✓ Success │ 145      │
│ ✗ Failed  │ 3        │
│ ⊘ Skipped │ 12       │
└───────────┴──────────┘

Running Tasks: [██░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░] 2/10 (20.0%)

⏱️  Task Performance
┌──────────────────┬───────────┬───────────┬──────────┐
│      Task        │ P50 (ms)  │ P99 (ms)  │  Status  │
├──────────────────┼───────────┼───────────┼──────────┤
│ install_packages │ 234.56    │ 567.89    │ 🟡 Normal│
│ check_service    │ 12.34     │ 45.67     │ 🟢 Fast  │
│ deploy_app       │ 1234.56   │ 5678.90   │ 🔴 Slow  │
└──────────────────┴───────────┴───────────┴──────────┘

🌐 gRPC Metrics
┌─────────────────┬───────────┬──────────────────┐
│     Method      │ Requests  │ Avg Latency (ms) │
├─────────────────┼───────────┼──────────────────┤
│ ExecuteTask     │ 156       │ 234.56           │
│ ExecuteCommand  │ 45        │ 12.34            │
└─────────────────┴───────────┴──────────────────┘

╔═══════════════════════════════════════════════════════╗
║                      📈 Summary                       ║
╠═══════════════════════════════════════════════════════╣
║ Total Tasks: 160 | Running: 2 | Memory: 78 MB |       ║
║ Goroutines: 342                                       ║
╚═══════════════════════════════════════════════════════╝
```

## 🔍 Detalhes Técnicos

### Parser de Métricas Prometheus
O visualizer parseia o formato de texto Prometheus:
- Suporta métricas com e sem labels
- Extrai valores de quantis (0.5, 0.99) para histogramas
- Agrupa métricas por tipo e label
- Calcula totais e agregações

### Bibliotecas Utilizadas
- **pterm**: Visualização rica no terminal
  - Tables com bordas e headers
  - Progress bars com cores
  - Boxes e sections
  - Styled text (cores, bold, etc.)

### Formato de Cores
- **Verde**: Sucesso, valores normais
- **Amarelo**: Avisos, valores médios
- **Vermelho**: Erros, valores altos
- **Ciano**: Headers e labels
- **Magenta**: Valores secundários

## 🧪 Testing

### Pré-requisitos
1. Agente rodando com telemetria:
   ```bash
   ./sloth-runner-new agent start --name lady-arch --telemetry --metrics-port 9090
   ```

2. Master rodando:
   ```bash
   ./sloth-runner master start
   ```

### Teste Básico
```bash
# Verificar se o agente está registrado
./sloth-runner-telemetry agent list

# Testar dashboard
./sloth-runner-telemetry agent metrics grafana lady-arch
```

### Teste com Carga
```bash
# Em um terminal: watch mode
./sloth-runner-telemetry agent metrics grafana lady-arch --watch --interval 2

# Em outro terminal: executar tarefas
./sloth-runner-telemetry agent run lady-arch "sleep 5"
./sloth-runner-telemetry agent run lady-arch "echo test"
```

Você verá as métricas sendo atualizadas em tempo real:
- Running tasks aumenta durante execução
- Tasks total incrementa ao completar
- Task duration é registrado
- Memory e goroutines variam

## 🚢 Deployment para lady-arch

### 1. Transfer Binary
```bash
# Método manual (devido a issues SSH)
# No macOS, copiar o binary:
ls -lh sloth-runner-linux-arm64-telemetry

# Usar incus file push (a partir do host igor@192.168.1.16):
ssh igor@192.168.1.16
incus file push ~/sloth-runner-linux-arm64-telemetry lady-arch/home/igor/sloth-runner-grafana
```

### 2. Deploy no Agent
```bash
# Entrar no container
incus exec lady-arch -- bash

# Parar agente antigo
pkill sloth-runner

# Dar permissão de execução
chmod +x sloth-runner-grafana

# Iniciar com telemetria
./sloth-runner-grafana agent start --name lady-arch --master 192.168.1.2:50053 --telemetry --metrics-port 9090 &

# Verificar se telemetria está rodando
curl http://localhost:9090/metrics
```

### 3. Testar do macOS
```bash
# Verificar endpoint
./sloth-runner-telemetry agent metrics prom lady-arch

# Visualizar dashboard
./sloth-runner-telemetry agent metrics grafana lady-arch

# Watch mode
./sloth-runner-telemetry agent metrics grafana lady-arch --watch
```

## 📚 Referências

- Formato de métricas: [Prometheus Text Format](https://prometheus.io/docs/instrumenting/exposition_formats/)
- pterm library: https://github.com/pterm/pterm
- Documentação completa: `TELEMETRY_TESTING.md`

## 🔮 Próximos Passos

Possíveis melhorias futuras:
1. **Historical Trends**: Gráficos de tendência ao longo do tempo
2. **Sparklines**: Mini-gráficos inline para métricas
3. **Alertas**: Destacar métricas que excedem thresholds
4. **Export**: Salvar snapshot do dashboard em arquivo
5. **Comparação**: Comparar métricas entre múltiplos agentes
6. **Filtros**: Filtrar por grupo de tarefas ou período
