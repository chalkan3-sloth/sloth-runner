package webui

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/chalkan3-sloth/sloth-runner/internal/webui/handlers"
	"github.com/chalkan3-sloth/sloth-runner/internal/webui/middleware"
)

//go:embed static/css static/js templates
var embeddedFS embed.FS

// Server represents the web UI server
type Server struct {
	router     *gin.Engine
	port       int
	wsHub      *handlers.WebSocketHub
	agentDB    *handlers.AgentDBWrapper
	slothRepo  *handlers.SlothRepoWrapper
	hookRepo   *handlers.HookRepoWrapper
	secretsSvc *handlers.SecretsServiceWrapper
	sshDB      *handlers.SSHDBWrapper
	httpServer *http.Server
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

	server := &Server{
		router:     router,
		port:       cfg.Port,
		wsHub:      wsHub,
		agentDB:    agentDB,
		slothRepo:  slothRepo,
		hookRepo:   hookRepo,
		secretsSvc: secretsSvc,
		sshDB:      sshDB,
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
		dashboard := handlers.NewDashboardHandler(s.agentDB, s.slothRepo, s.hookRepo)
		api.GET("/dashboard", dashboard.GetStats)

		// Agents
		agentHandler := handlers.NewAgentHandler(s.agentDB)
		agents := api.Group("/agents")
		{
			agents.GET("", agentHandler.List)
			agents.GET("/:name", agentHandler.Get)
			agents.DELETE("/:name", agentHandler.Delete)
			agents.GET("/:name/stats", agentHandler.GetStats)
			agents.POST("/:name/update", agentHandler.Update)
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
			hooks.GET("/:id", hookHandler.Get)
			hooks.POST("", hookHandler.Create)
			hooks.PUT("/:id", hookHandler.Update)
			hooks.DELETE("/:id", hookHandler.Delete)
			hooks.POST("/:id/enable", hookHandler.Enable)
			hooks.POST("/:id/disable", hookHandler.Disable)
			hooks.GET("/:id/history", hookHandler.GetHistory)
		}

		// Events
		eventHandler := handlers.NewEventHandler(s.hookRepo)
		events := api.Group("/events")
		{
			events.GET("", eventHandler.List)
			events.GET("/pending", eventHandler.ListPending)
			events.GET("/:id", eventHandler.Get)
			events.POST("/:id/retry", eventHandler.Retry)
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

		// Workflow Executions
		execHandler := handlers.NewWorkflowExecutionHandler(s.wsHub)
		executions := api.Group("/executions")
		{
			executions.POST("", execHandler.ExecuteWorkflow)
			executions.GET("", execHandler.ListExecutions)
			executions.GET("/:id", execHandler.GetExecution)
			executions.POST("/:id/cancel", execHandler.CancelExecution)
			executions.GET("/:id/logs", execHandler.GetExecutionLogs)
		}

		// Metrics
		metricsHandler := handlers.NewMetricsHandler(s.wsHub)
		api.GET("/metrics", metricsHandler.GetMetrics)
		api.GET("/metrics/history", metricsHandler.GetHistoricalMetrics)

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
	s.router.GET("/workflows", s.servePage("workflows.html"))
	s.router.GET("/hooks", s.servePage("hooks.html"))
	s.router.GET("/events", s.servePage("events.html"))
	s.router.GET("/secrets", s.servePage("secrets.html"))
	s.router.GET("/ssh", s.servePage("ssh.html"))
	s.router.GET("/executions", s.servePage("executions.html"))
	s.router.GET("/agent-dashboard", s.servePage("agent-dashboard.html"))
	s.router.GET("/metrics", s.servePage("metrics.html"))
	s.router.GET("/logs", s.servePage("logs.html"))
	s.router.GET("/scheduler", s.servePage("scheduler.html"))
	s.router.GET("/terminal", s.servePage("terminal.html"))
	s.router.GET("/backup", s.servePage("backup.html"))

	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"time":   time.Now().Unix(),
		})
	})
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

	// Close WebSocket hub
	s.wsHub.Shutdown()

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
