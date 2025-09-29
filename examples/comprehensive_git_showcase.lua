-- Comprehensive Git Showcase - Modern DSL
-- Demonstrates advanced Git operations using Modern DSL syntax

-- Task 1: Clone repository
local clone_repo_task = task("clone_repository")
    :description("Clone a Git repository with Modern DSL")
    :command(function(params, deps)
        log.info("🔄 Cloning repository...")
        
        local repo_url = params.repo_url or "https://github.com/chalkan3-sloth/sloth-runner.git"
        local target_dir = params.target_dir or "./cloned-repo"
        
        -- Use git module for enhanced operations
        local repo = git.clone(repo_url, target_dir)
        
        if not repo then
            return false, "Failed to clone repository", {}
        end
        
        local status = git.status(target_dir)
        
        return true, "Repository cloned successfully", {
            repo_path = target_dir,
            current_branch = status.branch,
            commit_hash = status.commit,
            repo_url = repo_url
        }
    end)
    :timeout("5m")
    :artifacts({"cloned-repo"})
    :on_success(function(params, output)
        log.info("✅ Repository cloned to: " .. output.repo_path)
        log.info("📍 Current branch: " .. output.current_branch)
    end)
    :build()

-- Task 2: Check repository status
local check_status_task = task("check_git_status")
    :description("Check Git repository status and information")
    :depends_on({"clone_repository"})
    :consumes({"cloned-repo"})
    :command(function(params, deps)
        log.info("📊 Checking Git status...")
        
        local repo_path = "./cloned-repo"
        
        -- Get comprehensive Git information
        local status = git.status(repo_path)
        local log_info = git.log(repo_path, { limit = 5 })
        local branches = git.list_branches(repo_path)
        
        log.info("Current branch: " .. status.branch)
        log.info("Latest commit: " .. status.commit)
        log.info("Repository is clean: " .. tostring(status.clean))
        
        return true, "Git status checked", {
            status = status,
            recent_commits = log_info,
            branches = branches,
            repo_info = {
                has_uncommitted_changes = not status.clean,
                total_branches = #branches,
                latest_commit_message = log_info[1] and log_info[1].message or "N/A"
            }
        }
    end)
    :build()

-- Task 3: Create and switch to new branch (if needed)
local create_branch_task = task("create_feature_branch")
    :description("Create and switch to a new feature branch")
    :depends_on({"check_git_status"})
    :consumes({"cloned-repo"})
    :command(function(params, deps)
        log.info("🌿 Creating feature branch...")
        
        local repo_path = "./cloned-repo"
        local branch_name = params.branch_name or "feature/modern-dsl-demo"
        
        -- Create and checkout new branch
        local result = git.checkout(repo_path, branch_name, { create = true })
        
        if not result.success then
            log.warn("Branch might already exist, trying to switch...")
            result = git.checkout(repo_path, branch_name)
        end
        
        if result.success then
            local current_status = git.status(repo_path)
            return true, "Feature branch created/switched", {
                branch_name = branch_name,
                current_branch = current_status.branch,
                operation = "branch_created"
            }
        else
            return false, "Failed to create/switch to branch: " .. result.error
        end
    end)
    :run_if(function(params, deps)
        -- Only create branch if we're on main/master
        local current_branch = deps.check_git_status.status.branch
        return current_branch == "main" or current_branch == "master"
    end)
    :build()

-- Task 4: Make some changes and commit
local make_changes_task = task("make_git_changes")
    :description("Make changes and create a commit")
    :depends_on({"create_feature_branch"})
    :consumes({"cloned-repo"})
    :command(function(params, deps)
        log.info("✏️ Making changes and committing...")
        
        local repo_path = "./cloned-repo"
        
        -- Create a demo file
        local demo_content = "# Modern DSL Demo\n\nThis file was created by the Modern DSL Git showcase at " .. os.date()
        fs.write_file(repo_path .. "/MODERN_DSL_DEMO.md", demo_content)
        
        -- Stage the changes
        local add_result = git.add(repo_path, "MODERN_DSL_DEMO.md")
        if not add_result.success then
            return false, "Failed to stage changes: " .. add_result.error
        end
        
        -- Commit the changes
        local commit_result = git.commit(repo_path, {
            message = "Add Modern DSL demo file",
            author = "Sloth Runner <sloth@runner.dev>"
        })
        
        if commit_result.success then
            local new_status = git.status(repo_path)
            return true, "Changes committed successfully", {
                commit_hash = new_status.commit,
                files_changed = {"MODERN_DSL_DEMO.md"},
                commit_message = "Add Modern DSL demo file"
            }
        else
            return false, "Failed to commit changes: " .. commit_result.error
        end
    end)
    :build()

-- Task 5: Show final repository state
local show_final_state_task = task("show_final_state")
    :description("Display final repository state and summary")
    :depends_on({"make_git_changes"})
    :consumes({"cloned-repo"})
    :command(function(params, deps)
        log.info("📋 Showing final repository state...")
        
        local repo_path = "./cloned-repo"
        
        -- Get final status
        local final_status = git.status(repo_path)
        local final_log = git.log(repo_path, { limit = 3 })
        
        log.info("=== Final Repository State ===")
        log.info("Current branch: " .. final_status.branch)
        log.info("Latest commit: " .. final_status.commit)
        log.info("Repository clean: " .. tostring(final_status.clean))
        
        log.info("=== Recent Commits ===")
        for i, commit in ipairs(final_log) do
            log.info(i .. ". " .. commit.hash .. " - " .. commit.message)
        end
        
        return true, "Repository showcase completed", {
            final_status = final_status,
            recent_commits = final_log,
            showcase_summary = {
                operations_performed = {"clone", "status_check", "branch_creation", "file_changes", "commit"},
                files_created = {"MODERN_DSL_DEMO.md"},
                final_branch = final_status.branch
            }
        }
    end)
    :on_success(function(params, output)
        log.info("🎉 Git showcase completed successfully!")
        log.info("📊 Summary: " .. #output.showcase_summary.operations_performed .. " operations performed")
    end)
    :build()

-- Define the comprehensive Git workflow
workflow.define("git_showcase_workflow", {
    description = "Comprehensive Git operations showcase - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"git", "vcs", "showcase", "modern-dsl"},
        complexity = "intermediate",
        estimated_duration = "5m"
    },
    
    tasks = {
        clone_repo_task,
        check_status_task,
        create_branch_task,
        make_changes_task,
        show_final_state_task
    },
    
    config = {
        timeout = "15m",
        retry_policy = "exponential",
        max_parallel_tasks = 1, -- Git operations should be sequential
        create_workdir_before_run = true,
        clean_workdir_after_run = false -- Keep the cloned repo for inspection
    },
    
    on_start = function()
        log.info("🚀 Starting comprehensive Git showcase...")
        log.info("📚 This showcase demonstrates Modern DSL with Git operations")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ Git showcase workflow completed successfully!")
            log.info("🎯 All Git operations executed with Modern DSL")
            log.info("📁 Check './cloned-repo' for the results")
        else
            log.error("❌ Git showcase workflow failed!")
            log.error("🔧 Check individual task results for details")
        end
        return true
    end
})
