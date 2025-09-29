package luainterface

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	lua "github.com/yuin/gopher-lua"
)

// ObservabilityModule provides distributed tracing and observability
type ObservabilityModule struct {
	traces    map[string]*Trace
	spans     map[string]*Span
	metrics   map[string]*Metric
	traceMutex sync.RWMutex
	spanMutex  sync.RWMutex
	metricMutex sync.RWMutex
}

// Trace represents a distributed trace
type Trace struct {
	ID        string
	Name      string
	StartTime time.Time
	EndTime   *time.Time
	Status    string
	Tags      map[string]string
	Spans     []*Span
}

// Span represents a span within a trace
type Span struct {
	ID         string
	TraceID    string
	Name       string
	StartTime  time.Time
	EndTime    *time.Time
	Duration   time.Duration
	Status     string
	Tags       map[string]string
	Events     []SpanEvent
	ParentID   string
}

// SpanEvent represents an event within a span
type SpanEvent struct {
	Name      string
	Timestamp time.Time
	Tags      map[string]string
}

// Metric represents a metric measurement
type Metric struct {
	Name      string
	Type      string
	Value     float64
	Timestamp time.Time
	Tags      map[string]string
}

// NewObservabilityModule creates a new observability module
func NewObservabilityModule() *ObservabilityModule {
	return &ObservabilityModule{
		traces:  make(map[string]*Trace),
		spans:   make(map[string]*Span),
		metrics: make(map[string]*Metric),
	}
}

// RegisterObservabilityModule registers the observability module with the Lua state
func RegisterObservabilityModule(L *lua.LState) {
	module := NewObservabilityModule()
	
	// Create the observability table
	observabilityTable := L.NewTable()
	
	// Trace management
	L.SetField(observabilityTable, "start_trace", L.NewFunction(module.luaStartTrace))
	L.SetField(observabilityTable, "end_trace", L.NewFunction(module.luaEndTrace))
	L.SetField(observabilityTable, "get_trace", L.NewFunction(module.luaGetTrace))
	L.SetField(observabilityTable, "list_traces", L.NewFunction(module.luaListTraces))
	
	// Span management
	L.SetField(observabilityTable, "start_span", L.NewFunction(module.luaStartSpan))
	L.SetField(observabilityTable, "end_span", L.NewFunction(module.luaEndSpan))
	L.SetField(observabilityTable, "add_span_event", L.NewFunction(module.luaAddSpanEvent))
	L.SetField(observabilityTable, "add_span_tag", L.NewFunction(module.luaAddSpanTag))
	
	// Metrics
	L.SetField(observabilityTable, "counter", L.NewFunction(module.luaCounter))
	L.SetField(observabilityTable, "gauge", L.NewFunction(module.luaGauge))
	L.SetField(observabilityTable, "histogram", L.NewFunction(module.luaHistogram))
	L.SetField(observabilityTable, "timer_start", L.NewFunction(module.luaTimerStart))
	L.SetField(observabilityTable, "timer_end", L.NewFunction(module.luaTimerEnd))
	
	// Export functions
	L.SetField(observabilityTable, "export_jaeger", L.NewFunction(module.luaExportJaeger))
	L.SetField(observabilityTable, "export_prometheus", L.NewFunction(module.luaExportPrometheus))
	L.SetField(observabilityTable, "export_json", L.NewFunction(module.luaExportJSON))
	
	// Health and monitoring
	L.SetField(observabilityTable, "health_check", L.NewFunction(module.luaHealthCheck))
	L.SetField(observabilityTable, "system_metrics", L.NewFunction(module.luaSystemMetrics))
	
	// Store module reference
	ud := L.NewUserData()
	ud.Value = module
	L.SetGlobal("__observability_module", ud)
	
	// Register the observability table globally
	L.SetGlobal("observability", observabilityTable)
}

