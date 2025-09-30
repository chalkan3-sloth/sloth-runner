# GitOps Example: Deploy Git + Terraform

This comprehensive example demonstrates how to create a complete GitOps workflow using Sloth Runner. It showcases the integration of Git repository management with Terraform infrastructure deployment, all configured through external YAML files.

## 📋 **Overview**

The `deploy_git_terraform.sloth` example implements a production-ready GitOps workflow that:

1. **Clones a Git repository** containing Terraform infrastructure code
2. **Initializes Terraform** automatically in the correct directory
3. **Loads configuration** from external `values.yaml` file
4. **Plans and applies** infrastructure changes with comprehensive error handling
5. **Provides detailed logging** throughout the entire process

## 🗂️ **Files Structure**

```
examples/
├── deploy_git_terraform.sloth    # Main workflow definition
├── values.yaml                   # External configuration
└── README.md                     # Documentation
```

## 🚀 **Quick Start**

### 1. Run the Example

```bash
# Clone the Sloth Runner repository
git clone https://github.com/chalkan3-sloth/sloth-runner.git
cd sloth-runner

# Execute the GitOps workflow
sloth-runner run -f examples/deploy_git_terraform.sloth -v examples/values.yaml deploy_git_terraform
```

### 2. Watch the Execution

```
Loading values from: examples/values.yaml
📡 Starting git.clone()...
✅ git.clone() executed!
🔄 Running terraform init...
✅ Terraform init successful
📋 Loading Terraform configuration from values.yaml
📄 Created tfvars file: terraform.tfvars
📊 Terraform Plan Output:
... (plan details)
🚀 Running terraform apply...
🎉 Infrastructure deployed successfully!
```

## 📝 **Workflow Definition**

### Task 1: Clone Repository

```lua
local clone_repo_task = task("clone_digitalocean_repo")
    :description("Clone Git repository with Terraform infrastructure")
    :workdir("/tmp/digitalocean-droplets")
    :command(function(this, params)
        local workdir_ensured = this.workdir.ensure()
        if not workdir_ensured then 
          return false, "Workdir problem"
        end 

        local git = require("git")
        
        log.info("📡 Starting git.clone()...")
        
        local git_repository = git.clone(
            "https://github.com/chalkan3/terraform-do-droplet",
            this.workdir.get()
        )
        
        log.info("✅ git.clone() executed!")
        log.info("📊 Repo object: " .. tostring(git_repository))
            
        return true, "Git clone successful using modern DSL", { 
            git_module_used = true,
            modern_dsl_used = true,
            repository_url = "https://github.com/chalkan3/terraform-do-droplet",
            clone_destination = this.workdir.get()
        }
    end)
    :timeout("5m")
    :on_success(function(this, params, output)
        log.info("✅ === CLONE SUCCESS ===")
    end)
    :on_fail(function(this, params, output)
        this.workdir:cleanup()
    end)
    :build()
```

### Task 2: Deploy Infrastructure

```lua
local deploy_terraform_task = task("deploy-terraform")
    :description("Deploy infrastructure using Terraform with values.yaml configuration")
    :workdir("/tmp/digitalocean-droplets/environments/dev/")
    :command(function(this, params)
        log.info("🔧 Starting Terraform deployment...")
        
        -- Check if terraform files exist
        local workdir = this.workdir:get()
        local main_tf = io.open(workdir .. "/main.tf", "r")
        local has_terraform_files = main_tf ~= nil
        if main_tf then main_tf:close() end
        
        if not has_terraform_files then
            log.warn("⚠️ No main.tf file found in repository")
            return true, "No Terraform files found - skipping Terraform operations", {
                terraform_used = false,
                skipped = true,
                reason = "no_terraform_files"
            }
        end
        
        local terraform = require("terraform")
        
        -- Create terraform client - this automatically runs terraform init
        log.info("🔄 Running terraform init...")
        local client = terraform.init(workdir)
        
        -- Load configuration from values.yaml with safety checks
        local terraform_config = {}
        if values and values.terraform then
            log.info("📋 Loading Terraform configuration from values.yaml")
            terraform_config = {
                do_token = values.terraform.do_token or "",
                droplet_name = values.terraform.droplet_name or "sloth-runner-demo",
                droplet_region = values.terraform.droplet_region or "nyc3",
                droplet_size = values.terraform.droplet_size or "s-1vcpu-1gb",
                environment = values.terraform.environment or "dev",
                project_name = values.terraform.project_name or "sloth-demo",
                enable_backups = values.terraform.enable_backups or false,
                droplet_tags = values.terraform.droplet_tags or { "sloth-runner", "demo", "terraform" }
            }
        else
            log.warn("⚠️ Values table not found, using default values")
            terraform_config = {
                do_token = "",
                droplet_name = "sloth-runner-demo",
                droplet_region = "nyc3",
                droplet_size = "s-1vcpu-1gb",
                environment = "dev",
                project_name = "sloth-demo",
                enable_backups = false,
                droplet_tags = { "sloth-runner", "demo", "terraform" }
            }
        end
        
        -- Create terraform.tfvars from configuration
        local tfvars = client:create_tfvars("terraform.tfvars", terraform_config)
        log.info("📄 Created tfvars file: " .. tfvars.filename)
        
        -- Execute terraform plan
        local plan_result = client:plan({ var_file = tfvars.filename })
        
        if plan_result.success then
            log.info("📊 Terraform Plan Output:")
            log.print(plan_result.output)
            
            -- Execute terraform apply after successful plan
            log.info("🚀 Running terraform apply...")
            local apply_result = client:apply({
                var_file = tfvars.filename,
                auto_approve = true  -- Auto-approve for automation
            })
            
            if apply_result.success then
                log.info("🎉 Terraform Apply Output:")
                log.print(apply_result.output)
                return true, "Terraform plan and apply successful", {
                    terraform_used = true,
                    init_success = true,
                    plan_success = true,
                    apply_success = true,
                    infrastructure_deployed = true,
                    config_source = "values.yaml"
                }
            else
                log.error("❌ Terraform apply failed: " .. apply_result.error)
                return false, "Terraform apply failed", {
                    terraform_used = true,
                    init_success = true,
                    plan_success = true,
                    apply_success = false,
                    error = apply_result.error
                }
            end
        else
            log.error("❌ Terraform plan failed: " .. plan_result.error)
            return false, "Terraform plan failed", {
                terraform_used = true,
                init_success = true,
                plan_success = false,
                error = plan_result.error
            }
        end
    end)
    :timeout("15m")
    :build()
```

