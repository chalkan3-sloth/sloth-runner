# Sloth Runner DSL - Modern Neovim Plugin

🦥 **Complete IDE experience for Sloth Runner DSL in Neovim**

> **📝 Important Note:** This plugin provides first-class support for `.sloth` files with modern Neovim features including LSP-style completion, Telescope integration, and health checks.

## ✨ Features

- 🦥 **Welcome Banner** - Friendly sloth emoji greets you when opening `.sloth` files
- 🎨 **Rich Syntax Highlighting** - Enhanced colors for DSL keywords, methods, and modules
- 📁 **Smart File Detection** - Auto-detects `.sloth` files and workflow patterns
- ⚡ **Intelligent Completion** - nvim-cmp integration with context-aware suggestions
- 🔭 **Telescope Integration** - Interactive task and workflow pickers
- 🔧 **Integrated Runner** - Execute workflows directly from Neovim with floating output
- 📋 **Quick Templates** - Snippets for rapid task and workflow creation
- 🎯 **Text Objects** - Navigate and manipulate DSL constructs easily
- 🔍 **Code Formatting** - Auto-format with stylua or built-in formatter
- 🏥 **Health Checks** - Comprehensive `:checkhealth` integration
- 🔌 **Modular Architecture** - Clean, maintainable Lua codebase

## 📋 Requirements

