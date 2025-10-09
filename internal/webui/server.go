package webui

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/chalkan3-sloth/sloth-runner/internal/config"
	"github.com/chalkan3-sloth/sloth-runner/internal/metrics"
	"github.com/chalkan3-sloth/sloth-runner/internal/webui/handlers"
	"github.com/chalkan3-sloth/sloth-runner/internal/webui/middleware"
	"github.com/chalkan3-sloth/sloth-runner/internal/webui/services"
)

//go:embed static/css static/js templates
var embeddedFS embed.FS

// Server represents the web UI server
type Server struct {
	router           *gin.Engine
	port             int
	wsHub            *handlers.WebSocketHub
	agentDB          *handlers.AgentDBWrapper
	slothRepo        *handlers.SlothRepoWrapper
	hookRepo         *handlers.HookRepoWrapper
	secretsSvc       *handlers.SecretsServiceWrapper
	sshDB            *handlers.SSHDBWrapper
	stackHandler     *handlers.StackHandler
	httpServer       *http.Server
	agentClient      *services.AgentClient
	metricsDB        *metrics.MetricsDB
	metricsCollector *metrics.Collector
}

// Config holds server configuration
type Config struct {
	Port           int
	Debug          bool
	AgentDBPath    string
	SlothDBPath    string
	HookDBPath     string
	SecretsDBPath  string
	SSHDBPath      string
	StackDBPath    string
	MetricsDBPath  string
	EnableAuth     bool
	Username       string
	Password       string
}

// NewServer creates a new web UI server
func NewServer(cfg *Config) (*Server, error) {
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Initialize WebSocket hub
	wsHub := handlers.NewWebSocketHub()
	go wsHub.Run()

	// Initialize database wrappers
	agentDB, err := handlers.NewAgentDBWrapper(cfg.AgentDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize agent DB: %w", err)
	}

	// Initialize agent groups schema
	if err := agentDB.InitializeAgentGroupsSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize agent groups schema: %w", err)
	}

	// Initialize advanced agent groups schema
	if err := agentDB.InitializeAdvancedSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize advanced agent groups schema: %w", err)
	}

	slothRepo, err := handlers.NewSlothRepoWrapper(cfg.SlothDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sloth repository: %w", err)
	}

	hookRepo, err := handlers.NewHookRepoWrapper(cfg.HookDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize hook repository: %w", err)
	}

	secretsSvc, err := handlers.NewSecretsServiceWrapper(cfg.SecretsDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize secrets service: %w", err)
	}

	sshDB, err := handlers.NewSSHDBWrapper(cfg.SSHDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize SSH database: %w", err)
	}

	stackHandler, err := handlers.NewStackHandler(cfg.StackDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize stack handler: %w", err)
	}

	// Initialize agent client for gRPC communication
	agentClient := services.NewAgentClient()

	// NOTE: Metrics database and collector are now initialized in the master server
	// The UI simply reads from the existing metrics database
	metricsDBPath := cfg.MetricsDBPath
	if metricsDBPath == "" {
		metricsDBPath = config.GetMetricsDBPath()
	}

	slog.Info("Connecting to metrics database", "path", metricsDBPath)
	metricsDB, err := metrics.NewMetricsDB(metricsDBPath)
	if err != nil {
		slog.Warn("Failed to connect to metrics database", "error", err)
		slog.Info("Metrics features will be unavailable. Make sure the master server is running.")
		metricsDB = nil
	}

	server := &Server{
		router:       router,
		port:         cfg.Port,
		wsHub:        wsHub,
		agentDB:      agentDB,
		slothRepo:    slothRepo,
		hookRepo:     hookRepo,
		secretsSvc:   secretsSvc,
		sshDB:        sshDB,
		stackHandler: stackHandler,
		agentClient:  agentClient,
		metricsDB:    metricsDB,
	}

	// Setup routes
	server.setupRoutes(cfg)

	return server, nil
}

