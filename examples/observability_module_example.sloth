-- Observability Module Examples

print("📊 OBSERVABILITY MODULE SHOWCASE")
print("=" .. string.rep("=", 50))

-- 1. Distributed Tracing
print("\n🕸️ Distributed Tracing:")

-- Start a trace for deployment pipeline
local trace_id = observability.start_trace("deployment-pipeline", {
    environment = "production",
    service = "api",
    version = "v3.1.0"
})

print("🎯 Started trace:", trace_id)

-- Create spans for different stages
local build_span = observability.start_span(trace_id, "build-application", "", {
    stage = "build",
    tool = "docker"
})

-- Simulate build process
print("🔨 Building application...")
time.sleep(0.1)  -- Simulate build time

-- Add events to span
observability.add_span_event(build_span, "docker-pull", {
    image = "node:18-alpine",
    status = "completed"
})

observability.add_span_event(build_span, "npm-install", {
    packages = "157",
    duration = "45s"
})

-- Add tags to span
observability.add_span_tag(build_span, "build.success", "true")
observability.add_span_tag(build_span, "build.duration", "2m15s")

-- End build span
local build_success, build_duration = observability.end_span(build_span, "completed")
print("✅ Build completed in " .. build_duration .. "ms")

-- Start test span
local test_span = observability.start_span(trace_id, "run-tests", build_span, {
    stage = "test",
    type = "unit+integration"
})

print("🧪 Running tests...")
time.sleep(0.05)  -- Simulate test time

observability.add_span_event(test_span, "unit-tests", {
    passed = "45",
    failed = "0",
    coverage = "92%"
})

observability.add_span_event(test_span, "integration-tests", {
    passed = "12",
    failed = "0"
})

local test_success, test_duration = observability.end_span(test_span, "completed")
print("✅ Tests completed in " .. test_duration .. "ms")

-- Start deployment span
local deploy_span = observability.start_span(trace_id, "deploy-to-production", test_span, {
    stage = "deploy",
    target = "kubernetes"
})

print("🚀 Deploying to production...")
time.sleep(0.08)  -- Simulate deploy time

observability.add_span_event(deploy_span, "image-push", {
    registry = "gcr.io",
    tag = "v3.1.0"
})

observability.add_span_event(deploy_span, "k8s-rollout", {
    replicas = "3",
    strategy = "rolling"
})

local deploy_success, deploy_duration = observability.end_span(deploy_span, "completed")
print("✅ Deployment completed in " .. deploy_duration .. "ms")

-- End the trace
local trace_success, total_duration = observability.end_trace(trace_id, "completed")
print("🎯 Trace completed in " .. total_duration .. "ms")

-- 2. Metrics Collection
print("\n📈 Metrics Collection:")

-- Counter metrics
observability.counter("deployments_total", 1, {
    environment = "production",
    service = "api",
    status = "success"
})

observability.counter("http_requests_total", 150, {
    method = "GET",
    endpoint = "/api/users",
    status_code = "200"
})

observability.counter("errors_total", 2, {
    type = "timeout",
    service = "database"
})

print("📊 Counter metrics recorded:")
print("   • Deployment successful")
print("   • HTTP requests: 150")
print("   • Errors: 2 (timeouts)")

-- Gauge metrics
observability.gauge("cpu_usage_percent", 75.5, {
    host = "web-server-01",
    region = "us-east-1"
})

observability.gauge("memory_usage_bytes", 2147483648, {  -- 2GB
    host = "web-server-01",
    type = "rss"
})

observability.gauge("active_connections", 42, {
    service = "api",
    port = "8080"
})

print("📊 Gauge metrics recorded:")
print("   • CPU usage: 75.5%")
print("   • Memory: 2GB")
print("   • Active connections: 42")

-- Histogram metrics
observability.histogram("request_duration_ms", 125.5, {
    endpoint = "/api/users",
    method = "GET"
})

observability.histogram("database_query_duration_ms", 45.2, {
    query_type = "SELECT",
    table = "users"
})

print("📊 Histogram metrics recorded:")
print("   • Request duration: 125.5ms")
print("   • DB query: 45.2ms")

-- 3. Timer Measurements
print("\n⏱️ Timer Measurements:")

-- Start a timer for API processing
local api_timer = observability.timer_start("api_processing_time")
print("⏱️ Started API processing timer")

-- Simulate API processing
time.sleep(0.05)

-- End timer
local timer_success, duration = observability.timer_end(api_timer)
if timer_success then
    print("⏱️ API processing completed in " .. duration .. "ms")
end

-- Start a timer for database operation
local db_timer = observability.timer_start("database_operation")
time.sleep(0.02)
local db_timer_success, db_duration = observability.timer_end(db_timer)
if db_timer_success then
    print("⏱️ Database operation completed in " .. db_duration .. "ms")
end

-- 4. Trace Inspection
print("\n🔍 Trace Inspection:")

