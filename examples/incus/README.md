# Incus Examples

This directory contains practical examples demonstrating the Incus module capabilities.

## Prerequisites

1. **Incus installed** on target host(s)
2. **Sloth Runner agent** running on Incus host
3. **Network connectivity** between master and agents

## Quick Start

### Setup Incus Host

```bash
# On your Incus host
curl -fsSL https://pkgs.zabbly.com/key.asc | gpg --show-keys --fingerprint
sh -c 'cat <<EOF > /etc/apt/sources.list.d/zabbly-incus-stable.sources
Enabled: yes
Types: deb
URIs: https://pkgs.zabbly.com/incus/stable
Suites: $(. /etc/os-release && echo ${VERSION_CODENAME})
Components: main
Architectures: $(dpkg --print-architecture)
Signed-By: /usr/share/keyrings/zabbly.asc
EOF'

apt update
apt install -y incus

# Initialize Incus
incus admin init --auto
```

### Start Sloth Runner Agent

```bash
# On Incus host
sloth-runner agent start \
  --name incus-host-01 \
  --master your-master-ip:50053 \
  --daemon
```

## Examples

### 1. Simple Container (`simple.sloth`)

The most basic example - launch a single Ubuntu container.

**What it does:**
- Launches an Ubuntu 22.04 container
- Installs basic packages (curl, wget, vim)
- Shows the container IP address

**Run:**
```bash
sloth-runner run simple-incus-container --file examples/incus/simple.sloth
```

### 2. Complete Infrastructure (`provision.sloth`)

A comprehensive example showing all Incus features.

**What it does:**
- Creates a custom network with bridge and NAT
- Creates a storage pool
- Creates a custom profile with resource limits
- Launches multiple containers in parallel
- Launches a VM with custom configuration
- Configures networking
- Creates snapshots of all instances
- Verifies deployment health

**Components created:**
- Network: `incus-dmz` (10.20.0.1/24)
- Storage: `local-pool`
- Profile: `web-server`
- Containers: `web-01`, `web-02` (with nginx)
- VM: `db-01` (with PostgreSQL)

**Run:**
```bash
sloth-runner run incus-infrastructure --file examples/incus/provision.sloth
```

### 3. Web Cluster with Load Balancer (`web-cluster.sloth`)

Production-ready example of a complete web infrastructure.

**What it does:**
- Creates an isolated network for the cluster
- Creates a dedicated storage pool
- Deploys 3 web servers in parallel (nginx)
- Deploys a PostgreSQL database server
- Deploys and configures HAProxy load balancer
- Automatically discovers and configures backend servers
- Provides access URLs and health checks

**Architecture:**
```
                    Internet
                       |
                   [HAProxy]
                    /  |  \
                   /   |   \
              [web-01][web-02][web-03]
                       |
                  [PostgreSQL]
```

**Run:**
```bash
sloth-runner run web-cluster --file examples/incus/web-cluster.sloth
```

**Access after deployment:**
- Load balancer: `http://<lb-ip>`
- HAProxy stats: `http://<lb-ip>/stats`

### 4. Legacy Web Cluster (`web_cluster.sloth`)

Deploy a complete nginx web server cluster with:
- Custom network (DMZ)
- Optimized profiles
- Parallel deployment
- Automatic snapshots

**Deploy:**
```bash
sloth-runner run -f examples/incus/web_cluster.sloth deploy-web-cluster
```

**Backup:**
```bash
sloth-runner run -f examples/incus/web_cluster.sloth backup-web-cluster
```

**Scale:**
```bash
sloth-runner run -f examples/incus/web_cluster.sloth scale-web-cluster \
  --values count=3
```

**Restore:**
```bash
sloth-runner run -f examples/incus/web_cluster.sloth restore-web-server \
  --values instance=web-01 \
  --values snapshot=backup-20241002-120000
```

**Cleanup:**
```bash
sloth-runner run -f examples/incus/web_cluster.sloth cleanup-web-cluster
```

### 5. CI/CD Pipeline (`ci_cd.sloth`)

Automated testing and deployment pipeline with:
- Isolated test environments
- Parallel test matrix
- Staging deployment
- Health checks
- Rollback support

**Create test environment:**
```bash
sloth-runner run -f examples/incus/ci_cd.sloth create-test-environment \
  --values branch=feature/new-api \
  --values repo=https://github.com/your-org/your-app.git
```

**Run test matrix:**
```bash
sloth-runner run -f examples/incus/ci_cd.sloth parallel-test-matrix
```

**Deploy to staging:**
```bash
sloth-runner run -f examples/incus/ci_cd.sloth deploy-staging
```

