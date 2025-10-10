//go:build windows

package agent

import (
	"log/slog"
)

// checkDisk checks for disk space threshold events (Windows stub implementation)
func (w *Watcher) checkDisk(eventWorker *EventWorker) {
	// Windows disk monitoring is not yet implemented
	// TODO: Implement using Windows-specific APIs (GetDiskFreeSpaceEx, etc.)
	slog.Debug("Disk monitoring not yet implemented on Windows", "watcher_id", w.config.ID)
}
