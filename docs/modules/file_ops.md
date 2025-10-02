# File Operations Module

O módulo `file_ops` fornece operações de gerenciamento de arquivos similares aos módulos do Ansible, permitindo manipular arquivos, templates, archives e muito mais.

## Funções Disponíveis

### copy()
Copia um arquivo de origem para destino, preservando ou definindo permissões.

**Sintaxe:**
```lua
result, err = file_ops.copy(src, dst [, options])
```

**Parâmetros:**
- `src` (string): Caminho do arquivo de origem
- `dst` (string): Caminho do arquivo de destino
- `options` (table, opcional): Opções adicionais
  - `mode` (string): Permissões do arquivo (formato octal, ex: "0644")

**Retorno:**
- `result` (table): Informações sobre a operação
  - `changed` (boolean): Se houve mudança
  - `src` (string): Arquivo de origem
  - `dest` (string): Arquivo de destino
  - `size` (number): Tamanho do arquivo copiado
- `err` (string): Mensagem de erro, se houver

**Exemplos:**

#### Exemplo Básico
```lua
task("copy_config", function()
  local result = file_ops.copy(
    "/etc/nginx/nginx.conf",
    "/tmp/nginx.conf.backup"
  )
  
  if result then
    log.info("Config copied successfully")
  end
end)
```

#### Cópia com Permissões Específicas
```lua
task("copy_with_permissions", function()
  local result = file_ops.copy(
    "/app/config.json",
    "/backup/config.json",
    {mode = "0600"}
  )
  
  if result and result.changed then
    log.info(string.format("Copied %d bytes", result.size))
  end
end)
```

#### Cópia com delegate_to
```lua
task("copy_to_remote", function()
  :delegate_to("prod-server")
  
  local result = file_ops.copy(
    "/local/app.conf",
    "/remote/app.conf",
    {mode = "0644"}
  )
  
  if result then
    log.info("Config deployed to prod-server")
  end
end)
```

---

### fetch()
Busca um arquivo do agente remoto para a máquina local.

**Sintaxe:**
```lua
result, err = file_ops.fetch(src, dst)
```

**Parâmetros:**
- `src` (string): Caminho do arquivo remoto
- `dst` (string): Caminho do arquivo local

**Retorno:**
- `result` (table): Informações sobre a operação
  - `changed` (boolean): Se houve mudança
  - `src` (string): Arquivo de origem
  - `dest` (string): Arquivo de destino
  - `size` (number): Tamanho do arquivo

**Exemplos:**

#### Fetch Básico
```lua
task("fetch_logs", function()
  local result = file_ops.fetch(
    "/var/log/app.log",
    "/local/logs/app.log"
  )
  
  if result then
    log.info(string.format("Fetched %d bytes", result.size))
  end
end)
```

#### Fetch de Múltiplos Servidores
```lua
local servers = {"web1", "web2", "web3"}

for _, server in ipairs(servers) do
  task("fetch_from_" .. server, function()
    :delegate_to(server)
    
    local result = file_ops.fetch(
      "/var/log/nginx/access.log",
      string.format("/logs/%s/access.log", server)
    )
    
    if result then
      log.info(string.format("%s: fetched log", server))
    end
  end)
end
```

---

### template()
Renderiza um template com variáveis e salva no destino.

**Sintaxe:**
```lua
result, err = file_ops.template(src, dst, vars)
```

**Parâmetros:**
- `src` (string): Caminho do arquivo template (formato Go template)
- `dst` (string): Caminho do arquivo de saída
- `vars` (table): Variáveis para renderização

**Retorno:**
- `result` (table): Informações sobre a operação
  - `changed` (boolean): Se houve mudança
  - `src` (string): Template de origem
  - `dest` (string): Arquivo de destino

**Exemplos:**

#### Template Simples
```lua
task("configure_app", function()
  local vars = {
    AppName = "MyApp",
    Port = 8080,
    Environment = "production",
    Debug = false
  }
  
  local result = file_ops.template(
    "/templates/config.json.tmpl",
    "/app/config.json",
    vars
  )
  
  if result and result.changed then
    log.info("Configuration updated")
  end
end)
```

