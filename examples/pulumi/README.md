# Pulumi Examples

This directory contains examples demonstrating Pulumi integration with Sloth Runner.

## 📁 Available Examples

### `deploy_git_pulumi.sloth`
Complete GitOps workflow that clones a Git repository and deploys infrastructure using Pulumi.

**Features:**
- Git repository cloning
- Pulumi login and stack management
- Configuration loading from `values_pulumi.yaml`
- Preview with change detection
- Conditional deployment based on changes
- Environment variables support

**Usage:**
```bash
sloth-runner run -f examples/pulumi/deploy_git_pulumi.sloth -v examples/pulumi/values_pulumi.yaml deploy_git_pulumi
```

### `pulumi_config_example.sloth`
Step-by-step Pulumi workflow demonstrating configuration management and preview functionality.

**Features:**
- Repository cloning from specific Go project
- Local Pulumi backend configuration
- Comprehensive configuration setup
- Preview execution with environment variables
- Detailed logging and error handling

**Usage:**
```bash
sloth-runner run -f examples/pulumi/pulumi_config_example.sloth pulumi_complete_example
```

**What it does:**
1. 📡 Clones `https://github.com/chalkan3/go-do-droplet`
2. 🔐 Pulumi login to local backend (`file://.`)
3. 📋 Creates/selects stack `dev`
4. ⚙️ Sets configuration values:
   - `dropletName: sloth-runner`
   - `region: nyc3`
   - `size: s-1vcpu-1gb`
   - `image: ubuntu-22-04-x64`
   - `environment: dev`
   - `project: main`
5. 🔍 Runs `pulumi preview` with environment variables

## 📄 Configuration Files

### `values_pulumi.yaml`
External configuration file containing Pulumi settings.

**Example structure:**
```yaml
pulumi:
  dropletName: "sloth-runner"
  region: "nyc3"
  size: "s-1vcpu-1gb"
  image: "ubuntu-22-04-x64"
  environment: "dev"
  project: "main"

git:
  repository: "https://github.com/chalkan3/go-do-droplet"
  branch: "main"
```

## 🌍 Environment Variables Support

All Pulumi operations support environment variables:

```lua
local envs = {
    PULUMI_CONFIG_PASSPHRASE = "",
    DIGITALOCEAN_TOKEN = "dop_v1_your_token",
    TF_LOG = "INFO"
}

client:preview({ envs = envs })
client:up({ auto_approve = true, envs = envs })
client:destroy({ auto_approve = true, envs = envs })
```

## 🎯 Use Cases

- **GitOps Workflows**: Automated infrastructure deployment from Git repositories
- **Multi-stack Management**: Deploy to different environments (dev, staging, prod)
- **Infrastructure as Code**: Modern infrastructure management with Pulumi
- **Cloud Provider Integration**: Secure authentication with environment variables
- **Preview-driven Deployments**: Review changes before applying

## ⚙️ Prerequisites

- Pulumi CLI installed
- Git installed
- Cloud provider credentials (if deploying real infrastructure)
- Sloth Runner compiled and available

## 🚀 Getting Started

### Quick Start (Config Example)
```bash
# Run the step-by-step configuration example
sloth-runner run -f examples/pulumi/pulumi_config_example.sloth pulumi_complete_example
```

### Full GitOps Workflow
```bash
# Edit configuration if needed
vim examples/pulumi/values_pulumi.yaml

# Run the complete GitOps workflow
sloth-runner run -f examples/pulumi/deploy_git_pulumi.sloth -v examples/pulumi/values_pulumi.yaml deploy_git_pulumi
```

## 🔧 Customization

### Adding Environment Variables
```lua
local envs = {
    PULUMI_CONFIG_PASSPHRASE = "",
    DIGITALOCEAN_TOKEN = os.getenv("DO_TOKEN") or "",
    AWS_ACCESS_KEY_ID = values.secrets.aws_key,
    DEBUG = "1"
}
```

### Using Different Backends
```lua
-- Local backend
client = pulumi.login("file://.", { login_local = true })

-- Pulumi Cloud
client = pulumi.login("urllogin", { login_local = false })

-- Custom backend
client = pulumi.login("s3://my-bucket", { login_local = false })
```

### Multiple Stacks
```lua
-- Development
client:stack("dev", { create = true })

-- Production
client:stack("production", { create = true })
```

## 🛡️ Security Best Practices

1. **Never hardcode secrets:**
   ```lua
   -- ❌ Bad
   DIGITALOCEAN_TOKEN = "dop_v1_actual_token"
   
   -- ✅ Good
   DIGITALOCEAN_TOKEN = os.getenv("DO_TOKEN") or values.secrets.do_token
   ```

2. **Use values.yaml for configuration:**
   ```yaml
   secrets:
     digitalocean_token: "${DO_TOKEN}"
   ```

3. **Environment-specific configurations:**
   ```bash
   sloth-runner run -f example.sloth -v values-prod.yaml workflow
   ```

## 📚 Related Documentation

- [Pulumi Module Documentation](../../docs/modules/pulumi.md)
- [Environment Variables Guide](../../docs/pulumi-environment-variables.md)
- [Values Configuration Guide](../../docs/configuration/values.md)
- [GitOps Workflows](../../docs/workflows/gitops.md)