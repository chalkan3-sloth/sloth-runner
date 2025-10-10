package process

import (
	"fmt"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessManager interface para gerenciamento de processos
type ProcessManager interface {
	List(options ListOptions) ([]*ProcessInfo, error)
	Kill(pid int32, signal string) error
	Info(pid int32) (*ProcessDetail, error)
	Monitor(pid int32, duration time.Duration) (*ProcessMetrics, error)
}

// ListOptions opções para listar processos
type ListOptions struct {
	SortBy    string // cpu, memory, name, pid
	Top       int    // limite de resultados
	Filter    string // filtro por nome
	UserFilter string // filtro por usuário
}

// ProcessInfo informações básicas do processo
type ProcessInfo struct {
	PID         int32
	Name        string
	Username    string
	CPUPercent  float64
	MemoryMB    float64
	MemoryPercent float32
	Status      string
	CreateTime  int64
	NumThreads  int32
	Cmdline     string
}

// ProcessDetail informações detalhadas do processo
type ProcessDetail struct {
	*ProcessInfo
	ParentPID   int32
	Nice        int32
	IOCounters  *process.IOCountersStat
	NumFDs      int32
	Connections []string
	OpenFiles   []string
	Environ     []string
}

// ProcessMetrics métricas de processo ao longo do tempo
type ProcessMetrics struct {
	PID           int32
	Duration      time.Duration
	Samples       []*ProcessSnapshot
	AvgCPU        float64
	MaxCPU        float64
	AvgMemory     float64
	MaxMemory     float64
}

// ProcessSnapshot snapshot de métricas
type ProcessSnapshot struct {
	Timestamp     time.Time
	CPUPercent    float64
	MemoryMB      float64
	MemoryPercent float32
	NumThreads    int32
}

// SystemProcessManager implementação padrão
type SystemProcessManager struct{}

// NewProcessManager cria um novo process manager
func NewProcessManager() ProcessManager {
	return &SystemProcessManager{}
}

// List lista processos
func (m *SystemProcessManager) List(options ListOptions) ([]*ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var results []*ProcessInfo

	for _, p := range processes {
		info, err := m.getProcessInfo(p)
		if err != nil {
			continue // Skip processos que não conseguimos acessar
		}

		// Aplica filtros
		if options.Filter != "" {
			if !contains(info.Name, options.Filter) && !contains(info.Cmdline, options.Filter) {
				continue
			}
		}

		if options.UserFilter != "" && info.Username != options.UserFilter {
			continue
		}

		results = append(results, info)
	}

	// Ordena
	m.sortProcesses(results, options.SortBy)

	// Limita resultados
	if options.Top > 0 && len(results) > options.Top {
		results = results[:options.Top]
	}

	return results, nil
}

// Kill mata um processo
func (m *SystemProcessManager) Kill(pid int32, signal string) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return fmt.Errorf("process not found: %v", err)
	}

	name, _ := p.Name()

	switch signal {
	case "SIGTERM", "TERM", "15":
		err = p.Terminate()
	case "SIGKILL", "KILL", "9":
		err = p.Kill()
	case "SIGINT", "INT", "2":
		err = p.SendSignal(2)
	case "SIGHUP", "HUP", "1":
		err = p.SendSignal(1)
	default:
		return fmt.Errorf("unknown signal: %s", signal)
	}

	if err != nil {
		return fmt.Errorf("failed to kill process %s (PID %d): %v", name, pid, err)
	}

	return nil
}

// Info obtém informações detalhadas de um processo
func (m *SystemProcessManager) Info(pid int32) (*ProcessDetail, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("process not found: %v", err)
	}

	info, err := m.getProcessInfo(p)
	if err != nil {
		return nil, err
	}

	detail := &ProcessDetail{
		ProcessInfo: info,
	}

	// Parent PID
	if ppid, err := p.Ppid(); err == nil {
		detail.ParentPID = ppid
	}

	// Nice
	if nice, err := p.Nice(); err == nil {
		detail.Nice = nice
	}

	// IO Counters
	if io, err := p.IOCounters(); err == nil {
		detail.IOCounters = io
	}

	// Num FDs
	if fds, err := p.NumFDs(); err == nil {
		detail.NumFDs = int32(fds)
	}

	// Connections
	if conns, err := p.Connections(); err == nil {
		for _, conn := range conns {
			detail.Connections = append(detail.Connections,
				fmt.Sprintf("%s:%d -> %s:%d", conn.Laddr.IP, conn.Laddr.Port, conn.Raddr.IP, conn.Raddr.Port))
		}
	}

	// Open Files
	if files, err := p.OpenFiles(); err == nil {
		for _, file := range files {
			detail.OpenFiles = append(detail.OpenFiles, file.Path)
		}
	}

	// Environment
	if env, err := p.Environ(); err == nil {
		detail.Environ = env
	}

	return detail, nil
}

