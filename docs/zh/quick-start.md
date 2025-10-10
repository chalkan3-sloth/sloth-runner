# ⚡ 快速教程

完整的中文文档，请访问：

## 🚀 快速开始

### 安装

```bash
# 下载
curl -sSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh | bash

# 或通过 Go
go install github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner@latest
```

### 第一个工作流

创建文件 `hello.sloth`:

```lua
local hello_task = task("hello")
    :description("我的第一个任务")
    :command(function(this, params)
        print("🦥 你好，来自 Sloth Runner!")
        return true, "成功完成"
    end)
    :build()

workflow
    .define("hello_world")
    :description("我的第一个工作流")
    :version("1.0.0")
    :tasks({hello_task})
```

运行:

```bash
sloth-runner run -f hello.sloth
```

## 📚 下一步

- [核心概念](./core-concepts.md)
- [高级示例](./advanced-examples.md)
- [高级功能](./advanced-features.md)

完整教程，请参阅：[主教程](../TUTORIAL.md)
