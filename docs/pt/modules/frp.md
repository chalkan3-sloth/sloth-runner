# Módulo FRP - Fast Reverse Proxy

O módulo FRP fornece integração com [FRP (Fast Reverse Proxy)](https://github.com/fatedier/frp) para expor servidores locais por trás de NAT e firewalls. Oferece uma API fluente para gerenciar tanto servidores FRP (frps) quanto clientes (frpc) com suporte para instalação, configuração e gerenciamento de ciclo de vida.

## Recursos

- **Gerenciamento de Instalação**: Download e instalação automática dos binários FRP
- **Configuração de Servidor**: Configure servidor FRP com configurações baseadas em TOML
- **Configuração de Cliente**: Configure cliente FRP com múltiplas configurações de proxy
- **Integração Systemd**: Gerencie FRP como serviços systemd
- **TOML para Tabelas Lua**: Conversão automática entre configurações TOML e tabelas Lua
- **API Fluente**: Encadeie métodos de configuração para código limpo e legível
- **Delegação de Agentes**: Execute operações em agentes remotos via `delegate_to()`

## Instalação

O módulo FRP é registrado automaticamente de forma global. Não é necessário `require()`.

```lua
-- FRP está disponível imediatamente
local server = frp.server("meu-servidor")
local client = frp.client("meu-cliente")
```

## Funções do Módulo

### frp.install(version?, target?)

Instala os binários FRP (tanto frps quanto frpc) local ou remotamente.

**Parâmetros:**
- `version` (string, opcional): Versão do FRP a instalar (padrão: "latest")
- `target` (string, opcional): Nome do agente remoto para instalação

**Retorna:** `(resultado, erro)`

**Exemplo:**

```lua
-- Instalar última versão localmente
local resultado, erro = frp.install()

-- Instalar versão específica em agente remoto
local resultado, erro = frp.install("0.52.0", "meu-agente")
```

### frp.server(nome?)

Cria uma nova instância do construtor de servidor FRP.

**Parâmetros:**
- `nome` (string, opcional): Nome do servidor (padrão: "frps")

**Retorna:** Instância do construtor `FrpServer`

**Exemplo:**

```lua
local servidor = frp.server("servidor-producao")
```

### frp.client(nome?)

Cria uma nova instância do construtor de cliente FRP.

**Parâmetros:**
- `nome` (string, opcional): Nome do cliente (padrão: "frpc")

**Retorna:** Instância do construtor `FrpClient`

**Exemplo:**

```lua
local cliente = frp.client("meu-cliente")
```

## API do Servidor FRP

### server:config(tabela_config)

Define configuração do servidor a partir de tabela Lua (método fluente).

**Parâmetros:**
- `tabela_config` (table): Opções de configuração

**Opções Comuns de Configuração:**
- `bindAddr` - Endereço de bind do servidor (padrão: "0.0.0.0")
- `bindPort` - Porta de bind do servidor (padrão: 7000)
- `vhostHTTPPort` - Porta vhost HTTP (padrão: 80)
- `vhostHTTPSPort` - Porta vhost HTTPS (padrão: 443)
- `auth.method` - Método de autenticação ("token")
- `auth.token` - Token de autenticação
- `webServer` - Configuração do painel web
- `log` - Configuração de logging

**Retorna:** `(self, nil)` para encadeamento de métodos

**Exemplo:**

```lua
local servidor = frp.server("meu-servidor")
    :config({
        bindAddr = "0.0.0.0",
        bindPort = 7000,
        vhostHTTPPort = 8080,
        auth = {
            method = "token",
            token = "meu_token_secreto"
        },
        webServer = {
            addr = "0.0.0.0",
            port = 7500,
            user = "admin",
            password = "admin123"
        },
        log = {
            to = "/var/log/frp/frps.log",
            level = "info",
            maxDays = 7
        }
    })
```

### server:config_path(caminho)

Define o caminho do arquivo de configuração (método fluente).

**Parâmetros:**
- `caminho` (string): Caminho para arquivo de configuração TOML

**Retorna:** `(self, nil)`

### server:version(versao)

Define a versão do FRP para instalação (método fluente).

**Parâmetros:**
- `versao` (string): String de versão (ex: "0.52.0" ou "latest")

**Retorna:** `(self, nil)`

### server:delegate_to(agente)

Executa operações em um agente remoto (método fluente).

**Parâmetros:**
- `agente` (string): Nome do agente

**Retorna:** `(self, nil)`

### server:save_config()

Salva configuração em arquivo TOML (método de ação).

**Retorna:** `(resultado, erro)`

### server:load_config()

Carrega configuração de arquivo TOML (método de ação).

**Retorna:** `(tabela_config, erro)`

### server:install()

Instala binários FRP e serviço systemd (método de ação).

**Retorna:** `(resultado, erro)`

### server:start()

Inicia o serviço do servidor FRP (método de ação).

**Retorna:** `(resultado, erro)`

### server:stop()

Para o serviço do servidor FRP (método de ação).

**Retorna:** `(resultado, erro)`

### server:restart()

Reinicia o serviço do servidor FRP (método de ação).

**Retorna:** `(resultado, erro)`

### server:status()

Obtém status do serviço do servidor FRP (método de ação).

**Retorna:** `(saida_status, erro)`

### server:enable()

Habilita servidor FRP para iniciar no boot (método de ação).

**Retorna:** `(resultado, erro)`

### server:disable()

Desabilita servidor FRP de iniciar no boot (método de ação).

**Retorna:** `(resultado, erro)`

## API do Cliente FRP

### client:config(tabela_config)

Define configuração do cliente a partir de tabela Lua (método fluente).

**Parâmetros:**
- `tabela_config` (table): Opções de configuração

**Retorna:** `(self, nil)`

### client:server(endereco, porta)

Define detalhes de conexão ao servidor FRP (método fluente).

**Parâmetros:**
- `endereco` (string): Endereço do servidor
- `porta` (number): Porta do servidor

**Retorna:** `(self, nil)`

**Exemplo:**

```lua
local cliente = frp.client()
    :server("frp.exemplo.com", 7000)
```

### client:proxy(config_proxy)

Adiciona uma configuração de proxy (método fluente). Pode ser chamado múltiplas vezes.

**Parâmetros:**
- `config_proxy` (table): Configuração do proxy

**Opções Comuns de Proxy:**
- `name` - Nome único do proxy
- `type` - Tipo do proxy ("tcp", "http", "https", "udp")
- `localIP` - IP do serviço local (padrão: "127.0.0.1")
- `localPort` - Porta do serviço local
- `remotePort` - Porta remota (para TCP/UDP)
- `customDomains` - Domínios customizados (para HTTP/HTTPS)

**Retorna:** `(self, nil)`

**Exemplo:**

```lua
local cliente = frp.client()
    :proxy({
        name = "web",
        type = "http",
        localIP = "127.0.0.1",
        localPort = 3000,
        customDomains = {"meuapp.exemplo.com"}
    })
    :proxy({
        name = "ssh",
        type = "tcp",
        localPort = 22,
        remotePort = 6000
    })
```

## Exemplos Completos

### Exemplo 1: Configuração de Servidor FRP

```lua
local tarefa_instalacao = task("instalar_frp")
    :command(function()
        local resultado, erro = frp.install("latest")
        if erro then
            return false, erro
        end
        return true, "FRP instalado"
    end)
    :build()

local tarefa_config = task("configurar_servidor")
    :depends_on({"instalar_frp"})
    :command(function()
        local servidor = frp.server("producao")
            :config({
                bindAddr = "0.0.0.0",
                bindPort = 7000,
                auth = {
                    method = "token",
                    token = "token_seguro_123"
                }
            })

        local resultado, erro = servidor:save_config()
        if erro then
            return false, erro
        end
        return true, "Configuração salva"
    end)
    :build()

local tarefa_iniciar = task("iniciar_servidor")
    :depends_on({"configurar_servidor"})
    :command(function()
        local servidor = frp.server()

        local _, erro = servidor:enable()
        if erro then return false, erro end

        local _, erro = servidor:start()
        if erro then return false, erro end

        return true, "Servidor iniciado"
    end)
    :build()

workflow
    .define("configuracao_servidor_frp")
    :tasks({tarefa_instalacao, tarefa_config, tarefa_iniciar})
```

### Exemplo 2: Cliente FRP com Múltiplos Serviços

```lua
local configurar_cliente = task("configurar_cliente_frp")
    :command(function()
        local cliente = frp.client("servicos")
            :server("frp.exemplo.com", 7000)
            :config({
                auth = {
                    method = "token",
                    token = "token_seguro_123"
                }
            })
            -- Aplicação web
            :proxy({
                name = "web",
                type = "http",
                localPort = 3000,
                customDomains = {"meuapp.exemplo.com"}
            })
            -- Acesso SSH
            :proxy({
                name = "ssh",
                type = "tcp",
                localPort = 22,
                remotePort = 6000
            })
            -- Banco de dados PostgreSQL
            :proxy({
                name = "postgres",
                type = "tcp",
                localPort = 5432,
                remotePort = 6432
            })

        local resultado, erro = cliente:save_config()
        if erro then
            return false, erro
        end

        local _, erro = cliente:install()
        if erro then return false, erro end

        local _, erro = cliente:enable()
        if erro then return false, erro end

        local _, erro = cliente:start()
        if erro then return false, erro end

        return true, "Cliente configurado e iniciado"
    end)
    :build()
```

### Exemplo 3: Deploy em Agente Remoto

```lua
local deploy_servidor = task("deploy_servidor_frp")
    :command(function()
        local servidor = frp.server("servidor-remoto")
            :delegate_to("agente-producao")
            :version("0.52.0")
            :config({
                bindAddr = "0.0.0.0",
                bindPort = 7000,
                vhostHTTPPort = 80,
                vhostHTTPSPort = 443
            })

        local _, erro = servidor:install()
        if erro then return false, "Instalação falhou: " .. erro end

        local _, erro = servidor:save_config()
        if erro then return false, "Configuração falhou: " .. erro end

        local _, erro = servidor:enable()
        if erro then return false, "Habilitar falhou: " .. erro end

        local _, erro = servidor:start()
        if erro then return false, "Iniciar falhou: " .. erro end

        return true, "Servidor implantado em agente remoto"
    end)
    :build()
```

## Configuração TOML

O módulo FRP trata automaticamente arquivos de configuração TOML. Configurações fornecidas via tabelas Lua são convertidas para formato TOML ao salvar.

### Exemplo de Configuração do Servidor (frps.toml)

```toml
bindAddr = "0.0.0.0"
bindPort = 7000
vhostHTTPPort = 80
vhostHTTPSPort = 443

[auth]
method = "token"
token = "token_seguro_123"

[webServer]
addr = "0.0.0.0"
port = 7500
user = "admin"
password = "admin123"

[log]
to = "/var/log/frp/frps.log"
level = "info"
maxDays = 7
```

### Exemplo de Configuração do Cliente (frpc.toml)

```toml
serverAddr = "frp.exemplo.com"
serverPort = 7000

[auth]
method = "token"
token = "token_seguro_123"

[[proxies]]
name = "web"
type = "http"
localIP = "127.0.0.1"
localPort = 3000
customDomains = ["meuapp.exemplo.com"]

[[proxies]]
name = "ssh"
type = "tcp"
localPort = 22
remotePort = 6000
```

## Melhores Práticas

1. **Use Autenticação**: Sempre configure autenticação baseada em token para produção
2. **Fixe Versões**: Especifique versão exata do FRP ao invés de "latest" para reprodutibilidade
3. **Deploy Remoto**: Use `delegate_to()` para gerenciar FRP em máquinas remotas
4. **Backup de Configuração**: Faça backup das configurações antes de atualizações
5. **Habilite no Boot**: Use `:enable()` para garantir que o FRP inicie após reinicializações
6. **Monitoramento**: Verifique status do serviço regularmente com `:status()`
7. **Logging**: Configure arquivos de log apropriados para troubleshooting
8. **Gerenciamento de Serviços**: Use integração systemd para gerenciamento confiável

## Troubleshooting

### Problemas de Instalação

```lua
-- Verificar se binários estão instalados
local resultado = exec.run("which frps frpc")
log.info(resultado)

-- Verificar versão
local versao = exec.run("frps --version")
log.info(versao)
```

### Problemas de Configuração

```lua
-- Validar sintaxe TOML
local config, erro = servidor:load_config()
if erro then
    log.error("Validação de config falhou: " .. erro)
end

-- Verificar permissões de arquivo
local perms = exec.run("ls -l /etc/frp/frps.toml")
log.info(perms)
```

### Problemas de Serviço

```lua
-- Verificar status do serviço
local status, _ = servidor:status()
log.info(status)

-- Ver logs
local logs = exec.run("journalctl -u frps -n 50")
log.info(logs)

-- Verificar se porta está em uso
local verif_porta = exec.run("ss -tlnp | grep 7000")
log.info(verif_porta)
```

## Veja Também

- [Módulo Incus](./incus.md) - Gerenciamento de containers/VMs
- [Módulo Systemd](./systemd.md) - Gerenciamento de serviços
- [Módulo Exec](./exec.md) - Execução de comandos
- [Documentação Oficial do FRP](https://github.com/fatedier/frp)
