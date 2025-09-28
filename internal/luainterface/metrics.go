package luainterface

import (
	"fmt"
	"runtime"
	"time"

	lua "github.com/yuin/gopher-lua"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/disk"
)

// MetricsModule provides access to system and custom metrics from Lua
type MetricsModule struct {
	customMetrics map[string]interface{}
}

// NewMetricsModule creates a new metrics module
func NewMetricsModule() *MetricsModule {
	return &MetricsModule{
		customMetrics: make(map[string]interface{}),
	}
}

// Loader is the module loader function
func (m *MetricsModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"system_cpu":        m.luaSystemCPU,
		"system_memory":     m.luaSystemMemory,
		"system_disk":       m.luaSystemDisk,
		"runtime_info":      m.luaRuntimeInfo,
		"gauge":             m.luaGauge,
		"counter":           m.luaCounter,
		"histogram":         m.luaHistogram,
		"timer":             m.luaTimer,
		"get_custom":        m.luaGetCustom,
		"list_custom":       m.luaListCustom,
		"alert":             m.luaAlert,
		"health_status":     m.luaHealthStatus,
	})
	L.Push(mod)
	return 1
}

// MetricsLoader is the global loader function
func MetricsLoader(L *lua.LState) int {
	return NewMetricsModule().Loader(L)
}

