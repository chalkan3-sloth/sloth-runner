# 📊 指标和监控模块

**指标和监控**模块提供全面的系统监控、自定义指标收集和健康检查功能。它实现了对系统资源和应用程序性能的实时观察能力。

## 🚀 核心特性

- **系统指标**: 自动收集CPU、内存、磁盘和网络指标
- **运行时指标**: Go运行时信息（协程、堆、GC）
- **自定义指标**: 计量器、计数器、直方图和计时器
- **健康检查**: 自动系统健康监控
- **HTTP端点**: 兼容Prometheus的指标导出
- **告警系统**: 基于阈值的告警
- **JSON API**: 完整的指标数据用于集成

## 📊 系统指标

### CPU、内存和磁盘监控

```lua
-- 获取当前CPU使用率
local cpu_usage = metrics.system_cpu()
log.info("CPU使用率: " .. string.format("%.1f%%", cpu_usage))

-- 获取内存信息
local memory_info = metrics.system_memory()
log.info("内存: " .. string.format("%.1f%% (%.0f/%.0f MB)", 
    memory_info.percent, memory_info.used_mb, memory_info.total_mb))

-- 获取磁盘使用情况
local disk_info = metrics.system_disk("/")
log.info("磁盘: " .. string.format("%.1f%% (%.1f/%.1f GB)", 
    disk_info.percent, disk_info.used_gb, disk_info.total_gb))

-- 检查特定磁盘路径
local var_disk = metrics.system_disk("/var")
log.info("/var 磁盘使用率: " .. string.format("%.1f%%", var_disk.percent))
```

### 运行时信息

```lua
-- 获取Go运行时指标
local runtime = metrics.runtime_info()
log.info("运行时信息:")
log.info("  协程数: " .. runtime.goroutines)
log.info("  CPU核心: " .. runtime.num_cpu)
log.info("  堆已分配: " .. string.format("%.1f MB", runtime.heap_alloc_mb))
log.info("  堆系统: " .. string.format("%.1f MB", runtime.heap_sys_mb))
log.info("  GC次数: " .. runtime.num_gc)
log.info("  Go版本: " .. runtime.go_version)
```

## 📈 自定义指标

### 计量器指标（当前值）

```lua
-- 设置简单的计量器值
metrics.gauge("cpu_temperature", 65.4)
metrics.gauge("active_connections", 142)
metrics.gauge("queue_size", 23)

-- 带标签设置计量器
metrics.gauge("memory_usage", memory_percent, {
    server = "web-01",
    environment = "production",
    region = "us-east-1"
})

-- 更新部署状态
metrics.gauge("deployment_progress", 75.5, {
    app = "frontend",
    version = "v2.1.0"
})
```

### 计数器指标（增量值）

```lua
-- 增量计数器
local total_requests = metrics.counter("http_requests_total", 1)
local error_count = metrics.counter("http_errors_total", 1, {
    status_code = "500",
    endpoint = "/api/users"
})

-- 批量增量
local processed = metrics.counter("messages_processed", 50, {
    queue = "user_notifications",
    priority = "high"
})

log.info("处理的总请求数: " .. total_requests)
```

### 直方图指标（值分布）

```lua
-- 记录响应时间
metrics.histogram("response_time_ms", 245.6, {
    endpoint = "/api/users",
    method = "GET"
})

-- 记录负载大小
metrics.histogram("payload_size_bytes", 1024, {
    content_type = "application/json"
})

-- 记录批处理大小
metrics.histogram("batch_size", 150, {
    operation = "bulk_insert",
    table = "user_events"
})
```

### 计时器指标（函数执行时间）

```lua
-- 自动计时函数执行
local duration = metrics.timer("database_query", function()
    -- 模拟数据库查询
    local result = exec.run("sleep 0.1")
    return result
end, {
    query_type = "select",
    table = "users"
})

log.info("数据库查询耗时: " .. string.format("%.2f ms", duration))

-- 计时复杂操作
local processing_time = metrics.timer("data_processing", function()
    -- 处理大数据集
    local data = {}
    for i = 1, 100000 do
        data[i] = math.sqrt(i) * 2.5
    end
    return #data
end, {
    operation = "mathematical_computation",
    size = "large"
})

log.info("数据处理完成用时: " .. string.format("%.2f ms", processing_time))
```

## 🏥 健康监控

### 自动健康状态

