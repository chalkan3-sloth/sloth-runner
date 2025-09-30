# Sloth Runner DSL - Neovim Plugin

🦥 **Syntax highlighting and IDE features for Sloth Runner DSL in Neovim**

## ✨ Features

- **🎨 Rich Syntax Highlighting** - Custom colors for DSL keywords, methods, and modules
- **📁 Smart File Detection** - Auto-detects Sloth DSL files 
- **⚡ Code Completion** - Intelligent completion for DSL methods and modules
- **🔧 Integrated Runner** - Run workflows directly from Neovim
- **📋 Code Snippets** - Quick templates for tasks and workflows
- **🔄 Folding Support** - Collapsible task and workflow blocks
- **🎯 Text Objects** - Navigate and select DSL constructs easily

## 🚀 Installation

### Using [lazy.nvim](https://github.com/folke/lazy.nvim)

```lua
{
  dir = "/path/to/sloth-runner/nvim-plugin", -- Local plugin path
  name = "sloth-runner",
  ft = { "sloth", "lua" },
  config = function()
    require("sloth-runner").setup({
      runner = {
        command = "sloth-runner", -- Path to your sloth-runner binary
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
      }
    })
  end
}
```

### Manual Installation

1. **Copy plugin files to your Neovim config:**

```bash
# Copy to your Neovim configuration directory
cp -r nvim-plugin/* ~/.config/nvim/
```

2. **Add to your init.lua:**

```lua
require("sloth-runner").setup()
```

## 📁 File Detection

The plugin automatically detects Sloth DSL files based on:

- **File extensions**: `*.sloth.lua`
- **File patterns**: `*task*.lua`, `*workflow*.lua`
- **Content detection**: Files containing `task(` or `workflow.define`

## 🎨 Syntax Highlighting

### Highlighted Elements

- **Keywords**: `task`, `workflow`, `local`, `function`, etc.
- **DSL Methods**: `:command()`, `:description()`, `:build()`, etc.
- **Modules**: `exec`, `fs`, `state`, `aws`, `kubernetes`, etc.
- **Strings**: With special handling for templates and paths
- **Comments**: Including TODO/FIXME highlighting
- **Environment Variables**: `${VAR}` and `$VAR` patterns

### Color Scheme

The plugin provides optimized colors for modern terminals:

```lua
-- Modern terminal colors (256-color and GUI)
DSL Keywords  → Bright Blue (#569cd6)
Modules       → Purple (#c586c0) 
Methods       → Teal (#4ec9b0)
Functions     → Yellow (#dcdcaa)
Env Variables → Red (#ff6b6b)
File Paths    → Cyan (#98d8c8)
```

## ⌨️ Key Mappings

| Key | Action | Description |
|-----|--------|-------------|
| `<leader>sr` | Run File | Execute current workflow file |
| `<leader>sl` | List Tasks | Show tasks in current file |
| `<leader>st` | Dry Run | Test workflow without execution |
| `<leader>sd` | Debug | Run with verbose debugging |

## 🔧 Commands

| Command | Description |
|---------|-------------|
| `:SlothRun` | Run current file |
| `:SlothList` | List tasks in file |
| `:SlothDryRun` | Dry run current file |
| `:SlothDebug` | Debug current workflow |
| `:SlothTaskSnippet` | Insert task template |
| `:SlothWorkflowSnippet` | Insert workflow template |

## 📋 Code Snippets

### Task Template

Type `_task` and expand:

```lua
local task_name = task("task_name")
    :description("Task description")
    :command(function(params, deps)
        -- TODO: implement task logic
        return true
    end)
    :build()
```

### Workflow Template  

Type `_workflow` and expand:

```lua
workflow.define("workflow_name", {
    description = "Workflow description",
    version = "1.0.0",
    
    tasks = {
        -- Add tasks here
    },
    
    on_success = function(results)
        -- Success handler
    end,
    
    on_failure = function(error, context)
        -- Error handler  
    end
})
```

## 🎯 Text Objects

