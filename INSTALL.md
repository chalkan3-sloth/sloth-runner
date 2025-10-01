# 🦥 Sloth Runner - Installation Guide

Complete guide to install Sloth Runner on any platform.

## 📋 Table of Contents

- [Quick Install](#-quick-install)
- [Installation Methods](#-installation-methods)
- [Platform-Specific Instructions](#-platform-specific-instructions)
- [Verification](#-verification)
- [Configuration](#-configuration)
- [Troubleshooting](#-troubleshooting)
- [Uninstallation](#-uninstallation)

## 🚀 Quick Install

### One-Line Installation

```bash
curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh | bash
```

This will:
- ✅ Detect your platform automatically
- ✅ Download the latest version
- ✅ Install to `/usr/local/bin` (with sudo) or `$HOME/.local/bin` (without)
- ✅ Verify the installation

### Manual Download

Visit the [releases page](https://github.com/chalkan3-sloth/sloth-runner/releases) and download the appropriate binary for your platform.

## 📦 Installation Methods

### Method 1: Automated Script (Recommended)

The installation script supports multiple options:

```bash
# Install latest version
bash install.sh

# Install specific version
bash install.sh --version v3.23.1

# Install without sudo (user directory)
bash install.sh --no-sudo

# Install to custom directory
bash install.sh --install-dir /opt/bin

# Force overwrite existing installation
bash install.sh --force

# Show help
bash install.sh --help
```

### Method 2: GitHub CLI

If you have GitHub CLI installed:

```bash
gh release download --repo chalkan3-sloth/sloth-runner --pattern '*linux_amd64.tar.gz'
tar -xzf sloth-runner_*_linux_amd64.tar.gz
sudo mv sloth-runner /usr/local/bin/
```

### Method 3: Direct Download

```bash
# Set version and platform
VERSION="v3.23.1"
OS="linux"  # or darwin
ARCH="amd64"  # or arm64

# Download
curl -LO "https://github.com/chalkan3-sloth/sloth-runner/releases/download/${VERSION}/sloth-runner_${VERSION}_${OS}_${ARCH}.tar.gz"

# Extract
tar -xzf sloth-runner_${VERSION}_${OS}_${ARCH}.tar.gz

# Install
sudo mv sloth-runner /usr/local/bin/
chmod +x /usr/local/bin/sloth-runner
```

### Method 4: Build from Source

```bash
# Clone repository
git clone https://github.com/chalkan3-sloth/sloth-runner.git
cd sloth-runner

# Build
go build -o sloth-runner ./cmd/sloth-runner

# Install
sudo mv sloth-runner /usr/local/bin/
```

## 💻 Platform-Specific Instructions

### Linux

#### Ubuntu/Debian

```bash
# Using install script
curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh | bash

# Verify
sloth-runner --version
```

#### CentOS/RHEL/Fedora

```bash
# Using install script
curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh | bash

# Or manually
sudo yum install -y tar gzip
curl -LO https://github.com/chalkan3-sloth/sloth-runner/releases/latest/download/sloth-runner_v*_linux_amd64.tar.gz
tar -xzf sloth-runner_*.tar.gz
sudo mv sloth-runner /usr/local/bin/
```

#### Arch Linux

```bash
# Using install script
curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh | bash
```

### macOS

#### Intel (x86_64)

```bash
# Using install script
curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh | bash

# Or using Homebrew (if available)
# brew install chalkan3-sloth/tap/sloth-runner
```

#### Apple Silicon (ARM64)

```bash
# Using install script
curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh | bash
```

### ARM/Raspberry Pi

```bash
# ARM64
curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh | bash

# The script will automatically detect ARM64 architecture
```

## ✅ Verification

After installation, verify that Sloth Runner is working:

```bash
# Check version
sloth-runner --version

# Expected output:
# sloth-runner version 3.23.1
# Git commit: 9b8b45da99d547712a100eb28e48d02d7a3ba041
# Build date: 2025-10-01T19:26:34Z

# Check help
sloth-runner --help

# Run a simple workflow
sloth-runner run -f examples/hello_world.sloth
```

### Verify Installation Path

```bash
# Check where sloth-runner is installed
which sloth-runner

# Check if it's in PATH
echo $PATH | grep -o '/usr/local/bin\|$HOME/.local/bin'
```

## ⚙️ Configuration

### PATH Configuration

If the binary is not in your PATH, add it:

#### Bash (~/.bashrc)

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

#### Zsh (~/.zshrc)

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

#### Fish (~/.config/fish/config.fish)

```fish
set -Ua fish_user_paths $HOME/.local/bin
```

### Shell Completion

Enable shell completion for better experience:

```bash
# Bash
sloth-runner completion bash > /etc/bash_completion.d/sloth-runner

# Zsh
sloth-runner completion zsh > "${fpath[1]}/_sloth-runner"

# Fish
sloth-runner completion fish > ~/.config/fish/completions/sloth-runner.fish
```

## 🔧 Troubleshooting

### Binary Not Found After Installation

**Problem:** `command not found: sloth-runner`

**Solution:**
```bash
# Check if binary exists
ls -l $HOME/.local/bin/sloth-runner
ls -l /usr/local/bin/sloth-runner

# Add to PATH
export PATH="$HOME/.local/bin:$PATH"

# Or move to a directory in PATH
sudo mv $HOME/.local/bin/sloth-runner /usr/local/bin/
```

### Permission Denied

**Problem:** `Permission denied` when running

**Solution:**
```bash
# Make executable
chmod +x $(which sloth-runner)
```

### Download Fails

**Problem:** Cannot download from GitHub

**Solution:**
```bash
# Check internet connection
curl -I https://github.com

# Use alternative method (direct download)
wget https://github.com/chalkan3-sloth/sloth-runner/releases/latest/download/sloth-runner_v3.23.1_linux_amd64.tar.gz

# Or use gh CLI
gh release download --repo chalkan3-sloth/sloth-runner
```

### Architecture Not Supported

**Problem:** `Unsupported architecture`

**Solution:**
```bash
# Check your architecture
uname -m

# Supported architectures:
# - x86_64 (amd64)
# - aarch64/arm64 (ARM 64-bit)

# Build from source if needed
git clone https://github.com/chalkan3-sloth/sloth-runner.git
cd sloth-runner
go build -o sloth-runner ./cmd/sloth-runner
```

### macOS Gatekeeper Warning

**Problem:** "sloth-runner cannot be opened because the developer cannot be verified"

**Solution:**
```bash
# Remove quarantine attribute
xattr -d com.apple.quarantine $(which sloth-runner)

# Or allow in System Preferences
# System Preferences > Security & Privacy > General > Allow
```

## 🗑️ Uninstallation

### Remove Binary

```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/sloth-runner

# If installed to ~/.local/bin
rm ~/.local/bin/sloth-runner
```

### Remove Configuration (Optional)

```bash
# Remove configuration files
rm -rf ~/.sloth-runner

# Remove cache
rm -rf ~/.cache/sloth-runner

# Remove logs
rm -rf ~/.local/share/sloth-runner
```

### Remove Shell Completion

```bash
# Bash
sudo rm /etc/bash_completion.d/sloth-runner

# Zsh
rm "${fpath[1]}/_sloth-runner"

# Fish
rm ~/.config/fish/completions/sloth-runner.fish
```

## 📚 Next Steps

After installation:

1. **Read the Quick Start Guide**
   ```bash
   sloth-runner docs quickstart
   ```

2. **Try Examples**
   ```bash
   cd examples
   sloth-runner run -f hello_world.sloth
   ```

3. **Start Master Server**
   ```bash
   sloth-runner master start --port 50053 --daemon
   ```

4. **Configure Agents**
   ```bash
   sloth-runner agent start --name myagent --master localhost:50053
   ```

5. **Explore Documentation**
   - [Online Docs](https://chalkan3.github.io/sloth-runner/)
   - [Examples](https://github.com/chalkan3-sloth/sloth-runner/tree/master/examples)
   - [Modern DSL Guide](https://chalkan3.github.io/sloth-runner/en/modern-dsl/overview/)

## 🆘 Getting Help

- 📖 [Documentation](https://chalkan3.github.io/sloth-runner/)
- 💬 [GitHub Discussions](https://github.com/chalkan3-sloth/sloth-runner/discussions)
- 🐛 [Report Issues](https://github.com/chalkan3-sloth/sloth-runner/issues)
- 📧 Email: support@sloth-runner.dev

## 🔄 Updating

To update to the latest version:

```bash
# Using install script
curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/install.sh | bash --force

# Or manually
bash install.sh --force
```

## 📊 System Requirements

- **OS**: Linux, macOS
- **Architecture**: x86_64 (amd64), ARM64
- **Disk Space**: ~20-30 MB
- **Memory**: Minimum 64 MB RAM
- **Dependencies**: None (statically compiled)

## 🎉 Installation Complete!

You're ready to use Sloth Runner! 🚀

For your first workflow:
```bash
echo 'workflow.define("hello")
  :tasks({
    task("greet")
      :command(function()
        log.info("Hello from Sloth Runner!")
        return true
      end)
      :build()
  })' > hello.sloth

sloth-runner run -f hello.sloth
```

Happy automating! 🦥