- **Neovim** >= 0.9.0
- **sloth-runner** executable in PATH
- **Optional**: [nvim-cmp](https://github.com/hrsh7th/nvim-cmp) for completion
- **Optional**: [telescope.nvim](https://github.com/nvim-telescope/telescope.nvim) for pickers
- **Optional**: [which-key.nvim](https://github.com/folke/which-key.nvim) for keymap hints

## 🚀 Installation

### Using [lazy.nvim](https://github.com/folke/lazy.nvim) (Recommended)

```lua
{
  dir = "/path/to/sloth-runner/nvim-plugin", -- Or use git URL when published
  name = "sloth-runner",
  ft = "sloth",
  dependencies = {
    "hrsh7th/nvim-cmp",              -- Optional: for completion
    "nvim-telescope/telescope.nvim", -- Optional: for pickers
    "folke/which-key.nvim",          -- Optional: for keymap hints
  },
  config = function()
    require("sloth-runner").setup({
      -- Your configuration here (see Configuration section)
    })
  end,
}
```

### Using [packer.nvim](https://github.com/wbthomason/packer.nvim)

```lua
use {
  "/path/to/sloth-runner/nvim-plugin",
  ft = "sloth",
  requires = {
    "hrsh7th/nvim-cmp",
    "nvim-telescope/telescope.nvim",
  },
  config = function()
    require("sloth-runner").setup()
  end,
}
```

### Manual Installation

```bash
# Copy to your Neovim configuration
cp -r nvim-plugin/* ~/.config/nvim/

# Add to your init.lua
lua require("sloth-runner").setup()
```

## ⚙️ Configuration

### Default Configuration

```lua
require("sloth-runner").setup({
  -- Runner configuration
  runner = {
    cmd = "sloth-runner",           -- Runner executable
    default_args = {},              -- Default CLI arguments
    use_float = true,               -- Use floating window for output
    auto_close_on_success = false,  -- Auto-close window on success
    notify = true,                  -- Show notifications
  },

  -- Formatter configuration
  formatter = {
    format_on_save = false,         -- Auto-format on save
    cmd = "stylua",                 -- External formatter (stylua, etc.)
    args = {                        -- Formatter arguments
      "--indent-type", "Spaces",
      "--indent-width", "2"
    },
    use_builtin = true,             -- Fallback to built-in formatter
  },

  -- Completion configuration
  completion = {
    enabled = true,                 -- Enable nvim-cmp integration
    priority = 100,                 -- Completion source priority
    show_docs = true,               -- Show documentation in completion
  },

  -- Keymap configuration
  keymaps = {
    enabled = true,                 -- Enable default keymaps
    prefix = "<leader>s",           -- Keymap prefix
    run = "r",                      -- Run workflow
    list = "l",                     -- List tasks
    test = "t",                     -- Dry run
    validate = "v",                 -- Validate file
    format = "f",                   -- Format file
    task_textobj = "it",            -- Task text object
    workflow_textobj = "iw",        -- Workflow text object
  },

  -- Telescope integration
  telescope = {
    enabled = true,                 -- Enable telescope integration
    theme = "dropdown",             -- Picker theme
    layout_config = {
      width = 0.8,
      height = 0.6,
    },
  },

  -- UI configuration
  ui = {
    icons = {
      task = "📋",
      workflow = "🔄",
      running = "⚡",
      success = "✓",
      error = "✗",
      warning = "⚠",
    },
    float = {
      border = "rounded",           -- Border style
      title_pos = "center",         -- Title position
      width = 0.8,
      height = 0.8,
    },
    -- Welcome banner
    show_welcome = true,            -- Show 🦥 when opening .sloth files
    welcome_style = "notification", -- "notification", "banner", "large", "float"
  },

  -- Debug configuration
  debug = {
    enabled = false,                -- Enable debug logging
    log_file = vim.fn.stdpath("cache") .. "/sloth-runner.log",
  },
})
```

### Minimal Configuration

```lua
-- Plugin works out of the box with defaults
require("sloth-runner").setup()
```

### Custom Configuration Example

```lua
require("sloth-runner").setup({
  runner = {
    use_float = false,              -- Use terminal split instead
    notify = false,                 -- Disable notifications
  },
  formatter = {
    format_on_save = true,          -- Auto-format on save
  },
  keymaps = {
    prefix = "<leader>r",           -- Use different prefix
  },
})
```

## ⌨️ Default Keymaps

All keymaps use `<leader>s` prefix by default (configurable):

| Keymap | Command | Description |
|--------|---------|-------------|
| `<leader>sr` | `:SlothRun` | Run workflow in current file |
| `<leader>sl` | `:SlothList` | List all tasks and workflows |
| `<leader>st` | `:SlothTest` | Dry run (test without execution) |
| `<leader>sv` | `:SlothValidate` | Validate file syntax |
| `<leader>sf` | `:SlothFormat` | Format current file |

### Text Objects

| Text Object | Description |
|-------------|-------------|
| `it` | Inner task block |
| `iw` | Inner workflow block |

**Examples:**
- `vit` - Visually select task block
- `dit` - Delete task block
- `yiw` - Yank workflow block
- `ciw` - Change workflow block

## 🔧 Commands

| Command | Description |
|---------|-------------|
| `:SlothRun [task]` | Run workflow or specific task |
| `:SlothList` | List all tasks/workflows |
| `:SlothTest [task]` | Dry run workflow or task |
| `:SlothValidate` | Validate current file |
| `:SlothFormat` | Format current file |
| `:SlothInfo` | Show plugin information |
| `:SlothWelcome` | Show welcome message with 🦥 |
| `:SlothAnimate` | Show sloth animation (easter egg) |
| `:SlothTasks` | Open Telescope task picker |
| `:SlothWorkflows` | Open Telescope workflow picker |

## 🔭 Telescope Integration

When telescope.nvim is installed, you get interactive pickers:

### Task Picker

```vim
:SlothTasks
" or
:Telescope sloth tasks
```

**Keymaps in picker:**
- `<CR>` - Run selected task
- `<C-g>` - Go to task definition

### Workflow Picker

```vim
:SlothWorkflows
" or
:Telescope sloth workflows
```

**Keymaps in picker:**
- `<CR>` - Run selected workflow
- `<C-g>` - Go to workflow definition

### Combined Picker

```vim
:Telescope sloth
```

Shows both tasks and workflows in one picker.

## ⚡ Completion

When nvim-cmp is installed, the plugin provides intelligent completion:

- **DSL Keywords**: `task`, `workflow`, `define`
- **Method Chaining**: `:command`, `:description`, `:timeout`, `:build`
- **Modules**: `exec`, `fs`, `net`, `aws`, `docker`, `kubernetes`
- **Context-Aware**: Suggests appropriate completions based on cursor position

The completion source is automatically registered with nvim-cmp when both are installed.

## 📋 Snippets & Templates

### Task Template

Type `_task` in insert mode and expand:

```lua
local task_name = task("task-name")
  :description("Task description")
  :command(function(params, deps)
    -- TODO: implement
    return true
  end)
  :build()
```

### Workflow Template

Type `_workflow` in insert mode and expand:

```lua
workflow.define("workflow-name", {
  description = "Workflow description",
  version = "1.0.0",
  tasks = {
    -- tasks here
  }
})
```

## 🎨 Syntax Highlighting

The plugin provides rich syntax highlighting for:

- **Keywords**: `task`, `workflow`, `local`, `function`
- **DSL Methods**: `:command()`, `:description()`, `:build()` (golden highlight)
- **Modules**: `exec`, `fs`, `aws`, `kubernetes` (purple)
- **Strings**: With template interpolation support
- **Environment Variables**: `${VAR}` and `$VAR` patterns
- **Comments**: With TODO/FIXME highlighting

### Color Scheme

```
DSL Keywords      → Bright Blue (#569cd6)
Modules           → Purple (#c586c0)
Chain Methods (:) → Bright Golden (#f9e79f)
Functions         → Yellow (#dcdcaa)
Env Variables     → Red (#ff6b6b)
File Paths        → Cyan (#98d8c8)
```

## 🔄 Folding

Intelligent folding for:

- Task definitions (from `task(` to `:build()`)
- Workflow definitions (from `workflow.define(` to closing brace)
- Function blocks

**Fold display:**
```
📋 Task: deploy_application (15 lines) ⋯⋯⋯⋯⋯⋯⋯
🔄 Workflow: ci_pipeline (42 lines) ⋯⋯⋯⋯⋯⋯⋯⋯
⚡ Function: deploy_to_env (8 lines) ⋯⋯⋯⋯⋯⋯⋯
```

## 🏥 Health Check

Check plugin status and dependencies:

```vim
:checkhealth sloth-runner
```

The health check verifies:
- ✅ Plugin initialization
- ✅ sloth-runner executable availability
- ✅ Optional dependencies (nvim-cmp, telescope)
- ✅ Configuration validity
- ✅ Formatter availability

## 🔌 API

### Lua API

```lua
local sloth = require("sloth-runner")

-- Initialize plugin
sloth.setup({ ... })

-- Run workflow
sloth.run({ file = "path/to/file.sloth", task = "build" })

-- List tasks
sloth.list({ file = "path/to/file.sloth" })

-- Validate file
local success, err = sloth.validate()

-- Format file
sloth.format()

-- Get configuration
local config = sloth.get_config()
```

## 📚 Example Workflow

Create a file `deploy.sloth`:

```lua
-- This file is auto-detected as Sloth DSL

local build_task = task("build")
  :description("Build the application")
  :command(function(params, deps)
    local result = exec.run("go build -o app ./cmd/main.go")
    return result.success, result.stdout, { artifact = "app" }
  end)
  :timeout("5m")
  :build()

local test_task = task("test")
  :description("Run tests")
  :depends_on({"build"})
  :command(function(params, deps)
    local result = exec.run("go test ./...")
    return result.success
  end)
  :build()

local deploy_task = task("deploy")
  :description("Deploy to production")
  :depends_on({"build", "test"})
  :command(function(params, deps)
    local app = deps.build.artifact
    log.info("Deploying " .. app)

    local result = exec.run("kubectl apply -f deployment.yaml")
    return result.success
  end)
  :build()

workflow.define("production_deploy", {
  description = "Complete production deployment pipeline",
  version = "1.0.0",

  tasks = { build_task, test_task, deploy_task },

  on_success = function(results)
    log.info("🚀 Deployment completed successfully!")
    notification.send({
      title = "Deployment Success",
      message = "Production deployment completed"
    })
  end,

  on_failure = function(error, context)
    log.error("❌ Deployment failed: " .. error.message)
    notification.send({
      title = "Deployment Failed",
      message = error.message,
      urgency = "critical"
    })
  end
})
```

**Now you can:**
- `<leader>sr` - Run the entire workflow
- `<leader>sl` - List all tasks
- `<leader>st` - Dry run to test
- `:SlothTasks` - Pick and run individual tasks
- `vit` - Select a task block
- `viw` - Select the workflow

## 🦥 Welcome Banner

When you open a `.sloth` file for the first time, you'll see a friendly sloth emoji! You can customize the welcome style:

**Notification Style** (default):
```lua
ui = { welcome_style = "notification" }
-- Shows: 🦥 Sloth Runner DSL
```

**Banner Style**:
```lua
ui = { welcome_style = "banner" }
-- Shows:
--     🦥 Sloth Runner DSL
--   ⚡ Ready to automate workflows!
```

**Large Banner Style**:
```lua
ui = { welcome_style = "large" }
-- Shows:
--  ╔══════════════════════════════════╗
--  ║     🦥  SLOTH RUNNER DSL  🦥     ║
--  ╚══════════════════════════════════╝
--   ⚡ Workflow automation made easy
--   💡 Tip: Use <leader>sr to run workflow
```

**Float Window Style**:
```lua
ui = { welcome_style = "float" }
-- Shows banner in a floating window (auto-closes after 3s)
```

**Disable Welcome**:
```lua
ui = { show_welcome = false }
```

**Easter Egg**: Try `:SlothAnimate` for a fun sloth animation! 🦥💤⚡✨

## 📖 Documentation

Full documentation available in:
- `:help sloth-runner` - Complete help documentation
- `:checkhealth sloth-runner` - Health check and diagnostics
- `:SlothInfo` - Quick reference
- `:SlothWelcome` - See the welcome banner again

## 🤝 Contributing

Contributions are welcome! The plugin uses a modern modular architecture:

```
nvim-plugin/
├── lua/sloth-runner/
│   ├── init.lua          # Main entry point
│   ├── config.lua        # Configuration management
│   ├── commands.lua      # Command definitions
│   ├── keymaps.lua       # Keymap setup
│   ├── runner.lua        # Workflow execution
│   ├── formatter.lua     # Code formatting
│   ├── completion.lua    # nvim-cmp integration
│   ├── telescope.lua     # Telescope integration
│   ├── health.lua        # Health check
│   ├── welcome.lua       # Welcome banner (🦥)
│   └── utils.lua         # Utility functions
├── plugin/
│   └── sloth-runner.vim  # Plugin initialization
├── ftdetect/
│   └── sloth.vim         # Filetype detection
├── ftplugin/
│   └── sloth.vim         # Filetype plugin
├── syntax/
│   └── sloth.vim         # Syntax highlighting
└── doc/
    └── sloth-runner.txt  # Help documentation
```

## 📝 License

MIT License - See LICENSE file for details

---

**Happy coding with Sloth Runner DSL! 🦥⚡**
