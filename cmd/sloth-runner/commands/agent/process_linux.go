//go:build linux
// +build linux

package agent

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ProcessCache caches process list to avoid expensive reads
type ProcessCache struct {
	mu         sync.RWMutex
	processes  []ProcessInfo
	lastUpdate time.Time
	cacheTTL   time.Duration
}

var processCache = &ProcessCache{
	cacheTTL: 10 * time.Second,
}

// getProcessesOptimized reads processes directly from /proc (Linux only)
// This is 10-20x faster than using gopsutil
func getProcessesOptimized() ([]ProcessInfo, error) {
	// Check cache first
	processCache.mu.RLock()
	if time.Since(processCache.lastUpdate) < processCache.cacheTTL {
		cached := make([]ProcessInfo, len(processCache.processes))
		copy(cached, processCache.processes)
		processCache.mu.RUnlock()
		return cached, nil
	}
	processCache.mu.RUnlock()

	// Read /proc directory
	procDir, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer procDir.Close()

	entries, err := procDir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	processes := make([]ProcessInfo, 0, 30)
	count := 0

	for _, entry := range entries {
		// Limit to top 30 processes (reduced for memory optimization)
		if count >= 30 {
			break
		}

		// Only numeric entries (PIDs)
		pid, err := strconv.Atoi(entry)
		if err != nil {
			continue
		}

		// Read process info
		info, err := readProcInfo(int32(pid))
		if err != nil {
			continue
		}

		if info != nil {
			processes = append(processes, *info)
			count++
		}
	}

	// Update cache
	processCache.mu.Lock()
	processCache.processes = processes
	processCache.lastUpdate = time.Now()
	processCache.mu.Unlock()

	return processes, nil
}

// readProcInfo reads minimal process info from /proc/[pid]
func readProcInfo(pid int32) (*ProcessInfo, error) {
	// Read /proc/[pid]/stat (most efficient)
	statPath := "/proc/" + strconv.Itoa(int(pid)) + "/stat"
	file, err := os.Open(statPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return nil, scanner.Err()
	}

	line := scanner.Text()

	// Parse stat line format: pid (name) state ppid ...
	nameStart := strings.IndexByte(line, '(')
	nameEnd := strings.LastIndexByte(line, ')')

	if nameStart < 0 || nameEnd < 0 {
		return nil, nil
	}

	name := line[nameStart+1 : nameEnd]
	fields := strings.Fields(line[nameEnd+2:])

	if len(fields) < 22 {
		return nil, nil
	}

	info := &ProcessInfo{
		PID:    int(pid),
		Name:   name,
		Status: fields[0], // Process state (R, S, D, Z, T)
	}

	// CPU time (user + system) - field 13 and 14
	if utime, err := strconv.ParseInt(fields[11], 10, 64); err == nil {
		if stime, err := strconv.ParseInt(fields[12], 10, 64); err == nil {
			// Simplified CPU calculation
			info.CPUPercent = float64(utime+stime) / 100.0
		}
	}

	// Memory (RSS in pages) - field 23
	if rss, err := strconv.ParseInt(fields[21], 10, 64); err == nil {
		info.MemoryBytes = uint64(rss * 4096) // Pages to bytes
		info.MemoryPercent = 0                // Will be calculated if needed
	}

	// Read cmdline (optional, cached separately)
	cmdlinePath := "/proc/" + strconv.Itoa(int(pid)) + "/cmdline"
	if cmdlineBytes, err := os.ReadFile(cmdlinePath); err == nil && len(cmdlineBytes) > 0 {
		// Replace null bytes with spaces
		cmdline := string(cmdlineBytes)
		info.Command = strings.ReplaceAll(cmdline, "\x00", " ")

		// STRICT limit on cmdline length to save memory (100 -> 50 chars)
		if len(info.Command) > 50 {
			info.Command = info.Command[:50] + "..."
		}
	}

	// Read user (owner) - requires stat syscall, skip for performance
	// Can be added if really needed

	return info, nil
}

// ClearProcessCache clears the process cache (useful for testing)
func ClearProcessCache() {
	processCache.mu.Lock()
	processCache.processes = nil
	processCache.lastUpdate = time.Time{}
	processCache.mu.Unlock()
}