// setupRoutes configures all HTTP routes
func (s *Server) setupRoutes(cfg *Config) {
	// Serve static files from embedded FS
	staticFS, _ := fs.Sub(embeddedFS, "static")
	s.router.StaticFS("/static", http.FS(staticFS))

	// Serve templates from embedded FS
	templatesFS, _ := fs.Sub(embeddedFS, "templates")
	s.router.StaticFS("/templates", http.FS(templatesFS))

	// Authentication middleware (optional)
	if cfg.EnableAuth {
		auth := middleware.NewBasicAuth(cfg.Username, cfg.Password)
		s.router.Use(auth.Middleware())
	}

	// CORS middleware
	s.router.Use(middleware.CORS())

	// API v1 routes
	api := s.router.Group("/api/v1")
	{
		// Dashboard
		dashboard := handlers.NewDashboardHandler(s.agentDB, s.slothRepo, s.hookRepo, s.agentClient)
		api.GET("/dashboard", dashboard.GetStats)

		// Agents
		agentHandler := handlers.NewAgentHandler(s.agentDB, s.agentClient)
		agents := api.Group("/agents")
		{
			agents.GET("", agentHandler.List)
			agents.GET("/:name", agentHandler.Get)
			agents.DELETE("/:name", agentHandler.Delete)
			agents.GET("/:name/stats", agentHandler.GetStats)
			agents.POST("/:name/update", agentHandler.Update)

			// Advanced Agent Management
			agents.GET("/:name/resources", agentHandler.GetResourceUsage)
			agents.GET("/:name/processes", agentHandler.GetProcessList)
			agents.GET("/:name/network", agentHandler.GetNetworkInfo)
			agents.GET("/:name/disk", agentHandler.GetDiskInfo)
			agents.POST("/:name/command", agentHandler.ExecuteCommand)
			agents.POST("/:name/restart", agentHandler.RestartAgent)
			agents.POST("/:name/shutdown", agentHandler.ShutdownAgent)
			agents.GET("/:name/logs/stream", agentHandler.StreamLogs)
			agents.GET("/:name/metrics/stream", agentHandler.StreamMetrics)

			// Enhanced Troubleshooting & Diagnostics
			diagnosticsHandler := handlers.NewAgentDiagnosticsHandler(s.agentDB, s.agentClient)
			agents.GET("/:name/metrics/detailed", diagnosticsHandler.GetDetailedMetrics)
			agents.GET("/:name/logs", diagnosticsHandler.GetRecentLogs)
			agents.GET("/:name/connections", diagnosticsHandler.GetActiveConnections)
			agents.GET("/:name/errors", diagnosticsHandler.GetSystemErrors)
			agents.GET("/:name/performance/history", diagnosticsHandler.GetPerformanceHistory)
			agents.GET("/:name/health/diagnose", diagnosticsHandler.DiagnoseHealth)

			// Persistent Metrics History (new metrics database)
			metricsHistoryHandler := handlers.NewMetricsHistoryHandler(s.metricsDB)
			agents.GET("/:name/metrics/history", metricsHistoryHandler.GetAgentMetricsHistory)

			// Bulk operations
			agents.POST("/bulk/execute", agentHandler.BulkExecute)
			agents.POST("/bulk/status", agentHandler.GetMultipleStatus)
		}

		// Agent Groups
		groupHandler := handlers.NewAgentGroupHandler(s.agentDB)
		advancedGroupHandler := handlers.NewAgentGroupAdvancedHandler(s.agentDB)
		groups := api.Group("/agent-groups")
		{
			// Basic operations
			groups.GET("", groupHandler.List)
			groups.GET("/:name", groupHandler.Get)
			groups.POST("", groupHandler.Create)
			groups.DELETE("/:name", groupHandler.Delete)
			groups.POST("/:name/agents", groupHandler.AddAgents)
			groups.DELETE("/:name/agents", groupHandler.RemoveAgents)
			groups.GET("/:name/metrics", groupHandler.GetAggregatedMetrics)

			// Advanced operations
			// Bulk operations
			groups.POST("/bulk-operation", advancedGroupHandler.ExecuteBulkOperation)

			// Templates
			groups.GET("/templates", advancedGroupHandler.ListTemplates)
			groups.GET("/templates/:id", advancedGroupHandler.GetTemplate)
			groups.POST("/templates", advancedGroupHandler.CreateTemplate)
			groups.DELETE("/templates/:id", advancedGroupHandler.DeleteTemplate)
			groups.POST("/templates/:id/apply", advancedGroupHandler.ApplyTemplate)

			// Hierarchy
			groups.POST("/:name/hierarchy", advancedGroupHandler.SetGroupHierarchy)
			groups.GET("/:name/hierarchy", advancedGroupHandler.GetGroupHierarchy)
			groups.DELETE("/:name/hierarchy", advancedGroupHandler.RemoveGroupHierarchy)
			groups.GET("/:name/children", advancedGroupHandler.GetGroupChildren)

			// Auto-discovery
			groups.GET("/auto-discovery", advancedGroupHandler.ListAutoDiscoveryConfigs)
			groups.GET("/auto-discovery/:id", advancedGroupHandler.GetAutoDiscoveryConfig)
			groups.POST("/auto-discovery", advancedGroupHandler.CreateAutoDiscoveryConfig)
			groups.PUT("/auto-discovery/:id", advancedGroupHandler.UpdateAutoDiscoveryConfig)
			groups.DELETE("/auto-discovery/:id", advancedGroupHandler.DeleteAutoDiscoveryConfig)
			groups.POST("/auto-discovery/:id/run", advancedGroupHandler.RunAutoDiscovery)

			// Webhooks
			groups.GET("/webhooks", advancedGroupHandler.ListWebhooks)
			groups.GET("/webhooks/:id", advancedGroupHandler.GetWebhook)
			groups.POST("/webhooks", advancedGroupHandler.CreateWebhook)
			groups.PUT("/webhooks/:id", advancedGroupHandler.UpdateWebhook)
			groups.DELETE("/webhooks/:id", advancedGroupHandler.DeleteWebhook)
			groups.GET("/webhooks/:id/logs", advancedGroupHandler.GetWebhookLogs)
		}

		// Workflows (Sloths)
		slothHandler := handlers.NewSlothHandler(s.slothRepo)
		sloths := api.Group("/sloths")
		{
			sloths.GET("", slothHandler.List)
			sloths.GET("/:name", slothHandler.Get)
			sloths.POST("", slothHandler.Create)
			sloths.PUT("/:name", slothHandler.Update)
			sloths.DELETE("/:name", slothHandler.Delete)
			sloths.POST("/:name/activate", slothHandler.Activate)
			sloths.POST("/:name/deactivate", slothHandler.Deactivate)
			sloths.POST("/:name/run", slothHandler.Run)
		}

		// Hooks
		hookHandler := handlers.NewHookHandler(s.hookRepo)
		hooks := api.Group("/hooks")
		{
			hooks.GET("", hookHandler.List)
			hooks.GET("/stats", hookHandler.GetStatistics)
			hooks.GET("/by-event-type", hookHandler.ListByEventType)
			hooks.GET("/by-stack", hookHandler.ListByStack)
			hooks.GET("/:id", hookHandler.Get)
			hooks.POST("", hookHandler.Create)
			hooks.PUT("/:id", hookHandler.Update)
			hooks.DELETE("/:id", hookHandler.Delete)
			hooks.POST("/:id/enable", hookHandler.Enable)
			hooks.POST("/:id/disable", hookHandler.Disable)
			hooks.GET("/:id/history", hookHandler.GetHistory)
			hooks.GET("/:id/execution-stats", hookHandler.GetExecutionStats)
		}

		// Events
		eventHandler := handlers.NewEventHandler(s.hookRepo)
		events := api.Group("/events")
		{
			events.GET("", eventHandler.List)
			events.GET("/stats", eventHandler.GetStatistics)
			events.GET("/recent", eventHandler.GetRecentActivity)
			events.GET("/pending", eventHandler.ListPending)
			events.GET("/by-type", eventHandler.ListByType)
			events.GET("/by-status", eventHandler.ListByStatus)
			events.GET("/by-agent", eventHandler.ListByAgent)
			events.GET("/hook-executions/by-agent", eventHandler.ListHookExecutionsByAgent)
			events.GET("/:id", eventHandler.Get)
			events.POST("/:id/retry", eventHandler.Retry)
		}

		// Watchers
		watcherHandler := handlers.NewWatcherHandler(s.agentDB, s.agentClient)
		watchers := api.Group("/watchers")
		{
			watchers.GET("", watcherHandler.ListAllWatchers)
			watchers.GET("/stats", watcherHandler.GetStatistics)
			watchers.GET("/agent/:agent", watcherHandler.ListByAgent)
			watchers.GET("/agent/:agent/:id", watcherHandler.GetByAgent)
			watchers.POST("/agent/:agent", watcherHandler.CreateForAgent)
			watchers.DELETE("/agent/:agent/:id", watcherHandler.DeleteFromAgent)
		}

		// Network Metrics
		networkHandler := handlers.NewNetworkHandler(s.agentDB, s.agentClient)
		network := api.Group("/network")
		{
			network.GET("/summary", networkHandler.GetNetworkSummary)
			network.GET("/all", networkHandler.GetAllNetworkStats)
			network.GET("/topology", networkHandler.GetNetworkTopology)
			network.GET("/top", networkHandler.GetTopAgentsByNetwork)
			network.GET("/agent/:agent", networkHandler.GetNetworkStats)
			network.GET("/agent/:agent/interface/:interface", networkHandler.GetInterfaceDetails)
		}

		// Secrets
		secretHandler := handlers.NewSecretHandler(s.secretsSvc)
		secrets := api.Group("/secrets")
		{
			secrets.GET("/:stack", secretHandler.List)
			secrets.POST("/:stack", secretHandler.Add)
			secrets.DELETE("/:stack/:name", secretHandler.Delete)
		}

		// SSH Profiles
		sshHandler := handlers.NewSSHHandler(s.sshDB)
		ssh := api.Group("/ssh")
		{
			ssh.GET("", sshHandler.List)
			ssh.GET("/:name", sshHandler.Get)
			ssh.POST("", sshHandler.Create)
			ssh.PUT("/:name", sshHandler.Update)
			ssh.DELETE("/:name", sshHandler.Delete)
			ssh.GET("/:name/audit", sshHandler.GetAuditLogs)
		}

		// Stacks
		stacks := api.Group("/stacks")
		{
			stacks.GET("", s.stackHandler.List)
			stacks.GET("/:name", s.stackHandler.Get)
			stacks.POST("", s.stackHandler.Create)
			stacks.PUT("/:name", s.stackHandler.Update)
			stacks.DELETE("/:name", s.stackHandler.Delete)
			stacks.POST("/:name/variables", s.stackHandler.AddVariable)
			stacks.DELETE("/:name/variables/:key", s.stackHandler.DeleteVariable)
		}

		// Workflow Executions
		execHandler := handlers.NewWorkflowExecutionHandler(s.wsHub)
		executions := api.Group("/executions")
		{
			executions.POST("", execHandler.ExecuteWorkflow)
			executions.GET("", handlers.ListExecutionsHandler)
			executions.GET("/stats", handlers.GetExecutionStatsHandler)
			executions.GET("/:id", handlers.GetExecutionHandler)
			executions.POST("/:id/cancel", execHandler.CancelExecution)
			executions.GET("/:id/logs", execHandler.GetExecutionLogs)
			executions.DELETE("/cleanup", handlers.DeleteOldExecutionsHandler)
		}

		// Metrics
		metricsHandler := handlers.NewMetricsHandler(s.wsHub)
		api.GET("/metrics", metricsHandler.GetMetrics)
		api.GET("/metrics/history", metricsHandler.GetHistoricalMetrics)

		// Metrics History (new persistent metrics)
		metricsHistoryHandler := handlers.NewMetricsHistoryHandler(s.metricsDB)
		api.GET("/metrics/all", metricsHistoryHandler.GetAllAgentsMetrics)
		api.GET("/metrics/stats", metricsHistoryHandler.GetMetricsStats)

		// Logs
		logsHandler := handlers.NewLogsHandler(s.wsHub)
		logs := api.Group("/logs")
		{
			logs.GET("", logsHandler.ListLogFiles)
			logs.GET("/:filename", logsHandler.GetLogFile)
			logs.GET("/:filename/stream", logsHandler.StreamLogs)
		}

		// Scheduler
		schedulerHandler := handlers.NewSchedulerHandler(s.wsHub)
		scheduler := api.Group("/scheduler")
		{
			scheduler.GET("", schedulerHandler.ListSchedules)
			scheduler.GET("/:id", schedulerHandler.GetSchedule)
			scheduler.POST("", schedulerHandler.CreateSchedule)
			scheduler.PUT("/:id", schedulerHandler.UpdateSchedule)
			scheduler.DELETE("/:id", schedulerHandler.DeleteSchedule)
			scheduler.POST("/:id/enable", schedulerHandler.EnableSchedule)
			scheduler.POST("/:id/disable", schedulerHandler.DisableSchedule)
			scheduler.POST("/:id/trigger", schedulerHandler.TriggerSchedule)
		}

		// Backup & Restore
		backupHandler := handlers.NewBackupHandler()
		backup := api.Group("/backup")
		{
			backup.GET("", backupHandler.ListBackups)
			backup.POST("/create", backupHandler.CreateBackup)
			backup.POST("/restore", backupHandler.RestoreBackup)
		}

		// Terminal
		terminalHandler := handlers.NewTerminalHandler()
		terminal := api.Group("/terminal")
		{
			terminal.POST("", terminalHandler.CreateSession)
			terminal.GET("", terminalHandler.ListSessions)
			terminal.GET("/:id/ws", terminalHandler.ConnectTerminal)
			terminal.DELETE("/:id", terminalHandler.CloseSession)
		}

		// WebSocket
		api.GET("/ws", func(c *gin.Context) {
			handlers.ServeWebSocket(s.wsHub, c.Writer, c.Request)
		})
	}

	// Serve main UI pages
	s.router.GET("/", s.servePage("index.html"))
	s.router.GET("/agents", s.servePage("agents.html"))
	s.router.GET("/agent-control", s.servePage("agent-control.html"))
	s.router.GET("/agent-groups", s.servePage("agent-groups.html"))
	s.router.GET("/workflows", s.servePage("workflows.html"))
	s.router.GET("/stacks", s.servePage("stacks.html"))
	s.router.GET("/hooks", s.servePage("hooks.html"))
	s.router.GET("/events", s.servePage("events.html"))
	s.router.GET("/watchers", s.servePage("watchers.html"))
	s.router.GET("/network", s.servePage("network.html"))
	s.router.GET("/network/topology", s.servePage("network-topology.html"))
	s.router.GET("/secrets", s.servePage("secrets.html"))
	s.router.GET("/ssh", s.servePage("ssh.html"))
	s.router.GET("/executions", s.servePage("executions.html"))
	s.router.GET("/agent-dashboard", s.servePage("agent-dashboard.html"))
	s.router.GET("/metrics", s.servePage("metrics.html"))
	s.router.GET("/logs", s.servePage("logs.html"))
	s.router.GET("/history", s.servePage("history.html"))
	s.router.GET("/scheduler", s.servePage("scheduler.html"))
	s.router.GET("/terminal", s.servePage("terminal.html"))
	s.router.GET("/backup", s.servePage("backup.html"))
	s.router.GET("/chart-test", s.servePage("chart-test.html"))

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Unix(),
		})
	})
}

