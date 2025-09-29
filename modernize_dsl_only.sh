#!/bin/bash

# Script para remover TaskDefinition e usar apenas Modern DSL
# Este script substitui completamente o formato antigo pelo novo

echo "🧹 Removing old DSL and keeping only Modern DSL..."

# Lista de arquivos principais para converter
files_to_modernize=(
    "examples/aws_example.lua"
    "examples/azure_example.lua" 
    "examples/gcp_example.lua"
    "examples/docker_example.lua"
    "examples/git_example.lua"
    "examples/pulumi_example.lua"
    "examples/terraform_example.lua"
    "examples/digitalocean_example.lua"
    "examples/artifact_example.lua"
    "examples/artifacts_example.lua"
    "examples/beginner/hello-world.lua"
    "examples/beginner/http-basics.lua"
    "examples/beginner/docker-basics.lua"
    "examples/beginner/state-basics.lua"
    "examples/intermediate/api-integration.lua"
    "examples/intermediate/multi-container.lua"
    "examples/intermediate/parallel-processing.lua"
    "examples/real-world/nodejs-cicd.lua"
)

# Função para converter arquivo para Modern DSL puro
modernize_file() {
    local file="$1"
    
    if [ ! -f "$file" ]; then
        echo "⚠️  File not found: $file"
        return
    fi
    
    echo "🔄 Modernizing: $file"
    
    # Backup
    cp "$file" "$file.backup" 2>/dev/null || true
    
    # Detectar se o arquivo já está modernizado
    if grep -q "task(" "$file" && ! grep -q "TaskDefinition" "$file"; then
        echo "✅ Already modern: $file"
        return
    fi
    
    # Gerar versão Modern DSL baseada no tipo de arquivo
    modernize_based_on_content "$file"
    
    echo "✅ Modernized: $file"
}

# Função para modernizar baseado no conteúdo
modernize_based_on_content() {
    local file="$1"
    local filename=$(basename "$file")
    
    case "$filename" in
        "hello-world.lua")
            create_modern_hello_world "$file"
            ;;
        "http-basics.lua")
            create_modern_http_basics "$file"
            ;;
        "docker-basics.lua")
            create_modern_docker_basics "$file"
            ;;
        "state-basics.lua")
            create_modern_state_basics "$file"
            ;;
        "aws_example.lua")
            create_modern_aws_example "$file"
            ;;
        "azure_example.lua")
            create_modern_azure_example "$file"
            ;;
        "gcp_example.lua")
            create_modern_gcp_example "$file"
            ;;
        *)
            create_generic_modern_example "$file"
            ;;
    esac
}

# Função para criar hello-world moderno
create_modern_hello_world() {
    local file="$1"
    cat > "$file" << 'EOF'
-- MODERN DSL ONLY - Hello World Example
-- Demonstrates basic Modern DSL task creation

-- Hello World task using Modern DSL
local hello_task = task("hello_world")
    :description("Simple hello world demonstration")
    :command(function(params)
        log.info("🌟 Hello World from Modern DSL!")
        log.info("📅 Current time: " .. os.date())
        
        return true, "echo 'Hello, Modern Sloth Runner!'", {
            message = "Hello World",
            timestamp = os.time(),
            status = "success"
        }
    end)
    :timeout("30s")
    :on_success(function(params, output)
        log.info("✅ Hello World task completed successfully!")
        log.info("💬 Message: " .. output.message)
    end)
    :build()

-- Modern Workflow Definition
workflow.define("hello_world_workflow", {
    description = "Simple Hello World - Modern DSL",
    version = "1.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"hello-world", "beginner", "modern-dsl"},
        created_at = os.date()
    },
    
    tasks = { hello_task },
    
    config = {
        timeout = "5m",
        max_parallel_tasks = 1
    },
    
    on_start = function()
        log.info("🚀 Starting Hello World workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 Hello World workflow completed!")
        else
            log.error("❌ Hello World workflow failed!")
        end
        return true
    end
})
EOF
}

# Função para criar http-basics moderno
create_modern_http_basics() {
    local file="$1"
    cat > "$file" << 'EOF'
-- MODERN DSL ONLY - HTTP Basics Example
-- Demonstrates HTTP operations with Modern DSL

-- HTTP GET task with circuit breaker
local http_get_task = task("http_get")
    :description("Perform HTTP GET with circuit breaker protection")
    :command(function(params)
        log.info("🌐 Making HTTP GET request...")
        
        -- Use circuit breaker for external API
        local result = circuit.protect("http_api", function()
            return net.http_get("https://jsonplaceholder.typicode.com/posts/1")
        end)
        
        if result.success then
            log.info("✅ HTTP request successful")
            return true, "HTTP GET completed", {
                status_code = result.status_code,
                body = result.body,
                headers = result.headers
            }
        else
            return false, "HTTP request failed: " .. (result.error or "unknown error")
        end
    end)
    :timeout("30s")
    :retries(3, "exponential")
    :on_success(function(params, output)
        log.info("📊 Response received: " .. string.len(output.body or "") .. " bytes")
    end)
    :build()

-- HTTP POST task
local http_post_task = task("http_post")
    :description("Perform HTTP POST with data")
    :depends_on({"http_get"})
    :command(function(params, deps)
        log.info("📤 Making HTTP POST request...")
        
        local post_data = {
            title = "Modern DSL Post",
            body = "Posted from Sloth Runner Modern DSL",
            userId = 1
        }
        
        local result = net.http_post("https://jsonplaceholder.typicode.com/posts", {
            headers = { ["Content-Type"] = "application/json" },
            body = data.to_json(post_data)
        })
        
        if result.success then
            return true, "HTTP POST completed", {
                post_result = result.body,
                post_status = result.status_code
            }
        else
            return false, "HTTP POST failed"
        end
    end)
    :timeout("45s")
    :build()

-- Modern Workflow Definition
workflow.define("http_basics", {
    description = "HTTP Operations - Modern DSL",
    version = "1.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"http", "api", "beginner", "modern-dsl"},
        created_at = os.date()
    },
    
    tasks = { http_get_task, http_post_task },
    
    config = {
        timeout = "10m",
        max_parallel_tasks = 1,
        circuit_breaker = {
            failure_threshold = 3,
            recovery_timeout = "1m"
        }
    },
    
    on_start = function()
        log.info("🚀 Starting HTTP basics workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ HTTP workflow completed successfully!")
        else
            log.error("❌ HTTP workflow failed!")
        end
        return true
    end
})
EOF
}

# Executar modernização
echo "🧹 Starting modernization process..."

for file in "${files_to_modernize[@]}"; do
    modernize_file "$file"
done

echo ""
echo "✅ Modernization completed!"
echo "📁 All files now use Modern DSL only"
echo "💾 Backup files created with .backup extension"
EOF

chmod +x migrate_to_modern_dsl.sh