-- System Module Examples

print("💻 SYSTEM MODULE SHOWCASE")
print("=" .. string.rep("=", 40))

-- 1. CPU Information
print("\n🧠 CPU Information:")

local cpu_info = system.cpu_info()
if cpu_info and #cpu_info > 0 then
    local cpu = cpu_info[1]
    print("🔧 CPU Model:", cpu.model or "Unknown")
    print("⚡ CPU Speed:", (cpu.speed_mhz or 0) .. " MHz")
    print("🏭 CPU Cores:", cpu.cores or "Unknown")
    print("🏢 Vendor:", cpu.vendor_id or "Unknown")
end

local cpu_usage = system.cpu_usage(1)
if cpu_usage then
    print("📊 CPU Usage:", string.format("%.1f%%", cpu_usage))
end

local cpu_count = system.cpu_count(true)
if cpu_count then
    print("💾 Logical CPUs:", cpu_count)
end

-- 2. Memory Information
print("\n🧠 Memory Information:")

local memory = system.memory_info()
if memory then
    print("💾 Total Memory:", string.format("%.2f GB", memory.total / 1024 / 1024 / 1024))
    print("✅ Available Memory:", string.format("%.2f GB", memory.available / 1024 / 1024 / 1024))
    print("⚡ Memory Usage:", string.format("%.1f%%", memory.percent))
    print("📊 Used Memory:", string.format("%.2f GB", memory.used / 1024 / 1024 / 1024))
    
    if memory.buffers and memory.buffers > 0 then
        print("🗃️ Buffers:", string.format("%.2f MB", memory.buffers / 1024 / 1024))
    end
    if memory.cached and memory.cached > 0 then
        print("💿 Cached:", string.format("%.2f MB", memory.cached / 1024 / 1024))
    end
end

local swap = system.swap_info()
if swap and swap.total > 0 then
    print("🔄 Swap Total:", string.format("%.2f GB", swap.total / 1024 / 1024 / 1024))
    print("🔄 Swap Used:", string.format("%.1f%%", swap.percent))
end

-- 3. Disk Information
print("\n💾 Disk Information:")

local disk_usage = system.disk_usage("/")
if disk_usage then
    print("🗂️ Root Filesystem:")
    print("   Total Space:", string.format("%.2f GB", disk_usage.total / 1024 / 1024 / 1024))
    print("   Used Space:", string.format("%.2f GB", disk_usage.used / 1024 / 1024 / 1024))
    print("   Free Space:", string.format("%.2f GB", disk_usage.free / 1024 / 1024 / 1024))
    print("   Usage:", string.format("%.1f%%", disk_usage.percent))
end

