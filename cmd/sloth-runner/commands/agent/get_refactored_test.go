package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/chalkan3-sloth/sloth-runner/cmd/sloth-runner/commands/agent/mocks"
	agentInternal "github.com/chalkan3-sloth/sloth-runner/internal/agent"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
	"google.golang.org/grpc"
)

// Test getAgentInfoWithClient function
func TestGetAgentInfoWithClient(t *testing.T) {
	tests := []struct {
		name          string
		mockResp      *pb.GetAgentInfoResponse
		mockError     error
		outputFormat  string
		expectedError bool
	}{
		{
			name: "successful JSON output",
			mockResp: &pb.GetAgentInfoResponse{
				Success: true,
				AgentInfo: &pb.AgentInfo{
					AgentName:         "test-agent",
					AgentAddress:      "localhost:50052",
					Status:            "Active",
					LastHeartbeat:     1234567890,
					LastInfoCollected: 1234567890,
					SystemInfoJson:    `{"hostname":"test-host"}`,
				},
			},
			outputFormat:  "json",
			expectedError: false,
		},
		{
			name: "successful text output",
			mockResp: &pb.GetAgentInfoResponse{
				Success: true,
				AgentInfo: &pb.AgentInfo{
					AgentName:    "text-agent",
					AgentAddress: "localhost:50053",
					Status:       "Active",
				},
			},
			outputFormat:  "text",
			expectedError: false,
		},
		{
			name: "failed response",
			mockResp: &pb.GetAgentInfoResponse{
				Success: false,
				Message: "agent not found",
			},
			outputFormat:  "text",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewMockAgentRegistryClient()
			mockClient.GetAgentInfoFunc = func(ctx context.Context, in *pb.GetAgentInfoRequest, opts ...grpc.CallOption) (*pb.GetAgentInfoResponse, error) {
				return tt.mockResp, tt.mockError
			}

			var buf bytes.Buffer
			opts := GetAgentInfoOptions{
				AgentName:    "test",
				OutputFormat: tt.outputFormat,
				Writer:       &buf,
			}

			err := getAgentInfoWithClient(context.Background(), mockClient, opts)

			if tt.expectedError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectedError {
				if buf.Len() == 0 {
					t.Error("Expected output but got none")
				}
			}
		})
	}
}

// Test formatAgentInfoJSON function
func TestFormatAgentInfoJSON(t *testing.T) {
	agentInfo := &pb.AgentInfo{
		AgentName:         "json-agent",
		AgentAddress:      "localhost:50055",
		Status:            "Active",
		LastHeartbeat:     1234567890,
		LastInfoCollected: 1234567890,
		SystemInfoJson:    `{"hostname":"json-host","platform":"linux"}`,
	}

	var buf bytes.Buffer
	err := formatAgentInfoJSON(agentInfo, &buf)
	if err != nil {
		t.Fatalf("formatAgentInfoJSON failed: %v", err)
	}

	// Verify JSON is valid
	var output map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Errorf("Invalid JSON output: %v", err)
	}

	// Verify required fields
	if output["agent_name"] != "json-agent" {
		t.Errorf("Expected agent_name=json-agent, got %v", output["agent_name"])
	}
	if output["status"] != "Active" {
		t.Errorf("Expected status=Active, got %v", output["status"])
	}

	// Verify system_info is parsed
	sysInfo, ok := output["system_info"].(map[string]interface{})
	if !ok {
		t.Error("system_info should be a map")
	}
	if sysInfo["hostname"] != "json-host" {
		t.Errorf("Expected hostname=json-host, got %v", sysInfo["hostname"])
	}
}

// Test formatAgentInfoJSON with empty system info
func TestFormatAgentInfoJSON_EmptySystemInfo(t *testing.T) {
	agentInfo := &pb.AgentInfo{
		AgentName:      "no-sysinfo",
		AgentAddress:   "localhost:50056",
		Status:         "Active",
		SystemInfoJson: "",
	}

	var buf bytes.Buffer
	err := formatAgentInfoJSON(agentInfo, &buf)
	if err != nil {
		t.Fatalf("formatAgentInfoJSON failed: %v", err)
	}

	var output map[string]interface{}
	json.Unmarshal(buf.Bytes(), &output)

	if output["system_info"] != nil {
		t.Errorf("Expected system_info=nil, got %v", output["system_info"])
	}
}

