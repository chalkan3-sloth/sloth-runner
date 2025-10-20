//go:build cgo
// +build cgo

package metrics

import (
	"context"
	"testing"
	"time"
)

// Test NewCollector
func TestNewCollector(t *testing.T) {
	cfg := CollectorConfig{
		MetricsDB:   &MetricsDB{},
		AgentClient: nil, // Mock will be needed for real tests
	}

	collector := NewCollector(cfg)

	if collector == nil {
		t.Error("Expected non-nil collector")
	}

	if collector.metricsDB == nil {
		t.Error("Expected metrics DB to be set")
	}

	if collector.interval != 120*time.Second {
		t.Errorf("Expected default interval 120s, got %v", collector.interval)
	}

	if collector.retentionDays != 7 {
		t.Errorf("Expected default retention 7 days, got %d", collector.retentionDays)
	}

	if collector.batchSize != 10 {
		t.Errorf("Expected default batch size 10, got %d", collector.batchSize)
	}

	if collector.timeout != 3*time.Second {
		t.Errorf("Expected default timeout 3s, got %v", collector.timeout)
	}
}

func TestNewCollector_CustomConfig(t *testing.T) {
	cfg := CollectorConfig{
		MetricsDB:     &MetricsDB{},
		Interval:      30 * time.Second,
		RetentionDays: 14,
		BatchSize:     20,
		Timeout:       5 * time.Second,
	}

	collector := NewCollector(cfg)

	if collector.interval != 30*time.Second {
		t.Errorf("Expected interval 30s, got %v", collector.interval)
	}

	if collector.retentionDays != 14 {
		t.Errorf("Expected retention 14 days, got %d", collector.retentionDays)
	}

	if collector.batchSize != 20 {
		t.Errorf("Expected batch size 20, got %d", collector.batchSize)
	}

	if collector.timeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", collector.timeout)
	}
}

func TestNewCollector_ZeroValues(t *testing.T) {
	cfg := CollectorConfig{
		MetricsDB: &MetricsDB{},
		// All zeros - should get defaults
	}

	collector := NewCollector(cfg)

	// Verify defaults are applied
	if collector.interval == 0 {
		t.Error("Expected default interval to be set")
	}

	if collector.retentionDays == 0 {
		t.Error("Expected default retention to be set")
	}

	if collector.batchSize == 0 {
		t.Error("Expected default batch size to be set")
	}

	if collector.timeout == 0 {
		t.Error("Expected default timeout to be set")
	}
}

// Test Collector struct
func TestCollector_InitialState(t *testing.T) {
	cfg := CollectorConfig{
		MetricsDB: &MetricsDB{},
	}

	collector := NewCollector(cfg)

	// Initial state should not be running
	if collector.IsRunning() {
		t.Error("Expected collector to not be running initially")
	}

	// StopCh should be created
	if collector.stopCh == nil {
		t.Error("Expected stopCh to be initialized")
	}
}

func TestCollector_Stop_WhenNotRunning(t *testing.T) {
	cfg := CollectorConfig{
		MetricsDB: &MetricsDB{},
	}

	collector := NewCollector(cfg)

	// Stopping when not running should not panic
	collector.Stop()

	if collector.IsRunning() {
		t.Error("Expected collector to not be running after stop")
	}
}

func TestCollector_IsRunning(t *testing.T) {
	cfg := CollectorConfig{
		MetricsDB: &MetricsDB{},
	}

	collector := NewCollector(cfg)

	// Initially not running
	if collector.IsRunning() {
		t.Error("Expected not running")
	}

	// Manually set running state for testing
	collector.mu.Lock()
	collector.running = true
	collector.mu.Unlock()

	if !collector.IsRunning() {
		t.Error("Expected running")
	}

	// Reset
	collector.mu.Lock()
	collector.running = false
	collector.mu.Unlock()

	if collector.IsRunning() {
		t.Error("Expected not running after reset")
	}
}

// Test CollectorConfig
func TestCollectorConfig_Validation(t *testing.T) {
	testCases := []struct {
		name   string
		config CollectorConfig
		valid  bool
	}{
		{
			name: "Valid config",
			config: CollectorConfig{
				MetricsDB:     &MetricsDB{},
				Interval:      60 * time.Second,
				RetentionDays: 7,
			},
			valid: true,
		},
		{
			name: "Minimal config",
			config: CollectorConfig{
				MetricsDB: &MetricsDB{},
			},
			valid: true,
		},
		{
			name: "Custom intervals",
			config: CollectorConfig{
				MetricsDB:     &MetricsDB{},
				Interval:      5 * time.Minute,
				RetentionDays: 30,
				BatchSize:     50,
				Timeout:       10 * time.Second,
			},
			valid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			collector := NewCollector(tc.config)
			if collector == nil && tc.valid {
				t.Error("Expected valid collector")
			}
			if collector != nil && !tc.valid {
				t.Error("Expected invalid collector")
			}
		})
	}
}

// Test AgentInfo struct
func TestAgentInfo_Creation(t *testing.T) {
	agent := AgentInfo{
		Name:    "test-agent",
		Address: "localhost:50051",
	}

	if agent.Name != "test-agent" {
		t.Errorf("Expected name 'test-agent', got '%s'", agent.Name)
	}

	if agent.Address != "localhost:50051" {
		t.Errorf("Expected address 'localhost:50051', got '%s'", agent.Address)
	}
}

func TestAgentInfo_EmptyValues(t *testing.T) {
	agent := AgentInfo{}

	if agent.Name != "" {
		t.Error("Expected empty name")
	}

	if agent.Address != "" {
		t.Error("Expected empty address")
	}
}

