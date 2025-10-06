# Secrets Management

Sloth Runner provides a secure secrets management system similar to Pulumi, allowing you to encrypt and store sensitive data for your stacks.

## Overview

- **Strong Encryption**: AES-256-GCM with Argon2 key derivation
- **Per-Stack Secrets**: Each stack has its own encrypted secrets
- **Password-Based**: Secrets are encrypted with a user-provided password
- **Global Access**: Secrets are available in Lua workflows via the `secrets` global table

## Architecture

### Encryption

- **Algorithm**: AES-256-GCM (authenticated encryption)
- **Key Derivation**: Argon2 (memory-hard, resistant to GPU attacks)
  - Memory: 64MB
  - Threads: 4
  - Iterations: 1
  - Key length: 32 bytes (256 bits)
- **Salt**: 32-byte random salt per stack, stored in stack metadata
- **Nonce**: 12-byte random nonce per encryption, prepended to ciphertext

### Storage

- **Database**: SQLite database at `~/.sloth-runner/secrets.db`
- **Table**: `secrets` with columns:
  - `id`: Primary key
  - `stack_id`: Foreign key to stack
  - `name`: Secret name
  - `encrypted_value`: Base64-encoded ciphertext
  - `created_at`: Unix timestamp
  - `updated_at`: Unix timestamp
  - Unique constraint on `(stack_id, name)`

## Commands

### Add Secrets

#### From YAML File

```bash
# Add multiple secrets from YAML file
echo 'mypassword' | sloth-runner secrets add \
  --stack my-app \
  --from-yaml secrets.yaml \
  --password-stdin
```

**secrets.yaml format:**
```yaml
secrets:
  api_key: "sk-abc123456789"
  db_password: "SuperSecret123!"
  aws_access_key: "AKIAIOSFODNN7EXAMPLE"
```

#### From File

```bash
# Add secret from file contents
echo 'mypassword' | sloth-runner secrets add api_key \
  --stack my-app \
  --from-file /path/to/secret.txt \
  --password-stdin
```

#### Interactive

```bash
# Add secret interactively (prompts for value and password)
sloth-runner secrets add api_key --stack my-app
```

### List Secrets

```bash
# List all secret names for a stack (values not shown)
sloth-runner secrets list --stack my-app
```

Output:
```
NAME              CREATED              UPDATED
----              -------              -------
api_key           2025-10-06 15:10:00  2025-10-06 15:10:00
db_password       2025-10-06 15:10:00  2025-10-06 15:10:00
aws_access_key    2025-10-06 15:10:00  2025-10-06 15:10:00
```

### Get Secret

```bash
# Decrypt and display a secret value
echo 'mypassword' | sloth-runner secrets get api_key \
  --stack my-app \
  --password-stdin
```

**⚠️ WARNING**: This displays the decrypted secret value in plain text!

### Remove Secrets

```bash
# Remove a specific secret
sloth-runner secrets remove api_key --stack my-app

# Remove all secrets from a stack
sloth-runner secrets remove --stack my-app --all
```

## Usage in Workflows

### Accessing Secrets

Secrets are automatically loaded and decrypted when you run a workflow with the `--password-stdin` flag:

```bash
# Run workflow with secrets
echo 'mypassword' | sloth-runner run my-app \
  --file deploy.sloth \
  --password-stdin \
  --yes
```

### In Lua Code

Secrets are available via the global `secrets` table:

```lua
task("deploy_app", function()
    -- Access secrets directly
    local api_key = secrets.api_key
    local db_password = secrets.db_password

    -- Use in commands
    exec(string.format("curl -H 'Authorization: Bearer %s' https://api.example.com", api_key))

    -- Configure database
    exec(string.format("psql -c \"ALTER USER app WITH PASSWORD '%s'\"", db_password))
end)
```

### Check if Secrets Exist

```lua
task("check_secrets", function()
    if secrets then
        print("Secrets are loaded")
        if secrets.api_key then
            print("API key is available")
        end
    else
        print("No secrets available")
    end
end)
```

## Security Best Practices

### Password Management

1. **Use Strong Passwords**: Minimum 16 characters, mix of upper/lower/numbers/symbols
2. **Don't Hardcode**: Never put passwords in scripts or version control
3. **Use Environment Variables**: Store passwords in secure environment variables
4. **Rotate Regularly**: Change passwords periodically

Example with environment variable:
```bash
echo "$SECRETS_PASSWORD" | sloth-runner run my-app \
  --file deploy.sloth \
  --password-stdin \
  --yes
```

### Secret Storage

1. **Don't Commit**: Never commit secrets to version control
2. **Use .gitignore**: Add `*.secret`, `secrets.yaml`, etc. to `.gitignore`
3. **Separate Secrets**: Use different passwords for different stacks/environments
4. **Backup Securely**: If backing up secrets database, encrypt the backup

### Workflow Security

1. **Avoid Logging**: Don't print secret values in logs
2. **Use Variables**: Store secrets in variables instead of using directly
3. **Clean Up**: Clear sensitive variables after use
4. **Limit Access**: Only load secrets when needed

Example secure workflow:
```lua
task("deploy_secure", function()
    -- Store in local variable
    local token = secrets.api_token

    -- Use without exposing
    exec(string.format("deploy --token '%s' --quiet", token))

    -- Clear from memory (Lua will garbage collect)
    token = nil
end)
```

