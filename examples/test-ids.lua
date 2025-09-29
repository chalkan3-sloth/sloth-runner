-- Simple test workflow for demonstrating list command with IDs

-- Legacy TaskDefinitions format
TaskDefinitions = {
    test_group = {
        description = "Simple test group for demonstrating task IDs",
        tasks = {
            {
                name = "hello_task",
                description = "Says hello",
                command = "echo 'Hello from task!'",
            },
            {
                name = "world_task", 
                description = "Says world",
                command = "echo 'World from task!'",
                depends_on = {"hello_task"}
            }
        }
    },

    deploy_group = {
        description = "Deployment tasks with IDs",
        tasks = {
            {
                name = "build",
                description = "Build the application",
                command = "echo 'Building...'",
            },
            {
                name = "test",
                description = "Run tests",
                command = "echo 'Testing...'",
                depends_on = {"build"}
            },
            {
                name = "deploy",
                description = "Deploy to production",
                command = "echo 'Deploying...'",
                depends_on = {"build", "test"}
            }
        }
    }
}