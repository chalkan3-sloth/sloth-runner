package handlers

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/chalkan3-sloth/sloth-runner/internal/webui/services"
	"github.com/gin-gonic/gin"
	pb "github.com/chalkan3-sloth/sloth-runner/proto"
)

// NetworkHandler handles network metrics operations
type NetworkHandler struct {
	agentDB     *AgentDBWrapper
	agentClient *services.AgentClient
}

// NewNetworkHandler creates a new network handler
func NewNetworkHandler(agentDB *AgentDBWrapper, agentClient *services.AgentClient) *NetworkHandler {
	return &NetworkHandler{
		agentDB:     agentDB,
		agentClient: agentClient,
	}
}

// NetworkStats represents aggregated network statistics
type NetworkStats struct {
	AgentName      string                 `json:"agent_name"`
	TotalRxBytes   uint64                 `json:"total_rx_bytes"`
	TotalTxBytes   uint64                 `json:"total_tx_bytes"`
	TotalRxMB      float64                `json:"total_rx_mb"`
	TotalTxMB      float64                `json:"total_tx_mb"`
	Interfaces     []InterfaceStats       `json:"interfaces"`
	Timestamp      time.Time              `json:"timestamp"`
}

// InterfaceStats represents statistics for a single network interface
type InterfaceStats struct {
	Name         string   `json:"name"`
	IPAddresses  []string `json:"ip_addresses"`
	MACAddress   string   `json:"mac_address"`
	BytesSent    uint64   `json:"bytes_sent"`
	BytesRecv    uint64   `json:"bytes_recv"`
	PacketsSent  uint32   `json:"packets_sent"`
	PacketsRecv  uint32   `json:"packets_recv"`
	IsUp         bool     `json:"is_up"`
	MBSent       float64  `json:"mb_sent"`
	MBRecv       float64  `json:"mb_recv"`
}

// BandwidthHistory represents bandwidth usage over time
type BandwidthHistory struct {
	Timestamp time.Time `json:"timestamp"`
	RxMBps    float64   `json:"rx_mbps"`
	TxMBps    float64   `json:"tx_mbps"`
	RxBytes   uint64    `json:"rx_bytes"`
	TxBytes   uint64    `json:"tx_bytes"`
}

// GetNetworkStats returns network statistics for a specific agent
func (h *NetworkHandler) GetNetworkStats(c *gin.Context) {
	agentName := c.Param("agent")
	if agentName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent parameter is required"})
		return
	}

	ctx := context.Background()
	agent, err := h.agentDB.GetAgent(ctx, agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Get agent client
	client, err := h.agentClient.GetClient(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to agent"})
		return
	}

	// Get network info
	resp, err := client.GetNetworkInfo(ctx, &pb.NetworkInfoRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get network info"})
		return
	}

	// Calculate totals and convert to stats
	var totalRx, totalTx uint64
	interfaces := make([]InterfaceStats, 0, len(resp.Interfaces))

	for _, iface := range resp.Interfaces {
		totalRx += iface.BytesRecv
		totalTx += iface.BytesSent

		interfaces = append(interfaces, InterfaceStats{
			Name:        iface.Name,
			IPAddresses: iface.IpAddresses,
			MACAddress:  iface.MacAddress,
			BytesSent:   iface.BytesSent,
			BytesRecv:   iface.BytesRecv,
			PacketsSent: iface.PacketsSent,
			PacketsRecv: iface.PacketsRecv,
			IsUp:        iface.IsUp,
			MBSent:      float64(iface.BytesSent) / 1024 / 1024,
			MBRecv:      float64(iface.BytesRecv) / 1024 / 1024,
		})
	}

	stats := NetworkStats{
		AgentName:    agentName,
		TotalRxBytes: totalRx,
		TotalTxBytes: totalTx,
		TotalRxMB:    float64(totalRx) / 1024 / 1024,
		TotalTxMB:    float64(totalTx) / 1024 / 1024,
		Interfaces:   interfaces,
		Timestamp:    time.Now(),
	}

	c.JSON(http.StatusOK, stats)
}

