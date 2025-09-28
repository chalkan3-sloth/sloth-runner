# 📊 Módulo de Métricas e Monitoramento

O módulo **Métricas e Monitoramento** fornece capacidades abrangentes de monitoramento do sistema, coleta de métricas customizadas e verificação de saúde. Ele habilita observabilidade em tempo real tanto dos recursos do sistema quanto da performance da aplicação.

## 🚀 Recursos Principais

- **Métricas do Sistema**: Coleta automática de métricas de CPU, memória, disco e rede
- **Métricas de Runtime**: Informações do runtime Go (goroutines, heap, GC)
- **Métricas Customizadas**: Gauges, contadores, histogramas e timers
- **Verificações de Saúde**: Monitoramento automático da saúde do sistema
- **Endpoints HTTP**: Export de métricas compatível com Prometheus
- **Sistema de Alertas**: Alertas baseados em thresholds
- **API JSON**: Dados completos de métricas para integrações

## 📊 Métricas do Sistema

### Monitoramento de CPU, Memória e Disco

```lua
-- Obter uso atual de CPU
local uso_cpu = metrics.system_cpu()
log.info("Uso de CPU: " .. string.format("%.1f%%", uso_cpu))

-- Obter informações de memória
local info_memoria = metrics.system_memory()
log.info("Memória: " .. string.format("%.1f%% (%.0f/%.0f MB)", 
    info_memoria.percent, info_memoria.used_mb, info_memoria.total_mb))

-- Obter uso de disco
local info_disco = metrics.system_disk("/")
log.info("Disco: " .. string.format("%.1f%% (%.1f/%.1f GB)", 
    info_disco.percent, info_disco.used_gb, info_disco.total_gb))

-- Verificar caminho específico do disco
local disco_var = metrics.system_disk("/var")
log.info("Uso do disco /var: " .. string.format("%.1f%%", disco_var.percent))
```

### Informações de Runtime

```lua
-- Obter métricas do runtime Go
local runtime = metrics.runtime_info()
log.info("Informações de Runtime:")
log.info("  Goroutines: " .. runtime.goroutines)
log.info("  Núcleos de CPU: " .. runtime.num_cpu)
log.info("  Heap alocado: " .. string.format("%.1f MB", runtime.heap_alloc_mb))
log.info("  Heap do sistema: " .. string.format("%.1f MB", runtime.heap_sys_mb))
log.info("  Ciclos de GC: " .. runtime.num_gc)
log.info("  Versão do Go: " .. runtime.go_version)
```

## 📈 Métricas Customizadas

### Métricas Gauge (Valores Atuais)

```lua
-- Definir valores simples de gauge
metrics.gauge("temperatura_cpu", 65.4)
metrics.gauge("conexoes_ativas", 142)
metrics.gauge("tamanho_fila", 23)

-- Definir gauge com tags
metrics.gauge("uso_memoria", percentual_memoria, {
    servidor = "web-01",
    ambiente = "producao",
    regiao = "us-east-1"
})

-- Atualizar status de deployment
metrics.gauge("progresso_deployment", 75.5, {
    app = "frontend",
    versao = "v2.1.0"
})
```

### Métricas Counter (Valores Incrementais)

```lua
-- Incrementar contadores
local total_requisicoes = metrics.counter("requisicoes_http_total", 1)
local contador_erros = metrics.counter("erros_http_total", 1, {
    codigo_status = "500",
    endpoint = "/api/usuarios"
})

-- Incremento em lote
local processados = metrics.counter("mensagens_processadas", 50, {
    fila = "notificacoes_usuario",
    prioridade = "alta"
})

log.info("Total de requisições processadas: " .. total_requisicoes)
```

### Métricas Histogram (Distribuição de Valores)

```lua
-- Registrar tempos de resposta
metrics.histogram("tempo_resposta_ms", 245.6, {
    endpoint = "/api/usuarios",
    metodo = "GET"
})

-- Registrar tamanhos de payload
metrics.histogram("tamanho_payload_bytes", 1024, {
    tipo_conteudo = "application/json"
})

-- Registrar tamanhos de lote
metrics.histogram("tamanho_lote", 150, {
    operacao = "insercao_lote",
    tabela = "eventos_usuario"
})
```

### Métricas Timer (Tempo de Execução de Funções)

```lua
-- Cronometrar execução de função automaticamente
local duracao = metrics.timer("consulta_banco", function()
    -- Simular consulta ao banco
    local resultado = exec.run("sleep 0.1")
    return resultado
end, {
    tipo_consulta = "select",
    tabela = "usuarios"
})

log.info("Consulta ao banco levou: " .. string.format("%.2f ms", duracao))

-- Cronometrar operações complexas
local tempo_processamento = metrics.timer("processamento_dados", function()
    -- Processar dataset grande
    local dados = {}
    for i = 1, 100000 do
        dados[i] = math.sqrt(i) * 2.5
    end
    return #dados
end, {
    operacao = "computacao_matematica",
    tamanho = "grande"
})

log.info("Processamento de dados concluído em: " .. string.format("%.2f ms", tempo_processamento))
```

