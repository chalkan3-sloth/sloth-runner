package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestObservabilityModule(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterObservabilityModule(L)

	tests := []struct {
		name   string
		script string
	}{
		{
			name: "observability.log()",
			script: `
				observability.log("info", "test message")
			`,
		},
		{
			name: "observability.info()",
			script: `
				observability.info("info message")
			`,
		},
		{
			name: "observability.warn()",
			script: `
				observability.warn("warning message")
			`,
		},
		{
			name: "observability.error()",
			script: `
				observability.error("error message")
			`,
		},
		{
			name: "observability.debug()",
			script: `
				observability.debug("debug message")
			`,
		},
		{
			name: "observability.trace()",
			script: `
				observability.trace("trace message")
			`,
		},
		{
			name: "observability.metric()",
			script: `
				observability.metric("counter", "test_metric", 1)
			`,
		},
		{
			name: "observability.counter()",
			script: `
				observability.counter("test_counter", 1)
			`,
		},
		{
			name: "observability.gauge()",
			script: `
				observability.gauge("test_gauge", 100)
			`,
		},
		{
			name: "observability.histogram()",
			script: `
				observability.histogram("test_histogram", 0.5)
			`,
		},
		{
			name: "observability.timer_start()",
			script: `
				local timer_id = observability.timer_start("test_timer")
				assert(type(timer_id) == "string", "timer_start() should return a string")
			`,
		},
		{
			name: "observability.timer_end()",
			script: `
				local timer_id = observability.timer_start("test_timer")
				local duration = observability.timer_end(timer_id)
				assert(type(duration) == "number", "timer_end() should return a number")
			`,
		},
		{
			name: "observability.span_start()",
			script: `
				local span_id = observability.span_start("test_span")
				assert(type(span_id) == "string", "span_start() should return a string")
			`,
		},
		{
			name: "observability.span_end()",
			script: `
				local span_id = observability.span_start("test_span")
				observability.span_end(span_id)
			`,
		},
		{
			name: "observability.add_tag()",
			script: `
				local span_id = observability.span_start("test_span")
				observability.add_tag(span_id, "key", "value")
			`,
		},
		{
			name: "observability.add_event()",
			script: `
				local span_id = observability.span_start("test_span")
				observability.add_event(span_id, "test_event")
			`,
		},
		{
			name: "observability.set_status()",
			script: `
				local span_id = observability.span_start("test_span")
				observability.set_status(span_id, "ok")
			`,
		},
		{
			name: "observability.checkpoint()",
			script: `
				observability.checkpoint("test_checkpoint")
			`,
		},
		{
			name: "observability.measure()",
			script: `
				local duration = observability.measure(function()
					-- Some work
					local sum = 0
					for i = 1, 100 do
						sum = sum + i
					end
				end)
				assert(type(duration) == "number", "measure() should return a number")
			`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := L.DoString(tt.script); err != nil {
				t.Fatalf("Failed to execute script: %v", err)
			}
		})
	}
}

func TestObservabilityTimerFlow(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterObservabilityModule(L)

	script := `
		local timer_id1 = observability.timer_start("operation1")
		local timer_id2 = observability.timer_start("operation2")
		
		-- Simulate some work
		local sum = 0
		for i = 1, 1000 do
			sum = sum + i
		end
		
		local duration1 = observability.timer_end(timer_id1)
		local duration2 = observability.timer_end(timer_id2)
		
		assert(duration1 >= 0, "duration1 should be non-negative")
		assert(duration2 >= 0, "duration2 should be non-negative")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestObservabilitySpanFlow(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterObservabilityModule(L)

	script := `
		local span_id = observability.span_start("main_operation")
		
		observability.add_tag(span_id, "environment", "test")
		observability.add_tag(span_id, "version", "1.0.0")
		
		observability.add_event(span_id, "step_1_complete")
		observability.add_event(span_id, "step_2_complete")
		
		observability.set_status(span_id, "ok")
		
		observability.span_end(span_id)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestObservabilityNestedSpans(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterObservabilityModule(L)

	script := `
		local parent_span = observability.span_start("parent_operation")
		
		local child_span1 = observability.span_start("child_operation_1")
		observability.span_end(child_span1)
		
		local child_span2 = observability.span_start("child_operation_2")
		observability.span_end(child_span2)
		
		observability.span_end(parent_span)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestObservabilityMeasureFunction(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterObservabilityModule(L)

	script := `
		local result = nil
		local duration = observability.measure(function()
			-- Calculate fibonacci
			local function fib(n)
				if n <= 1 then return n end
				return fib(n-1) + fib(n-2)
			end
			result = fib(10)
		end)
		
		assert(result == 55, "fibonacci(10) should be 55")
		assert(duration >= 0, "duration should be non-negative")
		assert(type(duration) == "number", "duration should be a number")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestObservabilityMetricTypes(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterObservabilityModule(L)

	script := `
		-- Test counter
		observability.counter("requests_total", 1)
		observability.counter("requests_total", 5)
		observability.counter("requests_total", 10)
		
		-- Test gauge
		observability.gauge("memory_usage", 1024)
		observability.gauge("memory_usage", 2048)
		observability.gauge("cpu_usage", 75.5)
		
		-- Test histogram
		observability.histogram("request_duration", 0.1)
		observability.histogram("request_duration", 0.5)
		observability.histogram("request_duration", 1.2)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestObservabilityLogLevels(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterObservabilityModule(L)

	script := `
		observability.trace("trace level message")
		observability.debug("debug level message")
		observability.info("info level message")
		observability.warn("warning level message")
		observability.error("error level message")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestObservabilityCheckpoint(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterObservabilityModule(L)

	script := `
		observability.checkpoint("start")
		
		-- Do some work
		local sum = 0
		for i = 1, 100 do
			sum = sum + i
		end
		
		observability.checkpoint("middle")
		
		-- Do more work
		for i = 1, 100 do
			sum = sum + i
		end
		
		observability.checkpoint("end")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}

func TestObservabilityWithContext(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterObservabilityModule(L)

	script := `
		local span_id = observability.span_start("test_operation")
		
		-- Add multiple tags
		observability.add_tag(span_id, "service", "test-service")
		observability.add_tag(span_id, "environment", "test")
		observability.add_tag(span_id, "region", "us-west-2")
		
		-- Add multiple events
		observability.add_event(span_id, "validation_started")
		observability.add_event(span_id, "validation_complete")
		observability.add_event(span_id, "processing_started")
		
		-- Set final status
		observability.set_status(span_id, "ok")
		
		observability.span_end(span_id)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Failed to execute script: %v", err)
	}
}
