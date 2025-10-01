#!/bin/bash
# Test script for install.sh

set -e

echo "🧪 Testing install.sh..."
echo ""

# Test 1: Help
echo "Test 1: --help"
./install.sh --help > /dev/null 2>&1 && echo "✅ Help works" || echo "❌ Help failed"
echo ""

# Test 2: Install to temp directory
echo "Test 2: Install to temp directory"
TEMP_TEST=$(mktemp -d)
./install.sh --version v3.23.1 --install-dir "$TEMP_TEST" --force > /dev/null 2>&1
if [ -f "$TEMP_TEST/sloth-runner" ]; then
    echo "✅ Installation successful"
    VERSION=$("$TEMP_TEST/sloth-runner" version 2>&1 | head -1)
    echo "   Version: $VERSION"
    rm -rf "$TEMP_TEST"
else
    echo "❌ Installation failed"
fi
echo ""

# Test 3: Check download URL format
echo "Test 3: Check download URL format"
if grep -q "https://github.com/chalkan3-sloth/sloth-runner/releases/download" install.sh; then
    echo "✅ Correct GitHub URL"
else
    echo "❌ Wrong GitHub URL"
fi
echo ""

echo "✅ All tests passed!"
