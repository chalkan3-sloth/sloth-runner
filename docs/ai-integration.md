# 🤖 AI/ML Integration

Sloth Runner provides **built-in artificial intelligence** capabilities that enable smart automation, intelligent decision making, and advanced text processing for modern workflows.

## 🧠 Core AI Features

### OpenAI Integration
**Direct integration** with OpenAI's powerful language models:
- **Text generation** and completion
- **Code generation** assistance  
- **Intelligent analysis** of logs and data
- **Automated decision** making
- **Natural language** task descriptions

### Smart Automation
**AI-powered automation** that learns from your workflows:
- **Pattern recognition** in task failures
- **Automatic retry** strategy suggestions
- **Performance optimization** recommendations
- **Anomaly detection** in execution patterns

## 🚀 Getting Started

### Prerequisites
```bash
# Set OpenAI API key
export OPENAI_API_KEY="your-api-key-here"

# Or use configuration file
echo "openai_api_key: your-api-key" > ai-config.yaml
```

### Basic AI Module Usage
```lua
-- Load AI module
local ai = require("ai")

-- Simple text completion
local response = ai.openai.complete({
    prompt = "Generate a bash script to deploy a Node.js application",
    max_tokens = 200,
    temperature = 0.7
})

log.info("Generated script: " .. response.text)
```

## 🎯 AI-Powered Tasks

### Intelligent Code Generation
```lua
task("generate_dockerfile")
    :description("AI-generated Dockerfile for the application")
    :command(function(params, deps)
        local ai = require("ai")
        
        -- Analyze project structure
        local project_files = exec.run("find . -type f -name '*.js' -o -name '*.py' -o -name '*.go'").stdout
        
        -- Generate appropriate Dockerfile
        local dockerfile_content = ai.openai.complete({
            prompt = "Create a Dockerfile for this project with files: " .. project_files,
            max_tokens = 300,
            temperature = 0.3  -- Lower temperature for more deterministic code
        })
        
        -- Write generated Dockerfile
        local file = io.open("Dockerfile", "w")
        file:write(dockerfile_content.text)
        file:close()
        
        return true, "Dockerfile generated successfully", {
            dockerfile_size = #dockerfile_content.text,
            ai_confidence = dockerfile_content.confidence
        }
    end)
    :build()
```

### Smart Log Analysis
```lua
task("analyze_logs")
    :description("AI-powered log analysis for error detection")
    :command(function(params, deps)
        local ai = require("ai")
        
        -- Get recent logs
        local logs = exec.run("tail -100 /var/log/application.log").stdout
        
        -- AI analysis
        local analysis = ai.openai.complete({
            prompt = "Analyze these application logs and identify any errors or issues:\n" .. logs,
            max_tokens = 150
        })
        
        -- Make decisions based on analysis
        local decision = ai.decide({
            analysis = analysis.text,
            threshold = 0.8
        })
        
        if decision.severity == "high" then
            -- Trigger alert workflow
            notifications.send_alert("Critical issue detected: " .. decision.issue)
        end
        
        return true, "Log analysis completed", {
            analysis = analysis.text,
            severity = decision.severity,
            recommendations = decision.actions
        }
    end)
    :build()
```

### Intelligent Environment Selection
```lua
task("smart_deploy")
    :description("AI chooses the best deployment environment")
    :command(function(params, deps)
        local ai = require("ai")
        
        -- Gather deployment context
        local context = {
            time_of_day = os.date("%H"),
            day_of_week = os.date("%w"),
            recent_deployments = state.get("recent_deployments") or {},
            current_load = monitoring.get_load_metrics(),
            test_results = deps.run_tests or {}
        }
        
        -- AI decision making
        local deployment_plan = ai.decide({
            context = context,
            options = {"staging", "production", "canary"},
            criteria = {
                "minimize_risk",
                "optimize_performance", 
                "consider_load"
            }
        })
        
        log.info("AI recommends deployment to: " .. deployment_plan.environment)
        log.info("Reasoning: " .. deployment_plan.reasoning)
        
        -- Execute recommended deployment
        return exec.run("deploy.sh " .. deployment_plan.environment)
    end)
    :build()
```

## 🔧 AI Module API Reference

### OpenAI Integration

