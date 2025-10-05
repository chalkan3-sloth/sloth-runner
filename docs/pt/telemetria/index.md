# 📊 Telemetria e Observabilidade

## Visão Geral

O Sloth Runner fornece recursos abrangentes de telemetria e observabilidade através de **integração nativa com Prometheus** e um **dashboard estilo Grafana diretamente no terminal**. Monitore sua frota de agentes, rastreie métricas de execução de tarefas, analise performance e obtenha insights profundos sobre sua automação de infraestrutura.

!!! success "Observabilidade Enterprise"
    Servidor de métricas Prometheus embutido com auto-descoberta, dashboards em tempo real e configuração zero.

## Recursos Principais

### 🎯 Integração Prometheus

- **Exportador Nativo**: Servidor HTTP embutido expondo métricas compatíveis com Prometheus
- **Auto-Descoberta**: Endpoint de métricas automaticamente configurado no startup do agente
- **Formato Padrão**: Compatível com Prometheus, Grafana e todas as ferramentas de observabilidade
- **Configuração Zero**: Telemetria habilitada por padrão com valores sensatos

### 📊 Dashboard no Terminal

- **Visualização Rica**: Dashboard bonito no terminal com tabelas, gráficos e barras de progresso
- **Atualizações em Tempo Real**: Modo watch com intervalos de refresh configuráveis
- **Métricas Abrangentes**: Recursos do sistema, performance de tarefas, estatísticas gRPC e rastreamento de erros
- **Insights Coloridos**: Indicadores visuais para performance e status de saúde

### 📈 Categorias de Métricas

#### Métricas de Tarefas
- Total de tarefas executadas (por status: success, failed, skipped)
- Tarefas executando atualmente
- Histogramas de duração de tarefas (latências P50, P99)
- Rastreamento de performance por tarefa e por grupo

#### Métricas do Sistema
- Uptime do agente
- Alocação de memória
- Contador de goroutines
- Versão do agente e informações de build

#### Métricas gRPC
- Contagem de requisições por método
- Histogramas de duração de requisições
- Taxas de sucesso/erro

#### Rastreamento de Erros
- Contagem de erros por tipo
- Rastreamento de tarefas falhadas
- Monitoramento de erros do sistema

## Início Rápido

### Habilitar Telemetria no Agente

Telemetria está **habilitada por padrão**. Simplesmente inicie seu agente:

```bash
./sloth-runner agent start --name meu-agente --master localhost:50053
```

Para configurar explicitamente a telemetria:

```bash
# Habilitar telemetria com porta customizada
./sloth-runner agent start \
  --name meu-agente \
  --master localhost:50053 \
  --telemetry \
  --metrics-port 9090
```

Para desabilitar telemetria:

```bash
./sloth-runner agent start \
  --name meu-agente \
  --master localhost:50053 \
  --telemetry=false
```

### Acessar Métricas

#### Obter Endpoint Prometheus

```bash
./sloth-runner agent metrics prom meu-agente
```

Saída:
```
✅ Metrics Endpoint:
  URL: http://192.168.1.100:9090/metrics

📝 Usage:
  # Visualizar métricas no navegador:
  open http://192.168.1.100:9090/metrics

  # Buscar métricas via curl:
  curl http://192.168.1.100:9090/metrics
```

#### Visualizar Snapshot

```bash
./sloth-runner agent metrics prom meu-agente --snapshot
```

### Visualizar Dashboard

#### Visualização Única

```bash
./sloth-runner agent metrics grafana meu-agente
```

#### Modo Watch (Auto-Refresh)

```bash
# Refresh a cada 5 segundos (padrão)
./sloth-runner agent metrics grafana meu-agente --watch

# Intervalo de refresh customizado (2 segundos)
./sloth-runner agent metrics grafana meu-agente --watch --interval 2
```

## Casos de Uso

### Desenvolvimento

Monitore suas tarefas durante o desenvolvimento:

```bash
# Terminal 1: Watch dashboard
./sloth-runner agent metrics grafana dev-agent --watch --interval 1

# Terminal 2: Executar tarefas
./sloth-runner run -f deploy.sloth --values dev.yaml
```

### Monitoramento em Produção

Integre com Prometheus e Grafana:

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'sloth-runner-producao'
    static_configs:
      - targets:
          - 'agent1:9090'
          - 'agent2:9090'
          - 'agent3:9090'
        labels:
          environment: production
```

### Análise de Performance

Identifique tarefas lentas e gargalos:

```bash
# Visualizar métricas detalhadas de performance
./sloth-runner agent metrics grafana prod-agent

# Verificar latências P99 na seção Task Performance
# Tarefas com indicador 🔴 Slow precisam de otimização
```

## Próximos Passos

- [Referência de Métricas Prometheus](prometheus-metrics.md) - Documentação completa de métricas
- [Guia do Dashboard Grafana](grafana-dashboard.md) - Funcionalidades e uso do dashboard
- [Guia de Deployment](deployment.md) - Deploy em produção e integração

## Plataformas Suportadas

- ✅ Linux (amd64, arm64)
- ✅ macOS (Intel, Apple Silicon)
- ✅ Windows (via WSL2)
- ✅ Containers (Docker, Incus/LXC)
- ✅ Kubernetes (via DaemonSet)

## Impacto na Performance

A telemetria tem **overhead mínimo de performance**:

- Memória: ~10-20MB adicional
- CPU: <1% sob carga normal
- Rede: Métricas servidas apenas sob demanda (modelo pull)
- Armazenamento: Métricas armazenadas em memória, sem persistência

## Considerações de Segurança

!!! warning "Exposição de Rede"
    O endpoint de métricas é exposto em todas as interfaces de rede por padrão. Em produção:

    - Use regras de firewall para restringir acesso
    - Considere bind apenas em localhost e use reverse proxy
    - Habilite autenticação via reverse proxy (Prometheus não suporta auth nativamente)

## Próximas Leituras

- [Documentação Prometheus](https://prometheus.io/docs/)
- [Documentação Grafana](https://grafana.com/docs/)
- [Biblioteca pterm](https://github.com/pterm/pterm) (usada para visualização no terminal)
