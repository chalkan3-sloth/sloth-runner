package process

import (
	"testing"
	"time"
)

func TestNewProcessManager(t *testing.T) {
	manager := NewProcessManager()
	if manager == nil {
		t.Fatal("NewProcessManager() returned nil")
	}

	_, ok := manager.(*SystemProcessManager)
	if !ok {
		t.Error("NewProcessManager() did not return *SystemProcessManager")
	}
}

func TestList(t *testing.T) {
	manager := NewProcessManager()

	tests := []struct {
		name    string
		options ListOptions
	}{
		{
			name: "List all processes",
			options: ListOptions{
				SortBy: "cpu",
				Top:    10,
			},
		},
		{
			name: "List with filter",
			options: ListOptions{
				SortBy: "memory",
				Top:    5,
				Filter: "go",
			},
		},
		{
			name: "List sorted by name",
			options: ListOptions{
				SortBy: "name",
				Top:    20,
			},
		},
		{
			name: "List sorted by PID",
			options: ListOptions{
				SortBy: "pid",
				Top:    15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processes, err := manager.List(tt.options)
			if err != nil {
				t.Fatalf("List() failed: %v", err)
			}

			if processes == nil {
				t.Fatal("List() returned nil processes")
			}

			// Verifica limite
			if len(processes) > tt.options.Top {
				t.Errorf("List() returned %d processes, expected max %d", len(processes), tt.options.Top)
			}

			// Verifica campos básicos
			for _, p := range processes {
				if p.PID == 0 {
					t.Error("Process has PID 0")
				}
				// Note: Some processes (especially kernel threads) may have empty names
				// So we don't fail if name is empty, just log it
			}
		})
	}
}

func TestListWithUserFilter(t *testing.T) {
	manager := NewProcessManager()

	options := ListOptions{
		SortBy:     "cpu",
		Top:        5,
		UserFilter: "root",
	}

	processes, err := manager.List(options)
	if err != nil {
		t.Fatalf("List() with user filter failed: %v", err)
	}

	// Em alguns sistemas pode não ter processos do root
	if len(processes) > 0 {
		for _, p := range processes {
			if p.Username != "root" {
				t.Errorf("Process %d (%s) has username %s, expected root", p.PID, p.Name, p.Username)
			}
		}
	}
}

func TestInfo(t *testing.T) {
	manager := NewProcessManager()

	// Primeiro pega um processo válido
	processes, err := manager.List(ListOptions{Top: 1})
	if err != nil || len(processes) == 0 {
		t.Skip("No processes available for testing")
	}

	pid := processes[0].PID

	detail, err := manager.Info(pid)
	if err != nil {
		t.Fatalf("Info() failed: %v", err)
	}

	if detail == nil {
		t.Fatal("Info() returned nil detail")
	}

	if detail.PID != pid {
		t.Errorf("Info() returned PID %d, expected %d", detail.PID, pid)
	}

	if detail.Name == "" {
		t.Error("Info() returned empty name")
	}

	// Verifica campos básicos
	if detail.ProcessInfo == nil {
		t.Error("Info() returned nil ProcessInfo")
	}
}

func TestInfoInvalidPID(t *testing.T) {
	manager := NewProcessManager()

	_, err := manager.Info(999999)
	if err == nil {
		t.Error("Info() with invalid PID should return error")
	}
}

