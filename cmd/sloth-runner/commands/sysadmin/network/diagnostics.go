package network

import (
	"fmt"
	"net"
	"time"

	"github.com/go-ping/ping"
)

// NetworkDiagnostics interface para diagnósticos de rede
type NetworkDiagnostics interface {
	Ping(host string, count int, timeout time.Duration) (*PingResult, error)
	CheckPort(host string, port int, timeout time.Duration) (*PortResult, error)
	CheckPorts(host string, ports []int, timeout time.Duration) ([]*PortResult, error)
}

// PingResult contém resultados de um ping
type PingResult struct {
	Host        string
	PacketsSent int
	PacketsRecv int
	PacketLoss  float64
	MinRTT      time.Duration
	AvgRTT      time.Duration
	MaxRTT      time.Duration
	StdDevRTT   time.Duration
}

// PortResult contém resultado de verificação de porta
type PortResult struct {
	Host    string
	Port    int
	Open    bool
	Service string
	Latency time.Duration
	Error   error
}

// SystemDiagnostics implementação padrão de NetworkDiagnostics
type SystemDiagnostics struct{}

// NewDiagnostics cria um novo diagnostics
func NewDiagnostics() NetworkDiagnostics {
	return &SystemDiagnostics{}
}

// Ping executa ping ICMP para um host
func (d *SystemDiagnostics) Ping(host string, count int, timeout time.Duration) (*PingResult, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return nil, fmt.Errorf("failed to create pinger: %w", err)
	}

	// Configure pinger
	pinger.Count = count
	pinger.Timeout = timeout
	pinger.SetPrivileged(false) // Use unprivileged mode (UDP)

	// Run ping
	err = pinger.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to run ping: %w", err)
	}

	stats := pinger.Statistics()

	result := &PingResult{
		Host:        host,
		PacketsSent: stats.PacketsSent,
		PacketsRecv: stats.PacketsRecv,
		PacketLoss:  stats.PacketLoss,
		MinRTT:      stats.MinRtt,
		AvgRTT:      stats.AvgRtt,
		MaxRTT:      stats.MaxRtt,
		StdDevRTT:   stats.StdDevRtt,
	}

	return result, nil
}

// CheckPort verifica se uma porta está aberta
func (d *SystemDiagnostics) CheckPort(host string, port int, timeout time.Duration) (*PortResult, error) {
	result := &PortResult{
		Host: host,
		Port: port,
	}

	address := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()

	conn, err := net.DialTimeout("tcp", address, timeout)
	latency := time.Since(start)

	if err != nil {
		result.Open = false
		result.Error = err
		result.Latency = latency
		return result, nil
	}

	defer conn.Close()

	result.Open = true
	result.Latency = latency
	result.Service = getServiceName(port)

	return result, nil
}

// CheckPorts verifica múltiplas portas
func (d *SystemDiagnostics) CheckPorts(host string, ports []int, timeout time.Duration) ([]*PortResult, error) {
	results := make([]*PortResult, len(ports))

	for i, port := range ports {
		result, err := d.CheckPort(host, port, timeout)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

// getServiceName retorna o nome do serviço para uma porta conhecida
func getServiceName(port int) string {
	services := map[int]string{
		20:    "FTP Data",
		21:    "FTP",
		22:    "SSH",
		23:    "Telnet",
		25:    "SMTP",
		53:    "DNS",
		80:    "HTTP",
		110:   "POP3",
		143:   "IMAP",
		443:   "HTTPS",
		465:   "SMTPS",
		587:   "SMTP (submission)",
		993:   "IMAPS",
		995:   "POP3S",
		3000:  "Development Server",
		3306:  "MySQL",
		5432:  "PostgreSQL",
		5672:  "RabbitMQ",
		6379:  "Redis",
		8000:  "HTTP Alt",
		8080:  "HTTP Proxy",
		8443:  "HTTPS Alt",
		9090:  "Prometheus",
		9200:  "Elasticsearch",
		27017: "MongoDB",
		50053: "gRPC (sloth-runner)",
	}

	if name, ok := services[port]; ok {
		return name
	}
	return "Unknown"
}

// TCPPing executa um TCP ping (alternativa ao ICMP ping)
func TCPPing(host string, port int, count int, timeout time.Duration) (*PingResult, error) {
	result := &PingResult{
		Host:        host,
		PacketsSent: count,
	}

	var totalRTT time.Duration
	var minRTT time.Duration = time.Hour
	var maxRTT time.Duration
	var successCount int

	for i := 0; i < count; i++ {
		address := fmt.Sprintf("%s:%d", host, port)
		start := time.Now()

		conn, err := net.DialTimeout("tcp", address, timeout)
		rtt := time.Since(start)

		if err == nil {
			conn.Close()
			successCount++
			totalRTT += rtt

			if rtt < minRTT {
				minRTT = rtt
			}
			if rtt > maxRTT {
				maxRTT = rtt
			}
		}

		if i < count-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	result.PacketsRecv = successCount
	result.PacketLoss = float64(count-successCount) / float64(count) * 100.0

	if successCount > 0 {
		result.MinRTT = minRTT
		result.AvgRTT = totalRTT / time.Duration(successCount)
		result.MaxRTT = maxRTT
	}

	return result, nil
}
