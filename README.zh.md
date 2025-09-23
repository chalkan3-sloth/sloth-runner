[English](./README.md) | [Português](./README.pt.md) | [中文](./README.zh.md)

# 🦥 Sloth Runner 🚀

一个用 Go 编写、由 Lua 脚本驱动的灵活且可扩展的任务运行器应用程序。`sloth-runner` 允许您通过简单的 Lua 脚本定义复杂的工作流、管理任务依赖关系以及与外部系统集成。

[![Go CI](https://github.com/chalkan3/sloth-runner/actions/workflows/go.yml/badge.svg)](https://github.com/chalkan3/sloth-runner/actions/workflows/go.yml)

---

## ✨ 功能特性

*   **📜 Lua 脚本:** 使用强大而灵活的 Lua 脚本定义任务和工作流。
*   **🔗 依赖管理:** 指定任务依赖关系，确保复杂管道的有序执行。
*   **⚡ 异步任务执行:** 并发运行任务以提高性能。
*   **🪝 执行前后钩子:** 定义自定义 Lua 函数，在任务命令之前和之后运行。
*   **⚙️ 丰富的 Lua API:** 直接从 Lua 任务访问系统功能：
    *   **`exec` 模块:** 执行 shell 命令。
    *   **`fs` 模块:** 执行文件系统操作（读、写、追加、检查存在、创建目录、删除、递归删除、列出）。
    *   **`net` 模块:** 发出 HTTP 请求（GET、POST）和下载文件。
    *   **`data` 模块:** 解析和序列化 JSON 和 YAML 数据。
    *   **`log` 模块:** 以不同的严重级别（info、warn、error、debug）记录消息。
    *   **`salt` 模块:** 直接执行 SaltStack 命令（`salt`、`salt-call`）。
    *   **`gcp` 模块:** 执行谷歌云 (`gcloud`) 命令行指令。
*   **⏰ 任务调度器:** 使用 cron 语法自动执行您的 Lua 任务，作为持久的后台进程运行。
*   **📦 任务产物管理:** 自动收集和存储任务生成的文件和目录，用于审计和重用。
*   **📝 `values.yaml` 集成:** 通过 `values.yaml` 文件将配置值传递给您的 Lua 任务，类似于 Helm。
*   **💻 命令行界面 (CLI):**
    *   `run`: 从 Lua 配置文件执行任务。
    *   `list`: 列出所有可用的任务组和任务及其描述和依赖关系。
    *   `validate`: 验证 Lua 任务文件的语法和结构。
    *   `test`: 执行任务工作流的 Lua 测试文件。
    *   `repl`: 启动交互式 REPL 会话。
    *   `version`: 打印 sloth-runner 的版本号。
    *   `scheduler`: 管理后台任务调度器，包括启用、禁用、列出和删除调度任务。
    *   `template list`: 列出所有可用模板。
    *   `new`: 从模板生成新的任务定义文件。
    *   `check dependencies`: 检查所需的外部 CLI 工具。

## 📚 完整文档

要获取更详细的文档、使用指南和高级示例，请访问我们的[完整文档](./docs/zh/index.md)。

---

## 🚀 开始使用

### 安装

要在您的系统上安装 `sloth-runner`，您可以使用提供的 `install.sh` 脚本。该脚本会自动检测您的操作系统和架构，从 GitHub 下载最新的发布版本，并将 `sloth-runner` 可执行文件放置在 `/usr/local/bin` 中。

```bash
bash <(curl -sL https://raw.githubusercontent.com/chalkan3/sloth-runner/master/install.sh)
```

**注意:** `install.sh` 脚本需要 `sudo` 权限才能将可执行文件移动到 `/usr/local/bin`。

### 基本用法

要运行一个 Lua 任务文件：

```bash
sloth-runner run -f examples/basic_pipeline.lua
```

要列出文件中的任务：

```bash
sloth-runner list -f examples/basic_pipeline.lua
```

---

## 📜 在 Lua 中定义任务

任务在 Lua 文件中定义，通常在一个 `TaskDefinitions` 表中。每个任务可以有 `name`、`description`、`command`（shell 命令字符串或 Lua 函数）、`async`（布尔值）、`pre_exec`（Lua 函数钩子）、`post_exec`（Lua 函数钩子）和 `depends_on`（字符串或字符串表）。

示例 (`examples/basic_pipeline.lua`):

```lua
-- 从另一个文件导入可重用的任务。路径是相对的。
local docker_tasks = import("examples/shared/docker.lua")

TaskDefinitions = {
    full_pipeline_demo = {
        description = "一个演示各种功能的综合管道。",
        tasks = {
            -- 任务 1: 获取数据，异步运行。
            fetch_data = {
                name = "fetch_data",
                description = "从 API 获取原始数据。",
                async = true,
                command = function(params)
                    log.info("正在获取数据...")
                    -- 模拟 API 调用
                    return true, "echo '获取了原始数据'", { raw_data = "api_data" }
                end,
            },

            -- 任务 2: 一个不稳定的任务，失败时会重试。
            flaky_task = {
                name = "flaky_task",
                description = "这个任务会间歇性失败，并且会重试。",
                retries = 3,
                command = function()
                    if math.random() > 0.5 then
                        log.info("不稳定的任务成功。")
                        return true, "echo '成功!'"
                    else
                        log.error("不稳定的任务失败，将重试...")
                        return false, "随机失败"
                    end
                end,
            },

            -- 任务 3: 处理数据，依赖于 fetch_data 和 flaky_task 的成功完成。
            process_data = {
                name = "process_data",
                description = "处理获取的数据。",
                depends_on = { "fetch_data", "flaky_task" },
                command = function(params, deps)
                    local raw_data = deps.fetch_data.raw_data
                    log.info("正在处理数据: " .. raw_data)
                    return true, "echo '处理了数据'", { processed_data = "processed_" .. raw_data }
                end,
            },

            -- 任务 4: 一个带超时的长时间运行任务。
            long_running_task = {
                name = "long_running_task",
                description = "一个如果运行时间过长将被终止的任务。",
                timeout = "5s",
                command = "echo '开始长任务...'; sleep 10; echo '这不会被打印出来。';",
            },

            -- 任务 5: 一个清理任务，如果 long_running_task 失败则运行。
            cleanup_on_fail = {
                name = "cleanup_on_fail",
                description = "仅在长时间运行的任务失败时运行。",
                next_if_fail = "long_running_task",
                command = "echo '由于先前的失败，清理任务已执行。'",
            },

            -- 任务 6: 使用从导入的 docker.lua 模块中可重用的任务。
            build_image = {
                uses = docker_tasks.build,
                description = "构建应用程序的 Docker 镜像。",
                params = {
                    image_name = "my-awesome-app",
                    tag = "v1.2.3",
                    context = "./app_context"
                }
            },

            -- 任务 7: 一个条件任务，仅在文件存在时运行。
            conditional_deploy = {
                name = "conditional_deploy",
                description = "仅在构建产物存在时部署应用程序。",
                depends_on = "build_image",
                run_if = "test -f ./app_context/artifact.txt", -- Shell 命令条件
                command = "echo '正在部署应用程序...'",
            },

            -- 任务 8: 如果满足条件，此任务将中止整个工作流。
            gatekeeper_check = {
                name = "gatekeeper_check",
                description = "如果关键条件未满足，则中止工作流。",
                abort_if = function(params, deps)
                    -- Lua 函数条件
                    log.warn("正在检查守门员条件...")
                    if params.force_proceed ~= "true" then
                        log.error("守门员检查失败。正在中止工作流。")
                        return true -- 中止
                    end
                    return false -- 不中止
                end,
                command = "echo '如果中止，此命令将不会运行。'"
            }
        }
    }
}
```

---

## 📄 模板

`sloth-runner` 提供了几个模板，可以快速搭建新的任务定义文件。

| 模板名称           | 描述                                                                    |
| :----------------- | :----------------------------------------------------------------------------- |
| `simple`           | 生成一个包含“hello world”任务的单一组。非常适合入门。                     |
| `python`           | 创建一个用于设置 Python 环境、安装依赖项和运行脚本的管道。                 |
| `parallel`         | 演示如何并发运行多个任务。                                                  |
| `python-pulumi`    | 使用 Python 管理的 Pulumi 基础设施部署管道。                                |
| `python-pulumi-salt` | 使用 Pulumi 预置基础设施并使用 SaltStack 进行配置。                       |
| `git-python-pulumi` | CI/CD 管道：克隆仓库，设置环境，并使用 Pulumi 进行部署。                   |
| `dummy`            | 生成一个什么都不做的虚拟任务。                                              |

---

## 💻 CLI 命令

`sloth-runner` 提供了一个简单而强大的命令行界面。

### `sloth-runner run`

执行在 Lua 模板文件中定义的任务。

**用法:** `sloth-runner run [flags]`

**描述:**
`run` 命令执行在 Lua 模板文件中定义的任务。
您可以指定文件、环境变量，并针对特定的任务或任务组。

**标志:**

*   `-f, --file string`: Lua 任务配置文件路径 (默认: "examples/basic_pipeline.lua")
*   `-e, --env string`: 任务环境 (例如: Development, Production) (默认: "Development")
*   `-p, --prod`: 设置为 true 表示生产环境 (默认: false)
*   `--shards string`: 分片编号的逗号分隔列表 (例如: 1,2,3) (默认: "1,2,3")
*   `-t, --tasks string`: 要运行的特定任务的逗号分隔列表 (例如: task1,task2)
*   `-g, --group string`: 仅运行特定任务组中的任务
*   `-v, --values string`: 包含要传递给 Lua 任务的值的 YAML 文件路径
*   `-d, --dry-run`: 模拟任务执行而不实际运行它们 (默认: false)
*   `--return`: 以 JSON 格式返回目标任务的输出 (默认: false)
*   `-y, --yes`: 绕过交互式任务选择并运行所有任务 (默认: false)

### `sloth-runner list`

列出所有可用的任务组和任务。

**用法:** `sloth-runner list [flags]`

**描述:**
`list` 命令显示所有任务组及其各自的任务，以及它们的描述和依赖关系。

**标志:**

*   `-f, --file string`: Lua 任务配置文件路径 (默认: "examples/basic_pipeline.lua")
*   `-e, --env string`: 任务环境 (例如: Development, Production) (默认: "Development")
*   `-p, --prod`: 设置为 true 表示生产环境 (默认: false)
*   `--shards string`: 分片编号的逗号分隔列表 (例如: 1,2,3) (默认: "1,2,3")
*   `-v, --values string`: 包含要传递给 Lua 任务的值的 YAML 文件路径

### `sloth-runner validate`

验证 Lua 任务文件的语法和结构。

**用法:** `sloth-runner validate [flags]`

**描述:**
`validate` 命令检查 Lua 任务文件的语法错误，并确保 `TaskDefinitions` 表结构正确。

**标志:**

*   `-f, --file string`: Lua 任务配置文件路径 (默认: "examples/basic_pipeline.lua")
*   `-e, --env string`: 任务环境 (例如: Development, Production) (默认: "Development")
*   `-p, --prod`: 设置为 true 表示生产环境 (默认: false)
*   `--shards string`: 分片编号的逗号分隔列表 (例如: 1,2,3) (默认: "1,2,3")
*   `-v, --values string`: 包含要传递给 Lua 任务的值的 YAML 文件路径

### `sloth-runner test`

执行任务工作流的 Lua 测试文件。

**用法:** `sloth-runner test -w <workflow-file> -f <test-file>`

**描述:**
`test` 命令针对工作流运行指定的 Lua 测试文件。
在测试文件中，您可以使用 `test` 和 `assert` 模块来验证任务行为。

**标志:**

*   `-f, --file string`: Lua 测试文件路径 (必需)
*   `-w, --workflow string`: 要测试的 Lua 工作流文件路径 (必需)

### `sloth-runner repl`

启动交互式 REPL 会话。

**用法:** `sloth-runner repl [flags]`

**描述:**
`repl` 命令启动一个交互式 Read-Eval-Print Loop，允许您执行 Lua 代码并与所有内置的 sloth-runner 模块进行交互。
您可以选择加载一个工作流文件以使其上下文可用。

**标志:**

*   `-f, --file string`: 要加载到 REPL 会话中的 Lua 工作流文件路径

### `sloth-runner version`

打印 sloth-runner 的版本号。

**用法:** `sloth-runner version`

**描述:**
所有软件都有版本。这是 sloth-runner 的版本。

### `sloth-runner template list`

列出所有可用模板。

**用法:** `sloth-runner template list`

**描述:**
显示一个表格，列出所有可用于 `new` 命令的模板。

### `sloth-runner new <group-name>`

从模板生成新的任务定义文件。

**用法:** `sloth-runner new <group-name> [flags]`

**描述:**
`new` 命令创建一个样板 Lua 任务定义文件。
您可以从不同的模板中选择并指定输出文件。
运行 `sloth-runner template list` 查看所有可用模板。

**参数:**

*   `<group-name>`: 要生成的任务组的名称。

**标志:**

*   `-o, --output string`: 输出文件路径 (默认: stdout)
*   `-t, --template string`: 要使用的模板。请参阅 `template list` 获取选项。 (默认: "simple")
*   `--set key=value`: 传递键值对到模板，用于动态内容生成。

### `sloth-runner check dependencies`

检查所需的外部 CLI 工具。

**用法:** `sloth-runner check dependencies`

**描述:**
验证各种模块使用的所有外部命令行工具 (例如: docker, aws, doctl) 是否已安装并在系统的 PATH 中可用。

