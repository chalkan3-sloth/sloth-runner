-- examples/id_demo.lua
-- Demonstration of task IDs and workflow structure

TaskDefinitions = {
    -- Build and Deploy Pipeline
    build_pipeline = {
        description = "Build and deployment pipeline with unique IDs",
        tasks = {
            {
                name = "setup",
                description = "Setup build environment",
                command = "echo 'Setting up build environment...'"
            },
            {
                name = "compile",
                description = "Compile the application",
                command = "echo 'Compiling application...'",
                depends_on = {"setup"}
            },
            {
                name = "test",
                description = "Run unit tests",
                command = "echo 'Running tests...'",
                depends_on = {"compile"}
            },
            {
                name = "package",
                description = "Package the application",
                command = "echo 'Packaging application...'",
                depends_on = {"compile", "test"}
            }
        }
    },
    
    -- Deployment Group
    deploy_pipeline = {
        description = "Deployment tasks with environment management",
        tasks = {
            {
                name = "deploy_staging",
                description = "Deploy to staging environment",
                command = "echo 'Deploying to staging...'",
                depends_on = {"package"}
            },
            {
                name = "integration_test",
                description = "Run integration tests",
                command = "echo 'Running integration tests...'",
                depends_on = {"deploy_staging"}
            },
            {
                name = "deploy_production",
                description = "Deploy to production environment", 
                command = "echo 'Deploying to production...'",
                depends_on = {"integration_test"}
            }
        }
    }
}