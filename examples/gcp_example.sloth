-- MODERN DSL ONLY - Google Cloud Platform Example
-- Demonstrates GCP operations using Modern DSL

-- GCP Project Info task
local gcp_project_info = task("gcp_project_info")
    :description("Get current GCP project information")
    :command(function(params)
        log.info("🌐 Getting GCP project information...")
        
        local result = exec.run("gcloud config get-value project", {
            timeout = "30s",
            capture_output = true
        })
        
        if result.success then
            local project_id = string.gsub(result.output or "", "%s+", "")
            log.info("📋 Current project: " .. project_id)
            
            return true, result.output, {
                project_id = project_id,
                config_type = "project_info"
            }
        else
            return false, "Failed to get project info"
        end
    end)
    :timeout("60s")
    :retries(2, "exponential")
    :build()

-- GCP Compute Instances task
local gcp_list_instances = task("gcp_list_compute_instances")
    :description("List GCP Compute Engine instances")
    :depends_on({"gcp_project_info"})
    :command(function(params, deps)
        log.info("🖥️  Listing GCP Compute Engine instances...")
        
        local project_id = deps.gcp_project_info.project_id
        log.info("🔍 Searching in project: " .. project_id)
        
        local result = exec.run("gcloud compute instances list --format='table(name,zone,machineType,status)'", {
            timeout = "90s",
            capture_output = true
        })
        
        if result.success then
            return true, result.output, {
                instances = result.output,
                project_context = project_id
            }
        else
            return false, "Failed to list compute instances"
        end
    end)
    :timeout("120s")
    :build()

-- GCP Storage Buckets task
local gcp_list_buckets = task("gcp_list_storage_buckets")
    :description("List GCS storage buckets")
    :command(function(params)
        log.info("🪣 Listing GCS storage buckets...")
        
        local result = exec.run("gsutil ls", {
            timeout = "45s",
            capture_output = true
        })
        
        if result.success then
            local bucket_count = select(2, string.gsub(result.output or "", "gs://", ""))
            return true, result.output, {
                buckets = result.output,
                bucket_count = bucket_count
            }
        else
            return false, "Failed to list storage buckets"
        end
    end)
    :timeout("90s")
    :on_success(function(params, output)
        log.info("📊 Found " .. (output.bucket_count or 0) .. " storage buckets")
    end)
    :build()

-- GCP Kubernetes Clusters task
local gcp_list_gke = task("gcp_list_gke_clusters")
    :description("List GKE clusters")
    :depends_on({"gcp_project_info"})
    :command(function(params, deps)
        log.info("☸️  Listing GKE clusters...")
        
        local result = exec.run("gcloud container clusters list --format='table(name,location,status,nodeCount)'", {
            timeout = "60s",
            capture_output = true
        })
        
        if result.success then
            return true, result.output, {
                clusters = result.output,
                cluster_type = "gke",
                project_id = deps.gcp_project_info.project_id
            }
        else
            return false, "Failed to list GKE clusters"
        end
    end)
    :timeout("120s")
    :build()

-- Modern Workflow Definition
workflow.define("gcp_operations", {
    description = "Google Cloud Platform Operations - Modern DSL",
    version = "2.0.0",
    
    metadata = {
        author = "Sloth Runner Team",
        tags = {"gcp", "google-cloud", "compute", "storage", "gke", "modern-dsl"},
        created_at = os.date(),
        prerequisites = "gcloud CLI authenticated and project configured"
    },
    
    tasks = {
        gcp_project_info,
        gcp_list_instances,
        gcp_list_buckets,
        gcp_list_gke
    },
    
    config = {
        timeout = "30m",
        retry_policy = "exponential",
        max_parallel_tasks = 3,
        cleanup_on_failure = true
    },
    
    on_start = function()
        log.info("🚀 Starting GCP operations workflow...")
        log.info("🔑 Ensure gcloud is authenticated and project is set")
        return true
    end,
    
    on_complete = function(success, results)
        if success then
            log.info("✅ GCP operations workflow completed successfully!")
            log.info("☁️  GCP resources have been inventoried")
        else
            log.error("❌ GCP operations workflow failed!")
            log.warn("🔍 Check gcloud authentication and project configuration")
        end
        return true
    end
})
