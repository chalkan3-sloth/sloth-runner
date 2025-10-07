package pooling

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// ConnectionPool manages reusable gRPC connections to agents
type ConnectionPool struct {
	mu          sync.RWMutex
	connections map[string]*pooledConnection
	maxIdle     time.Duration
	maxAge      time.Duration
}

type pooledConnection struct {
	conn         *grpc.ClientConn
	lastUsed     time.Time
	createdAt    time.Time
	useCount     int64
	mu           sync.Mutex
}

// Global connection pool
var GlobalConnPool = NewConnectionPool()

// NewConnectionPool creates a new connection pool
func NewConnectionPool() *ConnectionPool {
	pool := &ConnectionPool{
		connections: make(map[string]*pooledConnection),
		maxIdle:     30 * time.Minute, // Close idle connections after 30min
		maxAge:      2 * time.Hour,     // Recycle connections after 2h
	}

	// Start cleanup goroutine
	go pool.cleanupLoop()

	return pool
}

// GetConnection gets or creates a connection to an agent
func (p *ConnectionPool) GetConnection(ctx context.Context, address string) (*grpc.ClientConn, error) {
	p.mu.RLock()
	pc, exists := p.connections[address]
	p.mu.RUnlock()

	// Check if existing connection is still good
	if exists {
		pc.mu.Lock()
		defer pc.mu.Unlock()

		// Check if connection is ready
		state := pc.conn.GetState()
		if state == connectivity.Ready || state == connectivity.Idle {
			// Check age limits
			if time.Since(pc.createdAt) < p.maxAge {
				pc.lastUsed = time.Now()
				pc.useCount++
				return pc.conn, nil
			}
		}

		// Connection is bad or too old, close it
		pc.conn.Close()
		p.mu.Lock()
		delete(p.connections, address)
		p.mu.Unlock()
	}

	// Create new connection
	dialCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		dialCtx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		// Connection pool settings - REDUCED for memory optimization
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1*1024*1024), // 1MB (was 4MB)
			grpc.MaxCallSendMsgSize(1*1024*1024), // 1MB (was 4MB)
		),
	)
	if err != nil {
		return nil, err
	}

	// Store in pool
	pc = &pooledConnection{
		conn:      conn,
		lastUsed:  time.Now(),
		createdAt: time.Now(),
		useCount:  1,
	}

	p.mu.Lock()
	p.connections[address] = pc
	p.mu.Unlock()

	return conn, nil
}

// CloseConnection explicitly closes a connection
func (p *ConnectionPool) CloseConnection(address string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if pc, exists := p.connections[address]; exists {
		pc.conn.Close()
		delete(p.connections, address)
	}
}

// CloseAll closes all connections in the pool
func (p *ConnectionPool) CloseAll() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for address, pc := range p.connections {
		pc.conn.Close()
		delete(p.connections, address)
	}
}

// GetStats returns pool statistics
func (p *ConnectionPool) GetStats() map[string]interface{} {
	p.mu.RLock()
	defer p.mu.RUnlock()

	stats := map[string]interface{}{
		"total_connections": len(p.connections),
		"connections":       make([]map[string]interface{}, 0, len(p.connections)),
	}

	for addr, pc := range p.connections {
		pc.mu.Lock()
		connStats := map[string]interface{}{
			"address":    addr,
			"state":      pc.conn.GetState().String(),
			"age":        time.Since(pc.createdAt).String(),
			"idle_time":  time.Since(pc.lastUsed).String(),
			"use_count":  pc.useCount,
		}
		pc.mu.Unlock()

		stats["connections"] = append(stats["connections"].([]map[string]interface{}), connStats)
	}

	return stats
}

// cleanupLoop periodically cleans up idle connections
func (p *ConnectionPool) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		p.cleanup()
	}
}

// cleanup removes idle or old connections
func (p *ConnectionPool) cleanup() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	for address, pc := range p.connections {
		pc.mu.Lock()

		shouldClose := false
		reason := ""

		// Check idle time
		if now.Sub(pc.lastUsed) > p.maxIdle {
			shouldClose = true
			reason = "idle timeout"
		}

		// Check age
		if now.Sub(pc.createdAt) > p.maxAge {
			shouldClose = true
			reason = "max age exceeded"
		}

		// Check connection state
		state := pc.conn.GetState()
		if state == connectivity.Shutdown || state == connectivity.TransientFailure {
			shouldClose = true
			reason = "bad state: " + state.String()
		}

		if shouldClose {
			pc.conn.Close()
			delete(p.connections, address)
			// Log cleanup (could use slog here)
			_ = reason // Use for logging if needed
		}

		pc.mu.Unlock()
	}
}