// Trace management
func (o *ObservabilityModule) luaStartTrace(L *lua.LState) int {
	name := L.CheckString(1)
	tags := L.OptTable(2, L.NewTable())
	
	traceID := uuid.New().String()
	
	// Convert tags table to Go map
	tagMap := make(map[string]string)
	tags.ForEach(func(key, value lua.LValue) {
		tagMap[lua.LVAsString(key)] = lua.LVAsString(value)
	})
	
	trace := &Trace{
		ID:        traceID,
		Name:      name,
		StartTime: time.Now(),
		Status:    "active",
		Tags:      tagMap,
		Spans:     make([]*Span, 0),
	}
	
	o.traceMutex.Lock()
	o.traces[traceID] = trace
	o.traceMutex.Unlock()
	
	L.Push(lua.LString(traceID))
	return 1
}

func (o *ObservabilityModule) luaEndTrace(L *lua.LState) int {
	traceID := L.CheckString(1)
	status := L.OptString(2, "completed")
	
	o.traceMutex.Lock()
	defer o.traceMutex.Unlock()
	
	trace, exists := o.traces[traceID]
	if !exists {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("trace not found"))
		return 2
	}
	
	endTime := time.Now()
	trace.EndTime = &endTime
	trace.Status = status
	
	L.Push(lua.LBool(true))
	L.Push(lua.LNumber(endTime.Sub(trace.StartTime).Milliseconds()))
	return 2
}

func (o *ObservabilityModule) luaGetTrace(L *lua.LState) int {
	traceID := L.CheckString(1)
	
	o.traceMutex.RLock()
	trace, exists := o.traces[traceID]
	o.traceMutex.RUnlock()
	
	if !exists {
		L.Push(lua.LNil)
		L.Push(lua.LString("trace not found"))
		return 2
	}
	
	result := o.traceToLuaTable(L, trace)
	L.Push(result)
	return 1
}

func (o *ObservabilityModule) luaListTraces(L *lua.LState) int {
	status := L.OptString(1, "all")
	
	o.traceMutex.RLock()
	defer o.traceMutex.RUnlock()
	
	result := L.NewTable()
	index := 1
	
	for _, trace := range o.traces {
		if status == "all" || trace.Status == status {
			traceInfo := L.NewTable()
			L.SetField(traceInfo, "id", lua.LString(trace.ID))
			L.SetField(traceInfo, "name", lua.LString(trace.Name))
			L.SetField(traceInfo, "status", lua.LString(trace.Status))
			L.SetField(traceInfo, "start_time", lua.LNumber(trace.StartTime.Unix()))
			
			if trace.EndTime != nil {
				L.SetField(traceInfo, "duration_ms", lua.LNumber(trace.EndTime.Sub(trace.StartTime).Milliseconds()))
			}
			
			result.RawSetInt(index, traceInfo)
			index++
		}
	}
	
	L.Push(result)
	return 1
}

// Span management
func (o *ObservabilityModule) luaStartSpan(L *lua.LState) int {
	traceID := L.CheckString(1)
	name := L.CheckString(2)
	parentID := L.OptString(3, "")
	tags := L.OptTable(4, L.NewTable())
	
	spanID := uuid.New().String()
	
	// Convert tags table to Go map
	tagMap := make(map[string]string)
	tags.ForEach(func(key, value lua.LValue) {
		tagMap[lua.LVAsString(key)] = lua.LVAsString(value)
	})
	
	span := &Span{
		ID:        spanID,
		TraceID:   traceID,
		Name:      name,
		StartTime: time.Now(),
		Status:    "active",
		Tags:      tagMap,
		Events:    make([]SpanEvent, 0),
		ParentID:  parentID,
	}
	
	o.spanMutex.Lock()
	o.spans[spanID] = span
	o.spanMutex.Unlock()
	
	// Add span to trace
	o.traceMutex.Lock()
	if trace, exists := o.traces[traceID]; exists {
		trace.Spans = append(trace.Spans, span)
	}
	o.traceMutex.Unlock()
	
	L.Push(lua.LString(spanID))
	return 1
}

func (o *ObservabilityModule) luaEndSpan(L *lua.LState) int {
	spanID := L.CheckString(1)
	status := L.OptString(2, "completed")
	
	o.spanMutex.Lock()
	defer o.spanMutex.Unlock()
	
	span, exists := o.spans[spanID]
	if !exists {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("span not found"))
		return 2
	}
	
	endTime := time.Now()
	span.EndTime = &endTime
	span.Duration = endTime.Sub(span.StartTime)
	span.Status = status
	
	L.Push(lua.LBool(true))
	L.Push(lua.LNumber(span.Duration.Milliseconds()))
	return 2
}