```lua
-- 获取全面的健康状态
local health = metrics.health_status()
log.info("整体健康状态: " .. health.overall)

-- 检查各个组件
local components = {"cpu", "memory", "disk"}
for _, component in ipairs(components) do
    local comp_info = health[component]
    if comp_info then
        local status_icon = "✅"
        if comp_info.status == "warning" then
            status_icon = "⚠️"
        elseif comp_info.status == "critical" then
            status_icon = "❌"
        end
        
        log.info(string.format("  %s %s: %.1f%% (%s)", 
            status_icon, component:upper(), comp_info.usage, comp_info.status))
    end
end
```

### 自定义健康检查

```lua
-- 创建应用健康检查函数
function check_application_health()
    local health_score = 100
    local issues = {}
    
    -- 检查数据库连通性
    local db_result = exec.run("pg_isready -h localhost -p 5432")
    if db_result ~= "" then
        health_score = health_score - 20
        table.insert(issues, "数据库连接失败")
    end
    
    -- 检查磁盘空间
    local disk = metrics.system_disk("/")
    if disk.percent > 90 then
        health_score = health_score - 30
        table.insert(issues, "磁盘空间严重不足: " .. string.format("%.1f%%", disk.percent))
    end
    
    -- 检查内存使用
    local memory = metrics.system_memory()
    if memory.percent > 85 then
        health_score = health_score - 25
        table.insert(issues, "内存使用率高: " .. string.format("%.1f%%", memory.percent))
    end
    
    -- 记录健康得分
    metrics.gauge("application_health_score", health_score)
    
    if health_score < 70 then
        metrics.alert("application_health", {
            level = "warning",
            message = "应用健康状况下降: " .. table.concat(issues, ", "),
            score = health_score
        })
    end
    
    return health_score >= 70
end

-- 在任务中使用
local health_check = task("health_check")
    :description("执行健康检查")
    :command(function(this, params)
        local healthy = check_application_health()
        return healthy, healthy and "系统健康" or "检测到系统健康问题"
    end)
    :build()

workflow
    .define("health_monitoring")
    :description("健康监控")
    :version("1.0.0")
    :tasks({health_check})
```

## 🚨 告警系统

### 创建告警

```lua
-- 简单阈值告警
local cpu = metrics.system_cpu()
if cpu > 80 then
    metrics.alert("high_cpu_usage", {
        level = "warning",
        message = "CPU使用率过高: " .. string.format("%.1f%%", cpu),
        threshold = 80,
        value = cpu,
        severity = "medium"
    })
end

-- 多条件复杂告警
local memory = metrics.system_memory()
local disk = metrics.system_disk()

if memory.percent > 90 and disk.percent > 85 then
    metrics.alert("resource_exhaustion", {
        level = "critical",
        message = string.format("资源使用严重 - 内存: %.1f%%, 磁盘: %.1f%%", 
            memory.percent, disk.percent),
        memory_usage = memory.percent,
        disk_usage = disk.percent,
        recommended_action = "立即扩展资源"
    })
end

-- 应用特定告警
local queue_size = state.get("task_queue_size", 0)
if queue_size > 1000 then
    metrics.alert("queue_backlog", {
        level = "warning", 
        message = "检测到任务队列积压: " .. queue_size .. " 个项目",
        queue_size = queue_size,
        estimated_processing_time = queue_size * 2 .. " 秒"
    })
end
```

## 🔍 指标管理

### 检索自定义指标

```lua
-- 获取特定自定义指标
local cpu_metric = metrics.get_custom("cpu_temperature")
if cpu_metric then
    log.info("CPU温度指标: " .. data.to_json(cpu_metric))
end

-- 列出所有自定义指标
local all_metrics = metrics.list_custom()
log.info("自定义指标总数: " .. #all_metrics)
for i, metric_name in ipairs(all_metrics) do
    log.info("  " .. i .. ". " .. metric_name)
end
```

### 性能监控示例

