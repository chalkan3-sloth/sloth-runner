
---

## Example 3: Generating Tasks with Dynamic Data using Templates

This example demonstrates how to use the `sloth-runner new` command with the `--set` flag to generate a task definition file where content is dynamically injected from the command line. This allows for highly reusable templates that can be customized without modification.

**To generate and run this example:**

1.  **Generate the task file:**
    ```bash
    sloth-runner new templated-task --template simple --set custom_message="This is a custom message from the CLI!" -o examples/templated_task.sloth
    ```
    This command uses the `simple` template and injects `custom_message` into the generated sloth file.

2.  **Run the generated task:**
    ```bash
    sloth-runner run -f examples/templated_task.sloth -g templated-task -t hello_task
    ```
    Observe the output, which should include the custom message you provided.

---

### **Pipeline: `examples/templated_task.sloth`**

```lua
-- examples/templated_task.sloth
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

---

## Example 4: Exploring Workflows with Task IDs using List Command

This example demonstrates how to use the new `sloth-runner list` command to inspect workflow structure, view task relationships, and explore unique IDs for debugging and observability.

### **Creating a Sample Workflow**

First, let's create a comprehensive workflow file to explore:

```lua
-- examples/id_demo.sloth
-- Demonstration of task IDs and workflow structure

TaskDefinitions = {
    -- Build and Deploy Pipeline
    build_pipeline = {
        description = "Build and deployment pipeline with unique IDs",
        tasks = {
            {
                name = "setup",
                description = "Setup build environment",
                command = "echo 'Setting up build environment...'"
            },
            {
                name = "compile",
                description = "Compile the application",
                command = "echo 'Compiling application...'",
                depends_on = {"setup"}
            },
            {
                name = "test",
                description = "Run unit tests",
                command = "echo 'Running tests...'",
                depends_on = {"compile"}
            },
            {
                name = "package",
                description = "Package the application",
                command = "echo 'Packaging application...'",
                depends_on = {"compile", "test"}
            }
        }
    },
    
    -- Deployment Group
    deploy_pipeline = {
        description = "Deployment tasks with environment management",
        tasks = {
            {
                name = "deploy_staging",
                description = "Deploy to staging environment",
                command = "echo 'Deploying to staging...'",
                depends_on = {"package"}
            },
            {
                name = "integration_test",
                description = "Run integration tests",
                command = "echo 'Running integration tests...'",
                depends_on = {"deploy_staging"}
            },
            {
                name = "deploy_production",
                description = "Deploy to production environment", 
                command = "echo 'Deploying to production...'",
                depends_on = {"integration_test"}
            }
        }
    }
}
```

### **Using the List Command**

**1. Basic workflow inspection:**
```bash
sloth-runner list -f examples/id_demo.sloth
```

**Expected output:**
```
Workflow Tasks and Groups

## Task Group: build_pipeline
ID: a1b2c3d4-e5f6-7890-abcd-ef1234567890
Description: Build and deployment pipeline with unique IDs

Tasks:
NAME     ID           DESCRIPTION                 DEPENDS ON
----     --           -----------                 ----------
setup    12345678...  Setup build environment     -
compile  abcdef12...  Compile the application     setup
test     98765432...  Run unit tests              compile
package  fedcba09...  Package the application     compile, test

## Task Group: deploy_pipeline  
ID: f9e8d7c6-b5a4-3210-9876-543210fedcba
Description: Deployment tasks with environment management

Tasks:
NAME                ID           DESCRIPTION                      DEPENDS ON
----                --           -----------                      ----------
deploy_staging      11223344...  Deploy to staging environment   package
integration_test    55667788...  Run integration tests           deploy_staging
deploy_production   99aabbcc...  Deploy to production             integration_test
```

### **Benefits of Task IDs**

**🆔 Unique Identification:**
- Each task and group has a persistent UUID
- IDs remain consistent across executions
- Perfect for debugging and observability

**📊 Enhanced Debugging:**
- Trace specific tasks in logs using IDs
- Identify problematic tasks across multiple runs
- Better correlation with monitoring systems

**🔍 Workflow Inspection:**
- Understand task relationships at a glance
- Verify dependency chains before execution
- Plan execution strategies based on structure

### **Integration with Stack Management**

**Run with stack and inspect:**
```bash
# Run the workflow with a stack
sloth-runner run demo-stack -f examples/id_demo.sloth --output enhanced

# List stacks to see execution history
sloth-runner stack list

# Inspect the workflow structure
sloth-runner list -f examples/id_demo.sloth

# View detailed stack information
sloth-runner stack show demo-stack
```

This workflow demonstrates how task IDs integrate seamlessly with stack management, providing complete traceability from workflow definition to execution history.
