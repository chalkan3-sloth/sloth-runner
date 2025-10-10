# User Module 👤

O módulo **user** fornece funcionalidades completas de gerenciamento de usuários e grupos em sistemas Linux/Unix. Ele permite criar, modificar, deletar e consultar usuários e grupos de forma programática.

## 📦 Importação

```lua
```

## 🚀 Funcionalidades Principais

### Gerenciamento de Usuários

#### **user.create(username, options)**
Cria um novo usuário no sistema.

**Parâmetros:**
- `username` (string): Nome do usuário a ser criado
- `options` (table, opcional): Opções de configuração do usuário
  - `password`: Senha do usuário (será configurada automaticamente)
  - `home`: Diretório home do usuário
  - `shell`: Shell padrão do usuário
  - `uid`: UID específico para o usuário
  - `gid`: GID do grupo primário
  - `groups`: Lista de grupos secundários (separados por vírgula)
  - `comment`: Comentário/GECOS do usuário
  - `system`: Marcar como usuário de sistema
  - `create_home`: Criar diretório home
  - `no_create_home`: Não criar diretório home
  - `expiry`: Data de expiração (formato: YYYY-MM-DD)

**Retorna:** `success (boolean), message (string)`

**Exemplo:**

```lua
task("create-user", {
    action = function()
        
        -- Criar usuário simples
        local ok, msg = user.create("john")
        if not ok then
            error("Failed to create user: " .. msg)
        end
        
        -- Criar usuário com opções avançadas e senha
        local ok, msg = user.create("devops", {
            password = "SecureP@ssw0rd!",
            home = "/home/devops",
            shell = "/bin/bash",
            groups = "docker,wheel",
            comment = "DevOps Engineer",
            create_home = true
        })
        
        print("User created successfully!")
    end
})
```

**Exemplo com delegate_to:**

```lua
task("create-remote-user", {
    action = function()
        
        -- Criar usuário em servidor remoto com senha
        delegate_to("production-server", function()
            local ok, msg = user.create("appuser", {
                password = "MyS3cretP@ss",
                shell = "/bin/bash",
                groups = "www-data",
                system = true,
                no_create_home = true
            })
            
            if ok then
                print("User created on remote server")
            end
        end)
    end
})
```

#### **user.delete(username, remove_home)**
Remove um usuário do sistema.

**Parâmetros:**
- `username` (string): Nome do usuário a ser removido
- `remove_home` (boolean, opcional): Remover também o diretório home (padrão: false)

**Retorna:** `success (boolean), message (string)`

**Exemplo:**

```lua
task("cleanup-users", {
    action = function()
        
        -- Deletar usuário mantendo o home
        user.delete("tempuser")
        
        -- Deletar usuário e seu diretório home
        user.delete("olduser", true)
    end
})
```

[... continues with full documentation ...]
