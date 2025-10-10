# Interactive REPL

The `sloth-runner repl` command drops you into an interactive Read-Eval-Print Loop (REPL) session. This is a powerful tool for debugging, exploration, and quick experimentation with the sloth-runner modules.

## Starting the REPL

To start a session, simply run:
```bash
sloth-runner repl
```

You can also pre-load a workflow file to have its tasks, workflows, and helper functions available in the session. This is incredibly useful for debugging an existing pipeline.

```bash
sloth-runner repl -f /path/to/your/pipeline.sloth
```

## Features

### Live Environment
The REPL provides a live Lua environment where you can execute any Lua code. All the built-in sloth-runner modules (`aws`, `docker`, `fs`, `log`, etc.) are pre-loaded and ready to use.

```
sloth> log.info("Hello from the REPL!")
sloth> result = fs.read("README.md")
sloth> print(string.sub(result, 1, 50))
```

### Autocompletion
The REPL has a sophisticated autocompletion system.
- Start typing the name of a global variable or module (e.g., `aws`) and press `Tab` to see suggestions.
- Type a module name followed by a dot (e.g., `docker.`) and press `Tab` to see all the functions available in that module.

### History
The REPL keeps a history of your commands. Use the up and down arrow keys to navigate through previous commands.

## Example Sessions

### Testing Module Functions

Here is an example of using the REPL to debug a Docker command.

```bash
$ sloth-runner repl
Sloth-Runner Interactive REPL
Type 'exit' or 'quit' to leave.
sloth> result = docker.exec({"ps", "-a"})
sloth> print(result.stdout)
CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
sloth> -- Now let's try to build an image
sloth> build_result = docker.build({tag="my-test", path="./examples/docker"})
sloth> print(build_result.success)
true
sloth> exit
Bye!
```

### Testing Modern DSL Tasks

You can also test and debug tasks using the modern DSL pattern in the REPL:

```bash
$ sloth-runner repl
Sloth-Runner Interactive REPL
Type 'exit' or 'quit' to leave.
sloth> -- Create a task using modern DSL
sloth> my_task = task("test_deploy")
sloth>     :description("Test deployment task")
sloth>     :command(function(this, params)
sloth>         log.info("Running deployment...")
sloth>         return true, "Deployment successful"
sloth>     end)
sloth>     :timeout("5m")
sloth>     :build()
sloth>
sloth> -- Test the task execution
sloth> result = my_task:execute()
sloth> print(result)
sloth> exit
Bye!
```

### Testing Workflows

You can build and test complete workflows interactively:

```bash
$ sloth-runner repl
Sloth-Runner Interactive REPL
Type 'exit' or 'quit' to leave.
sloth> -- Define a simple task
sloth> check_task = task("check")
sloth>     :command(function(this, params)
sloth>         log.info("Checking environment...")
sloth>         return true, "Environment OK"
sloth>     end)
sloth>     :build()
sloth>
sloth> -- Create a workflow with the task
sloth> my_workflow = workflow.define("test_workflow")
sloth>     :description("Test workflow")
sloth>     :version("1.0.0")
sloth>     :tasks({ check_task })
sloth>
sloth> -- Execute the workflow
sloth> my_workflow:run()
sloth> exit
Bye!
```