**Rollback staging:**
```bash
sloth-runner run -f examples/incus/ci_cd.sloth rollback-staging \
  --values snapshot=pre-deploy-20241002-120000
```

**Keep test environment for debugging:**
```bash
sloth-runner run -f examples/incus/ci_cd.sloth keep-test-env \
  --values branch=feature/new-api
```

**Cleanup test environment:**
```bash
sloth-runner run -f examples/incus/ci_cd.sloth cleanup-test-env \
  --values branch=feature/new-api
```

## Key Features Demonstrated

### Fluent API

```lua
-- Chain operations naturally
incus.instance({name = "myapp", image = "ubuntu:22.04"})
     :create()
     :start()
     :wait_running()
```

### Remote Execution

```lua
-- Execute on remote Incus host
task({
    name = "deploy",
    delegate_to = "incus-host-01",  -- Run on remote agent
    run = function()
        incus.instance({name = "app"})
             :create()
             :start()
    end
})
```

### Parallel Operations

```lua
-- Deploy multiple instances simultaneously
goroutine.map({"web-01", "web-02", "web-03"}, function(name)
    incus.instance({name = name, image = "nginx:latest"})
         :create()
         :start()
         :wait_running()
end)
```

### State Management

```lua
-- Create snapshots before changes
incus.snapshot({
    instance = "production-db",
    name = "pre-upgrade",
    stateful = true  -- Include memory state
}):create()

-- Make changes...

-- Restore if needed
incus.snapshot({
    instance = "production-db",
    name = "pre-upgrade"
}):restore()
```

### Network Configuration

```lua
-- Create and configure networks
incus.network({
    name = "dmz",
    type = "bridge"
}):set_config({
    ["ipv4.address"] = "10.0.0.1/24",
    ["ipv4.nat"] = "true"
}):create()

-- Attach instances
incus.network({name = "dmz"}):attach("web-server")
```

### File Operations

```lua
-- Push configuration files
instance:file_push("./nginx.conf", "/etc/nginx/nginx.conf")

-- Execute commands
instance:exec("systemctl restart nginx")

-- Pull logs
instance:file_pull("/var/log/app.log", "./logs/app.log")
```

## Common Patterns

### Deploy with Validation

```lua
-- Create instance
local app = incus.instance({name = "app", image = "ubuntu:22.04"})
app:create():start():wait_running()

-- Deploy application
app:exec("apt update && apt install -y docker.io")
app:file_push("./docker-compose.yml", "/opt/app/docker-compose.yml")
app:exec("cd /opt/app && docker-compose up -d")

-- Validate
local health = app:exec("curl -sf http://localhost:8080/health")
if health:find("OK") then
    log.info("✅ Application is healthy!")
else
    error("❌ Health check failed!")
end
```

### Blue-Green Deployment

```lua
-- Create new version (green)
incus.instance({name = "app-green", image = "ubuntu:22.04"})
     :create():start():wait_running()

-- Deploy and test green
-- ... deployment logic ...

-- If successful, stop old version (blue)
incus.instance({name = "app-blue"}):stop()

-- Rename green to blue for next deployment
incus.exec("app-green", "echo 'Now serving production'")
```

### Disaster Recovery

```lua
-- Regular backups
local timestamp = os.date("%Y%m%d-%H%M%S")
incus.snapshot({
    instance = "critical-app",
    name = "backup-" .. timestamp,
    stateful = true
}):create()

-- In case of disaster
incus.snapshot({
    instance = "critical-app",
    name = "backup-20241002-120000"  -- Last known good
}):restore()
```

## Tips & Best Practices

1. **Always create snapshots** before major changes
2. **Use profiles** for common configurations
3. **Leverage goroutines** for parallel operations
4. **Set resource limits** to prevent resource exhaustion
5. **Use stateful snapshots** for databases and stateful apps
6. **Organize networks** by function (DMZ, internal, management)
7. **Tag instances** with metadata for better management
8. **Monitor resources** to prevent over-provisioning

## Troubleshooting

### Instance won't start

```lua
-- Check instance status
local info = incus.info("instance", "myapp")
log.info(info)

-- Check logs
local logs = incus.exec("myapp", "dmesg")
log.info(logs)
```

### Network issues

```lua
-- List networks
local networks = incus.list("networks")
log.info(networks)

-- Check network config
local net_info = incus.info("network", "br0")
log.info(net_info)
```

### Cleanup stuck instances

```bash
# Force stop and delete via command line
incus stop --force myapp
incus delete myapp
```

## See Also

- [Incus Documentation](../../../docs/modules/incus.md)
- [Incus Official Docs](https://linuxcontainers.org/incus/docs/main/)
- [Goroutine Module](../../docs/modules/goroutine.md)