```lua
local monitor_api_performance = task("monitor_api_performance")
    :description("监控API性能")
    :command(function(this, params)
        -- 开始监控会话
        log.info("开始API性能监控...")

        -- 模拟API调用并测量性能
        for i = 1, 10 do
            local api_time = metrics.timer("api_call_" .. i, function()
                -- 模拟API调用
                exec.run("curl -s -o /dev/null -w '%{time_total}' https://api.example.com/health")
            end, {
                endpoint = "health",
                call_number = tostring(i)
            })

            -- 记录响应时间
            metrics.histogram("api_response_time", api_time, {
                endpoint = "health"
            })

            -- 检查响应时间是否可接受
            if api_time > 1000 then -- 1秒
                metrics.counter("slow_api_calls", 1, {
                    endpoint = "health"
                })

                metrics.alert("slow_api_response", {
                    level = "warning",
                    message = string.format("API响应慢: %.2f ms", api_time),
                    response_time = api_time,
                    threshold = 1000
                })
            end

            -- 调用间短暂延迟
            exec.run("sleep 0.1")
        end

        -- 获取汇总统计
        local system_health = metrics.health_status()
        log.info("API测试后系统健康: " .. system_health.overall)

        return true, "API性能监控完成"
    end)
    :build()

workflow
    .define("performance_monitoring")
    :description("性能监控")
    :version("1.0.0")
    :tasks({monitor_api_performance})
```

## 🌐 HTTP端点

指标模块自动为外部监控系统公开HTTP端点：

### Prometheus格式 (`/metrics`)
```bash
# 访问兼容Prometheus的指标
curl http://agent:8080/metrics

# 示例输出:
# sloth_agent_cpu_usage_percent 15.4
# sloth_agent_memory_usage_mb 2048.5
# sloth_agent_disk_usage_percent 67.2
# sloth_agent_tasks_total 142
```

### JSON格式 (`/metrics/json`)
```bash
# 获取JSON格式的完整指标
curl http://agent:8080/metrics/json

# 示例响应:
{
  "agent_name": "myagent1",
  "timestamp": "2024-01-15T10:30:00Z",
  "system": {
    "cpu_usage_percent": 15.4,
    "memory_usage_mb": 2048.5,
    "disk_usage_percent": 67.2
  },
  "runtime": {
    "num_goroutines": 25,
    "heap_alloc_mb": 45.2
  },
  "custom": {
    "api_response_time": {...},
    "deployment_progress": 85.5
  }
}
```

### 健康检查 (`/health`)
```bash
# 检查代理健康状态
curl http://agent:8080/health

# 示例响应:
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "checks": {
    "cpu": {"usage": 15.4, "status": "healthy"},
    "memory": {"usage": 45.8, "status": "healthy"},
    "disk": {"usage": 67.2, "status": "healthy"}
  }
}
```

## 📋 API参考

### 系统指标
| 函数 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| `metrics.system_cpu()` | - | usage: number | 获取当前CPU使用百分比 |
| `metrics.system_memory()` | - | info: table | 获取内存使用信息 |
| `metrics.system_disk(path?)` | path?: string | info: table | 获取路径的磁盘使用情况 (默认: "/") |
| `metrics.runtime_info()` | - | info: table | 获取Go运行时信息 |

### 自定义指标
| 函数 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| `metrics.gauge(name, value, tags?)` | name: string, value: number, tags?: table | success: boolean | 设置计量器指标 |
| `metrics.counter(name, increment?, tags?)` | name: string, increment?: number, tags?: table | new_value: number | 增量计数器 |
| `metrics.histogram(name, value, tags?)` | name: string, value: number, tags?: table | success: boolean | 记录直方图值 |
| `metrics.timer(name, function, tags?)` | name: string, func: function, tags?: table | duration: number | 计时函数执行 |

### 健康和监控
| 函数 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| `metrics.health_status()` | - | status: table | 获取全面健康状态 |
| `metrics.alert(name, data)` | name: string, data: table | success: boolean | 创建告警 |

### 实用工具
| 函数 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| `metrics.get_custom(name)` | name: string | metric: table \| nil | 按名称获取自定义指标 |
| `metrics.list_custom()` | - | names: table | 列出所有自定义指标名称 |

## 🎯 最佳实践

1. **使用合适的指标类型** - 计量器用于当前值，计数器用于总计，直方图用于分布
2. **添加有意义的标签** 来分类和过滤指标
3. **设置合理的告警阈值** 以避免告警疲劳
4. **监控广泛指标收集的性能影响**
5. **对性能关键操作使用计时器** 来识别瓶颈
6. **为所有关键系统组件实施健康检查**
7. **将指标导出到外部系统** 如Prometheus进行长期存储

**指标和监控**模块为您的分布式sloth-runner环境提供全面的可观测性! 📊🚀