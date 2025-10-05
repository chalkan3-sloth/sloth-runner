package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server manages the Prometheus metrics HTTP server
type Server struct {
	httpServer *http.Server
	metrics    *Metrics
	registry   *prometheus.Registry
	port       int
	enabled    bool
}

// NewServer creates a new telemetry server
func NewServer(port int, enabled bool) *Server {
	if port == 0 {
		port = 9090 // Default Prometheus port
	}

	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// Add Go collector for runtime metrics
	registry.MustRegister(prometheus.NewGoCollector())
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	return &Server{
		metrics:  metrics,
		registry: registry,
		port:     port,
		enabled:  enabled,
	}
}

// Start starts the metrics HTTP server
func (s *Server) Start() error {
	if !s.enabled {
		slog.Info("Telemetry server is disabled")
		return nil
	}

	mux := http.NewServeMux()

	// Prometheus metrics endpoint
	mux.Handle("/metrics", promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	}))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Info endpoint
	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"service":"sloth-runner","metrics_port":%d,"metrics_endpoint":"/metrics"}`, s.port)
	})

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Starting telemetry server", "port", s.port, "endpoint", "/metrics")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Telemetry server failed", "error", err)
		}
	}()

	// Start runtime metrics updater
	go s.updateRuntimeMetricsLoop()

	return nil
}

// Stop gracefully shuts down the metrics server
func (s *Server) Stop(ctx context.Context) error {
	if !s.enabled || s.httpServer == nil {
		return nil
	}

	slog.Info("Stopping telemetry server")
	return s.httpServer.Shutdown(ctx)
}

// GetMetrics returns the metrics instance
func (s *Server) GetMetrics() *Metrics {
	return s.metrics
}

// GetPort returns the server port
func (s *Server) GetPort() int {
	return s.port
}

// GetEndpoint returns the full metrics endpoint URL
func (s *Server) GetEndpoint() string {
	return fmt.Sprintf("http://localhost:%d/metrics", s.port)
}

// IsEnabled returns whether telemetry is enabled
func (s *Server) IsEnabled() bool {
	return s.enabled
}

// updateRuntimeMetricsLoop periodically updates runtime metrics
func (s *Server) updateRuntimeMetricsLoop() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		s.metrics.UpdateRuntimeMetrics()
	}
}

// GetRegistry returns the Prometheus registry (useful for testing)
func (s *Server) GetRegistry() *prometheus.Registry {
	return s.registry
}
