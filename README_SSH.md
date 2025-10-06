# 🚀 Sloth Runner - SSH Remote Execution

## Overview

Sloth Runner is a powerful task automation tool that combines Lua scripting capabilities with secure SSH remote execution. This document focuses on the SSH management and remote execution features, emphasizing security best practices for credential handling.

## 🔐 Security First Design

### Core Security Principles

1. **No Password Storage**: Passwords are NEVER stored in the database
2. **Key-Based Authentication Preferred**: SSH keys are the recommended authentication method
3. **Secure Password Input**: When passwords are required, they're read from stdin, never from command-line arguments
4. **Encrypted Communication**: All remote connections use SSH protocol encryption

## 📋 Quick Start

### 1. Add SSH Profile

Register a new SSH connection profile with key-based authentication:

```bash
sloth-runner ssh add production-server \
  --host 192.168.1.100 \
  --user ubuntu \
  --port 22 \
  --key ~/.ssh/id_rsa
```

### 2. Execute Remote Commands

#### Using SSH Key (Recommended)
```bash
sloth-runner run my-stack \
  --file deploy.sloth \
  --ssh production-server \
  'docker ps -a'
```

#### Using Password Authentication (When Required)
```bash
# Create secure password file (no trailing newline)
echo -n "SecureP@ssw0rd" > pass.txt
chmod 600 pass.txt

# Execute with password from stdin
sloth-runner run my-stack \
  --file deploy.sloth \
  --ssh production-server \
  --ssh-password-stdin - < pass.txt \
  'ls -la /var/log'

# Clean up immediately
shred -u pass.txt
```

## 🔧 Command Reference

### SSH Profile Management

#### `sloth-runner ssh add`

Add a new SSH connection profile to the local SQLite database.

**Syntax:**
```bash
sloth-runner ssh add <profile-name> \
  --host <hostname-or-ip> \
  --user <username> \
  [--port <port>] \
  --key <path-to-private-key>
```

**Parameters:**
- `<profile-name>`: Unique identifier for this SSH profile
- `--host`: Target hostname or IP address
- `--user`: SSH username
- `--port`: SSH port (default: 22)
- `--key`: Path to private SSH key file

**Example:**
```bash
sloth-runner ssh add staging \
  --host staging.example.com \
  --user deploy \
  --key ~/.ssh/staging_key
```

#### `sloth-runner ssh list`

List all registered SSH profiles.

```bash
sloth-runner ssh list
```

**Output:**
```
NAME            HOST                USER      PORT   KEY PATH
production      192.168.1.100      ubuntu    22     ~/.ssh/id_rsa
staging         staging.example.com deploy    22     ~/.ssh/staging_key
development     localhost          dev       2222   ~/.ssh/dev_key
```

#### `sloth-runner ssh remove`

Remove an SSH profile from the database.

```bash
sloth-runner ssh remove <profile-name>
```

#### `sloth-runner ssh update`

Update an existing SSH profile.

```bash
sloth-runner ssh update <profile-name> \
  [--host <new-host>] \
  [--user <new-user>] \
  [--port <new-port>] \
  [--key <new-key-path>]
```

### Remote Execution

#### `sloth-runner run` with SSH

Execute tasks defined in a Sloth file on a remote host via SSH.

**Syntax with Key Authentication:**
```bash
sloth-runner run <stack-name> \
  --file <sloth-file> \
  --ssh <ssh-profile-name> \
  '<remote-command>'
```

**Syntax with Password Authentication:**
```bash
sloth-runner run <stack-name> \
  --file <sloth-file> \
  --ssh <ssh-profile-name> \
  --ssh-password-stdin - < password-file \
  '<remote-command>'
```

**Parameters:**
- `<stack-name>`: Name of the execution stack for state management
- `--file`: Path to the Sloth Runner configuration file
- `--ssh`: Name of the SSH profile to use
- `--ssh-password-stdin`: Read password from stdin (use `-` to indicate stdin)
- `<remote-command>`: Command to execute on the remote host

## 🛡️ Security Best Practices

### 1. Password File Management

When password authentication is unavoidable:

```bash
# Create password file securely
touch pass.txt
chmod 600 pass.txt
echo -n "YourPassword" > pass.txt

# Use the password
sloth-runner run deployment \
  --file deploy.sloth \
  --ssh production \
  --ssh-password-stdin - < pass.txt \
  'sudo systemctl restart nginx'

# Immediately destroy the password file
shred -vzu pass.txt
```

### 2. Use Environment Variables

