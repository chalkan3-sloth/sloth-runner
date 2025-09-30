-- Simple test workflow

-- Task 1: Simple echo task
local hello_task = task("hello")
    :description("Says hello")
    :command(function(params)
        log.info("Saying hello...")
        return true, "echo 'Hello World'", { message = "Hello from Sloth Runner!" }
    end)
    :build()

-- Task 2: Another task
local goodbye_task = task("goodbye")
    :description("Says goodbye")
    :depends_on({"hello"})
    :command(function(params, deps)
        log.info("Saying goodbye...")
        return true, "echo 'Goodbye'", { message = "Goodbye from Sloth Runner!", 
                                         previous = deps.hello.message }
    end)
    :build()

-- Define workflow
workflow.define("simple_test", {
    description = "Simple test workflow",
    version = "1.0.0",
    tasks = { hello_task, goodbye_task }
})

-- Export outputs (this should appear in the JSON output)
outputs = {
    test_result = "success",
    timestamp = os.date(),
    app_url = "https://example.com",
    version = "1.2.3"
}