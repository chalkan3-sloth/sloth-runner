# ⚡ Tutorial Rápido

Para documentação completa em português, visite:

## 🚀 Início Rápido

### Instalação

```bash
# Download
curl -sSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh | bash

# Ou via Go
go install github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner@latest
```

### Primeiro Workflow

Crie um arquivo `hello.sloth`:

```lua
local hello_task = task("hello")
    :description("Minha primeira task")
    :command(function(this, params)
        log.info("🦥 Olá do Sloth Runner!")
        return true, "Task executada com sucesso"
    end)
    :build()

workflow.define("hello_world")
    :description("Meu primeiro workflow")
    :version("1.0.0")
    :tasks({ hello_task })
    :on_complete(function(success, results)
        if success then
            log.info("✅ Workflow concluído!")
        end
    end)
```

Execute:

```bash
sloth-runner run -f hello.sloth
```

## 📚 Próximos Passos

- [Conceitos Fundamentais](./core-concepts.md)
- [Exemplos Avançados](./advanced-examples.md)
- [Recursos Avançados](./advanced-features.md)

Para o tutorial completo, veja: [Tutorial Principal](../TUTORIAL.md)
