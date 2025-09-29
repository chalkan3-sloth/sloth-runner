-- MODERN DSL ONLY
workflow.define("ci_cd_pipeline", {
    description = "Complete CI/CD pipeline with enhanced features",
    version = "2.0.0",
    
    -- Pipeline stages
    stages = {
        {
            name = "preparation",
            tasks = chain({
                "setup_workspace",
                "validate_environment",
                "load_secrets"
            })
        },
        
        {
            name = "build_and_test",
            tasks = workflow.parallel({
                "build_application",
                "run_tests",
                "security_scan",
                "quality_analysis"
            }, {
                max_workers = 4,
                fail_fast = true,
                timeout = "15m"
            })
        },
        
        {
            name = "deployment",
            condition = when("test.success && build.success")
                :then("deploy_staging")
                :else("notify_failure"),
            
            tasks = {
                {
                    name = "deploy_staging",
                    command = function(params, deps)
                        -- Use circuit breaker for external service
                        local deploy_result, deploy_err = flow.circuit_breaker("deployment_service", function()
                            return exec.run("kubectl apply -f k8s/staging/")
                        end)
                        
                        if deploy_err then
                            return false, "Deployment failed: " .. deploy_err
                        end
                        
                        -- Create checkpoint for potential rollback
                        task.checkpoint("pre_production_deploy", {
                            staging_version = params.version,
                            timestamp = os.time(),
                            config = params.deploy_config
                        })
                        
                        return true, "Staging deployment successful"
                    end,
                    
                    saga = {
                        compensation = function(ctx)
                            log.warn("Rolling back staging deployment")
                            return exec.run("kubectl rollout undo deployment/app -n staging")
                        end
                    }
                }
            }
        }
    },
    
    -- Enhanced error handling
    error_handling = {
        strategy = "retry_with_backoff",
        max_attempts = 3,
        
        on_failure = function(ctx, error)
            -- Advanced recovery with multiple strategies
            local recovery_result, recovery_err = error.try(
                function()
                    -- Primary recovery: restart from last checkpoint
                    local checkpoint = ctx.checkpoints["pre_production_deploy"]
                    if checkpoint then
                        return task.restore_checkpoint(checkpoint)
                    end
                    return false, "No checkpoint available"
                end,
                function(err)
                    log.error("Primary recovery failed", err)
                    -- Fallback: manual intervention notification
                    return notifications.send({
                        type = "critical",
                        message = "Pipeline failed, manual intervention required",
                        error = err,
                        pipeline = ctx.pipeline_id
                    })
                end,
                function()
                    log.info("Recovery completed")
                end
            )
            
            return recovery_result, recovery_err
        end
    },
    
    -- Resource management
    resources = {
        cpu = {
            request = "500m",
            limit = "2000m"
        },
        memory = {
            request = "1Gi",
            limit = "4Gi"
        },
        disk = {
            size = "20Gi",
            type = "ssd"
        }
    },
    
    -- Security policies
    security = {
        rbac = {
            roles = {"ci-runner", "deployer"},
            service_account = "sloth-runner-sa"
        },
        
        secrets = {
            mount_path = "/etc/secrets",
            keys = {"database_password", "api_key", "signing_cert"}
        },
        
        network_policy = {
            ingress = {
                from = {"ci-namespace", "monitoring-namespace"}
            },
            egress = {
                to = {"docker-registry", "kubernetes-api", "external-apis"}
            }
        }
    }
})