#### Template para Nginx
```lua
task("configure_nginx", function()
  :delegate_to("web-server")
  
  local vars = {
    ServerName = "example.com",
    Port = 80,
    RootDir = "/var/www/html",
    AccessLog = "/var/log/nginx/access.log",
    ErrorLog = "/var/log/nginx/error.log"
  }
  
  local result = file_ops.template(
    "/templates/nginx.conf.tmpl",
    "/etc/nginx/sites-available/mysite",
    vars
  )
  
  if result then
    -- Reload nginx
    cmd.run("nginx -t && systemctl reload nginx")
  end
end)
```

#### Template Condicional
```lua
task("configure_database", function()
  local db_vars = {
    Host = "localhost",
    Port = 5432,
    Database = "myapp",
    User = "appuser",
    MaxConnections = 100,
    EnableSSL = os.getenv("ENV") == "production"
  }
  
  file_ops.template(
    "/templates/database.conf.tmpl",
    "/etc/myapp/database.conf",
    db_vars
  )
end)
```

---

### lineinfile()
Garante que uma linha específica existe (ou não existe) em um arquivo.

**Sintaxe:**
```lua
result, err = file_ops.lineinfile(path, line [, options])
```

**Parâmetros:**
- `path` (string): Caminho do arquivo
- `line` (string): Linha a ser inserida/removida
- `options` (table, opcional): Opções
  - `state` (string): "present" ou "absent" (padrão: "present")
  - `regexp` (string): Expressão regular para encontrar linha existente

**Retorno:**
- `result` (table): Informações sobre a operação
  - `changed` (boolean): Se houve mudança
  - `path` (string): Caminho do arquivo

**Exemplos:**

#### Adicionar Linha
```lua
task("add_host_entry", function()
  local result = file_ops.lineinfile(
    "/etc/hosts",
    "192.168.1.100 myserver.local"
  )
  
  if result.changed then
    log.info("Host entry added")
  end
end)
```

#### Substituir com Regexp
```lua
task("update_config", function()
  -- Atualiza a linha que começa com "port="
  local result = file_ops.lineinfile(
    "/app/config.ini",
    "port=8080",
    {regexp = "^port="}
  )
  
  if result.changed then
    log.info("Port updated to 8080")
  end
end)
```

#### Remover Linha
```lua
task("remove_deprecated_config", function()
  local result = file_ops.lineinfile(
    "/etc/app.conf",
    "enable_legacy_mode=true",
    {state = "absent"}
  )
  
  if result.changed then
    log.info("Legacy mode disabled")
  end
end)
```

#### Configuração de SSH
```lua
task("secure_ssh", function()
  :delegate_to("all-servers")
  
  -- Desabilita autenticação por senha
  file_ops.lineinfile(
    "/etc/ssh/sshd_config",
    "PasswordAuthentication no",
    {regexp = "^#?PasswordAuthentication"}
  )
  
  -- Desabilita root login
  file_ops.lineinfile(
    "/etc/ssh/sshd_config",
    "PermitRootLogin no",
    {regexp = "^#?PermitRootLogin"}
  )
  
  -- Restart SSH
  systemd.restart("sshd")
end)
```

---

### blockinfile()
Insere, atualiza ou remove um bloco de linhas em um arquivo.

**Sintaxe:**
```lua
result, err = file_ops.blockinfile(path, block [, options])
```

**Parâmetros:**
- `path` (string): Caminho do arquivo
- `block` (string): Bloco de texto a ser gerenciado
- `options` (table, opcional): Opções
  - `state` (string): "present" ou "absent" (padrão: "present")
  - `marker_begin` (string): Marcador inicial (padrão: "# BEGIN MANAGED BLOCK")
  - `marker_end` (string): Marcador final (padrão: "# END MANAGED BLOCK")

**Retorno:**
- `result` (table): Informações sobre a operação
  - `changed` (boolean): Se houve mudança
  - `path` (string): Caminho do arquivo

**Exemplos:**

#### Adicionar Bloco de Configuração
```lua
task("add_cron_jobs", function()
  local cron_block = [[
0 2 * * * /usr/local/bin/backup.sh
0 3 * * * /usr/local/bin/cleanup.sh
0 4 * * 0 /usr/local/bin/weekly-report.sh]]
  
  local result = file_ops.blockinfile(
    "/etc/cron.d/myapp",
    cron_block
  )
  
  if result.changed then
    log.info("Cron jobs configured")
  end
end)
```

