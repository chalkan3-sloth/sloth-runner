# 🔒 Security

Enterprise-grade security features for production deployments.

## Overview

Sloth Runner provides comprehensive security features:

- 🔐 Secrets management
- 🛡️ Access control
- 📝 Audit logging
- 🔒 Encryption

## Key Features

### Secrets Management
Secure storage and injection of sensitive data.

```lua
local secret = require("secrets")

local deploy_task = task("secure_deploy")
    :command(function()
        local api_key = secret.get("API_KEY")
        -- Use securely
    end)
    :build()
```

### Access Control
Role-based access control (RBAC) for workflows and resources.

### Audit Trail
Complete logging of all actions for compliance.

### Encryption
Data encryption at rest and in transit.

## Best Practices

- ✅ Use secret management for credentials
- ✅ Enable audit logging
- ✅ Implement least privilege access
- ✅ Regular security audits
- ✅ Encrypt sensitive data

## Learn More

- [Enterprise Features](./enterprise-features.md)
- [Best Practices](./best-practices.md)