-- Get trace details
local trace_details = observability.get_trace(trace_id)
if trace_details then
    print("🎯 Trace details:")
    print("   ID:", trace_details.id)
    print("   Name:", trace_details.name)
    print("   Status:", trace_details.status)
    print("   Duration:", trace_details.duration_ms .. "ms")
    print("   Spans:", #trace_details.spans)
    
    if trace_details.tags then
        print("   Tags:")
        for key, value in pairs(trace_details.tags) do
            print("     " .. key .. ": " .. value)
        end
    end
end

-- List all traces
local all_traces = observability.list_traces("completed")
if all_traces then
    print("📋 All completed traces:")
    for i = 1, #all_traces do
        local trace = all_traces[i]
        print("   " .. i .. ". " .. trace.name .. " (" .. (trace.duration_ms or 0) .. "ms)")
    end
end

-- 5. System Metrics
print("\n🖥️ System Metrics:")

local sys_metrics = observability.system_metrics()
if sys_metrics then
    print("📊 Runtime metrics:")
    
    if sys_metrics.memory then
        print("   Memory:")
        print("     Allocated: " .. string.format("%.2f MB", sys_metrics.memory.alloc / 1024 / 1024))
        print("     Total allocated: " .. string.format("%.2f MB", sys_metrics.memory.total_alloc / 1024 / 1024))
        print("     System: " .. string.format("%.2f MB", sys_metrics.memory.sys / 1024 / 1024))
        print("     GC runs: " .. sys_metrics.memory.gc_runs)
    end
    
    if sys_metrics.runtime then
        print("   Runtime:")
        print("     Goroutines: " .. sys_metrics.runtime.goroutines)
        print("     CPUs: " .. sys_metrics.runtime.cpus)
        print("     Go version: " .. sys_metrics.runtime.go_version)
    end
end

-- 6. Health Check
print("\n🏥 Health Check:")

local health = observability.health_check()
if health then
    print("💚 Observability health:")
    print("   Status:", health.status)
    print("   Active traces:", health.active_traces)
    print("   Total traces:", health.total_traces)
    print("   Total metrics:", health.total_metrics)
    print("   Timestamp:", os.date("%H:%M:%S", health.timestamp))
end

-- 7. Data Export
print("\n📤 Data Export:")

-- Export to JSON
local json_export = observability.export_json()
if json_export then
    print("📊 JSON export:")
    print("   Traces exported:", #json_export.traces)
    print("   Metrics exported:", #json_export.metrics)
    print("   Export timestamp:", os.date("%H:%M:%S", json_export.exported_at))
end

-- Export to Jaeger (simulation)
local jaeger_export = observability.export_jaeger("http://localhost:14268/api/traces")
if jaeger_export then
    print("🕸️ Jaeger export (simulated):")
    print("   Success:", jaeger_export.success and "Yes" or "No")
    print("   Endpoint:", jaeger_export.endpoint)
    print("   Exported traces:", jaeger_export.exported_traces)
    print("   Note:", jaeger_export.note)
end

-- Export to Prometheus (simulation)
local prometheus_export = observability.export_prometheus("http://localhost:9090/api/v1/write")
if prometheus_export then
    print("📈 Prometheus export (simulated):")
    print("   Success:", prometheus_export.success and "Yes" or "No")
    print("   Endpoint:", prometheus_export.endpoint)
    print("   Exported metrics:", prometheus_export.exported_metrics)
    print("   Note:", prometheus_export.note)
end

-- 8. Advanced Observability Patterns
print("\n🚀 Advanced Patterns:")

print("💡 Advanced observability capabilities:")
print("   • Distributed trace correlation")
print("   • Custom metric aggregations")
print("   • Real-time dashboards")
print("   • Alerting based on traces/metrics")
print("   • Performance bottleneck detection")
print("   • Service dependency mapping")
print("   • Error rate tracking")
print("   • SLA monitoring")

print("\n📋 Integration Examples:")
print("🔗 Works seamlessly with:")
print("   • Jaeger for distributed tracing")
print("   • Prometheus for metrics storage")
print("   • Grafana for visualization")
print("   • AlertManager for notifications")
print("   • OpenTelemetry standards")
print("   • Custom monitoring solutions")

print("\n📊 Use Cases:")
print("🎯 Perfect for:")
print("   • Microservice monitoring")
print("   • Performance optimization")
print("   • Error tracking and debugging")
print("   • SRE and DevOps workflows")
print("   • Compliance and auditing")
print("   • Capacity planning")

-- 9. Sample Dashboard Data
print("\n📊 Sample Dashboard Data:")

print("🎛️ Key metrics summary:")
print("   • Total requests: 1,247")
print("   • Error rate: 0.16%")
print("   • Avg response time: 125ms")
print("   • P95 response time: 450ms")
print("   • Successful deployments: 1")
print("   • Active traces: 1")
print("   • System health: Green")

print("\n✅ Observability module demonstration completed!")
print("📊 Enterprise-grade monitoring and tracing system ready!")
print("🔍 Full visibility into your distributed systems achieved!")