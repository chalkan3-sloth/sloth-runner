# Core Concepts - Modern DSL

This document explains the fundamental concepts of Sloth-Runner using the **Modern DSL**, helping you understand how tasks are defined and executed with the new fluent API.

## Modern DSL Task Definition

Tasks in Sloth-Runner are now defined using the **Modern DSL** fluent API, which provides a more intuitive and powerful way to create workflows.

### Task Builder Pattern

Each task is created using the `task()` function and the fluent API:

```lua
local my_task = task("task_name")
    :description("Task description")
    :command(function(params, deps)
        -- Task logic here
        return true, "Success message", { output_data = "value" }
    end)
    :timeout("30s")
    :retries(3, "exponential")
    :build()
```

### Workflow Definition Structure

Workflows are defined using `workflow.define()` with comprehensive configuration:

```lua
workflow.define("workflow_name", {
    description = "Workflow description - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Team Name",
        tags = {"tag1", "tag2"},
        created_at = os.date()
    },
    
    tasks = { task1, task2, task3 },
    
    config = {
        timeout = "30m",
        retry_policy = "exponential",
        max_parallel_tasks = 4
    },
    
    on_start = function()
        log.info("Starting workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("Workflow completed successfully!")
        end
        return true
    end
})
```

## Task Properties and Methods

Each task can use the following fluent API methods:

### Basic Properties
-   `:name(string)` - Task name (usually set in task() constructor)
-   `:description(string)` - Brief description of what the task does
-   `:command(function|string)` - Task execution logic:
    -   `function(params, deps)` - Lua function with parameters and dependencies
    -   `string` - Shell command to execute
    
### Execution Control
-   `:timeout(string)` - Execution timeout (e.g., "10s", "1m", "30m")
-   `:retries(number, strategy)` - Retry count and strategy ("exponential", "linear", "fixed")
-   `:async(boolean)` - Whether to execute asynchronously
-   `:depends_on(table)` - Array of task names this task depends on

### Conditional Execution
-   `:run_if(function|string)` - Execute only if condition is true
-   `:abort_if(function|string)` - Abort entire workflow if condition is true
-   `:condition(function)` - Advanced conditional logic

### Lifecycle Hooks
-   `:pre_hook(function)` - Execute before main command
-   `:post_hook(function)` - Execute after main command  
-   `:on_success(function)` - Execute on successful completion
-   `:on_failure(function)` - Execute on failure
-   `:on_timeout(function)` - Execute on timeout

### Artifact Management
-   `:artifacts(table)` - Files to save as artifacts after execution
-   `:consumes(table)` - Artifacts from other tasks to consume

### Advanced Features
-   `:circuit_breaker(config)` - Circuit breaker configuration
-   `:performance_monitoring(boolean)` - Enable performance tracking
-   `:environment(table)` - Environment variables for task execution

## Artifact Management

Sloth-Runner allows tasks to share files through an artifact mechanism using the Modern DSL. A task can "produce" one or more files as artifacts, and subsequent tasks can "consume" those artifacts.

This is useful for CI/CD pipelines where a build step might generate a binary (artifact) that is then used by a test or deployment step.

### How It Works

1.  **Producing Artifacts:** Use the `:artifacts()` method in your task definition. The value can be a single file pattern (e.g., `"report.txt"`) or a list (e.g., `{"*.log", "app.bin"}`). After the task completes successfully, the runner will look for files in the task's workdir that match these patterns and copy them to shared artifact storage.

2.  **Consuming Artifacts:** Use the `:consumes()` method in another task definition (which typically `:depends_on()` the producer task). The value should be the artifact file name you want to use (e.g., `"report.txt"`). Before this task executes, the runner will copy the named artifact from shared storage to this task's workdir.

### Modern DSL Artifact Example

```lua
-- Producer task that creates artifacts
local build_task = task("build_app")
    :description("Build application and create artifacts")
    :command(function()
        -- Build the application
        local result = exec.run("go build -o myapp ./cmd/main.go")
        if not result.success then
            return false, "Build failed: " .. result.stderr
        end
        
        -- Create a report file
        fs.write_file("build-report.txt", "Build completed at " .. os.date())
        
        return true, "Build completed successfully", {
            binary_size = fs.size("myapp"),
            build_time = result.duration
        }
    end)
    :artifacts({"myapp", "build-report.txt"}) -- These files will be saved as artifacts
    :timeout("5m")
    :build()

-- Consumer task that uses artifacts
local test_task = task("test_app")
    :description("Test the built application")
    :depends_on({"build_app"})
    :consumes({"myapp", "build-report.txt"}) -- These artifacts will be available
    :command(function(params, deps)
        -- The artifacts are now available in the workdir
        log.info("Testing application: myapp")
        log.info("Build report: " .. fs.read_file("build-report.txt"))
        
        -- Run tests on the binary
        local result = exec.run("./myapp --version")
        return result.success, "Tests completed", {
            version_output = result.stdout
        }
    end)
    :build()

-- Deploy task that also uses the binary
local deploy_task = task("deploy_app")
    :description("Deploy the built application")
    :depends_on({"test_app"})
    :consumes({"myapp"})
    :command(function()
        log.info("Deploying application...")
        -- Copy binary to deployment location
        local result = exec.run("cp myapp /usr/local/bin/")
        return result.success, "Deployment completed"
    end)
    :build()

-- Complete workflow
workflow.define("ci_cd_pipeline", {
    description = "CI/CD pipeline with artifact management",
    version = "1.0.0",
    
    tasks = { build_task, test_task, deploy_task },
    
    config = {
        create_workdir_before_run = true,
        cleanup_artifacts_after = "7d"
    }
})
```

