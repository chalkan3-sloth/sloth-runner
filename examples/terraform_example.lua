-- MODERN DSL ONLY - Terraform Lifecycle Pipeline
-- Complete Modern DSL implementation for Terraform workflows

local tf_workdir = "./examples/terraform"

-- Task 1: Initialize Terraform project
local init_task = task("init")
    :description("Initializes the Terraform project with Modern DSL")
    :command(function(params)
        log.info("🔧 Running terraform init...")
        local result = terraform.init({workdir = tf_workdir})
        if not result.success then
            log.error("Terraform init failed: " .. result.stderr)
            return false, "Terraform init failed."
        end
        log.info("✅ Terraform init successful.")
        return true, "Terraform initialized.", { workdir = tf_workdir }
    end)
    :timeout("5m")
    :retries(2, "exponential")
    :build()

-- Task 2: Create Terraform plan
local plan_task = task("plan")
    :description("Creates a Terraform execution plan")
    :depends_on({"init"})
    :command(function(params, deps)
        log.info("📋 Running terraform plan...")
        local result = terraform.plan({workdir = tf_workdir})
        if not result.success then
            log.error("Terraform plan failed: " .. result.stderr)
            return false, "Terraform plan failed."
        end
        log.info("✅ Terraform plan successful.")
        print(result.stdout)
        return true, "Terraform plan created.", { 
            plan_output = result.stdout,
            workdir = tf_workdir 
        }
    end)
    :timeout("10m")
    :build()

-- Task 3: Apply Terraform plan
local apply_task = task("apply")
    :description("Applies the Terraform plan")
    :depends_on({"plan"})
    :command(function(params, deps)
        log.info("🚀 Running terraform apply...")
        local result = terraform.apply({workdir = tf_workdir, auto_approve = true})
        if not result.success then
            log.error("Terraform apply failed: " .. result.stderr)
            return false, "Terraform apply failed."
        end
        log.info("✅ Terraform apply successful.")
        return true, "Terraform apply complete.", { 
            apply_output = result.stdout 
        }
    end)
    :timeout("30m")
    :on_success(function(params, output)
        log.info("🎉 Infrastructure deployed successfully!")
    end)
    :build()

-- Task 4: Get Terraform outputs
local get_output_task = task("get_output")
    :description("Reads the output variables from Terraform")
    :depends_on({"apply"})
    :command(function(params, deps)
        log.info("📤 Reading Terraform output...")
        local filename, err = terraform.output({workdir = tf_workdir, name = "report_filename"})
        if not filename then
            log.error("Failed to get Terraform output: " .. err)
            return false, "Terraform output failed."
        end
        
        log.info("📄 Got filename from output: " .. filename)
        
        -- Read the content of the file created by Terraform
        local content, read_err = fs.read(filename)
        if read_err then
            log.error("Failed to read the report file: " .. read_err)
            return false, "Could not read artifact."
        end

        log.info("✅ Successfully read content from Terraform-generated file:")
        print("--- Report Content ---")
        print(content)
        print("----------------------")

        return true, "Terraform output processed.", { 
            output_file = filename,
            content = content 
        }
    end)
    :artifacts({tf_workdir .. "/terraform.tfstate"})
    :build()

-- Task 5: Destroy resources (optional cleanup)
local destroy_task = task("destroy")
    :description("Destroys the Terraform-managed resources")
    :depends_on({"get_output"})
    :command(function(params, deps)
        log.warn("💥 Running terraform destroy...")
        local result = terraform.destroy({workdir = tf_workdir, auto_approve = true})
        if not result.success then
            log.error("Terraform destroy failed: " .. result.stderr)
            return false, "Terraform destroy failed."
        end
        log.info("✅ Terraform destroy successful.")
        return true, "Terraform resources destroyed.", {
            destroyed_at = os.time()
        }
    end)
    :timeout("20m")
    :on_success(function(params, output)
        log.info("🧹 Infrastructure cleanup completed!")
    end)
    :build()

-- Modern Workflow Definition
workflow.define("terraform_lifecycle", {
    description = "Complete Terraform lifecycle management - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"terraform", "infrastructure", "iac", "modern-dsl"},
        complexity = "advanced",
        estimated_duration = "45m"
    },
    
    tasks = {
        init_task,
        plan_task,
        apply_task,
        get_output_task,
        destroy_task
    },
    
    config = {
        timeout = "60m",
        retry_policy = "exponential",
        max_parallel_tasks = 1, -- Sequential execution for Terraform
        fail_fast = true
    },
    
    on_start = function()
        log.info("🚀 Starting Terraform lifecycle pipeline...")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("🎉 Terraform pipeline completed successfully!")
            log.info("📊 Infrastructure state updated")
        else
            log.error("❌ Terraform pipeline failed!")
            log.warn("🔍 Check terraform state and resolve issues")
        end
        return true
    end
})
