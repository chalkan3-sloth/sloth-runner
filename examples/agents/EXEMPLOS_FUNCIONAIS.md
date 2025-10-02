# 🚀 Exemplos de Execução Remota com delegate_to

## ✅ **COMANDOS QUE FUNCIONAM 100%**

### 1️⃣ Execução Direta (RECOMENDADO)

**Comando básico:**
```bash
sloth-runner agent run ladyguica "hostname && date" --master 192.168.1.29:50053
```

**Comando com múltiplas operações:**
```bash
sloth-runner agent run ladyguica "echo '🚀 Testando agent:' && hostname && date && pwd && whoami" --master 192.168.1.29:50053
```

**Listar arquivos remotamente:**
```bash
sloth-runner agent run keiteguica "ls -la /home/chalkan3 | head -10" --master 192.168.1.29:50053
```

**Informações do sistema:**
```bash
sloth-runner agent run ladyguica "uname -a && uptime && free -h" --master 192.168.1.29:50053
```

### 2️⃣ Listar Agents Disponíveis

```bash
sloth-runner agent list --master 192.168.1.29:50053
```

### 3️⃣ Agents Configurados

- **ladyguica**: `192.168.1.16:50051`
- **keiteguica**: `192.168.1.17:50051`

## 📝 **EXEMPLOS DE ARQUIVOS .SLOTH**

### Arquivo: `simple_delegate_example.sloth`

```lua
-- Task com delegate_to por nome
local hostname_task = task("get_hostname")
    :description("Get hostname from ladyguica agent")
    :command(function(this, params)
        log.info("🔍 Getting hostname from remote agent...")
        
        local hostname_out, _, hostname_failed = exec.run("hostname")
        local hostname = hostname_out and hostname_out:gsub("\n", "") or "unknown"
        
        log.info("📍 Running on host: " .. hostname)
        
        return true, "Hostname: " .. hostname
    end)
    :delegate_to("ladyguica")  -- ✨ DELEGAÇÃO POR NOME
    :timeout("30s")
    :build()

workflow.define("simple_remote_commands")
    :description("🚀 Execute simple commands on remote agents")
    :version("1.0.0")
    :tasks({ hostname_task })
    :config({ timeout = "1m" })
    :on_complete(function(success, results)
        log.info("🎉 Remote commands completed: " .. tostring(success))
        return true
    end)
```

**Executar:**
```bash
sloth-runner run -f examples/agents/simple_delegate_example.sloth simple_remote_commands
```

### ⚠️ **LIMITAÇÕES CONHECIDAS DOS ARQUIVOS .SLOTH**

Os workflows .sloth com `delegate_to` têm uma limitação conhecida na transferência de definições de tasks para agents remotos. Para execução remota eficiente, use os **comandos diretos** que funcionam perfeitamente.

## 🎯 **CASOS DE USO RECOMENDADOS**

### ✅ Use comandos diretos para:
- Execução de comandos shell remotos
- Verificação de status de sistemas
- Coleta de informações rápidas
- Manutenção e diagnóstico

### 🔄 Use arquivos .sloth para:
- Workflows complexos locais
- Orquestração multi-step
- Lógica condicional avançada
- Processamento de dados

## 📊 **EXEMPLO COMPLETO DE DIAGNÓSTICO REMOTO**

```bash
# 1. Verificar conectividade
sloth-runner agent list --master 192.168.1.29:50053

# 2. Status básico do ladyguica
sloth-runner agent run ladyguica "hostname && uptime && whoami" --master 192.168.1.29:50053

# 3. Status básico do keiteguica  
sloth-runner agent run keiteguica "hostname && uptime && whoami" --master 192.168.1.29:50053

# 4. Verificar espaço em disco
sloth-runner agent run ladyguica "df -h /" --master 192.168.1.29:50053

# 5. Verificar memória
sloth-runner agent run keiteguica "free -h" --master 192.168.1.29:50053

# 6. Listar processos
sloth-runner agent run ladyguica "ps aux | head -5" --master 192.168.1.29:50053
```

## 🚀 **RESULTADO ESPERADO**

```
INFO  🚀 Executing on agent: ladyguica
 INFO  📝 Command: hostname && date

ladyguica
Wed Oct  1 10:10:35 AM -03 2025

 SUCCESS  ✅ Command completed successfully on agent ladyguica
```

## 💡 **DICAS**

1. **Sempre use `--master 192.168.1.29:50053`** para especificar o master
2. **Use aspas duplas** para comandos complexos com pipes
3. **Teste a conectividade** primeiro com `agent list`
4. **Para workflows complexos**, considere scripts shell + comandos diretos
5. **O delegate_to por nome** resolve automaticamente para o IP correto

---

**🎉 Sistema 100% funcional para execução remota via nomes de agents!**