# 🦥 Plugin Neovim

**Suporte completo para Sloth Runner DSL no Neovim/LunarVim**

O plugin Neovim do Sloth Runner oferece recursos completos de IDE para trabalhar com arquivos de workflow `.sloth`, incluindo syntax highlighting, autocompletar e execução integrada de tarefas.

## ✨ Recursos

### 🎨 Syntax Highlighting Avançado
- **Cores customizadas** para palavras-chave DSL, métodos e módulos
- **Interpolação de strings** destacada com sintaxe `${variavel}`
- **Detecção de caminhos** para arquivos de script e configuração
- **Variáveis de ambiente** em destaque
- **Suporte a comentários** com verificação ortográfica

### 📁 Detecção Inteligente de Arquivos
- Auto-detecta arquivos `.sloth` e aplica highlighting adequado
- Compatibilidade com extensão `.lua` para retrocompatibilidade
- Ícones personalizados (🦥) em exploradores de arquivo

### ⚡ Autocompletar de Código
- **Completion inteligente** para métodos DSL: `command`, `description`, `timeout`, etc.
- **Completion de módulos** para módulos internos: `exec`, `fs`, `net`, `aws`, etc.
- **Completion de funções** para padrões comuns: `task()`, `workflow.define()`

### 🔧 Executor Integrado
- **Executar workflows** diretamente do Neovim com `<leader>sr`
- **Listar tarefas** no arquivo atual com `<leader>sl`
- **Suporte a dry-run** para testar workflows

### 📋 Snippets & Templates de Código
- **Criação rápida de tarefas** com abreviação `_task`
- **Templates de workflow** com abreviação `_workflow`
- **Templates de função** com abreviação `_cmd`
- **Templates auto-gerados** para novos arquivos `.sloth`

### 🎯 Text Objects & Navegação
- **Selecionar blocos de tarefa** com `vit` (visual in task)
- **Selecionar blocos de workflow** com `viw` (visual in workflow)
- **Dobramento inteligente** para seções de código recolhíveis
- **Indentação inteligente** para encadeamento DSL

## 🚀 Configuração Rápida

### Para Usuários do LunarVim

Adicione ao seu `~/.config/lvim/config.lua`:

```lua
-- Desabilitar formatação automática (recomendado)
lvim.format_on_save.enabled = false

-- Configurar ícones de arquivos sloth
require('nvim-web-devicons').setup {
  override_by_extension = {
    ["sloth"] = {
      icon = "🦥",
      color = "#8B4513",
      name = "SlothDSL"
    }
  }
}

-- Mapeamentos de teclas para sloth runner
lvim.keys.normal_mode["<leader>sr"] = function()
  local file = vim.api.nvim_buf_get_name(0)
  if file:match("%.sloth$") then
    vim.cmd("split | terminal sloth-runner run -f " .. vim.fn.shellescape(file))
  end
end

lvim.keys.normal_mode["<leader>sl"] = function()
  local file = vim.api.nvim_buf_get_name(0)
  if file:match("%.sloth$") then
    vim.cmd("split | terminal sloth-runner list -f " .. vim.fn.shellescape(file))
  end
end

-- Comando de formatação manual
lvim.keys.normal_mode["<leader>sf"] = ":SlothFormat<CR>"
```

### Para Neovim Padrão

