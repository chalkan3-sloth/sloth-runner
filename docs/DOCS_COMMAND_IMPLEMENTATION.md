# Implementação do Comando `docs` - Manual Pages Style

**Status:** 🚧 1/13 Implementado | Template Pronto para os Demais

## 📚 Visão Geral

Cada comando sysadmin agora possui um subcomando `docs` que funciona como páginas `man` do Linux, fornecendo documentação completa com exemplos, opções e referências cruzadas.

## ✅ Implementado

### 1. services (COMPLETO)
**Comando:** `sloth-runner sysadmin services docs`

**Output:**
```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
         SLOTH-RUNNER SYSADMIN SERVICES(1)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

# NAME
    sloth-runner sysadmin services - Manage systemd/init.d services

# SYNOPSIS
    sloth-runner sysadmin services [subcommand] [options]

# OPTIONS
    list
        List all services on the specified agent...

    status [service-name]
        Check the current status of a specific service...

# EXAMPLES
    Example 1: List all services
    $ sloth-runner sysadmin services list --agent web-01
        Shows all services with status (active/inactive/failed)
...
```

**Arquivo:** `cmd/sloth-runner/commands/sysadmin/services/services.go`

## 🔨 Estrutura do Comando `docs`

### Componentes

1. **Subcomando docs** - Adicionado ao comando principal
2. **Função show{Command}Docs()** - Gera o conteúdo
3. **Função showDocs()** - Renderiza a documentação formatada

### Template

```go
// Adicionar ao comando principal
cmd.AddCommand(&cobra.Command{
    Use:   "docs",
    Short: "Show detailed documentation (man page style)",
    Long:  `Display comprehensive documentation with examples`,
    Run: func(cmd *cobra.Command, args []string) {
        show{Command}Docs()
    },
})

// Função de documentação
func show{Command}Docs() {
    title := "SLOTH-RUNNER SYSADMIN {COMMAND}(1)"
    description := "sloth-runner sysadmin {command} - Description"
    synopsis := "sloth-runner sysadmin {command} [subcommand] [options]"

    options := [][]string{
        {"subcommand", "Description of subcommand"},
        // ... more options
    }

    examples := [][]string{
        {
            "Example title",
            "sloth-runner sysadmin {command} ...",
            "What this example does",
        },
        // ... more examples
    }

    seeAlso := []string{
        "Related command 1",
        "Related command 2",
    }

    showDocs(title, description, synopsis, options, examples, seeAlso)
}

// Função de renderização (copiar do services.go)
func showDocs(title, description, synopsis string, options [][]string, examples [][]string, seeAlso []string) {
    // ... implementação completa em services/services.go
}
```

## 📋 Comandos Restantes

### Prioridade Alta (2 comandos)
1. **packages** - ✅ Template gerado em `/tmp/generate_docs_commands.py`
2. **resources** - ✅ Template gerado em `/tmp/generate_docs_commands.py`

### Prioridade Média (5 comandos)
3. **config** - Gerenciamento de configuração
4. **maintenance** - Manutenção do sistema
5. **network** - Diagnósticos de rede
6. **performance** - Monitoramento de performance
7. **backup** - Backup e restore

### Prioridade Baixa (2 comandos)
8. **security** - Auditoria de segurança
9. **deployment** - Deploy e rollback

### Já Implementados (3 comandos)
10. **logs** - Usa sistema próprio de documentação
11. **health** - Usa sistema próprio de documentação
12. **debug** - Usa sistema próprio de documentação

## 🛠️ Como Adicionar `docs` a um Comando

### Passo 1: Gerar Template
```bash
# Rodar o gerador Python
python3 /tmp/generate_docs_commands.py

# Output tem a função show{Command}Docs() pronta
```

### Passo 2: Adicionar ao Arquivo .go

1. **Abrir** `cmd/sloth-runner/commands/sysadmin/{comando}/{comando}.go`

2. **Adicionar imports** (se necessário):
```go
import (
    "fmt"
    "github.com/pterm/pterm"
    "github.com/spf13/cobra"
)
```

3. **Adicionar subcomando docs** antes do `return cmd`:
```go
// docs subcommand
cmd.AddCommand(&cobra.Command{
    Use:   "docs",
    Short: "Show detailed documentation (man page style)",
    Long:  `Display comprehensive documentation for the {command} command.`,
    Run: func(cmd *cobra.Command, args []string) {
        show{Command}Docs()
    },
})
```

4. **Copiar função show{Command}Docs()** do output do gerador

5. **Copiar função showDocs()** de `services/services.go` linhas 211-266

### Passo 3: Build e Test
```bash
# Build
go build -o sloth-runner-docs ./cmd/sloth-runner

# Test
./sloth-runner-docs sysadmin {comando} docs
```

## 📊 Estatísticas

### Implementação Atual
- **Comandos com docs:** 1/13 (7.7%)
- **Comandos prioritários:** 1/3 (33.3%)
- **Linhas de código:** ~60 por comando
- **Tempo estimado:** 10-15 min por comando

### Estimativa de Conclusão
- **Total de comandos restantes:** 9
- **Tempo total estimado:** 90-135 minutos
- **Complexidade:** Baixa (código template pronto)

## 🎯 Exemplo Completo: packages

### Código Completo para packages.go

