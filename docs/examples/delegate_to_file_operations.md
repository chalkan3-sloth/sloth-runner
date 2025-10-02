# Exemplo: Deploy de Configuração com delegate_to

Este exemplo demonstra como usar o `delegate_to` junto com operações de arquivo
para fazer deploy de configurações em servidores remotos.

## Arquitetura

```
Master (local)
├── deploy.sloth          # Script de deploy
├── templates/
│   ├── nginx.conf.tmpl   # Template de configuração
│   └── app.env.tmpl      # Template de variáveis de ambiente
└── files/
    └── app.tar.gz        # Aplicação para deploy

Remote Agent (production-server)
└── /etc/nginx/           # Destino do deploy
    └── sites-available/
```

## Script de Deploy (deploy.sloth)

```lua
local file_ops = require("file_ops")
local exec = require("exec")

-- Task: Deploy de configuração do Nginx
local deploy_nginx_task = task("deploy_nginx_config")
    :description("Deploy Nginx configuration to production server")
    :command(function(this, params)
        log.info("🚀 Deploying Nginx configuration...")
        
        -- 1. Renderizar template do Nginx
        local result, err = file_ops.template(
            "templates/nginx.conf.tmpl",      -- Template no master
            "/tmp/nginx.conf",                -- Destino temporário no agente
            {
                server_name = params.server_name or "example.com",
                port = params.port or "80",
                root_path = "/var/www/html",
                proxy_pass = params.backend_url or "http://localhost:3000"
            }
        )
        
        if not result then
            return false, "Failed to render template: " .. tostring(err)
        end
        
        log.info("✅ Template rendered successfully")
        
        -- 2. Copiar para o local final (requer sudo)
        local copy_result = exec.run("sudo", "cp", "/tmp/nginx.conf", 
                                     "/etc/nginx/sites-available/myapp")
        
        if copy_result.exit_code ~= 0 then
            return false, "Failed to copy config: " .. copy_result.stderr
        end
        
        -- 3. Criar symlink se não existir
        exec.run("sudo", "ln", "-sf", 
                "/etc/nginx/sites-available/myapp",
                "/etc/nginx/sites-enabled/myapp")
        
        -- 4. Testar configuração
        local test_result = exec.run("sudo", "nginx", "-t")
        
        if test_result.exit_code ~= 0 then
            log.error("❌ Nginx config test failed!")
            log.error(test_result.stderr)
            return false, "Nginx configuration is invalid"
        end
        
        log.info("✅ Nginx configuration test passed")
        
        -- 5. Reload Nginx
        local reload_result = exec.run("sudo", "systemctl", "reload", "nginx")
        
        if reload_result.exit_code ~= 0 then
            return false, "Failed to reload Nginx: " .. reload_result.stderr
        end
        
        log.info("✅ Nginx reloaded successfully")
        
        return true
    end)
    :delegate_to("production-server")  -- Executa no servidor remoto
    :timeout("2m")
    :retry(3)
    :build()

-- Task: Deploy de aplicação
local deploy_app_task = task("deploy_application")
    :description("Deploy application files to production server")
    :depends_on({ "deploy_nginx_config" })
    :command(function(this, params)
        log.info("📦 Deploying application...")
        
        -- 1. Extrair aplicação
        local result, err = file_ops.unarchive(
            "files/app.tar.gz",              -- Arquivo no master
            "/var/www/html",                 -- Destino no agente
            { remote_src = false }           -- Arquivo vem do master
        )
        
        if not result then
            return false, "Failed to extract app: " .. tostring(err)
        end
        
        log.info("✅ Application extracted successfully")
        
        -- 2. Configurar permissões
        exec.run("sudo", "chown", "-R", "www-data:www-data", "/var/www/html")
        exec.run("sudo", "chmod", "-R", "755", "/var/www/html")
        
        log.info("✅ Permissions configured")
        
        return true
    end)
    :delegate_to("production-server")
    :timeout("3m")
    :build()

-- Task: Verificar deploy
local verify_deploy_task = task("verify_deployment")
    :description("Verify deployment is working")
    :depends_on({ "deploy_application" })
    :command(function(this, params)
        log.info("🔍 Verifying deployment...")
        
        local http = require("http")
        
        -- Testar endpoint
        local response, err = http.get("http://localhost")
        
        if not response or response.status_code ~= 200 then
            return false, "Health check failed: " .. tostring(err)
        end
        
        log.info("✅ Deployment verified successfully")
        log.info("Status code: " .. response.status_code)
        
        return true
    end)
    :delegate_to("production-server")
    :timeout("1m")
    :build()

-- Workflow de deploy
workflow.define("production_deployment")
    :description("Complete production deployment workflow")
    :version("1.0.0")
    :tasks({
        deploy_nginx_task,
        deploy_app_task,
        verify_deploy_task
    })
    :config({
        timeout = "10m",
        on_error = "rollback"
    })
    :on_complete(function(success, results)
        if success then
            log.info("🎉 Deployment completed successfully!")
        else
            log.error("❌ Deployment failed!")
            -- Aqui você poderia implementar rollback
        end
        return true
    end)
```

## Template de Nginx (templates/nginx.conf.tmpl)

```nginx
server {
    listen {{.port}};
    server_name {{.server_name}};
    
    root {{.root_path}};
    index index.html index.htm;
    
    location / {
        try_files $uri $uri/ @proxy;
    }
    
    location @proxy {
        proxy_pass {{.proxy_pass}};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # Logs
    access_log /var/log/nginx/{{.server_name}}_access.log;
    error_log /var/log/nginx/{{.server_name}}_error.log;
}
```

## Executando o Deploy

```bash
# Deploy para produção
sloth-runner run deploy.sloth

# Deploy com parâmetros customizados
sloth-runner run deploy.sloth --set server_name=myapp.com --set port=8080
```

## Como Funciona

1. **Master**: O `sloth-runner run` é executado no master
2. **Empacotamento**: O workspace inteiro (incluindo templates e arquivos) é empacotado em um tarball
3. **Envio**: O tarball é enviado para o agente `production-server` via gRPC
4. **Extração**: O agente extrai o tarball em um diretório temporário
5. **Mudança de Diretório**: O agente muda para o workspace extraído (`os.Chdir`)
6. **Execução**: As tasks são executadas, e as operações de arquivo (template, copy) funcionam corretamente porque:
   - Os arquivos de origem (`templates/nginx.conf.tmpl`) existem no workspace extraído
   - Os paths relativos funcionam corretamente
7. **Retorno**: O workspace atualizado é empacotado e retornado ao master

## Vantagens

✅ **Centralizado**: Templates e arquivos ficam no master
✅ **Versionado**: Tudo pode ser versionado no Git
✅ **Seguro**: Não é necessário ter arquivos nos agentes
✅ **Consistente**: Mesmos templates para múltiplos ambientes
✅ **Auditável**: Todas as mudanças são rastreáveis

## Boas Práticas

1. **Use paths relativos**: Sempre use paths relativos ao workspace
2. **Template everything**: Use templates para configurações variáveis
3. **Verifique antes de aplicar**: Sempre teste configurações (ex: `nginx -t`)
4. **Implemente rollback**: Em caso de falha, reverta para a versão anterior
5. **Use dependências**: Garanta ordem de execução com `depends_on()`
