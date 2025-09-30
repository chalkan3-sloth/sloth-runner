# 🦥 Neovim Plugin

**IDE-grade support for Sloth Runner DSL in Neovim/LunarVim**

The Sloth Runner Neovim plugin provides comprehensive IDE features for working with `.sloth` workflow files, including syntax highlighting, code completion, and integrated task execution.

## ✨ Features

### 🎨 Rich Syntax Highlighting
- **Custom colors** for DSL keywords, methods, and modules
- **String interpolation** highlighting with `${variable}` syntax
- **File path detection** for script and configuration files
- **Environment variable** highlighting
- **Comment support** with proper spell checking

### 📁 Smart File Detection
- Auto-detects `.sloth` files and applies proper highlighting
- Backward compatibility with `.lua` extension
- Custom file icons (🦥) in file explorers

### ⚡ Code Completion
- **Intelligent completion** for DSL methods: `command`, `description`, `timeout`, etc.
- **Module completion** for built-in modules: `exec`, `fs`, `net`, `aws`, etc.
- **Function completion** for common patterns: `task()`, `workflow.define()`

### 🔧 Integrated Runner
- **Run workflows** directly from Neovim with `<leader>sr`
- **List tasks** in current file with `<leader>sl`
- **Dry-run support** for testing workflows

### 📋 Code Snippets & Templates
- **Quick task creation** with `_task` abbreviation
- **Workflow templates** with `_workflow` abbreviation
- **Function templates** with `_cmd` abbreviation
- **Auto-generated templates** for new `.sloth` files

### 🎯 Text Objects & Navigation
- **Select task blocks** with `vit` (visual in task)
- **Select workflow blocks** with `viw` (visual in workflow)
- **Smart folding** for collapsible code sections
- **Intelligent indentation** for DSL chaining

## 🚀 Quick Setup

### For LunarVim Users

Add to your `~/.config/lvim/config.lua`:

```lua
-- Disable auto-formatting (recommended)
lvim.format_on_save.enabled = false

-- Configure sloth file icons
require('nvim-web-devicons').setup {
  override_by_extension = {
    ["sloth"] = {
      icon = "🦥",
      color = "#8B4513",
      name = "SlothDSL"
    }
  }
}

-- Key mappings for sloth runner
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

-- Manual formatting command
lvim.keys.normal_mode["<leader>sf"] = ":SlothFormat<CR>"
```

### For Standard Neovim

