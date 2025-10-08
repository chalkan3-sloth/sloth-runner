package handlers

import (
	"context"
	"net/http"

	"github.com/chalkan3-sloth/sloth-runner/internal/webui/services"
	"github.com/gin-gonic/gin"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
)

// DashboardHandler handles dashboard operations
type DashboardHandler struct {
	agentDB     *AgentDBWrapper
	slothRepo   *SlothRepoWrapper
	hookRepo    *HookRepoWrapper
	agentClient *services.AgentClient
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(agentDB *AgentDBWrapper, slothRepo *SlothRepoWrapper, hookRepo *HookRepoWrapper, agentClient *services.AgentClient) *DashboardHandler {
	return &DashboardHandler{
		agentDB:     agentDB,
		slothRepo:   slothRepo,
		hookRepo:    hookRepo,
		agentClient: agentClient,
	}
}

// GetStats returns dashboard statistics
func (h *DashboardHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	// Get agent stats
	agents, err := h.agentDB.ListAgents(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	activeAgents := 0
	for _, agent := range agents {
		if agent.Status == "Active" {
			activeAgents++
		}
	}

	// Get sloth stats
	sloths, err := h.slothRepo.List(ctx, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	activeSloths := 0
	for _, s := range sloths {
		if s.IsActive {
			activeSloths++
		}
	}

	// Get hook stats
	hooks, err := h.hookRepo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	enabledHooks := 0
	for _, hook := range hooks {
		if hook.Enabled {
			enabledHooks++
		}
	}

	// Get event queue stats
	eventQueue := h.hookRepo.GetEventQueue()
	pendingEvents, _ := eventQueue.GetPendingEvents(1000)
	allEvents, _ := eventQueue.ListEvents("", "", 1000)

	// Count events by status
	completedEvents := 0
	failedEvents := 0
	processingEvents := 0
	for _, event := range allEvents {
		switch event.Status {
		case "completed":
			completedEvents++
		case "failed":
			failedEvents++
		case "processing":
			processingEvents++
		}
	}

	// Get watcher stats from all agents
	totalWatchers := 0
	watchersByType := make(map[string]int)
	for _, agent := range agents {
		// Get agent client
		client, err := h.agentClient.GetClient(ctx, agent.Address)
		if err != nil {
			continue // Skip agents that can't be reached
		}

		resp, err := client.ListWatchers(context.Background(), &pb.ListWatchersRequest{})
		if err != nil {
			continue
		}

		totalWatchers += len(resp.Watchers)
		for _, w := range resp.Watchers {
			watchersByType[w.Type]++
		}
	}

	// Calculate hook execution stats
	totalHookRuns := int64(0)
	for _, hook := range hooks {
		totalHookRuns += hook.RunCount
	}

	c.JSON(http.StatusOK, gin.H{
		"agents": gin.H{
			"total":  len(agents),
			"active": activeAgents,
		},
		"workflows": gin.H{
			"total":  len(sloths),
			"active": activeSloths,
		},
		"hooks": gin.H{
			"total":       len(hooks),
			"enabled":     enabledHooks,
			"total_runs":  totalHookRuns,
		},
		"events": gin.H{
			"total":      len(allEvents),
			"pending":    len(pendingEvents),
			"processing": processingEvents,
			"completed":  completedEvents,
			"failed":     failedEvents,
		},
		"watchers": gin.H{
			"total":   totalWatchers,
			"by_type": watchersByType,
		},
	})
}
