# 🔌 Desenvolvimento de Plugins

**Construindo Extensões para a Plataforma Sloth Runner**

O Sloth Runner fornece um sistema de plugins poderoso que permite aos desenvolvedores estender a plataforma com funcionalidades personalizadas. Este guia cobre tudo que você precisa saber para desenvolver seus próprios plugins.

## 🏗️ Arquitetura de Plugins

### Tipos de Plugins

O Sloth Runner suporta vários tipos de plugins:

1. **🌙 Módulos Lua** - Estendem a DSL com novas funções e capacidades
2. **⚡ Processadores de Comando** - Adicionam novos comandos CLI e operações
3. **🎨 Extensões de UI** - Melhoram o dashboard web e interface
4. **🔗 Integrações** - Conectam com ferramentas e serviços externos
5. **🦥 Plugins de Editor** - Extensões para IDE/Editor (como nosso plugin Neovim)

### Componentes Principais

```
sloth-runner/
├── plugins/
│   ├── lua-modules/       # Extensões DSL Lua
│   ├── commands/          # Plugins de comando CLI
│   ├── ui/               # Extensões de UI web
│   ├── integrations/     # Integrações de terceiros
│   └── editors/          # Plugins de editor/IDE
└── internal/
    └── plugin/           # Core do sistema de plugins
```

## 🌙 Desenvolvendo Plugins de Módulo Lua

### Estrutura Básica

Crie um novo módulo Lua que estende a DSL:

```lua
-- plugins/lua-modules/meu-modulo/init.lua
local M = {}

-- Metadados do módulo
M._NAME = "meu-modulo"
M._VERSION = "1.0.0"
M._DESCRIPTION = "Funcionalidade customizada para Sloth Runner"

-- API Pública
function M.ola(nome)
    return string.format("Olá, %s do meu módulo customizado!", nome or "Mundo")
end

function M.tarefa_customizada(config)
    return {
        execute = function(params)
            log.info("🔌 Executando tarefa customizada: " .. config.name)
            -- Lógica da tarefa customizada aqui
            return true
        end,
        validate = function()
            return config.name ~= nil
        end
    }
end

-- Registrar funções do módulo
function M.register()
    -- Tornar funções disponíveis na DSL
    _G.meu_modulo = M
    
    -- Registrar tipo de tarefa customizada
    task.register_type("customizada", M.tarefa_customizada)
end

return M
```

### Usando Módulos Customizados em Workflows

```lua
-- workflow.sloth
local minha_tarefa = task("teste_customizado")
    :type("customizada", { name = "teste" })
    :description("Testando plugin customizado")
    :build()

-- Uso direto do módulo
local saudacao = meu_modulo.ola("Desenvolvedor")
log.info(saudacao)

workflow.define("teste_plugin", {
    description = "Testando plugin customizado",
    tasks = { minha_tarefa }
})
```

### Registro de Plugin

Crie um manifesto de plugin:

```yaml
# plugins/lua-modules/meu-modulo/plugin.yaml
name: meu-modulo
version: 1.0.0
description: Funcionalidade customizada para Sloth Runner
type: lua-module
author: Seu Nome
license: MIT

entry_point: init.lua
dependencies:
  - sloth-runner: ">=1.0.0"

permissions:
  - filesystem.read
  - network.http
  - system.exec

configuration:
  settings:
    api_key:
      type: string
      required: false
      description: "Chave da API para serviço externo"
```

## ⚡ Desenvolvimento de Plugin de Comando

### Estrutura de Comando CLI

```go
// plugins/commands/meu-comando/main.go
package main

import (
    "github.com/spf13/cobra"
    "github.com/chalkan3-sloth/sloth-runner/pkg/plugin"
)

type MeuComandoPlugin struct {
    config *MinhaConfig
}

type MinhaConfig struct {
    Configuracao1 string `json:"configuracao1"`
    Configuracao2 int    `json:"configuracao2"`
}

func (p *MeuComandoPlugin) Name() string {
    return "meu-comando"
}

func (p *MeuComandoPlugin) Command() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "meu-comando",
        Short: "Funcionalidade de comando customizado",
        Long:  "Descrição estendida do comando customizado",
        RunE:  p.execute,
    }
    
    cmd.Flags().StringVar(&p.config.Configuracao1, "config1", "", "Configuração customizada")
    cmd.Flags().IntVar(&p.config.Configuracao2, "config2", 0, "Outra configuração")
    
    return cmd
}

func (p *MeuComandoPlugin) execute(cmd *cobra.Command, args []string) error {
    log.Info("🔌 Executando comando customizado com configurações:", 
        "config1", p.config.Configuracao1,
        "config2", p.config.Configuracao2)
    
    // Lógica do comando customizado aqui
    return nil
}

func main() {
    plugin := &MeuComandoPlugin{
        config: &MinhaConfig{},
    }
    
    plugin.Register()
}
```

## 🛠️ Ferramentas de Desenvolvimento de Plugin

### Gerador de Plugin

Crie novos plugins rapidamente com o gerador:

