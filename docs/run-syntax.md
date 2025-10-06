# Sloth Runner `run` Command - Complete Syntax Guide

## Executive Summary

The `sloth-runner run` command is the core execution engine that orchestrates task automation both locally and remotely via SSH. This document provides the **definitive syntax reference** with special emphasis on secure password handling via stdin.

## ⚡ Quick Reference

### Basic Syntax Structure

```bash
sloth-runner run <stack-name> \
  --file <sloth-file> \
  [--ssh <profile-name>] \
  [--ssh-password-stdin -] \
  [< password-source] \
  ['<remote-command>']
```

### Command Anatomy

```
sloth-runner run production-deploy    # Stack identifier
  --file deploy.sloth                 # Task definition file
  --ssh prod-server                   # SSH profile name
  --ssh-password-stdin -               # Password from stdin flag
  < pass.txt                           # Password source (redirection)
  'sudo systemctl restart nginx'       # Remote command to execute
```

## 📋 Complete Parameter Reference

### Positional Arguments

#### `<stack-name>` (Required)

The stack name serves as a unique identifier for this execution context, maintaining state across runs.

**Format:** Alphanumeric characters, hyphens, and underscores
**Length:** 1-50 characters
**Pattern:** `^[a-zA-Z][a-zA-Z0-9_-]*$`

**Examples:**
```bash
production-deploy
staging-2024-01-15
maintenance-task-001
db-backup-daily
```

**Best Practice:** Use timestamp for uniqueness
```bash
STACK="deploy-$(date +%Y%m%d-%H%M%S)"
sloth-runner run "$STACK" --file deploy.sloth
```

#### `'<remote-command>'` (Optional with SSH)

The command to execute on the remote host. Required when using `--ssh`.

**Important:** Always quote complex commands to prevent shell interpretation issues.

```bash
# Simple command
'ls -la'

# Complex command with pipes
'ps aux | grep nginx | grep -v grep'

# Multi-line command
'cd /app && \
 git pull origin main && \
 npm install && \
 npm run build'

# Command with special characters
'echo "Status: $(systemctl is-active nginx)"'
```

### Required Flags

#### `--file <sloth-file>` (Required)

Path to the Sloth Runner configuration file containing task definitions.

**Path Resolution:**
- Relative: Resolved from current working directory
- Absolute: Used as-is
- Home expansion: `~/` is expanded to user home

**Examples:**
```bash
--file deploy.sloth              # Relative path
--file ./tasks/deploy.sloth      # Relative with directory
--file /opt/tasks/deploy.sloth   # Absolute path
--file ~/sloth/deploy.sloth      # Home directory
--file $TASK_DIR/deploy.sloth    # Environment variable
```

### SSH-Related Flags

#### `--ssh <profile-name>` (Required for Remote Execution)

Specifies the SSH profile to use for remote execution. The profile must exist in the SQLite database.

**Validation:**
- Profile must exist (`sloth-runner ssh list` to verify)
- Profile must have valid connection details
- Key file (if specified) must be accessible

**Example:**
```bash
--ssh production-web
--ssh staging-database
--ssh bastion-host
```

#### `--ssh-password-stdin -` (Optional)

Instructs the command to read the SSH password from standard input. The hyphen (`-`) is **mandatory** and indicates stdin as the source.

**Critical Requirements:**
1. Must be followed by exactly one hyphen: `-`
2. Cannot be used with key-based authentication
3. Must be paired with input redirection or pipe
4. Password must not contain trailing newline

**Syntax Variations:**
```bash
# From file
--ssh-password-stdin - < pass.txt

# From pipe
echo -n "password" | sloth-runner run stack --file task.sloth --ssh server --ssh-password-stdin -

# From process substitution
--ssh-password-stdin - < <(get-password-command)

# From here-doc (NOT RECOMMENDED - visible in process list)
--ssh-password-stdin - <<< "password"
```

## 🔐 Password Handling - Critical Security Section

### The Golden Rules

1. **NEVER** put passwords directly in the command line
2. **ALWAYS** use `echo -n` or `printf` to avoid newlines
3. **IMMEDIATELY** destroy password files after use
4. **NEVER** commit password files to version control

### Method 1: File Redirection (Most Common)

