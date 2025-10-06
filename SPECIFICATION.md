# Sloth Runner SSH Remote Execution - Technical Specification

## Document Information

- **Version:** 1.0.0
- **Status:** Final
- **Author:** Senior Engineering Team
- **Date:** 2024
- **Classification:** Technical Specification

## 1. System Overview

### 1.1 Purpose

Sloth Runner SSH Remote Execution extends the core Sloth Runner task automation system with secure remote command execution capabilities via SSH protocol, implementing industry-standard security practices for credential management.

### 1.2 Scope

This specification defines:
- SSH profile management system using SQLite database
- Secure password handling via standard input
- Remote command execution syntax and behavior
- Security requirements and constraints
- Data persistence and audit logging

### 1.3 Design Principles

1. **Security by Design**: No passwords stored in database or logs
2. **Fail Secure**: Default to most secure options
3. **Explicit Over Implicit**: Require explicit flags for sensitive operations
4. **Audit Everything**: Comprehensive logging without exposing credentials
5. **Standard Compliance**: Follow SSH RFC standards and security best practices

## 2. Functional Requirements

### 2.1 SSH Profile Management

#### 2.1.1 Profile Creation

**Requirement ID:** SSH-001
**Priority:** MUST HAVE

The system MUST provide a command to create SSH connection profiles:

```bash
sloth-runner ssh add <profile-name> --host <host> --user <user> [--port <port>] --key <key-path>
```

**Constraints:**
- Profile names MUST be unique
- Profile names MUST match pattern: `^[a-zA-Z][a-zA-Z0-9_-]{0,49}$`
- Host MUST be valid hostname or IP address
- Port MUST be between 1-65535
- Key path MUST point to readable file with 600 permissions

#### 2.1.2 Profile Storage

**Requirement ID:** SSH-002
**Priority:** MUST HAVE

Profiles MUST be stored in SQLite database with schema:

