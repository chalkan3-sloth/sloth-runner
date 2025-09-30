-- MODERN DSL - Git Module Showcase
-- This example demonstrates comprehensive Git operations using Modern DSL

-- Task 1: Clone repository with authentication
local clone_repo = task("clone_repository")
    :description("Clones a Git repository with modern features")
    :command(function(params)
        local repo_url = params.repo_url or "https://github.com/example/project.git"
        local branch = params.branch or "main"
        local target_dir = params.target_dir or "project"
        
        log.info("📦 Cloning repository: " .. repo_url)
        log.info("🌿 Branch: " .. branch)
        
        -- Use modern git module with enhanced features
        local result = git.clone({
            url = repo_url,
            branch = branch,
            directory = target_dir,
            depth = 1,  -- Shallow clone for performance
            tags = false
        })
        
        if result.success then
            log.info("✅ Repository cloned successfully")
            return true, "Repository cloned", {
                directory = target_dir,
                commit = result.commit_hash,
                branch = branch
            }
        else
            log.error("❌ Failed to clone repository: " .. result.error)
            return false, "Clone failed: " .. result.error
        end
    end)
    :timeout("5m")
    :retries(2, "exponential")
    :build()

-- Task 2: Check repository status and information
local check_status = task("check_git_status")
    :description("Examines Git repository status and information")
    :depends_on({"clone_repository"})
    :command(function(params, deps)
        local repo_dir = deps.clone_repository.directory
        
        log.info("🔍 Checking Git status in: " .. repo_dir)
        
        -- Check repository status
        local status = git.status({
            directory = repo_dir
        })
        
        if status.success then
            log.info("📊 Repository Status:")
            log.info("    • Current branch: " .. (status.branch or "unknown"))
            log.info("    • Clean working directory: " .. tostring(status.clean))
            log.info("    • Modified files: " .. #(status.modified or {}))
            log.info("    • Untracked files: " .. #(status.untracked or {}))
            
            -- Get commit information
            local commit_info = git.log({
                directory = repo_dir,
                max_count = 1
            })
            
            if commit_info.success and #commit_info.commits > 0 then
                local latest = commit_info.commits[1]
                log.info("📝 Latest commit:")
                log.info("    • Hash: " .. latest.hash)
                log.info("    • Author: " .. latest.author)
                log.info("    • Message: " .. latest.message)
                log.info("    • Date: " .. latest.date)
            end
            
            return true, "Status checked successfully", {
                branch = status.branch,
                clean = status.clean,
                latest_commit = commit_info.commits[1]
            }
        else
            log.error("❌ Failed to check status: " .. status.error)
            return false, "Status check failed"
        end
    end)
    :timeout("1m")
    :build()

-- Task 3: Create and manage branches
local branch_operations = task("branch_operations")
    :description("Demonstrates Git branch operations")
    :depends_on({"check_git_status"})
    :command(function(params, deps)
        local repo_dir = deps.clone_repository.directory
        local feature_branch = "feature/modern-dsl-demo"
        
        log.info("🌿 Performing branch operations...")
        
        -- List all branches
        local branches = git.branch({
            directory = repo_dir,
            all = true
        })
        
        if branches.success then
            log.info("📋 Available branches:")
            for _, branch in ipairs(branches.branches) do
                local marker = branch.current and "* " or "  "
                log.info("    " .. marker .. branch.name)
            end
        end
        
        -- Create new feature branch
        local create_branch = git.checkout({
            directory = repo_dir,
            branch = feature_branch,
            create = true
        })
        
        if create_branch.success then
            log.info("✅ Created and switched to branch: " .. feature_branch)
            
            return true, "Branch operations completed", {
                current_branch = feature_branch,
                total_branches = #(branches.branches or {})
            }
        else
            log.warn("⚠️  Branch operations had issues, continuing...")
            return true, "Branch operations completed with warnings", {}
        end
    end)
    :timeout("2m")
    :build()

-- Task 4: File operations and staging
local file_operations = task("file_operations")
    :description("Demonstrates Git file operations and staging")
    :depends_on({"branch_operations"})
    :command(function(params, deps)
        local repo_dir = deps.clone_repository.directory
        
        log.info("📝 Performing file operations...")
        
        -- Create a test file
        local test_file = repo_dir .. "/SLOTH_RUNNER_DEMO.md"
        local content = [[# Sloth Runner Modern DSL Demo

This file was created by the Sloth Runner Git module showcase.

## Features Demonstrated

- Repository cloning
- Status checking
- Branch operations
- File staging and commits
- Remote operations

Generated at: ]] .. os.date() .. [[

## Modern DSL Benefits

- Fluent, chainable API
- Enhanced error handling
- Built-in retry logic
- Comprehensive logging
]]
        
        -- Write file using fs module
        local write_result = fs.write(test_file, content)
        
        if write_result then
            log.info("✅ Created demo file: " .. test_file)
            
            -- Add file to staging
            local add_result = git.add({
                directory = repo_dir,
                files = {"SLOTH_RUNNER_DEMO.md"}
            })
            
            if add_result.success then
                log.info("✅ File staged successfully")
                
                return true, "File operations completed", {
                    file_created = test_file,
                    staged = true
                }
            else
                log.warn("⚠️  Failed to stage file: " .. add_result.error)
                return true, "File created but not staged", {
                    file_created = test_file,
                    staged = false
                }
            end
        else
            log.error("❌ Failed to create demo file")
            return false, "File creation failed"
        end
    end)
    :timeout("1m")
    :build()

-- Task 5: Commit and push operations
local commit_and_push = task("commit_and_push")
    :description("Demonstrates commit and push operations")
    :depends_on({"file_operations"})
    :command(function(params, deps)
        local repo_dir = deps.clone_repository.directory
        
        log.info("💾 Performing commit operations...")
        
        -- Check if there are changes to commit
        local status = git.status({
            directory = repo_dir
        })
        
        if status.success and not status.clean then
            -- Commit the changes
            local commit_result = git.commit({
                directory = repo_dir,
                message = "Add Sloth Runner Modern DSL demonstration file\n\nGenerated by: sloth-runner git module showcase\nTimestamp: " .. os.date(),
                author = "Sloth Runner <demo@sloth-runner.dev>"
            })
            
            if commit_result.success then
                log.info("✅ Changes committed successfully")
                log.info("📝 Commit hash: " .. commit_result.commit_hash)
                
                -- Note: In a real scenario, you might want to push to a remote
                -- For demo purposes, we'll just simulate the push operation
                log.info("🔄 Push operation would be performed here")
                log.info("    (Skipped in demo to avoid unauthorized pushes)")
                
                return true, "Commit completed", {
                    commit_hash = commit_result.commit_hash,
                    pushed = false,
                    reason = "Demo mode - push skipped"
                }
            else
                log.error("❌ Failed to commit: " .. commit_result.error)
                return false, "Commit failed"
            end
        else
            log.info("ℹ️  No changes to commit")
            return true, "No changes to commit", {
                clean_working_directory = true
            }
        end
    end)
    :timeout("2m")
    :build()

-- Task 6: Cleanup operations
local cleanup_repo = task("cleanup_demo")
    :description("Cleans up demo repository")
    :depends_on({"commit_and_push"})
    :command(function(params, deps)
        local repo_dir = deps.clone_repository.directory
        
        log.info("🧹 Cleaning up demo repository...")
        
        -- Get final repository information
        local final_status = git.status({
            directory = repo_dir
        })
        
        if final_status.success then
            log.info("📊 Final Repository State:")
            log.info("    • Directory: " .. repo_dir)
            log.info("    • Branch: " .. (final_status.branch or "unknown"))
            log.info("    • Clean: " .. tostring(final_status.clean))
        end
        
        -- Note: In production, you might want to remove the cloned directory
        -- For demo purposes, we'll leave it for inspection
        log.info("📁 Repository preserved for inspection at: " .. repo_dir)
        log.info("🔍 You can manually inspect the repository state")
        
        return true, "Cleanup completed", {
            repository_path = repo_dir,
            preserved = true
        }
    end)
    :timeout("1m")
    :build()

-- Modern Workflow Definition for Git Operations
workflow.define("git_module_showcase", {
    description = "Comprehensive Git module showcase using Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"git", "version-control", "modern-dsl", "showcase"},
        complexity = "intermediate",
        estimated_duration = "10m"
    },
    
    tasks = {
        clone_repo,
        check_status,
        branch_operations,
        file_operations,
        commit_and_push,
        cleanup_repo
    },
    
    config = {
        timeout = "20m",
        retry_policy = "exponential",
        max_parallel_tasks = 1,  -- Sequential execution for Git operations
        fail_fast = false
    },
    
    on_start = function()
        log.info("🚀 Starting Git Module Showcase...")
        log.info("🔧 This workflow demonstrates:")
        log.info("   • Repository cloning with authentication")
        log.info("   • Status checking and information gathering")
        log.info("   • Branch creation and management")
        log.info("   • File operations and staging")
        log.info("   • Commit operations with metadata")
        log.info("   • Repository cleanup")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 Git Module Showcase completed successfully!")
            log.info("✅ All Git operations demonstrated:")
            
            -- Log summary of operations
            for task_name, result in pairs(results) do
                if result.success then
                    log.info("   ✓ " .. task_name .. ": " .. (result.output or "completed"))
                else
                    log.warn("   ⚠ " .. task_name .. ": " .. (result.error or "failed"))
                end
            end
            
            log.info("📚 Git module features successfully showcased!")
        else
            log.error("❌ Git Module Showcase encountered issues!")
            log.info("🔍 Check individual task logs for details")
        end
        
        return true
    end
})

-- Usage examples:
-- Basic run: sloth-runner run -f git_module_showcase.lua
-- With custom repo: sloth-runner run -f git_module_showcase.lua --set repo_url=https://github.com/your/repo.git
-- With custom branch: sloth-runner run -f git_module_showcase.lua --set branch=develop