### Key Benefits

- **🔄 Automatic Dependency Resolution**: Tasks automatically get the artifacts they need
- **📦 Efficient Storage**: Artifacts are stored in a shared space, reducing duplication
- **🧹 Automatic Cleanup**: Artifacts can be automatically cleaned up after a specified period
- **📊 Rich Metadata**: Artifacts include metadata like creation time, size, and source task

## Modern DSL vs Legacy Format

The Modern DSL provides several advantages over the legacy TaskDefinitions format:

| Feature | Legacy Format | Modern DSL |
|---------|---------------|------------|
| **Syntax** | Table-based, procedural | Fluent API, chainable |
| **Type Safety** | Runtime discovery | Compile-time validation |
| **Error Handling** | Basic | Enhanced with context |
| **Metadata** | Limited | Rich, structured |
| **Retry Logic** | Manual implementation | Built-in strategies |
| **Dependencies** | Simple strings | Advanced with conditions |
| **Lifecycle Hooks** | Basic pre/post | Rich event handling |
| **Testing** | Manual | Integrated framework |

## Next Steps

- **📚 Learn More**: Check out the [Modern DSL Introduction](modern-dsl/introduction.md)
- **🎯 API Reference**: See [Task Definition API](modern-dsl/task-api.md) for complete reference
- **📝 Examples**: Browse [Examples](../examples/) for real-world Modern DSL workflows
- **🔧 Migration**: Use [Migration Guide](modern-dsl/migration-guide.md) to convert existing workflows

The Modern DSL represents the future of workflow automation in Sloth Runner - more powerful, intuitive, and maintainable!
          })

          local spoke_config = values.pulumi.spoke.config
          spoke_config.hub_network_self_link = hub_outputs.network_self_link -- Use hub output in spoke config

          spoke_stack:select():config_map(spoke_config)
          local spoke_result = spoke_stack:up({ yes = true })
          if not spoke_result.success then
            log.error("Spoke stack deployment failed: " .. spoke_result.stdout)
            return false, "Spoke stack deployment failed."
          end
          log.info("Spoke stack deployed successfully.")
          local spoke_outputs = spoke_stack:outputs()
          return true, "Spoke stack deployed.", { outputs = spoke_outputs }
        end
      },
      {
          name = "final_summary",
          depends_on = "deploy_spoke_stack", -- Depends on the final deployment task
          command = function(inputs)
              log.info("GCP Hub and Spoke orchestration completed successfully!")
              -- You can access outputs from dependencies like this:
              -- local hub_outputs = inputs.deploy_hub_stack.outputs
              -- local spoke_outputs = inputs.deploy_spoke_stack.outputs
              return true, "Orchestration successful."
          end
      }
    }
  }
}
```

## Parâmetros e Outputs

*   **Parâmetros (`params`):** Podem ser passados para as tarefas via linha de comando ou definidos na própria tarefa. A função `command` e as funções `run_if`/`abort_if` podem acessá-los.
*   **Outputs (`deps`):** As funções Lua de `command` podem retornar uma tabela de outputs. Tarefas que dependem desta tarefa podem acessar esses outputs através do argumento `deps`.

## Exportando Dados para a CLI

Além dos outputs de tarefas, o `sloth-runner` fornece uma função global `export()` que permite passar dados de dentro de um script diretamente para a saída da linha de comando.

### `export(tabela)`

*   **`tabela`**: Uma tabela Lua cujos pares de chave-valor serão exportados.

Quando você executa uma tarefa com a flag `--return`, os dados passados para a função `export()` serão mesclados com o output da tarefa final e impressos como um único objeto JSON. Se houver chaves duplicadas, o valor da função `export()` terá precedência.

Isso é útil para extrair informações importantes de qualquer ponto do seu script, não apenas do valor de retorno da última tarefa.

**Exemplo:**

```lua
command = function(params, deps)
  -- Lógica da tarefa...
  local some_data = {
    info = "Este é um dado importante",
    timestamp = os.time()
  }
  
  -- Exporta a tabela
  export(some_data)
  
  -- A tarefa pode continuar e retornar seu próprio output
  return true, "Tarefa concluída", { status = "ok" }
end
```

Executando com `--return` resultaria em uma saída JSON como:
```json
{
  "info": "Este é um dado importante",
  "timestamp": 1678886400,
  "status": "ok"
}
```

## Módulos Built-in

O Sloth-Runner expõe várias funcionalidades Go como módulos Lua, permitindo que suas tarefas interajam com o sistema e serviços externos. Além dos módulos básicos (`exec`, `fs`, `net`, `data`, `log`, `import`, `parallel`), o Sloth-Runner agora inclui módulos avançados para Git, Pulumi e Salt.

Esses módulos oferecem uma API fluente e intuitiva para automação complexa.

*   **`exec` module:** Para executar comandos de shell arbitrários.
*   **`fs` module:** Para operações de sistema de arquivos (leitura, escrita, etc.).
*   **`net` module:** Para fazer requisições HTTP e downloads.
*   **`data` module:** Para parsear e serializar JSON e YAML.
*   **`log` module:** Para registrar mensagens no console do Sloth-Runner.
*   **`import` function:** Para importar outros arquivos Lua e reutilizar tarefas.
*   **`parallel` function:** Para executar tarefas em paralelo.
*   **`git` module:** Para interagir com repositórios Git.
*   **`pulumi` module:** Para orquestrar stacks do Pulumi.
*   **`salt` module:** Para executar comandos SaltStack.

Para detalhes sobre cada módulo, consulte suas respectivas seções na documentação.

---

[Voltar ao Índice](./index.md)