```go
// No final do arquivo, antes do último }

// docs subcommand
cmd.AddCommand(&cobra.Command{
    Use:   "docs",
    Short: "Show detailed documentation (man page style)",
    Long:  `Display comprehensive documentation for the packages command with examples.`,
    Run: func(cmd *cobra.Command, args []string) {
        showPackagesDocs()
    },
})

func showPackagesDocs() {
    title := "SLOTH-RUNNER SYSADMIN PACKAGES(1)"
    description := "sloth-runner sysadmin packages - Manage system packages on remote agents"
    synopsis := "sloth-runner sysadmin packages [subcommand] [options]"

    options := [][]string{
        {"list", "List all installed packages on the agent."},
        {"search [package-name]", "Search for available packages in repositories."},
        {"install [package-name]", "Install a package. Supports dependency resolution."},
        {"remove [package-name]", "Remove an installed package. Can purge configurations."},
        {"update", "Update package repository lists (apt update, yum check-update)."},
        {"upgrade", "Upgrade all or specific packages to latest versions."},
        {"check-updates", "Check which packages have updates available."},
        {"info [package-name]", "Display detailed information about a package."},
        {"history", "Show history of package installations and updates."},
        {"docs", "Show this documentation page."},
    }

    examples := [][]string{
        {
            "List installed packages",
            "sloth-runner sysadmin packages list --agent web-01",
            "Shows all installed packages with versions",
        },
        {
            "Search for nginx",
            "sloth-runner sysadmin pkg search nginx --agent web-01",
            "Searches repositories for nginx packages",
        },
        {
            "Install package",
            "sloth-runner sysadmin packages install nginx --agent web-01",
            "Installs nginx with dependency resolution",
        },
        {
            "Rolling upgrade",
            "sloth-runner sysadmin packages upgrade --all-agents --strategy rolling",
            "Upgrades packages across all agents safely",
        },
    }

    seeAlso := []string{
        "sloth-runner sysadmin services - Service management",
        "sloth-runner sysadmin maintenance - System maintenance",
    }

    showDocs(title, description, synopsis, options, examples, seeAlso)
}

// showDocs function (copy from services.go lines 211-266)
func showDocs(title, description, synopsis string, options [][]string, examples [][]string, seeAlso []string) {
    // ... (copiar implementação completa)
}
```

## 🎨 Formato de Saída

### Cores e Formatação
- **Header:** Branco sobre fundo preto, largura total
- **Seções:** Amarelo bold (NAME, SYNOPSIS, OPTIONS, etc)
- **Opções:** Ciano (nomes dos subcomandos)
- **Exemplos:** Amarelo (títulos), Verde (comandos)
- **Footer:** Cinza

### Estrutura
1. **Header** - Título estilo man page
2. **NAME** - Nome e descrição curta
3. **SYNOPSIS** - Sintaxe básica
4. **OPTIONS** - Lista de subcomandos com descrições
5. **EXAMPLES** - Exemplos práticos com output esperado
6. **SEE ALSO** - Comandos relacionados
7. **Footer** - Informação de versão e help

## 📚 Benefícios

### Para Usuários
- ✅ Documentação sempre disponível offline
- ✅ Exemplos práticos e testados
- ✅ Formato familiar (similar ao `man`)
- ✅ Referências cruzadas entre comandos
- ✅ Não precisa sair do terminal

### Para Desenvolvedores
- ✅ Template padronizado
- ✅ Fácil manutenção
- ✅ Documentação versionada com código
- ✅ Geração automática via script Python

### Para Projeto
- ✅ Profissionalismo e qualidade
- ✅ Reduz suporte/documentação externa
- ✅ Facilita onboarding
- ✅ Consistency across commands

## 🚀 Próximos Passos

### Imediato
1. ✅ Template implementado em services
2. ✅ Gerador Python funcional
3. ✅ Documentação de processo completa

### Curto Prazo
1. [ ] Adicionar docs aos 2 comandos prioritários (packages, resources)
2. [ ] Adicionar docs aos 5 comandos de prioridade média
3. [ ] Adicionar docs aos 2 comandos restantes

### Médio Prazo
1. [ ] Criar tests unitários para comandos docs
2. [ ] Adicionar flag --format (text, markdown, html)
3. [ ] Gerar man pages reais (.1 files)
4. [ ] Integração com help system

## 📖 Referências

- **Arquivo principal:** `cmd/sloth-runner/commands/sysadmin/services/services.go`
- **Gerador:** `/tmp/generate_docs_commands.py`
- **Helper (futuro):** `cmd/sloth-runner/commands/sysadmin/docs_helper.go`

## ✨ Exemplo de Uso

```bash
# Ver documentação do services
$ sloth-runner sysadmin services docs

# Ver documentação do packages (quando implementado)
$ sloth-runner sysadmin packages docs

# Ver documentação do resources (quando implementado)
$ sloth-runner sysadmin resources docs

# Buscar na documentação (futuro)
$ sloth-runner sysadmin services docs | grep restart

# Export para markdown (futuro)
$ sloth-runner sysadmin services docs --format markdown > services.md
```

---

**Status:** ✅ Template Completo e Funcionando
**Próximo:** Aplicar template aos 9 comandos restantes
**Tempo Estimado:** 2-3 horas para todos os comandos
