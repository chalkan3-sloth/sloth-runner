-- Simple working workflow

-- Task 1: Simple echo task
local hello_task = task("hello")
    :description("Says hello")
    :command(function(params)
        log.info("Saying hello...")
        return true, "echo 'Hello World'", { message = "Hello from Sloth Runner!" }
    end)
    :build()

-- Define workflow
workflow.define("simple_working", {
    description = "Simple working workflow",
    version = "1.0.0",
    tasks = { hello_task }
})

-- Export outputs (this should appear in the JSON output)
outputs = {
    status = "success",
    timestamp = os.date(),
    app_url = "https://myapp.example.com",
    version = "1.2.3",
    environment = "production"
}