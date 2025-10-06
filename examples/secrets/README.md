# Secrets Management Examples

This directory contains examples for using the secrets management feature in Sloth Runner.

## Quick Start

### 1. Add Secrets

Copy the example secrets file and edit with your values:
```bash
cp secrets.yaml.example secrets.yaml
# Edit secrets.yaml with your actual secret values
```

Add secrets to a stack:
```bash
echo 'mypassword' | sloth-runner secrets add \
  --stack my-app \
  --from-yaml secrets.yaml \
  --password-stdin
```

### 2. List Secrets

```bash
sloth-runner secrets list --stack my-app
```

### 3. Run Workflow with Secrets

```bash
echo 'mypassword' | sloth-runner run my-app \
  --file secrets_example.sloth \
  --password-stdin \
  --yes
```

## Files

- **secrets_example.sloth** - Complete example workflow using secrets
- **secrets.yaml.example** - Template for secrets YAML file

## Security Note

⚠️ **NEVER commit secrets.yaml to version control!**

Add to `.gitignore`:
```
secrets.yaml
*.secret
.env
```
