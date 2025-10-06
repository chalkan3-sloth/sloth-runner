package agent

import (
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    uint64
		expected string
	}{
		{
			name:     "zero bytes",
			bytes:    0,
			expected: "0 B",
		},
		{
			name:     "less than 1 KiB",
			bytes:    512,
			expected: "512 B",
		},
		{
			name:     "exactly 1 KiB",
			bytes:    1024,
			expected: "1.0 KiB",
		},
		{
			name:     "1.5 KiB",
			bytes:    1536,
			expected: "1.5 KiB",
		},
		{
			name:     "exactly 1 MiB",
			bytes:    1024 * 1024,
			expected: "1.0 MiB",
		},
		{
			name:     "2.5 MiB",
			bytes:    2621440,
			expected: "2.5 MiB",
		},
		{
			name:     "exactly 1 GiB",
			bytes:    1024 * 1024 * 1024,
			expected: "1.0 GiB",
		},
		{
			name:     "4.7 GiB",
			bytes:    5046586573,
			expected: "4.7 GiB",
		},
		{
			name:     "exactly 1 TiB",
			bytes:    1024 * 1024 * 1024 * 1024,
			expected: "1.0 TiB",
		},
		{
			name:     "exactly 1 PiB",
			bytes:    1024 * 1024 * 1024 * 1024 * 1024,
			expected: "1.0 PiB",
		},
		{
			name:     "exactly 1 EiB",
			bytes:    1024 * 1024 * 1024 * 1024 * 1024 * 1024,
			expected: "1.0 EiB",
		},
		{
			name:     "large value",
			bytes:    9223372036854775807, // Max int64
			expected: "8.0 EiB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("formatBytes(%d) = %s; want %s", tt.bytes, result, tt.expected)
			}
		})
	}
}

func BenchmarkFormatBytes(b *testing.B) {
	testCases := []uint64{
		0,
		1024,
		1024 * 1024,
		1024 * 1024 * 1024,
		1024 * 1024 * 1024 * 1024,
	}

	for _, size := range testCases {
		b.Run(formatBytes(size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				formatBytes(size)
			}
		})
	}
}