Using [lazy.nvim](https://github.com/folke/lazy.nvim):

```lua
{
  dir = "/path/to/sloth-runner/nvim-plugin",
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

## 📝 Key Mappings

| Key | Action | Description |
|-----|--------|-------------|
| `<leader>sr` | Run File | Execute current `.sloth` workflow |
| `<leader>sl` | List Tasks | Show all tasks in current file |
| `<leader>st` | Dry Run | Test workflow without execution |
| `<leader>sd` | Debug | Run with debug output |
| `<leader>sf` | Format | Format current file (manual) |

## 🎨 Code Snippets

### Quick Task Creation
Type `_task` and press Tab:

```lua
local task_name = task("")
    :description("")
    :command(function(params, deps)
        -- TODO: implement
        return true
    end)
    :build()
```

### Quick Workflow Creation
Type `_workflow` and press Tab:

```lua
workflow.define("", {
    description = "",
    version = "1.0.0",
    tasks = {
        -- tasks here
    }
})
```

### Quick Command Function
Type `_cmd` and press Tab:

```lua
:command(function(params, deps)
    -- TODO: implement
    return true
end)
```

## 🔧 Advanced Configuration

### Custom Syntax Highlighting

The plugin provides custom highlight groups:

```lua
-- Customize colors (add to your config)
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

### File Tree Integration

For nvim-tree users:

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

### Status Line Integration

Show sloth icon in status line for `.sloth` files:

```lua
vim.api.nvim_create_autocmd({"BufEnter", "BufWinEnter"}, {
  pattern = "*.sloth",
  callback = function()
    vim.b.sloth_file = true
    -- Your status line will pick this up
  end
})
```

## 🛠️ Manual Installation

1. **Clone or copy the plugin files:**
   ```bash
   cp -r /path/to/sloth-runner/nvim-plugin ~/.config/nvim/
   ```

2. **Add to your Neovim configuration:**
   ```lua
   -- Add to init.lua or init.vim
   vim.opt.runtimepath:append("~/.config/nvim/nvim-plugin")
   ```

3. **Restart Neovim and open a `.sloth` file**

## 🐛 Troubleshooting

### Syntax Highlighting Not Working
- Ensure the file has `.sloth` extension
- Run `:set filetype=sloth` manually if needed
- Check if the plugin files are in the correct location

### Key Mappings Not Working
- Verify `sloth-runner` is in your PATH
- Check if keys are conflicting with other plugins
- Use `:map <leader>sr` to verify mapping exists

### Code Completion Not Showing
- Ensure completion is enabled: `:set completeopt=menu,menuone,noselect`
- Try triggering manually with `<C-x><C-o>`
- Check if omnifunc is set: `:set omnifunc?`

### Treesitter Conflicts (FIXED)
If you encounter treesitter errors like "attempt to call method 'parent' (a nil value)":
- **This is already fixed** in the current configuration
- The plugin automatically disables treesitter for `.sloth` files
- Uses traditional syntax highlighting to avoid conflicts
- If issues persist, restart Neovim

### Formatting Issues
- **Auto-formatting is disabled by default** to prevent conflicts
- Use manual formatting: `<leader>sf` or `:SlothFormat`
- For stylua formatting, ensure stylua is installed and configured

## 📖 Examples

### Basic Workflow File

```lua
-- deployment.sloth
local deploy_task = task("deploy_app")
    :description("Deploy application to production")
    :command(function(params, deps)
        local result = exec.run("kubectl apply -f deployment.yaml")
        if not result.success then
            log.error("Deployment failed: " .. result.stderr)
            return false
        end
        
        log.info("🚀 Deployment successful!")
        return true
    end)
    :timeout(300)
    :retries(3)
    :build()

workflow.define("production_deployment", {
    description = "Production deployment workflow",
    version = "1.0.0",
    tasks = { deploy_task }
})
```

With the plugin installed, this file will have:
- **Syntax highlighting** for keywords, functions, and strings
- **Code completion** when typing method names
- **Quick execution** with `<leader>sr`
- **Task listing** with `<leader>sl`

### Working with Environment Variables

```lua
-- config.sloth  
local config_task = task("setup_config")
    :description("Setup application configuration")
    :command(function()
        local env = "${NODE_ENV}" -- Highlighted as environment variable
        local config_path = "/app/config/${env}.json" -- Path highlighting
        
        if not fs.exists(config_path) then
            log.error("Config file not found: " .. config_path)
            return false
        end
        
        return true
    end)
    :build()
```

## 🔗 Integration with Other Tools

### With Telescope (File Finder)

```lua
-- Add to telescope config
require('telescope').setup({
  defaults = {
    file_ignore_patterns = { "%.git/", "node_modules/" },
    -- .sloth files will be found and properly highlighted
  }
})
```

### With LSP (Language Server Protocol)

While there's no dedicated LSP for Sloth DSL yet, you can use Lua LSP for basic support:

```lua
require('lspconfig').lua_ls.setup({
  settings = {
    Lua = {
      workspace = {
        checkThirdParty = false,
        library = {
          -- Add sloth-runner globals if needed
        }
      }
    }
  }
})
```

### With Git Integration

The plugin works well with git plugins like vim-fugitive:

```lua
-- .sloth files will be properly handled in git diffs
-- Syntax highlighting works in :Gdiff views
```

## 🚀 Next Steps

- **Learn the DSL:** Check out [Core Concepts](../core-concepts.md)
- **Try Examples:** See [Examples Guide](../EXAMPLES.md)
- **Advanced Features:** Explore [Advanced Features](../advanced-features.md)
- **API Reference:** Read [Lua API Documentation](../LUA_API.md)

---

The Neovim plugin makes writing Sloth workflows a breeze with full IDE support. Start creating powerful automation workflows with confidence! 🦥✨