## Examples

### Complete Workflow

```lua
-- deploy.sloth
task("setup_infrastructure", function()
    -- Use cloud provider credentials from secrets
    exec(string.format([[
        export AWS_ACCESS_KEY_ID=%s
        export AWS_SECRET_ACCESS_KEY=%s
        terraform apply -auto-approve
    ]], secrets.aws_access_key_id, secrets.aws_secret_key))
end)

task("deploy_application", function()
    :depends_on("setup_infrastructure")

    -- Use database credentials
    exec(string.format([[
        echo "DB_PASSWORD=%s" > .env
        docker-compose up -d
    ]], secrets.db_password))
end)

task("configure_monitoring", function()
    :depends_on("deploy_application")

    -- Use monitoring API key
    exec(string.format([[
        curl -X POST https://monitoring.example.com/api/config \
             -H "X-API-Key: %s" \
             -d '{"service": "my-app"}'
    ]], secrets.monitoring_api_key))
end)
```

### Setup Script

```bash
#!/bin/bash
# setup-secrets.sh

STACK_NAME="production"
PASSWORD_FILE=".secrets-password"

# Add secrets from YAML
cat <<EOF > /tmp/secrets.yaml
secrets:
  aws_access_key_id: "${AWS_ACCESS_KEY_ID}"
  aws_secret_key: "${AWS_SECRET_ACCESS_KEY}"
  db_password: "$(openssl rand -base64 32)"
  monitoring_api_key: "${MONITORING_API_KEY}"
EOF

# Add to sloth-runner
cat "$PASSWORD_FILE" | sloth-runner secrets add \
  --stack "$STACK_NAME" \
  --from-yaml /tmp/secrets.yaml \
  --password-stdin

# Clean up
rm /tmp/secrets.yaml

echo "Secrets added successfully!"
```

### CI/CD Integration

```yaml
# .github/workflows/deploy.yml
name: Deploy with Secrets

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Install sloth-runner
        run: |
          curl -L https://github.com/your-org/sloth-runner/releases/download/v1.0.0/sloth-runner -o sloth-runner
          chmod +x sloth-runner

      - name: Deploy with secrets
        env:
          SECRETS_PASSWORD: ${{ secrets.SLOTH_SECRETS_PASSWORD }}
        run: |
          echo "$SECRETS_PASSWORD" | ./sloth-runner run production \
            --file deploy.sloth \
            --password-stdin \
            --yes
```

## Troubleshooting

### "failed to decrypt secret"

- **Cause**: Wrong password or corrupted data
- **Solution**: Verify password, re-add secret if corrupted

### "stack not found"

- **Cause**: Stack doesn't exist yet
- **Solution**: Run workflow once to create stack, then add secrets

### "secret not found"

- **Cause**: Secret name doesn't exist for this stack
- **Solution**: List secrets to verify name, add if missing

### Secrets not available in workflow

- **Cause**: Forgot `--password-stdin` flag
- **Solution**: Add `--password-stdin` and pipe password via stdin

## Migration

### From Plain Text Files

```bash
# Convert existing plain text secrets to encrypted secrets
for secret_file in secrets/*.txt; do
    name=$(basename "$secret_file" .txt)
    value=$(cat "$secret_file")

    echo 'mypassword' | sloth-runner secrets add "$name" \
      --stack my-app \
      --from-file "$secret_file" \
      --password-stdin
done
```

### From Environment Variables

```bash
# Convert environment variables to secrets
cat <<EOF | sloth-runner secrets add --stack my-app --from-yaml /dev/stdin --password-stdin <<< 'mypassword'
secrets:
  api_key: "${API_KEY}"
  db_password: "${DB_PASSWORD}"
EOF
```

## Technical Details

### File Locations

- Secrets database: `~/.sloth-runner/secrets.db`
- Stack database: `~/.sloth-runner/stacks.db`
- Salt storage: Stack metadata in stacks database

### Encryption Process

1. Generate or retrieve 32-byte salt for stack
2. Derive 32-byte encryption key from password using Argon2
3. Generate 12-byte random nonce
4. Encrypt plaintext with AES-256-GCM
5. Prepend nonce to ciphertext
6. Base64-encode result for storage

### Decryption Process

1. Retrieve encrypted value and salt from database
2. Derive encryption key from password using Argon2
3. Base64-decode ciphertext
4. Extract nonce from first 12 bytes
5. Decrypt with AES-256-GCM
6. Return plaintext

### Memory Security

- Passwords are cleared from memory after use (overwritten with 'x')
- Secrets are only loaded when needed
- Lua garbage collector cleans up secret values

## FAQ

**Q: Can I use different passwords for different stacks?**
A: Yes, each stack has its own salt and secrets. You can use different passwords.

**Q: What happens if I forget my password?**
A: Secrets cannot be recovered without the password. You'll need to re-add them.

**Q: Are secrets encrypted at rest?**
A: Yes, secrets are always encrypted in the database using AES-256-GCM.

**Q: Can I rotate passwords?**
A: Not directly. You need to decrypt all secrets with old password and re-encrypt with new password.

**Q: Are secrets encrypted in transit?**
A: Secrets are decrypted in memory before passing to workflows. Use TLS for network transmission.

**Q: Can I share secrets between stacks?**
A: No, secrets are per-stack. You can add the same secret to multiple stacks.
