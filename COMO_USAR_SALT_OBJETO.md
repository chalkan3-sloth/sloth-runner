# 🧂 Como Usar o Módulo Salt como Objeto

Baseado na análise do código, o módulo Salt no task-runner **está implementado como um objeto orientado** e pode ser usado de duas formas:

## 📋 **Duas Formas de Usar o Salt**

### 1️⃣ **Tradicional (Funcional)**
```lua
local salt = require("salt")
local result = salt.ping("*")
```

### 2️⃣ **Orientado a Objetos (Recomendado)**
```lua
local Salt = require("salt_object_oriented")
local client = Salt({config})
local result = client:ping("*")
```

## 🔧 **Uso do Salt como Objeto**

### **Criação do Cliente Salt**
```lua
-- Importar o construtor Salt
local Salt = require("salt_object_oriented")

-- Criar instância com configurações personalizadas
local salt_client = Salt({
    master_host = "localhost",
    master_port = 4506,
    timeout = 30,
    retries = 3,
    output_format = "json",
    cache_dir = "/tmp/salt_cache",
    log_level = "info",
    batch_size = 10
})
```

### **Usando Métodos do Objeto**

#### 🔌 **Conectividade**
```lua
-- SINTAXE CORRETA: objeto:método() com ":"
local ping_result = salt_client:ping("*", {timeout = 30})
local version_info = salt_client:version("*")
local test_result = salt_client:test("*", "version")
```

#### 🏗️ **Gerenciamento de Estados**
```lua
-- Aplicar estado específico
local state_result = salt_client:state_apply("web*", "nginx", {
    test = true,
    pillar = {
        nginx = {
            worker_processes = 4,
            worker_connections = 1024
        }
    }
})

-- Executar highstate completo
local highstate = salt_client:state_highstate("*")

-- Mostrar SLS
local sls_info = salt_client:state_show_sls("*", "nginx")
```

#### 📦 **Gerenciamento de Pacotes**
```lua
-- Instalar pacotes
local install_result = salt_client:pkg_install("web*", "nginx")

-- Listar pacotes
local package_list = salt_client:pkg_list("*")

-- Atualizar repositórios
local refresh_result = salt_client:pkg_refresh("*")

-- Atualizar pacotes
local upgrade_result = salt_client:pkg_upgrade("*")
```

#### ⚙️ **Gerenciamento de Serviços**
```lua
-- Controlar serviços
local start_result = salt_client:service_start("web*", "nginx")
local stop_result = salt_client:service_stop("web*", "apache2")
local restart_result = salt_client:service_restart("*", "ssh")

-- Verificar status
local status_result = salt_client:service_status("*", "nginx")

-- Habilitar/desabilitar
local enable_result = salt_client:service_enable("*", "nginx")
local disable_result = salt_client:service_disable("*", "apache2")
```

#### 🌾 **Operações com Grains**
```lua
-- Obter grains específicos
local os_info = salt_client:grains_get("*", "os_family")

-- Obter todos os grains
local all_grains = salt_client:grains_items("*")

-- Definir grains customizados
local set_result = salt_client:grains_set("*", "environment", "production")

-- Adicionar a listas
local append_result = salt_client:grains_append("*", "roles", "webserver")
```

#### 🔑 **Gerenciamento de Chaves**
```lua
-- Listar chaves
local keys = salt_client:key_list("all")

-- Aceitar chaves
local accept_result = salt_client:key_accept("minion-01")

-- Obter fingerprints
local fingerprints = salt_client:key_finger("*")
```

### **Múltiplos Clientes para Diferentes Ambientes**

```lua
-- Cliente para produção
local salt_prod = Salt({
    master_host = "prod-salt-master.company.com",
    timeout = 60,
    log_level = "warning"
})

-- Cliente para desenvolvimento
local salt_dev = Salt({
    master_host = "dev-salt-master.company.com",
    timeout = 30,
    log_level = "debug"
})

-- Usar clientes específicos
local prod_result = salt_prod:ping("prod-*")
local dev_result = salt_dev:ping("dev-*")
```

### **Operações Avançadas**

#### 🚀 **Execução Assíncrona**
```lua
-- Iniciar job assíncrono
local async_job = salt_client:async("*", "cmd", "run", "long-command")
print("Job ID:", async_job.jid)

-- Verificar status
local job_status = salt_client:job_lookup(async_job.jid)
```

#### 📊 **Execução em Lotes**
```lua
-- Executar em 25% dos minions por vez
local batch_result = salt_client:batch("*", "25%", "pkg", "upgrade")
```

#### 🚨 **Monitoramento com Beacons**
```lua
-- Configurar beacon de monitoramento de disco
local beacon_result = salt_client:beacon_add("*", "diskusage", {
    interval = 300,
    threshold = 85
})

-- Listar beacons
local beacon_list = salt_client:beacon_list("*")
```

#### ⏰ **Agendamento de Tarefas**
```lua
-- Agendar backup diário
local schedule_result = salt_client:schedule_add("*", "daily-backup", {
    function = "cmd.run",
    args = ["/usr/local/bin/backup.sh"],
    hours = 2,
    minutes = 0
})
```

## 🎯 **Principais Diferenças**

| Aspecto | Tradicional | Orientado a Objetos |
|---------|-------------|-------------------|
| **Importação** | `require("salt")` | `require("salt_object_oriented")` |
| **Uso** | `salt.método()` | `client:método()` |
| **Configuração** | Global | Por instância |
| **Múltiplos Ambientes** | Difícil | Fácil |
| **Sintaxe** | `módulo.função()` | `objeto:método()` |

## ✅ **Vantagens do Objeto Salt**

- **🎯 Configuração Isolada**: Cada cliente tem suas próprias configurações
- **🏢 Múltiplos Ambientes**: Diferentes clientes para prod/dev/test
- **🔧 Encapsulamento**: Estado e configuração encapsulados
- **🚀 Sintaxe Limpa**: Uso de `:` para métodos de objeto
- **♻️ Reutilização**: Reutilização eficiente de conexões
- **🔒 Flexibilidade**: Configuração flexível por contexto

## 🚀 **Exemplo Completo**

```lua
-- Importar e criar cliente
local Salt = require("salt_object_oriented")
local salt_client = Salt({
    master_host = "salt-master.company.com",
    timeout = 30,
    output_format = "json"
})

-- Verificar conectividade
local ping_result = salt_client:ping("*")
if ping_result.success then
    print("✅ Salt conectado!")
    
    -- Atualizar sistema
    salt_client:pkg_refresh("*")
    salt_client:pkg_upgrade("*")
    
    -- Instalar nginx
    salt_client:pkg_install("web*", "nginx")
    salt_client:state_apply("web*", "nginx.config")
    salt_client:service_enable("web*", "nginx")
    salt_client:service_start("web*", "nginx")
    
    print("✅ Deploy concluído!")
else
    print("❌ Falha na conexão:", ping_result.error)
end
```

O módulo Salt está **perfeitamente implementado como objeto** e oferece uma API limpa e orientada a objetos para gerenciar infraestrutura SaltStack! 🎉