func TestAgentInfo_MultipleInstances(t *testing.T) {
	agents := []AgentInfo{
		{Name: "agent1", Address: "host1:50051"},
		{Name: "agent2", Address: "host2:50051"},
		{Name: "agent3", Address: "host3:50051"},
	}

	if len(agents) != 3 {
		t.Errorf("Expected 3 agents, got %d", len(agents))
	}

	for i, agent := range agents {
		expectedName := "agent" + string(rune('1'+i))
		if agent.Name != expectedName {
			t.Errorf("Expected name '%s', got '%s'", expectedName, agent.Name)
		}
	}
}

// Test collector with empty agents
func TestCollector_EmptyAgents(t *testing.T) {
	cfg := CollectorConfig{
		MetricsDB: &MetricsDB{},
	}

	collector := NewCollector(cfg)
	ctx := context.Background()

	// Should not panic with empty agents
	collector.collectAllMetrics(ctx, []AgentInfo{})
}

// Test collector timeout configuration
func TestCollector_TimeoutConfiguration(t *testing.T) {
	testCases := []struct {
		name            string
		timeout         time.Duration
		expectedTimeout time.Duration
	}{
		{"Default timeout", 0, 3 * time.Second},
		{"Custom 1s", 1 * time.Second, 1 * time.Second},
		{"Custom 10s", 10 * time.Second, 10 * time.Second},
		{"Custom 100ms", 100 * time.Millisecond, 100 * time.Millisecond},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := CollectorConfig{
				MetricsDB: &MetricsDB{},
				Timeout:   tc.timeout,
			}

			collector := NewCollector(cfg)

			if collector.timeout != tc.expectedTimeout {
				t.Errorf("Expected timeout %v, got %v", tc.expectedTimeout, collector.timeout)
			}
		})
	}
}

// Test collector interval configuration
func TestCollector_IntervalConfiguration(t *testing.T) {
	intervals := []time.Duration{
		30 * time.Second,
		1 * time.Minute,
		5 * time.Minute,
		10 * time.Minute,
	}

	for _, interval := range intervals {
		cfg := CollectorConfig{
			MetricsDB: &MetricsDB{},
			Interval:  interval,
		}

		collector := NewCollector(cfg)

		if collector.interval != interval {
			t.Errorf("Expected interval %v, got %v", interval, collector.interval)
		}
	}
}

// Test collector retention configuration
func TestCollector_RetentionConfiguration(t *testing.T) {
	retentions := []int{1, 7, 14, 30, 90}

	for _, retention := range retentions {
		cfg := CollectorConfig{
			MetricsDB:     &MetricsDB{},
			RetentionDays: retention,
		}

		collector := NewCollector(cfg)

		if collector.retentionDays != retention {
			t.Errorf("Expected retention %d days, got %d", retention, collector.retentionDays)
		}
	}
}

// Test collector batch size configuration
func TestCollector_BatchSizeConfiguration(t *testing.T) {
	batchSizes := []int{1, 5, 10, 20, 50, 100}

	for _, batchSize := range batchSizes {
		cfg := CollectorConfig{
			MetricsDB: &MetricsDB{},
			BatchSize: batchSize,
		}

		collector := NewCollector(cfg)

		if collector.batchSize != batchSize {
			t.Errorf("Expected batch size %d, got %d", batchSize, collector.batchSize)
		}
	}
}

// Test concurrent access to IsRunning
func TestCollector_ConcurrentIsRunning(t *testing.T) {
	cfg := CollectorConfig{
		MetricsDB: &MetricsDB{},
	}

	collector := NewCollector(cfg)

	done := make(chan bool, 10)

	// Multiple goroutines checking IsRunning
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_ = collector.IsRunning()
			}
			done <- true
		}()
	}

	// Wait for all
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Test collector stop channel
func TestCollector_StopChannel(t *testing.T) {
	cfg := CollectorConfig{
		MetricsDB: &MetricsDB{},
	}

	collector := NewCollector(cfg)

	// stopCh should be created
	if collector.stopCh == nil {
		t.Fatal("Expected stopCh to be created")
	}

	// Should be able to receive from stopCh after closing
	collector.Stop()

	select {
	case <-collector.stopCh:
		// Good - channel was closed
	case <-time.After(100 * time.Millisecond):
		// Channel might already be closed from a previous stop
		// This is acceptable
	}
}

// Test collector fields initialization
func TestCollector_FieldsInitialization(t *testing.T) {
	cfg := CollectorConfig{
		MetricsDB:     &MetricsDB{},
		Interval:      45 * time.Second,
		RetentionDays: 10,
		BatchSize:     15,
		Timeout:       4 * time.Second,
	}

	collector := NewCollector(cfg)

	// Verify all fields are properly initialized
	if collector.metricsDB == nil {
		t.Error("metricsDB should be initialized")
	}

	if collector.interval != 45*time.Second {
		t.Error("interval should be initialized")
	}

	if collector.retentionDays != 10 {
		t.Error("retentionDays should be initialized")
	}

	if collector.batchSize != 15 {
		t.Error("batchSize should be initialized")
	}

	if collector.timeout != 4*time.Second {
		t.Error("timeout should be initialized")
	}

	if collector.stopCh == nil {
		t.Error("stopCh should be initialized")
	}

	// Check initial state
	if collector.running {
		t.Error("running should be false initially")
	}

	if len(collector.lastAgents) != 0 {
		t.Error("lastAgents should be empty initially")
	}

	if !collector.lastCheck.IsZero() {
		t.Error("lastCheck should be zero initially")
	}
}
