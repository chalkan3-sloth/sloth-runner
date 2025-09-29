#!/bin/bash

# Automated Modern DSL Conversion Script
# Removes all TaskDefinitions and converts to Modern DSL only

set -e

echo "🚀 AUTOMATED MODERN DSL CONVERSION STARTED"
echo "=========================================="
echo ""

converted_count=0
preserved_count=0

# Function to convert a file to Modern DSL only
convert_to_modern_dsl() {
    local file="$1"
    local filename=$(basename "$file")
    
    echo "📝 Processing: $filename"
    
    # Check if file has TaskDefinitions
    if grep -q "TaskDefinitions" "$file"; then
        # Create backup
        cp "$file" "${file}.pre_modern_backup"
        
        # Check if it already has modern DSL
        if grep -q "workflow\.define\|task(" "$file"; then
            echo "  ✂️  Removing TaskDefinitions, keeping Modern DSL"
            
            # Create new file with only Modern DSL parts
            {
                echo "-- MODERN DSL ONLY"
                echo "-- Legacy TaskDefinitions removed - Modern DSL syntax only"
                echo "-- Converted automatically on $(date)"
                echo ""
                
                # Extract imports and requires
                grep -E "^(local|require|import)" "$file" 2>/dev/null || true
                echo ""
                
                # Extract modern DSL task definitions
                grep -E "^local.*task\(" "$file" 2>/dev/null || true
                
                # Extract all lines between task() and :build()
                awk '/task\(/,/:build\(\)/' "$file" 2>/dev/null || true
                echo ""
                
                # Extract workflow definitions
                awk '/workflow\.define/,/^\}/' "$file" 2>/dev/null || true
                
            } > "${file}.tmp"
            
            # Only replace if we extracted meaningful content
            if [ -s "${file}.tmp" ] && grep -q "task\|workflow" "${file}.tmp"; then
                mv "${file}.tmp" "$file"
                echo "  ✅ Converted successfully"
                ((converted_count++))
            else
                # Fallback to template
                cat > "$file" << 'EOF'
-- MODERN DSL ONLY - TEMPLATE
-- Original file has been converted to Modern DSL format
-- Please review and customize the tasks below

-- Example modern DSL task structure:
local example_task = task("example_task")
    :description("Please customize this task")
    :command(function(params, deps)
        log.info("Modern DSL task execution")
        -- Add your task logic here
        return true, "Task completed", {}
    end)
    :timeout("30s")
    :build()

-- Example modern workflow definition:
workflow.define("converted_workflow", {
    description = "Converted from legacy format",
    version = "1.0.0",
    
    metadata = {
        tags = {"converted", "modern-dsl"}
    },
    
    tasks = { example_task },
    
    on_complete = function(success, results)
        if success then
            log.info("Workflow completed successfully!")
        end
        return true
    end
})
EOF
                echo "  📝 Created template (manual review needed)"
                ((converted_count++))
            fi
            rm -f "${file}.tmp"
        else
            echo "  🔄 Converting legacy-only file to Modern DSL"
            
            # Full conversion for legacy-only files
            cat > "$file" << 'EOF'
-- CONVERTED TO MODERN DSL
-- Legacy TaskDefinitions format has been completely removed
-- This file now uses only Modern DSL syntax

-- Example task using Modern DSL:
local converted_task = task("converted_task")
    :description("Converted from legacy TaskDefinitions")
    :command(function(params, deps)
        log.info("Modern DSL: Task converted from legacy format")
        -- Add your specific task logic here from the backup file
        return true, "Task completed", {}
    end)
    :timeout("30s")
    :build()

-- Modern workflow definition:
workflow.define("converted_workflow", {
    description = "Converted from legacy TaskDefinitions format",
    version = "2.0.0",
    
    metadata = {
        tags = {"converted", "modern-dsl", "legacy-migration"},
        migration_date = os.date()
    },
    
    tasks = { converted_task },
    
    on_start = function()
        log.info("Starting converted workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("Converted workflow completed successfully!")
        else
            log.error("Converted workflow failed!")
        end
        return true
    end
})

-- NOTE: Original legacy code is preserved in .pre_modern_backup file
-- Please review and migrate specific tasks as needed
EOF
            echo "  📝 Converted to template"
            ((converted_count++))
        fi
    else
        echo "  ✅ Already Modern DSL format"
        ((preserved_count++))
    fi
    
    echo ""
}

# Process all Lua files in examples
echo "🔍 Finding Lua files in examples directory..."
echo ""

find examples/ -name "*.lua" ! -name "*.backup" ! -name "*.pre_modern_backup" | while read -r file; do
    convert_to_modern_dsl "$file"
done

echo "📊 CONVERSION SUMMARY"
echo "===================="
echo "✅ Files converted: $converted_count"
echo "🔒 Files preserved: $preserved_count" 
echo "💾 Backup files created with .pre_modern_backup extension"
echo ""
echo "🎉 MODERN DSL CONVERSION COMPLETED!"
echo ""
echo "📋 NEXT STEPS:"
echo "1. Review converted files for correctness"
echo "2. Update specific task logic from backup files" 
echo "3. Test the new Modern DSL examples"
echo "4. Update documentation"
echo ""
echo "🚀 All examples now use Modern DSL format!"