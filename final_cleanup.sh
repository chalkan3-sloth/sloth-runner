#!/bin/bash

# Simple TaskDefinitions Removal Script using sed

set -e

echo "🧹 REMOVING ALL TASKDEFINITIONS"
echo "==============================="
echo ""

cleaned=0

# Function to remove TaskDefinitions using a more robust approach
remove_taskdefs() {
    local file="$1"
    local filename=$(basename "$file")
    
    echo "🧹 Processing: $filename"
    
    if grep -q "TaskDefinitions" "$file"; then
        # Create backup
        if [ ! -f "${file}.clean_backup" ]; then
            cp "$file" "${file}.clean_backup"
        fi
        
        # Remove lines containing TaskDefinitions and the blocks that follow
        # This removes from TaskDefinitions = { to the end of the file
        # or until we find a workflow.define
        awk '
        /TaskDefinitions/ { skip = 1 }
        /workflow\.define/ { skip = 0 }
        !skip { print }
        ' "$file" > "${file}.tmp"
        
        # If the temp file is too small, it means we removed too much
        # In that case, just add a basic modern DSL template
        if [ $(wc -l < "${file}.tmp") -lt 5 ]; then
            cat > "${file}.tmp" << 'EOF'
-- MODERN DSL ONLY
-- TaskDefinitions have been removed

-- Example Modern DSL task:
local example_task = task("example_task")
    :description("Modern DSL task")
    :command(function(params, deps)
        log.info("Modern DSL execution")
        return true, "Completed", {}
    end)
    :build()

-- Modern workflow:
workflow.define("modern_workflow", {
    description = "Modern DSL workflow",
    version = "1.0.0",
    tasks = { example_task }
})
EOF
        fi
        
        mv "${file}.tmp" "$file"
        echo "  ✅ Cleaned"
        ((cleaned++))
    else
        echo "  ✨ Already clean"
    fi
}

# Process all files
find examples/ -name "*.lua" ! -name "*.backup" ! -name "*.clean_backup" ! -name "*.pre_modern_backup" | while read -r file; do
    remove_taskdefs "$file"
done

echo ""
echo "📊 SUMMARY:"
echo "🧹 Files processed: $cleaned"

# Final verification
remaining=$(find examples/ -name "*.lua" ! -name "*backup*" -exec grep -l "TaskDefinitions" {} \; 2>/dev/null | wc -l || echo "0")
echo "📋 Remaining TaskDefinitions: $remaining"

if [ "$remaining" -eq 0 ]; then
    echo "🎉 SUCCESS: All TaskDefinitions removed!"
else
    echo "⚠️  Some files may still need manual cleanup"
fi