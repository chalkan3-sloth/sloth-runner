# ☁️ Multi-Cloud Excellence

Sloth Runner provides **comprehensive cloud provider support** with advanced automation capabilities across AWS, GCP, Azure, and DigitalOcean. Built-in modules enable infrastructure management, cost optimization, and security compliance.

## 🌟 Supported Cloud Providers

### ☁️ Amazon Web Services (AWS)
**Complete AWS ecosystem** integration with 200+ services:
- **EC2** - Virtual machines and auto-scaling
- **S3** - Object storage and lifecycle management  
- **RDS** - Managed databases (MySQL, PostgreSQL, Oracle)
- **Lambda** - Serverless functions
- **EKS** - Kubernetes clusters
- **CloudFormation** - Infrastructure as Code
- **IAM** - Identity and access management
- **VPC** - Virtual private clouds and networking

### 🌩️ Google Cloud Platform (GCP)
**Native GCP integration** with advanced features:
- **Compute Engine** - Virtual machines and instance groups
- **Cloud Storage** - Object storage with global distribution
- **Cloud SQL** - Managed relational databases
- **GKE** - Google Kubernetes Engine
- **Cloud Functions** - Event-driven serverless
- **Cloud Deployment Manager** - Infrastructure automation
- **Cloud IAM** - Identity and access management
- **VPC** - Virtual private cloud networking

### 🔷 Microsoft Azure
**Enterprise-grade Azure** support:
- **Virtual Machines** - Compute instances and scale sets
- **Storage Accounts** - Blob, file, and disk storage
- **Azure SQL Database** - Managed SQL services
- **AKS** - Azure Kubernetes Service
- **Azure Functions** - Serverless computing
- **ARM Templates** - Infrastructure deployment
- **Azure AD** - Identity services
- **Virtual Networks** - Software-defined networking

### 🌊 DigitalOcean
**Developer-friendly** cloud operations:
- **Droplets** - Virtual private servers
- **Spaces** - Object storage compatible with S3
- **Managed Databases** - PostgreSQL, MySQL, Redis
- **Kubernetes** - Managed container orchestration
- **App Platform** - Platform-as-a-Service
- **Load Balancers** - Traffic distribution
- **Firewalls** - Network security

## 🚀 AWS Advanced Module

### EC2 Management
```lua
local aws = require("aws")

task("manage_ec2_fleet")
    :command(function(params, deps)
        -- List instances with filtering
        local instances = aws.ec2.list_instances({
            filters = {
                ["instance-state-name"] = "running",
                ["tag:Environment"] = "production"
            }
        })
        
        -- Auto-scaling based on metrics
        for _, instance in ipairs(instances) do
            local metrics = aws.cloudwatch.get_metrics({
                instance_id = instance.id,
                metric_name = "CPUUtilization",
                period = "5m"
            })
            
            if metrics.average > 80 then
                log.info("High CPU detected, scaling up: " .. instance.id)
                aws.autoscaling.scale_out({
                    auto_scaling_group = instance.asg_name,
                    desired_capacity = instance.desired_capacity + 1
                })
            end
        end
        
        return true, "EC2 fleet management completed"
    end)
    :build()
```

### S3 Advanced Operations
```lua
task("s3_lifecycle_management")
    :command(function(params, deps)
        local aws = require("aws")
        
        -- Intelligent tiering and cleanup
        local buckets = aws.s3.list_buckets()
        
        for _, bucket in ipairs(buckets) do
            -- Analyze storage classes
            local analysis = aws.s3.analyze_storage_classes(bucket.name)
            
            -- Apply cost optimization
            if analysis.savings_potential > 1000 then  -- $1000 monthly savings
                aws.s3.apply_lifecycle_policy({
                    bucket = bucket.name,
                    rules = {
                        {
                            name = "auto_tiering",
                            transitions = {
                                {days = 30, storage_class = "STANDARD_IA"},
                                {days = 90, storage_class = "GLACIER"},
                                {days = 365, storage_class = "DEEP_ARCHIVE"}
                            }
                        }
                    }
                })
                log.info("Applied lifecycle policy to: " .. bucket.name)
            end
        end
        
        return true, "S3 optimization completed"
    end)
    :build()
```

