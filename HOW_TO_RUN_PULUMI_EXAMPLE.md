# Como Rodar o Exemplo Pulumi

## 🚀 **Execução Rápida**

### **Comando Principal:**
```bash
sloth-runner run -f examples/pulumi_config_example.sloth pulumi_complete_example
```

## ⚙️ **Pré-requisitos**

### **1. Pulumi CLI Instalado:**
```bash
# Instalar Pulumi (se não estiver instalado)
curl -fsSL https://get.pulumi.com | sh

# Verificar instalação
pulumi version
```

### **2. Git Instalado:**
```bash
# Verificar se git está disponível
git --version
```

### **3. Sloth Runner Compilado:**
```bash
# Se não estiver compilado, execute:
cd /Users/chalkan3/.projects/task-runner
go build -o sloth-runner ./cmd/sloth-runner
cp sloth-runner $HOME/.local/bin/

# Verificar instalação
sloth-runner --version
```

## 📋 **O que o Exemplo Faz**

### **Task 1: Clone do Repositório**
- 📁 Clona `https://github.com/chalkan3/go-do-droplet`
- 🗂️ Para o diretório `/tmp/pulumi-project`

### **Task 2: Configuração e Preview**
- 🔐 Login Pulumi local: `pulumi login file://.`
- 📋 Criar stack: `pulumi stack init dev`
- ⚙️ Configurar valores:
  ```
  pulumi config set dropletName sloth-runner
  pulumi config set region nyc3
  pulumi config set size s-1vcpu-1gb
  pulumi config set image ubuntu-22-04-x64
  pulumi config set environment dev
  pulumi config set project main
  ```
- 🔍 Executar preview com environment variables

## 🎯 **Comandos de Execução**

### **Execução Completa:**
```bash
# Comando principal
sloth-runner run -f examples/pulumi_config_example.sloth pulumi_complete_example
```

### **Com Verbose (para debug):**
```bash
# Se quiser ver mais detalhes
sloth-runner run -f examples/pulumi_config_example.sloth pulumi_complete_example --verbose
```

### **Executar apenas uma task específica:**
```bash
# Apenas clone
sloth-runner run -f examples/pulumi_config_example.sloth clone_go_droplet_repo

# Apenas configuração (requer clone primeiro)
sloth-runner run -f examples/pulumi_config_example.sloth setup_pulumi_config
```

## 📊 **Output Esperado**

### **1. Clone Task:**
```
📡 Cloning Go DigitalOcean droplet repository...
✅ Repository cloned successfully!
✅ === CLONE SUCCESS ===
Repository available at: /tmp/pulumi-project
```

### **2. Setup Task:**
```
🔐 Pulumi login (local backend)...
✅ Pulumi login successful
📋 Setting up Pulumi stack...
✅ Stack 'dev' ready
⚙️ Setting Pulumi configuration values...
  ✅ dropletName: sloth-runner
  ✅ region: nyc3
  ✅ size: s-1vcpu-1gb
  ✅ image: ubuntu-22-04-x64
  ✅ environment: dev
  ✅ project: main
📊 === PULUMI PREVIEW OUTPUT ===
[preview output aqui]
📊 === END PREVIEW OUTPUT ===
```

## 🛠️ **Troubleshooting**

### **Erro: "pulumi not found"**
```bash
# Instalar Pulumi
curl -fsSL https://get.pulumi.com | sh
export PATH=$PATH:$HOME/.pulumi/bin
```

### **Erro: "git not found"**
```bash
# macOS
brew install git

# Ubuntu/Debian
sudo apt-get install git
```

### **Erro: "sloth-runner not found"**
```bash
# Compilar e instalar
cd /Users/chalkan3/.projects/task-runner
go build -o sloth-runner ./cmd/sloth-runner
cp sloth-runner $HOME/.local/bin/
export PATH=$PATH:$HOME/.local/bin
```

### **Erro: "permission denied" no /tmp**
```bash
# Verificar permissões
ls -la /tmp
sudo chmod 755 /tmp
```

### **Erro: "repository already exists"**
```bash
# Limpar diretório anterior
rm -rf /tmp/pulumi-project
```

## 🔧 **Customização**

### **Alterar Configurações:**
Edite `examples/pulumi_config_example.sloth` para mudar:
- Workdir (padrão: `/tmp/pulumi-project`)
- Repositório (padrão: `https://github.com/chalkan3/go-do-droplet`)
- Stack name (padrão: `dev`)
- Valores de configuração

### **Adicionar Environment Variables:**
```lua
local preview_success, preview_output = client:preview({ 
    envs = { 
        PULUMI_CONFIG_PASSPHRASE = "",
        DIGITALOCEAN_TOKEN = "seu_token_aqui",
        TF_LOG = "DEBUG"
    }
})
```

## 📁 **Estrutura de Arquivos Após Execução**

```
/tmp/pulumi-project/
├── main.go                 # Código Go do Pulumi
├── go.mod                  # Dependências Go
├── Pulumi.yaml            # Configuração do projeto
├── Pulumi.dev.yaml        # Configuração do stack dev
└── .pulumi/               # Estado local do Pulumi
```

## ✅ **Sucesso Esperado**

Se tudo correr bem, você verá:
1. ✅ Clone do repositório concluído
2. ✅ Login Pulumi realizado
3. ✅ Stack criado/selecionado
4. ✅ 6 configurações aplicadas
5. ✅ Preview executado com sucesso
6. 📊 Output do preview mostrando recursos planejados

## 🎉 **Próximos Passos**

Após o exemplo funcionar, você pode:
1. Modificar as configurações
2. Adicionar `pulumi up` para fazer deploy real
3. Usar com tokens reais do DigitalOcean
4. Integrar com CI/CD pipelines