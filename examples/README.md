# Sloth Runner Examples

This directory contains comprehensive examples demonstrating various Sloth Runner capabilities organized by technology and use case.

## 📁 **Directory Structure**

```
examples/
├── terraform/          # Terraform infrastructure examples
├── pulumi/             # Pulumi infrastructure examples
├── systemd/            # Systemd service management examples
└── README.md          # This file
```

## 🚀 **Quick Start Examples**

### **Terraform - Infrastructure as Code**
```bash
# GitOps workflow with Terraform
sloth-runner run -f examples/terraform/deploy_git_terraform.sloth -v examples/terraform/values.yaml deploy_git_terraform
```

### **Pulumi - Modern Infrastructure**
```bash
# Complete Pulumi workflow with preview
sloth-runner run -f examples/pulumi/pulumi_config_example.sloth pulumi_complete_example

# GitOps with Pulumi
sloth-runner run -f examples/pulumi/deploy_git_pulumi.sloth -v examples/pulumi/values_pulumi.yaml deploy_git_pulumi
```

### **Systemd - Service Management**
```bash
# Comprehensive service management demo
sudo sloth-runner run -f examples/systemd/systemd_demo.sloth systemd_demo_workflow
```

## 📋 **Examples by Category**

### 🏗️ **Infrastructure as Code**

| Technology | Example | Description |
|------------|---------|-------------|
| **Terraform** | `terraform/deploy_git_terraform.sloth` | Complete GitOps workflow with Terraform |
| **Pulumi** | `pulumi/deploy_git_pulumi.sloth` | GitOps workflow with Pulumi and change detection |
| **Pulumi** | `pulumi/pulumi_config_example.sloth` | Step-by-step Pulumi configuration and preview |

### ⚙️ **Service Management**

| Technology | Example | Description |
|------------|---------|-------------|
| **Systemd** | `systemd/systemd_demo.sloth` | Complete service lifecycle management |

## 🎯 **Use Cases Demonstrated**

### **GitOps Workflows**
- ✅ Git repository cloning
- ✅ Infrastructure deployment automation
- ✅ Configuration management via values files
- ✅ Error handling and rollback scenarios

### **Infrastructure Management**
- ✅ Terraform plan and apply operations
- ✅ Pulumi preview and up with environment variables
- ✅ Multi-stack deployments
- ✅ Backend configuration (local, cloud, custom)

### **Service Operations**
- ✅ Service creation and configuration
- ✅ Lifecycle management (start, stop, restart)
- ✅ Health monitoring and status checks
- ✅ Blue-green deployment patterns

## 🔧 **Common Patterns**

### **Configuration Management**
All examples support external configuration via values files:
```bash
# Terraform
sloth-runner run -f example.sloth -v values.yaml workflow_name

# Pulumi
sloth-runner run -f example.sloth -v values_pulumi.yaml workflow_name
```

### **Environment Variables**
Secure credential management:
```lua
-- In your workflows
local envs = {
    CLOUD_TOKEN = os.getenv("CLOUD_TOKEN"),
    DATABASE_URL = values.secrets.database_url
}
```

### **Error Handling**
Robust error handling with cleanup:
```lua
:on_fail(function(this, params, output)
    log.error("Task failed, cleaning up...")
    this.workdir:cleanup()
end)
```

## ⚙️ **Prerequisites**

### **General Requirements**
- Sloth Runner compiled and available in PATH
- Git installed and configured

### **Technology-Specific Requirements**

#### **Terraform Examples**
- Terraform CLI installed (`terraform --version`)
- Cloud provider credentials configured

#### **Pulumi Examples**
- Pulumi CLI installed (`pulumi version`)
- Cloud provider credentials or tokens

#### **Systemd Examples**
- Linux system with systemd
- Sudo privileges for service management

## 📚 **Documentation**

Each subdirectory contains detailed README files with:
- 📖 **Comprehensive guides** for each example
- ⚙️ **Configuration options** and customization
- 🎯 **Use case scenarios** and best practices
- 🔧 **Troubleshooting** and common issues

### **Quick Links**
- [Terraform Examples](terraform/README.md)
- [Pulumi Examples](pulumi/README.md)  
- [Systemd Examples](systemd/README.md)

## 🛠️ **Development Workflow**

### **Testing Examples**
```bash
# Test a specific example
sloth-runner run -f examples/category/example.sloth workflow_name

# With custom configuration
sloth-runner run -f examples/category/example.sloth -v custom-values.yaml workflow_name
```

### **Creating New Examples**
1. Choose the appropriate category directory
2. Create your `.sloth` file
3. Add corresponding values file if needed
4. Update the category README.md
5. Test thoroughly with different scenarios

## 🎉 **Getting Started**

1. **Choose your use case** from the categories above
2. **Navigate to the relevant directory** (terraform/, pulumi/, systemd/)
3. **Read the specific README** for detailed instructions
4. **Run the example** with the provided commands
5. **Customize** for your specific needs

## 🔗 **Related Resources**

- [Sloth Runner Documentation](../docs/)
- [Module Reference](../docs/modules/)
- [Configuration Guide](../docs/configuration/)
- [Best Practices](../docs/best-practices/)

---

**Start with any example that matches your infrastructure needs and explore the comprehensive capabilities of Sloth Runner!** 🚀