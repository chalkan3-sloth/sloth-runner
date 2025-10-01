# Pulumi Module Environment Variables Documentation

The Pulumi module now supports environment variables for all major operations: `preview`, `up`, `destroy`, and `refresh`.

## 🌍 **Environment Variables Support**

All Pulumi operations accept an `envs` parameter to set environment variables:

```lua
local pulumi = require("pulumi")
local client = pulumi.login("file://.", { login_local = true })

-- Environment variables for all operations
local envs = {
    PULUMI_CONFIG_PASSPHRASE = "",
    DIGITALOCEAN_TOKEN = "your-token-here",
    AWS_ACCESS_KEY_ID = "your-key",
    AWS_SECRET_ACCESS_KEY = "your-secret"
}
```

## 📋 **API Reference**

### `client:preview({ envs = {...} })`
Preview infrastructure changes with environment variables.

```lua
local success, output = client:preview({ 
    envs = {
        PULUMI_CONFIG_PASSPHRASE = "",
        DIGITALOCEAN_TOKEN = "dop_v1_your_token"
    }
})
```

### `client:up({ auto_approve = true, envs = {...} })`
Deploy infrastructure with environment variables.

```lua
local success, output = client:up({ 
    auto_approve = true,
    envs = {
        PULUMI_CONFIG_PASSPHRASE = "",
        DIGITALOCEAN_TOKEN = "dop_v1_your_token"
    }
})
```

### `client:destroy({ auto_approve = true, envs = {...} })`
Destroy infrastructure with environment variables.

```lua
local success, output = client:destroy({ 
    auto_approve = true,
    envs = {
        PULUMI_CONFIG_PASSPHRASE = "",
        DIGITALOCEAN_TOKEN = "dop_v1_your_token"
    }
})
```

### `client:refresh({ auto_approve = true, envs = {...} })`
Refresh infrastructure state with environment variables.

```lua
local success, output = client:refresh({ 
    auto_approve = true,
    envs = {
        PULUMI_CONFIG_PASSPHRASE = "",
        DIGITALOCEAN_TOKEN = "dop_v1_your_token"
    }
})
```

## 💡 **Complete Example**

```lua
local deploy_with_envs = task("deploy_with_envs")
    :command(function(this, params)
        local pulumi = require("pulumi")
        
        -- Setup client
        local client = pulumi.login("file://.", { login_local = true })
        client:stack("production", { create = true })
        
        -- Common environment variables
        local envs = {
            PULUMI_CONFIG_PASSPHRASE = "",
            DIGITALOCEAN_TOKEN = "dop_v1_your_production_token",
            TF_LOG = "INFO"
        }
        
        -- Preview with environment
        local preview_ok, preview_output = client:preview({ envs = envs })
        if not preview_ok then
            return false, "Preview failed: " .. preview_output
        end
        
        log.info("Preview successful, proceeding with deployment...")
        
        -- Deploy with environment
        local up_ok, up_output = client:up({ 
            auto_approve = true, 
            envs = envs 
        })
        
        if up_ok then
            log.info("Deployment successful!")
            return true, "Infrastructure deployed"
        else
            log.error("Deployment failed: " .. up_output)
            return false, "Deployment failed"
        end
    end)
    :build()
```

## 🔧 **Common Environment Variables**

### **Pulumi Core:**
- `PULUMI_CONFIG_PASSPHRASE` - Passphrase for state encryption
- `PULUMI_BACKEND_URL` - Backend URL override
- `PULUMI_DEBUG_COMMANDS` - Debug command execution

### **Cloud Providers:**
- `DIGITALOCEAN_TOKEN` - DigitalOcean API token
- `AWS_ACCESS_KEY_ID` - AWS access key
- `AWS_SECRET_ACCESS_KEY` - AWS secret key
- `GOOGLE_CREDENTIALS` - GCP service account JSON
- `ARM_CLIENT_ID` - Azure client ID
- `ARM_CLIENT_SECRET` - Azure client secret

### **Debugging:**
- `TF_LOG` - Terraform log level
- `PULUMI_LOG_LEVEL` - Pulumi log level
- `DEBUG` - General debug flag

## 🛡️ **Security Best Practices**

1. **Never hardcode tokens in code:**
```lua
-- ❌ Bad
local envs = { DIGITALOCEAN_TOKEN = "dop_v1_actual_token" }

-- ✅ Good
local envs = { DIGITALOCEAN_TOKEN = os.getenv("DO_TOKEN") or "" }
```

2. **Use values.yaml for configuration:**
```yaml
# values.yaml
secrets:
  digitalocean_token: "dop_v1_your_token"
  aws_access_key: "AKIA..."
```

```lua
-- Load from values
local envs = {
    DIGITALOCEAN_TOKEN = values.secrets.digitalocean_token,
    AWS_ACCESS_KEY_ID = values.secrets.aws_access_key
}
```

3. **Use CI/CD environment variables:**
```lua
-- In CI/CD pipeline
local envs = {
    DIGITALOCEAN_TOKEN = os.getenv("DO_TOKEN"),
    PULUMI_CONFIG_PASSPHRASE = os.getenv("PULUMI_PASSPHRASE")
}
```

## 📊 **Environment Priority**

The environment variables are merged with the system environment:

1. **System environment** (existing vars)
2. **Provided envs** (override system vars)

```lua
-- System has: DIGITALOCEAN_TOKEN=old_token
local envs = { DIGITALOCEAN_TOKEN = "new_token" }
-- Pulumi will use: new_token
```