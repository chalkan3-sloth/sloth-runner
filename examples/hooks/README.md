# Event Watcher Hooks

Example hooks responding to watcher events.

## Available Hooks

### File Watchers
- `file_changed_alert.lua` - Alert on file changes
- `file_created_deploy.lua` - Deploy on new files
- `file_deleted_cleanup.lua` - Cleanup on deletion

### Process/Resource Watchers
- `process_monitor.lua` - Process lifecycle monitoring
- `cpu_alert.lua` - High CPU alerts
- `log_pattern_action.lua` - Log pattern responses
- `port_scanner.lua` - Port security monitoring

## Usage

```bash
# Register hooks
sloth-runner hook register examples/hooks/file_changed_alert.lua

# List hooks
sloth-runner hook list

# Test with watchers
sloth-runner run file_watcher_test --file examples/watchers/01_file_watcher.sloth --yes
```
