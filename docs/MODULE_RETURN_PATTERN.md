# Padrão de Retorno para Módulos Lua

## Visão Geral

Este documento descreve o padrão consistente para retorno de valores e erros em módulos Lua do sloth-runner.

## Problema

Anteriormente, os módulos tinham padrões inconsistentes:
- Alguns retornavam apenas 1 valor (resultado OU erro)
- Outros retornavam 2 valores (resultado, erro)
- Usuários não sabiam quando algo deu errado

## Solução: Padrão (result, error)

Todos os módulos devem seguir o padrão **(result, error)** do Lua:

```lua
local result, err = module.function(args)
if err then
    print("Erro: " .. err)
    return
end

-- Usar result
print("Sucesso: " .. result.message)
```

### Regras

1. **Sempre retornar 2 valores**: `(result, error)`
2. **Em caso de sucesso**: `(result, nil)`
3. **Em caso de erro**: `(nil, "mensagem de erro")`
4. **Fluent API**: Retornar `(self, nil)` para permitir encadeamento

## Exemplos

### 1. Operação com Sucesso

```go
func luaCreateResource(L *lua.LState) int {
    // ... lógica ...

    result := L.NewTable()
    result.RawSetString("changed", lua.LBool(true))
    result.RawSetString("message", lua.LString("Recurso criado"))

    L.Push(result)
    L.Push(lua.LNil) // Sempre incluir nil para erro
    return 2
}
```

```lua
-- Uso em Lua
local result, err = module.create_resource({name = "test"})
if err then
    error("Falha ao criar: " .. err)
end
print("Criado: " .. result.message)
```

### 2. Operação com Erro

```go
func luaDeleteResource(L *lua.LState) int {
    // ... lógica ...

    if err != nil {
        L.Push(lua.LNil)
        L.Push(lua.LString(fmt.Sprintf("Falha ao deletar: %v", err)))
        return 2
    }

    // ... sucesso ...
}
```

```lua
-- Uso em Lua
local result, err = module.delete_resource({name = "test"})
if err then
    print("Erro: " .. err)
    return
end
```

### 3. Idempotência (sem mudanças)

```go
func luaEnsureResource(L *lua.LState) int {
    // ... verificar se existe ...

    if alreadyExists {
        result := L.NewTable()
        result.RawSetString("changed", lua.LBool(false))
        result.RawSetString("message", lua.LString("Recurso já existe"))

        L.Push(result)
        L.Push(lua.LNil) // Sempre retornar nil para erro
        return 2
    }

    // ... criar recurso ...
}
```

```lua
-- Uso em Lua
local result, err = module.ensure_resource({name = "test"})
if err then
    error("Falha: " .. err)
end

if result.changed then
    print("Recurso criado")
else
    print("Recurso já existia")
end
```

### 4. Fluent API

```go
func instanceImage(L *lua.LState) int {
    instance := checkIncusInstance(L, 1)
    image := L.CheckString(2)
    instance.config["image"] = image

    L.Push(L.Get(1)) // Retornar self
    L.Push(lua.LNil) // Sempre incluir nil para erro
    return 2
}
```

```lua
-- Uso em Lua com encadeamento
local instance, err = incus.instance("myvm")
    :image("ubuntu/22.04")
    :config({memory = "2GB"})
    :create()

if err then
    error("Falha ao criar instância: " .. err)
end
```

## Helper Functions

Use as helper functions em `internal/modules/helpers.go`:

```go
import "github.com/chalkan3-sloth/sloth-runner/internal/modules"

// Retornar sucesso
func luaMyFunction(L *lua.LState) int {
    result := L.NewTable()
    result.RawSetString("data", lua.LString("valor"))
    return modules.Helpers.ReturnSuccess(L, result)
}

// Retornar erro
func luaMyFunction(L *lua.LState) int {
    return modules.Helpers.ReturnError(L, "algo deu errado")
}

// Retornar fluent (self, nil)
func builderMethod(L *lua.LState) int {
    // ... modificar self ...
    return modules.Helpers.ReturnFluentSuccess(L, L.Get(1))
}

// Retornar resultado com changed=true
func luaCreateResource(L *lua.LState) int {
    fields := map[string]lua.LValue{
        "id":   lua.LString("123"),
        "name": lua.LString("resource"),
    }
    return modules.Helpers.ReturnChanged(L, "Recurso criado", fields)
}

// Retornar resultado idempotente (changed=false)
func luaEnsureResource(L *lua.LState) int {
    return modules.Helpers.ReturnIdempotent(L, "Recurso já existe")
}
```

## Módulos Atualizados

Os seguintes módulos já seguem este padrão:

- ✅ `internal/modules/core/sloth.go` - Módulo sloth-runner
- ✅ `internal/modules/infra/incus.go` - Módulo Incus (com Fluent API)
- ✅ `internal/luainterface/modules/exec/exec.go` - Módulo exec

## Migrando Módulos Antigos

Se você tem um módulo usando o padrão antigo (retornando apenas 1 valor):

### Antes

```go
func luaOldFunction(L *lua.LState) int {
    if err != nil {
        errorTable := L.NewTable()
        errorTable.RawSetString("error", lua.LString("falhou"))
        L.Push(errorTable)
        return 1 // ❌ Retorna apenas 1 valor
    }

    resultTable := L.NewTable()
    resultTable.RawSetString("data", lua.LString("ok"))
    L.Push(resultTable)
    return 1 // ❌ Retorna apenas 1 valor
}
```

### Depois

```go
func luaNewFunction(L *lua.LState) int {
    if err != nil {
        L.Push(lua.LNil)
        L.Push(lua.LString("falhou"))
        return 2 // ✅ Sempre retorna 2 valores
    }

    resultTable := L.NewTable()
    resultTable.RawSetString("data", lua.LString("ok"))
    L.Push(resultTable)
    L.Push(lua.LNil) // ✅ Sempre inclui nil para erro
    return 2 // ✅ Sempre retorna 2 valores
}
```

## Benefícios

1. **Consistência**: Todos os módulos usam o mesmo padrão
2. **Clareza**: Usuários sempre sabem como verificar erros
3. **Previsibilidade**: `if err then` funciona em todos os lugares
4. **Compatibilidade**: Segue convenção padrão do Lua
5. **Fluent API**: Permite encadeamento elegante com `(self, nil)`

## Referências

- [Lua Error Handling Best Practices](https://www.lua.org/pil/8.4.html)
- [Go-Lua Documentation](https://pkg.go.dev/github.com/yuin/gopher-lua)
- Helper Functions: `internal/modules/helpers.go`
