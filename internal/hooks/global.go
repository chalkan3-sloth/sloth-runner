package hooks

import (
	"log/slog"
	"sync"
)

var (
	globalDispatcher *Dispatcher
	globalRepo       *Repository
	globalMu         sync.RWMutex
)

// InitializeGlobalDispatcher initializes the global hook dispatcher
// This should be called at application startup
func InitializeGlobalDispatcher() error {
	globalMu.Lock()
	defer globalMu.Unlock()

	if globalDispatcher != nil {
		return nil // Already initialized
	}

	repo, err := NewRepository()
	if err != nil {
		return err
	}

	globalRepo = repo
	globalDispatcher = NewDispatcher(repo)

	slog.Info("global hook dispatcher initialized")
	return nil
}

// GetGlobalDispatcher returns the global dispatcher instance
func GetGlobalDispatcher() *Dispatcher {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalDispatcher
}

// CleanupGlobalDispatcher cleans up the global dispatcher
func CleanupGlobalDispatcher() {
	globalMu.Lock()
	defer globalMu.Unlock()

	if globalDispatcher != nil {
		globalDispatcher.StopEventProcessor()
		globalDispatcher = nil
	}

	if globalRepo != nil {
		globalRepo.Close()
		globalRepo = nil
	}

	slog.Info("global hook dispatcher cleaned up")
}
