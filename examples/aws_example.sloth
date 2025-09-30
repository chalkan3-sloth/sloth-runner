-- MODERN DSL ONLY - AWS Integration Example
-- Demonstrates AWS operations using Modern DSL

-- AWS S3 List Buckets task
local aws_list_buckets = task("aws_list_buckets")
    :description("List S3 buckets using AWS CLI")
    :command(function(params)
        log.info("🪣 Listing AWS S3 buckets...")
        
        -- Execute AWS CLI command with error handling
        local result = exec.run("aws s3 ls", {
            timeout = "30s",
            capture_output = true
        })
        
        if result.success then
            log.info("✅ Successfully listed S3 buckets")
            return true, result.output, {
                buckets = result.output,
                bucket_count = select(2, string.gsub(result.output or "", "\n", ""))
            }
        else
            return false, "Failed to list S3 buckets: " .. (result.error or "unknown error")
        end
    end)
    :timeout("60s")
    :retries(2, "exponential")
    :on_success(function(params, output)
        log.info("📊 Found " .. (output.bucket_count or 0) .. " S3 buckets")
    end)
    :build()

-- AWS EC2 Instance List task
local aws_list_ec2 = task("aws_list_ec2")
    :description("List EC2 instances")
    :command(function(params)
        log.info("🖥️  Listing AWS EC2 instances...")
        
        local result = exec.run("aws ec2 describe-instances --query 'Reservations[*].Instances[*].[InstanceId,State.Name,InstanceType]' --output table", {
            timeout = "45s",
            capture_output = true
        })
        
        if result.success then
            return true, result.output, {
                instances_info = result.output,
                command_executed = "describe-instances"
            }
        else
            return false, "Failed to list EC2 instances"
        end
    end)
    :timeout("90s")
    :build()

-- AWS CloudFormation Stacks task
local aws_list_cf_stacks = task("aws_list_cf_stacks")
    :description("List CloudFormation stacks")
    :depends_on({"aws_list_buckets"})
    :command(function(params, deps)
        log.info("☁️  Listing CloudFormation stacks...")
        
        local result = exec.run("aws cloudformation list-stacks --stack-status-filter CREATE_COMPLETE UPDATE_COMPLETE", {
            timeout = "30s",
            capture_output = true
        })
        
        if result.success then
            return true, result.output, {
                stacks_info = result.output,
                previous_buckets = deps.aws_list_buckets.bucket_count
            }
        else
            return false, "Failed to list CloudFormation stacks"
        end
    end)
    :timeout("60s")
    :build()

-- Modern Workflow Definition
workflow.define("aws_operations", {
    description = "AWS Operations Workflow - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"aws", "cloud", "s3", "ec2", "cloudformation", "modern-dsl"},
        created_at = os.date(),
        prerequisites = "AWS CLI configured with proper credentials"
    },
    
    tasks = {
        aws_list_buckets,
        aws_list_ec2,
        aws_list_cf_stacks
    },
    
    config = {
        timeout = "20m",
        retry_policy = "exponential",
        max_parallel_tasks = 2,
        cleanup_on_failure = true
    },
    
    on_start = function()
        log.info("🚀 Starting AWS operations workflow...")
        log.info("🔑 Ensure AWS credentials are configured")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ AWS operations workflow completed successfully!")
            log.info("📊 AWS resources have been listed and analyzed")
        else
            log.error("❌ AWS operations workflow failed!")
            log.warn("🔍 Check AWS CLI configuration and permissions")
        end
        return true
    end
})