#### Atualizar Bloco Existente
```lua
task("update_firewall_rules", function()
  local rules = [[
-A INPUT -p tcp --dport 80 -j ACCEPT
-A INPUT -p tcp --dport 443 -j ACCEPT
-A INPUT -p tcp --dport 22 -j ACCEPT]]
  
  file_ops.blockinfile(
    "/etc/iptables/rules.v4",
    rules,
    {
      marker_begin = "# BEGIN SLOTH RULES",
      marker_end = "# END SLOTH RULES"
    }
  )
end)
```

#### Remover Bloco
```lua
task("cleanup_old_config", function()
  local result = file_ops.blockinfile(
    "/etc/app.conf",
    "",
    {state = "absent"}
  )
  
  if result.changed then
    log.info("Old configuration block removed")
  end
end)
```

#### Configurar Hosts File
```lua
task("configure_internal_hosts", function()
  local hosts_block = [[
192.168.1.10 db-primary.internal
192.168.1.11 db-replica.internal
192.168.1.20 cache-01.internal
192.168.1.21 cache-02.internal]]
  
  file_ops.blockinfile(
    "/etc/hosts",
    hosts_block,
    {
      marker_begin = "# BEGIN INTERNAL HOSTS",
      marker_end = "# END INTERNAL HOSTS"
    }
  )
end)
```

---

### replace()
Substitui todas as ocorrências de um padrão em um arquivo usando expressões regulares.

**Sintaxe:**
```lua
result, err = file_ops.replace(path, pattern, replacement)
```

**Parâmetros:**
- `path` (string): Caminho do arquivo
- `pattern` (string): Expressão regular para busca
- `replacement` (string): Texto de substituição

**Retorno:**
- `result` (table): Informações sobre a operação
  - `changed` (boolean): Se houve mudança
  - `path` (string): Caminho do arquivo

**Exemplos:**

#### Substituição Simples
```lua
task("update_version", function()
  local result = file_ops.replace(
    "/app/version.txt",
    "version=1\\.0\\.0",
    "version=2.0.0"
  )
  
  if result.changed then
    log.info("Version updated")
  end
end)
```

#### Atualizar Múltiplas Ocorrências
```lua
task("update_api_endpoint", function()
  local result = file_ops.replace(
    "/app/config.json",
    "http://old-api\\.example\\.com",
    "https://new-api.example.com"
  )
  
  if result.changed then
    log.info("API endpoint updated")
  end
end)
```

#### Replace com Captura de Grupos
```lua
task("update_database_config", function()
  -- Substitui host do banco mantendo o resto da string
  file_ops.replace(
    "/etc/myapp/db.conf",
    "host=(\\w+)\\.old\\.domain",
    "host=$1.new.domain"
  )
end)
```

---

### unarchive()
Extrai arquivos compactados (.zip, .tar, .tar.gz, .tgz).

**Sintaxe:**
```lua
result, err = file_ops.unarchive(src, dst)
```

**Parâmetros:**
- `src` (string): Caminho do arquivo compactado
- `dst` (string): Diretório de destino para extração

**Retorno:**
- `result` (table): Informações sobre a operação
  - `changed` (boolean): Se houve mudança
  - `src` (string): Arquivo de origem
  - `dest` (string): Diretório de destino

**Exemplos:**

#### Extrair ZIP
```lua
task("extract_release", function()
  local result = file_ops.unarchive(
    "/tmp/app-v2.0.0.zip",
    "/opt/myapp"
  )
  
  if result then
    log.info("Release extracted successfully")
  end
end)
```

#### Extrair TAR.GZ
```lua
task("deploy_backup", function()
  :delegate_to("backup-server")
  
  file_ops.unarchive(
    "/backups/data-20240101.tar.gz",
    "/var/restore"
  )
  
  log.info("Backup extracted")
end)
```

#### Deploy de Aplicação
```lua
task("deploy_application", function()
  -- Download release
  http.download(
    "https://releases.example.com/app-v3.0.0.tar.gz",
    "/tmp/app-v3.0.0.tar.gz"
  )
  
  -- Extract
  file_ops.unarchive(
    "/tmp/app-v3.0.0.tar.gz",
    "/opt/myapp"
  )
  
  -- Restart service
  systemd.restart("myapp")
  
  log.info("Application deployed: v3.0.0")
end)
```

