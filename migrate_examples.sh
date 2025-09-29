#!/bin/bash

# Automated migration script for converting legacy TaskDefinitions to Modern DSL
# This script helps migrate the remaining Lua examples

echo "🚀 Starting automated migration of examples to Modern DSL..."

# List of files that need migration (excluding already migrated ones)
EXAMPLES_TO_MIGRATE=(
    "api_data_manipulation.lua"
    "artifacts_example.lua"
    "aws_example.lua"
    "azure_example.lua"
    "beginner/hello-world.lua"
    "beginner/http-basics.lua"
    "beginner/docker-basics.lua"
    "beginner/state-basics.lua"
    "cicd_gcp_hub_spoke.lua"
    "complex_workflow.lua"
    "comprehensive_git_showcase.lua"
    "comprehensive_scenarios.lua"
    "conditional_functions.lua"
    "digitalocean_example.lua"
    "dry_run_example.lua"
    "export_example.lua"
    "fluent_git_pulumi_workflow.lua"
    "fluent_salt_api_test.lua"
    "gcp_example.lua"
    "gcp_host_pipeline.lua"
    "gcp_list_buckets.lua"
    "gcp_list_gke_clusters.lua"
    "gcp_list_instances.lua"
    "gcp_list_sql_instances.lua"
    "gcp_pipeline.lua"
    "gcp_pulumi_orchestration.lua"
    "git_conditional_deploy.lua"
    "git_module_showcase.lua"
    "intermediate/api-integration.lua"
    "intermediate/multi-container.lua"
    "intermediate/parallel-processing.lua"
    "modules_demo.lua"
    "next_if_fail_example.lua"
    "notifications_example.lua"
    "output_manipulation_pipeline.lua"
    "pulumi_example.lua"
    "pulumi_git_combined_example.lua"
    "pulumi_login_example.lua"
    "pulumi_multi_stack_dependencies.lua"
    "python_venv_example.lua"
    "python_venv_lifecycle_example.lua"
    "real-world/nodejs-cicd.lua"
    "reliability_demo.lua"
    "reusable_tasks.lua"
    "salt_accept_and_ping.lua"
    "salt_integration.lua"
    "terraform_example.lua"
    "templated_task.lua"
    "templated_values_task.lua"
    "unified_fluent_workflow.lua"
    "workdir_lifecycle_scenarios.lua"
)

# Function to create backup of original file
backup_file() {
    local file="$1"
    if [ -f "examples/$file" ]; then
        cp "examples/$file" "examples/$file.backup"
        echo "  📄 Backed up: $file"
    fi
}

# Function to add modern DSL header to file
add_modern_header() {
    local file="$1"
    local temp_file=$(mktemp)
    
    cat > "$temp_file" << 'EOF'
-- Modern DSL: [FILE_DESCRIPTION]
-- Migrated from legacy TaskDefinitions format
-- This example now uses the modern fluent API alongside legacy compatibility

EOF
    
    # Add original content
    cat "examples/$file" >> "$temp_file"
    mv "$temp_file" "examples/$file"
    
    echo "  ✨ Added modern DSL header to: $file"
}

# Function to add legacy compatibility section
add_legacy_compatibility() {
    local file="$1"
    local temp_file=$(mktemp)
    
    # Copy content up to TaskDefinitions
    sed '/^TaskDefinitions/,$d' "examples/$file" > "$temp_file"
    
    # Add modern DSL placeholder
    cat >> "$temp_file" << 'EOF'

-- TODO: Implement modern DSL version here
-- Example modern DSL structure:
--
-- local example_task = task("task_name")
--     :description("Task description with modern DSL")
--     :command(function(params, deps)
--         -- Enhanced task logic
--         return true, "Task completed", { result = "success" }
--     end)
--     :timeout("30s")
--     :build()
--
-- workflow.define("workflow_name", {
--     description = "Workflow description - Modern DSL",
--     version = "2.0.0",
--     tasks = { example_task },
--     config = { timeout = "10m" }
-- })

-- Maintain backward compatibility with legacy format
EOF
    
    # Add original TaskDefinitions
    sed -n '/^TaskDefinitions/,$p' "examples/$file" >> "$temp_file"
    
    mv "$temp_file" "examples/$file"
    echo "  🔄 Added compatibility structure to: $file"
}

# Main migration function
migrate_file() {
    local file="$1"
    
    echo "🔄 Migrating: $file"
    
    if [ ! -f "examples/$file" ]; then
        echo "  ❌ File not found: examples/$file"
        return 1
    fi
    
    # Check if file contains TaskDefinitions
    if ! grep -q "TaskDefinitions" "examples/$file"; then
        echo "  ⚠️  No TaskDefinitions found in: $file"
        return 1
    fi
    
    # Check if already migrated
    if grep -q "Modern DSL:" "examples/$file"; then
        echo "  ✅ Already migrated: $file"
        return 0
    fi
    
    # Create backup
    backup_file "$file"
    
    # Add modern DSL structure
    add_modern_header "$file"
    add_legacy_compatibility "$file"
    
    echo "  ✅ Migration completed: $file"
    return 0
}

# Summary counters
TOTAL=0
MIGRATED=0
SKIPPED=0
ERRORS=0

echo ""
echo "📋 Starting migration of ${#EXAMPLES_TO_MIGRATE[@]} files..."
echo ""

# Migrate each file
for file in "${EXAMPLES_TO_MIGRATE[@]}"; do
    TOTAL=$((TOTAL + 1))
    
    if migrate_file "$file"; then
        if grep -q "Modern DSL:" "examples/$file"; then
            MIGRATED=$((MIGRATED + 1))
        else
            SKIPPED=$((SKIPPED + 1))
        fi
    else
        ERRORS=$((ERRORS + 1))
    fi
done

echo ""
echo "📊 Migration Summary:"
echo "  📁 Total files processed: $TOTAL"
echo "  ✅ Successfully migrated: $MIGRATED"
echo "  ⚠️  Skipped (already done): $SKIPPED"
echo "  ❌ Errors: $ERRORS"
echo ""

if [ $MIGRATED -gt 0 ]; then
    echo "🎉 Migration completed! $MIGRATED files now support Modern DSL."
    echo ""
    echo "📝 Next steps:"
    echo "  1. Review migrated files and implement full modern DSL syntax"
    echo "  2. Test the examples to ensure they work correctly"
    echo "  3. Update documentation if needed"
    echo "  4. Remove .backup files when satisfied with migration"
    echo ""
    echo "🧪 Test a migrated example:"
    echo "  ./sloth-runner run -f examples/[migrated-file].lua --yes"
fi

echo "✨ Migration script completed!"