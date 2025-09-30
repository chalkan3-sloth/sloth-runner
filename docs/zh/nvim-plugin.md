# 🦥 Neovim 插件

**为 Neovim/LunarVim 提供完整的 Sloth Runner DSL 支持**

Sloth Runner Neovim 插件为处理 `.sloth` 工作流文件提供全面的 IDE 功能，包括语法高亮、代码补全和集成任务执行。

## ✨ 功能特性

### 🎨 丰富的语法高亮
- **自定义颜色** 用于 DSL 关键字、方法和模块
- **字符串插值** 高亮显示 `${variable}` 语法
- **文件路径检测** 用于脚本和配置文件
- **环境变量** 高亮显示
- **注释支持** 带有拼写检查

### 📁 智能文件检测
- 自动检测 `.sloth` 文件并应用适当的高亮
- 向后兼容 `.lua` 扩展名
- 文件浏览器中的自定义文件图标 (🦥)

### ⚡ 代码补全
- **智能补全** DSL 方法：`command`、`description`、`timeout` 等
- **模块补全** 内置模块：`exec`、`fs`、`net`、`aws` 等
- **函数补全** 常用模式：`task()`、`workflow.define()`

### 🔧 集成执行器
- **运行工作流** 直接在 Neovim 中使用 `<leader>sr`
- **列出任务** 当前文件中的任务使用 `<leader>sl`
- **试运行支持** 用于测试工作流

### 📋 代码片段和模板
- **快速任务创建** 使用 `_task` 缩写
- **工作流模板** 使用 `_workflow` 缩写
- **函数模板** 使用 `_cmd` 缩写
- **自动生成模板** 用于新的 `.sloth` 文件

### 🎯 文本对象和导航
- **选择任务块** 使用 `vit` (visual in task)
- **选择工作流块** 使用 `viw` (visual in workflow)
- **智能折叠** 可折叠的代码段
- **智能缩进** 用于 DSL 链式调用

## 🚀 快速设置

### LunarVim 用户

添加到你的 `~/.config/lvim/config.lua`：

```lua
-- 禁用自动格式化（推荐）
lvim.format_on_save.enabled = false

-- 配置 sloth 文件图标
require('nvim-web-devicons').setup {
  override_by_extension = {
    ["sloth"] = {
      icon = "🦥",
      color = "#8B4513",
      name = "SlothDSL"
    }
  }
}

-- sloth runner 按键映射
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

-- 手动格式化命令
lvim.keys.normal_mode["<leader>sf"] = ":SlothFormat<CR>"
```

### 标准 Neovim

使用 [lazy.nvim](https://github.com/folke/lazy.nvim)：

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

## 📝 按键映射

| 按键 | 动作 | 描述 |
|------|------|------|
| `<leader>sr` | 运行文件 | 执行当前 `.sloth` 工作流 |
| `<leader>sl` | 列出任务 | 显示当前文件中的所有任务 |
| `<leader>st` | 试运行 | 测试工作流而不执行 |
| `<leader>sd` | 调试 | 运行并输出调试信息 |
| `<leader>sf` | 格式化 | 格式化当前文件（手动） |

## 🎨 代码片段

### 快速任务创建
输入 `_task` 并按 Tab：

```lua
local task_name = task("")
    :description("")
    :command(function(params, deps)
        -- TODO: 实现
        return true
    end)
    :build()
```

### 快速工作流创建
输入 `_workflow` 并按 Tab：

```lua
workflow.define("", {
    description = "",
    version = "1.0.0",
    tasks = {
        -- 任务在这里
    }
})
```

### 快速命令函数
输入 `_cmd` 并按 Tab：

```lua
:command(function(params, deps)
    -- TODO: 实现
    return true
end)
```

## 🛠️ 手动安装

1. **克隆或复制插件文件：**
   ```bash
   cp -r /path/to/sloth-runner/nvim-plugin ~/.config/nvim/
   ```

2. **添加到你的 Neovim 配置：**
   ```lua
   -- 添加到 init.lua 或 init.vim
   vim.opt.runtimepath:append("~/.config/nvim/nvim-plugin")
   ```

3. **重启 Neovim 并打开 `.sloth` 文件**

## 🐛 故障排除

### 语法高亮不工作
- 确保文件有 `.sloth` 扩展名
- 如果需要，手动运行 `:set filetype=sloth`
- 检查插件文件是否在正确位置

### 按键映射不工作
- 验证 `sloth-runner` 在你的 PATH 中
- 检查按键是否与其他插件冲突
- 使用 `:map <leader>sr` 验证映射是否存在

### 代码补全不显示
- 确保补全已启用：`:set completeopt=menu,menuone,noselect`
- 尝试手动触发 `<C-x><C-o>`
- 检查 omnifunc 是否已设置：`:set omnifunc?`

## 📖 示例

### 基本工作流文件

```lua
-- deployment.sloth
local deploy_task = task("deploy_app")
    :description("部署应用到生产环境")
    :command(function(params, deps)
        local result = exec.run("kubectl apply -f deployment.yaml")
        if not result.success then
            log.error("部署失败: " .. result.stderr)
            return false
        end
        
        log.info("🚀 部署成功!")
        return true
    end)
    :timeout(300)
    :retries(3)
    :build()

workflow.define("production_deployment", {
    description = "生产环境部署工作流",
    version = "1.0.0",
    tasks = { deploy_task }
})
```

安装插件后，此文件将具有：
- **语法高亮** 用于关键字、函数和字符串
- **代码补全** 在输入方法名时
- **快速执行** 使用 `<leader>sr`
- **任务列表** 使用 `<leader>sl`

## 🚀 下一步

- **学习 DSL：** 查看 [核心概念](../zh/core-concepts.md)
- **尝试示例：** 参见 [示例指南](../zh/advanced-examples.md)
- **高级功能：** 探索 [高级功能](../zh/advanced-features.md)
- **API 参考：** 阅读 [Lua API 文档](../zh/lua-api-overview.md)

---

Neovim 插件让编写 Sloth 工作流变得轻而易举，提供完整的 IDE 支持。开始自信地创建强大的自动化工作流！🦥✨