// GetAllNetworkStats returns network statistics for all agents
func (h *NetworkHandler) GetAllNetworkStats(c *gin.Context) {
	ctx := context.Background()
	agents, err := h.agentDB.ListAgents(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	allStats := make([]NetworkStats, 0, len(agents))
	var grandTotalRx, grandTotalTx uint64

	for _, agent := range agents {
		client, err := h.agentClient.GetClient(ctx, agent.Address)
		if err != nil {
			continue // Skip failed connections
		}

		resp, err := client.GetNetworkInfo(ctx, &pb.NetworkInfoRequest{})
		if err != nil {
			continue
		}

		var totalRx, totalTx uint64
		interfaces := make([]InterfaceStats, 0, len(resp.Interfaces))

		for _, iface := range resp.Interfaces {
			totalRx += iface.BytesRecv
			totalTx += iface.BytesSent

			interfaces = append(interfaces, InterfaceStats{
				Name:        iface.Name,
				IPAddresses: iface.IpAddresses,
				MACAddress:  iface.MacAddress,
				BytesSent:   iface.BytesSent,
				BytesRecv:   iface.BytesRecv,
				PacketsSent: iface.PacketsSent,
				PacketsRecv: iface.PacketsRecv,
				IsUp:        iface.IsUp,
				MBSent:      float64(iface.BytesSent) / 1024 / 1024,
				MBRecv:      float64(iface.BytesRecv) / 1024 / 1024,
			})
		}

		grandTotalRx += totalRx
		grandTotalTx += totalTx

		allStats = append(allStats, NetworkStats{
			AgentName:    agent.Name,
			TotalRxBytes: totalRx,
			TotalTxBytes: totalTx,
			TotalRxMB:    float64(totalRx) / 1024 / 1024,
			TotalTxMB:    float64(totalTx) / 1024 / 1024,
			Interfaces:   interfaces,
			Timestamp:    time.Now(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"agents":         allStats,
		"total_agents":   len(allStats),
		"total_rx_bytes": grandTotalRx,
		"total_tx_bytes": grandTotalTx,
		"total_rx_gb":    float64(grandTotalRx) / 1024 / 1024 / 1024,
		"total_tx_gb":    float64(grandTotalTx) / 1024 / 1024 / 1024,
		"timestamp":      time.Now(),
	})
}

// GetTopAgentsByNetwork returns top agents by network usage
func (h *NetworkHandler) GetTopAgentsByNetwork(c *gin.Context) {
	ctx := context.Background()
	agents, err := h.agentDB.ListAgents(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type AgentNetworkUsage struct {
		AgentName  string  `json:"agent_name"`
		TotalBytes uint64  `json:"total_bytes"`
		TotalGB    float64 `json:"total_gb"`
		RxBytes    uint64  `json:"rx_bytes"`
		TxBytes    uint64  `json:"tx_bytes"`
	}

	usage := make([]AgentNetworkUsage, 0, len(agents))

	for _, agent := range agents {
		client, err := h.agentClient.GetClient(ctx, agent.Address)
		if err != nil {
			continue
		}

		resp, err := client.GetNetworkInfo(ctx, &pb.NetworkInfoRequest{})
		if err != nil {
			continue
		}

		var totalRx, totalTx uint64
		for _, iface := range resp.Interfaces {
			totalRx += iface.BytesRecv
			totalTx += iface.BytesSent
		}

		total := totalRx + totalTx
		usage = append(usage, AgentNetworkUsage{
			AgentName:  agent.Name,
			TotalBytes: total,
			TotalGB:    float64(total) / 1024 / 1024 / 1024,
			RxBytes:    totalRx,
			TxBytes:    totalTx,
		})
	}

	// Sort by total usage
	sort.Slice(usage, func(i, j int) bool {
		return usage[i].TotalBytes > usage[j].TotalBytes
	})

	// Limit to top 10
	if len(usage) > 10 {
		usage = usage[:10]
	}

	c.JSON(http.StatusOK, gin.H{
		"top_agents": usage,
		"count":      len(usage),
	})
}

// GetInterfaceDetails returns detailed information about a specific interface
func (h *NetworkHandler) GetInterfaceDetails(c *gin.Context) {
	agentName := c.Param("agent")
	interfaceName := c.Param("interface")

	if agentName == "" || interfaceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "agent and interface parameters are required"})
		return
	}

	ctx := context.Background()
	agent, err := h.agentDB.GetAgent(ctx, agentName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	client, err := h.agentClient.GetClient(ctx, agent.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to agent"})
		return
	}

	resp, err := client.GetNetworkInfo(ctx, &pb.NetworkInfoRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get network info"})
		return
	}

	// Find the specific interface
	for _, iface := range resp.Interfaces {
		if iface.Name == interfaceName {
			stats := InterfaceStats{
				Name:        iface.Name,
				IPAddresses: iface.IpAddresses,
				MACAddress:  iface.MacAddress,
				BytesSent:   iface.BytesSent,
				BytesRecv:   iface.BytesRecv,
				PacketsSent: iface.PacketsSent,
				PacketsRecv: iface.PacketsRecv,
				IsUp:        iface.IsUp,
				MBSent:      float64(iface.BytesSent) / 1024 / 1024,
				MBRecv:      float64(iface.BytesRecv) / 1024 / 1024,
			}

			c.JSON(http.StatusOK, gin.H{
				"agent":     agentName,
				"interface": stats,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Interface not found"})
}

// GetNetworkTopology returns network topology information
func (h *NetworkHandler) GetNetworkTopology(c *gin.Context) {
	ctx := context.Background()
	agents, err := h.agentDB.ListAgents(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type Node struct {
		ID         string   `json:"id"`
		Name       string   `json:"name"`
		IPAddresses []string `json:"ip_addresses"`
		Status     string   `json:"status"`
		Type       string   `json:"type"`
	}

	type Link struct {
		Source string `json:"source"`
		Target string `json:"target"`
		Type   string `json:"type"`
	}

	nodes := make([]Node, 0, len(agents))
	links := make([]Link, 0)

	// Create nodes for each agent
	for _, agent := range agents {
		client, err := h.agentClient.GetClient(ctx, agent.Address)
		if err != nil {
			nodes = append(nodes, Node{
				ID:     agent.Name,
				Name:   agent.Name,
				Status: "offline",
				Type:   "agent",
			})
			continue
		}

		resp, err := client.GetNetworkInfo(ctx, &pb.NetworkInfoRequest{})
		if err != nil {
			nodes = append(nodes, Node{
				ID:     agent.Name,
				Name:   agent.Name,
				Status: "error",
				Type:   "agent",
			})
			continue
		}

		// Collect all IP addresses
		var ips []string
		for _, iface := range resp.Interfaces {
			if iface.IsUp {
				ips = append(ips, iface.IpAddresses...)
			}
		}

		nodes = append(nodes, Node{
			ID:          agent.Name,
			Name:        agent.Name,
			IPAddresses: ips,
			Status:      agent.Status,
			Type:        "agent",
		})
	}

	// TODO: Detect connections between agents (could use ping, traceroute, or connection logs)
	// For now, create links to master
	for _, node := range nodes {
		if node.Status == "Active" {
			links = append(links, Link{
				Source: "master",
				Target: node.ID,
				Type:   "grpc",
			})
		}
	}

	// Add master node
	masterNode := Node{
		ID:     "master",
		Name:   "Master",
		Status: "Active",
		Type:   "master",
	}
	nodes = append([]Node{masterNode}, nodes...)

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
		"links": links,
		"stats": gin.H{
			"total_nodes": len(nodes),
			"total_links": len(links),
			"active_agents": func() int {
				count := 0
				for _, n := range nodes {
					if n.Status == "Active" {
						count++
					}
				}
				return count
			}(),
		},
	})
}

// GetNetworkSummary returns a summary of network metrics
func (h *NetworkHandler) GetNetworkSummary(c *gin.Context) {
	ctx := context.Background()
	agents, err := h.agentDB.ListAgents(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var totalInterfaces int
	var totalActiveInterfaces int
	var grandTotalRx, grandTotalTx uint64

	for _, agent := range agents {
		client, err := h.agentClient.GetClient(ctx, agent.Address)
		if err != nil {
			continue
		}

		resp, err := client.GetNetworkInfo(ctx, &pb.NetworkInfoRequest{})
		if err != nil {
			continue
		}

		totalInterfaces += len(resp.Interfaces)

		for _, iface := range resp.Interfaces {
			if iface.IsUp {
				totalActiveInterfaces++
			}
			grandTotalRx += iface.BytesRecv
			grandTotalTx += iface.BytesSent
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_agents":           len(agents),
		"total_interfaces":       totalInterfaces,
		"active_interfaces":      totalActiveInterfaces,
		"total_rx_bytes":         grandTotalRx,
		"total_tx_bytes":         grandTotalTx,
		"total_rx_gb":            float64(grandTotalRx) / 1024 / 1024 / 1024,
		"total_tx_gb":            float64(grandTotalTx) / 1024 / 1024 / 1024,
		"total_traffic_gb":       float64(grandTotalRx+grandTotalTx) / 1024 / 1024 / 1024,
		"timestamp":              time.Now(),
	})
}