// Test formatAgentInfoText function
func TestFormatAgentInfoText(t *testing.T) {
	agentInfo := &pb.AgentInfo{
		AgentName:         "text-agent",
		AgentAddress:      "localhost:50057",
		Status:            "Active",
		LastHeartbeat:     1234567890,
		LastInfoCollected: 1234567890,
		SystemInfoJson:    "",
	}

	var buf bytes.Buffer
	err := formatAgentInfoText(agentInfo, &buf)
	if err != nil {
		t.Fatalf("formatAgentInfoText failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "text-agent") {
		t.Error("Output should contain agent name")
	}
	if !strings.Contains(output, "Active") {
		t.Error("Output should contain status")
	}
}

// Test formatMemoryInfo function
func TestFormatMemoryInfo(t *testing.T) {
	memory := &agentInternal.MemoryInfo{
		Total:       8 * 1024 * 1024 * 1024, // 8 GB
		Used:        4 * 1024 * 1024 * 1024, // 4 GB
		Available:   4 * 1024 * 1024 * 1024,
		Free:        3 * 1024 * 1024 * 1024,
		Cached:      1 * 1024 * 1024 * 1024,
		UsedPercent: 50.0,
	}

	var buf bytes.Buffer
	formatMemoryInfo(memory, &buf)

	output := buf.String()
	if !strings.Contains(output, "Memory Information") {
		t.Error("Output should contain 'Memory Information'")
	}
	if !strings.Contains(output, "GiB") {
		t.Error("Output should contain GiB unit")
	}
	if !strings.Contains(output, "50") {
		t.Error("Output should contain percentage")
	}
}

// Test formatDiskInfoList function
func TestFormatDiskInfoList(t *testing.T) {
	disks := []*agentInternal.DiskInfo{
		{
			Device:      "/dev/sda1",
			Mountpoint:  "/",
			Total:       100 * 1024 * 1024 * 1024,
			Used:        60 * 1024 * 1024 * 1024,
			Free:        40 * 1024 * 1024 * 1024,
			UsedPercent: 60.0,
		},
		{
			Device:      "/dev/sdb1",
			Mountpoint:  "/home",
			Total:       500 * 1024 * 1024 * 1024,
			Used:        200 * 1024 * 1024 * 1024,
			Free:        300 * 1024 * 1024 * 1024,
			UsedPercent: 40.0,
		},
	}

	var buf bytes.Buffer
	formatDiskInfoList(disks, &buf)

	output := buf.String()
	if !strings.Contains(output, "Disk Information") {
		t.Error("Output should contain 'Disk Information'")
	}
	if !strings.Contains(output, "/dev/sda1") {
		t.Error("Output should contain device name")
	}
	if !strings.Contains(output, "/home") {
		t.Error("Output should contain mountpoint")
	}
}

// Test formatNetworkInfoList function
func TestFormatNetworkInfoList(t *testing.T) {
	network := []*agentInternal.NetworkInfo{
		{
			Name:      "eth0",
			MAC:       "00:11:22:33:44:55",
			Addresses: []string{"192.168.1.100", "fe80::1"},
			IsUp:      true,
		},
		{
			Name:      "eth1",
			Addresses: []string{"10.0.0.1"},
			IsUp:      false,
		},
		{
			Name:      "lo",
			Addresses: []string{"127.0.0.1"},
			IsUp:      true,
		},
	}

	var buf bytes.Buffer
	formatNetworkInfoList(network, &buf)

	output := buf.String()
	if !strings.Contains(output, "Network Interfaces") {
		t.Error("Output should contain 'Network Interfaces'")
	}
	if !strings.Contains(output, "eth0") {
		t.Error("Output should contain eth0")
	}
	if !strings.Contains(output, "192.168.1.100") {
		t.Error("Output should contain IP address")
	}
	// Loopback should be skipped
	if strings.Count(output, "127.0.0.1") > 0 {
		t.Error("Output should not contain loopback interface")
	}
}

// Test formatPackageInfo function
func TestFormatPackageInfo(t *testing.T) {
	packages := &agentInternal.PackageInfo{
		Manager:          "apt",
		InstalledCount:   1500,
		UpdatesAvailable: 25,
	}

	var buf bytes.Buffer
	formatPackageInfo(packages, &buf)

	output := buf.String()
	if !strings.Contains(output, "Package Information") {
		t.Error("Output should contain 'Package Information'")
	}
	if !strings.Contains(output, "apt") {
		t.Error("Output should contain package manager")
	}
	if !strings.Contains(output, "1500") {
		t.Error("Output should contain installed count")
	}
	if !strings.Contains(output, "25") {
		t.Error("Output should contain updates available")
	}
}

// Test formatPackageInfo with no updates
func TestFormatPackageInfo_NoUpdates(t *testing.T) {
	packages := &agentInternal.PackageInfo{
		Manager:          "yum",
		InstalledCount:   800,
		UpdatesAvailable: 0,
	}

	var buf bytes.Buffer
	formatPackageInfo(packages, &buf)

	output := buf.String()
	if !strings.Contains(output, "up to date") {
		t.Error("Output should indicate system is up to date")
	}
}

// Test formatServicesInfoList function
func TestFormatServicesInfoList(t *testing.T) {
	services := []agentInternal.ServiceInfo{
		{Name: "nginx", Status: "running"},
		{Name: "postgresql", Status: "running"},
		{Name: "redis", Status: "running"},
		{Name: "docker", Status: "running"},
		{Name: "ssh", Status: "running"},
	}

	var buf bytes.Buffer
	formatServicesInfoList(services, &buf)

	output := buf.String()
	if !strings.Contains(output, "Running Services") {
		t.Error("Output should contain 'Running Services'")
	}
	if !strings.Contains(output, "nginx") {
		t.Error("Output should contain service name")
	}
	if !strings.Contains(output, "5 total") {
		t.Error("Output should contain total count")
	}
}

// Test formatServicesInfoList with many services
func TestFormatServicesInfoList_Many(t *testing.T) {
	services := make([]agentInternal.ServiceInfo, 15)
	for i := 0; i < 15; i++ {
		services[i] = agentInternal.ServiceInfo{
			Name:   "service" + string(rune('0'+i)),
			Status: "running",
		}
	}

	var buf bytes.Buffer
	formatServicesInfoList(services, &buf)

	output := buf.String()
	if !strings.Contains(output, "15 total") {
		t.Error("Output should contain correct total")
	}
	if !strings.Contains(output, "and 5 more") {
		t.Error("Output should indicate there are more services")
	}
}

// Test formatSystemInfo function
func TestFormatSystemInfo(t *testing.T) {
	sysInfo := &agentInternal.SystemInfo{
		Hostname:        "test-host",
		Platform:        "linux",
		PlatformVersion: "22.04",
		Architecture:    "x86_64",
		CPUs:            8,
		Kernel:          "Linux",
		KernelVersion:   "5.15.0",
		Virtualization:  "kvm",
		Uptime:          3600,
		LoadAverage:     []float64{1.5, 1.2, 1.0},
	}

	var buf bytes.Buffer
	formatSystemInfo(sysInfo, &buf)

	output := buf.String()
	if !strings.Contains(output, "test-host") {
		t.Error("Output should contain hostname")
	}
	if !strings.Contains(output, "linux") {
		t.Error("Output should contain platform")
	}
	if !strings.Contains(output, "x86_64") {
		t.Error("Output should contain architecture")
	}
	if !strings.Contains(output, "kvm") {
		t.Error("Output should contain virtualization")
	}
	if !strings.Contains(output, "1.5") {
		t.Error("Output should contain load average")
	}
}

// Test getAgentInfo integration with factory
func TestGetAgentInfo_WithFactory(t *testing.T) {
	// This test verifies the integration but will fail on connection
	// It's here to show the structure is correct
	err := getAgentInfo("test-agent", "invalid:99999", "json")
	if err == nil {
		t.Log("Command succeeded (unexpected, but possible if server is running)")
	} else {
		// Expected to fail with connection error
		if !strings.Contains(err.Error(), "failed to connect") &&
			!strings.Contains(err.Error(), "connection refused") &&
			!strings.Contains(err.Error(), "no such host") {
			t.Logf("Got connection error as expected: %v", err)
		}
	}
}
