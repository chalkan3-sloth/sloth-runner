# 🦥 Sloth Runner Examples

Esta pasta contém uma coleção abrangente de exemplos que demonstram as capacidades do Sloth Runner, organizados por nível de complexidade.

## 📁 Estrutura dos Exemplos

### 🟢 [`beginner/`](./beginner/) - Exemplos para Iniciantes
Exemplos simples e diretos que introduzem conceitos básicos:
- ✨ Tarefas básicas e comandos simples
- 🔄 Workflow linear básico
- 📊 Estado simples
- 🛠️ Ferramentas básicas (exec, fs, log)

### 🟡 [`intermediate/`](./intermediate/) - Exemplos Intermediários
Exemplos que combinam múltiplas funcionalidades:
- 🌐 Integração com APIs HTTP
- 🐳 Automação Docker
- ☁️ Operações na nuvem básicas
- 🔄 Workflows condicionais
- ⚡ Execução paralela

### 🔴 [`advanced/`](./advanced/) - Exemplos Avançados
Exemplos complexos que demonstram recursos avançados:
- 🏗️ Arquiteturas distribuídas
- 🛡️ Padrões de reliability e retry
- 🔐 Gerenciamento de segredos
- 📊 Monitoramento e métricas
- 🚀 Pipelines CI/CD complexos

### 🌍 [`real-world/`](./real-world/) - Casos de Uso Reais
Exemplos práticos do mundo real:
- 🚀 Deploy de aplicações
- 🏗️ Infraestrutura como código
- 📦 Build e release pipelines
- 🔄 Data processing workflows
- 🏥 Health checks e monitoring

### 🔌 [`integrations/`](./integrations/) - Integrações
Exemplos de integração com serviços externos:
- ☁️ AWS, Azure, GCP
- 🐳 Docker & Kubernetes
- 📊 Bancos de dados
- 📧 Notificações e alertas
- 🔧 Ferramentas DevOps

## 🚀 Como Executar os Exemplos

### Requisitos
```bash
# Instalar Sloth Runner
./install.sh

# Ou compilar do código fonte
go build -o sloth-runner ./cmd/sloth-runner
```

### Execução Básica
```bash
# Executar um exemplo específico
sloth-runner -f examples/beginner/hello-world.lua

# Executar com parâmetros
sloth-runner -f examples/intermediate/http-api.lua -p "{'api_key': 'your-key'}"

# Modo debug para ver detalhes
sloth-runner -f examples/advanced/retry-patterns.lua --debug
```

### Sistema de Help Integrado
```lua
-- Dentro de qualquer script
help()                    -- Ajuda geral
help.modules()           -- Lista todos os módulos disponíveis
help.module("http")      -- Help do módulo HTTP
help.search("docker")    -- Busca por funcionalidades
help.examples("http")    -- Exemplos do módulo HTTP
```

## 📚 Exemplos por Módulo

### 🌐 HTTP Client
- [`beginner/http-basics.lua`](./beginner/http-basics.lua) - GET/POST básicos
- [`intermediate/api-integration.lua`](./intermediate/api-integration.lua) - Integração com API
- [`advanced/http-reliability.lua`](./advanced/http-reliability.lua) - Retry e circuit breaker

### 🐳 Docker
- [`beginner/docker-basics.lua`](./beginner/docker-basics.lua) - Build e run básicos
- [`intermediate/multi-container.lua`](./intermediate/multi-container.lua) - Multi-container setup
- [`advanced/docker-compose.lua`](./advanced/docker-compose.lua) - Orchestration completa

### ☁️ Cloud Providers
- [`intermediate/aws-s3.lua`](./intermediate/aws-s3.lua) - Operações S3
- [`intermediate/gcp-storage.lua`](./intermediate/gcp-storage.lua) - Google Cloud Storage
- [`advanced/multi-cloud.lua`](./advanced/multi-cloud.lua) - Deploy multi-cloud

### 💾 State Management
- [`beginner/state-basics.lua`](./beginner/state-basics.lua) - Operações básicas de estado
- [`intermediate/distributed-state.lua`](./intermediate/distributed-state.lua) - Estado distribuído
- [`advanced/state-patterns.lua`](./advanced/state-patterns.lua) - Padrões avançados

## 🎯 Exemplos por Caso de Uso

### 🚀 CI/CD Pipelines
- [`real-world/nodejs-cicd.lua`](./real-world/nodejs-cicd.lua) - Pipeline Node.js completo
- [`real-world/go-microservice.lua`](./real-world/go-microservice.lua) - Deploy de microserviço Go
- [`real-world/frontend-deploy.lua`](./real-world/frontend-deploy.lua) - Deploy de aplicação React

### 🏗️ Infrastructure as Code
- [`real-world/terraform-aws.lua`](./real-world/terraform-aws.lua) - Infraestrutura AWS
- [`real-world/pulumi-kubernetes.lua`](./real-world/pulumi-kubernetes.lua) - Deploy Kubernetes
- [`real-world/monitoring-stack.lua`](./real-world/monitoring-stack.lua) - Stack de monitoramento

### 📊 Data Processing
- [`real-world/etl-pipeline.lua`](./real-world/etl-pipeline.lua) - Pipeline ETL
- [`real-world/data-validation.lua`](./real-world/data-validation.lua) - Validação de dados
- [`real-world/backup-restore.lua`](./real-world/backup-restore.lua) - Backup e restore

## 💡 Dicas para Aprender

1. **Comece pelo Básico**: Inicie pelos exemplos em `beginner/`
2. **Use o Help**: `help()` é seu melhor amigo
3. **Experimente**: Modifique os exemplos para entender melhor
4. **Debug Mode**: Use `--debug` para ver o que acontece internamente
5. **Combine Módulos**: Os exemplos avançados mostram como combinar funcionalidades

## 🤝 Contribuindo

Quer adicionar um exemplo? Siga estas diretrizes:

1. **Escolha a Categoria Certa**: beginner, intermediate, advanced, real-world, integrations
2. **Documente Bem**: Comentários claros e README quando necessário
3. **Teste Tudo**: Certifique-se que o exemplo funciona
4. **Siga o Padrão**: Use a estrutura similar aos exemplos existentes

## 📞 Suporte

- 📚 **Documentação**: [docs/](../docs/)
- 🐛 **Issues**: [GitHub Issues](https://github.com/chalkan3/sloth-runner/issues)
- 💬 **Discussões**: [GitHub Discussions](https://github.com/chalkan3/sloth-runner/discussions)
- 📧 **Email**: support@sloth-runner.dev

---

**Happy Automating! 🦥✨**