```bash
# Step 1: Create password file securely
touch pass.txt
chmod 600 pass.txt  # CRITICAL: Set permissions BEFORE writing password
echo -n "MySecureP@ssw0rd" > pass.txt

# Step 2: Use the password
sloth-runner run deployment \
  --file deploy.sloth \
  --ssh production \
  --ssh-password-stdin - < pass.txt \
  'sudo systemctl restart application'

# Step 3: Securely destroy the file
shred -vzu pass.txt  # Overwrite and remove
```

### Method 2: Environment Variable (CI/CD)

```bash
# In CI/CD pipeline (e.g., GitHub Actions, Jenkins)
export SSH_PASSWORD="$SECRET_SSH_PASSWORD"  # From secure vault

# Execute command
sloth-runner run "$BUILD_NUMBER" \
  --file deploy.sloth \
  --ssh production \
  --ssh-password-stdin - < <(echo -n "$SSH_PASSWORD") \
  'deployment-script.sh'

# Clear from environment
unset SSH_PASSWORD
```

### Method 3: Password Manager Integration

```bash
# Using 'pass' password manager
pass show ssh/production | sloth-runner run deployment \
  --file deploy.sloth \
  --ssh production \
  --ssh-password-stdin - \
  'restart-services.sh'

# Using 1Password CLI
op get item "Production SSH" --fields password | \
  sloth-runner run deployment \
  --file deploy.sloth \
  --ssh production \
  --ssh-password-stdin - \
  'deployment-task.sh'

# Using HashiCorp Vault
vault kv get -field=password secret/ssh/production | \
  sloth-runner run deployment \
  --file deploy.sloth \
  --ssh production \
  --ssh-password-stdin - \
  'secure-task.sh'
```

### Method 4: Interactive Input (Manual Only)

```bash
# For manual, one-time execution
sloth-runner run manual-task \
  --file task.sloth \
  --ssh server \
  --ssh-password-stdin -

# Terminal will wait for input
# Type password and press Ctrl+D (EOF)
# Password will not be echoed
```

### Common Password File Errors and Solutions

#### Error: Authentication Failed Due to Newline

**Problem:**
```bash
echo "password" > pass.txt  # WRONG: adds newline
```

**Solution:**
```bash
echo -n "password" > pass.txt  # CORRECT: no newline
# OR
printf "password" > pass.txt   # CORRECT: no newline
```

**Verification:**
```bash
# Check for newline (0a in hex)
hexdump -C pass.txt
# Last byte should NOT be 0a

# Check file size
wc -c pass.txt
# Should match exact password length
```

#### Error: Password Visible in Process List

**Problem:**
```bash
sloth-runner run stack --ssh server --password "visible"  # NEVER DO THIS
```

**Solution:**
```bash
sloth-runner run stack --ssh server --ssh-password-stdin - < pass.txt
```

#### Error: Password File Permissions Too Open

**Problem:**
```bash
echo -n "password" > pass.txt  # Created with default permissions (644)
```

**Solution:**
```bash
touch pass.txt
chmod 600 pass.txt  # Set permissions FIRST
echo -n "password" > pass.txt
```

## 🎯 Complete Usage Examples

### Example 1: Basic Local Execution

```bash
# No SSH, local execution only
sloth-runner run local-task \
  --file maintenance.sloth
```

### Example 2: Remote Execution with SSH Key

```bash
# Using SSH key authentication (recommended)
sloth-runner run production-deploy \
  --file deploy.sloth \
  --ssh prod-web \
  'cd /app && git pull && docker-compose up -d'
```

### Example 3: Remote Execution with Password

```bash
#!/bin/bash
# secure-deploy.sh

set -euo pipefail

# Configuration
STACK="deploy-$(date +%Y%m%d-%H%M%S)"
TASK_FILE="deploy.sloth"
SSH_PROFILE="legacy-server"
PASSWORD_FILE="/tmp/.ssh_pass_$$"  # Unique temp file

# Create secure password file
touch "$PASSWORD_FILE"
chmod 600 "$PASSWORD_FILE"
echo -n "SecurePassword123!" > "$PASSWORD_FILE"

# Execute deployment
sloth-runner run "$STACK" \
  --file "$TASK_FILE" \
  --ssh "$SSH_PROFILE" \
  --ssh-password-stdin - < "$PASSWORD_FILE" \
  'sudo /opt/scripts/deploy.sh'

# Clean up immediately
shred -u "$PASSWORD_FILE"
```

### Example 4: Multi-Command Execution

