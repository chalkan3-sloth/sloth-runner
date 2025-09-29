-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:31 -03

local build_task = task("build_app")

local build_task = task("build_app")
--    task("name"):description("desc"):command(func):timeout("30s"):build()
--   🔄 Full task() and workflow() function implementations
--   🔄 Enhanced error handling and validation
--   🔄 Performance optimizations for new DSL
--   🔄 Advanced features (saga patterns, circuit breakers)

-- 💡 USAGE EXAMPLES
-- =================

-- Modern DSL Syntax (Target Implementation):
--[[
local build_task = task("build_app")
    :description("Build application with modern DSL")
    :command(function(params, deps)
        log.info("Building application...")
        return exec.run("go build -o app ./cmd/main.go")
    end)
    :timeout("5m")
    :retries(2, "exponential")
    :artifacts({"app"})
    :on_success(function(params, output)
        log.info("Build completed successfully!")
    end)
    :build()

--    workflow.define("name", { description, tasks, config })

-- 3. Enhanced Features:
--    • :depends_on() - Modern dependency management
--    • :condition() - Advanced conditional execution  
--    • :retries() - Sophisticated retry strategies
--    • :timeout() - Granular timeout control
--    • :on_success() / :on_failure() - Lifecycle hooks
--    • :artifacts() - Enhanced artifact management
--    • :async() - Modern asynchronous execution

-- 4. Workflow-Level Features:
--    • metadata - Rich workflow metadata
--    • config - Centralized configuration
--    • on_start/on_complete - Workflow lifecycle hooks
--    • pre_conditions - Workflow prerequisites

-- 5. Backward Compatibility:
--    • All legacy TaskDefinitions still work
--    • Gradual migration path provided
--    • Dual syntax support in same file

-- 📋 CURRENT STATUS
-- =================

-- ✅ COMPLETED:
--   ✅ New DSL syntax design and specification
--   ✅ Core infrastructure for DSL support  
--   ✅ Legacy compatibility layer maintained
--   ✅ All examples structured for modern DSL
--   ✅ Key examples fully implemented with modern syntax
--   ✅ Automated migration tooling created
--   ✅ Documentation and examples updated

-- 🚧 IN PROGRESS / NEXT STEPS:
--   🔄 Complete modern DSL runtime implementation
--   🔄 Full task() and workflow() function implementations
--   🔄 Enhanced error handling and validation
--   🔄 Performance optimizations for new DSL
--   🔄 Advanced features (saga patterns, circuit breakers)

-- 💡 USAGE EXAMPLES
-- =================

-- Modern DSL Syntax (Target Implementation):
--[[
local build_task = task("build_app")
    :description("Build application with modern DSL")
    :command(function(params, deps)
        log.info("Building application...")
        return exec.run("go build -o app ./cmd/main.go")
    end)
    :timeout("5m")
    :retries(2, "exponential")
    :artifacts({"app"})
    :on_success(function(params, output)
        log.info("Build completed successfully!")
    end)
    :build()

workflow.define("ci_pipeline", {
    description = "Continuous Integration Pipeline",
    version = "2.0.0",
    tasks = { build_task },
    config = {
        timeout = "30m",
        retry_policy = "exponential"
    }
})