func (o *ObservabilityModule) luaAddSpanEvent(L *lua.LState) int {
	spanID := L.CheckString(1)
	eventName := L.CheckString(2)
	tags := L.OptTable(3, L.NewTable())
	
	o.spanMutex.Lock()
	defer o.spanMutex.Unlock()
	
	span, exists := o.spans[spanID]
	if !exists {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("span not found"))
		return 2
	}
	
	// Convert tags table to Go map
	tagMap := make(map[string]string)
	tags.ForEach(func(key, value lua.LValue) {
		tagMap[lua.LVAsString(key)] = lua.LVAsString(value)
	})
	
	event := SpanEvent{
		Name:      eventName,
		Timestamp: time.Now(),
		Tags:      tagMap,
	}
	
	span.Events = append(span.Events, event)
	
	L.Push(lua.LBool(true))
	return 1
}

func (o *ObservabilityModule) luaAddSpanTag(L *lua.LState) int {
	spanID := L.CheckString(1)
	key := L.CheckString(2)
	value := L.CheckString(3)
	
	o.spanMutex.Lock()
	defer o.spanMutex.Unlock()
	
	span, exists := o.spans[spanID]
	if !exists {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("span not found"))
		return 2
	}
	
	span.Tags[key] = value
	
	L.Push(lua.LBool(true))
	return 1
}

// Metrics
func (o *ObservabilityModule) luaCounter(L *lua.LState) int {
	name := L.CheckString(1)
	value := L.OptNumber(2, 1)
	tags := L.OptTable(3, L.NewTable())
	
	// Convert tags table to Go map
	tagMap := make(map[string]string)
	tags.ForEach(func(key, val lua.LValue) {
		tagMap[lua.LVAsString(key)] = lua.LVAsString(val)
	})
	
	metric := &Metric{
		Name:      name,
		Type:      "counter",
		Value:     float64(value),
		Timestamp: time.Now(),
		Tags:      tagMap,
	}
	
	o.metricMutex.Lock()
	metricKey := fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
	o.metrics[metricKey] = metric
	o.metricMutex.Unlock()
	
	L.Push(lua.LBool(true))
	return 1
}

func (o *ObservabilityModule) luaGauge(L *lua.LState) int {
	name := L.CheckString(1)
	value := L.CheckNumber(2)
	tags := L.OptTable(3, L.NewTable())
	
	// Convert tags table to Go map
	tagMap := make(map[string]string)
	tags.ForEach(func(key, val lua.LValue) {
		tagMap[lua.LVAsString(key)] = lua.LVAsString(val)
	})
	
	metric := &Metric{
		Name:      name,
		Type:      "gauge",
		Value:     float64(value),
		Timestamp: time.Now(),
		Tags:      tagMap,
	}
	
	o.metricMutex.Lock()
	metricKey := fmt.Sprintf("%s_gauge", name)
	o.metrics[metricKey] = metric
	o.metricMutex.Unlock()
	
	L.Push(lua.LBool(true))
	return 1
}

func (o *ObservabilityModule) luaHistogram(L *lua.LState) int {
	name := L.CheckString(1)
	value := L.CheckNumber(2)
	tags := L.OptTable(3, L.NewTable())
	
	// Convert tags table to Go map
	tagMap := make(map[string]string)
	tags.ForEach(func(key, val lua.LValue) {
		tagMap[lua.LVAsString(key)] = lua.LVAsString(val)
	})
	
	metric := &Metric{
		Name:      name,
		Type:      "histogram",
		Value:     float64(value),
		Timestamp: time.Now(),
		Tags:      tagMap,
	}
	
	o.metricMutex.Lock()
	metricKey := fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
	o.metrics[metricKey] = metric
	o.metricMutex.Unlock()
	
	L.Push(lua.LBool(true))
	return 1
}

