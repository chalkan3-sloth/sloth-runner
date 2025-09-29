-- High Priority Modules Integration Showcase
-- Demonstrates Network, System, Security, Queue, and Observability modules working together

print("🚀 HIGH-PRIORITY MODULES INTEGRATION SHOWCASE")
print("=" .. string.rep("=", 60))
print("Demonstrating enterprise-grade monitoring, security, and infrastructure automation")

-- Initialize observability for the entire operation
local main_trace = observability.start_trace("infrastructure-health-check", {
    operation = "full_system_audit",
    priority = "high",
    automated = "true"
})

print("\n🎯 Started distributed trace:", main_trace)

-- 1. SYSTEM HEALTH ASSESSMENT
print("\n💻 PHASE 1: SYSTEM HEALTH ASSESSMENT")
print("-" .. string.rep("-", 45))

local system_span = observability.start_span(main_trace, "system-health-assessment", "", {
    component = "system",
    check_type = "comprehensive"
})

-- Get comprehensive system snapshot
local perf_snapshot = system.performance_snapshot()
local health_check = system.system_health()

if perf_snapshot and health_check then
    print("📊 System Performance Snapshot:")
    print("   CPU Usage:", string.format("%.1f%%", perf_snapshot.cpu_percent or 0))
    print("   Memory Usage:", string.format("%.1f%% (%.1f/%.1fGB)", 
        perf_snapshot.memory.percent or 0,
        perf_snapshot.memory.used_gb or 0,
        perf_snapshot.memory.total_gb or 0))
    
    if perf_snapshot.disk then
        print("   Disk Usage:", string.format("%.1f%%", perf_snapshot.disk.percent or 0))
    end
    
    print("🏥 Health Score:", string.format("%.0f/100", health_check.score or 0))
    print("📈 Status:", health_check.status or "unknown")
    
    -- Record metrics
    observability.gauge("system_health_score", health_check.score or 0, {
        status = health_check.status or "unknown"
    })
    
    observability.gauge("cpu_usage", perf_snapshot.cpu_percent or 0, {
        host = "current_system"
    })
    
    observability.gauge("memory_usage", perf_snapshot.memory.percent or 0, {
        host = "current_system"
    })
end

observability.add_span_event(system_span, "health-check-completed", {
    score = tostring(health_check.score or 0),
    status = health_check.status or "unknown"
})

observability.end_span(system_span, "completed")

-- 2. NETWORK CONNECTIVITY ASSESSMENT
print("\n🌐 PHASE 2: NETWORK CONNECTIVITY ASSESSMENT")
print("-" .. string.rep("-", 48))

local network_span = observability.start_span(main_trace, "network-connectivity-check", "", {
    component = "network",
    targets = "critical_services"
})

-- Create queue for network test results
queue.create("network_tests", 50)

-- Test critical network services
local network_targets = {
    {host = "google.com", port = 80, service = "HTTP"},
    {host = "8.8.8.8", port = 53, service = "DNS"},
    {host = "github.com", port = 443, service = "HTTPS"}
}

local network_results = {}

for _, target in ipairs(network_targets) do
    local timer = observability.timer_start("network_connectivity_test")
    
    -- Test connectivity
    local connected, msg = network.port_check(target.host, target.port, 3)
    local test_duration = observability.timer_end(timer)
    
    local result = {
        target = target.host .. ":" .. target.port,
        service = target.service,
        status = connected and "UP" or "DOWN",
        response_time = test_duration,
        message = msg or "N/A"
    }
    
    table.insert(network_results, result)
    
    -- Queue the result for processing
    queue.publish("network_tests", data.json_encode(result))
    
    -- Record connectivity metric
    observability.counter("network_connectivity_tests", 1, {
        target = target.host,
        service = target.service,
        status = connected and "success" or "failure"
    })
    
    print("🔌 " .. target.service .. " (" .. target.host .. ":" .. target.port .. "): " .. 
          (connected and "✅ UP" or "❌ DOWN") .. " (" .. test_duration .. "ms)")
end

-- Get local network information
local local_ip = network.local_ip()
if local_ip then
    print("🏠 Local IP Address:", local_ip)
    observability.add_span_tag(network_span, "local_ip", local_ip)
end

