# NixOS Configuration Workflow

Este workflow demonstra como transformar uma configuração NixOS tradicional do GitHub em um workflow Sloth totalmente automatizado e declarativo.

## 📦 Baseado em

Repositório original: [chalkan3/nixos-config](https://github.com/chalkan3/nixos-config)

Este workflow replica completamente a configuração desse repositório usando as funções avançadas do módulo NixOS do Sloth Runner.

## 🚀 Como Usar

### Pré-requisitos

1. Sloth Runner instalado e configurado
2. Agent rodando no host NixOS alvo
3. Acesso SSH ao host (se remoto)

### Execução Básica

```bash
# Executar no host local
sloth-runner run nixos_complete_setup \
    --file nixos_config_complete_setup.sloth \
    --yes

# Executar em host remoto via agent
sloth-runner run nixos_complete_setup \
    --file nixos_config_complete_setup.sloth \
    --delegate-to my-nixos-host \
    --yes
```

### Modo Dry-Run (Validação)

Para apenas validar a configuração sem aplicar:

```bash
# Comentar a task 'apply_configuration' no workflow
# Ou executar apenas até a validação
sloth-runner run nixos_complete_setup \
    --file nixos_config_complete_setup.sloth \
    --delegate-to my-nixos-host \
    --yes
```

## 📋 O Que Este Workflow Faz

### 1. **Boot & Hardware** (Task 1)
- Configura systemd-boot como bootloader
- Habilita modificação de variáveis EFI
- Define parâmetros do kernel (quiet, splash)
- Configura tmpfs

### 2. **Sistema** (Task 2)
- Define hostname: `nixos-qemu`
- Timezone: `America/Sao_Paulo`
- Locale: `en_US.UTF-8`
- Keymap do console: `us`

### 3. **Usuários** (Task 3)
Cria 3 usuários:

| Usuário   | Grupos                 | Shell | Senha Inicial    |
|-----------|------------------------|-------|------------------|
| chalkan3  | wheel, networkmanager  | zsh   | changeme123!     |
| nixos     | wheel, networkmanager  | zsh   | changeme123!     |
| root      | -                      | -     | root             |

⚠️ **IMPORTANTE**: Mude essas senhas padrão após a primeira execução!

### 4. **Pacotes** (Task 4)
Instala e configura:

**Ferramentas Core:**
- `vim`, `neovim`
- `wget`, `curl`, `git`

**Shell & Terminal:**
- `zsh` (com completion)
- `kitty.terminfo`

**Utilitários:**
- `btop` - Monitor de sistema
- `lsd` - ls melhorado
- `fzf` - Fuzzy finder
- `gh` - GitHub CLI

**Configurações:**
- EDITOR=nvim
- VISUAL=nvim
- ZSH completion habilitado
- nix-ld habilitado

### 5. **Networking** (Task 5)
- NetworkManager habilitado
- Firewall ativo
- Portas TCP liberadas: 22 (SSH), 50051

### 6. **Serviços** (Task 6)
- **OpenSSH**:
  - Habilitado
  - Root login permitido
  - Autenticação por senha habilitada
  - Firewall aberto automaticamente

- **QEMU Guest Agent**:
  - Habilitado (útil para VMs)

### 7. **Virtualização** (Task 7 - Opcional)
- libvirt/KVM habilitado
- QEMU OVMF (UEFI) suportado
- QEMU não roda como root

### 8. **Performance** (Task 8)
- CPU Governor: `performance`
- zram swap: 50% da RAM
- Kernel params para melhor performance:
  - `mitigations=off` (menos seguro, mais rápido)
  - `nowatchdog`

### 9. **Backup** (Task 9)
- Cria backup timestamped da configuração atual
- Formato: `/etc/nixos/configuration.nix.backup-YYYYMMDD-HHMMSS`

### 10. **Validação** (Task 10)
- Executa `nixos-rebuild dry-build`
- Valida configuração sem aplicar
- Para execução se houver erros

### 11. **Aplicação** (Task 11)
- Executa `nixos-rebuild switch`
- Aplica a configuração
- Ativa a nova geração

### 12. **Verificação** (Task 12)
- Lista gerações do sistema
- Verifica estado final

## 🔧 Customização

### Modificar Timezone

```lua
local ok, msg = nixos.configure_system({
    timezone = "America/New_York",  -- Mudar aqui
    locale = "en_US.UTF-8",
    -- ...
})
```

### Adicionar Mais Pacotes

```lua
local ok, msg = nixos.configure_environment({
    packages = {
        -- Pacotes existentes...
        "htop",
        "tree",
        "ncdu",
        -- Seus pacotes aqui
    },
    -- ...
})
```

### Configurar SSH Keys

```lua
local ok, msg = nixos.configure_user({
    username = "chalkan3",
    -- ... outras configs
    ssh_keys = {
        "ssh-ed25519 AAAAC3Nza... seu-email@example.com",
        "ssh-rsa AAAAB3NzaC... outro-email@example.com"
    }
})
```

### Ajustar Firewall

```lua
local ok, msg = nixos.configure_networking({
    firewall = {
        enable = true,
        tcp_ports = {22, 80, 443, 8080},  -- Adicionar portas
        udp_ports = {53, 123}
    }
})
```

### Hardening do SSH

```lua
local ok, msg = nixos.configure_service({
    service = "openssh",
    enable = true,
    settings = {
        permitRootLogin = "no",              -- Desabilitar root
        passwordAuthentication = false,      -- Apenas chaves SSH
        kbdInteractiveAuthentication = false
    }
})
```

## 🛡️ Segurança

### ⚠️ Avisos de Segurança

Este workflow usa configurações **NÃO SEGURAS** para facilitar setup inicial:

1. **Senhas Padrão**: Todas as senhas são padrão e conhecidas
2. **Root Login SSH**: Permitido (não recomendado em produção)
3. **Password Auth SSH**: Habilitado (chaves SSH são mais seguras)
4. **Kernel Mitigations Off**: Desabilita proteções de segurança

### ✅ Recomendações para Produção

Após a primeira execução:

1. **Mudar todas as senhas**:
```bash
passwd chalkan3
passwd nixos
sudo passwd root
```

2. **Adicionar chaves SSH**:
```bash
ssh-copy-id chalkan3@nixos-qemu
```

3. **Hardening SSH** - Modificar o workflow:
```lua
settings = {
    permitRootLogin = "no",
    passwordAuthentication = false,
    kbdInteractiveAuthentication = false
}
```

4. **Remover senhas plaintext** - Usar hashed passwords:
```lua
hashed_password = "$6$rounds=4096$saltsaltsa$hash..."
```

5. **Habilitar kernel mitigations**:
```lua
kernel_params = {
    "quiet",
    "splash"
    -- Remover "mitigations=off"
}
```

## 🔄 Gerenciamento de Gerações

### Listar Gerações

```bash
# Via Sloth
nixos.list_generations({})

# Diretamente
nix-env --list-generations -p /nix/var/nix/profiles/system
```

### Rollback

```bash
# Via Sloth
nixos.rollback({use_sudo = true})

# Diretamente
sudo nixos-rebuild switch --rollback
```

### Mudar para Geração Específica

```lua
nixos.switch_generation({
    generation = 42,
    use_sudo = true
})
```

### Limpar Gerações Antigas

```lua
nixos.delete_generations({
    older_than = "30d",  -- Mais de 30 dias
    use_sudo = true
})
```

## 📊 Workflow Features

### Idempotência

Todas as funções são idempotentes - você pode executar o workflow múltiplas vezes com segurança.

### Error Handling

- Validação antes de aplicar
- Backup automático
- Continue-on-error desabilitado (para)
- Mensagens de erro claras

### Logging

- Emojis para facilitar leitura
- Logs estruturados
- Output colorido (via pterm)

### Timeout

- Workflow timeout: 30 minutos
- Adequado para rebuilds grandes

## 🎯 Use Cases

### 1. Setup Inicial de VM

```bash
sloth-runner run nixos_complete_setup \
    --file nixos_config_complete_setup.sloth \
    --delegate-to new-vm \
    --yes
```

### 2. Replicar Configuração

Clone este workflow e customize para criar configurações padronizadas para múltiplas máquinas.

### 3. Disaster Recovery

Use o workflow para reconstruir um sistema NixOS do zero rapidamente.

### 4. CI/CD Testing

Teste configurações NixOS em VMs efêmeras antes de aplicar em produção.

## 🔍 Troubleshooting

### Erro de Validação

Se a validação falhar:

```bash
# Check manualmente
sudo nixos-rebuild dry-build

# Ver erros detalhados
journalctl -xe
```

### Rebuild Falha

Se o rebuild falhar:

```bash
# Rollback imediato
sudo nixos-rebuild switch --rollback

# Ver logs
journalctl -u nixos-rebuild
```

### SSH Não Funciona

Verifique:

1. Firewall: porta 22 liberada?
2. Serviço SSH: `systemctl status sshd`
3. Usuário tem permissão: está no grupo `wheel`?

### NetworkManager Não Inicia

```bash
# Verificar conflitos
systemctl status NetworkManager

# Ver configuração
cat /etc/nixos/configuration.nix | grep -A5 networking
```

## 📚 Referências

- [Repositório Original](https://github.com/chalkan3/nixos-config)
- [NixOS Manual](https://nixos.org/manual/nixos/stable/)
- [Sloth Runner Documentation](https://github.com/chalkan3-sloth/sloth-runner)
- [NixOS Options Search](https://search.nixos.org/options)

## 🤝 Contribuindo

Este é um exemplo demonstrativo. Para melhorias:

1. Abra issue no repositório
2. Faça fork e customize
3. Compartilhe suas variações

## 📄 Licença

Este workflow é baseado na configuração pública de [chalkan3/nixos-config](https://github.com/chalkan3/nixos-config).

---

**Criado com ❤️ usando Sloth Runner**