## ⚙️ **Configuration Management**

### values.yaml Structure

```yaml
# Values configuration for deploy_git_terraform.sloth example
terraform:
  # DigitalOcean API Token (set your real token here for actual deployment)
  do_token: "your-digitalocean-token-here"
  
  # Droplet Configuration
  droplet_name: "sloth-runner-demo"
  droplet_region: "nyc3"
  droplet_size: "s-1vcpu-1gb"
  
  # Environment and Project Settings
  environment: "demo"
  project_name: "sloth-demo"
  
  # Backup Configuration
  enable_backups: false
  
  # Tags for the droplet
  droplet_tags:
    - "sloth-runner"
    - "demo"
    - "terraform"
    - "gitops"

# Git Configuration (for future extensions)
git:
  repository: "https://github.com/chalkan3/terraform-do-droplet"
  branch: "main"
  
# Workflow Configuration
workflow:
  timeout: "20m"
  max_parallel_tasks: 1
  environment: "demo"
```

### Loading Configuration in Workflow

```lua
-- Safe configuration loading with fallbacks
local terraform_config = {}

if values and values.terraform then
    log.info("📋 Loading Terraform configuration from values.yaml")
    terraform_config = {
        do_token = values.terraform.do_token or "",
        droplet_name = values.terraform.droplet_name or "sloth-runner-demo",
        droplet_region = values.terraform.droplet_region or "nyc3",
        droplet_size = values.terraform.droplet_size or "s-1vcpu-1gb",
        environment = values.terraform.environment or "dev",
        project_name = values.terraform.project_name or "sloth-demo",
        enable_backups = values.terraform.enable_backups or false,
        droplet_tags = values.terraform.droplet_tags or { "sloth-runner", "demo", "terraform" }
    }
else
    log.warn("⚠️ Values table not found, using default values")
    terraform_config = {
        -- Default configuration fallbacks
    }
end
```

## 🔧 **Key Features Demonstrated**

### 1. **Git Integration**
- Repository cloning with error handling
- Workspace management
- Automatic cleanup on failure

### 2. **Terraform Integration**
- Automatic `terraform init` execution
- Dynamic `terraform.tfvars` generation
- Plan and apply lifecycle management
- Comprehensive error handling

### 3. **External Configuration**
- YAML-based configuration management
- Safe value loading with fallbacks
- Environment-specific settings

### 4. **Error Handling**
- Graceful error handling and reporting
- Automatic cleanup on failure
- Detailed logging throughout execution

### 5. **Modern DSL Features**
- Fluent API task definitions
- Workflow orchestration
- Comprehensive success/failure handlers

## 🎯 **Customization**

### For Your Infrastructure

1. **Update the Git repository URL:**
   ```lua
   local git_repository = git.clone(
       "https://github.com/your-org/your-terraform-repo",
       this.workdir.get()
   )
   ```

2. **Modify the values.yaml:**
   ```yaml
   terraform:
     do_token: "your-actual-token"
     droplet_name: "your-droplet-name"
     # ... other configuration
   ```

3. **Adjust working directories:**
   ```lua
   :workdir("/tmp/your-infrastructure/path/to/terraform/")
   ```

## 🔒 **Security Considerations**

1. **Never commit real API tokens** to version control
2. **Use environment variables** for sensitive values:
   ```yaml
   terraform:
     do_token: "${DO_TOKEN}"
   ```
3. **Implement proper access controls** for your Git repositories
4. **Use separate values.yaml files** for different environments

## 📚 **Learning Outcomes**

After studying this example, you'll understand:

- How to create multi-task GitOps workflows
- How to integrate Git and Terraform modules effectively
- How to use external configuration files with Sloth Runner
- How to implement comprehensive error handling
- How to structure production-ready automation workflows
- How to leverage the Modern DSL for infrastructure automation

## 🔗 **Related Documentation**

- [Git Module Reference](modules/git.md)
- [Terraform Module Reference](modules/terraform.md)
- [Values Configuration Guide](configuration/values.md)
- [Modern DSL Syntax](modern-dsl/syntax.md)
- [Error Handling Patterns](advanced-features/error-handling.md)