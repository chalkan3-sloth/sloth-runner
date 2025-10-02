#!/bin/bash

# Sloth Runner Installation Setup
# Adds $HOME/.local/bin to PATH if not already present

SLOTH_BIN="$HOME/.local/bin/sloth-runner"
LOCAL_BIN="$HOME/.local/bin"

echo "🛠️  SLOTH-RUNNER INSTALLATION SETUP"
echo "=================================="
echo ""

# Check if binary exists
if [ -f "$SLOTH_BIN" ]; then
    echo "✅ sloth-runner binary found at: $SLOTH_BIN"
    echo "   Size: $(du -h "$SLOTH_BIN" | cut -f1)"
else
    echo "❌ sloth-runner binary not found at: $SLOTH_BIN"
    echo "   Run: go build -o \$HOME/.local/bin/sloth-runner ./cmd/sloth-runner/"
    exit 1
fi

# Check if PATH includes $HOME/.local/bin
if echo "$PATH" | grep -q "$LOCAL_BIN"; then
    echo "✅ $LOCAL_BIN is in PATH"
else
    echo "⚠️  $LOCAL_BIN is NOT in PATH"
    echo ""
    echo "📝 To add it permanently, add this to your shell profile:"
    echo ""
    
    # Detect shell and provide appropriate instructions
    if [ -n "$ZSH_VERSION" ]; then
        echo "   # For Zsh (~/.zshrc):"
        echo "   export PATH=\"\$HOME/.local/bin:\$PATH\""
        echo ""
        echo "   # Run this command to add it now:"
        echo "   echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.zshrc"
        echo "   source ~/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        echo "   # For Bash (~/.bashrc or ~/.bash_profile):"
        echo "   export PATH=\"\$HOME/.local/bin:\$PATH\""
        echo ""
        echo "   # Run this command to add it now:"
        echo "   echo 'export PATH=\"\$HOME/.local/bin:\$PATH\"' >> ~/.bashrc"
        echo "   source ~/.bashrc"
    else
        echo "   # For your shell configuration file:"
        echo "   export PATH=\"\$HOME/.local/bin:\$PATH\""
    fi
    
    echo ""
    echo "   # Or add it temporarily for this session:"
    echo "   export PATH=\"\$HOME/.local/bin:\$PATH\""
fi

echo ""
echo "🧪 TESTING INSTALLATION:"
echo "---"

# Test with full path
echo "Testing with full path:"
if "$SLOTH_BIN" version 2>/dev/null; then
    echo "✅ Binary works correctly"
else
    echo "❌ Binary test failed"
fi

echo ""

# Test if available in PATH
echo "Testing if available in PATH:"
if command -v sloth-runner >/dev/null 2>&1; then
    echo "✅ sloth-runner is available in PATH"
    echo "   Location: $(which sloth-runner)"
    echo "   Version: $(sloth-runner version 2>/dev/null | head -1)"
else
    echo "⚠️  sloth-runner not found in PATH"
    echo "   Use full path: $SLOTH_BIN"
fi

echo ""
echo "🔧 QUICK COMMANDS:"
echo "---"
echo "• Test installation:     $SLOTH_BIN version"
echo "• List agents:          $SLOTH_BIN agent list --master 192.168.1.29:50053"
echo "• Run workflow:         $SLOTH_BIN run -f examples/agents/ls_delegate_simple.sloth"
echo "• Start master:         $SLOTH_BIN master --port 50053 --daemon"
echo ""
echo "📁 EXAMPLE FILES AVAILABLE:"
echo "• $(pwd)/examples/agents/ls_delegate_simple.sloth"
echo "• $(pwd)/examples/agents/demo_remote_execution.sh"
echo "• $(pwd)/examples/agents/README_SQLITE.md"
echo ""
echo "✅ INSTALLATION COMPLETE!"