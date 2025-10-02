# 🚀 Módulo Goroutine - Índice Completo

## 📚 Documentação Disponível

### 🎯 Para Começar
1. **[Quick Start](../../GOROUTINE_QUICKSTART.md)** - Comece em 30 segundos
   - Instalação rápida
   - Exemplo básico
   - Casos de uso comuns

### 📖 Documentação Completa
2. **[Referência Completa](goroutine.md)** - API detalhada
   - Todas as funções documentadas
   - Parâmetros e retornos
   - Exemplos para cada função
   - Melhores práticas
   - Performance e limitações
   - Troubleshooting

### 🎨 README Destacado
3. **[README](GOROUTINE_README.md)** - Visão geral
   - Características principais
   - Casos de uso reais (ETL, Health Check, CI/CD)
   - Benchmarks de performance
   - Boas práticas vs anti-patterns
   - Troubleshooting avançado

### 📝 Informações Técnicas
4. **[Resumo da Implementação](../../GOROUTINE_MODULE_SUMMARY.md)**
   - Arquitetura interna
   - Estrutura de classes
   - Fluxo de execução
   - Segurança e thread-safety
   - Checklist de implementação

## 🎯 Exemplos Práticos

### Exemplos Disponíveis
- **test_goroutine.sloth** - 7 exemplos completos
  - Spawn simples e múltiplo
  - Worker pools
  - Async/await
  - WaitGroups
  - Timeouts
  - Processamento paralelo

- **test_goroutine_simple.sloth** - Teste rápido
  - Testes básicos
  - Validação de funcionalidade

## 🔗 Links Rápidos

| Recurso | Descrição | Link |
|---------|-----------|------|
| 📖 API Reference | Documentação técnica completa | [goroutine.md](goroutine.md) |
| 🚀 Quick Start | Comece em 30 segundos | [GOROUTINE_QUICKSTART.md](../../GOROUTINE_QUICKSTART.md) |
| 📝 README | Visão geral e exemplos | [GOROUTINE_README.md](GOROUTINE_README.md) |
| 🏗️ Implementação | Detalhes técnicos | [GOROUTINE_MODULE_SUMMARY.md](../../GOROUTINE_MODULE_SUMMARY.md) |
| 🎬 Exemplos | Código executável | Ver arquivos .sloth |

## 📊 Funcionalidades

### ✅ O Que Você Pode Fazer

| Funcionalidade | Descrição | Função Principal |
|----------------|-----------|------------------|
| **Spawn** | Executar funções em goroutines | `goroutine.spawn()` |
| **Worker Pools** | Gerenciar pools de workers | `goroutine.pool_*()` |
| **Async/Await** | Programação assíncrona | `goroutine.async()`, `await()` |
| **Sincronização** | WaitGroups | `goroutine.wait_group()` |
| **Timeout** | Limites de tempo | `goroutine.timeout()` |
| **Sleep** | Pausar execução | `goroutine.sleep()` |

### 🎯 Casos de Uso

1. **Processamento Paralelo de Dados**
   - Processar arquivos em lote
   - ETL paralelo
   - Transformação de dados

2. **Operações I/O Assíncronas**
   - Múltiplas requisições HTTP
   - Health checks distribuídos
   - Download/upload paralelo

3. **Pipelines Complexos**
   - CI/CD paralelo
   - Build multi-stage
   - Deploy distribuído

4. **Monitoramento e Alertas**
   - Check de múltiplos serviços
   - Coleta de métricas
   - Aggregação de dados

## 🚀 Como Usar

### Passo 1: Importar
```lua
local goroutine = require("goroutine")
```

### Passo 2: Escolher Padrão

#### Para Processamento em Lote → Use Worker Pool
```lua
goroutine.pool_create("mypool", { workers = 10 })
-- submeter tarefas
goroutine.pool_wait("mypool")
goroutine.pool_close("mypool")
```

#### Para Operações I/O → Use Async/Await
```lua
local handle = goroutine.async(function()
    return fetch_data()
end)
local success, data = goroutine.await(handle)
```

#### Para Fire-and-Forget → Use Spawn
```lua
goroutine.spawn(function()
    background_task()
end)
```

#### Para Sincronização → Use WaitGroup
```lua
local wg = goroutine.wait_group()
wg:add(n)
-- spawn tasks
wg:wait()
```

## 📈 Performance

| Métrica | Valor |
|---------|-------|
| Overhead por goroutine | ~2-3 KB |
| Tempo de criação | ~100 ns |
| Context switch | ~200 ns |
| Goroutines simultâneas | Milhares |

## 🎓 Tutoriais

### Nível Iniciante
1. Leia o [Quick Start](../../GOROUTINE_QUICKSTART.md)
2. Execute `test_goroutine_simple.sloth`
3. Modifique o exemplo para seu caso de uso

### Nível Intermediário
1. Estude a [Referência Completa](goroutine.md)
2. Execute `test_goroutine.sloth` (todos os exemplos)
3. Implemente worker pools em seus workflows

### Nível Avançado
1. Leia o [Resumo da Implementação](../../GOROUTINE_MODULE_SUMMARY.md)
2. Estude os padrões no [README](GOROUTINE_README.md)
3. Crie pipelines complexos com múltiplos padrões

## 🛠️ Desenvolvimento

### Arquivos Fonte
- **Implementação**: `internal/modules/core/goroutine.go`
- **Registro**: `internal/modules/init.go`
- **Documentação**: `docs/modules/goroutine.md`

### Como Contribuir
1. Fork o repositório
2. Crie uma branch para sua feature
3. Adicione testes
4. Faça PR

## 🐛 Problemas Comuns

| Problema | Solução | Link |
|----------|---------|------|
| Pool não encontrado | Criar antes de usar | [Troubleshooting](goroutine.md#troubleshooting) |
| Tarefas não terminam | Usar timeout | [Troubleshooting](goroutine.md#troubleshooting) |
| Race conditions | Usar WaitGroup | [Troubleshooting](goroutine.md#troubleshooting) |
| Memory leak | Fechar pools | [Best Practices](goroutine.md#melhores-práticas) |

## 📞 Suporte

- 🐛 **Issues**: [GitHub Issues](https://github.com/chalkan3-sloth/sloth-runner/issues)
- 📖 **Documentação**: Este índice
- 💬 **Comunidade**: Discord (em breve)

## ✅ Status

- ✅ **Implementação**: Completa
- ✅ **Testes**: Funcionais
- ✅ **Documentação**: Completa
- ✅ **Exemplos**: Disponíveis
- ✅ **Performance**: Otimizada
- ✅ **Thread-Safety**: Garantida

---

**Última atualização**: 2025-01-10  
**Versão**: 1.0.0  
**Autor**: Sloth Runner Team