### Lambda Deployment
```lua
task("deploy_lambda_functions")
    :command(function(params, deps)
        local aws = require("aws")
        
        -- Deploy multiple functions with dependencies
        local functions = {
            {
                name = "api-handler",
                runtime = "nodejs18.x",
                code = "build/api-handler.zip",
                environment = {
                    NODE_ENV = "production",
                    DB_CONNECTION = aws.ssm.get_parameter("/app/db/connection")
                }
            },
            {
                name = "data-processor", 
                runtime = "python3.9",
                code = "build/data-processor.zip",
                memory = 512,
                timeout = 300
            }
        }
        
        for _, func in ipairs(functions) do
            -- Create or update function
            local result = aws.lambda.deploy_function(func)
            
            -- Set up triggers
            if func.name == "api-handler" then
                aws.apigateway.create_integration({
                    api_id = params.api_gateway_id,
                    lambda_function = result.function_arn
                })
            end
            
            log.info("Deployed function: " .. func.name .. " (ARN: " .. result.function_arn .. ")")
        end
        
        return true, "Lambda deployment completed"
    end)
    :build()
```

## 🌩️ GCP Advanced Module

### Compute Engine Automation
```lua
task("gcp_compute_management")
    :command(function(params, deps)
        local gcp = require("gcp")
        
        -- Create managed instance group with auto-scaling
        local instance_template = gcp.compute.create_instance_template({
            name = "web-server-template-v2",
            machine_type = "e2-medium",
            source_image = "projects/ubuntu-os-cloud/global/images/family/ubuntu-2004-lts",
            startup_script = [[
                #!/bin/bash
                apt-get update
                apt-get install -y nginx
                systemctl start nginx
            ]],
            tags = {"web-server", "http-server"},
            metadata = {
                ["startup-script-url"] = "gs://my-scripts/startup.sh"
            }
        })
        
        -- Create managed instance group
        local mig = gcp.compute.create_managed_instance_group({
            name = "web-servers",
            base_instance_name = "web-server",
            template = instance_template.self_link,
            target_size = 3,
            zone = "us-central1-a"
        })
        
        -- Configure auto-scaling
        gcp.compute.create_autoscaler({
            name = "web-servers-autoscaler",
            target = mig.self_link,
            autoscaling_policy = {
                min_replicas = 2,
                max_replicas = 10,
                cpu_utilization = {
                    target = 0.7
                }
            }
        })
        
        return true, "GCP compute resources created"
    end)
    :build()
```

### GKE Cluster Management
```lua
task("manage_gke_cluster")
    :command(function(params, deps)
        local gcp = require("gcp")
        local k8s = require("kubernetes")
        
        -- Create GKE cluster with advanced configuration
        local cluster = gcp.gke.create_cluster({
            name = "production-cluster",
            location = "us-central1",
            initial_node_count = 3,
            
            node_config = {
                machine_type = "e2-standard-4",
                disk_size_gb = 100,
                oauth_scopes = {
                    "https://www.googleapis.com/auth/devstorage.read_only",
                    "https://www.googleapis.com/auth/logging.write",
                    "https://www.googleapis.com/auth/monitoring"
                }
            },
            
            addons_config = {
                horizontal_pod_autoscaling = {disabled = false},
                http_load_balancing = {disabled = false},
                network_policy_config = {disabled = false}
            },
            
            network_policy = {
                enabled = true,
                provider = "CALICO"
            }
        })
        
        -- Configure kubectl context
        gcp.gke.get_credentials({
            cluster_name = cluster.name,
            location = cluster.location
        })
        
        -- Deploy application
        k8s.apply_manifest("k8s/production/")
        
        return true, "GKE cluster configured and application deployed"
    end)
    :build()
```

## 🔷 Azure Advanced Module