---

### stat()
Obtém informações detalhadas sobre um arquivo ou diretório.

**Sintaxe:**
```lua
result = file_ops.stat(path)
```

**Parâmetros:**
- `path` (string): Caminho do arquivo ou diretório

**Retorno:**
- `result` (table): Informações do arquivo
  - `exists` (boolean): Se o arquivo existe
  - `path` (string): Caminho do arquivo
  - `size` (number): Tamanho em bytes
  - `mode` (string): Permissões (formato octal)
  - `is_dir` (boolean): Se é um diretório
  - `is_file` (boolean): Se é um arquivo regular
  - `mtime` (number): Timestamp de modificação
  - `checksum` (string): Checksum SHA256
  - `uid` (number): User ID do proprietário (Unix)
  - `gid` (number): Group ID do proprietário (Unix)

**Exemplos:**

#### Verificar Existência
```lua
task("check_config", function()
  local info = file_ops.stat("/etc/myapp/config.json")
  
  if not info.exists then
    log.error("Configuration file missing!")
    return
  end
  
  log.info("Config file exists")
end)
```

#### Verificar Tamanho
```lua
task("check_log_size", function()
  local info = file_ops.stat("/var/log/app.log")
  
  if info.exists and info.size > 100 * 1024 * 1024 then
    log.warn("Log file is over 100MB, rotating...")
    cmd.run("logrotate /etc/logrotate.d/myapp")
  end
end)
```

#### Verificar Permissões
```lua
task("audit_permissions", function()
  local sensitive_files = {
    "/etc/ssl/private/server.key",
    "/etc/myapp/secrets.conf",
    "/root/.ssh/id_rsa"
  }
  
  for _, file in ipairs(sensitive_files) do
    local info = file_ops.stat(file)
    
    if info.exists then
      if info.mode ~= "600" then
        log.error(string.format(
          "Insecure permissions on %s: %s (expected 600)",
          file, info.mode
        ))
      else
        log.info(string.format("%s: OK", file))
      end
    end
  end
end)
```

#### Comparar Checksums
```lua
task("verify_deployment", function()
  local local_file = file_ops.stat("/local/app.jar")
  local remote_file = file_ops.stat("/opt/app/app.jar")
  
  if local_file.checksum == remote_file.checksum then
    log.info("Deployment verified: checksums match")
  else
    log.error("Checksum mismatch! Deployment may be corrupted")
  end
end)
```

---

## Exemplos Avançados

### Pipeline de Deploy Completo
```lua
local config = {
  version = "v2.5.0",
  servers = {"web1", "web2", "web3"},
  app_dir = "/opt/myapp",
  backup_dir = "/var/backups/myapp"
}

-- Backup da versão atual
task("backup_current", function()
  for _, server in ipairs(config.servers) do
    :delegate_to(server)
    
    local timestamp = os.date("%Y%m%d_%H%M%S")
    local backup_name = string.format("backup_%s.tar.gz", timestamp)
    
    cmd.run(string.format(
      "tar czf %s/%s -C %s .",
      config.backup_dir,
      backup_name,
      config.app_dir
    ))
    
    log.info(string.format("%s: backup created", server))
  end
end)

-- Deploy nova versão
task("deploy_new_version", function()
  :depends_on("backup_current")
  
  for _, server in ipairs(config.servers) do
    :delegate_to(server)
    
    -- Download release
    local release_url = string.format(
      "https://releases.example.com/app-%s.tar.gz",
      config.version
    )
    local tmp_file = "/tmp/app-release.tar.gz"
    
    http.download(release_url, tmp_file)
    
    -- Extract
    file_ops.unarchive(tmp_file, config.app_dir)
    
    -- Update configuration
    local vars = {
      ServerID = server,
      Environment = "production",
      Version = config.version
    }
    
    file_ops.template(
      "/templates/app.conf.tmpl",
      config.app_dir .. "/config/app.conf",
      vars
    )
    
    log.info(string.format("%s: deployed %s", server, config.version))
  end
end)

-- Restart services
task("restart_services", function()
  :depends_on("deploy_new_version")
  
  for _, server in ipairs(config.servers) do
    :delegate_to(server)
    
    systemd.restart("myapp")
    
    -- Wait for healthcheck
    local healthy = false
    for i = 1, 10 do
      sleep(2)
      local response = http.get("http://localhost:8080/health")
      if response and response.status == 200 then
        healthy = true
        break
      end
    end
    
    if healthy then
      log.info(string.format("%s: service healthy", server))
    else
      log.error(string.format("%s: healthcheck failed!", server))
    end
  end
end)
```