Usando [lazy.nvim](https://github.com/folke/lazy.nvim):

```lua
{
  dir = "/caminho/para/sloth-runner/nvim-plugin",
  name = "sloth-runner",
  ft = { "sloth" },
  config = function()
    require("sloth-runner").setup({
      runner = {
        command = "sloth-runner",
        keymaps = {
          run_file = "<leader>sr",
          list_tasks = "<leader>sl",
          dry_run = "<leader>st",
          debug = "<leader>sd",
        }
      },
      completion = {
        enable = true,
        snippets = true,
      },
      folding = {
        enable = true,
        auto_close = false,
      }
    })
  end,
}
```

## 📝 Mapeamentos de Teclas

| Tecla | Ação | Descrição |
|-------|------|-----------|
| `<leader>sr` | Executar Arquivo | Executa o workflow `.sloth` atual |
| `<leader>sl` | Listar Tarefas | Mostra todas as tarefas no arquivo atual |
| `<leader>st` | Dry Run | Testa workflow sem execução |
| `<leader>sd` | Debug | Executa com saída de debug |
| `<leader>sf` | Formatar | Formata arquivo atual (manual) |

## 🎨 Snippets de Código

### Criação Rápida de Tarefa
Digite `_task` e pressione Tab:

```lua
local nome_tarefa = task("")
    :description("")
    :command(function(params, deps)
        -- TODO: implementar
        return true
    end)
    :build()
```

### Criação Rápida de Workflow
Digite `_workflow` e pressione Tab:

```lua
workflow.define("", {
    description = "",
    version = "1.0.0",
    tasks = {
        -- tarefas aqui
    }
})
```

### Função de Comando Rápida
Digite `_cmd` e pressione Tab:

```lua
:command(function(params, deps)
    -- TODO: implementar
    return true
end)
```

## 🔧 Configuração Avançada

### Syntax Highlighting Personalizado

O plugin fornece grupos de highlight personalizados:

```lua
-- Personalizar cores (adicione à sua config)
vim.api.nvim_create_autocmd("ColorScheme", {
  pattern = "*",
  callback = function()
    vim.api.nvim_set_hl(0, "slothDSLKeyword", { fg = '#8B4513', bold = true })
    vim.api.nvim_set_hl(0, "slothModule", { fg = '#228B22', bold = true })
    vim.api.nvim_set_hl(0, "slothMethod", { fg = '#4682B4' })
    vim.api.nvim_set_hl(0, "slothFunction", { fg = '#DAA520' })
    vim.api.nvim_set_hl(0, "slothEnvVar", { fg = '#DC143C', bold = true })
    vim.api.nvim_set_hl(0, "slothPath", { fg = '#20B2AA' })
  end
})
```

### Integração com Árvore de Arquivos

Para usuários do nvim-tree:

```lua
require("nvim-tree").setup({
  renderer = {
    icons = {
      glyphs = {
        extension = {
          sloth = "🦥"
        }
      }
    }
  }
})
```

## 🛠️ Instalação Manual

1. **Clone ou copie os arquivos do plugin:**
   ```bash
   cp -r /caminho/para/sloth-runner/nvim-plugin ~/.config/nvim/
   ```

2. **Adicione à sua configuração do Neovim:**
   ```lua
   -- Adicione ao init.lua ou init.vim
   vim.opt.runtimepath:append("~/.config/nvim/nvim-plugin")
   ```

3. **Reinicie o Neovim e abra um arquivo `.sloth`**

## 🐛 Solução de Problemas

### Syntax Highlighting Não Funciona
- Certifique-se de que o arquivo tem extensão `.sloth`
- Execute `:set filetype=sloth` manualmente se necessário
- Verifique se os arquivos do plugin estão no local correto

### Mapeamentos de Teclas Não Funcionam
- Verifique se `sloth-runner` está no seu PATH
- Verifique se as teclas estão conflitando com outros plugins
- Use `:map <leader>sr` para verificar se o mapeamento existe

### Autocompletar Não Aparece
- Certifique-se de que completion está habilitado: `:set completeopt=menu,menuone,noselect`
- Tente disparar manualmente com `<C-x><C-o>`
- Verifique se omnifunc está definido: `:set omnifunc?`

### Problemas de Formatação
- **Formatação automática está desabilitada por padrão** para evitar conflitos
- Use formatação manual: `<leader>sf` ou `:SlothFormat`
- Para formatação com stylua, certifique-se de que está instalado e configurado

## 📖 Exemplos

### Arquivo de Workflow Básico

```lua
-- deployment.sloth
local tarefa_deploy = task("deploy_app")
    :description("Deploy da aplicação para produção")
    :command(function(params, deps)
        local resultado = exec.run("kubectl apply -f deployment.yaml")
        if not resultado.success then
            log.error("Deploy falhou: " .. resultado.stderr)
            return false
        end
        
        log.info("🚀 Deploy realizado com sucesso!")
        return true
    end)
    :timeout(300)
    :retries(3)
    :build()

workflow.define("deployment_producao", {
    description = "Workflow de deploy para produção",
    version = "1.0.0",
    tasks = { tarefa_deploy }
})
```

Com o plugin instalado, este arquivo terá:
- **Syntax highlighting** para palavras-chave, funções e strings
- **Autocompletar** ao digitar nomes de métodos
- **Execução rápida** com `<leader>sr`
- **Listagem de tarefas** com `<leader>sl`

## 🚀 Próximos Passos

- **Aprenda a DSL:** Veja [Conceitos Fundamentais](../pt/core-concepts.md)
- **Teste Exemplos:** Consulte [Guia de Exemplos](../pt/advanced-examples.md)
- **Recursos Avançados:** Explore [Recursos Avançados](../pt/advanced-features.md)
- **Referência da API:** Leia [Documentação da API Lua](../pt/lua-api-overview.md)

---

O plugin Neovim torna a criação de workflows Sloth muito mais fácil com suporte completo de IDE. Comece a criar automações poderosas com confiança! 🦥✨