## 🏥 Monitoramento de Saúde

### Status de Saúde Automático

```lua
-- Obter status abrangente de saúde
local saude = metrics.health_status()
log.info("Status Geral de Saúde: " .. saude.overall)

-- Verificar componentes individuais
local componentes = {"cpu", "memory", "disk"}
for _, componente in ipairs(componentes) do
    local info_comp = saude[componente]
    if info_comp then
        local icone_status = "✅"
        if info_comp.status == "warning" then
            icone_status = "⚠️"
        elseif info_comp.status == "critical" then
            icone_status = "❌"
        end
        
        log.info(string.format("  %s %s: %.1f%% (%s)", 
            icone_status, componente:upper(), info_comp.usage, info_comp.status))
    end
end
```

### Verificações de Saúde Customizadas

```lua
-- Criar função de verificação de saúde
function verificar_saude_aplicacao()
    local pontuacao_saude = 100
    local problemas = {}
    
    -- Verificar conectividade do banco
    local resultado_bd = exec.run("pg_isready -h localhost -p 5432")
    if resultado_bd ~= "" then
        pontuacao_saude = pontuacao_saude - 20
        table.insert(problemas, "Falha na conexão com o banco de dados")
    end
    
    -- Verificar espaço em disco
    local disco = metrics.system_disk("/")
    if disco.percent > 90 then
        pontuacao_saude = pontuacao_saude - 30
        table.insert(problemas, "Espaço em disco crítico: " .. string.format("%.1f%%", disco.percent))
    end
    
    -- Verificar uso de memória
    local memoria = metrics.system_memory()
    if memoria.percent > 85 then
        pontuacao_saude = pontuacao_saude - 25
        table.insert(problemas, "Uso de memória alto: " .. string.format("%.1f%%", memoria.percent))
    end
    
    -- Registrar pontuação de saúde
    metrics.gauge("pontuacao_saude_aplicacao", pontuacao_saude)
    
    if pontuacao_saude < 70 then
        metrics.alert("saude_aplicacao", {
            level = "warning",
            message = "Saúde da aplicação degradada: " .. table.concat(problemas, ", "),
            pontuacao = pontuacao_saude
        })
    end
    
    return pontuacao_saude >= 70
end

-- Usar em tasks
TaskDefinitions = {
    monitoramento_saude = {
        tasks = {
            verificacao_saude = {
                command = function()
                    local saudavel = verificar_saude_aplicacao()
                    return saudavel, saudavel and "Sistema saudável" or "Problemas de saúde detectados"
                end
            }
        }
    }
}
```

## 🚨 Sistema de Alertas

### Criando Alertas

```lua
-- Alerta simples por threshold
local cpu = metrics.system_cpu()
if cpu > 80 then
    metrics.alert("uso_alto_cpu", {
        level = "warning",
        message = "Uso de CPU está alto: " .. string.format("%.1f%%", cpu),
        threshold = 80,
        value = cpu,
        severidade = "media"
    })
end

-- Alerta complexo com múltiplas condições
local memoria = metrics.system_memory()
local disco = metrics.system_disk()

if memoria.percent > 90 and disco.percent > 85 then
    metrics.alert("esgotamento_recursos", {
        level = "critical",
        message = string.format("Uso crítico de recursos - Memória: %.1f%%, Disco: %.1f%%", 
            memoria.percent, disco.percent),
        uso_memoria = memoria.percent,
        uso_disco = disco.percent,
        acao_recomendada = "Escalar recursos imediatamente"
    })
end

-- Alertas específicos da aplicação
local tamanho_fila = state.get("tamanho_fila_tarefas", 0)
if tamanho_fila > 1000 then
    metrics.alert("acumulo_fila", {
        level = "warning", 
        message = "Acúmulo detectado na fila de tarefas: " .. tamanho_fila .. " itens",
        tamanho_fila = tamanho_fila,
        tempo_processamento_estimado = tamanho_fila * 2 .. " segundos"
    })
end
```

## 🔍 Gerenciamento de Métricas

### Recuperando Métricas Customizadas

```lua
-- Obter métrica customizada específica
local metrica_cpu = metrics.get_custom("temperatura_cpu")
if metrica_cpu then
    log.info("Métrica de temperatura da CPU: " .. data.to_json(metrica_cpu))
end

-- Listar todas as métricas customizadas
local todas_metricas = metrics.list_custom()
log.info("Total de métricas customizadas: " .. #todas_metricas)
for i, nome_metrica in ipairs(todas_metricas) do
    log.info("  " .. i .. ". " .. nome_metrica)
end
```

### Exemplo de Monitoramento de Performance

