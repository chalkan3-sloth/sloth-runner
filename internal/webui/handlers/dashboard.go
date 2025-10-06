package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard operations
type DashboardHandler struct {
	agentDB   *AgentDBWrapper
	slothRepo *SlothRepoWrapper
	hookRepo  *HookRepoWrapper
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(agentDB *AgentDBWrapper, slothRepo *SlothRepoWrapper, hookRepo *HookRepoWrapper) *DashboardHandler {
	return &DashboardHandler{
		agentDB:   agentDB,
		slothRepo: slothRepo,
		hookRepo:  hookRepo,
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
	pendingEvents, _ := eventQueue.GetPendingEvents(100)

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
			"total":   len(hooks),
			"enabled": enabledHooks,
		},
		"events": gin.H{
			"pending": len(pendingEvents),
		},
	})
}
