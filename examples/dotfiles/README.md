# Dotfiles Management Examples

This directory contains examples of using the `stow` module for managing dotfiles with Sloth Runner.

## What is GNU Stow?

GNU Stow is a symlink farm manager. It helps you manage your dotfiles by creating symbolic links from a centralized dotfiles directory to your home directory.

## Examples

### 1. simple_stow.sloth

Quick start example showing the basics:
- Stowing a single package
- Stowing multiple packages
- Removing configurations
- Listing stowed packages
- Verifying package integrity
- Remote deployment

**Run it:**
```bash
sloth-runner run my-dotfiles --file examples/dotfiles/simple_stow.sloth setup-vim
```

### 2. dotfiles_management.sloth

Complete production-ready example with:
- Idempotent deployment
- Conflict detection and adoption
- Comprehensive verification
- Safe configuration rotation with rollback
- Parallel deployment to multiple servers
- Audit reporting
- Cleanup operations

**Run it:**
```bash
# Deploy locally
sloth-runner run dotfiles-production --file examples/dotfiles/dotfiles_management.sloth setup-local-dotfiles

# Audit current state
sloth-runner run dotfiles-production --file examples/dotfiles/dotfiles_management.sloth audit-dotfiles

# Deploy to team servers
sloth-runner run dotfiles-production --file examples/dotfiles/dotfiles_management.sloth deploy-team-dotfiles
```

## Prerequisites

### Install GNU Stow

**On macOS:**
```bash
brew install stow
```

**On Arch Linux:**
```bash
sudo pacman -S stow
```

**On Ubuntu/Debian:**
```bash
sudo apt install stow
```

### Set up your dotfiles directory

1. Create a dotfiles directory:
```bash
mkdir -p ~/.dotfiles
```

2. Organize your configs as packages:
```bash
~/.dotfiles/
├── vim/
│   └── .vimrc
├── zsh/
│   ├── .zshrc
│   └── .zshenv
├── tmux/
│   └── .tmux.conf
└── git/
    └── .gitconfig
```

## Features Demonstrated

### Idempotency
The `stow` module tracks state and only makes changes when necessary:
```lua
-- Safe to run multiple times
stow.stow({ package = "vim" })
stow.stow({ package = "vim" })  -- No action, already stowed
```

### Verification
Check integrity of stowed configurations:
```lua
local verify = stow.verify({ package = "vim" })
if verify.is_valid and verify.is_complete then
    print("✓ All configs properly linked")
end
```

### Conflict Handling
Adopt existing files or detect conflicts:
```lua
local check = stow.check({ package = "vim" })
if not check.would_succeed then
    -- Conflicts detected, adopt existing files
    stow.adopt({ package = "vim" })
end
```

### Remote Deployment
Deploy configurations to remote servers:
```lua
stow.stow({
    package = "vim",
    delegate_to = "my-server"
})
```

### Parallel Execution
Deploy to multiple servers in parallel:
```lua
local servers = {"web-01", "web-02", "web-03"}
goroutine.map(servers, function(server)
    stow.stow({
        package = "vim",
        delegate_to = server
    })
end)
```

## Best Practices

1. **Version Control**: Keep your dotfiles in git
2. **Test Locally**: Always test configurations locally first
3. **Use Verification**: Always verify after stowing critical configs
4. **Backup**: Keep backups before major changes
5. **Incremental**: Deploy one package at a time initially
6. **Use State**: Run with a stack for state tracking and rollback

## Troubleshooting

### Stow Conflicts
If you get conflicts, you can:
- Use `adopt()` to move existing files into your dotfiles
- Remove conflicting files manually
- Use `check()` to see what would happen first

### Broken Symlinks
Use `verify()` to detect broken symlinks:
```lua
local verify = stow.verify({ package = "vim" })
for _, broken in ipairs(verify.broken_links) do
    print("Broken: " .. broken)
end
```

### State Issues
If state gets out of sync:
```bash
# View current state
sloth-runner state list dotfiles-production

# Clean state if needed
sloth-runner state clean dotfiles-production
```

## See Also

- [Stow Module Documentation](../../docs/modules/stow.md)
- [GNU Stow Manual](https://www.gnu.org/software/stow/manual/)
- [State Management Guide](../../docs/state-management.md)