observability.add_span_event(network_span, "connectivity-tests-completed", {
    tests_count = tostring(#network_targets),
    successful = tostring(#network_results)
})

observability.end_span(network_span, "completed")

-- 3. SECURITY ASSESSMENT
print("\n🔐 PHASE 3: SECURITY ASSESSMENT")
print("-" .. string.rep("-", 35))

local security_span = observability.start_span(main_trace, "security-assessment", "", {
    component = "security",
    scan_type = "baseline"
})

-- Create queue for security findings
queue.create("security_alerts", 100)

-- Security baseline check
local baseline = security.security_baseline()
local security_issues = 0

if baseline then
    print("🛡️ Security Baseline Assessment:")
    print("   Compliance Rate:", string.format("%.1f%%", baseline.compliance_percentage or 0))
    print("   Checks Passed:", (baseline.passed_checks or 0) .. "/" .. (baseline.total_checks or 0))
    
    -- Record security metrics
    observability.gauge("security_compliance_percentage", baseline.compliance_percentage or 0)
    observability.gauge("security_checks_passed", baseline.passed_checks or 0)
    
    if baseline.compliance_percentage and baseline.compliance_percentage < 80 then
        security_issues = security_issues + 1
        queue.publish("security_alerts", "Low security compliance: " .. 
                     string.format("%.1f%%", baseline.compliance_percentage))
    end
end

-- Check firewall status
local firewall = security.firewall_status()
if firewall then
    local firewall_active = firewall.iptables_active or firewall.ufw_active or firewall.firewalld_active
    print("🔥 Firewall Status:", firewall_active and "✅ ACTIVE" or "❌ INACTIVE")
    
    if not firewall_active then
        security_issues = security_issues + 1
        queue.publish("security_alerts", "Firewall is not active")
    end
    
    observability.counter("security_firewall_check", 1, {
        status = firewall_active and "active" or "inactive"
    })
end

-- Password strength assessment
local test_passwords = {"admin123", "MySecureP@ssw0rd!2024"}
for i, pwd in ipairs(test_passwords) do
    local strength = security.password_strength(pwd)
    if strength then
        print("🔑 Password " .. i .. " Strength:", strength.strength .. " (" .. strength.score .. "/100)")
        
        if strength.score < 70 then
            security_issues = security_issues + 1
            queue.publish("security_alerts", "Weak password detected: " .. strength.strength)
        end
        
        observability.histogram("password_strength_score", strength.score, {
            strength_level = strength.strength
        })
    end
end

-- SSL certificate check
local ssl_check = security.check_ssl_cert("https://google.com")
if ssl_check and ssl_check.valid then
    local days_until_expiry = ssl_check.expires_in_days or 0
    print("🔒 SSL Certificate Check: ✅ VALID (expires in " .. math.floor(days_until_expiry) .. " days)")
    
    if days_until_expiry < 30 then
        security_issues = security_issues + 1
        queue.publish("security_alerts", "SSL certificate expires soon: " .. math.floor(days_until_expiry) .. " days")
    end
    
    observability.gauge("ssl_certificate_days_remaining", days_until_expiry, {
        domain = "google.com"
    })
end

print("⚠️ Security Issues Found:", security_issues)

observability.add_span_event(security_span, "security-scan-completed", {
    issues_found = tostring(security_issues),
    compliance_rate = tostring(baseline.compliance_percentage or 0)
})

observability.end_span(security_span, "completed")

-- 4. QUEUE PROCESSING & ALERTING
print("\n📬 PHASE 4: QUEUE PROCESSING & ALERTING")
print("-" .. string.rep("-", 45))

local queue_span = observability.start_span(main_trace, "alert-processing", "", {
    component = "queue",
    processor = "alert_handler"
})

-- Process network test results
local network_queue_size = queue.size("network_tests")
print("📊 Network Test Results Queue:", network_queue_size, "items")

if network_queue_size > 0 then
    local network_messages = queue.consume_batch("network_tests", network_queue_size, 2)
    if network_messages then
        print("📈 Processing network test results:")
        for i = 1, #network_messages do
            local result = data.json_decode(network_messages[i].payload)
            if result and result.status == "DOWN" then
                print("   🚨 Service DOWN: " .. result.service .. " (" .. result.target .. ")")
            end
        end
    end
end

-- Process security alerts
local security_queue_size = queue.size("security_alerts")
print("🚨 Security Alerts Queue:", security_queue_size, "items")

if security_queue_size > 0 then
    local security_messages = queue.consume_batch("security_alerts", security_queue_size, 2)
    if security_messages then
        print("🔐 Processing security alerts:")
        for i = 1, #security_messages do
            local alert = security_messages[i]
            print("   ⚠️ ALERT: " .. alert.payload)
            
            -- Record alert metric
            observability.counter("security_alerts_processed", 1, {
                severity = "medium"
            })
        end
    end
end

-- Queue statistics
local all_queues = queue.list()
if all_queues then
    print("📋 Queue Status Summary:")
    for i = 1, #all_queues do
        local q = all_queues[i]
        print("   " .. q.name .. ": " .. q.size .. "/" .. q.capacity .. " messages")
        
        observability.gauge("queue_size", q.size, {
            queue_name = q.name
        })
    end
end

observability.add_span_event(queue_span, "alert-processing-completed", {
    network_alerts = tostring(network_queue_size),
    security_alerts = tostring(security_queue_size)
})

observability.end_span(queue_span, "completed")

-- 5. OBSERVABILITY & REPORTING
print("\n📊 PHASE 5: OBSERVABILITY & REPORTING")
print("-" .. string.rep("-", 42))

local reporting_span = observability.start_span(main_trace, "generate-report", "", {
    component = "observability",
    report_type = "health_summary"
})

-- Generate comprehensive health report
local health_summary = {
    timestamp = os.date("%Y-%m-%d %H:%M:%S"),
    system_health = health_check,
    network_tests = network_results,
    security_issues = security_issues,
    total_alerts = security_queue_size
}

-- Export observability data
local obs_health = observability.health_check()
if obs_health then
    print("📊 Observability Health:")
    print("   Status:", obs_health.status)
    print("   Active Traces:", obs_health.active_traces)
    print("   Total Metrics:", obs_health.total_metrics)
end

-- System metrics
local sys_metrics = observability.system_metrics()
if sys_metrics and sys_metrics.memory then
    print("💾 Runtime Metrics:")
    print("   Memory Allocated:", string.format("%.2f MB", sys_metrics.memory.alloc / 1024 / 1024))
    print("   Goroutines:", sys_metrics.runtime.goroutines)
    print("   GC Runs:", sys_metrics.memory.gc_runs)
end

-- Export data (simulated)
local json_export = observability.export_json()
print("📤 Data Export:")
print("   Traces exported:", #json_export.traces)
print("   Metrics exported:", #json_export.metrics)

observability.add_span_event(reporting_span, "report-generated", {
    system_score = tostring(health_check.score or 0),
    security_issues = tostring(security_issues),
    network_services = tostring(#network_targets)
})

observability.end_span(reporting_span, "completed")

-- 6. FINAL ASSESSMENT & CLEANUP
print("\n🎯 FINAL ASSESSMENT")
print("-" .. string.rep("-", 25))

-- End main trace
local trace_success, total_duration = observability.end_trace(main_trace, "completed")

-- Calculate overall health score
local overall_score = (health_check.score or 0)
if security_issues > 0 then
    overall_score = overall_score - (security_issues * 10)  -- Deduct for security issues
end

-- Determine status
local overall_status = "HEALTHY"
if overall_score < 60 then
    overall_status = "CRITICAL"
elseif overall_score < 80 then
    overall_status = "WARNING"
end

print("🏆 INFRASTRUCTURE HEALTH REPORT")
print("   Overall Score:", string.format("%.0f/100", overall_score))
print("   Status:", overall_status)
print("   Assessment Duration:", total_duration .. "ms")
print("   Security Issues:", security_issues)
print("   Network Services:", #network_targets)

-- Record final metrics
observability.gauge("infrastructure_health_score", overall_score, {
    status = overall_status
})

observability.counter("health_assessments_completed", 1, {
    status = overall_status,
    duration_ms = tostring(total_duration)
})

-- Cleanup queues
queue.delete("network_tests")
queue.delete("security_alerts")

-- 7. INTEGRATION SUMMARY
print("\n🚀 INTEGRATION SUMMARY")
print("-" .. string.rep("-", 30))

print("✅ Successfully demonstrated:")
print("   🖥️  System monitoring and health assessment")
print("   🌐 Network connectivity testing and diagnostics")
print("   🔐 Security scanning and compliance checking")
print("   📬 Queue-based alert processing and management")
print("   📊 Distributed tracing and metrics collection")

print("\n💡 Enterprise capabilities achieved:")
print("   • Real-time infrastructure monitoring")
print("   • Automated security compliance checking")
print("   • Distributed system observability")
print("   • Event-driven alert processing")
print("   • Performance metrics and analytics")
print("   • Comprehensive health reporting")

print("\n🎯 Use cases enabled:")
print("   • DevOps monitoring dashboards")
print("   • Security operations centers (SOC)")
print("   • Site reliability engineering (SRE)")
print("   • Compliance and audit automation")
print("   • Infrastructure as code validation")
print("   • Multi-cloud environment monitoring")

print("\n✅ HIGH-PRIORITY MODULES INTEGRATION COMPLETED!")
print("🚀 Enterprise-grade infrastructure automation and monitoring system ready!")
print("📊 Total execution time: " .. total_duration .. "ms")