### Virtual Machine Scale Sets
```lua
task("azure_vmss_management")
    :command(function(params, deps)
        local azure = require("azure")
        
        -- Create VM Scale Set with custom image
        local vmss = azure.compute.create_vmss({
            name = "web-servers-vmss",
            resource_group = "production-rg",
            location = "East US",
            
            sku = {
                name = "Standard_B2s",
                tier = "Standard",
                capacity = 3
            },
            
            virtual_machine_profile = {
                os_profile = {
                    computer_name_prefix = "web",
                    admin_username = "azureuser",
                    custom_data = base64.encode([[
                        #!/bin/bash
                        apt-get update
                        apt-get install -y docker.io
                        docker run -d -p 80:80 nginx
                    ]])
                },
                
                storage_profile = {
                    image_reference = {
                        publisher = "Canonical",
                        offer = "UbuntuServer",
                        sku = "18.04-LTS",
                        version = "latest"
                    }
                },
                
                network_profile = {
{% raw %}                    network_interface_configurations = {{
                        name = "web-nic",
                        primary = true,
                        ip_configurations = {{
                            name = "internal",
                            subnet = {
                                id = "/subscriptions/.../subnets/web-subnet"
                            },
                            load_balancer_backend_address_pools = {{
                                id = "/subscriptions/.../backendAddressPools/web-backend"
                            }}
                        }}{% endraw %}
                    }}
                }
            },
            
            -- Auto-scaling configuration
            upgrade_policy = {
                mode = "Rolling",
                rolling_upgrade_policy = {
                    max_batch_instance_percent = 20,
                    max_unhealthy_instance_percent = 20,
                    max_unhealthy_upgraded_instance_percent = 20,
                    pause_time_between_batches = "PT0S"
                }
            }
        })
        
        -- Configure auto-scaling rules
        azure.monitor.create_autoscale_settings({
            name = "web-servers-autoscale",
            target_resource_id = vmss.id,
{% raw %}            profiles = {{{% endraw %}
                name = "default",
                capacity = {
                    minimum = "2",
                    maximum = "10", 
                    default = "3"
                },
                rules = {
                    {
                        metric_trigger = {
                            metric_name = "Percentage CPU",
                            metric_namespace = "microsoft.compute/virtualmachinescalesets",
                            time_grain = "PT1M",
                            statistic = "Average",
                            time_window = "PT5M",
                            time_aggregation = "Average",
                            operator = "GreaterThan",
                            threshold = 75
                        },
                        scale_action = {
                            direction = "Increase",
                            type = "ChangeCount",
                            value = "1",
                            cooldown = "PT5M"
                        }
                    }
                }
            }}
        })
        
        return true, "Azure VMSS created and configured"
    end)
    :build()
```

## 🌊 DigitalOcean Advanced Module

### Kubernetes Cluster with Apps
```lua
task("do_kubernetes_deployment")
    :command(function(params, deps)
        local do_client = require("digitalocean")
        
        -- Create Kubernetes cluster
        local cluster = do_client.kubernetes.create_cluster({
            name = "production-k8s",
            region = "nyc1",
            version = "1.25.4-do.0",
            
{% raw %}            node_pools = {{{% endraw %}
                size = "s-2vcpu-2gb",
                count = 3,
                name = "worker-pool",
                tags = {"production", "web"},
                auto_scale = true,
                min_nodes = 2,
                max_nodes = 8
            }},
            
            maintenance_policy = {
                start_time = "00:00",
                day = "sunday"
            },
            
            auto_upgrade = true,
            surge_upgrade = true,
            ha = true
        })
        
        -- Wait for cluster to be ready
        do_client.kubernetes.wait_for_cluster(cluster.id, "running", 600)
        
        -- Get kubeconfig
        local kubeconfig = do_client.kubernetes.get_kubeconfig(cluster.id)
        fs.write_file("~/.kube/do-config", kubeconfig)
        
        -- Deploy applications
        local k8s = require("kubernetes")
        k8s.set_context("do-production")
        
        -- Create namespace and deploy app
        k8s.create_namespace("production")
        k8s.apply_yaml([[
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
  namespace: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web-app
  template:
    metadata:
      labels:
        app: web-app
    spec:
      containers:
      - name: web
        image: nginx:latest
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: web-service
  namespace: production
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 80
  selector:
    app: web-app
        ]])
        
        -- Create DigitalOcean Load Balancer
        local lb = do_client.load_balancers.create({
            name = "web-lb",
            algorithm = "round_robin",
            status = "active",
            
{% raw %}            forwarding_rules = {{{% endraw %}
                entry_protocol = "http",
                entry_port = 80,
                target_protocol = "http",
                target_port = 80,
                certificate_id = "",
                tls_passthrough = false
            }},
            
            health_check = {
                protocol = "http",
                port = 80,
                path = "/health",
                check_interval_seconds = 10,
                response_timeout_seconds = 5,
                unhealthy_threshold = 3,
                healthy_threshold = 2
            },
            
            sticky_sessions = {
                type = "cookies",
                cookie_name = "lb",
                cookie_ttl_seconds = 300
            },
            
            region = "nyc1",
            tag = "production",
            droplet_ids = {},  -- Will be populated by k8s service
            redirect_http_to_https = false,
            enable_proxy_protocol = false
        })
        
        return true, "DigitalOcean Kubernetes cluster deployed"
    end)
    :build()
```