```bash
# Execute multiple commands in sequence
sloth-runner run maintenance \
  --file maintenance.sloth \
  --ssh production \
  'cd /app && \
   echo "Starting maintenance..." && \
   docker-compose down && \
   docker system prune -f && \
   docker-compose up -d && \
   echo "Maintenance complete"'
```

### Example 5: Conditional Execution

```bash
# Execute based on condition
sloth-runner run conditional-deploy \
  --file deploy.sloth \
  --ssh staging \
  'if [ -f /app/.deploy_lock ]; then \
     echo "Deploy locked, aborting"; \
     exit 1; \
   else \
     touch /app/.deploy_lock && \
     /opt/deploy.sh && \
     rm /app/.deploy_lock; \
   fi'
```

### Example 6: Pipeline Integration

```bash
# GitHub Actions example
- name: Deploy to Production
  env:
    SSH_PASSWORD: ${{ secrets.SSH_PASSWORD }}
  run: |
    # Create temporary password file
    PASS_FILE="$(mktemp)"
    chmod 600 "$PASS_FILE"
    echo -n "$SSH_PASSWORD" > "$PASS_FILE"

    # Run deployment
    sloth-runner run "deploy-${{ github.run_id }}" \
      --file .sloth/deploy.sloth \
      --ssh production \
      --ssh-password-stdin - < "$PASS_FILE" \
      'sudo systemctl restart application'

    # Cleanup
    shred -u "$PASS_FILE"
```

## 🔄 Execution Flow Diagram

```
┌──────────────────┐
│  User Executes   │
│  sloth-runner    │
│      run         │
└────────┬─────────┘
         │
         v
┌──────────────────┐
│  Parse Command   │
│   Line Args      │
└────────┬─────────┘
         │
         v
┌──────────────────┐     ┌──────────────────┐
│  Load Stack      │────>│  Create New      │
│  Configuration   │     │  Stack if Needed │
└────────┬─────────┘     └──────────────────┘
         │
         v
┌──────────────────┐
│  Load Sloth      │
│  File (.sloth)   │
└────────┬─────────┘
         │
         v
      ┌──┴──┐
      │ SSH │ No
      │Flag?├──────────┐
      └──┬──┘          │
         │Yes          v
         │      ┌──────────────────┐
         v      │  Execute Tasks   │
┌──────────────────┐   │    Locally       │
│  Load SSH        │   └──────────────────┘
│  Profile from    │
│  SQLite          │
└────────┬─────────┘
         │
         v
    ┌────┴────┐
    │Password │ No
    │ Flag?   ├──────────┐
    └────┬────┘          │
         │Yes            v
         v        ┌──────────────────┐
┌──────────────────┐    │  Use SSH Key     │
│  Read Password   │    │  from Profile    │
│  from STDIN      │    └──────┬───────────┘
└────────┬─────────┘            │
         │                       │
         v                       v
    ┌────┴──────────────────┬────┘
    │  Establish SSH        │
    │  Connection           │
    └───────────┬───────────┘
                │
                v
    ┌───────────────────────┐
    │  Execute Remote       │
    │  Command              │
    └───────────┬───────────┘
                │
                v
    ┌───────────────────────┐
    │  Capture Output       │
    │  and Exit Code        │
    └───────────┬───────────┘
                │
                v
    ┌───────────────────────┐
    │  Update Stack State   │
    │  and Audit Log        │
    └───────────────────────┘
```

## ⚙️ Advanced Configuration

### Stack State Management

Stack state is persisted across executions:

```bash
# Initial run creates stack
sloth-runner run my-deployment --file deploy.sloth --ssh prod 'echo "v1.0"'

# Subsequent runs use existing stack
sloth-runner run my-deployment --file deploy.sloth --ssh prod 'echo "v1.1"'

# View stack state
sloth-runner stack show my-deployment
```

### Parallel Execution

Execute on multiple hosts simultaneously:

```bash
#!/bin/bash
# parallel-deploy.sh

HOSTS=("web1" "web2" "web3")
PIDS=()

for host in "${HOSTS[@]}"; do
    sloth-runner run "deploy-$host" \
      --file deploy.sloth \
      --ssh "$host" \
      'deployment-script.sh' &
    PIDS+=($!)
done

# Wait for all to complete
for pid in "${PIDS[@]}"; do
    wait "$pid"
    echo "Process $pid completed with status $?"
done
```