func (o *ObservabilityModule) luaTimerStart(L *lua.LState) int {
	name := L.CheckString(1)
	
	timerID := uuid.New().String()
	
	// Store timer start time
	timer := &Metric{
		Name:      name,
		Type:      "timer_start",
		Value:     float64(time.Now().UnixNano()),
		Timestamp: time.Now(),
		Tags:      map[string]string{"timer_id": timerID},
	}
	
	o.metricMutex.Lock()
	o.metrics[timerID] = timer
	o.metricMutex.Unlock()
	
	L.Push(lua.LString(timerID))
	return 1
}

func (o *ObservabilityModule) luaTimerEnd(L *lua.LState) int {
	timerID := L.CheckString(1)
	
	o.metricMutex.Lock()
	defer o.metricMutex.Unlock()
	
	startTimer, exists := o.metrics[timerID]
	if !exists {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("timer not found"))
		return 2
	}
	
	startTime := time.Unix(0, int64(startTimer.Value))
	duration := time.Since(startTime)
	
	// Create duration metric
	durationMetric := &Metric{
		Name:      startTimer.Name,
		Type:      "timer_duration",
		Value:     float64(duration.Milliseconds()),
		Timestamp: time.Now(),
		Tags:      startTimer.Tags,
	}
	
	metricKey := fmt.Sprintf("%s_duration_%d", startTimer.Name, time.Now().UnixNano())
	o.metrics[metricKey] = durationMetric
	
	// Remove start timer
	delete(o.metrics, timerID)
	
	L.Push(lua.LBool(true))
	L.Push(lua.LNumber(duration.Milliseconds()))
	return 2
}

// Export functions
func (o *ObservabilityModule) luaExportJaeger(L *lua.LState) int {
	endpoint := L.OptString(1, "http://localhost:14268/api/traces")
	
	o.traceMutex.RLock()
	traceCount := len(o.traces)
	o.traceMutex.RUnlock()
	
	// Placeholder implementation - would normally export to Jaeger
	result := L.NewTable()
	L.SetField(result, "success", lua.LBool(true))
	L.SetField(result, "endpoint", lua.LString(endpoint))
	L.SetField(result, "exported_traces", lua.LNumber(traceCount))
	L.SetField(result, "note", lua.LString("Jaeger export requires Jaeger client library"))
	
	L.Push(result)
	return 1
}

func (o *ObservabilityModule) luaExportPrometheus(L *lua.LState) int {
	endpoint := L.OptString(1, "http://localhost:9090/api/v1/write")
	
	o.metricMutex.RLock()
	metricCount := len(o.metrics)
	o.metricMutex.RUnlock()
	
	// Placeholder implementation - would normally export to Prometheus
	result := L.NewTable()
	L.SetField(result, "success", lua.LBool(true))
	L.SetField(result, "endpoint", lua.LString(endpoint))
	L.SetField(result, "exported_metrics", lua.LNumber(metricCount))
	L.SetField(result, "note", lua.LString("Prometheus export requires Prometheus client library"))
	
	L.Push(result)
	return 1
}

func (o *ObservabilityModule) luaExportJSON(L *lua.LState) int {
	result := L.NewTable()
	
	// Export traces
	traces := L.NewTable()
	o.traceMutex.RLock()
	traceIndex := 1
	for _, trace := range o.traces {
		traceTable := o.traceToLuaTable(L, trace)
		traces.RawSetInt(traceIndex, traceTable)
		traceIndex++
	}
	o.traceMutex.RUnlock()
	
	// Export metrics
	metrics := L.NewTable()
	o.metricMutex.RLock()
	metricIndex := 1
	for _, metric := range o.metrics {
		metricTable := L.NewTable()
		L.SetField(metricTable, "name", lua.LString(metric.Name))
		L.SetField(metricTable, "type", lua.LString(metric.Type))
		L.SetField(metricTable, "value", lua.LNumber(metric.Value))
		L.SetField(metricTable, "timestamp", lua.LNumber(metric.Timestamp.Unix()))
		
		// Tags
		tagsTable := L.NewTable()
		for k, v := range metric.Tags {
			tagsTable.RawSetString(k, lua.LString(v))
		}
		L.SetField(metricTable, "tags", tagsTable)
		
		metrics.RawSetInt(metricIndex, metricTable)
		metricIndex++
	}
	o.metricMutex.RUnlock()
	
	L.SetField(result, "traces", traces)
	L.SetField(result, "metrics", metrics)
	L.SetField(result, "exported_at", lua.LNumber(time.Now().Unix()))
	
	L.Push(result)
	return 1
}

