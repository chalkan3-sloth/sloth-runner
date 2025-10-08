# Event Watchers Examples

This directory contains examples of all available event watchers in sloth-runner. Watchers monitor specific conditions on agents and automatically send events to the master when those conditions are met.

## Available Watchers

### 1. File Watcher (`01_file_watcher.sloth`)
Monitors individual files for changes.

**Example:**
```lua
event.register.file({
    file_path = "/tmp/watch-test.txt",
    when = {'created', 'changed', 'deleted'},
    check_hash = true,
    interval = "3s"
})
```

**Events Generated:**
- `file.created` - File was created
- `file.changed` - File was modified (size, mtime, or hash changed)
- `file.deleted` - File was deleted

### 2. Directory Watcher (`02_directory_watcher.sloth`)
Monitors directories for file additions/removals.

**Example:**
```lua
event.register.directory({
    directory_path = "/tmp/watch-dir",
    when = {'created', 'deleted'},
    pattern = "*.txt",
    interval = "3s"
})
```

**Events Generated:**
- `directory.file_created` - New file added to directory
- `directory.file_deleted` - File removed from directory

### 3. Process Watcher (`03_process_watcher.sloth`)
Monitors process start/stop.

**Example:**
```lua
event.register.process({
    process_name = "nginx",
    when = {'created', 'deleted'},
    interval = "2s"
})
```

**Events Generated:**
- `process.started` - Process started
- `process.stopped` - Process stopped

### 4. Log Watcher (`04_log_watcher.sloth`)
Monitors log files for pattern matches.

**Example:**
```lua
event.register.log({
    log_path = "/var/log/app.log",
    pattern = "ERROR|CRITICAL",
    when = {'matches'},
    follow = true,
    interval = "2s"
})
```

**Events Generated:**
- `log.pattern_matched` - Regex pattern matched in log line
- `log.line_matched` - String found in log line

### 5. CPU Watcher (`05_cpu_watcher.sloth`)
Monitors CPU usage thresholds.

**Example:**
```lua
event.register.cpu({
    threshold = 80,
    when = {'above'},
    interval = "5s"
})
```

**Events Generated:**
- `cpu.high_usage` - CPU usage above threshold

### 6. Command Watcher (`06_command_watcher.sloth`)
Executes commands and monitors output.

**Example:**
```lua
event.register.command({
    command = "df /tmp | tail -1 | awk '{print $5}' | tr -d '%'",
    threshold = 80,
    when = {'above'},
    interval = "10s"
})
```

**Events Generated:**
- `command.unexpected_exit` - Command exited with unexpected code
- `command.output_matched` - Output matched pattern
- `command.threshold_exceeded` - Numeric output exceeded threshold
- `command.value_increased` - Value increased since last check

## Additional Watchers

### Port Watcher
```lua
event.register.port({
    port = 8080,
    when = {'created', 'deleted'},
    interval = "5s"
})
```

### Memory Watcher
```lua
event.register.memory({
    threshold = 85,
    when = {'above'},
    interval = "10s"
})
```

### Disk Watcher
```lua
event.register.disk({
    mount_point = "/",
    threshold = 90,
    when = {'above'},
    interval = "30s"
})
```

### Network Watcher
```lua
event.register.network({
    threshold = 1000000,  -- bytes/sec
    when = {'above'},
    interval = "10s"
})
```

### Connection Watcher
```lua
event.register.connection({
    protocol = "tcp",
    state = "ESTABLISHED",
    when = {'changed'},
    interval = "10s"
})
```

### User Watcher
```lua
event.register.user({
    username = "root",
    when = {'created', 'deleted'},
    interval = "5s"
})
```

### Package Watcher
```lua
event.register.package({
    package_name = "nginx",
    when = {'created', 'deleted'},
    interval = "60s"
})
```

### Custom Watcher
```lua
event.register.custom({
    check = function()
        -- Your custom logic here
        local result = some_check()
        if result > threshold then
            return true, {value = result, message = "Threshold exceeded"}
        end
        return false, nil
    end,
    interval = "30s"
})
```

## Running Examples

### On a Specific Agent
```bash
sloth-runner run file_watcher_test --file examples/watchers/01_file_watcher.sloth --delegate-to lady-guica --yes
```

### Check Events on Master
```bash
sqlite3 /etc/sloth-runner/hooks.db "SELECT * FROM events ORDER BY created_at DESC LIMIT 10"
```

## Event Conditions

- `created` - Resource was created
- `deleted` - Resource was deleted
- `changed` - Resource changed
- `exists` - Resource exists
- `above` - Value above threshold
- `below` - Value below threshold
- `matches` - Pattern matches (regex)
- `contains` - Contains string
- `increased` - Value increased
- `decreased` - Value decreased

## Notes

- Watchers run continuously on agents
- Events are buffered (batch size: 50) and sent to master
- All events are stored in `/etc/sloth-runner/hooks.db` on master
- Watchers survive workflow completion (until agent restart)
- Use appropriate intervals to avoid excessive checking
