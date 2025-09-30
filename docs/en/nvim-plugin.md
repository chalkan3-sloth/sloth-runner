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

## 📝 Key Mappings

| Key | Action | Description |
|-----|--------|-------------|
| `<leader>sr` | Run File | Execute current `.sloth` workflow |
| `<leader>sl` | List Tasks | Show all tasks in current file |
| `<leader>st` | Dry Run | Test workflow without execution |
| `<leader>sd` | Debug | Run with debug output |
| `<leader>sf` | Format | Format current file (manual) |

## 🛠️ Installation

### Manual Installation

1. **Copy plugin files to your Neovim config:**
   ```bash
   cp -r /path/to/sloth-runner/nvim-plugin ~/.config/nvim/sloth-runner
   ```

2. **Add to your Neovim configuration:**
   ```lua
   -- Add to init.lua
   vim.opt.runtimepath:append("~/.config/nvim/sloth-runner")
   ```

3. **Restart Neovim and open a `.sloth` file**

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

## 🔧 Configuration Options

The plugin automatically configures itself when you open `.sloth` files. Key features include:

- **Automatic filetype detection** for `.sloth` extensions
- **Syntax highlighting** with custom color scheme
- **Code completion** using omnifunc
- **Smart indentation** for DSL method chaining
- **Code folding** for task and workflow blocks
- **Key mappings** for common operations

## 🐛 Troubleshooting

### Syntax Highlighting Not Working
- Ensure the file has `.sloth` extension
- Run `:set filetype=sloth` manually if needed
- Check if plugin files are in correct location

### Key Mappings Not Working
- Verify `sloth-runner` is in your PATH
- Check for conflicts with other plugins
- Use `:map <leader>sr` to verify mapping exists

### Code Completion Not Showing
- Ensure completion is enabled: `:set completeopt=menu,menuone,noselect`
- Try triggering manually with `<C-x><C-o>`
- Check omnifunc setting: `:set omnifunc?`

## 📖 Example Workflow

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

---

The Neovim plugin makes writing Sloth workflows a breeze with full IDE support. Start creating powerful automation workflows with confidence! 🦥✨