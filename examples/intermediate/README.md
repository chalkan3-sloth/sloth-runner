# 🟡 Exemplos Intermediários

Esta pasta contém exemplos que combinam múltiplas funcionalidades do Sloth Runner, ideais para quem já domina o básico.

## 📋 Exemplos Disponíveis

### 🚀 [`api-integration.lua`](./api-integration.lua)
**Integração completa com APIs REST**
- 🔧 Configuração avançada de HTTP client
- 🔄 Retry automático com diferentes estratégias
- ✅ Validação de dados com schemas
- 📊 Geração de relatórios
- 🔗 Composição de múltiplas chamadas de API

```bash
sloth-runner -f examples/intermediate/api-integration.lua
```

**Funcionalidades demonstradas:**
- Busca e processamento de dados de APIs
- Validação de esquemas JSON
- Criação e atualização de recursos
- Relatórios automatizados

### 🐳 [`multi-container.lua`](./multi-container.lua)
**Orquestração de múltiplos containers**
- 🌐 Criação de redes Docker customizadas
- 💾 Gerenciamento de volumes persistentes
- 🔗 Comunicação entre containers
- 📊 Monitoramento de recursos
- 🧹 Limpeza automatizada

```bash
sloth-runner -f examples/intermediate/multi-container.lua
```

**Funcionalidades demonstradas:**
- Setup de ambiente multi-container
- Redis + Nginx com comunicação
- Health checks e monitoring
- Gerenciamento do ciclo de vida

### ⚡ [`parallel-processing.lua`](./parallel-processing.lua)
**Processamento paralelo eficiente**
- 🔄 Comparação sequencial vs paralelo
- 📡 Requisições HTTP paralelas
- 📁 Operações de arquivo simultâneas
- 🔢 Cálculos matemáticos distribuídos
- 📈 Análise de performance

```bash
sloth-runner -f examples/intermediate/parallel-processing.lua
```

**Funcionalidades demonstradas:**
- Execução paralela de tarefas independentes
- Benchmarking de performance
- Otimização de throughput
- Relatórios comparativos

## 🎓 Conceitos Avançados

### 🏗️ **Arquitetura**
- Composição de múltiplos módulos
- Padrões de configuração avançados
- Separação de responsabilidades
- Reutilização de componentes

### 🔄 **Fluxo de Dados**
- Pipeline de transformação de dados
- Validação em múltiplas etapas
- Propagação de resultados entre tarefas
- Estado compartilhado complexo

### ⚡ **Performance**
- Identificação de gargalos
- Paralelização eficiente
- Otimização de recursos
- Monitoramento de métricas

### 🛡️ **Confiabilidade**
- Tratamento de falhas graceful
- Retry com backoff
- Validação robusta de dados
- Logs estruturados

## 💡 Dicas Intermediárias

1. **Combine Módulos**: Use múltiplos módulos juntos
2. **Valide Dados**: Sempre valide entradas e saídas
3. **Use Estado**: Aproveite persistência para workflows complexos
4. **Monitore**: Acompanhe métricas e performance
5. **Paralelização**: Identifique tarefas independentes

## 🔧 Padrões Demonstrados

### 📊 **API Integration Pattern**
```lua
-- Configuração centralizada
local api_config = state.get("api_config")

-- HTTP com retry e validação
local result = http.get({
    url = config.base_url .. "/users",
    max_retries = 3,
    timeout = 10
})

-- Validação de esquemas
local validation = validate.schema(data, schema)
```

### 🐳 **Container Orchestration Pattern**
```lua
-- Rede personalizada
docker.exec({"network", "create", "app-network"})

-- Containers com dependências
local redis_result = docker.exec({
    "run", "-d", "--network", "app-network",
    "--name", "app-redis", "redis:alpine"
})
```

### ⚡ **Parallel Execution Pattern**
```lua
-- Tarefas paralelas
local tasks = {}
for i, url in ipairs(urls) do
    table.insert(tasks, {
        name = "request_" .. i,
        command = create_http_task(url)
    })
end

local results = parallel(tasks)
```

## 🚀 Casos de Uso Práticos

- **🔄 Data Pipeline**: Coleta, processamento e armazenamento
- **🧪 Testing Suite**: Testes automatizados multi-ambiente
- **📦 Build System**: Build paralelo com múltiplas targets
- **🌍 Multi-Service Deploy**: Deploy coordenado de microserviços
- **📊 Data Analytics**: Processamento paralelo de grandes datasets

## ➡️ Próximos Passos

Quando dominar estes exemplos:

- 🔴 [`../advanced/`](../advanced/) - Padrões de reliability e arquiteturas complexas
- 🌍 [`../real-world/`](../real-world/) - Casos de uso completos de produção
- 🔌 [`../integrations/`](../integrations/) - Integrações com serviços externos

## 📚 Recursos Adicionais

- **Performance Tuning**: Como otimizar workflows
- **Error Handling**: Padrões robustos de tratamento de erro
- **State Management**: Estratégias avançadas de estado
- **Module Composition**: Como combinar funcionalidades

---

**Keep Building! 🦥⚡**