```lua
TaskDefinitions = {
    monitoramento_performance = {
        tasks = {
            monitorar_performance_api = {
                command = function()
                    -- Iniciar sessão de monitoramento
                    log.info("Iniciando monitoramento de performance da API...")
                    
                    -- Simular chamadas de API e medir performance
                    for i = 1, 10 do
                        local tempo_api = metrics.timer("chamada_api_" .. i, function()
                            -- Simular chamada de API
                            exec.run("curl -s -o /dev/null -w '%{time_total}' https://api.exemplo.com/health")
                        end, {
                            endpoint = "health",
                            numero_chamada = tostring(i)
                        })
                        
                        -- Registrar tempo de resposta
                        metrics.histogram("tempo_resposta_api", tempo_api, {
                            endpoint = "health"
                        })
                        
                        -- Verificar se o tempo de resposta é aceitável
                        if tempo_api > 1000 then -- 1 segundo
                            metrics.counter("chamadas_api_lentas", 1, {
                                endpoint = "health"
                            })
                            
                            metrics.alert("resposta_api_lenta", {
                                level = "warning",
                                message = string.format("Resposta lenta da API: %.2f ms", tempo_api),
                                tempo_resposta = tempo_api,
                                threshold = 1000
                            })
                        end
                        
                        -- Breve atraso entre chamadas
                        exec.run("sleep 0.1")
                    end
                    
                    -- Obter estatísticas resumidas
                    local saude_sistema = metrics.health_status()
                    log.info("Saúde do sistema após testes da API: " .. saude_sistema.overall)
                    
                    return true, "Monitoramento de performance da API concluído"
                end
            }
        }
    }
}
```

## 🌐 Endpoints HTTP

O módulo de métricas expõe automaticamente endpoints HTTP para sistemas de monitoramento externos:

### Formato Prometheus (`/metrics`)
```bash
# Acessar métricas compatíveis com Prometheus
curl http://agente:8080/metrics

# Exemplo de saída:
# sloth_agent_cpu_usage_percent 15.4
# sloth_agent_memory_usage_mb 2048.5
# sloth_agent_disk_usage_percent 67.2
# sloth_agent_tasks_total 142
```

### Formato JSON (`/metrics/json`)
```bash
# Obter métricas completas em formato JSON
curl http://agente:8080/metrics/json

# Exemplo de resposta:
{
  "agent_name": "meuagente1",
  "timestamp": "2024-01-15T10:30:00Z",
  "system": {
    "cpu_usage_percent": 15.4,
    "memory_usage_mb": 2048.5,
    "disk_usage_percent": 67.2
  },
  "runtime": {
    "num_goroutines": 25,
    "heap_alloc_mb": 45.2
  },
  "custom": {
    "tempo_resposta_api": {...},
    "progresso_deployment": 85.5
  }
}
```

### Verificação de Saúde (`/health`)
```bash
# Verificar status de saúde do agente
curl http://agente:8080/health

# Exemplo de resposta:
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "checks": {
    "cpu": {"usage": 15.4, "status": "healthy"},
    "memory": {"usage": 45.8, "status": "healthy"},
    "disk": {"usage": 67.2, "status": "healthy"}
  }
}
```

## 📋 Referência da API

### Métricas do Sistema
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `metrics.system_cpu()` | - | uso: number | Obter percentual atual de uso de CPU |
| `metrics.system_memory()` | - | info: table | Obter informações de uso de memória |
| `metrics.system_disk(caminho?)` | caminho?: string | info: table | Obter uso de disco para caminho (padrão: "/") |
| `metrics.runtime_info()` | - | info: table | Obter informações do runtime Go |

### Métricas Customizadas
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `metrics.gauge(nome, valor, tags?)` | nome: string, valor: number, tags?: table | sucesso: boolean | Definir métrica gauge |
| `metrics.counter(nome, incremento?, tags?)` | nome: string, incremento?: number, tags?: table | novo_valor: number | Incrementar contador |
| `metrics.histogram(nome, valor, tags?)` | nome: string, valor: number, tags?: table | sucesso: boolean | Registrar valor de histograma |
| `metrics.timer(nome, funcao, tags?)` | nome: string, funcao: function, tags?: table | duracao: number | Cronometrar execução de função |

### Saúde e Monitoramento
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `metrics.health_status()` | - | status: table | Obter status abrangente de saúde |
| `metrics.alert(nome, dados)` | nome: string, dados: table | sucesso: boolean | Criar alerta |

### Utilitários
| Função | Parâmetros | Retorno | Descrição |
|--------|------------|---------|-----------|
| `metrics.get_custom(nome)` | nome: string | metrica: table \| nil | Obter métrica customizada por nome |
| `metrics.list_custom()` | - | nomes: table | Listar todos os nomes de métricas customizadas |

## 🎯 Melhores Práticas

1. **Use tipos apropriados de métricas** - gauges para valores atuais, contadores para totais, histogramas para distribuições
2. **Adicione tags significativas** para categorizar e filtrar métricas
3. **Defina thresholds razoáveis para alertas** para evitar fadiga de alertas
4. **Monitore o impacto na performance** da coleta extensiva de métricas
5. **Use timers para operações críticas** para identificar gargalos
6. **Implemente health checks** para todos os componentes críticos do sistema
7. **Exporte métricas para sistemas externos** como Prometheus para armazenamento de longo prazo

O módulo **Métricas e Monitoramento** fornece observabilidade abrangente para seu ambiente distribuído sloth-runner! 📊🚀