// Monitor monitora um processo ao longo do tempo
func (m *SystemProcessManager) Monitor(pid int32, duration time.Duration) (*ProcessMetrics, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("process not found: %v", err)
	}

	metrics := &ProcessMetrics{
		PID:      pid,
		Duration: duration,
		Samples:  []*ProcessSnapshot{},
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	timeout := time.After(duration)

	var sumCPU, sumMemory float64

	for {
		select {
		case <-timeout:
			// Calcula médias
			if len(metrics.Samples) > 0 {
				metrics.AvgCPU = sumCPU / float64(len(metrics.Samples))
				metrics.AvgMemory = sumMemory / float64(len(metrics.Samples))
			}
			return metrics, nil

		case <-ticker.C:
			snapshot := &ProcessSnapshot{
				Timestamp: time.Now(),
			}

			if cpu, err := p.CPUPercent(); err == nil {
				snapshot.CPUPercent = cpu
				sumCPU += cpu
				if cpu > metrics.MaxCPU {
					metrics.MaxCPU = cpu
				}
			}

			if mem, err := p.MemoryInfo(); err == nil {
				snapshot.MemoryMB = float64(mem.RSS) / 1024 / 1024
				sumMemory += snapshot.MemoryMB
				if snapshot.MemoryMB > metrics.MaxMemory {
					metrics.MaxMemory = snapshot.MemoryMB
				}
			}

			if memPct, err := p.MemoryPercent(); err == nil {
				snapshot.MemoryPercent = memPct
			}

			if threads, err := p.NumThreads(); err == nil {
				snapshot.NumThreads = threads
			}

			metrics.Samples = append(metrics.Samples, snapshot)
		}
	}
}

// getProcessInfo obtém informações básicas de um processo
func (m *SystemProcessManager) getProcessInfo(p *process.Process) (*ProcessInfo, error) {
	info := &ProcessInfo{
		PID: p.Pid,
	}

	if name, err := p.Name(); err == nil {
		info.Name = name
	}

	if username, err := p.Username(); err == nil {
		info.Username = username
	}

	if cpu, err := p.CPUPercent(); err == nil {
		info.CPUPercent = cpu
	}

	if mem, err := p.MemoryInfo(); err == nil {
		info.MemoryMB = float64(mem.RSS) / 1024 / 1024
	}

	if memPct, err := p.MemoryPercent(); err == nil {
		info.MemoryPercent = memPct
	}

	if status, err := p.Status(); err == nil {
		if len(status) > 0 {
			info.Status = status[0]
		}
	}

	if createTime, err := p.CreateTime(); err == nil {
		info.CreateTime = createTime
	}

	if threads, err := p.NumThreads(); err == nil {
		info.NumThreads = threads
	}

	if cmdline, err := p.Cmdline(); err == nil {
		info.Cmdline = cmdline
	}

	return info, nil
}

// sortProcesses ordena processos
func (m *SystemProcessManager) sortProcesses(processes []*ProcessInfo, sortBy string) {
	switch sortBy {
	case "cpu":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].CPUPercent > processes[j].CPUPercent
		})
	case "memory", "mem":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].MemoryMB > processes[j].MemoryMB
		})
	case "name":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].Name < processes[j].Name
		})
	case "pid":
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].PID < processes[j].PID
		})
	default:
		// Default: ordenar por CPU
		sort.Slice(processes, func(i, j int) bool {
			return processes[i].CPUPercent > processes[j].CPUPercent
		})
	}
}

// contains verifica se uma string contém outra (case insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		   (s == substr || len(substr) == 0 ||
			(len(s) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if matchAt(s, substr, i) {
			return true
		}
	}
	return false
}

func matchAt(s, substr string, pos int) bool {
	for i := 0; i < len(substr); i++ {
		if toLower(s[pos+i]) != toLower(substr[i]) {
			return false
		}
	}
	return true
}

func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}
