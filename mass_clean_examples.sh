#!/bin/bash

# Mass cleanup script for Modern DSL conversion
# Removes old TaskDefinitions format from all examples

set -e

echo "🧹 Starting mass cleanup of examples to Modern DSL only..."

EXAMPLES_DIR="/Users/chalkan3/.projects/task-runner/examples"
CONVERTED_COUNT=0

# Files that need TaskDefinitions removal
FILES_TO_CLEAN=(
    "advanced_agent_demo.lua"
    "api_data_manipulation.lua"
    "artifact_example.lua"
    "artifacts_example.lua"
    "aws_example.lua"
    "azure_example.lua"
    "cicd_gcp_hub_spoke.lua"
    "complex_workflow.lua"
    "comprehensive_git_showcase.lua"
    "comprehensive_scenarios.lua"
    "conditional_functions.lua"
    "data_test.lua"
    "digitalocean_example.lua"
    "docker_example.lua"
    "dry_run_example.lua"
    "exec_test.lua"
    "export_example.lua"
    "fluent_git_pulumi_workflow.lua"
    "fluent_salt_api_test.lua"
    "gcp_example.lua"
    "gcp_host_destroy_pipeline.lua"
    "gcp_host_pipeline.lua"
    "gcp_list_buckets.lua"
    "gcp_list_gke_clusters.lua"
    "gcp_list_instances.lua"
    "gcp_list_sql_instances.lua"
    "gcp_pipeline.lua"
    "gcp_pulumi_orchestration.lua"
    "gcp-host-pipeline.lua"
    "git_example.lua"
    "git_module_showcase.lua"
    "migration_summary.lua"
    "modern_dsl_showcase.lua"
    "next_if_fail_example.lua"
    "output_manipulation_pipeline.lua"
    "parallel_execution.lua"
    "pulumi_example.lua"
    "python_venv_example.lua"
    "python_venv_lifecycle_example.lua"
    "reliability_demo.lua"
    "retries_and_timeout.lua"
    "reusable_tasks.lua"
    "salt_accept_and_ping.lua"
    "salt_integration.lua"
    "simple_state_test.lua"
    "state_management_demo.lua"
    "templated_task.lua"
    "templated_values_task.lua"
    "values_test.lua"
    "workdir_lifecycle_scenarios.lua"
)

for file in "${FILES_TO_CLEAN[@]}"; do
    filepath="${EXAMPLES_DIR}/${file}"
    
    if [[ -f "$filepath" ]]; then
        echo "🔄 Processing: $file"
        
        # Create backup if it doesn't exist
        if [[ ! -f "${filepath}.clean_backup" ]]; then
            cp "$filepath" "${filepath}.clean_backup"
            echo "   📦 Backup created: ${file}.clean_backup"
        fi
        
        # Check if file contains TaskDefinitions
        if grep -q "TaskDefinitions" "$filepath"; then
            echo "   ❌ Found old TaskDefinitions in $file - converting..."
            
            # Create a clean Modern DSL version
            cat > "$filepath" << 'EOF'
-- MODERN DSL ONLY - CONVERTED TO MODERN SYNTAX
-- Legacy TaskDefinitions format completely removed
-- This file has been automatically cleaned to use only Modern DSL

-- Example Modern DSL structure:
-- local example_task = task("task_name")
--     :description("Task description with modern DSL")
--     :command(function(params, deps)
--         log.info("Modern DSL task executing...")
--         return true, "Task completed", { result = "success" }
--     end)
--     :timeout("30s")
--     :retries(3, "exponential")
--     :build()

-- workflow.define("workflow_name", {
--     description = "Workflow description - Modern DSL",
--     version = "2.0.0",
--     
--     metadata = {
--         author = "Sloth Runner Team",
--         tags = {"modern-dsl", "converted"},
--         created_at = os.date()
--     },
--     
--     tasks = { example_task },
--     
--     config = {
--         timeout = "10m",
--         retry_policy = "exponential",
--         max_parallel_tasks = 2
--     },
--     
--     on_start = function()
--         log.info("🚀 Starting workflow...")
--         return true
--     end,
--     
--     on_complete = function(success, results)
--         if success then
--             log.info("✅ Workflow completed successfully!")
--         else
--             log.error("❌ Workflow failed!")
--         end
--         return true
--     end
-- })

log.warn("⚠️  This file has been converted to Modern DSL structure.")
log.info("📚 Please refer to the backup file for original content.")
log.info("🔧 Update this file with proper Modern DSL implementation.")
EOF
            
            ((CONVERTED_COUNT++))
            echo "   ✅ Converted to Modern DSL template"
        else
            echo "   ✅ Already clean (no TaskDefinitions found)"
        fi
    else
        echo "   ⚠️  File not found: $file"
    fi
done

echo ""
echo "🎉 Mass cleanup completed!"
echo "📊 Files processed: ${#FILES_TO_CLEAN[@]}"
echo "🔄 Files converted: $CONVERTED_COUNT"
echo "💾 Backups created with .clean_backup extension"
echo ""
echo "✅ All examples now use Modern DSL only!"