### Error Handling

```bash
#!/bin/bash
# robust-execution.sh

set -euo pipefail
trap 'echo "Error on line $LINENO"' ERR

execute_with_retry() {
    local max_attempts=3
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        echo "Attempt $attempt of $max_attempts"

        if sloth-runner run "task-$attempt" \
             --file task.sloth \
             --ssh server \
             "$1"; then
            echo "Success on attempt $attempt"
            return 0
        fi

        echo "Failed on attempt $attempt"
        ((attempt++))
        sleep 5
    done

    echo "All attempts failed"
    return 1
}

execute_with_retry 'critical-command.sh'
```

## 🚫 Security Anti-Patterns to Avoid

### ❌ NEVER: Password in Command Line

```bash
# WRONG - Password visible in process list
sloth-runner run task --ssh server --password "visible_password" 'command'

# WRONG - Password in command history
sloth-runner run task --ssh server --ssh-password "my_password" 'command'
```

### ❌ NEVER: Password in Environment Variable Without Care

```bash
# WRONG - Visible in /proc/*/environ
export PASSWORD="my_password"
sloth-runner run task --ssh server 'command'  # Even if tool reads from env
```

### ❌ NEVER: Unencrypted Password Files

```bash
# WRONG - Password file without permissions
echo "password" > pass.txt  # Default 644 permissions

# WRONG - Password file in shared location
echo -n "password" > /tmp/shared_password.txt
```

### ❌ NEVER: Password with Newline

```bash
# WRONG - Contains newline character
echo "password" > pass.txt

# WRONG - Here-string adds newline
cat <<< "password" > pass.txt
```

### ❌ NEVER: Leaving Password Files

```bash
# WRONG - Password file remains after use
sloth-runner run task --ssh server --ssh-password-stdin - < pass.txt
# Forgot to delete pass.txt
```

## 📊 Exit Codes

| Code | Meaning | Example Scenario |
|------|---------|------------------|
| 0 | Success | All tasks completed successfully |
| 1 | General failure | Unspecified error |
| 2 | Invalid arguments | Missing required flags |
| 3 | File not found | Sloth file doesn't exist |
| 4 | SSH connection failed | Network unreachable |
| 5 | Authentication failed | Wrong password or key |
| 6 | Remote command failed | Non-zero exit on remote |
| 7 | Stack error | Stack operation failed |
| 8 | Permission denied | Insufficient privileges |
| 9 | Timeout | Operation timed out |
| 10 | Interrupted | User cancelled (Ctrl+C) |

## 🔍 Debugging

### Enable Debug Output

```bash
# Set debug environment variable
export SLOTH_RUNNER_DEBUG=true
export SLOTH_RUNNER_LOG_LEVEL=trace

# Run with verbose output
sloth-runner run task \
  --file task.sloth \
  --ssh server \
  --debug \
  --verbose \
  'command' 2>&1 | tee debug.log
```

### Trace SSH Connection

```bash
# Add SSH debug flags
export SSH_DEBUG="-vvv"

sloth-runner run task \
  --file task.sloth \
  --ssh server \
  --ssh-debug \
  'command'
```

### Validate Without Execution

```bash
# Dry run - validate only
sloth-runner run task \
  --file task.sloth \
  --ssh server \
  --dry-run \
  'command'
```

## 🎓 Best Practices Summary

1. **Always quote remote commands** to prevent shell interpretation issues
2. **Use timestamps in stack names** for uniqueness
3. **Store passwords securely** and destroy immediately after use
4. **Prefer key-based authentication** over passwords
5. **Set restrictive permissions** on all sensitive files
6. **Validate password files** don't contain newlines
7. **Use process substitution** for dynamic passwords
8. **Implement retry logic** for network operations
9. **Log actions** but never log passwords
10. **Test in staging** before production execution

## 🔗 Related Documentation

- [SSH Profile Management](./ssh-management.md) - Setting up SSH profiles
- [Security Best Practices](./security.md) - Comprehensive security guide
- [Sloth File Syntax](./sloth-syntax.md) - Writing task definitions
- [Stack Management](./stack-management.md) - Understanding execution stacks

---

**Final Security Reminder:** The `--ssh-password-stdin -` pattern with the dash (`-`) is the ONLY secure way to pass passwords to the command. Any deviation from this exact syntax may compromise security.