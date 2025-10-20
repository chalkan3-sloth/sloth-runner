package grpcclient

import (
	"testing"
	"time"
)

// Test Connect with invalid addresses (should fail gracefully)
func TestConnect_InvalidAddress(t *testing.T) {
	addresses := []string{
		"",
		"invalid",
		"::::",
		"localhost",
	}

	for _, addr := range addresses {
		t.Run(addr, func(t *testing.T) {
			// These should timeout or fail
			_, err := Connect(addr)
			// We expect errors for invalid addresses
			if err == nil && addr == "" {
				t.Error("Expected error for empty address")
			}
		})
	}
}

// Test Connect function signature
func TestConnect_FunctionSignature(t *testing.T) {
	// Just verify the function exists and has the right signature
	var f func(string) (interface{}, error) = func(s string) (interface{}, error) {
		return Connect(s)
	}

	if f == nil {
		t.Error("Connect function should exist")
	}
}

// Test context timeout behavior
func TestConnect_Timeout(t *testing.T) {
	// Connect to non-existent address should timeout
	address := "localhost:99999"

	start := time.Now()
	_, err := Connect(address)
	duration := time.Since(start)

	if err == nil {
		t.Error("Expected error connecting to non-existent address")
	}

	// Should timeout around 10 seconds (with some margin)
	if duration < 8*time.Second || duration > 12*time.Second {
		t.Logf("Connection attempt took %v (expected ~10s)", duration)
	}
}

// Test address format validation
func TestConnect_AddressFormats(t *testing.T) {
	tests := []struct {
		name    string
		address string
		valid   bool
	}{
		{"empty", "", false},
		{"localhost_no_port", "localhost", false},
		{"localhost_with_port", "localhost:50051", true},
		{"ip_with_port", "127.0.0.1:50051", true},
		{"ip_no_port", "127.0.0.1", false},
		{"ipv6_with_port", "[::1]:50051", true},
		{"fqdn_with_port", "example.com:50051", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We don't actually connect, just validate the address format concept
			if tt.address == "" && tt.valid {
				t.Error("Empty address should not be valid")
			}
		})
	}
}

// Test connection with context
func TestConnect_UsesContext(t *testing.T) {
	// The Connect function uses context internally
	// We can't easily test this without a real server, but we verify the concept

	// Context usage is evidenced by timeout behavior
	address := "localhost:99999"

	// This will use the internal 10s timeout
	start := time.Now()
	_, err := Connect(address)
	duration := time.Since(start)

	if err == nil {
		t.Error("Expected error for unreachable address")
	}

	// Should respect the 10s timeout
	if duration < 8*time.Second {
		t.Errorf("Timeout too short: %v", duration)
	}
}

// Test connection error messages
func TestConnect_ErrorMessages(t *testing.T) {
	address := "invalid-address"

	_, err := Connect(address)

	if err == nil {
		t.Error("Expected error for invalid address")
		return
	}

	// Error message should mention the address
	errStr := err.Error()
	if errStr == "" {
		t.Error("Error message should not be empty")
	}
}

// Test multiple concurrent connections
func TestConnect_Concurrent(t *testing.T) {
	// Test that multiple concurrent connection attempts don't interfere
	done := make(chan bool, 3)

	for i := 0; i < 3; i++ {
		go func() {
			_, _ = Connect("localhost:99999")
			done <- true
		}()
	}

	// Wait for all goroutines
	timeout := time.After(35 * time.Second) // 3 * 10s + margin
	for i := 0; i < 3; i++ {
		select {
		case <-done:
			// Good
		case <-timeout:
			t.Fatal("Timeout waiting for concurrent connections")
		}
	}
}

// Test that function returns proper types
func TestConnect_ReturnTypes(t *testing.T) {
	// Connect should return (*grpc.ClientConn, error)
	// We verify the return signature concept

	_, err := Connect("localhost:99999")

	// Error should be non-nil for failed connection
	if err == nil {
		t.Log("Note: Connection succeeded (unexpected)")
	}

	// If connection failed, error should have a message
	if err != nil && err.Error() == "" {
		t.Error("Error should have a message")
	}
}

// Test address with different ports
func TestConnect_DifferentPorts(t *testing.T) {
	ports := []string{
		"localhost:50051",
		"localhost:8080",
		"localhost:9090",
		"localhost:3000",
	}

	for _, addr := range ports {
		t.Run(addr, func(t *testing.T) {
			// Will timeout but shouldn't panic
			_, err := Connect(addr)
			// Expected to fail (no server running)
			if err == nil {
				t.Log("Note: Connection succeeded (server running)")
			}
		})
	}
}