## 🔒 Security & Compliance

### Multi-Cloud Security Scanning
```lua
task("multi_cloud_security_scan")
    :command(function(params, deps)
        local security = require("security")
        local results = {}
        
        -- AWS Security Assessment
        local aws_scan = security.scan_aws({
            regions = {"us-east-1", "us-west-2"},
            services = {"ec2", "s3", "rds", "iam"},
            compliance_frameworks = {"SOC2", "PCI-DSS"}
        })
        
        -- GCP Security Assessment  
        local gcp_scan = security.scan_gcp({
            projects = {"prod-project", "staging-project"},
            services = {"compute", "storage", "iam"},
            compliance_frameworks = {"ISO27001", "SOC2"}
        })
        
        -- Azure Security Assessment
        local azure_scan = security.scan_azure({
            subscriptions = {"prod-subscription"},
            resource_groups = {"production-rg"},
            compliance_frameworks = {"HIPAA", "SOC2"}
        })
        
        -- Consolidate results
        results.aws = aws_scan
        results.gcp = gcp_scan
        results.azure = azure_scan
        
        -- Generate compliance report
        local report = security.generate_compliance_report({
            scans = results,
            format = "pdf",
            include_recommendations = true
        })
        
        -- Send to security team
        notifications.email.send({
            to = ["security@company.com"],
            subject = "Multi-Cloud Security Assessment Report",
            attachments = {report.file_path}
        })
        
        return true, "Security scan completed", results
    end)
    :build()
```

## 💰 Cost Optimization

### Cross-Cloud Cost Analysis
```lua
task("multi_cloud_cost_optimization")
    :command(function(params, deps)
        local cost_optimizer = require("cost_optimizer")
        
        -- Gather cost data from all providers
        local costs = {
            aws = aws.billing.get_costs({period = "30d", breakdown = "service"}),
            gcp = gcp.billing.get_costs({period = "30d", breakdown = "service"}),
            azure = azure.billing.get_costs({period = "30d", breakdown = "service"}),
            digitalocean = digitalocean.billing.get_costs({period = "30d"})
        }
        
        -- Analyze usage patterns
        local analysis = cost_optimizer.analyze_multi_cloud({
            costs = costs,
            utilization_threshold = 0.3,  -- 30% utilization minimum
            savings_threshold = 100       -- $100 minimum savings
        })
        
        -- Generate optimization recommendations
        local recommendations = cost_optimizer.generate_recommendations({
            analysis = analysis,
            strategies = {
                "rightsizing",
                "reserved_instances",
                "spot_instances", 
                "storage_tiering",
                "resource_consolidation"
            }
        })
        
        -- Auto-apply low-risk optimizations
        for _, rec in ipairs(recommendations) do
            if rec.risk_level == "low" and rec.estimated_savings > 50 then
                log.info("Auto-applying optimization: " .. rec.description)
                cost_optimizer.apply_recommendation(rec)
            end
        end
        
        -- Generate cost report
        local report = cost_optimizer.generate_report({
            costs = costs,
            recommendations = recommendations,
            format = "html"
        })
        
        return true, "Cost optimization completed", {
            total_monthly_cost = analysis.total_cost,
            potential_savings = analysis.potential_savings,
            optimizations_applied = #recommendations
        }
    end)
    :build()
```

## 🎯 Best Practices

### Infrastructure as Code
1. **Use version control** for all cloud configurations
2. **Implement proper tagging** strategies across providers
3. **Regular backup** and disaster recovery testing
4. **Automate security** scanning and compliance checks
5. **Monitor costs** and optimize regularly

### Security Guidelines
1. **Enable multi-factor authentication** on all cloud accounts
2. **Use least privilege** access principles
3. **Regular security** assessments and penetration testing
4. **Encrypt data** at rest and in transit
5. **Implement proper** logging and monitoring

### Performance Optimization
1. **Right-size resources** based on actual usage
2. **Use auto-scaling** for variable workloads
3. **Implement caching** strategies appropriately
4. **Monitor performance** metrics continuously
5. **Regular architecture** reviews and optimizations

---

Multi-Cloud Excellence with Sloth Runner enables organizations to **leverage the best of each cloud provider** while maintaining consistent automation, security, and cost optimization across their entire cloud infrastructure! ☁️✨