#### Text Completion
```lua
local ai = require("ai")

-- Basic completion
local result = ai.openai.complete("Explain CI/CD best practices")

-- Advanced completion with parameters
local result = ai.openai.complete({
    prompt = "Write a deployment script",
    max_tokens = 500,
    temperature = 0.5,
    top_p = 0.9,
    frequency_penalty = 0.1,
    presence_penalty = 0.1
})
```

#### Chat Completions
```lua
-- Multi-turn conversation
local chat_result = ai.openai.chat({
    messages = {
        {role = "system", content = "You are a DevOps expert"},
        {role = "user", content = "How do I optimize my Docker builds?"}
    },
    model = "gpt-4",
    max_tokens = 300
})
```

#### Code Generation
```lua
-- Specialized code generation
local code = ai.openai.generate_code({
    language = "bash",
    description = "Script to backup PostgreSQL database",
    parameters = {
        database_name = "myapp_db",
        backup_location = "/backups/"
    }
})
```

### Decision Making Engine

#### Smart Decisions
```lua
-- AI-powered decision making
local decision = ai.decide({
    situation = "Database CPU usage is 85%",
    options = {"scale_up", "optimize_queries", "add_replica"},
    constraints = {
        budget = "limited",
        downtime = "not_allowed"
    },
    history = previous_decisions
})

-- Returns structured decision with reasoning
print(decision.choice)      -- "add_replica"
print(decision.confidence)  -- 0.87
print(decision.reasoning)   -- "Adding replica provides..."
```

#### Pattern Recognition
```lua
-- Detect patterns in workflow data
local patterns = ai.analyze_patterns({
    data = workflow_execution_history,
    pattern_types = {"failure_correlation", "performance_trends", "resource_usage"},
    time_window = "30d"
})

for _, pattern in ipairs(patterns) do
    log.info("Detected pattern: " .. pattern.description)
    log.info("Confidence: " .. pattern.confidence)
    log.info("Recommendation: " .. pattern.recommendation)
end
```

## 🎯 Advanced Use Cases

### AI-Driven CI/CD Pipeline
```lua
-- Complete AI-powered CI/CD workflow
workflow.define("ai_cicd_pipeline", {
    description = "AI-enhanced CI/CD pipeline with smart decisions",
    
    tasks = {
        -- AI code review
        task("ai_code_review")
            :command(function(params, deps)
                local ai = require("ai")
                local git = require("git")
                
                -- Get diff for review
                local diff = git.get_diff("HEAD~1", "HEAD")
                
                -- AI code review
                local review = ai.openai.complete({
                    prompt = "Review this code diff for security issues, bugs, and best practices:\n" .. diff,
                    max_tokens = 400
                })
                
                -- Fail if critical issues found
                if string.find(review.text:lower(), "critical") then
                    return false, "Critical issues found in code review"
                end
                
                return true, "Code review passed", {review = review.text}
            end)
            :build(),
            
        -- Smart test selection
        task("ai_test_selection")
            :depends_on({"ai_code_review"})
            :command(function(params, deps)
                local ai = require("ai")
                
                -- Analyze changed files
                local changed_files = git.get_changed_files()
                
                -- AI decides which tests to run
                local test_plan = ai.decide({
                    changed_files = changed_files,
                    available_tests = {"unit", "integration", "e2e", "performance"},
                    time_budget = "10m"
                })
                
                -- Run selected tests
                for _, test_type in ipairs(test_plan.selected_tests) do
                    exec.run("npm run test:" .. test_type)
                end
                
                return true, "Smart test execution completed"
            end)
            :build(),
            
        -- Intelligent deployment strategy
        task("ai_deployment")
            :depends_on({"ai_test_selection"})
            :command(function(params, deps)
                local ai = require("ai")
                
                -- Deployment decision factors
                local factors = {
                    test_results = deps.ai_test_selection,
                    current_production_health = monitoring.get_health(),
                    deployment_history = state.get("deployment_history"),
                    time_context = {
                        hour = tonumber(os.date("%H")),
                        day_of_week = os.date("%A")
                    }
                }
                
                -- AI deployment strategy
                local strategy = ai.decide({
                    context = factors,
                    strategies = {"blue_green", "canary", "rolling", "immediate"},
                    risk_tolerance = "medium"
                })
                
                log.info("AI selected deployment strategy: " .. strategy.choice)
                log.info("Reasoning: " .. strategy.reasoning)
                
                -- Execute deployment
                return exec.run("deploy.sh --strategy=" .. strategy.choice)
            end)
            :build()
    }
})
```

