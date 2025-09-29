# 🟢 Exemplos para Iniciantes

Esta pasta contém exemplos simples e didáticos para quem está começando com o Sloth Runner.

## 📋 Exemplos Disponíveis

### 🌟 [`hello-world.lua`](./hello-world.lua)
**Seu primeiro script Sloth Runner**
- ✨ Estrutura básica de um TaskDefinitions
- 📝 Uso do sistema de logging
- 🔗 Dependências entre tarefas
- 📁 Operações básicas de arquivo
- ✅ Tratamento de retorno de funções

```bash
sloth-runner -f examples/beginner/hello-world.lua
```

### 🌐 [`http-basics.lua`](./http-basics.lua) 
**Primeiras requisições HTTP**
- 📡 GET requests simples
- 🔧 Headers customizados
- 📤 POST com dados JSON
- 🚨 Tratamento de erros HTTP
- ⏰ Configuração de timeouts

```bash
sloth-runner -f examples/beginner/http-basics.lua
```

### 💾 [`state-basics.lua`](./state-basics.lua)
**Gerenciamento de estado básico**
- 📝 Operações set, get, delete
- 🎯 Valores padrão
- 💿 Persistência entre execuções
- 🔢 Contadores e timestamps
- 🧹 Limpeza de dados

```bash
sloth-runner -f examples/beginner/state-basics.lua
```

### 🐳 [`docker-basics.lua`](./docker-basics.lua)
**Primeiros passos com Docker**
- ✅ Verificação de instalação
- 📋 Listagem de containers
- ⬇️ Pull de imagens
- 🚀 Execução de containers
- 🧹 Limpeza de recursos

```bash
sloth-runner -f examples/beginner/docker-basics.lua
```

## 🎓 Conceitos Aprendidos

Após executar todos os exemplos desta seção, você terá aprendido:

### 📚 **Conceitos Fundamentais**
- Como estruturar um arquivo de configuração Lua
- Sistema de dependências entre tarefas
- Uso do sistema de logging integrado
- Persistência de estado entre execuções

### 🛠️ **Ferramentas Básicas**
- Módulo `fs` para operações de arquivo
- Módulo `exec` para executar comandos
- Módulo `log` para logging estruturado
- Módulo `state` para persistência de dados

### 🌐 **Integração Externa**
- Fazer requisições HTTP GET/POST
- Configurar headers e timeouts
- Tratar erros de rede
- Trabalhar com dados JSON

### 🐳 **Containers**
- Verificar instalação do Docker
- Listar e gerenciar containers
- Baixar imagens do registry
- Executar containers de forma controlada

## 💡 Dicas para Iniciantes

1. **Comece Simples**: Execute `hello-world.lua` primeiro
2. **Leia os Comentários**: Cada exemplo tem comentários explicativos
3. **Use o Help**: `help()` no início de qualquer script
4. **Experimente**: Modifique os exemplos para entender melhor
5. **Debug Mode**: Use `--debug` para ver detalhes internos

## ➡️ Próximos Passos

Quando se sentir confortável com estes exemplos, passe para:

- 🟡 [`../intermediate/`](../intermediate/) - Exemplos que combinam múltiplas funcionalidades
- 📚 [`../README.md`](../README.md) - Visão geral completa de todos os exemplos

## 🤝 Precisa de Ajuda?

- 📖 **Help Integrado**: Execute `help()` em qualquer script
- 📚 **Documentação**: Veja a pasta `docs/`
- 💬 **Comunidade**: GitHub Discussions
- 🐛 **Issues**: GitHub Issues

---

**Happy Learning! 🦥✨**