### Configuração Centralizada
```lua
task("configure_all_servers", function()
  local servers = {
    {name = "web1", role = "web", ip = "192.168.1.10"},
    {name = "web2", role = "web", ip = "192.168.1.11"},
    {name = "db1", role = "database", ip = "192.168.1.20"}
  }
  
  -- Configure /etc/hosts em todos os servidores
  local hosts_block = ""
  for _, srv in ipairs(servers) do
    hosts_block = hosts_block .. string.format(
      "%s %s.internal %s\n",
      srv.ip, srv.name, srv.name
    )
  end
  
  for _, server in ipairs(servers) do
    :delegate_to(server.name)
    
    -- Update hosts
    file_ops.blockinfile(
      "/etc/hosts",
      hosts_block,
      {
        marker_begin = "# BEGIN CLUSTER HOSTS",
        marker_end = "# END CLUSTER HOSTS"
      }
    )
    
    -- Configure based on role
    if server.role == "web" then
      file_ops.template(
        "/templates/nginx.conf.tmpl",
        "/etc/nginx/nginx.conf",
        {ServerName = server.name, Role = server.role}
      )
    elseif server.role == "database" then
      file_ops.template(
        "/templates/postgresql.conf.tmpl",
        "/etc/postgresql/postgresql.conf",
        {ServerName = server.name, Role = server.role}
      )
    end
    
    log.info(string.format("%s: configured", server.name))
  end
end)
```

---

## Integração com Outros Módulos

### Com systemd
```lua
task("update_and_restart", function()
  :delegate_to("app-server")
  
  -- Update configuration
  file_ops.lineinfile(
    "/etc/myapp/app.conf",
    "workers=8",
    {regexp = "^workers="}
  )
  
  -- Restart if changed
  systemd.restart("myapp")
  
  -- Verify
  if systemd.is_active("myapp") then
    log.info("Service restarted successfully")
  end
end)
```

### Com pkg
```lua
task("install_and_configure", function()
  :delegate_to("new-server")
  
  -- Install package
  pkg.install("nginx")
  
  -- Configure
  file_ops.template(
    "/templates/nginx.conf.tmpl",
    "/etc/nginx/nginx.conf",
    {Port = 80, Workers = 4}
  )
  
  -- Enable and start
  systemd.enable("nginx")
  systemd.start("nginx")
end)
```

### Com user
```lua
task("setup_application_user", function()
  -- Create user
  user.create("appuser", {
    home = "/home/appuser",
    shell = "/bin/bash",
    system = true
  })
  
  -- Create app directory
  cmd.run("mkdir -p /opt/myapp")
  cmd.run("chown appuser:appuser /opt/myapp")
  
  -- Copy application files
  file_ops.copy(
    "/dist/app.jar",
    "/opt/myapp/app.jar"
  )
  
  -- Set permissions
  cmd.run("chown appuser:appuser /opt/myapp/app.jar")
  cmd.run("chmod 755 /opt/myapp/app.jar")
end)
```

---

## Melhores Práticas

1. **Sempre verifique o retorno**: Cheque se a operação foi bem-sucedida
2. **Use idempotência**: Os módulos são idempotentes por design
3. **Faça backup**: Sempre faça backup antes de modificações críticas
4. **Use templates**: Para arquivos complexos, prefira templates
5. **Valide com stat()**: Verifique o estado final dos arquivos
6. **Use delegate_to**: Para operações em servidores remotos
7. **Combine com systemd**: Reinicie serviços após mudanças de configuração

---

## Tratamento de Erros

```lua
task("safe_file_operation", function()
  local result, err = file_ops.copy(
    "/source/file.txt",
    "/dest/file.txt"
  )
  
  if not result then
    log.error("Failed to copy file: " .. tostring(err))
    return
  end
  
  if result.changed then
    log.info("File copied successfully")
  else
    log.info("File already up to date")
  end
end)
```
