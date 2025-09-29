
---

## Example 3: Generating Tasks with Dynamic Data using Templates

This example demonstrates how to use the `sloth-runner new` command with the `--set` flag to generate a task definition file where content is dynamically injected from the command line. This allows for highly reusable templates that can be customized without modification.

**To generate and run this example:**

1.  **Generate the task file:**
    ```bash
    sloth-runner new templated-task --template simple --set custom_message="This is a custom message from the CLI!" -o examples/templated_task.lua
    ```
    This command uses the `simple` template and injects `custom_message` into the generated Lua file.

2.  **Run the generated task:**
    ```bash
    sloth-runner run -f examples/templated_task.lua -g templated-task -t hello_task
    ```
    Observe the output, which should include the custom message you provided.

---

### **Pipeline: `examples/templated_task.lua`**

```lua
-- examples/templated_task.lua
--
-- This file is generated using 'sloth-runner new' with the --set flag.
-- It demonstrates how to inject dynamic data into templates using Modern DSL.

-- Define the hello task with Modern DSL
local hello_task = task("hello_task")
    :description("An example task with a custom message - Modern DSL")
    :command(function(params)
        local workdir = params.workdir or "."
        log.info("Running Modern DSL task in: " .. workdir)
        log.info("Custom message: This is a custom message from the CLI!")
        
        local result = exec.run("echo 'Hello from sloth-runner Modern DSL!'")
        if not result.success then
            log.error("Failed to run example task: " .. result.stderr)
            return false, "Task failed", { error = result.stderr }
        else
            log.info("Example task completed successfully")
            print("Command output: " .. result.stdout)
            return true, "Task executed successfully", { 
                output = result.stdout,
                custom_message = "This is a custom message from the CLI!"
            }
        end
    end)
    :timeout("30s")
    :build()

-- Define workflow with Modern DSL
workflow.define("templated-task", {
    description = "A task group generated with dynamic data - Modern DSL",
    version = "1.0.0",
    
    metadata = {
        author = "Sloth Runner CLI",
        tags = {"templated", "example", "modern-dsl"},
        template_source = "simple"
    },
    
    tasks = { hello_task },
    
    config = {
        timeout = "5m"
    },
    
    on_start = function()
        log.info("🚀 Starting templated task workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ Templated task workflow completed successfully!")
        end
        return true
    end
})
          end
        end
      }
    }
  }
}
```
