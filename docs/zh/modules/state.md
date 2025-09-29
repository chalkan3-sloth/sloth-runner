# 💾 状态管理模块

**状态管理**模块提供强大的持久化状态功能，包括原子操作、分布式锁和TTL（生存时间）功能。所有数据都使用SQLite的WAL模式在本地存储，以获得最大的性能和可靠性。

## 🚀 核心特性

- **SQLite 持久化**: 使用WAL模式的可靠存储
- **原子操作**: 线程安全的增量、比较交换、追加操作
- **分布式锁**: 带自动超时的临界区
- **TTL (生存时间)**: 自动键过期
- **数据类型**: 字符串、数字、布尔值、表、列表
- **模式匹配**: 通配符键搜索
- **自动清理**: 过期数据的后台清理
- **统计信息**: 使用情况和性能指标

## 📋 基本用法

### 设置和获取值

```lua
-- 设置值
state.set("app_version", "v1.2.3")
state.set("user_count", 1000)
state.set("config", {
    debug = true,
    max_connections = 100
})

-- 获取值
local version = state.get("app_version")
local count = state.get("user_count")
local config = state.get("config")

-- 带默认值获取
local theme = state.get("ui_theme", "dark")

-- 检查存在性
if state.exists("app_version") then
    log.info("应用版本已配置")
end

-- 删除键
state.delete("old_key")
```

### TTL (生存时间)

```lua
-- 设置带TTL (60秒)
state.set("session_token", "abc123", 60)

-- 为现有键设置TTL
state.set_ttl("user_session", 300) -- 5分钟

-- 检查剩余TTL
local ttl = state.get_ttl("session_token")
log.info("令牌在 " .. ttl .. " 秒后过期")
```

### 原子操作

```lua
-- 原子增量
local counter = state.increment("page_views", 1)
local bulk_counter = state.increment("downloads", 50)

-- 原子减量  
local remaining = state.decrement("inventory", 5)

-- 字符串追加
state.set("log_messages", "启动应用程序")
local new_length = state.append("log_messages", " -> 连接到数据库")

-- 原子比较交换
local old_version = state.get("config_version")
local success = state.compare_swap("config_version", old_version, old_version + 1)
if success then
    log.info("配置安全更新")
end
```

### 列表操作

```lua
-- 添加项目到列表
state.list_push("deployment_queue", {
    app = "frontend",
    version = "v2.1.0",
    environment = "staging"
})

-- 检查列表大小
local queue_size = state.list_length("deployment_queue")
log.info("队列中的项目: " .. queue_size)

-- 处理列表 (pop移除最后一项)
while state.list_length("deployment_queue") > 0 do
    local deployment = state.list_pop("deployment_queue")
    log.info("处理部署: " .. deployment.app)
    -- 处理部署...
end
```

### 分布式锁和临界区

```lua
-- 尝试获取锁 (不等待)
local lock_acquired = state.try_lock("deployment_lock", 30) -- 30秒TTL
if lock_acquired then
    -- 关键工作
    state.unlock("deployment_lock")
end

-- 带等待和超时的锁
local acquired = state.lock("database_migration", 60) -- 等待最多60秒
if acquired then
    -- 执行迁移
    state.unlock("database_migration")
end

-- 带自动锁管理的临界区
state.with_lock("critical_section", function()
    log.info("执行关键操作...")
    
    -- 更新全局计数器
    local counter = state.increment("global_counter", 1)
    
    -- 更新时间戳
    state.set("last_operation", os.time())
    
    log.info("关键操作完成 - 计数器: " .. counter)
    
    -- 函数返回时自动释放锁
    return "operation_success"
end, 15) -- 15秒超时
```

## 🔍 API参考

### 基本操作
| 函数 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| `state.set(key, value, ttl?)` | key: string, value: any, ttl?: number | success: boolean | 设置值，可选TTL |
| `state.get(key, default?)` | key: string, default?: any | value: any | 获取值或返回默认值 |
| `state.delete(key)` | key: string | success: boolean | 删除键 |
| `state.exists(key)` | key: string | exists: boolean | 检查键是否存在 |
| `state.clear(pattern?)` | pattern?: string | success: boolean | 按模式删除键 |

### TTL操作
| 函数 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| `state.set_ttl(key, seconds)` | key: string, seconds: number | success: boolean | 为现有键设置TTL |
| `state.get_ttl(key)` | key: string | ttl: number | 获取剩余TTL (-1 = 无TTL, -2 = 不存在) |

### 原子操作
| 函数 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| `state.increment(key, delta?)` | key: string, delta?: number | new_value: number | 原子增量值 |
| `state.decrement(key, delta?)` | key: string, delta?: number | new_value: number | 原子减量值 |
| `state.append(key, value)` | key: string, value: string | new_length: number | 原子追加字符串 |
| `state.compare_swap(key, old, new)` | key: string, old: any, new: any | success: boolean | 原子比较交换 |

### 列表操作
| 函数 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| `state.list_push(key, item)` | key: string, item: any | length: number | 添加项目到列表末尾 |
| `state.list_pop(key)` | key: string | item: any \| nil | 移除并返回最后一项 |
| `state.list_length(key)` | key: string | length: number | 获取列表长度 |

