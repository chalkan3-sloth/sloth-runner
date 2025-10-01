# 📦 Módulo de Gerenciamento de Pacotes

O módulo `pkg` fornece funcionalidades abrangentes de gerenciamento de pacotes multiplataforma. Ele detecta automaticamente o gerenciador de pacotes do sistema e fornece uma interface unificada.

## 🎯 Gerenciadores Suportados

- **apt / apt-get** (Debian/Ubuntu)
- **yum / dnf** (RHEL/CentOS/Fedora)
- **pacman** (Arch Linux)
- **zypper** (openSUSE)
- **brew** (macOS - Homebrew)

## 📚 Visão Geral das Funções

| Função | Descrição |
|--------|-----------|
| `pkg.install(pacotes)` | Instalar um ou mais pacotes |
| `pkg.remove(pacotes)` | Remover um ou mais pacotes |
| `pkg.update()` | Atualizar cache de pacotes |
| `pkg.upgrade()` | Atualizar todos os pacotes |
| `pkg.search(query)` | Buscar pacotes |
| `pkg.info(pacote)` | Obter informações do pacote |
| `pkg.list()` | Listar pacotes instalados |
| `pkg.is_installed(pacote)` | Verificar se está instalado |
| `pkg.get_manager()` | Obter gerenciador detectado |
| `pkg.clean()` | Limpar cache |
| `pkg.autoremove()` | Remover dependências não usadas |
| `pkg.which(executavel)` | Encontrar caminho do executável |
| `pkg.version(pacote)` | Obter versão do pacote |
| `pkg.deps(pacote)` | Listar dependências |
| `pkg.install_local(arquivo)` | Instalar de arquivo local |

## 📖 Documentação Detalhada

### Instalação e Remoção

#### `pkg.install(pacotes)`

Instala um ou mais pacotes.

**Parâmetros:**
- `pacotes`: String (pacote único) ou Tabela (múltiplos pacotes)

**Retorna:**
- `sucesso` (boolean): `true` em caso de sucesso
- `saida` (string): Saída do comando

**Exemplos:**

=== "DSL Moderno"
    ```lua
    local pkg = require("pkg")
    
    local instalar_ferramentas = task("instalar_ferramentas")
        :description("Instalar ferramentas de desenvolvimento")
        :command(function(this, params)
            log.info("Instalando ferramentas...")
            
            -- Instalar múltiplos pacotes
            local ferramentas = {"git", "curl", "wget", "vim"}
            local sucesso, saida = pkg.install(ferramentas)
            
            if sucesso then
                log.info("✅ Ferramentas instaladas com sucesso!")
                return true, "Instalado"
            else
                log.error("❌ Falha: " .. saida)
                return false, "Falhou"
            end
        end)
        :timeout("300s")
        :build()
    
    workflow.define("configurar")
        :tasks({ instalar_ferramentas })
    ```

=== "Com delegate_to"
    ```lua
    local pkg = require("pkg")
    
    local instalar_no_agente = task("instalar_no_agente")
        :description("Instalar pacotes no agente remoto")
        :command(function(this, params)
            log.info("Instalando no agente remoto...")
            
            local sucesso, saida = pkg.install({"htop", "ncdu"})
            
            if sucesso then
                log.info("✅ Instalado no agente!")
                return true, "OK"
            else
                return false, "Falhou"
            end
        end)
        :delegate_to("servidor-producao")
        :timeout("300s")
        :build()
    
    workflow.define("instalacao_remota")
        :tasks({ instalar_no_agente })
    ```

#### `pkg.remove(pacotes)`

Remove um ou mais pacotes.

**Exemplo:**

```lua
local pkg = require("pkg")

local limpeza = task("limpeza")
    :description("Remover pacotes desnecessários")
    :command(function(this, params)
        local pacotes = {"pacote1", "pacote2"}
        local sucesso, saida = pkg.remove(pacotes)
        
        if sucesso then
            log.info("✅ Pacotes removidos")
            return true, "Removido"
        end
        return false, "Falhou"
    end)
    :timeout("180s")
    :build()
```

### Informações de Pacotes

#### `pkg.search(query)`

Busca pacotes.

**Exemplo:**

```lua
local pkg = require("pkg")

local buscar_python = task("buscar_python")
    :description("Buscar pacotes Python")
    :command(function(this, params)
        local sucesso, resultados = pkg.search("python3")
        
        if sucesso then
            log.info("Resultados da busca:")
            local contador = 0
            for linha in resultados:gmatch("[^\r\n]+") do
                if contador < 10 then
                    log.info("  • " .. linha)
                end
                contador = contador + 1
            end
            return true, contador .. " resultados"
        end
        return false, "Busca falhou"
    end)
    :timeout("60s")
    :build()
```

#### `pkg.info(pacote)`

Obtém informações do pacote.

```lua
local sucesso, info = pkg.info("curl")
if sucesso then
    log.info("Info do pacote:\n" .. info)
end
```

#### `pkg.list()`

Lista pacotes instalados.

**Retorna:** `sucesso` (boolean), `pacotes` (tabela)

```lua
local sucesso, pacotes = pkg.list()
if sucesso and type(pacotes) == "table" then
    local contador = 0
    for _ in pairs(pacotes) do contador = contador + 1 end
    log.info("📦 Total: " .. contador .. " pacotes")
end
```

### Manutenção do Sistema

