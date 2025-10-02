# Incus Module

O módulo Incus fornece integração completa com o [Incus](https://github.com/lxc/incus), permitindo gerenciar containers e máquinas virtuais de forma programática.

## Índice

- [Instalação](#instalação)
- [Conceitos Básicos](#conceitos-básicos)
- [API Reference](#api-reference)
  - [Instâncias](#instâncias)
  - [Imagens](#imagens)
  - [Redes](#redes)
  - [Perfis](#perfis)
  - [Storage](#storage)
  - [Snapshots](#snapshots)
- [Exemplos Práticos](#exemplos-práticos)
- [Integração com delegate_to](#integração-com-delegate_to)

## Instalação

O módulo Incus está disponível globalmente em todos os scripts do Sloth Runner:

```lua
-- Não é necessário require, o módulo está global
incus.instance({
    name = "mycontainer",
    image = "ubuntu:22.04"
}):create():start()
```

## Conceitos Básicos

O módulo Incus utiliza uma **API fluente** que permite encadear operações de forma intuitiva e legível.

### Fluent API

```lua
-- Criar, iniciar e aguardar uma instância
incus.instance({
    name = "web-server",
    image = "ubuntu:22.04"
}):create()
  :start()
  :wait_running()

-- Configurar rede
incus.network({
    name = "br-dmz",
    type = "bridge"
}):set_config({
    ["ipv4.address"] = "10.0.0.1/24",
    ["ipv4.nat"] = "true"
}):create()
```

### Execução Remota

Todas as operações suportam execução em agentes remotos via `:delegate_to()`:

```lua
-- Criar container em host remoto
incus.instance({
    name = "db-server",
    image = "ubuntu:22.04"
}):delegate_to("db-host-01")
  :create()
  :start()
```

## API Reference

### Instâncias

Gerenciamento completo de containers e VMs.

#### incus.instance(options)

Cria um builder de instância.

**Parâmetros:**

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `name` | string | ✅ | Nome da instância |
| `image` | string | ❌ | Imagem a ser utilizada |
| `type` | string | ❌ | Tipo: `container` ou `virtual-machine` |
| `profiles` | table | ❌ | Lista de perfis |
| `config` | table | ❌ | Configurações da instância |
| `devices` | table | ❌ | Dispositivos da instância |

**Retorna:** `IncusInstance`

**Métodos da Instância:**

##### :create()

Cria a instância.

```lua
incus.instance({
    name = "web01",
    image = "ubuntu:22.04"
}):create()
```

##### :start()

Inicia a instância.

```lua
instance:start()
```

##### :stop([force])

Para a instância.

```lua
instance:stop()          -- Para gracefully
instance:stop(true)      -- Força a parada
```

##### :restart()

Reinicia a instância.

```lua
instance:restart()
```

##### :delete()

Deleta a instância.

```lua
instance:delete()
```

##### :wait_running([timeout])

Aguarda a instância estar em execução.

```lua
instance:wait_running()       -- Timeout padrão
instance:wait_running(120)    -- 120 segundos
```

##### :exec(command, [options])

Executa comando na instância.

```lua
instance:exec("apt update && apt upgrade -y")

-- Com opções
instance:exec("bash -c 'echo hello'", {
    user = "ubuntu",
    cwd = "/tmp",
    env = {
        ["PATH"] = "/usr/local/bin:/usr/bin:/bin"
    }
})
```

##### :file_push(local_path, remote_path)

Envia arquivo para a instância.

```lua
instance:file_push("/local/config.yaml", "/etc/app/config.yaml")
```

##### :file_pull(remote_path, local_path)

Baixa arquivo da instância.

```lua
instance:file_pull("/var/log/app.log", "./logs/app.log")
```

##### :set_config(config)

Define configurações da instância.

```lua
instance:set_config({
    ["limits.cpu"] = "2",
    ["limits.memory"] = "2GB"
})
```

##### :add_device(name, device_config)

Adiciona dispositivo à instância.

```lua
instance:add_device("eth0", {
    type = "nic",
    nictype = "bridged",
    parent = "br0"
})
```

##### :remove_device(name)

Remove dispositivo da instância.

```lua
instance:remove_device("eth0")
```

##### :delegate_to(agent)

Define em qual agente executar as operações.

```lua
instance:delegate_to("prod-host-01")
```

### Imagens

Gerenciamento de imagens de containers/VMs.

#### incus.image(options)

Cria um builder de imagem.

**Parâmetros:**

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `alias` | string | ✅ | Alias da imagem |
| `source` | string | ❌ | Fingerprint ou URL da imagem |
| `server` | string | ❌ | Servidor de imagens (padrão: `images:`) |

**Métodos:**

##### :copy()

Copia a imagem.

```lua
incus.image({
    alias = "ubuntu-custom",
    source = "ubuntu:22.04"
}):copy()
```

##### :delete()

Deleta a imagem.

```lua
incus.image({alias = "old-image"}):delete()
```

##### :export(path)

Exporta a imagem.

```lua
incus.image({alias = "my-image"}):export("/tmp/image.tar.gz")
```

### Redes

Gerenciamento de redes virtuais.

#### incus.network(options)

Cria um builder de rede.

**Parâmetros:**

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `name` | string | ✅ | Nome da rede |
| `type` | string | ❌ | Tipo de rede (bridge, macvlan, etc.) |
| `config` | table | ❌ | Configurações da rede |

**Métodos:**

##### :create()

Cria a rede.

```lua
incus.network({
    name = "br-dmz",
    type = "bridge"
}):create()
```

##### :delete()

Deleta a rede.

```lua
incus.network({name = "br-dmz"}):delete()
```

##### :attach(instance)

Anexa a rede a uma instância.

```lua
incus.network({name = "br-dmz"}):attach("web01")
```

##### :detach(instance)

Desanexa a rede de uma instância.

```lua
incus.network({name = "br-dmz"}):detach("web01")
```

##### :set_config(config)

Define configurações da rede.

```lua
incus.network({name = "br-dmz"}):set_config({
    ["ipv4.address"] = "10.0.0.1/24",
    ["ipv4.nat"] = "true",
    ["ipv6.address"] = "none"
})
```

### Perfis

Gerenciamento de perfis de configuração.

#### incus.profile(options)

Cria um builder de perfil.

**Parâmetros:**

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `name` | string | ✅ | Nome do perfil |
| `description` | string | ❌ | Descrição do perfil |
| `config` | table | ❌ | Configurações do perfil |
| `devices` | table | ❌ | Dispositivos do perfil |

**Métodos:**

##### :create()

Cria o perfil.

```lua
incus.profile({
    name = "high-performance",
    description = "High performance profile"
}):create()
```

##### :delete()

Deleta o perfil.

```lua
incus.profile({name = "old-profile"}):delete()
```

##### :apply(instance)

Aplica o perfil a uma instância.

```lua
incus.profile({name = "high-performance"}):apply("web01")
```

##### :remove_from(instance)

Remove o perfil de uma instância.

```lua
incus.profile({name = "high-performance"}):remove_from("web01")
```

##### :set_config(config)

Define configurações do perfil.

```lua
incus.profile({name = "high-performance"}):set_config({
    ["limits.cpu"] = "4",
    ["limits.memory"] = "8GB"
})
```

### Storage

Gerenciamento de storage pools.

#### incus.storage(options)

Cria um builder de storage.

**Parâmetros:**

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `name` | string | ✅ | Nome do storage pool |
| `driver` | string | ✅ | Driver (zfs, btrfs, dir, lvm, etc.) |
| `config` | table | ❌ | Configurações do storage |

**Métodos:**

##### :create()

Cria o storage pool.

```lua
incus.storage({
    name = "ssd-pool",
    driver = "zfs"
}):create()
```

##### :delete()

Deleta o storage pool.

```lua
incus.storage({name = "old-pool"}):delete()
```

##### :set_config(config)

Define configurações do storage.

```lua
incus.storage({name = "ssd-pool"}):set_config({
    ["size"] = "100GB",
    ["volume.zfs.use_refquota"] = "true"
})
```

### Snapshots

Gerenciamento de snapshots de instâncias.

#### incus.snapshot(options)

Cria um builder de snapshot.

**Parâmetros:**

| Parâmetro | Tipo | Obrigatório | Descrição |
|-----------|------|-------------|-----------|
| `instance` | string | ✅ | Nome da instância |
| `name` | string | ✅ | Nome do snapshot |
| `stateful` | boolean | ❌ | Snapshot com estado da memória |

**Métodos:**

##### :create()

Cria o snapshot.

```lua
incus.snapshot({
    instance = "web01",
    name = "before-upgrade"
}):create()

-- Snapshot stateful (com memória)
incus.snapshot({
    instance = "web01",
    name = "before-upgrade",
    stateful = true
}):create()
```

##### :restore()

Restaura o snapshot.

```lua
incus.snapshot({
    instance = "web01",
    name = "before-upgrade"
}):restore()
```

##### :delete()

Deleta o snapshot.

```lua
incus.snapshot({
    instance = "web01",
    name = "old-snapshot"
}):delete()
```

### Funções Utilitárias

#### incus.list([filter])

Lista recursos Incus.

```lua
-- Listar todas as instâncias
local instances = incus.list("instances")

-- Listar todas as redes
local networks = incus.list("networks")

-- Listar todos os perfis
local profiles = incus.list("profiles")

-- Listar storage pools
local pools = incus.list("storage-pools")
```

#### incus.info(type, name)

Obtém informações sobre um recurso.

```lua
-- Info de uma instância
local info = incus.info("instance", "web01")
print(info)

-- Info de uma rede
local net_info = incus.info("network", "br-dmz")
print(net_info)
```

#### incus.exec(instance, command, [options])

Executa comando em uma instância (função standalone).

```lua
incus.exec("web01", "systemctl status nginx")

-- Com opções
incus.exec("web01", "whoami", {
    user = "ubuntu"
})
```

#### incus.delete(type, name)

Deleta um recurso (função standalone).

```lua
-- Deletar instância
incus.delete("instance", "old-container")

-- Deletar rede
incus.delete("network", "old-bridge")

-- Deletar perfil
incus.delete("profile", "unused-profile")
```

## 🔥 Exemplo Destacado: Deploy de Web Cluster com Paralelismo

Este é um exemplo completo que demonstra o poder do módulo Incus combinado com goroutines para deploy paralelo de um cluster web completo.

```lua
task({
    name = "deploy-web-cluster",
    delegate_to = "incus-host-01",
    run = function()
        -- 🌐 Criar rede isolada para o cluster
        incus.network({
            name = "web-dmz",
            type = "bridge"
        }):set_config({
            ["ipv4.address"] = "10.10.0.1/24",
            ["ipv4.nat"] = "true",
            ["ipv6.address"] = "none"
        }):create()

        -- ⚡ Deploy paralelo de 3 servidores web
        goroutine.map({"web-01", "web-02", "web-03"}, function(name)
            incus.instance({
                name = name,
                image = "ubuntu:22.04"
            }):create()
              :start()
              :wait_running()
              :exec("apt update && apt install -y nginx")
            
            log.info("✅ " .. name .. " deployed and ready")
        end)
        
        log.info("🎉 Web cluster deployed successfully!")
    end
})
```

**Recursos demonstrados:**

- ✅ Criação de rede isolada com configuração customizada
- ✅ Deploy paralelo usando `goroutine.map()`
- ✅ Método fluente (chaining) para operações sequenciais
- ✅ Execução remota via `:delegate_to()`
- ✅ Instalação de pacotes dentro dos containers

## Exemplos Práticos

### Exemplo 1: Deploy de Web Server

```lua
task({
    name = "deploy-web-server",
    delegate_to = "incus-host-01",
    run = function()
        -- Criar rede
        incus.network({
            name = "web-dmz",
            type = "bridge"
        }):set_config({
            ["ipv4.address"] = "10.10.0.1/24",
            ["ipv4.nat"] = "true",
            ["ipv6.address"] = "none"
        }):create()
        
        -- Criar perfil web
        incus.profile({
            name = "web-server",
            description = "Profile for web servers"
        }):set_config({
            ["limits.cpu"] = "2",
            ["limits.memory"] = "2GB"
        }):create()
        
        -- Criar e configurar instância
        local web = incus.instance({
            name = "nginx-01",
            image = "ubuntu:22.04",
            profiles = {"default", "web-server"}
        })
        
        web:create()
        web:start()
        web:wait_running()
        
        -- Anexar à rede
        incus.network({name = "web-dmz"}):attach("nginx-01")
        
        -- Instalar nginx
        web:exec("apt update")
        web:exec("apt install -y nginx")
        
        -- Copiar configuração
        web:file_push("./nginx.conf", "/etc/nginx/sites-available/default")
        web:exec("systemctl restart nginx")
        
        -- Criar snapshot
        incus.snapshot({
            instance = "nginx-01",
            name = "initial-setup"
        }):create()
        
        log.info("Web server deployed successfully!")
    end
})
```

### Exemplo 2: Cluster de Aplicação

```lua
task({
    name = "deploy-app-cluster",
    delegate_to = "cluster-manager",
    run = function()
        local app_nodes = {"app-01", "app-02", "app-03"}
        
        -- Criar rede do cluster
        incus.network({
            name = "app-cluster",
            type = "bridge"
        }):set_config({
            ["ipv4.address"] = "172.20.0.1/24",
            ["ipv4.nat"] = "false"
        }):create()
        
        -- Criar perfil do app
        incus.profile({
            name = "app-node",
            description = "Application node profile"
        }):set_config({
            ["limits.cpu"] = "4",
            ["limits.memory"] = "4GB",
            ["boot.autostart"] = "true"
        }):create()
        
        -- Deploy de cada nó
        goroutine.map(app_nodes, function(node_name)
            local node = incus.instance({
                name = node_name,
                image = "ubuntu:22.04",
                profiles = {"default", "app-node"}
            })
            
            node:create()
            node:start()
            node:wait_running()
            
            -- Conectar à rede do cluster
            incus.network({name = "app-cluster"}):attach(node_name)
            
            -- Instalar dependências
            node:exec("apt update && apt install -y docker.io")
            
            -- Deploy da aplicação
            node:file_push("./app/docker-compose.yml", "/opt/app/docker-compose.yml")
            node:exec("cd /opt/app && docker-compose up -d")
            
            log.info("Node " .. node_name .. " deployed")
        end)
        
        log.info("Cluster deployed successfully!")
    end
})
```

### Exemplo 3: Backup e Restore

```lua
task({
    name = "backup-critical-instances",
    run = function()
        local instances = {"db-01", "web-01", "cache-01"}
        local timestamp = os.date("%Y%m%d-%H%M%S")
        
        goroutine.map(instances, function(instance)
            -- Criar snapshot
            local snap_name = "backup-" .. timestamp
            
            incus.snapshot({
                instance = instance,
                name = snap_name,
                stateful = true
            }):delegate_to("backup-host"):create()
            
            log.info("Snapshot created: " .. instance .. "/" .. snap_name)
        end)
    end
})

task({
    name = "restore-instance",
    run = function()
        local instance = values.instance or "db-01"
        local snapshot = values.snapshot or "backup-latest"
        
        log.info("Restoring " .. instance .. " from " .. snapshot)
        
        -- Parar instância
        incus.instance({name = instance}):stop(true)
        
        -- Restaurar snapshot
        incus.snapshot({
            instance = instance,
            name = snapshot
        }):restore()
        
        -- Reiniciar
        incus.instance({name = instance}):start():wait_running()
        
        log.info("Restore completed successfully!")
    end
})
```

### Exemplo 4: CI/CD Test Environment

```lua
task({
    name = "create-test-environment",
    run = function()
        local branch = values.branch or "main"
        local test_name = "test-" .. branch:gsub("[^%w]", "-")
        
        -- Criar instância de teste
        local test_env = incus.instance({
            name = test_name,
            image = "ubuntu:22.04"
        })
        
        test_env:delegate_to("ci-runner")
                :create()
                :start()
                :wait_running()
        
        -- Setup do ambiente
        test_env:exec("apt update && apt install -y git build-essential")
        
        -- Clonar código
        test_env:exec("git clone https://github.com/user/repo.git /app")
        test_env:exec("cd /app && git checkout " .. branch)
        
        -- Rodar testes
        local result = test_env:exec("cd /app && make test")
        
        -- Criar snapshot se testes passarem
        if result:find("All tests passed") then
            incus.snapshot({
                instance = test_name,
                name = "tests-passed"
            }):create()
            
            log.info("Tests passed! Snapshot created.")
        else
            log.error("Tests failed!")
        end
        
        -- Cleanup (opcional)
        -- test_env:stop():delete()
    end
})
```

### Exemplo 5: Multi-Host Deployment

```lua
task({
    name = "deploy-distributed-system",
    run = function()
        local hosts = {
            {name = "host-01", role = "database"},
            {name = "host-02", role = "application"},
            {name = "host-03", role = "cache"}
        }
        
        goroutine.map(hosts, function(host)
            local container_name = host.role .. "-server"
            
            incus.instance({
                name = container_name,
                image = "ubuntu:22.04"
            }):delegate_to(host.name)
              :set_config({
                  ["limits.cpu"] = "4",
                  ["limits.memory"] = "8GB"
              }):create()
                :start()
                :wait_running()
            
            -- Configuração específica por role
            if host.role == "database" then
                incus.exec(container_name, "apt install -y postgresql")
            elseif host.role == "application" then
                incus.exec(container_name, "apt install -y nodejs npm")
            elseif host.role == "cache" then
                incus.exec(container_name, "apt install -y redis-server")
            end
            
            log.info("Deployed " .. container_name .. " on " .. host.name)
        end)
    end
})
```

### Exemplo 6: Storage Management

```lua
task({
    name = "setup-storage-infrastructure",
    delegate_to = "storage-host",
    run = function()
        -- Criar storage pools
        incus.storage({
            name = "ssd-pool",
            driver = "zfs"
        }):set_config({
            ["size"] = "500GB",
            ["volume.zfs.use_refquota"] = "true"
        }):create()
        
        incus.storage({
            name = "hdd-pool",
            driver = "btrfs"
        }):set_config({
            ["size"] = "2TB"
        }):create()
        
        -- Criar perfil com storage customizado
        incus.profile({
            name = "fast-storage"
        }):set_config({
            ["root"] = {
                pool = "ssd-pool",
                type = "disk",
                path = "/"
            }
        }):create()
        
        log.info("Storage infrastructure ready!")
    end
})
```

## Integração com delegate_to

O módulo Incus suporta completamente execução remota via `:delegate_to()`:

```lua
-- Exemplo completo com delegate_to
task({
    name = "remote-incus-management",
    run = function()
        -- Criar instância em host remoto
        incus.instance({
            name = "remote-app",
            image = "ubuntu:22.04"
        }):delegate_to(values.target_host)
          :create()
          :start()
          :wait_running()
        
        -- Executar comando na instância remota
        incus.exec("remote-app", "hostname")
               :delegate_to(values.target_host)
    end
})
```

## Melhores Práticas

### 1. Use Perfis para Configurações Comuns

```lua
-- Definir perfil uma vez
incus.profile({
    name = "production",
    config = {
        ["limits.cpu"] = "4",
        ["limits.memory"] = "8GB",
        ["boot.autostart"] = "true"
    }
}):create()

-- Usar em múltiplas instâncias
incus.instance({name = "app-01", profiles = {"default", "production"}})
incus.instance({name = "app-02", profiles = {"default", "production"}})
```

### 2. Sempre Use Snapshots Antes de Mudanças Críticas

```lua
-- Snapshot antes de upgrade
incus.snapshot({
    instance = "prod-db",
    name = "pre-upgrade-" .. os.date("%Y%m%d"),
    stateful = true
}):create()

-- Fazer upgrade
instance:exec("apt upgrade -y")

-- Se der errado, restaurar
-- incus.snapshot({instance = "prod-db", name = "pre-upgrade-..."})restore()
```

### 3. Use Goroutines para Operações Paralelas

```lua
-- Deploy paralelo
goroutine.map({"web-01", "web-02", "web-03"}, function(name)
    incus.instance({name = name, image = "nginx:latest"})
         :create():start():wait_running()
end)
```

### 4. Organize Redes por Função

```lua
-- Rede DMZ
incus.network({name = "dmz", type = "bridge"}):create()

-- Rede interna
incus.network({name = "internal", type = "bridge"}):create()

-- Atribuir instâncias às redes apropriadas
incus.network({name = "dmz"}):attach("web-server")
incus.network({name = "internal"}):attach("database")
```

## Troubleshooting

### Verificar Status de Instâncias

```lua
local info = incus.info("instance", "my-container")
print("Status: " .. info)
```

### Listar Recursos

```lua
local instances = incus.list("instances")
local networks = incus.list("networks")
local profiles = incus.list("profiles")
```

### Logs de Execução

```lua
-- Executar com output detalhado
local result = incus.exec("my-container", "systemctl status nginx")
log.info("Command output: " .. result)
```

## Veja Também

- [Módulo SystemD](systemd.md) - Gerenciar serviços systemd
- [Módulo Pkg](pkg.md) - Gerenciar pacotes
- [Módulo User](user.md) - Gerenciar usuários
- [Módulo SSH](ssh.md) - Operações SSH
- [Goroutines](../core/goroutines.md) - Execução paralela
