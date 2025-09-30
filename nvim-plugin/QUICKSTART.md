# 🦥 Sloth Runner DSL - Neovim Plugin Quick Start

## 🚀 Quick Installation

```bash
# Clone the repository
git clone https://github.com/chalkan3-sloth/sloth-runner.git
cd sloth-runner

# Install the plugin
./nvim-plugin/install.sh
```

## 🎯 Test the Plugin

1. **Restart Neovim** or reload config:
   ```vim
   :source ~/.config/nvim/init.lua
   ```

2. **Open the example file**:
   ```bash
   nvim ~/example.sloth
   ```

3. **Verify syntax highlighting** - You should see:
   - `task` and `workflow` in **bright blue**
   - `:command()`, `:description()` methods in **teal**
   - `exec`, `fs`, `state` modules in **purple** 
   - Strings and comments properly colored

## ⌨️ Test Key Mappings

With a `.sloth` file open:

| Key | Action | Test |
|-----|--------|------|
| `<leader>sr` | Run File | Should attempt to run sloth-runner |
| `<leader>sl` | List Tasks | Shows tasks in current file |
| `<leader>st` | Dry Run | Tests workflow syntax |
| `<leader>sd` | Debug | Verbose execution |

## 📋 Test Snippets

In insert mode, type and expand:

- **`_task`** → Full task template
- **`_workflow`** → Full workflow template
- **`_cmd`** → Command function template

## 🎯 Test Text Objects

In normal/visual mode:

- **`vit`** - Select entire task block
- **`viw`** - Select entire workflow block
- **`dit`** - Delete task block
- **`diw`** - Delete workflow block

## 🔄 Test Folding

1. Open a file with multiple tasks/workflows
2. Use `zc` to close folds, `zo` to open
3. Should see nice fold text like:
   ```
   📋 Task: build_application (15 lines) ⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯
   🔄 Workflow: ci_pipeline (42 lines) ⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯⋯
   ```

## 🎨 Customization

Add to your `init.lua`:

```lua
require("sloth-runner").setup({
  runner = {
    command = "/path/to/your/sloth-runner", -- Custom binary path
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

-- Custom highlight colors
vim.api.nvim_set_hl(0, 'slothDSLKeyword', { fg = '#569cd6', bold = true })
vim.api.nvim_set_hl(0, 'slothModule', { fg = '#c586c0', bold = true })
vim.api.nvim_set_hl(0, 'slothMethod', { fg = '#4ec9b0' })
```

## 🐛 Troubleshooting

### Syntax highlighting not working?
```vim
:set filetype?  " Should show 'sloth'
:set filetype=sloth  " Force set if needed
```

### File not detected as Sloth DSL?
- Use `.sloth` extension
- Or add DSL keywords like `task(` or `workflow.define`

### Runner commands not working?
- Check if sloth-runner is in PATH: `:!which sloth-runner`
- Update binary path in configuration

## 🎯 Example File Types

The plugin auto-detects these files as Sloth DSL:

```lua
# example.sloth ✅
task("build"):command(function() end):build()

# ci-pipeline.sloth ✅  
workflow.define("ci", { tasks = {} })

# deploy-task.sloth ✅
local deploy = task("deploy")
    :description("Deploy app")
    :build()
```

---

**🎉 You now have a complete IDE experience for Sloth Runner DSL!**

Features working:
- ✅ Syntax highlighting
- ✅ File detection  
- ✅ Code completion
- ✅ Snippets
- ✅ Text objects
- ✅ Folding
- ✅ Runner integration
- ✅ Key mappings