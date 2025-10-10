//go:build windows

package resources

import (
	"fmt"
)

// GetDiskUsage retorna uso de disco para um path específico (Windows stub implementation)
func GetDiskUsage(path string) (*DiskStats, error) {
	// Windows disk usage is not yet implemented
	// TODO: Implement using Windows-specific APIs (GetDiskFreeSpaceEx)
	return nil, fmt.Errorf("disk usage monitoring not yet implemented on Windows")
}
