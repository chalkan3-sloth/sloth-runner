-- MODERN DSL - Dry Run Example
-- This example demonstrates dry run functionality with Modern DSL

-- Task 1: File operations that support dry run
local create_backup = task("create_backup")
    :description("Creates backup files (supports dry run)")
    :command(function(params)
        local dry_run = params.dry_run or false
        
        if dry_run then
            log.info("🔍 [DRY RUN] Would create backup of config files")
            log.info("🔍 [DRY RUN] Would backup: config.yaml -> config.yaml.backup")
            log.info("🔍 [DRY RUN] Would backup: settings.json -> settings.json.backup")
            return true, "Dry run: backup creation simulated", {
                files_to_backup = {"config.yaml", "settings.json"},
                backup_location = "/backups/" .. os.date("%Y%m%d")
            }
        else
            log.info("📁 Creating actual backup files...")
            -- In real scenario, this would create actual backups
            local backup_files = {"config.yaml.backup", "settings.json.backup"}
            
            for _, file in ipairs(backup_files) do
                log.info("✅ Created backup: " .. file)
            end
            
            return true, "Backup files created successfully", {
                backup_files = backup_files,
                backup_location = "/backups/" .. os.date("%Y%m%d")
            }
        end
    end)
    :timeout("30s")
    :build()

-- Task 2: Database migration task with dry run support
local database_migration = task("database_migration")
    :description("Runs database migration (supports dry run)")
    :depends_on({"create_backup"})
    :command(function(params, deps)
        local dry_run = params.dry_run or false
        
        if dry_run then
            log.info("🔍 [DRY RUN] Would execute database migration")
            log.info("🔍 [DRY RUN] Migration steps that would be executed:")
            log.info("    1. ALTER TABLE users ADD COLUMN email_verified BOOLEAN")
            log.info("    2. CREATE INDEX idx_users_email ON users(email)")
            log.info("    3. UPDATE users SET email_verified = false WHERE email_verified IS NULL")
            
            return true, "Dry run: migration steps validated", {
                migration_steps = 3,
                estimated_duration = "2m",
                affected_tables = {"users"}
            }
        else
            log.info("🗄️  Executing actual database migration...")
            
            -- Simulate migration steps
            local steps = {
                "ALTER TABLE users ADD COLUMN email_verified BOOLEAN",
                "CREATE INDEX idx_users_email ON users(email)", 
                "UPDATE users SET email_verified = false WHERE email_verified IS NULL"
            }
            
            for i, step in ipairs(steps) do
                log.info("📝 Step " .. i .. ": " .. step)
                -- In real scenario, execute SQL here
            end
            
            log.info("✅ Database migration completed successfully")
            return true, "Migration completed", {
                steps_executed = #steps,
                duration = "1m45s"
            }
        end
    end)
    :timeout("5m")
    :build()

-- Task 3: Service deployment with dry run
local deploy_service = task("deploy_service")
    :description("Deploys application service (supports dry run)")
    :depends_on({"database_migration"})
    :command(function(params, deps)
        local dry_run = params.dry_run or false
        local app_version = params.app_version or "v1.0.0"
        
        if dry_run then
            log.info("🔍 [DRY RUN] Would deploy service version: " .. app_version)
            log.info("🔍 [DRY RUN] Deployment plan:")
            log.info("    • Pull Docker image: myapp:" .. app_version)
            log.info("    • Stop current service")
            log.info("    • Start new service with updated configuration")
            log.info("    • Run health checks")
            log.info("    • Update load balancer")
            
            return true, "Dry run: deployment plan validated", {
                image = "myapp:" .. app_version,
                strategy = "rolling_update",
                estimated_downtime = "30s"
            }
        else
            log.info("🚀 Deploying service version: " .. app_version)
            
            local deployment_steps = {
                "Pulling Docker image: myapp:" .. app_version,
                "Stopping current service",
                "Starting new service", 
                "Running health checks",
                "Updating load balancer"
            }
            
            for i, step in ipairs(deployment_steps) do
                log.info("📦 " .. step .. "...")
                -- Simulate deployment time
            end
            
            log.info("✅ Service deployed successfully!")
            return true, "Deployment completed", {
                version = app_version,
                status = "running",
                health_check = "passed"
            }
        end
    end)
    :timeout("10m")
    :build()

-- Task 4: Notification task
local notify_completion = task("notify_completion")
    :description("Sends deployment notification")
    :depends_on({"deploy_service"})
    :command(function(params, deps)
        local dry_run = params.dry_run or false
        local deployment_result = deps.deploy_service
        
        if dry_run then
            log.info("🔍 [DRY RUN] Would send notification:")
            log.info("    • To: devops-team@company.com")
            log.info("    • Subject: Deployment Completed (DRY RUN)")
            log.info("    • Message: Service deployment dry run completed successfully")
            
            return true, "Dry run: notification prepared", {
                notification_type = "email",
                recipients = ["devops-team@company.com"]
            }
        else
            log.info("📧 Sending deployment notification...")
            log.info("✅ Notification sent to devops team")
            
            return true, "Notification sent", {
                sent_to = "devops-team@company.com",
                timestamp = os.time()
            }
        end
    end)
    :timeout("30s")
    :build()

-- Workflow with dry run support
workflow.define("deployment_with_dry_run", {
    description = "Application deployment pipeline with dry run support - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "DevOps Team",
        tags = {"deployment", "dry-run", "modern-dsl", "pipeline"},
        supports_dry_run = true
    },
    
    tasks = {
        create_backup,
        database_migration, 
        deploy_service,
        notify_completion
    },
    
    config = {
        timeout = "30m",
        retry_policy = "linear",
        max_parallel_tasks = 1,  -- Sequential execution for deployment
        fail_fast = true
    },
    
    on_start = function(params)
        local dry_run = params.dry_run or false
        
        if dry_run then
            log.info("🔍 STARTING DRY RUN MODE")
            log.info("    No actual changes will be made")
            log.info("    All operations will be simulated")
        else
            log.info("🚀 STARTING ACTUAL DEPLOYMENT")
            log.info("    Real changes will be made to systems")
        end
        
        return true
    end,
    
    on_complete = function(success, results, params)
        local dry_run = params.dry_run or false
        
        if dry_run then
            if success then
                log.info("✅ DRY RUN COMPLETED SUCCESSFULLY!")
                log.info("🎯 Summary:")
                log.info("    • All deployment steps validated")
                log.info("    • No actual changes were made")
                log.info("    • Ready for real deployment")
            else
                log.error("❌ DRY RUN FAILED!")
                log.warn("🚨 Issues found that need to be resolved before real deployment")
            end
        else
            if success then
                log.info("🎉 DEPLOYMENT COMPLETED SUCCESSFULLY!")
                log.info("🚀 Application is now running in production")
            else
                log.error("❌ DEPLOYMENT FAILED!")
                log.error("🚨 Please check logs and resolve issues")
            end
        end
        
        return true
    end
})

-- Usage examples:
-- For dry run: sloth-runner run -f dry_run_example.lua --dry-run
-- For actual deployment: sloth-runner run -f dry_run_example.lua