For CI/CD pipelines:

```bash
# Store password in environment variable (still visible in process list)
export SSH_PASS="SecurePassword"

# Use with process substitution
sloth-runner run deployment \
  --file deploy.sloth \
  --ssh production \
  --ssh-password-stdin - < <(echo -n "$SSH_PASS") \
  'deployment-script.sh'
```

### 3. SSH Key Management

Preferred method for authentication:

```bash
# Generate dedicated key for automation
ssh-keygen -t ed25519 -f ~/.ssh/sloth_runner_key -C "sloth-runner@automation"

# Add to SSH profile
sloth-runner ssh add production \
  --host prod.example.com \
  --user deploy \
  --key ~/.ssh/sloth_runner_key

# Set proper permissions
chmod 600 ~/.ssh/sloth_runner_key
chmod 644 ~/.ssh/sloth_runner_key.pub
```

## 📁 Data Storage

### SQLite Database Location

SSH profiles are stored in:
- **Linux/macOS**: `~/.sloth-runner/ssh_profiles.db`
- **Windows**: `%USERPROFILE%\.sloth-runner\ssh_profiles.db`

### Database Schema

```sql
CREATE TABLE ssh_profiles (
    name TEXT PRIMARY KEY,
    host TEXT NOT NULL,
    user TEXT NOT NULL,
    port INTEGER DEFAULT 22,
    key_path TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Note: Passwords are NEVER stored in the database
```

## 🔄 Workflow Examples

### Deployment Pipeline

```bash
#!/bin/bash
# deploy.sh - Secure deployment script

STACK="production-deploy-$(date +%Y%m%d-%H%M%S)"
DEPLOY_FILE="deploy.sloth"
SSH_PROFILE="production"

# Step 1: Add/Update SSH profile
sloth-runner ssh add $SSH_PROFILE \
  --host prod.example.com \
  --user deploy \
  --key ~/.ssh/deploy_key

# Step 2: Execute deployment
sloth-runner run $STACK \
  --file $DEPLOY_FILE \
  --ssh $SSH_PROFILE \
  'cd /app && git pull && docker-compose up -d'

# Step 3: Verify deployment
sloth-runner run $STACK \
  --file $DEPLOY_FILE \
  --ssh $SSH_PROFILE \
  'docker ps | grep myapp'
```

### Multi-Host Execution

```bash
#!/bin/bash
# multi-host.sh - Execute on multiple servers

HOSTS=("web1" "web2" "web3")
COMMAND="sudo systemctl restart nginx"

for host in "${HOSTS[@]}"; do
    echo "Executing on $host..."
    sloth-runner run "restart-nginx-$host" \
      --file maintenance.sloth \
      --ssh "$host" \
      "$COMMAND"
done
```

## ⚠️ Common Pitfalls and Solutions

### Issue: Password visible in process list
**Wrong:**
```bash
sloth-runner run stack --ssh profile --password "MyPassword"  # NEVER DO THIS
```

**Correct:**
```bash
sloth-runner run stack --ssh profile --ssh-password-stdin - < pass.txt
```

### Issue: Newline in password file
**Wrong:**
```bash
echo "password" > pass.txt  # Adds newline
```

**Correct:**
```bash
echo -n "password" > pass.txt  # No newline
# or
printf "password" > pass.txt
```

### Issue: Password file persists
**Wrong:**
```bash
cat pass.txt | sloth-runner run stack --ssh profile --ssh-password-stdin -
# pass.txt remains on disk
```

**Correct:**
```bash
sloth-runner run stack --ssh profile --ssh-password-stdin - < pass.txt && shred -u pass.txt
```

## 🚦 Exit Codes

- `0`: Success
- `1`: General error
- `2`: SSH connection failed
- `3`: Authentication failed
- `4`: Remote command execution failed
- `5`: Profile not found
- `6`: Invalid configuration

## 📚 Additional Resources

- [SSH Profile Management Guide](docs/ssh-management.md)
- [Run Command Complete Syntax](docs/run-syntax.md)
- [Security Best Practices](docs/security.md)
- [Troubleshooting Guide](docs/troubleshooting.md)

## 🤝 Contributing

When contributing SSH-related features, ensure:
1. No passwords in code or logs
2. Secure default configurations
3. Clear security warnings in documentation
4. Unit tests for security features

## 📄 License

This project follows security-first principles. See LICENSE for details.

---

**Remember:** Security is not optional. Always use SSH keys when possible, and handle passwords with extreme care when required.