### 分布式锁
| 函数 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| `state.try_lock(name, ttl)` | name: string, ttl: number | success: boolean | 尝试获取锁而不等待 |
| `state.lock(name, timeout?)` | name: string, timeout?: number | success: boolean | 带超时获取锁 |
| `state.unlock(name)` | name: string | success: boolean | 释放锁 |
| `state.with_lock(name, fn, timeout?)` | name: string, fn: function, timeout?: number | result: any | 使用自动锁执行函数 |

### 实用工具
| 函数 | 参数 | 返回值 | 描述 |
|------|------|--------|------|
| `state.keys(pattern?)` | pattern?: string | keys: table | 按模式列出键 |
| `state.stats()` | - | stats: table | 获取系统统计信息 |

## 💡 实际用例

### 1. 部署版本控制

```lua
Modern DSLs = {
    deployment_pipeline = {
        tasks = {
            prepare_deploy = {
                command = function()
                    -- 检查最后部署的版本
                    local last_version = state.get("last_deployed_version", "v0.0.0")
                    local new_version = "v1.2.3"
                    
                    -- 检查是否已部署
                    if last_version == new_version then
                        log.warn("版本 " .. new_version .. " 已部署")
                        return false, "版本已部署"
                    end
                    
                    -- 注册部署开始
                    state.set("deploy_status", "in_progress")
                    state.set("deploy_start_time", os.time())
                    state.increment("total_deploys", 1)
                    
                    return true, "部署准备完成"
                end
            },
            
            execute_deploy = {
                depends_on = "prepare_deploy",
                command = function()
                    -- 部署的临界区
                    return state.with_lock("deployment_lock", function()
                        log.info("使用锁执行部署...")
                        
                        -- 模拟部署
                        exec.run("sleep 5")
                        
                        -- 更新状态
                        state.set("last_deployed_version", "v1.2.3")
                        state.set("deploy_status", "completed")
                        state.set("deploy_end_time", os.time())
                        
                        -- 记录历史
                        state.list_push("deploy_history", {
                            version = "v1.2.3",
                            timestamp = os.time(),
                            duration = state.get("deploy_end_time") - state.get("deploy_start_time")
                        })
                        
                        return true, "部署成功完成"
                    end, 300) -- 5分钟超时
                end
            }
        }
    }
}
```

### 2. 带TTL的智能缓存

```lua
-- 缓存助手函数
function get_cached_data(cache_key, fetch_function, ttl)
    local cached = state.get(cache_key)
    if cached then
        log.info("缓存命中: " .. cache_key)
        return cached
    end
    
    log.info("缓存未命中: " .. cache_key .. " - 正在获取...")
    local data = fetch_function()
    state.set(cache_key, data, ttl or 300) -- 默认5分钟
    return data
end

-- 在任务中使用
Modern DSLs = {
    data_processing = {
        tasks = {
            fetch_user_data = {
                command = function()
                    local user_data = get_cached_data("user:123:profile", function()
                        -- 模拟昂贵的获取操作
                        return {
                            name = "张三",
                            email = "zhangsan@example.com",
                            preferences = {"dark_mode", "notifications"}
                        }
                    end, 600) -- 缓存10分钟
                    
                    log.info("用户数据: " .. data.to_json(user_data))
                    return true, "用户数据已获取"
                end
            }
        }
    }
}
```

### 3. 速率限制

```lua
function check_rate_limit(identifier, max_requests, window_seconds)
    local key = "rate_limit:" .. identifier
    local current_count = state.get(key, 0)
    
    if current_count >= max_requests then
        return false, "速率限制超出"
    end
    
    -- 增加计数器
    if current_count == 0 then
        -- 窗口中的第一个请求
        state.set(key, 1, window_seconds)
    else
        -- 增加现有计数器
        state.increment(key, 1)
    end
    
    return true, "请求允许"
end

-- 在任务中使用
Modern DSLs = {
    api_tasks = {
        tasks = {
            make_api_call = {
                command = function()
                    local allowed, msg = check_rate_limit("api_calls", 100, 3600) -- 100次调用/小时
                    
                    if not allowed then
                        log.error(msg)
                        return false, msg
                    end
                    
                    -- 进行API调用
                    log.info("进行API调用...")
                    return true, "API调用完成"
                end
            }
        }
    }
}
```

## ⚙️ 配置和存储

### 数据库位置

默认情况下，SQLite数据库创建在:
- **Linux/macOS**: `~/.sloth-runner/state.db`
- **Windows**: `%USERPROFILE%\.sloth-runner\state.db`

### 技术特性

- **引擎**: 带WAL模式的SQLite 3
- **并发访问**: 支持多个同时连接
- **自动清理**: 每5分钟自动清理过期数据
- **锁超时**: 过期锁自动清理
- **序列化**: 复杂对象使用JSON，简单类型使用原生格式

### 限制

- **本地范围**: 状态仅在本地机器上持久化
- **并发性**: 锁仅在本地进程内有效
- **大小**: 适合小到中型数据集 (< 1GB)

## 🔄 最佳实践

1. **对临时数据使用TTL** 以防止存储膨胀
2. **对临界区使用锁** 以避免竞态条件
3. **使用模式进行批量操作** 管理相关键
4. **使用`state.stats()`监控存储大小**
5. **使用原子操作** 而不是读-修改-写模式
6. **使用`state.clear(pattern)`定期清理过期键**

**状态管理**模块将sloth-runner转变为有状态的、可靠的复杂任务编排平台! 🚀