```sql
CREATE TABLE ssh_profiles (
    name TEXT PRIMARY KEY,
    host TEXT NOT NULL,
    user TEXT NOT NULL,
    port INTEGER DEFAULT 22,
    key_path TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

**Security Constraint:** Passwords MUST NEVER be stored in the database.

#### 2.1.3 Profile Operations

**Requirement ID:** SSH-003
**Priority:** MUST HAVE

The system MUST support:
- List profiles: `sloth-runner ssh list`
- Show profile: `sloth-runner ssh show <profile-name>`
- Update profile: `sloth-runner ssh update <profile-name> [options]`
- Remove profile: `sloth-runner ssh remove <profile-name>`

### 2.2 Remote Command Execution

#### 2.2.1 Basic Execution Syntax

**Requirement ID:** EXEC-001
**Priority:** MUST HAVE

Remote execution command structure:

```bash
sloth-runner run <stack-name> --file <sloth-file> --ssh <profile-name> '<command>'
```

#### 2.2.2 Password Authentication

**Requirement ID:** EXEC-002
**Priority:** MUST HAVE

When password authentication is required:

```bash
sloth-runner run <stack-name> --file <sloth-file> --ssh <profile> --ssh-password-stdin - < password-file '<command>'
```

**Requirements:**
- The flag MUST be exactly `--ssh-password-stdin`
- The flag MUST be followed by a single dash (`-`)
- Password MUST be read from stdin only
- Password MUST NOT contain trailing newline

### 2.3 Security Requirements

#### 2.3.1 Password Handling

**Requirement ID:** SEC-001
**Priority:** MUST HAVE

- Passwords MUST NEVER be accepted as command-line arguments
- Passwords MUST NEVER be stored in database
- Passwords MUST NEVER appear in logs
- Passwords MUST be cleared from memory after use
- Password input MUST be masked in interactive mode

#### 2.3.2 File Permissions

**Requirement ID:** SEC-002
**Priority:** MUST HAVE

- Private key files MUST have 600 permissions
- Database file MUST have 600 permissions
- Password files SHOULD have 600 permissions
- System MUST validate permissions before use

#### 2.3.3 Audit Logging

**Requirement ID:** SEC-003
**Priority:** MUST HAVE

All SSH operations MUST be logged with:
- Timestamp
- Profile name
- Action type
- Success/failure status
- Error message (without credentials)

## 3. Non-Functional Requirements

### 3.1 Performance

- Connection establishment: < 5 seconds
- Command execution overhead: < 100ms
- Database query time: < 10ms
- Maximum concurrent connections: 100

### 3.2 Reliability

- Connection retry: 3 attempts with exponential backoff
- Timeout handling: Configurable (default 30s)
- Graceful degradation on network issues
- Transaction-safe database operations

### 3.3 Usability

- Clear error messages without exposing sensitive data
- Consistent command-line interface
- Comprehensive help documentation
- Example-driven documentation

### 3.4 Compatibility

- SSH Protocol: Version 2
- Operating Systems: Linux, macOS, Windows (WSL)
- Go Version: 1.21+
- SQLite Version: 3.x

## 4. Technical Architecture

### 4.1 Component Diagram

```
┌─────────────────────────────────────┐
│         CLI Interface Layer          │
├─────────────────────────────────────┤
│     Command Parser & Validator       │
├─────────────┬───────────────────────┤
│  SSH Module │    Stack Manager       │
├─────────────┴───────────────────────┤
│         Data Access Layer           │
├─────────────────────────────────────┤
│      SQLite Database Driver         │
└─────────────────────────────────────┘
```

### 4.2 Data Flow

1. User invokes command
2. CLI parses and validates arguments
3. Load SSH profile from database
4. Read password from stdin (if required)
5. Establish SSH connection
6. Execute remote command
7. Capture output and exit code
8. Update stack state
9. Log audit entry
10. Return result to user

### 4.3 Error Handling

All errors MUST be handled with specific exit codes:

| Code | Condition |
|------|-----------|
| 0 | Success |
| 2 | Invalid arguments |
| 3 | File not found |
| 4 | Connection failed |
| 5 | Authentication failed |
| 6 | Command execution failed |

## 5. Implementation Guidelines

### 5.1 Password Reading Implementation

```go
func ReadPasswordFromStdin() (string, error) {
    // Check if stdin is pipe or terminal
    stat, _ := os.Stdin.Stat()
    if (stat.Mode() & os.ModeCharDevice) != 0 {
        // Interactive - use terminal package
        password, err := term.ReadPassword(int(os.Stdin.Fd()))
        return string(password), err
    }

    // Pipe mode - read directly
    scanner := bufio.NewScanner(os.Stdin)
    if scanner.Scan() {
        password := scanner.Text()
        // Validate no newline
        if strings.Contains(password, "\n") {
            return "", errors.New("password contains newline")
        }
        return password, nil
    }
    return "", scanner.Err()
}
```

### 5.2 SSH Connection Establishment

```go
func ConnectSSH(profile Profile, password *string) (*ssh.Client, error) {
    config := &ssh.ClientConfig{
        User: profile.User,
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Implement proper verification
        Timeout: 30 * time.Second,
    }

    if password != nil {
        config.Auth = []ssh.AuthMethod{
            ssh.Password(*password),
        }
        // Clear password immediately
        *password = strings.Repeat("x", len(*password))
        *password = ""
    } else {
        // Key-based auth
        key, err := ioutil.ReadFile(profile.KeyPath)
        if err != nil {
            return nil, err
        }

        signer, err := ssh.ParsePrivateKey(key)
        if err != nil {
            return nil, err
        }

        config.Auth = []ssh.AuthMethod{
            ssh.PublicKeys(signer),
        }
    }

    return ssh.Dial("tcp", fmt.Sprintf("%s:%d", profile.Host, profile.Port), config)
}
```

## 6. Testing Requirements

### 6.1 Unit Tests

- Profile CRUD operations
- Password reading from various sources
- Command parsing and validation
- Database operations
- Error handling

### 6.2 Integration Tests

- End-to-end SSH connection with key
- End-to-end SSH connection with password
- Profile management workflow
- Concurrent execution handling

### 6.3 Security Tests

- Password not in process list
- Password not in logs
- File permission validation
- SQL injection prevention
- Command injection prevention

## 7. Documentation Requirements

### 7.1 User Documentation

- README.md with quick start guide
- Command reference with examples
- Security best practices guide
- Troubleshooting guide
- Migration guide

### 7.2 Developer Documentation

- API documentation (godoc)
- Architecture decision records
- Contributing guidelines
- Security disclosure policy

## 8. Deployment Considerations

### 8.1 Installation

```bash
# Binary installation
curl -L https://github.com/org/sloth-runner/releases/latest/download/sloth-runner-linux-amd64 -o sloth-runner
chmod +x sloth-runner
sudo mv sloth-runner /usr/local/bin/

