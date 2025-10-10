# 🔌 插件开发

**为 Sloth Runner 平台构建扩展**

Sloth Runner 提供了强大的插件系统，允许开发者使用自定义功能扩展平台。本指南涵盖了开发自己的插件所需了解的一切。

## 🏗️ 插件架构

### 插件类型

Sloth Runner 支持多种类型的插件：

1. **🌙 Lua 模块** - 使用新功能和能力扩展 DSL
2. **⚡ 命令处理器** - 添加新的 CLI 命令和操作
3. **🎨 UI 扩展** - 增强 Web 仪表板和界面
4. **🔗 集成** - 连接外部工具和服务
5. **🦥 编辑器插件** - IDE/编辑器扩展（如我们的 Neovim 插件）

### 核心组件

```
sloth-runner/
├── plugins/
│   ├── lua-modules/       # Lua DSL 扩展
│   ├── commands/          # CLI 命令插件
│   ├── ui/               # Web UI 扩展
│   ├── integrations/     # 第三方集成
│   └── editors/          # 编辑器/IDE 插件
└── internal/
    └── plugin/           # 插件系统核心
```

## 🌙 开发 Lua 模块插件

### 基本结构

创建一个扩展 DSL 的新 Lua 模块：

```lua
-- plugins/lua-modules/my-module/init.lua
local M = {}

-- 模块元数据
M._NAME = "my-module"
M._VERSION = "1.0.0"
M._DESCRIPTION = "Sloth Runner 的自定义功能"

-- 公共 API
function M.hello(name)
    return string.format("你好，%s 来自我的自定义模块！", name or "世界")
end

function M.custom_task(config)
    return {
        execute = function(params)
            log.info("🔌 执行自定义任务: " .. config.name)
            -- 自定义任务逻辑
            return true
        end,
        validate = function()
            return config.name ~= nil
        end
    }
end

-- 注册模块函数
function M.register()
    -- 使函数在 DSL 中可用
    _G.my_module = M
    
    -- 注册自定义任务类型
    task.register_type("custom", M.custom_task)
end

return M
```

### 在工作流中使用自定义模块

```lua
-- workflow.sloth
local my_task = task("test_custom")
    :type("custom", { name = "test" })
    :description("测试自定义插件")
    :build()

-- 直接使用模块
local greeting = my_module.hello("开发者")
log.info(greeting)

workflow
    .define("plugin_test")
    :description("测试自定义插件")
    :version("1.0.0")
    :tasks({my_task})
```

### 插件注册

创建插件清单：

```yaml
# plugins/lua-modules/my-module/plugin.yaml
name: my-module
version: 1.0.0
description: Sloth Runner 的自定义功能
type: lua-module
author: 您的姓名
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
      description: "外部服务的 API 密钥"
```

## ⚡ 命令插件开发

### CLI 命令结构

```go
// plugins/commands/my-command/main.go
package main

import (
    "github.com/spf13/cobra"
    "github.com/chalkan3-sloth/sloth-runner/pkg/plugin"
)

type MyCommandPlugin struct {
    config *MyConfig
}

type MyConfig struct {
    Setting1 string `json:"setting1"`
    Setting2 int    `json:"setting2"`
}

func (p *MyCommandPlugin) Name() string {
    return "my-command"
}

func (p *MyCommandPlugin) Command() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "my-command",
        Short: "自定义命令功能",
        Long:  "自定义命令的扩展描述",
        RunE:  p.execute,
    }
    
    cmd.Flags().StringVar(&p.config.Setting1, "setting1", "", "自定义设置")
    cmd.Flags().IntVar(&p.config.Setting2, "setting2", 0, "另一个设置")
    
    return cmd
}

func (p *MyCommandPlugin) execute(cmd *cobra.Command, args []string) error {
    log.Info("🔌 执行自定义命令，设置:", 
        "setting1", p.config.Setting1,
        "setting2", p.config.Setting2)
    
    // 自定义命令逻辑
    return nil
}

func main() {
    plugin := &MyCommandPlugin{
        config: &MyConfig{},
    }
    
    plugin.Register()
}
```

## 🛠️ 插件开发工具

### 插件生成器

使用生成器快速创建新插件：

```bash
# 生成新的 Lua 模块插件
sloth-runner plugin generate --type=lua-module --name=my-module

# 生成 CLI 命令插件
sloth-runner plugin generate --type=command --name=my-command

# 生成 UI 扩展
sloth-runner plugin generate --type=ui --name=my-dashboard
```

