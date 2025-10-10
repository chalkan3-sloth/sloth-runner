# 🤝 贡献 Sloth Runner

**感谢您对贡献 Sloth Runner 的兴趣！**

我们欢迎所有技能水平的开发者的贡献。无论您是修复错误、添加功能、改进文档还是创建插件，您的帮助都会让 Sloth Runner 变得更好。

## 🚀 快速开始

### 前置条件

- **Go 1.21+** 用于核心开发
- **Node.js 18+** 用于 UI 开发  
- **Lua 5.4+** 用于 DSL 开发
- **Git** 用于版本控制

### 开发环境设置

```bash
# 克隆仓库
git clone https://github.com/chalkan3-sloth/sloth-runner.git
cd sloth-runner

# 安装依赖
go mod download
npm install  # 用于 UI 组件

# 运行测试
make test

# 构建项目
make build
```

## 📋 贡献方式

### 🐛 错误报告

发现了错误？请帮助我们修复：

1. **搜索现有 issue** 以避免重复
2. **使用我们的错误报告模板**，包含：
   - Sloth Runner 版本
   - 操作系统
   - 重现步骤
   - 预期行为 vs 实际行为
   - 错误日志（如果有）

### 💡 功能请求

有改进想法？

1. **检查路线图** 查看计划的功能
2. **开启功能请求**，包含：
   - 功能的清晰描述
   - 用例和好处
   - 可能的实现方法

### 🔧 代码贡献

准备编码？以下是步骤：

1. **Fork 仓库**
2. **创建功能分支** (`git checkout -b feature/amazing-feature`)
3. **进行更改** 遵循我们的编码标准
4. **为新功能添加测试**
5. **如需要更新文档**
6. **使用清晰的消息提交**
7. **推送并创建 Pull Request**

### 📚 文档

帮助改进我们的文档：

- 修复拼写错误和不清楚的说明
- 添加示例和教程
- 将内容翻译成其他语言
- 更新 API 文档

### 🔌 插件开发

为社区创建插件：

- 遵循我们的[插件开发指南](plugin-development.md)
- 提交到插件注册表
- 保持与核心版本的兼容性

## 📐 开发指南

### 代码风格

#### Go 代码

遵循标准 Go 约定：

```go
// 好的：清晰的函数名和注释
func ProcessWorkflowTasks(ctx context.Context, workflow *Workflow) error {
    if workflow == nil {
        return fmt.Errorf("workflow 不能为 nil")
    }
    
    for _, task := range workflow.Tasks {
        if err := processTask(ctx, task); err != nil {
            return fmt.Errorf("处理任务 %s 失败: %w", task.ID, err)
        }
    }
    
    return nil
}
```

#### Lua DSL

保持 DSL 代码清洁可读：

```lua
-- 好的：清晰的任务定义，适当的链式调用
local deploy_task = task("deploy_application")
    :description("将应用部署到生产环境")
    :command(function(this, params)
        local result = exec.run("kubectl apply -f deployment.yaml")
        if not result.success then
            log.error("部署失败: " .. result.stderr)
            return false, "部署失败"
        end
        return true, "部署成功"
    end)
    :timeout("300s")
    :retries(3)
    :build()
```

### 测试标准

#### 单元测试

为所有新功能编写测试：

```go
func TestProcessWorkflowTasks(t *testing.T) {
    tests := []struct {
        name     string
        workflow *Workflow
        wantErr  bool
    }{
        {
            name:     "nil workflow 应该返回错误",
            workflow: nil,
            wantErr:  true,
        },
        {
            name: "有效 workflow 应该成功处理",
            workflow: &Workflow{
                Tasks: []*Task{{ID: "test-task"}},
            },
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ProcessWorkflowTasks(context.Background(), tt.workflow)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessWorkflowTasks() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 文档标准

- **保持简单** - 使用清晰、简洁的语言
- **包含示例** - 展示而不仅仅是告诉
- **随更改更新** - 保持文档与代码同步
- **测试示例** - 确保所有代码示例都能工作

## 🔄 Pull Request 流程

### 提交前

- [ ] **运行测试** - `make test`
- [ ] **运行代码检查** - `make lint`
- [ ] **更新文档** - 如果添加/更改功能
- [ ] **添加更新日志条目** - 在 `CHANGELOG.md` 中
- [ ] **检查兼容性** - 与现有功能

### PR 模板

使用我们的 pull request 模板：

```markdown
## 描述
更改的简要描述

## 更改类型
- [ ] 错误修复
- [ ] 新功能
- [ ] 破坏性更改
- [ ] 文档更新

## 测试
- [ ] 单元测试已添加/更新
- [ ] 集成测试通过
- [ ] 手动测试完成

## 检查清单
- [ ] 代码遵循样式指南
- [ ] 文档已更新
- [ ] 更新日志已更新
```

## 🏗️ 项目结构

理解代码库：

```
sloth-runner/
├── cmd/                    # CLI 命令
├── internal/              # 内部包
│   ├── core/             # 核心业务逻辑
│   ├── dsl/              # DSL 实现
│   ├── execution/        # 任务执行引擎
│   └── plugins/          # 插件系统
├── pkg/                   # 公共包
├── plugins/              # 内置插件
├── docs/                 # 文档
├── web/                  # Web UI 组件
└── examples/             # 示例工作流
```

## 🎯 贡献领域

### 高优先级

- **🐛 错误修复** - 总是欢迎
- **📈 性能改进** - 优化机会
- **🧪 测试覆盖率** - 增加测试覆盖率
- **📚 文档** - 保持文档全面

### 中等优先级

- **✨ 新功能** - 遵循路线图优先级
- **🔌 插件生态系统** - 更多插件和集成
- **🎨 UI 改进** - 更好的用户体验

## 🏆 认可

贡献者在以下方面得到认可：

- **CONTRIBUTORS.md** - 列出所有贡献者
- **发布说明** - 突出显示主要贡献
- **社区展示** - 特色贡献
- **贡献者徽章** - GitHub 个人资料认可

## 📞 获取帮助

### 开发问题

- **💬 Discord** - `#development` 频道
- **📧 邮件列表** - dev@sloth-runner.io
- **📖 Wiki** - 开发指南和常见问题

### 指导

初次接触开源？我们提供指导：

- **👥 导师匹配** - 与有经验的贡献者配对
- **📚 学习资源** - 策划的学习材料
- **🎯 引导贡献** - 适合初学者的 issue

## 📜 行为准则

我们致力于提供一个欢迎和包容的环境。请阅读我们的[行为准则](https://github.com/chalkan3-sloth/sloth-runner/blob/main/CODE_OF_CONDUCT.md)。

### 我们的标准

- **🤝 互相尊重** - 尊重对待每个人
- **💡 建设性** - 提供有用的反馈
- **🌍 包容性** - 欢迎多元化的观点
- **📚 耐心** - 帮助他人学习和成长

---

**准备贡献？**

从探索我们的[新手友好的 Issues](https://github.com/chalkan3-sloth/sloth-runner/labels/good%20first%20issue) 开始，或加入我们的 [Discord 社区](https://discord.gg/sloth-runner) 介绍自己！

感谢您帮助让 Sloth Runner 变得更好！🦥✨