```bash
# Gerar um novo plugin de módulo Lua
sloth-runner plugin generate --type=lua-module --name=meu-modulo

# Gerar um plugin de comando CLI
sloth-runner plugin generate --type=command --name=meu-comando

# Gerar uma extensão de UI
sloth-runner plugin generate --type=ui --name=meu-dashboard
```

### Ambiente de Desenvolvimento

```bash
# Iniciar servidor de desenvolvimento com hot-reload de plugin
sloth-runner dev --plugins-dir=./plugins

# Testar plugin localmente
sloth-runner plugin test ./plugins/meu-plugin

# Construir plugin para distribuição
sloth-runner plugin build ./plugins/meu-plugin --output=dist/
```

### Teste de Plugin

```go
// plugins/meu-plugin/plugin_test.go
package main

import (
    "testing"
    "github.com/chalkan3-sloth/sloth-runner/pkg/plugin/testing"
)

func TestMeuPlugin(t *testing.T) {
    // Criar ambiente de teste
    env := plugintest.NewEnvironment(t)
    
    // Carregar plugin
    plugin, err := env.LoadPlugin("./")
    if err != nil {
        t.Fatal(err)
    }
    
    // Testar funcionalidade do plugin
    result, err := plugin.Execute(map[string]interface{}{
        "parametro_teste": "valor",
    })
    
    if err != nil {
        t.Fatal(err)
    }
    
    // Verificar resultados
    if result.Status != "success" {
        t.Errorf("Esperado sucesso, obtido %s", result.Status)
    }
}
```

## 📦 Distribuição de Plugin

### Registro de Plugin

Publique seu plugin no registro de plugins do Sloth Runner:

```bash
# Login no registro
sloth-runner registry login

# Publicar plugin
sloth-runner plugin publish ./meu-plugin

# Instalar plugin publicado
sloth-runner plugin install meu-usuario/meu-plugin
```

### Marketplace de Plugin

Navegue e descubra plugins:

```bash
# Buscar plugins
sloth-runner plugin search "kubernetes"

# Obter informações do plugin
sloth-runner plugin info kubernetes-operator

# Instalar do marketplace
sloth-runner plugin install --marketplace kubernetes-operator
```

## 🔒 Segurança e Melhores Práticas

### Diretrizes de Segurança

1. **🛡️ Princípio do Menor Privilégio** - Solicite apenas as permissões necessárias
2. **🔐 Validação de Entrada** - Sempre valide entrada do usuário e configuração
3. **🚫 Evitar Estado Global** - Mantenha o estado do plugin isolado
4. **📝 Tratamento de Erros** - Forneça mensagens de erro claras e logging
5. **🧪 Testes** - Escreva testes abrangentes para toda a funcionalidade

### Qualidade de Código

```go
// Bom: Tratamento claro de erros
func (p *MeuPlugin) Execute(params map[string]interface{}) (*Result, error) {
    value, ok := params["parametro_obrigatorio"].(string)
    if !ok {
        return nil, fmt.Errorf("parametro_obrigatorio deve ser uma string")
    }
    
    if value == "" {
        return nil, fmt.Errorf("parametro_obrigatorio não pode estar vazio")
    }
    
    // Processar com entrada validada
    result := p.process(value)
    return result, nil
}
```

### Padrões de Documentação

Cada plugin deve incluir:

- **📋 README.md** - Instruções de instalação e uso
- **📚 Documentação da API** - Documentação de função/método
- **📖 Exemplos** - Exemplos de código funcionais
- **🧪 Testes** - Testes unitários e de integração
- **📄 Licença** - Informações claras de licenciamento

## 📚 Exemplos e Templates

### Exemplo Completo de Plugin

Confira estes plugins de exemplo:

- **[Plugin Kubernetes Operator](https://github.com/sloth-runner/plugin-kubernetes)** - Gerenciar recursos K8s
- **[Plugin Integração Slack](https://github.com/sloth-runner/plugin-slack)** - Enviar notificações
- **[Plugin Dashboard Monitoramento](https://github.com/sloth-runner/plugin-monitoring)** - UI de métricas customizadas

### Templates de Plugin

Use templates oficiais para início rápido:

```bash
# Usar template
sloth-runner plugin init --template=lua-module meu-plugin
sloth-runner plugin init --template=go-command meu-comando
sloth-runner plugin init --template=react-ui meu-dashboard
```

## 💬 Comunidade e Suporte

### Obtendo Ajuda

- **📖 [Documentação da API de Plugin](https://docs.sloth-runner.io/plugin-api)**
- **💬 [Comunidade Discord](https://discord.gg/sloth-runner)** - #plugin-development
- **🐛 [Issues do GitHub](https://github.com/chalkan3-sloth/sloth-runner/issues)** - Relatórios de bug e solicitações de recurso
- **📧 [Lista de Email de Plugin](mailto:plugins@sloth-runner.io)** - Discussões de desenvolvimento

### Contribuindo

Recebemos contribuições de plugin! Veja nosso [Guia de Contribuição](contributing.md) para detalhes sobre:

- Processo de submissão de plugin
- Diretrizes de revisão de código
- Requisitos de documentação
- Padrões de teste

---

Comece a construir plugins incríveis para o Sloth Runner hoje! A arquitetura extensível da plataforma torna fácil adicionar exatamente a funcionalidade que você precisa. 🔌✨