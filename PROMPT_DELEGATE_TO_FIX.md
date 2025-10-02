# Prompt para IA ajustar delegate_to no Sloth Runner

## Problema
O delegate_to não está funcionando corretamente nos arquivos .sloth. A execução remota funciona perfeitamente com comandos diretos:

```bash
sloth-runner agent run ladyguica "ls -la \$HOME" --master 192.168.1.29:50053
```

Mas quando usamos um arquivo .sloth com delegate_to, ocorre erro: "attempt to call a non-function object"

## Estrutura do projeto
- `/cmd/sloth-runner/main.go` - Comando principal
- `/internal/luainterface/modern_dsl.go` - DSL moderno com task() e workflow.define()
- `/internal/taskrunner/taskrunner.go` - Executor de tasks

## Exemplo que deveria funcionar mas não funciona

Arquivo: `examples/agents/simple_delegate_ls.sloth`
```lua
local ls_task = task("list_files")
    :description("List files on remote agent")
    :command(function(this, params)
        local output, err, failed = exec.run("ls -la $HOME")
        if not failed then
            log.info("Files listed successfully:")
            log.info(output)
            return true, "Files listed"
        else
            log.error("Failed to list files: " .. (err or "unknown error"))
            return false, "Failed to list files"
        end
    end)
    :delegate_to("ladyguica")
    :timeout("30s")
    :build()

workflow.define("simple_ls_test")
    :description("Simple LS test on remote agent")
    :version("1.0.0")
    :tasks({ ls_task })
    :config({ timeout = "1m" })
    :on_complete(function(success, results)
        if success then
            log.info("✅ LS remoto executado com sucesso!")
        else
            log.error("❌ Falha na execução do LS remoto")
        end
        return true
    end)
```

## Comando que falha
```bash
sloth-runner run -f examples/agents/simple_delegate_ls.sloth simple_ls_test
```

## Erro obtido
```
rpc error: code = Unknown desc = failed to load lua script: <string>:2: attempt to call a non-function object
```

## Funcionando
- ✅ Resolução de nomes de agentes (ladyguica -> 192.168.1.16:50051)
- ✅ Conexão com master
- ✅ Comando direto: `sloth-runner agent run ladyguica "comando"`
- ✅ Parsing do delegate_to no workflow

## Não funcionando
- ❌ Execução de .sloth com delegate_to no agente remoto
- ❌ Script enviado para o agente não consegue executar

## Contexto técnico
1. O delegate_to está sendo corretamente parseado e o agent está sendo resolvido
2. O script está sendo enviado para o agente via gRPC
3. O agente recebe o script mas falha ao executar: "attempt to call a non-function object"
4. O problema parece ser na serialização/envio do script Lua para o agente

## Tarefa
Preciso que você analise o código e identifique:
1. Como o script é serializado/enviado para o agente remoto
2. Por que o script enviado não consegue executar no agente
3. Ajustar o código para que scripts .sloth com delegate_to funcionem corretamente

## Logs completos disponíveis
- Os logs mostram que delegate_to é parseado corretamente
- Agent é resolvido com sucesso (ladyguica -> 192.168.1.16:50051)  
- Conexão gRPC é estabelecida
- Script é enviado mas falha na execução no agente

Foque em descobrir como o script Lua é enviado para o agente e por que ele não consegue executar lá.