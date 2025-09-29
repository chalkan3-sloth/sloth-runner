package core

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/yuin/gopher-lua"
)

// MonitoringModule provides monitoring and metrics functionality
type MonitoringModule struct {
	info    CoreModuleInfo
	metrics map[string]*Metric
	mutex   sync.RWMutex
}

// Metric represents a monitoring metric
type Metric struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"` // counter, gauge, histogram, timer
	Value       float64                `json:"value"`
	Labels      map[string]string      `json:"labels,omitempty"`
	LastUpdated time.Time              `json:"last_updated"`
	History     []MetricPoint          `json:"history,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// MetricPoint represents a point in time for a metric
type MetricPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

// NewMonitoringModule creates a new monitoring module
func NewMonitoringModule() *MonitoringModule {
	info := CoreModuleInfo{
		Name:         "monitor",
		Version:      "1.0.0",
		Description:  "Monitoring and metrics collection with support for counters, gauges, histograms",
		Author:       "Sloth Runner Team",
		Category:     "core",
		Dependencies: []string{},
	}

	return &MonitoringModule{
		info:    info,
		metrics: make(map[string]*Metric),
	}
}

// Info returns module information
func (m *MonitoringModule) Info() CoreModuleInfo {
	return m.info
}

// Loader loads the monitoring module into Lua
func (m *MonitoringModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"counter_inc":     m.luaCounterInc,
		"counter_add":     m.luaCounterAdd,
		"gauge_set":       m.luaGaugeSet,
		"gauge_inc":       m.luaGaugeInc,
		"gauge_dec":       m.luaGaugeDec,
		"timer_start":     m.luaTimerStart,
		"timer_end":       m.luaTimerEnd,
		"histogram_observe": m.luaHistogramObserve,
		"get_metric":      m.luaGetMetric,
		"list_metrics":    m.luaListMetrics,
		"reset_metric":    m.luaResetMetric,
		"clear_all":       m.luaClearAll,
		"export_prometheus": m.luaExportPrometheus,
		"export_json":     m.luaExportJSON,
		"system_metrics":  m.luaSystemMetrics,
		"memory_stats":    m.luaMemoryStats,
	})

	L.Push(mod)
	return 1
}

// luaCounterInc increments a counter by 1
func (m *MonitoringModule) luaCounterInc(L *lua.LState) int {
	name := L.CheckString(1)
	labels := m.parseLabels(L, 2)
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	key := m.generateKey(name, labels)
	metric, exists := m.metrics[key]
	if !exists {
		metric = &Metric{
			Name:        name,
			Type:        "counter",
			Value:       0,
			Labels:      labels,
			LastUpdated: time.Now(),
			History:     []MetricPoint{},
			Metadata:    make(map[string]interface{}),
		}
		m.metrics[key] = metric
	}
	
	metric.Value++
	metric.LastUpdated = time.Now()
	m.addToHistory(metric)
	
	L.Push(lua.LNumber(metric.Value))
	return 1
}

// luaCounterAdd adds a value to a counter
func (m *MonitoringModule) luaCounterAdd(L *lua.LState) int {
	name := L.CheckString(1)
	value := float64(L.CheckNumber(2))
	labels := m.parseLabels(L, 3)
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	key := m.generateKey(name, labels)
	metric, exists := m.metrics[key]
	if !exists {
		metric = &Metric{
			Name:        name,
			Type:        "counter",
			Value:       0,
			Labels:      labels,
			LastUpdated: time.Now(),
			History:     []MetricPoint{},
			Metadata:    make(map[string]interface{}),
		}
		m.metrics[key] = metric
	}
	
	metric.Value += value
	metric.LastUpdated = time.Now()
	m.addToHistory(metric)
	
	L.Push(lua.LNumber(metric.Value))
	return 1
}

// luaGaugeSet sets a gauge value
func (m *MonitoringModule) luaGaugeSet(L *lua.LState) int {
	name := L.CheckString(1)
	value := float64(L.CheckNumber(2))
	labels := m.parseLabels(L, 3)
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	key := m.generateKey(name, labels)
	metric, exists := m.metrics[key]
	if !exists {
		metric = &Metric{
			Name:        name,
			Type:        "gauge",
			Value:       0,
			Labels:      labels,
			LastUpdated: time.Now(),
			History:     []MetricPoint{},
			Metadata:    make(map[string]interface{}),
		}
		m.metrics[key] = metric
	}
	
	metric.Value = value
	metric.LastUpdated = time.Now()
	m.addToHistory(metric)
	
	L.Push(lua.LNumber(metric.Value))
	return 1
}

// luaGaugeInc increments a gauge
func (m *MonitoringModule) luaGaugeInc(L *lua.LState) int {
	name := L.CheckString(1)
	value := 1.0
	if L.GetTop() > 1 {
		value = float64(L.CheckNumber(2))
	}
	labels := m.parseLabels(L, 3)
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	key := m.generateKey(name, labels)
	metric, exists := m.metrics[key]
	if !exists {
		metric = &Metric{
			Name:        name,
			Type:        "gauge",
			Value:       0,
			Labels:      labels,
			LastUpdated: time.Now(),
			History:     []MetricPoint{},
			Metadata:    make(map[string]interface{}),
		}
		m.metrics[key] = metric
	}
	
	metric.Value += value
	metric.LastUpdated = time.Now()
	m.addToHistory(metric)
	
	L.Push(lua.LNumber(metric.Value))
	return 1
}

// luaGaugeDec decrements a gauge
func (m *MonitoringModule) luaGaugeDec(L *lua.LState) int {
	name := L.CheckString(1)
	value := 1.0
	if L.GetTop() > 1 {
		value = float64(L.CheckNumber(2))
	}
	labels := m.parseLabels(L, 3)
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	key := m.generateKey(name, labels)
	metric, exists := m.metrics[key]
	if !exists {
		metric = &Metric{
			Name:        name,
			Type:        "gauge",
			Value:       0,
			Labels:      labels,
			LastUpdated: time.Now(),
			History:     []MetricPoint{},
			Metadata:    make(map[string]interface{}),
		}
		m.metrics[key] = metric
	}
	
	metric.Value -= value
	metric.LastUpdated = time.Now()
	m.addToHistory(metric)
	
	L.Push(lua.LNumber(metric.Value))
	return 1
}

// luaTimerStart starts a timer
func (m *MonitoringModule) luaTimerStart(L *lua.LState) int {
	name := L.CheckString(1)
	labels := m.parseLabels(L, 2)
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	key := m.generateKey(name, labels)
	metric, exists := m.metrics[key]
	if !exists {
		metric = &Metric{
			Name:        name,
			Type:        "timer",
			Value:       0,
			Labels:      labels,
			LastUpdated: time.Now(),
			History:     []MetricPoint{},
			Metadata:    make(map[string]interface{}),
		}
		m.metrics[key] = metric
	}
	
	metric.Metadata["start_time"] = time.Now()
	L.Push(lua.LString(key))
	return 1
}

// luaTimerEnd ends a timer and records the duration
func (m *MonitoringModule) luaTimerEnd(L *lua.LState) int {
	timerKey := L.CheckString(1)
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	metric, exists := m.metrics[timerKey]
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("Timer not found"))
		return 2
	}
	
	startTime, ok := metric.Metadata["start_time"].(time.Time)
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString("Timer not started"))
		return 2
	}
	
	duration := time.Since(startTime).Seconds()
	metric.Value = duration
	metric.LastUpdated = time.Now()
	m.addToHistory(metric)
	
	L.Push(lua.LNumber(duration))
	return 1
}

// luaHistogramObserve observes a value in a histogram
func (m *MonitoringModule) luaHistogramObserve(L *lua.LState) int {
	name := L.CheckString(1)
	value := float64(L.CheckNumber(2))
	labels := m.parseLabels(L, 3)
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	key := m.generateKey(name, labels)
	metric, exists := m.metrics[key]
	if !exists {
		metric = &Metric{
			Name:        name,
			Type:        "histogram",
			Value:       0,
			Labels:      labels,
			LastUpdated: time.Now(),
			History:     []MetricPoint{},
			Metadata:    map[string]interface{}{
				"count": 0.0,
				"sum":   0.0,
			},
		}
		m.metrics[key] = metric
	}
	
	count := metric.Metadata["count"].(float64)
	sum := metric.Metadata["sum"].(float64)
	
	count++
	sum += value
	metric.Value = sum / count // average
	metric.Metadata["count"] = count
	metric.Metadata["sum"] = sum
	metric.LastUpdated = time.Now()
	m.addToHistory(metric)
	
	L.Push(lua.LNumber(metric.Value))
	return 1
}

// luaGetMetric gets a metric value
func (m *MonitoringModule) luaGetMetric(L *lua.LState) int {
	name := L.CheckString(1)
	labels := m.parseLabels(L, 2)
	
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	key := m.generateKey(name, labels)
	metric, exists := m.metrics[key]
	if !exists {
		L.Push(lua.LNil)
		return 1
	}
	
	result := L.NewTable()
	result.RawSetString("name", lua.LString(metric.Name))
	result.RawSetString("type", lua.LString(metric.Type))
	result.RawSetString("value", lua.LNumber(metric.Value))
	result.RawSetString("last_updated", lua.LNumber(metric.LastUpdated.Unix()))
	
	// Add labels
	if len(metric.Labels) > 0 {
		labelsTable := L.NewTable()
		for k, v := range metric.Labels {
			labelsTable.RawSetString(k, lua.LString(v))
		}
		result.RawSetString("labels", labelsTable)
	}
	
	L.Push(result)
	return 1
}

// luaListMetrics lists all metrics
func (m *MonitoringModule) luaListMetrics(L *lua.LState) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	result := L.NewTable()
	index := 1
	
	for _, metric := range m.metrics {
		metricTable := L.NewTable()
		metricTable.RawSetString("name", lua.LString(metric.Name))
		metricTable.RawSetString("type", lua.LString(metric.Type))
		metricTable.RawSetString("value", lua.LNumber(metric.Value))
		metricTable.RawSetString("last_updated", lua.LNumber(metric.LastUpdated.Unix()))
		
		result.RawSetInt(index, metricTable)
		index++
	}
	
	L.Push(result)
	return 1
}

// luaResetMetric resets a metric to zero
func (m *MonitoringModule) luaResetMetric(L *lua.LState) int {
	name := L.CheckString(1)
	labels := m.parseLabels(L, 2)
	
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	key := m.generateKey(name, labels)
	metric, exists := m.metrics[key]
	if !exists {
		L.Push(lua.LBool(false))
		return 1
	}
	
	metric.Value = 0
	metric.LastUpdated = time.Now()
	metric.History = []MetricPoint{}
	metric.Metadata = make(map[string]interface{})
	
	L.Push(lua.LBool(true))
	return 1
}

// luaClearAll clears all metrics
func (m *MonitoringModule) luaClearAll(L *lua.LState) int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	
	m.metrics = make(map[string]*Metric)
	L.Push(lua.LBool(true))
	return 1
}

// luaExportPrometheus exports metrics in Prometheus format
func (m *MonitoringModule) luaExportPrometheus(L *lua.LState) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	var output strings.Builder
	
	for _, metric := range m.metrics {
		// HELP line
		fmt.Fprintf(&output, "# HELP %s %s\n", metric.Name, metric.Name)
		// TYPE line
		fmt.Fprintf(&output, "# TYPE %s %s\n", metric.Name, metric.Type)
		
		// Metric line
		if len(metric.Labels) > 0 {
			var labelPairs []string
			for k, v := range metric.Labels {
				labelPairs = append(labelPairs, fmt.Sprintf(`%s="%s"`, k, v))
			}
			fmt.Fprintf(&output, "%s{%s} %g %d\n", 
				metric.Name, 
				strings.Join(labelPairs, ","), 
				metric.Value, 
				metric.LastUpdated.Unix()*1000)
		} else {
			fmt.Fprintf(&output, "%s %g %d\n", 
				metric.Name, 
				metric.Value, 
				metric.LastUpdated.Unix()*1000)
		}
	}
	
	L.Push(lua.LString(output.String()))
	return 1
}

// luaExportJSON exports metrics in JSON format
func (m *MonitoringModule) luaExportJSON(L *lua.LState) int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	jsonData, err := json.MarshalIndent(m.metrics, "", "  ")
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(string(jsonData)))
	return 1
}

// luaSystemMetrics collects basic system metrics
func (m *MonitoringModule) luaSystemMetrics(L *lua.LState) int {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	result := L.NewTable()
	result.RawSetString("goroutines", lua.LNumber(runtime.NumGoroutine()))
	result.RawSetString("cpu_count", lua.LNumber(runtime.NumCPU()))
	result.RawSetString("memory_alloc", lua.LNumber(mem.Alloc))
	result.RawSetString("memory_total_alloc", lua.LNumber(mem.TotalAlloc))
	result.RawSetString("memory_sys", lua.LNumber(mem.Sys))
	result.RawSetString("gc_count", lua.LNumber(mem.NumGC))
	
	L.Push(result)
	return 1
}

// luaMemoryStats gets detailed memory statistics
func (m *MonitoringModule) luaMemoryStats(L *lua.LState) int {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	result := L.NewTable()
	result.RawSetString("alloc", lua.LNumber(mem.Alloc))
	result.RawSetString("total_alloc", lua.LNumber(mem.TotalAlloc))
	result.RawSetString("sys", lua.LNumber(mem.Sys))
	result.RawSetString("lookups", lua.LNumber(mem.Lookups))
	result.RawSetString("mallocs", lua.LNumber(mem.Mallocs))
	result.RawSetString("frees", lua.LNumber(mem.Frees))
	result.RawSetString("heap_alloc", lua.LNumber(mem.HeapAlloc))
	result.RawSetString("heap_sys", lua.LNumber(mem.HeapSys))
	result.RawSetString("heap_idle", lua.LNumber(mem.HeapIdle))
	result.RawSetString("heap_inuse", lua.LNumber(mem.HeapInuse))
	result.RawSetString("heap_released", lua.LNumber(mem.HeapReleased))
	result.RawSetString("heap_objects", lua.LNumber(mem.HeapObjects))
	result.RawSetString("stack_inuse", lua.LNumber(mem.StackInuse))
	result.RawSetString("stack_sys", lua.LNumber(mem.StackSys))
	result.RawSetString("next_gc", lua.LNumber(mem.NextGC))
	result.RawSetString("last_gc", lua.LNumber(mem.LastGC))
	result.RawSetString("gc_cpu_fraction", lua.LNumber(mem.GCCPUFraction))
	result.RawSetString("gc_count", lua.LNumber(mem.NumGC))
	
	L.Push(result)
	return 1
}

// Helper functions

func (m *MonitoringModule) parseLabels(L *lua.LState, index int) map[string]string {
	labels := make(map[string]string)
	if L.GetTop() >= index {
		if labelsTable := L.CheckTable(index); labelsTable != nil {
			labelsTable.ForEach(func(key, value lua.LValue) {
				labels[key.String()] = value.String()
			})
		}
	}
	return labels
}

func (m *MonitoringModule) generateKey(name string, labels map[string]string) string {
	if len(labels) == 0 {
		return name
	}
	
	var labelPairs []string
	for k, v := range labels {
		labelPairs = append(labelPairs, fmt.Sprintf("%s=%s", k, v))
	}
	return fmt.Sprintf("%s{%s}", name, strings.Join(labelPairs, ","))
}

func (m *MonitoringModule) addToHistory(metric *Metric) {
	point := MetricPoint{
		Timestamp: time.Now(),
		Value:     metric.Value,
	}
	
	metric.History = append(metric.History, point)
	
	// Keep only last 100 points
	if len(metric.History) > 100 {
		metric.History = metric.History[1:]
	}
}