### AI-Powered Infrastructure Management
```lua
task("ai_infrastructure_optimization")
    :description("AI optimizes infrastructure based on usage patterns")
    :command(function(params, deps)
        local ai = require("ai")
        local aws = require("aws")
        
        -- Gather infrastructure metrics
        local metrics = {
            ec2_utilization = aws.ec2.get_utilization_metrics("7d"),
            rds_performance = aws.rds.get_performance_metrics("7d"),
            costs = aws.billing.get_costs("30d")
        }
        
        -- AI analysis and recommendations
        local optimization = ai.openai.complete({
            prompt = "Analyze these AWS infrastructure metrics and provide optimization recommendations:\n" .. 
                     json.encode(metrics),
            max_tokens = 500
        })
        
        -- Parse recommendations and create action plan
        local recommendations = ai.parse_recommendations(optimization.text)
        
        for _, rec in ipairs(recommendations) do
            if rec.type == "scale_down" and rec.confidence > 0.8 then
                log.info("AI recommends scaling down: " .. rec.resource)
                -- Auto-execute high-confidence recommendations
                aws.ec2.modify_instance(rec.resource, {instance_type = rec.new_size})
            end
        end
        
        return true, "Infrastructure optimization completed", {
            recommendations = recommendations,
            estimated_savings = optimization.estimated_savings
        }
    end)
    :build()
```

## 🔒 Security & Best Practices

### API Key Management
```lua
-- Secure API key handling
local ai = require("ai")

-- Load from environment (recommended)
ai.configure({
    api_key = os.getenv("OPENAI_API_KEY"),
    timeout = 30,
    retry_attempts = 3
})

-- Or from encrypted configuration
ai.configure({
    config_file = "encrypted-ai-config.yaml",
    encryption_key = os.getenv("CONFIG_ENCRYPTION_KEY")
})
```

### Rate Limiting
```lua
-- Built-in rate limiting
ai.configure({
    rate_limit = {
        requests_per_minute = 60,
        tokens_per_minute = 40000
    }
})
```

### Data Privacy
```lua
-- Sanitize sensitive data before AI processing
task("ai_log_analysis")
    :command(function(params, deps)
        local ai = require("ai")
        
        -- Remove sensitive information
        local sanitized_logs = ai.sanitize_data(raw_logs, {
            remove_patterns = {
                "password=.*",
                "token=.*",
                "api_key=.*"
            }
        })
        
        local analysis = ai.openai.complete({
            prompt = "Analyze these sanitized logs: " .. sanitized_logs
        })
        
        return true, "Analysis completed safely"
    end)
    :build()
```

## 📊 Monitoring AI Usage

### Token Usage Tracking
```lua
-- Monitor AI API usage
local usage = ai.get_usage_stats()
print("Tokens used today: " .. usage.tokens_today)
print("API calls made: " .. usage.calls_today)
print("Estimated cost: $" .. usage.estimated_cost)
```

### Performance Metrics
```lua
-- AI response time monitoring
local start_time = os.time()
local result = ai.openai.complete("Generate deployment script")
local duration = os.time() - start_time

state.set("ai_response_time", duration)
```

## 🎨 Custom AI Models

### Local Model Integration
```lua
-- Use local AI models for sensitive data
ai.configure({
    provider = "local",
    model_path = "/models/custom-code-model",
    device = "gpu"  -- or "cpu"
})

local result = ai.local.complete("Generate Kubernetes deployment")
```

### Custom Fine-Tuned Models
```lua
-- Use organization-specific fine-tuned models
ai.configure({
    provider = "openai",
    model = "ft:gpt-3.5-turbo:company:devops-model:abc123"
})
```

## 🔮 Future AI Features

The AI integration is continuously evolving with planned features:
- **Multi-modal AI** (text + images + code)
- **Workflow learning** from execution patterns
- **Predictive failure** detection
- **Auto-healing** infrastructure
- **Natural language** workflow creation

---

AI integration transforms Sloth Runner into an **intelligent automation platform** that doesn't just execute tasks, but thinks about the best way to accomplish your goals! 🤖✨