- **`vit`** - Select task block (visual in task)
- **`viw`** - Select workflow block (visual in workflow)
- **`dit`** - Delete task block
- **`diw`** - Delete workflow block

## 🔄 Folding

The plugin provides intelligent folding for:

- **Task definitions** - From `task(` to `:build()`
- **Workflow definitions** - From `workflow.define(` to closing brace
- **Function blocks** - Standard Lua functions

### Fold Display

```
📋 Task: deploy_application (15 lines) ⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯
🔄 Workflow: ci_pipeline (42 lines) ⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯
⚡ Function: deploy_to_env (8 lines) ⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯
```

## ⚙️ Configuration

### Full Configuration Example

```lua
require("sloth-runner").setup({
  -- Syntax highlighting
  highlights = {
    enable = true,
    use_treesitter = true, -- Use tree-sitter if available
  },
  
  -- Code completion
  completion = {
    enable = true,
    snippets = true,
  },
  
  -- Runner integration
  runner = {
    command = "sloth-runner", -- Path to binary
    auto_run_on_save = false, -- Auto-run on file save
    keymaps = {
      run_file = "<leader>sr",
      list_tasks = "<leader>sl", 
      dry_run = "<leader>st",
      debug = "<leader>sd",
    }
  },
  
  -- Code folding
  folding = {
    enable = true,
    auto_close = false, -- Auto-close folds
  }
})
```

### Disable Features

```lua
require("sloth-runner").setup({
  completion = { enable = false },
  folding = { enable = false },
  runner = { keymaps = false }
})
```

## 🔌 Integration

### LSP Integration

The plugin works alongside Lua LSP for enhanced features:

```lua
-- If using nvim-lspconfig
require('lspconfig').lua_ls.setup({
  filetypes = { 'lua', 'sloth' }, -- Add sloth filetype
  settings = {
    Lua = {
      workspace = {
        library = {
          -- Add sloth-runner modules to workspace
          "/path/to/sloth-runner/lua-modules"
        }
      }
    }
  }
})
```

### Completion Integration (nvim-cmp)

```lua
require('cmp').setup.filetype('sloth', {
  sources = {
    { name = 'sloth' },      -- DSL-specific completion
    { name = 'lua_ls' },     -- Lua language server
    { name = 'luasnip' },    -- Snippets
    { name = 'buffer' },     -- Buffer words
  }
})
```

## 🎨 Color Customization

Override highlight groups in your config:

```lua
vim.api.nvim_set_hl(0, 'SlothKeyword', { fg = '#your-color', bold = true })
vim.api.nvim_set_hl(0, 'SlothModule', { fg = '#your-color', italic = true })
vim.api.nvim_set_hl(0, 'SlothMethod', { fg = '#your-color' })
```

## 📚 Example Usage

Create a file `deploy.sloth.lua`:

```lua
-- This file will be auto-detected as Sloth DSL

local build_task = task("build")
    :description("Build the application")
    :command(function(params, deps)
        local result = exec.run("go build -o app ./cmd/main.go")
        return result.success, result.stdout, { artifact = "app" }
    end)
    :timeout("5m")
    :build()

local deploy_task = task("deploy")
    :description("Deploy to production")
    :depends_on({"build"})
    :command(function(params, deps)
        local app_artifact = deps.build.artifact
        log.info("Deploying " .. app_artifact)
        
        local result = exec.run("kubectl apply -f deployment.yaml")
        return result.success
    end)
    :build()

workflow.define("production_deploy", {
    description = "Complete production deployment pipeline",
    version = "1.0.0",
    
    tasks = { build_task, deploy_task },
    
    on_success = function(results)
        log.info("🚀 Deployment completed successfully!")
    end,
    
    on_failure = function(error, context)
        log.error("❌ Deployment failed: " .. error.message)
    end
})
```

Now you can:
- **`<leader>sr`** - Run the entire workflow
- **`<leader>sl`** - See all tasks with their IDs
- **`<leader>st`** - Test without actually deploying
- Use **`vit`** to select the build task
- Use **`viw`** to select the entire workflow

---

**Happy coding with Sloth Runner DSL! 🦥⚡**