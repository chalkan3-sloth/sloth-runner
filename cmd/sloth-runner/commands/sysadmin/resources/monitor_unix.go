//go:build unix || linux || darwin || freebsd || openbsd || netbsd

package resources

import (
	"syscall"
)

// GetDiskUsage retorna uso de disco para um path específico (Unix/Linux/macOS implementation)
func GetDiskUsage(path string) (*DiskStats, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return nil, err
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	available := stat.Bavail * uint64(stat.Bsize)
	used := total - free

	usagePercent := 0.0
	if total > 0 {
		usagePercent = 100.0 * float64(used) / float64(total)
	}

	return &DiskStats{
		Filesystem:   "direct",
		MountPoint:   path,
		Total:        total,
		Used:         used,
		Available:    available,
		UsagePercent: usagePercent,
	}, nil
}