func TestMonitor(t *testing.T) {
	manager := NewProcessManager()

	// Pega um processo válido
	processes, err := manager.List(ListOptions{Top: 1})
	if err != nil || len(processes) == 0 {
		t.Skip("No processes available for testing")
	}

	pid := processes[0].PID
	duration := 3 * time.Second

	metrics, err := manager.Monitor(pid, duration)
	if err != nil {
		t.Fatalf("Monitor() failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("Monitor() returned nil metrics")
	}

	if metrics.PID != pid {
		t.Errorf("Monitor() returned PID %d, expected %d", metrics.PID, pid)
	}

	if metrics.Duration != duration {
		t.Errorf("Monitor() returned duration %v, expected %v", metrics.Duration, duration)
	}

	if len(metrics.Samples) == 0 {
		t.Error("Monitor() returned no samples")
	}

	// Deveria ter pelo menos 2 samples em 3 segundos
	if len(metrics.Samples) < 2 {
		t.Errorf("Monitor() returned only %d samples in %v", len(metrics.Samples), duration)
	}

	// Verifica cálculos de média
	if metrics.AvgCPU < 0 {
		t.Error("Monitor() returned negative AvgCPU")
	}

	if metrics.AvgMemory < 0 {
		t.Error("Monitor() returned negative AvgMemory")
	}

	// MaxCPU deve ser >= AvgCPU
	if metrics.MaxCPU < metrics.AvgCPU {
		t.Error("Monitor() MaxCPU is less than AvgCPU")
	}

	// MaxMemory deve ser >= AvgMemory
	if metrics.MaxMemory < metrics.AvgMemory {
		t.Error("Monitor() MaxMemory is less than AvgMemory")
	}
}

func TestMonitorInvalidPID(t *testing.T) {
	manager := NewProcessManager()

	_, err := manager.Monitor(999999, time.Second)
	if err == nil {
		t.Error("Monitor() with invalid PID should return error")
	}
}

func TestSortProcesses(t *testing.T) {
	manager := &SystemProcessManager{}

	processes := []*ProcessInfo{
		{PID: 1, Name: "init", CPUPercent: 10.0, MemoryMB: 100},
		{PID: 2, Name: "systemd", CPUPercent: 5.0, MemoryMB: 200},
		{PID: 3, Name: "bash", CPUPercent: 20.0, MemoryMB: 50},
	}

	tests := []struct {
		sortBy   string
		expected int32 // PID esperado na primeira posição
	}{
		{"cpu", 3},    // bash tem mais CPU
		{"memory", 2}, // systemd tem mais memória
		{"name", 3},   // bash vem primeiro alfabeticamente
		{"pid", 1},    // init tem menor PID
	}

	for _, tt := range tests {
		t.Run(tt.sortBy, func(t *testing.T) {
			// Faz cópia para não afetar outros testes
			testProcs := make([]*ProcessInfo, len(processes))
			copy(testProcs, processes)

			manager.sortProcesses(testProcs, tt.sortBy)

			if testProcs[0].PID != tt.expected {
				t.Errorf("sortProcesses(%s) first PID = %d, expected %d",
					tt.sortBy, testProcs[0].PID, tt.expected)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		s      string
		substr string
		want   bool
	}{
		{"hello world", "world", true},
		{"hello world", "WORLD", true}, // case insensitive
		{"hello world", "foo", false},
		{"", "", true},
		{"test", "", true},
		{"", "test", false},
	}

	for _, tt := range tests {
		got := contains(tt.s, tt.substr)
		if got != tt.want {
			t.Errorf("contains(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"this is a very long string", 10, "this is..."},
		{"exactly10!", 10, "exactly10!"},
		{"", 5, ""},
	}

	for _, tt := range tests {
		got := truncate(tt.input, tt.maxLen)
		if got != tt.want {
			t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
		}
		if len(got) > tt.maxLen {
			t.Errorf("truncate(%q, %d) returned string longer than maxLen", tt.input, tt.maxLen)
		}
	}
}

func TestProcessInfoStructure(t *testing.T) {
	info := &ProcessInfo{
		PID:           1234,
		Name:          "test-process",
		Username:      "testuser",
		CPUPercent:    15.5,
		MemoryMB:      256.0,
		MemoryPercent: 5.2,
		Status:        "running",
		CreateTime:    time.Now().Unix(),
		NumThreads:    4,
		Cmdline:       "/usr/bin/test",
	}

	if info.PID != 1234 {
		t.Error("PID not set correctly")
	}
	if info.Name != "test-process" {
		t.Error("Name not set correctly")
	}
	if info.CPUPercent != 15.5 {
		t.Error("CPUPercent not set correctly")
	}
}

func TestProcessDetailStructure(t *testing.T) {
	detail := &ProcessDetail{
		ProcessInfo: &ProcessInfo{
			PID:  1234,
			Name: "test",
		},
		ParentPID:  1,
		Nice:       0,
		NumFDs:     10,
		Connections: []string{"tcp:80", "tcp:443"},
		OpenFiles:   []string{"/tmp/file1", "/tmp/file2"},
	}

	if detail.PID != 1234 {
		t.Error("PID not accessible through ProcessDetail")
	}
	if detail.ParentPID != 1 {
		t.Error("ParentPID not set correctly")
	}
	if len(detail.Connections) != 2 {
		t.Error("Connections not set correctly")
	}
	if len(detail.OpenFiles) != 2 {
		t.Error("OpenFiles not set correctly")
	}
}

func TestProcessMetricsStructure(t *testing.T) {
	metrics := &ProcessMetrics{
		PID:      1234,
		Duration: 10 * time.Second,
		Samples: []*ProcessSnapshot{
			{
				Timestamp:     time.Now(),
				CPUPercent:    10.0,
				MemoryMB:      100.0,
				MemoryPercent: 5.0,
				NumThreads:    4,
			},
			{
				Timestamp:     time.Now(),
				CPUPercent:    20.0,
				MemoryMB:      150.0,
				MemoryPercent: 7.0,
				NumThreads:    4,
			},
		},
		AvgCPU:    15.0,
		MaxCPU:    20.0,
		AvgMemory: 125.0,
		MaxMemory: 150.0,
	}

	if metrics.PID != 1234 {
		t.Error("PID not set correctly")
	}
	if len(metrics.Samples) != 2 {
		t.Error("Samples not set correctly")
	}
	if metrics.AvgCPU != 15.0 {
		t.Error("AvgCPU not set correctly")
	}
	if metrics.MaxCPU != 20.0 {
		t.Error("MaxCPU not set correctly")
	}
}