// Health and monitoring
func (o *ObservabilityModule) luaHealthCheck(L *lua.LState) int {
	result := L.NewTable()
	
	// Check trace count
	o.traceMutex.RLock()
	activeTraces := 0
	for _, trace := range o.traces {
		if trace.Status == "active" {
			activeTraces++
		}
	}
	totalTraces := len(o.traces)
	o.traceMutex.RUnlock()
	
	// Check metric count
	o.metricMutex.RLock()
	totalMetrics := len(o.metrics)
	o.metricMutex.RUnlock()
	
	L.SetField(result, "status", lua.LString("healthy"))
	L.SetField(result, "active_traces", lua.LNumber(activeTraces))
	L.SetField(result, "total_traces", lua.LNumber(totalTraces))
	L.SetField(result, "total_metrics", lua.LNumber(totalMetrics))
	L.SetField(result, "timestamp", lua.LNumber(time.Now().Unix()))
	
	L.Push(result)
	return 1
}

func (o *ObservabilityModule) luaSystemMetrics(L *lua.LState) int {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	result := L.NewTable()
	
	// Memory metrics
	memory := L.NewTable()
	L.SetField(memory, "alloc", lua.LNumber(memStats.Alloc))
	L.SetField(memory, "total_alloc", lua.LNumber(memStats.TotalAlloc))
	L.SetField(memory, "sys", lua.LNumber(memStats.Sys))
	L.SetField(memory, "gc_runs", lua.LNumber(memStats.NumGC))
	L.SetField(result, "memory", memory)
	
	// Runtime metrics
	runtimeInfo := L.NewTable()
	L.SetField(runtimeInfo, "goroutines", lua.LNumber(runtime.NumGoroutine()))
	L.SetField(runtimeInfo, "cpus", lua.LNumber(runtime.NumCPU()))
	L.SetField(runtimeInfo, "go_version", lua.LString(runtime.Version()))
	L.SetField(result, "runtime", runtimeInfo)
	
	L.SetField(result, "timestamp", lua.LNumber(time.Now().Unix()))
	
	L.Push(result)
	return 1
}

// Helper functions
func (o *ObservabilityModule) traceToLuaTable(L *lua.LState, trace *Trace) *lua.LTable {
	traceTable := L.NewTable()
	L.SetField(traceTable, "id", lua.LString(trace.ID))
	L.SetField(traceTable, "name", lua.LString(trace.Name))
	L.SetField(traceTable, "start_time", lua.LNumber(trace.StartTime.Unix()))
	L.SetField(traceTable, "status", lua.LString(trace.Status))
	
	if trace.EndTime != nil {
		L.SetField(traceTable, "end_time", lua.LNumber(trace.EndTime.Unix()))
		L.SetField(traceTable, "duration_ms", lua.LNumber(trace.EndTime.Sub(trace.StartTime).Milliseconds()))
	}
	
	// Tags
	tags := L.NewTable()
	for k, v := range trace.Tags {
		tags.RawSetString(k, lua.LString(v))
	}
	L.SetField(traceTable, "tags", tags)
	
	// Spans
	spans := L.NewTable()
	for i, span := range trace.Spans {
		spanTable := L.NewTable()
		L.SetField(spanTable, "id", lua.LString(span.ID))
		L.SetField(spanTable, "name", lua.LString(span.Name))
		L.SetField(spanTable, "start_time", lua.LNumber(span.StartTime.Unix()))
		L.SetField(spanTable, "status", lua.LString(span.Status))
		L.SetField(spanTable, "parent_id", lua.LString(span.ParentID))
		
		if span.EndTime != nil {
			L.SetField(spanTable, "end_time", lua.LNumber(span.EndTime.Unix()))
			L.SetField(spanTable, "duration_ms", lua.LNumber(span.Duration.Milliseconds()))
		}
		
		spans.RawSetInt(i+1, spanTable)
	}
	L.SetField(traceTable, "spans", spans)
	
	return traceTable
}