#### `pkg.update()`

Atualiza cache de pacotes.

```lua
local atualizar_cache = task("atualizar_cache")
    :description("Atualizar cache de pacotes")
    :command(function(this, params)
        log.info("Atualizando...")
        return pkg.update()
    end)
    :timeout("120s")
    :build()
```

#### `pkg.upgrade()`

Atualiza todos os pacotes.

#### `pkg.clean()`

Limpa cache de pacotes.

#### `pkg.autoremove()`

Remove dependências não utilizadas.

**Exemplo:**

```lua
local manutencao = task("manutencao")
    :description("Manutenção do sistema")
    :command(function(this, params)
        -- Atualizar
        pkg.update()
        
        -- Fazer upgrade
        pkg.upgrade()
        
        -- Limpar
        pkg.clean()
        pkg.autoremove()
        
        return true, "Manutenção completa"
    end)
    :timeout("600s")
    :build()
```

### Funções Avançadas

#### `pkg.is_installed(pacote)`

Verifica se está instalado.

```lua
local pkg = require("pkg")

local verificar_requisitos = task("verificar_requisitos")
    :description("Verificar pacotes necessários")
    :command(function(this, params)
        local necessarios = {"git", "curl", "wget"}
        local faltando = {}
        
        for _, nome_pkg in ipairs(necessarios) do
            local instalado, _ = pkg.is_installed(nome_pkg)
            if not instalado then
                table.insert(faltando, nome_pkg)
            end
        end
        
        if #faltando > 0 then
            return false, "Faltando: " .. table.concat(faltando, ", ")
        end
        
        return true, "Tudo OK"
    end)
    :build()
```

#### `pkg.get_manager()`

Retorna nome do gerenciador.

```lua
local gerenciador, err = pkg.get_manager()
log.info("Gerenciador: " .. (gerenciador or "desconhecido"))
```

#### `pkg.which(executavel)`

Encontra caminho do executável.

```lua
local caminho, err = pkg.which("git")
if caminho then
    log.info("Git em: " .. caminho)
end
```

## 🎯 Exemplos Completos

### Configuração de Ambiente de Desenvolvimento

```lua
local pkg = require("pkg")

local atualizar = task("atualizar")
    :command(function() return pkg.update() end)
    :build()

local instalar_ferramentas = task("instalar_ferramentas")
    :command(function()
        local ferramentas = {"git", "curl", "wget", "vim", "htop"}
        return pkg.install(ferramentas)
    end)
    :depends_on({"atualizar"})
    :build()

local verificar = task("verificar")
    :command(function()
        for _, ferramenta in ipairs({"git", "curl"}) do
            if pkg.is_installed(ferramenta) then
                local caminho = pkg.which(ferramenta)
                log.info("✅ " .. ferramenta .. " (" .. caminho .. ")")
            end
        end
        return true, "OK"
    end)
    :depends_on({"instalar_ferramentas"})
    :build()

workflow.define("configurar_dev")
    :tasks({ atualizar, instalar_ferramentas, verificar })
```

### Gerenciamento Distribuído

```lua
local pkg = require("pkg")

local atualizar_servidores = task("atualizar_servidores")
    :command(function() return pkg.update() end)
    :delegate_to("servidor-prod-1")
    :build()

local instalar_monitoramento = task("instalar_monitoramento")
    :command(function()
        return pkg.install({"htop", "iotop", "nethogs"})
    end)
    :delegate_to("servidor-prod-1")
    :depends_on({"atualizar_servidores"})
    :build()

workflow.define("configurar_monitoramento")
    :tasks({ atualizar_servidores, instalar_monitoramento })
```

### Auditoria do Sistema

```lua
local pkg = require("pkg")

local auditoria = task("auditoria")
    :command(function()
        log.info("📊 Auditoria do Sistema")
        log.info(string.rep("=", 60))
        
        local gerenciador = pkg.get_manager()
        log.info("Gerenciador: " .. gerenciador)
        
        local _, pacotes = pkg.list()
        local contador = 0
        for _ in pairs(pacotes) do contador = contador + 1 end
        log.info("Pacotes: " .. contador)
        
        local criticos = {"openssl", "curl"}
        for _, p in ipairs(criticos) do
            local instalado = pkg.is_installed(p)
            log.info((instalado and "✅" or "❌") .. " " .. p)
        end
        
        return true, "OK"
    end)
    :build()

workflow.define("auditoria")
    :tasks({ auditoria })
```

## 🚀 Melhores Práticas

1. **Atualizar antes de instalar:**
   ```lua
   pkg.update()
   pkg.install("pacote")
   ```

2. **Verificar antes de instalar:**
   ```lua
   if not pkg.is_installed("git") then
       pkg.install("git")
   end
   ```

3. **Limpar após operações:**
   ```lua
   pkg.clean()
   pkg.autoremove()
   ```

4. **Usar delegate_to para remoto:**
   ```lua
   :delegate_to("nome-servidor")
   ```

## ⚠️ Notas de Plataforma

- **Linux**: Requer sudo
- **macOS**: Homebrew não precisa de sudo
- **Arch**: Usa sintaxe do pacman
- **openSUSE**: Usa zypper

## 🔗 Veja Também

- [Módulo exec](exec.md)
- [Guia DSL Moderno](../modern-dsl/overview.md)
- [Agentes Distribuídos](../distributed.md)
