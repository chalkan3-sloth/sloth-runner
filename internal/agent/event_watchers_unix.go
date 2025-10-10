//go:build unix || linux || darwin || freebsd || openbsd || netbsd

package agent

import (
	"syscall"
)

// checkDisk checks for disk space threshold events (Unix/Linux/macOS implementation)
func (w *Watcher) checkDisk(eventWorker *EventWorker) {
	// Use syscall to get disk stats
	var stat syscall.Statfs_t
	err := syscall.Statfs(w.config.FilePath, &stat)
	if err != nil {
		return
	}

	totalSpace := stat.Blocks * uint64(stat.Bsize)
	freeSpace := stat.Bfree * uint64(stat.Bsize)
	usedPercent := float64(totalSpace-freeSpace) / float64(totalSpace) * 100.0

	if w.hasCondition(ConditionAbove) && usedPercent > w.config.Threshold {
		eventWorker.SendEvent("disk.high_usage", w.config.Stack, w.config.RunID, map[string]interface{}{
			"path":         w.config.FilePath,
			"used_percent": usedPercent,
			"threshold":    w.config.Threshold,
			"watcher_id":   w.config.ID,
		})
	}
}
