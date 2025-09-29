#!/bin/bash

# Complete TaskDefinitions Removal Script
# This script removes ALL TaskDefinitions blocks from files

set -e

echo "🧹 COMPLETE TASKDEFINITIONS CLEANUP"
echo "===================================="
echo ""

cleaned_files=0

# Function to completely remove TaskDefinitions from a file
clean_taskdefinitions() {
    local file="$1"
    local filename=$(basename "$file")
    
    echo "🧹 Cleaning: $filename"
    
    if grep -q "TaskDefinitions" "$file"; then
        # Create backup if not already exists
        if [ ! -f "${file}.pre_modern_backup" ]; then
            cp "$file" "${file}.pre_modern_backup"
        fi
        
        # Use sed to remove TaskDefinitions blocks
        # Remove everything from "TaskDefinitions = {" to the matching closing brace
        # This handles multi-line TaskDefinitions blocks
        python3 << 'PYTHON' "$file"
import sys
import re

file_path = sys.argv[1]

with open(file_path, 'r') as f:
    content = f.read()

# Remove TaskDefinitions blocks - handle nested braces properly
def remove_taskdefinitions(text):
    lines = text.split('\n')
    result_lines = []
    in_taskdef = False
    brace_count = 0
    
    for line in lines:
        if not in_taskdef:
            # Check if this line starts a TaskDefinitions block
            if 'TaskDefinitions' in line and '=' in line and '{' in line:
                in_taskdef = True
                brace_count = line.count('{') - line.count('}')
                continue
            else:
                result_lines.append(line)
        else:
            # We're inside a TaskDefinitions block
            brace_count += line.count('{') - line.count('}')
            if brace_count <= 0:
                in_taskdef = False
                # Skip the closing line too
                continue
    
    return '\n'.join(result_lines)

# Remove TaskDefinitions
cleaned_content = remove_taskdefinitions(content)

# Remove empty lines at the end
cleaned_content = cleaned_content.rstrip() + '\n'

# Write back
with open(file_path, 'w') as f:
    f.write(cleaned_content)
PYTHON
        
        echo "  ✅ TaskDefinitions removed"
        ((cleaned_files++))
    else
        echo "  ✨ Already clean"
    fi
    echo ""
}

# Process all files that still contain TaskDefinitions
find examples/ -name "*.lua" ! -name "*.backup" ! -name "*.pre_modern_backup" -exec grep -l "TaskDefinitions" {} \; | while read -r file; do
    clean_taskdefinitions "$file"
done

echo "📊 CLEANUP SUMMARY"
echo "=================="
echo "🧹 Files cleaned: $cleaned_files"
echo "💾 Backup files: .pre_modern_backup"
echo ""
echo "✅ ALL TASKDEFINITIONS REMOVED!"
echo ""
echo "🔍 Verification:"
remaining=$(find examples/ -name "*.lua" ! -name "*.backup" ! -name "*.pre_modern_backup" -exec grep -l "TaskDefinitions" {} \; | wc -l)
echo "   Remaining TaskDefinitions: $remaining"

if [ "$remaining" -eq 0 ]; then
    echo "🎉 SUCCESS: No TaskDefinitions remain in active files!"
else
    echo "⚠️  Warning: $remaining files still contain TaskDefinitions"
fi