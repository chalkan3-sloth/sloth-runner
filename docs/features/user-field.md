# Task User Field

The `:user()` method allows you to specify which user should run the task on the agent. By default, tasks run as the `root` user (or the user running the agent).

## Syntax

```lua
task("task-name")
    :user("username")
    :command(function()
        -- Task code here
    end)
    :build()
```

## Parameters

- `username` (string): The username to run the task as. Must be a valid user on the target system.

## Example: Running Task as Specific User

```lua
task("deploy-app")
    :description("Deploy application as app user")
    :delegate_to("prod-server")
    :user("appuser")  -- Run as 'appuser' instead of root
    :command(function()
        -- This will execute as 'appuser'
        cmd.execute("cd /app && ./deploy.sh")
    end)
    :build()
```

## Example: Different Users for Different Tasks

```lua
workflow.define("multi-user-deploy")
    :description("Deploy with different users")
    :tasks({
        -- Database migration as postgres user
        task("db-migrate")
            :delegate_to("db-server")
            :user("postgres")
            :command(function()
                cmd.execute("psql -c 'SELECT version();'")
            end)
            :build(),
            
        -- Application deployment as app user
        task("app-deploy")
            :delegate_to("app-server")
            :user("appuser")
            :command(function()
                file_ops.copy({
                    src = "/tmp/app.tar.gz",
                    dest = "/opt/app/releases/",
                    mode = "0644"
                })
                cmd.execute("cd /opt/app && ./deploy.sh")
            end)
            :build(),
            
        -- Nginx config as root (default)
        task("nginx-reload")
            :delegate_to("web-server")
            -- No :user() specified, runs as root
            :command(function()
                cmd.execute("nginx -t && systemctl reload nginx")
            end)
            :build()
    })
    :on_complete(function()
        print("✅ Multi-user deployment completed")
    end)
```

## Implementation Details

- When `user` is specified and is not `root`, the agent uses `sudo -u <user>` to run the command
- The specified user must exist on the target system
- The agent must have permission to run commands as that user (usually requires agent to run as root)
- If no user is specified, the task runs as the agent's user (typically root)

## Security Considerations

- Ensure the specified user has appropriate permissions for the task
- Consider using dedicated service accounts for specific tasks
- Avoid running unnecessary tasks as root when a less privileged user would suffice
- The agent must be running as root or have sudo privileges to switch users

## Related

- `:delegate_to()` - Specify which agent runs the task
- `:workdir()` - Specify the working directory for the task
