# Terraform Examples

This directory contains examples demonstrating Terraform integration with Sloth Runner.

## 📁 Available Examples

### `deploy_git_terraform.sloth`
Complete GitOps workflow that clones a Git repository and deploys infrastructure using Terraform.

**Features:**
- Git repository cloning
- Terraform initialization and planning
- Configuration loading from `values.yaml`
- Comprehensive error handling
- Automatic cleanup on failure

**Usage:**
```bash
sloth-runner run -f examples/terraform/deploy_git_terraform.sloth -v examples/terraform/values.yaml deploy_git_terraform
```

**What it does:**
1. 📡 Clones the specified Git repository
2. 🔧 Initializes Terraform in the cloned directory
3. 📋 Loads configuration from `values.yaml`
4. 🔍 Runs `terraform plan` to preview changes
5. 🚀 Optionally runs `terraform apply` to deploy infrastructure

## 📄 Configuration Files

### `values.yaml`
External configuration file containing Terraform variables and settings.

**Example structure:**
```yaml
terraform:
  do_token: "your-digitalocean-token-here"
  droplet_name: "sloth-runner-demo"
  droplet_region: "nyc3"
  droplet_size: "s-1vcpu-1gb"
  environment: "dev"
  project_name: "sloth-demo"
  enable_backups: false
  droplet_tags:
    - "sloth-runner"
    - "demo"
    - "terraform"
```

## 🎯 Use Cases

- **GitOps Workflows**: Automated infrastructure deployment from Git repositories
- **CI/CD Integration**: Terraform deployments in continuous integration pipelines
- **Infrastructure as Code**: Declarative infrastructure management
- **Multi-environment Deployments**: Using different values files for different environments

## ⚙️ Prerequisites

- Terraform CLI installed
- Git installed
- Appropriate cloud provider credentials configured
- Sloth Runner compiled and available

## 🚀 Getting Started

1. **Edit configuration:**
   ```bash
   vim examples/terraform/values.yaml
   ```

2. **Run the example:**
   ```bash
   sloth-runner run -f examples/terraform/deploy_git_terraform.sloth -v examples/terraform/values.yaml deploy_git_terraform
   ```

3. **Monitor execution:**
   The workflow will show detailed logs for each step including Git clone, Terraform init, plan, and apply operations.

## 🔧 Customization

You can customize the examples by:
- Modifying the Git repository URL
- Changing Terraform variables in `values.yaml`
- Adding additional Terraform operations
- Implementing custom validation logic
- Adding notification hooks for success/failure scenarios

## 📚 Related Documentation

- [Terraform Module Documentation](../../docs/modules/terraform.md)
- [Values Configuration Guide](../../docs/configuration/values.md)
- [GitOps Workflows](../../docs/workflows/gitops.md)