### 开发环境

```bash
# 启动开发服务器，支持插件热重载
sloth-runner dev --plugins-dir=./plugins

# 本地测试插件
sloth-runner plugin test ./plugins/my-plugin

# 构建插件用于分发
sloth-runner plugin build ./plugins/my-plugin --output=dist/
```

### 插件测试

```go
// plugins/my-plugin/plugin_test.go
package main

import (
    "testing"
    "github.com/chalkan3-sloth/sloth-runner/pkg/plugin/testing"
)

func TestMyPlugin(t *testing.T) {
    // 创建测试环境
    env := plugintest.NewEnvironment(t)
    
    // 加载插件
    plugin, err := env.LoadPlugin("./")
    if err != nil {
        t.Fatal(err)
    }
    
    // 测试插件功能
    result, err := plugin.Execute(map[string]interface{}{
        "test_param": "value",
    })
    
    if err != nil {
        t.Fatal(err)
    }
    
    // 验证结果
    if result.Status != "success" {
        t.Errorf("期望成功，得到 %s", result.Status)
    }
}
```

## 📦 插件分发

### 插件注册表

将您的插件发布到 Sloth Runner 插件注册表：

```bash
# 登录注册表
sloth-runner registry login

# 发布插件
sloth-runner plugin publish ./my-plugin

# 安装已发布的插件
sloth-runner plugin install my-username/my-plugin
```

### 插件市场

浏览和发现插件：

```bash
# 搜索插件
sloth-runner plugin search "kubernetes"

# 获取插件信息
sloth-runner plugin info kubernetes-operator

# 从市场安装
sloth-runner plugin install --marketplace kubernetes-operator
```

## 🔒 安全性和最佳实践

### 安全指南

1. **🛡️ 最小权限原则** - 只请求必要的权限
2. **🔐 输入验证** - 始终验证用户输入和配置
3. **🚫 避免全局状态** - 保持插件状态隔离
4. **📝 错误处理** - 提供清晰的错误消息和日志记录
5. **🧪 测试** - 为所有功能编写全面的测试

### 代码质量

```go
// 好的：清晰的错误处理
func (p *MyPlugin) Execute(params map[string]interface{}) (*Result, error) {
    value, ok := params["required_param"].(string)
    if !ok {
        return nil, fmt.Errorf("required_param 必须是字符串")
    }
    
    if value == "" {
        return nil, fmt.Errorf("required_param 不能为空")
    }
    
    // 使用验证的输入进行处理
    result := p.process(value)
    return result, nil
}
```

### 文档标准

每个插件都应包括：

- **📋 README.md** - 安装和使用说明
- **📚 API 文档** - 函数/方法文档
- **📖 示例** - 工作代码示例
- **🧪 测试** - 单元测试和集成测试
- **📄 许可证** - 清晰的许可信息

## 📚 示例和模板

### 完整插件示例

查看这些示例插件：

- **[Kubernetes Operator Plugin](https://github.com/sloth-runner/plugin-kubernetes)** - 管理 K8s 资源
- **[Slack Integration Plugin](https://github.com/sloth-runner/plugin-slack)** - 发送通知
- **[Monitoring Dashboard Plugin](https://github.com/sloth-runner/plugin-monitoring)** - 自定义指标 UI

### 插件模板

使用官方模板快速开始：

```bash
# 使用模板
sloth-runner plugin init --template=lua-module my-plugin
sloth-runner plugin init --template=go-command my-command
sloth-runner plugin init --template=react-ui my-dashboard
```

## 💬 社区和支持

### 获取帮助

- **📖 [插件 API 文档](https://docs.sloth-runner.io/plugin-api)**
- **💬 [Discord 社区](https://discord.gg/sloth-runner)** - #plugin-development
- **🐛 [GitHub Issues](https://github.com/chalkan3-sloth/sloth-runner/issues)** - 错误报告和功能请求
- **📧 [插件邮件列表](mailto:plugins@sloth-runner.io)** - 开发讨论

### 贡献

我们欢迎插件贡献！请参阅我们的[贡献指南](contributing.md)了解以下详情：

- 插件提交流程
- 代码审查指南
- 文档要求
- 测试标准

---

今天就开始为 Sloth Runner 构建出色的插件！平台的可扩展架构使添加您需要的确切功能变得简单。🔌✨