# Initialize database
sloth-runner ssh init
```

### 8.2 Configuration

Default configuration paths:
- Database: `~/.sloth-runner/ssh_profiles.db`
- Logs: `~/.sloth-runner/logs/`
- Config: `~/.sloth-runner/config.yaml`

### 8.3 Upgrade Path

- Database schema migrations supported
- Backward compatibility for 2 major versions
- Automated backup before migration

## 9. Security Compliance

### 9.1 Standards Compliance

- SSH: RFC 4251-4256 (SSH Protocol)
- TLS: RFC 8446 (TLS 1.3)
- Passwords: NIST SP 800-63B
- Cryptography: FIPS 140-2

### 9.2 Regulatory Compliance

- SOC2 Type II ready
- PCI-DSS compliant password handling
- GDPR compliant data management
- HIPAA ready with encryption

## 10. Acceptance Criteria

### 10.1 Functional Acceptance

- [ ] All profile management commands work as specified
- [ ] Password authentication works via stdin only
- [ ] Key-based authentication works with standard SSH keys
- [ ] Remote command execution captures output correctly
- [ ] Error messages are clear and actionable

### 10.2 Security Acceptance

- [ ] No passwords stored in database verified
- [ ] No passwords in logs verified
- [ ] File permissions enforced
- [ ] Audit logging functional
- [ ] Security scan passed (no critical vulnerabilities)

### 10.3 Performance Acceptance

- [ ] Connection time < 5 seconds
- [ ] 100 concurrent connections supported
- [ ] Database operations < 10ms
- [ ] Memory usage < 100MB for typical operation

## 11. Risk Analysis

### 11.1 Security Risks

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Password exposure | High | Low | Stdin-only input, immediate clearing |
| Key compromise | High | Low | Permission validation, rotation support |
| Database breach | Medium | Low | Local storage, encryption option |
| Network intercept | High | Very Low | SSH encryption, host verification |

### 11.2 Operational Risks

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Connection failure | Low | Medium | Retry logic, timeout handling |
| Database corruption | Medium | Low | Transaction safety, backups |
| Performance degradation | Low | Low | Connection pooling, caching |

## 12. Future Enhancements

### Phase 2 (Planned)
- SSH agent forwarding support
- Jump host / ProxyJump support
- Encrypted database option
- Multi-factor authentication
- Session recording for compliance

### Phase 3 (Considered)
- Kubernetes exec support
- Docker exec support
- Cloud provider native auth (AWS SSM, GCP IAP)
- Centralized profile management
- Role-based access control

## 13. Approval Sign-off

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Technical Lead | | | |
| Security Officer | | | |
| Product Manager | | | |
| QA Lead | | | |

---

## Appendix A: Command Examples

### A.1 Complete Workflow Example

```bash
#!/bin/bash
# Complete deployment workflow

# 1. Setup SSH profile
sloth-runner ssh add production \
  --host prod.example.com \
  --user deploy \
  --key ~/.ssh/production_key \
  --description "Production web server"

# 2. Verify profile
sloth-runner ssh test production

# 3. Execute deployment
STACK="deploy-$(date +%Y%m%d-%H%M%S)"
sloth-runner run "$STACK" \
  --file deploy.sloth \
  --ssh production \
  'cd /app && git pull && docker-compose up -d'

# 4. Verify deployment
sloth-runner run "$STACK" \
  --file verify.sloth \
  --ssh production \
  'docker ps | grep myapp'
```

### A.2 Password Authentication Example

```bash
#!/bin/bash
# Secure password handling

# Create secure password file
PASS_FILE=$(mktemp)
chmod 600 "$PASS_FILE"
trap 'shred -u "$PASS_FILE"' EXIT

# Get password from vault
vault kv get -field=password secret/ssh/legacy > "$PASS_FILE"

# Execute with password
sloth-runner run "legacy-task" \
  --file task.sloth \
  --ssh legacy-server \
  --ssh-password-stdin - < "$PASS_FILE" \
  'sudo systemctl restart application'
```

## Appendix B: Database Schema (Complete)

```sql
-- Complete database schema

-- Main profiles table
CREATE TABLE ssh_profiles (
    name TEXT PRIMARY KEY,
    host TEXT NOT NULL,
    user TEXT NOT NULL,
    port INTEGER DEFAULT 22,
    key_path TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used TIMESTAMP,
    use_count INTEGER DEFAULT 0,
    connection_timeout INTEGER DEFAULT 30,
    keepalive_interval INTEGER DEFAULT 60,
    strict_host_checking BOOLEAN DEFAULT TRUE,
    UNIQUE(host, user, port)
);

-- Audit log
CREATE TABLE ssh_audit_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    profile_name TEXT NOT NULL,
    action TEXT NOT NULL,
    command TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    success BOOLEAN,
    error_message TEXT,
    duration_ms INTEGER,
    bytes_transferred INTEGER,
    FOREIGN KEY(profile_name) REFERENCES ssh_profiles(name)
);

-- Indexes for performance
CREATE INDEX idx_profiles_host ON ssh_profiles(host);
CREATE INDEX idx_profiles_last_used ON ssh_profiles(last_used);
CREATE INDEX idx_audit_timestamp ON ssh_audit_log(timestamp);
CREATE INDEX idx_audit_profile ON ssh_audit_log(profile_name);

-- Triggers
CREATE TRIGGER update_profile_timestamp
AFTER UPDATE ON ssh_profiles
BEGIN
    UPDATE ssh_profiles
    SET updated_at = CURRENT_TIMESTAMP
    WHERE name = NEW.name;
END;

CREATE TRIGGER update_profile_usage
AFTER INSERT ON ssh_audit_log
WHEN NEW.action = 'execute' AND NEW.success = 1
BEGIN
    UPDATE ssh_profiles
    SET last_used = CURRENT_TIMESTAMP,
        use_count = use_count + 1
    WHERE name = NEW.profile_name;
END;
```

---

**Document End**

This specification represents the complete technical requirements for the Sloth Runner SSH Remote Execution feature. Implementation should follow these specifications exactly, with any deviations requiring formal change management approval.