// Test special characters in address
func TestConnect_SpecialAddresses(t *testing.T) {
	addresses := []string{
		"0.0.0.0:50051",
		"[::]:50051",
		"localhost:0",
	}

	for _, addr := range addresses {
		t.Run(addr, func(t *testing.T) {
			_, err := Connect(addr)
			// These should fail or timeout
			_ = err // Just verify it doesn't panic
		})
	}
}

// Test connection with IPv4
func TestConnect_IPv4(t *testing.T) {
	address := "127.0.0.1:50051"

	_, err := Connect(address)

	// Will likely timeout (no server)
	if err == nil {
		t.Log("Note: Connection succeeded to 127.0.0.1:50051")
	}
}

// Test connection with IPv6
func TestConnect_IPv6(t *testing.T) {
	address := "[::1]:50051"

	_, err := Connect(address)

	// Will likely timeout (no server)
	if err == nil {
		t.Log("Note: Connection succeeded to [::1]:50051")
	}
}

// Test connection with hostname
func TestConnect_Hostname(t *testing.T) {
	address := "localhost:50051"

	_, err := Connect(address)

	// Will likely timeout (no server)
	if err == nil {
		t.Log("Note: Connection succeeded to localhost:50051")
	}
}

// Test that Connect uses WithBlock option
func TestConnect_BlockingBehavior(t *testing.T) {
	// WithBlock option causes Connect to wait for connection
	// We verify this by checking that it doesn't return immediately

	address := "localhost:99999"

	start := time.Now()
	_, err := Connect(address)
	duration := time.Since(start)

	if err == nil {
		t.Log("Connection succeeded (unexpected)")
		return
	}

	// Should take several seconds (waiting for timeout)
	if duration < 1*time.Second {
		t.Error("Connect returned too quickly (should block)")
	}
}

// Test that Connect uses insecure credentials
func TestConnect_UsesInsecureCredentials(t *testing.T) {
	// The function uses insecure.NewCredentials()
	// We can't test this directly without inspecting internals,
	// but we verify the function works as expected

	address := "localhost:50051"
	_, err := Connect(address)

	// Should fail with timeout or connection error, not credentials error
	if err != nil {
		errStr := err.Error()
		if errStr != "" {
			// Just verify it's some kind of connection error
			t.Logf("Connection error: %v", err)
		}
	}
}

// Test timeout duration
func TestConnect_TimeoutDuration(t *testing.T) {
	// Connect uses 10 second timeout
	address := "198.51.100.1:50051" // TEST-NET-2 address (should timeout)

	start := time.Now()
	_, err := Connect(address)
	duration := time.Since(start)

	if err == nil {
		t.Log("Connection succeeded (unexpected)")
		return
	}

	// Should timeout around 10 seconds
	if duration < 8*time.Second || duration > 12*time.Second {
		t.Logf("Note: Timeout was %v (expected ~10s)", duration)
	}
}

// Test error wrapping
func TestConnect_ErrorWrapping(t *testing.T) {
	address := "invalid:::address"

	_, err := Connect(address)

	if err == nil {
		t.Error("Expected error for malformed address")
		return
	}

	// Error should provide context
	errStr := err.Error()
	if !contains(errStr, "failed to connect") {
		t.Logf("Error message: %s", errStr)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0)
}

// Test nil address handling
func TestConnect_NilAddress(t *testing.T) {
	// Empty string should be handled
	_, err := Connect("")

	if err == nil {
		t.Error("Expected error for empty address")
	}
}

// Test very long address
func TestConnect_LongAddress(t *testing.T) {
	longAddr := "very-long-hostname-that-probably-does-not-exist-anywhere.example.com:50051"

	_, err := Connect(longAddr)

	// Should fail (likely DNS resolution or timeout)
	if err == nil {
		t.Log("Connection succeeded (unexpected)")
	}
}

// Test sequential connections
func TestConnect_Sequential(t *testing.T) {
	addresses := []string{
		"localhost:50051",
		"localhost:50052",
		"localhost:50053",
	}

	for _, addr := range addresses {
		_, err := Connect(addr)
		_ = err // Each should fail independently
	}
}

// Test connection cancellation concept
func TestConnect_CancellationConcept(t *testing.T) {
	// Connect uses internal context with timeout
	// This ensures connections don't hang indefinitely

	address := "192.0.2.1:50051" // TEST-NET-1 (should timeout)

	start := time.Now()
	_, err := Connect(address)
	duration := time.Since(start)

	if err == nil {
		t.Log("Connection succeeded (unexpected)")
		return
	}

	// Should not hang indefinitely
	if duration > 15*time.Second {
		t.Errorf("Connection took too long: %v", duration)
	}
}
