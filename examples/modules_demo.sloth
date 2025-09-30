-- Exemplo de uso do sistema de módulos melhorado

-- 1. Sistema de Help Interativo
print("=== 🆘 Sistema de Help ===")
help() -- Mostra help geral
help.modules() -- Lista todos os módulos
help.module("http") -- Help específico do módulo HTTP

-- 2. Módulo HTTP Melhorado
print("\n=== 🌐 HTTP Client Avançado ===")
local http = require("http")

-- GET simples com retry automático
local result = http.get({
    url = "https://api.github.com/zen",
    timeout = 10,
    max_retries = 3,
    retry_delay = 2
})

if result.success then
    print("GitHub Zen:", result.data.body)
    print("Elapsed:", result.data.elapsed_ms, "ms")
    print("Status:", result.data.status_code)
else
    print("Error:", result.error)
end

-- POST with JSON and validation
local api_result = http.post({
    url = "https://httpbin.org/post",
    json = {
        name = "Sloth Runner",
        version = "2.0",
        features = {"modules", "validation", "help"}
    },
    headers = {
        ["User-Agent"] = "Sloth-Runner/2.0"
    },
    timeout = 15
})

if api_result.success and api_result.data.json then
    print("JSON Response:", api_result.data.json.data.name)
end

-- 3. Sistema de Validação
print("\n=== ✅ Validação de Dados ===")
local validate = require("validate")

-- Validação de email
local email_check = validate.email("user@example.com")
print("Email válido:", email_check.valid)

-- Validação de URL
local url_check = validate.url("https://github.com/chalkan3/sloth-runner")
if url_check.valid then
    print("URL válida - Host:", url_check.host, "Scheme:", url_check.scheme)
end

-- Validação de schema complexo
local user_data = {
    name = "John Doe",
    email = "john@example.com",
    age = "25",
    website = "https://johndoe.com"
}

local schema = {
    name = { required = "true", type = "string" },
    email = { required = "true", type = "string" },
    age = { required = "true", type = "string" },
    website = { required = "false", type = "string" }
}

local schema_result = validate.schema(user_data, schema)
if schema_result.valid then
    print("✅ Dados válidos!")
else
    print("❌ Erros de validação:")
    for i, error in ipairs(schema_result.errors) do
        print("  -", error)
    end
end

-- Validação de comprimento
local length_check = validate.length("Hello World", {
    min = 5,
    max = 20
})
print("Comprimento válido:", length_check.valid, "- Length:", length_check.length)

-- Sanitização de dados
local sanitized = validate.sanitize("<script>alert('xss')</script>", "html")
print("Original:", sanitized.original)
print("Sanitized:", sanitized.sanitized)

-- 4. Descoberta de Módulos
print("\n=== 📚 Descoberta de Módulos ===")

-- Lista todos os módulos
local all_modules = modules()
print("Módulos disponíveis:", #all_modules)
for i, name in ipairs(all_modules) do
    local info = module_info(name)
    if info then
        print(string.format("  • %s v%s (%s) - %s", 
            info.name, info.version, info.category, info.description))
    end
end

-- Busca por funcionalidade
print("\n🔍 Busca por 'http':")
help.search("http")

-- 5. Exemplo de uso com tratamento de erros
print("\n=== 🛡️  Tratamento de Erros ===")

local function safe_http_call(url)
    local result = http.get({
        url = url,
        timeout = 5,
        max_retries = 2
    })
    
    if not result.success then
        print("❌ Falha na requisição:", result.error)
        return nil
    end
    
    if result.data.status_code >= 400 then
        print("❌ HTTP Error:", result.data.status_code)
        return nil
    end
    
    return result.data
end

-- Teste com URL válida
local good_result = safe_http_call("https://httpbin.org/json")
if good_result then
    print("✅ Requisição bem-sucedida!")
end

-- Teste com URL inválida
local bad_result = safe_http_call("https://nonexistent-domain-12345.com")
-- Deve mostrar erro tratado

-- 6. Exemplo de composição de módulos
print("\n=== 🔧 Composição de Módulos ===")

local function fetch_and_validate_user(user_id)
    -- Buscar dados do usuário
    local user_response = http.get({
        url = "https://jsonplaceholder.typicode.com/users/" .. user_id
    })
    
    if not user_response.success then
        return { success = false, error = "Failed to fetch user" }
    end
    
    local user = user_response.data.json
    if not user then
        return { success = false, error = "Invalid response format" }
    end
    
    -- Validar dados do usuário
    local email_validation = validate.email(user.email)
    local website_validation = validate.url(user.website)
    
    return {
        success = true,
        user = user,
        validations = {
            email_valid = email_validation.valid,
            website_valid = website_validation.valid
        }
    }
end

-- Testar com usuário ID 1
local user_result = fetch_and_validate_user(1)
if user_result.success then
    print("👤 Usuário:", user_result.user.name)
    print("📧 Email válido:", user_result.validations.email_valid)
    print("🌐 Website válido:", user_result.validations.website_valid)
end

print("\n🎉 Demo completed! Explore help.examples() for more usage patterns.")