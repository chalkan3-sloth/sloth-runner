#!/usr/bin/env python3

import os
import re
import shutil
from pathlib import Path

def create_modern_dsl_template(original_file_path, original_content):
    """Create a Modern DSL template based on the original file name and content"""
    
    file_name = Path(original_file_path).stem
    
    # Extract any meaningful info from original content
    description_match = re.search(r'description\s*=\s*["\']([^"\']+)["\']', original_content)
    description = description_match.group(1) if description_match else f"Modern DSL version of {file_name}"
    
    # Determine category based on file path
    if 'beginner' in original_file_path:
        category = 'beginner'
        complexity = 'simple'
    elif 'intermediate' in original_file_path:
        category = 'intermediate' 
        complexity = 'moderate'
    elif 'advanced' in original_file_path:
        category = 'advanced'
        complexity = 'complex'
    else:
        category = 'general'
        complexity = 'basic'
    
    # Create modern DSL template
    template = f'''-- MODERN DSL ONLY - {description}
-- Converted from legacy TaskDefinition format
-- Category: {category}, Complexity: {complexity}

-- Main task using Modern DSL
local main_task = task("{file_name}_task")
    :description("{description}")
    :command(function(params, deps)
        log.info("🚀 Executing {file_name} with Modern DSL...")
        
        -- TODO: Replace with actual implementation from original file
        -- Original logic should be migrated here
        
        return true, "Task completed successfully", {{
            task_name = "{file_name}",
            execution_time = os.time(),
            status = "success"
        }}
    end)
    :timeout("5m")
    :retries(2, "exponential")
    :on_success(function(params, output)
        log.info("✅ {file_name} task completed successfully")
    end)
    :on_failure(function(params, error)
        log.error("❌ {file_name} task failed: " .. error)
    end)
    :build()

-- Additional tasks can be added here following the same pattern
-- local secondary_task = task("{file_name}_secondary")
--     :description("Secondary task for {file_name}")
--     :depends_on({{"{file_name}_task"}})
--     :command(function(params, deps)
--         -- Secondary logic here
--         return true, "Secondary task completed", {{}}
--     end)
--     :build()

-- Modern Workflow Definition
workflow.define("{file_name}_workflow", {{
    description = "{description} - Modern DSL",
    version = "2.0.0",
    
    metadata = {{
        author = "Sloth Runner Team",
        category = "{category}",
        complexity = "{complexity}",
        tags = {{"{file_name}", "modern-dsl", "{category}"}},
        created_at = os.date(),
        migrated_from = "TaskDefinition format"
    }},
    
    tasks = {{
        main_task
        -- Add additional tasks here
    }},
    
    config = {{
        timeout = "15m",
        retry_policy = "exponential",
        max_parallel_tasks = 2,
        cleanup_on_failure = true
    }},
    
    on_start = function()
        log.info("🚀 Starting {file_name} workflow...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ {file_name} workflow completed successfully!")
        else
            log.error("❌ {file_name} workflow failed!")
        end
        return true
    end
}})

-- Migration Note:
-- This file has been converted from legacy TaskDefinition format to Modern DSL
-- TODO: Review and implement the original logic in the Modern DSL structure above
-- Original backup saved as {Path(original_file_path).name}.backup
'''
    
    return template

def convert_file_to_modern_dsl(file_path):
    """Convert a single file to Modern DSL"""
    
    try:
        # Read original content
        with open(file_path, 'r', encoding='utf-8') as f:
            original_content = f.read()
        
        # Check if already modern
        if 'task(' in original_content and 'TaskDefinition' not in original_content:
            print(f"✅ Already modern: {file_path}")
            return True
        
        # Create backup
        backup_path = f"{file_path}.backup"
        shutil.copy2(file_path, backup_path)
        
        # Generate modern DSL content
        modern_content = create_modern_dsl_template(file_path, original_content)
        
        # Write modern content
        with open(file_path, 'w', encoding='utf-8') as f:
            f.write(modern_content)
        
        print(f"🔄 Converted: {file_path}")
        return True
        
    except Exception as e:
        print(f"❌ Error converting {file_path}: {e}")
        return False

def main():
    """Main conversion function"""
    
    print("🧹 Starting mass conversion to Modern DSL...")
    
    # Find all .lua files with TaskDefinition
    examples_dir = Path("examples")
    lua_files = []
    
    for lua_file in examples_dir.rglob("*.lua"):
        try:
            with open(lua_file, 'r', encoding='utf-8') as f:
                content = f.read()
                if 'TaskDefinition' in content or 'task_definition' in content:
                    lua_files.append(str(lua_file))
        except Exception as e:
            print(f"⚠️  Error reading {lua_file}: {e}")
    
    print(f"📝 Found {len(lua_files)} files to convert")
    
    if not lua_files:
        print("✅ No files need conversion!")
        return
    
    # Convert each file
    converted = 0
    failed = 0
    
    for file_path in lua_files:
        if convert_file_to_modern_dsl(file_path):
            converted += 1
        else:
            failed += 1
    
    print(f"""
🎉 Conversion Summary:
✅ Converted: {converted} files
❌ Failed: {failed} files
📁 Total processed: {len(lua_files)} files

📚 Next steps:
1. Review converted files and implement original logic
2. Test the workflows with: ./sloth-runner run -f <file>
3. Update documentation if needed
4. Remove .backup files when satisfied with conversion
""")

if __name__ == "__main__":
    main()