// getAgentList returns a list of all agents for metrics collection
func (s *Server) getAgentList() []metrics.AgentInfo {
	ctx := context.Background()
	agents, err := s.agentDB.ListAgents(ctx)
	if err != nil {
		slog.Error("Failed to get agent list for metrics collection", "error", err)
		return nil
	}

	slog.Info("Retrieved agents for metrics collection", "count", len(agents))

	agentInfos := make([]metrics.AgentInfo, 0, len(agents))
	for _, agent := range agents {
		agentInfos = append(agentInfos, metrics.AgentInfo{
			Name:    agent.Name,
			Address: agent.Address,
		})
		slog.Info("Adding agent to metrics collection", "name", agent.Name, "address", agent.Address)
	}

	return agentInfos
}

// servePage returns a handler that serves an HTML page
func (s *Server) servePage(filename string) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := embeddedFS.ReadFile("templates/" + filename)
		if err != nil {
			c.String(http.StatusNotFound, "Page not found")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	}
}

// Start starts the web UI server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting Sloth Runner Web UI on http://localhost%s", addr)

	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down web UI server...")

	// Stop metrics collector
	if s.metricsCollector != nil {
		s.metricsCollector.Stop()
	}

	// Close WebSocket hub
	s.wsHub.Shutdown()

	// Close agent client connections
	if s.agentClient != nil {
		s.agentClient.CloseAll()
	}

	// Close database connections
	if s.agentDB != nil {
		s.agentDB.Close()
	}
	if s.slothRepo != nil {
		s.slothRepo.Close()
	}
	if s.hookRepo != nil {
		s.hookRepo.Close()
	}
	if s.secretsSvc != nil {
		s.secretsSvc.Close()
	}
	if s.sshDB != nil {
		s.sshDB.Close()
	}
	if s.metricsDB != nil {
		s.metricsDB.Close()
	}

	// Shutdown HTTP server
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}

	return nil
}

// GetWebSocketHub returns the WebSocket hub
func (s *Server) GetWebSocketHub() *handlers.WebSocketHub {
	return s.wsHub
}
