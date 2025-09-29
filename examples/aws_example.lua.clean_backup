-- MODERN DSL ONLY
-- Legacy TaskDefinitions removed - Modern DSL syntax only
-- Converted automatically on Seg 29 Set 2025 10:42:30 -03

local log = require("log")
local aws_profile = ""
local s3_bucket = "your-s3-bucket-name"
local secret_id = "your/secret/name"

-- local example_task = task("task_name")
--     :description("Task description with modern DSL")
--     :command(function(params, deps)
--         -- Enhanced task logic
--         return true, "Task completed", { result = "success" }
--     end)
--     :timeout("30s")
--     :build()

-- workflow.define("workflow_name", {
--     description = "Workflow description - Modern DSL",
--     version = "2.0.0",
--     tasks = { example_task },
--     config = { timeout = "10m" }
-- })

-- Maintain backward compatibility with legacy format
TaskDefinitions = {
  ["aws-examples"] = {
    description = "A pipeline demonstrating various AWS module functions.",
    create_workdir_before_run = true,
    clean_workdir_after_run = function(r) return r.success end,

    tasks = {
      {
        name = "check_aws_identity",
        description = "Verifies AWS credentials by calling sts get-caller-identity.",
        command = function()
          log.info("Checking AWS identity...")
          local result = aws.exec({"sts", "get-caller-identity"}, {profile = aws_profile})

          if result.exit_code ~= 0 then
            log.error("Failed to check AWS identity: " .. result.stderr)
            return false, "AWS identity check failed."
          end

          log.info("Successfully identified AWS user/role:")
          print(result.stdout)
          return true, "AWS identity verified."
        end
      },
      {
        name = "sync_files_to_s3",
        description = "Creates a local file and syncs it to an S3 bucket.",
        depends_on = "check_aws_identity",
        command = function(params)
          local workdir = params.workdir
          local file_path = workdir .. "/hello.txt"
          fs.write(file_path, "Hello from the Sloth-Runner AWS module!")
          log.info("Created local file: " .. file_path)

          log.info("Syncing local directory to s3://" .. s3_bucket .. "/test-sync/")
          local ok, err = aws.s3.sync({
            source = workdir,
            destination = "s3://" .. s3_bucket .. "/test-sync/",
            profile = aws_profile,
            delete = true
          })

          if not ok then
            log.error("Failed to sync to S3: " .. err)
            return false, "S3 sync failed."
          end

          log.info("S3 sync completed successfully.")
          return true, "Files synced to S3."
        end
      },
      {
        name = "get_secret_value",
        description = "Retrieves a secret from AWS Secrets Manager.",
        depends_on = "check_aws_identity",
        command = function()
          log.info("Attempting to retrieve secret: " .. secret_id)
          local secret_string, err = aws.secretsmanager.get_secret({
            secret_id = secret_id,
            profile = aws_profile
          })

          if not secret_string then
            log.error("Failed to retrieve secret: " .. err)
            return false, "Secret retrieval failed."
          end

          log.info("Successfully retrieved secret!")
          -- IMPORTANT: Be careful not to print the actual secret in production logs.
          -- This is just for demonstration.
          log.info("Secret Value (first 10 chars): " .. string.sub(secret_string, 1, 10) .. "...")
          
          -- You can now use this secret in subsequent steps.
          return true, "Secret retrieved."
        end
      }
    }
  }
}
