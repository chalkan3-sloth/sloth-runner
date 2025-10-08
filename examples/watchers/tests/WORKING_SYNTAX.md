# Watcher System - Working Syntax Guide

## ✅ Status: Watcher System is Functional!

The watcher system is fully operational. Watchers registered via `delegate_to` correctly register on the remote agent where the code executes.

## 🎯 Correct Workflow Syntax

### Task Definition Pattern

```lua
local task_name = task("task_identifier")
    :description("Task description")
    :command(function()
        -- Your code here
        -- This code runs on the agent where it's delegated

        -- Register watcher
        local watcher_id = watcher.register.file({
            file_path = '/path/to/file',
            when = {'created', 'changed', 'deleted'},
            check_hash = true,
            interval = '2s'
        })

        return true, "Success message"
    end)
    :depends_on({"previous_task"})  -- Optional
    :build()  -- REQUIRED!
```

### Workflow Definition Pattern

```lua
workflow
    .define("workflow_name")  -- Use .define() not ("name")
    :description("Workflow description")
    :version("1.0.0")
    :tasks({
        task1,
        task2,
        task3
    })
    :config({
        timeout = "5m"
    })
    :on_complete(function(success, results)
        if success then
            print("✅ Success!")
        else
            print("❌ Failed!")
        end
    end)
```

## 🚀 Usage

### Running with delegate_to

```bash
sloth-runner run workflow_name --file test.sloth --delegate-to lady-guica --yes
```

This will:
1. Send the workflow to lady-guica
2. Execute all tasks on lady-guica
3. Any `watcher.register.*()` calls will register on lady-guica automatically
4. No need to specify agent IP in watcher registration!

## ❌ Common Mistakes

### ❌ Wrong: Using workflow() instead of workflow.define()

```lua
workflow("name")  -- WRONG! Doesn't work
```

### ❌ Wrong: Not calling :build() on tasks

```lua
local task1 = task("task1")
    :command(function() ... end)
    -- Missing :build()!
```

### ❌ Wrong: Trying to use task() inside :tasks({})

```lua
workflow.define("test")
    :tasks({
        task("inline_task")  -- WRONG! task() not available here
            :command(...)
    })
```

### ❌ Wrong: Using delegate_to inside group (not supported)

```lua
group "test" {
    delegate_to = {"agent"}  -- WRONG! Syntax error
    task "..." { ... }
}
```

## ✅ Working Example

See: `examples/watchers/tests/working_file_watcher_test.sloth`

This example:
- ✅ Uses correct workflow.define() syntax
- ✅ Properly defines tasks with :build()
- ✅ Registers file watcher on remote agent
- ✅ Tests create, modify, delete events
- ✅ Works with --delegate-to flag

## 🔍 How Automatic Registration Works

When you use `--delegate-to agent-name`:

1. The entire Lua script is sent to the remote agent
2. The agent executes the Lua code locally
3. When `watcher.register.*()` is called, it stores in the local `_WATCHERS` table
4. The agent's watcher manager picks up watchers from `_WATCHERS`
5. Watchers are registered and start monitoring locally on that agent

**No manual agent specification needed!** The watcher automatically registers where the code executes.

## 📝 Next Steps

To create more watcher tests, copy `working_file_watcher_test.sloth` and modify:
- Change watcher type (cpu, memory, process, port, file)
- Adjust watcher parameters
- Modify test logic
- Keep the same workflow structure

## 🐛 Troubleshooting

If you get "no workflows found":
- Check you're using `workflow.define()` not `workflow()`
- Verify all tasks have `:build()` at the end
- Ensure tasks are defined as local variables before workflow

If watchers don't register:
- Verify agent is running: `ssh agent "ps aux | grep sloth-runner"`
- Check agent logs: `ssh agent "tail -f agent.log"`
- Confirm --delegate-to agent name is correct
- Test connection: `nc -zv agent-ip agent-port`