// luaSystemCPU returns current CPU usage percentage
func (m *MetricsModule) luaSystemCPU(L *lua.LState) int {
	cpuPercents, err := cpu.Percent(time.Second, false)
	if err != nil {
		L.Push(lua.LNumber(0))
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	if len(cpuPercents) > 0 {
		L.Push(lua.LNumber(cpuPercents[0]))
	} else {
		L.Push(lua.LNumber(0))
	}
	return 1
}

// luaSystemMemory returns memory usage information
func (m *MetricsModule) luaSystemMemory(L *lua.LState) int {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	memTable := L.NewTable()
	memTable.RawSetString("used_mb", lua.LNumber(float64(memInfo.Used)/1024/1024))
	memTable.RawSetString("total_mb", lua.LNumber(float64(memInfo.Total)/1024/1024))
	memTable.RawSetString("percent", lua.LNumber(memInfo.UsedPercent))
	memTable.RawSetString("available_mb", lua.LNumber(float64(memInfo.Available)/1024/1024))

	L.Push(memTable)
	return 1
}

// luaSystemDisk returns disk usage information
func (m *MetricsModule) luaSystemDisk(L *lua.LState) int {
	path := "/"
	if L.GetTop() >= 1 {
		path = L.CheckString(1)
	}

	diskInfo, err := disk.Usage(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	diskTable := L.NewTable()
	diskTable.RawSetString("used_gb", lua.LNumber(float64(diskInfo.Used)/1024/1024/1024))
	diskTable.RawSetString("total_gb", lua.LNumber(float64(diskInfo.Total)/1024/1024/1024))
	diskTable.RawSetString("percent", lua.LNumber(diskInfo.UsedPercent))
	diskTable.RawSetString("free_gb", lua.LNumber(float64(diskInfo.Free)/1024/1024/1024))

	L.Push(diskTable)
	return 1
}

// luaRuntimeInfo returns Go runtime information
func (m *MetricsModule) luaRuntimeInfo(L *lua.LState) int {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	runtimeTable := L.NewTable()
	runtimeTable.RawSetString("goroutines", lua.LNumber(runtime.NumGoroutine()))
	runtimeTable.RawSetString("num_cpu", lua.LNumber(runtime.NumCPU()))
	runtimeTable.RawSetString("heap_alloc_mb", lua.LNumber(float64(memStats.HeapAlloc)/1024/1024))
	runtimeTable.RawSetString("heap_sys_mb", lua.LNumber(float64(memStats.HeapSys)/1024/1024))
	runtimeTable.RawSetString("num_gc", lua.LNumber(memStats.NumGC))
	runtimeTable.RawSetString("go_version", lua.LString(runtime.Version()))

	L.Push(runtimeTable)
	return 1
}

// luaGauge sets a gauge metric (current value)
func (m *MetricsModule) luaGauge(L *lua.LState) int {
	name := L.CheckString(1)
	value := L.CheckNumber(2)
	
	tags := make(map[string]string)
	if L.GetTop() >= 3 {
		tagsTable := L.CheckTable(3)
		tagsTable.ForEach(func(key, val lua.LValue) {
			tags[key.String()] = val.String()
		})
	}

	metric := map[string]interface{}{
		"type":      "gauge",
		"value":     float64(value),
		"tags":      tags,
		"timestamp": time.Now(),
	}

	m.customMetrics[name] = metric
	L.Push(lua.LBool(true))
	return 1
}

// luaCounter increments a counter metric
func (m *MetricsModule) luaCounter(L *lua.LState) int {
	name := L.CheckString(1)
	increment := 1.0
	if L.GetTop() >= 2 {
		increment = float64(L.CheckNumber(2))
	}

	tags := make(map[string]string)
	if L.GetTop() >= 3 {
		tagsTable := L.CheckTable(3)
		tagsTable.ForEach(func(key, val lua.LValue) {
			tags[key.String()] = val.String()
		})
	}

	// Get existing counter or create new one
	currentValue := 0.0
	if existing, exists := m.customMetrics[name]; exists {
		if existingMetric, ok := existing.(map[string]interface{}); ok {
			if val, ok := existingMetric["value"].(float64); ok {
				currentValue = val
			}
		}
	}

	newValue := currentValue + increment
	metric := map[string]interface{}{
		"type":      "counter",
		"value":     newValue,
		"tags":      tags,
		"timestamp": time.Now(),
	}

	m.customMetrics[name] = metric
	L.Push(lua.LNumber(newValue))
	return 1
}

// luaHistogram records a histogram value
func (m *MetricsModule) luaHistogram(L *lua.LState) int {
	name := L.CheckString(1)
	value := L.CheckNumber(2)

	tags := make(map[string]string)
	if L.GetTop() >= 3 {
		tagsTable := L.CheckTable(3)
		tagsTable.ForEach(func(key, val lua.LValue) {
			tags[key.String()] = val.String()
		})
	}

	// Get existing histogram or create new one
	var samples []float64
	if existing, exists := m.customMetrics[name]; exists {
		if existingMetric, ok := existing.(map[string]interface{}); ok {
			if vals, ok := existingMetric["samples"].([]float64); ok {
				samples = vals
			}
		}
	}

	// Add new sample (keep last 100 samples)
	samples = append(samples, float64(value))
	if len(samples) > 100 {
		samples = samples[len(samples)-100:]
	}

	// Calculate stats
	sum := 0.0
	min := samples[0]
	max := samples[0]
	for _, s := range samples {
		sum += s
		if s < min {
			min = s
		}
		if s > max {
			max = s
		}
	}
	avg := sum / float64(len(samples))

	metric := map[string]interface{}{
		"type":      "histogram",
		"samples":   samples,
		"count":     len(samples),
		"sum":       sum,
		"avg":       avg,
		"min":       min,
		"max":       max,
		"tags":      tags,
		"timestamp": time.Now(),
	}

	m.customMetrics[name] = metric
	L.Push(lua.LBool(true))
	return 1
}

// luaTimer provides timing functionality
func (m *MetricsModule) luaTimer(L *lua.LState) int {
	name := L.CheckString(1)
	fn := L.CheckFunction(2)

	tags := make(map[string]string)
	if L.GetTop() >= 3 {
		tagsTable := L.CheckTable(3)
		tagsTable.ForEach(func(key, val lua.LValue) {
			tags[key.String()] = val.String()
		})
	}

	// Time the function execution
	start := time.Now()
	
	// Call the function
	L.Push(fn)
	err := L.PCall(0, lua.MultRet, nil)
	
	duration := time.Since(start)
	
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Function execution failed: %v", err)))
		return 2
	}

	// Record the timing as a histogram
	metric := map[string]interface{}{
		"type":      "timer",
		"duration_ms": float64(duration.Nanoseconds()) / 1000000,
		"tags":      tags,
		"timestamp": time.Now(),
	}

	m.customMetrics[name] = metric
	
	// Return the duration in milliseconds
	L.Push(lua.LNumber(float64(duration.Nanoseconds()) / 1000000))
	return 1
}

// luaGetCustom gets a custom metric value
func (m *MetricsModule) luaGetCustom(L *lua.LState) int {
	name := L.CheckString(1)

	if metric, exists := m.customMetrics[name]; exists {
		if metricMap, ok := metric.(map[string]interface{}); ok {
			// Convert to Lua table
			table := L.NewTable()
			for key, value := range metricMap {
				switch v := value.(type) {
				case string:
					table.RawSetString(key, lua.LString(v))
				case float64:
					table.RawSetString(key, lua.LNumber(v))
				case int:
					table.RawSetString(key, lua.LNumber(v))
				case bool:
					table.RawSetString(key, lua.LBool(v))
				case time.Time:
					table.RawSetString(key, lua.LString(v.Format(time.RFC3339)))
				case map[string]string:
					tagsTable := L.NewTable()
					for k, val := range v {
						tagsTable.RawSetString(k, lua.LString(val))
					}
					table.RawSetString(key, tagsTable)
				}
			}
			L.Push(table)
		} else {
			L.Push(lua.LNil)
		}
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

// luaListCustom lists all custom metrics
func (m *MetricsModule) luaListCustom(L *lua.LState) int {
	metricsTable := L.NewTable()
	
	i := 1
	for name := range m.customMetrics {
		metricsTable.RawSetInt(i, lua.LString(name))
		i++
	}
	
	L.Push(metricsTable)
	return 1
}

// luaAlert creates an alert based on metrics
func (m *MetricsModule) luaAlert(L *lua.LState) int {
	name := L.CheckString(1)
	alertData := L.CheckTable(2)

	alert := make(map[string]interface{})
	alertData.ForEach(func(key, value lua.LValue) {
		switch key.String() {
		case "level":
			alert["level"] = value.String()
		case "message":
			alert["message"] = value.String()
		case "threshold":
			alert["threshold"] = float64(lua.LVAsNumber(value))
		case "value":
			alert["value"] = float64(lua.LVAsNumber(value))
		default:
			alert[key.String()] = value.String()
		}
	})

	alert["timestamp"] = time.Now()
	alert["type"] = "alert"

	alertName := fmt.Sprintf("alert_%s", name)
	m.customMetrics[alertName] = alert

	// Log the alert
	level := alert["level"].(string)
	message := alert["message"].(string)
	
	fmt.Printf("[ALERT] %s: %s\n", level, message)

	L.Push(lua.LBool(true))
	return 1
}

// luaHealthStatus provides a health check based on system metrics
func (m *MetricsModule) luaHealthStatus(L *lua.LState) int {
	healthTable := L.NewTable()
	overallStatus := "healthy"

	// Check CPU
	if cpuPercents, err := cpu.Percent(time.Second, false); err == nil && len(cpuPercents) > 0 {
		cpuUsage := cpuPercents[0]
		cpuStatus := "healthy"
		if cpuUsage > 90 {
			cpuStatus = "critical"
			overallStatus = "critical"
		} else if cpuUsage > 70 {
			cpuStatus = "warning"
			if overallStatus == "healthy" {
				overallStatus = "warning"
			}
		}
		
		cpuTable := L.NewTable()
		cpuTable.RawSetString("usage", lua.LNumber(cpuUsage))
		cpuTable.RawSetString("status", lua.LString(cpuStatus))
		healthTable.RawSetString("cpu", cpuTable)
	}

	// Check Memory
	if memInfo, err := mem.VirtualMemory(); err == nil {
		memUsage := memInfo.UsedPercent
		memStatus := "healthy"
		if memUsage > 90 {
			memStatus = "critical"
			overallStatus = "critical"
		} else if memUsage > 80 {
			memStatus = "warning"
			if overallStatus == "healthy" {
				overallStatus = "warning"
			}
		}
		
		memTable := L.NewTable()
		memTable.RawSetString("usage", lua.LNumber(memUsage))
		memTable.RawSetString("status", lua.LString(memStatus))
		healthTable.RawSetString("memory", memTable)
	}

	// Check Disk
	if diskInfo, err := disk.Usage("/"); err == nil {
		diskUsage := diskInfo.UsedPercent
		diskStatus := "healthy"
		if diskUsage > 95 {
			diskStatus = "critical"
			overallStatus = "critical"
		} else if diskUsage > 85 {
			diskStatus = "warning"
			if overallStatus == "healthy" {
				overallStatus = "warning"
			}
		}
		
		diskTable := L.NewTable()
		diskTable.RawSetString("usage", lua.LNumber(diskUsage))
		diskTable.RawSetString("status", lua.LString(diskStatus))
		healthTable.RawSetString("disk", diskTable)
	}

	healthTable.RawSetString("overall", lua.LString(overallStatus))
	healthTable.RawSetString("timestamp", lua.LString(time.Now().Format(time.RFC3339)))

	L.Push(healthTable)
	return 1
}

// OpenMetrics initializes the metrics module in Lua
func OpenMetrics(L *lua.LState) {
	L.PreloadModule("metrics", MetricsLoader)
	if err := L.DoString(`metrics = require("metrics")`); err != nil {
		panic(err)
	}
}