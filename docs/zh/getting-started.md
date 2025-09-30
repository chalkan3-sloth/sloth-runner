# 快速入门

欢迎使用 Sloth-Runner！本指南将帮助您快速开始使用该工具。

> **📝 重要说明：** 从当前版本开始，Sloth Runner 工作流文件使用 `.sloth` 扩展名而不是 `.lua`。Lua 语法保持不变 - 只是文件扩展名更改为更好地识别 Sloth Runner DSL 文件。

## 安装

要在您的系统上安装 `sloth-runner`，您可以使用提供的 `install.sh` 脚本。此脚本会自动检测您的操作系统和架构，从 GitHub 下载最新版本，并将 `sloth-runner` 可执行文件放置在 `/usr/local/bin` 中。

```bash
bash <(curl -sL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh)
```

**注意：** `install.sh` 脚本需要 `sudo` 权限才能将可执行文件移动到 `/usr/local/bin`。

## 基本用法

### 堆栈管理

```bash
# 创建新堆栈
sloth-runner stack new my-app --description "应用程序部署堆栈"

# 在堆栈上运行工作流
sloth-runner run my-app -f examples/basic_pipeline.sloth

# 列出所有堆栈
sloth-runner stack list

# 查看堆栈详情
sloth-runner stack show my-app
```

### 直接工作流执行

要直接运行工作流文件：

```bash
sloth-runner run -f examples/basic_pipeline.sloth
```

要列出文件中的任务：

```bash
sloth-runner list -f examples/basic_pipeline.sloth
```

## 下一步

现在您已经安装并运行了 Sloth-Runner，请探索[核心概念](./core-concepts.md)以了解如何定义任务，或者直接深入了解新的[内置模块](../index.md#内置模块)以使用 Git、Pulumi 和 Salt 进行高级自动化。

---
[English](../en/getting-started.md) | [Português](../pt/getting-started.md) | [中文](./getting-started.md)