local partitions = system.disk_partitions()
if partitions then
    print("🗂️ Disk Partitions (" .. #partitions .. " found):")
    for i = 1, math.min(#partitions, 3) do  -- Show first 3
        local part = partitions[i]
        print("   - " .. (part.device or "unknown") .. " -> " .. (part.mountpoint or "unknown"))
        print("     Type: " .. (part.fstype or "unknown"))
    end
end

-- 4. Network Stats
print("\n🌐 Network Statistics:")

local net_stats = system.network_stats()
if net_stats and #net_stats > 0 then
    print("📡 Network Interfaces Statistics:")
    for i = 1, math.min(#net_stats, 2) do  -- Show first 2
        local stat = net_stats[i]
        if stat.name and stat.name ~= "lo" then  -- Skip loopback
            print("   Interface: " .. stat.name)
            print("     Bytes Sent: " .. string.format("%.2f MB", stat.bytes_sent / 1024 / 1024))
            print("     Bytes Received: " .. string.format("%.2f MB", stat.bytes_recv / 1024 / 1024))
            print("     Packets Sent: " .. stat.packets_sent)
            print("     Packets Received: " .. stat.packets_recv)
        end
    end
end

-- 5. System Information
print("\n🖥️ System Information:")

local host_info = system.host_info()
if host_info then
    print("🏠 Hostname:", host_info.hostname or "Unknown")
    print("🖥️ Operating System:", host_info.os or "Unknown")
    print("📋 Platform:", host_info.platform or "Unknown")
    print("🏗️ Architecture:", host_info.kernel_arch or "Unknown")
    print("🔢 Kernel Version:", host_info.kernel_version or "Unknown")
    print("🔄 Total Processes:", host_info.procs or "Unknown")
    
    if host_info.virtualization_system and host_info.virtualization_system ~= "" then
        print("☁️ Virtualization:", host_info.virtualization_system)
    end
end

local uptime = system.uptime()
if uptime then
    print("⏰ System Uptime:")
    print("   Days:", uptime.days or 0)
    print("   Hours:", uptime.hours or 0)
    print("   Minutes:", uptime.minutes or 0)
    print("   Human readable:", uptime.human or "Unknown")
end

-- 6. Load Average
print("\n📊 System Load:")

local load = system.load_average()
if load then
    print("⚖️ Load Average:")
    print("   1 minute:", string.format("%.2f", load.load1))
    print("   5 minutes:", string.format("%.2f", load.load5))
    print("   15 minutes:", string.format("%.2f", load.load15))
end

-- 7. Process Information
print("\n🔄 Process Information:")

-- Get information about current process (approximate PID)
local current_pid = 1  -- You would get actual PID in real usage
local proc_info = system.process_info(current_pid)
if proc_info then
    print("🔧 Process " .. current_pid .. ":")
    print("   Name:", proc_info.name or "Unknown")
    print("   Status:", proc_info.status or "Unknown")
    if proc_info.memory then
        print("   Memory RSS:", string.format("%.2f MB", (proc_info.memory.rss or 0) / 1024 / 1024))
    end
end

-- 8. Performance Snapshot
print("\n🎯 Performance Snapshot:")

local snapshot = system.performance_snapshot()
if snapshot then
    print("📸 Current System Performance:")
    if snapshot.cpu_percent then
        print("   CPU Usage:", string.format("%.1f%%", snapshot.cpu_percent))
    end
    if snapshot.memory then
        print("   Memory Usage:", string.format("%.1f%% (%.1f/%.1f GB)", 
            snapshot.memory.percent, snapshot.memory.used_gb, snapshot.memory.total_gb))
    end
    if snapshot.disk then
        print("   Disk Usage:", string.format("%.1f%% (%.1f/%.1f GB)", 
            snapshot.disk.percent, snapshot.disk.used_gb, snapshot.disk.total_gb))
    end
    if snapshot.load then
        print("   Load Average:", string.format("%.2f, %.2f, %.2f", 
            snapshot.load.load1, snapshot.load.load5, snapshot.load.load15))
    end
end

-- 9. System Health Check
print("\n🏥 System Health Check:")

local health = system.system_health()
if health then
    print("🎯 Health Score:", string.format("%.0f/100", health.score))
    print("📊 Status:", health.status or "unknown")
    
    if health.issues and #health.issues > 0 then
        print("⚠️ Issues Found:")
        for i = 1, math.min(#health.issues, 3) do  -- Show first 3 issues
            print("   - " .. health.issues[i])
        end
    else
        print("✅ No critical issues detected")
    end
end

-- 10. Environment Variables (sample)
print("\n🌍 Environment Variables (sample):")

local env = system.environment()
if env then
    local sample_vars = {"HOME", "PATH", "USER", "SHELL"}
    for _, var in ipairs(sample_vars) do
        if env[var] then
            local value = env[var]
            if #value > 50 then
                value = string.sub(value, 1, 47) .. "..."
            end
            print("   " .. var .. ":", value)
        end
    end
end

print("\n✅ System module demonstration completed!")
print("📊 Comprehensive system monitoring and information gathering ready!")