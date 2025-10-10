# Plano de Implementação - Comandos Sysadmin

**Data:** 2025-10-10
**Status:** 🚧 Em Progresso

## 📊 Status Atual

### Comandos Implementados ✅
1. **packages** - Gerenciamento de pacotes (apt, yum, dnf)
   - ✅ list, search, install, remove, update
   - Arquivo: `manager.go` com abstração de package managers

2. **services** - Gerenciamento de serviços (systemd, init.d)
   - ✅ list, status, start, stop, restart, enable, disable
   - Arquivo: `manager.go` com abstração de service managers

### Comandos Pendentes 🔨
3. **resources** - Monitoramento de recursos do sistema
4. **network** - Diagnósticos de rede
5. **performance** - Monitoramento de performance
6. **maintenance** - Manutenção do sistema
7. **config** - Gerenciamento de configuração
8. **backup** - Backup e restore
9. **deployment** - Deploy e rollback
10. **security** - Auditoria de segurança

## 🎯 Plano de Implementação

### Prioridade 1: Resources (Alta Prioridade)

**Funcionalidades:**
- `overview` - Visão geral de CPU, memória, disco, rede
- `cpu` - Uso detalhado de CPU (por core, load average)
- `memory` - Estatísticas de memória (RAM, swap, buffers, cache)
- `disk` - Uso de disco por filesystem
- `io` - Estatísticas de I/O (read/write, IOPS)
- `network` - Estatísticas de interface de rede
- `top` - Processos que mais consomem recursos

**Implementação:**
```go
// monitor.go - Interface para monitoramento de recursos
type ResourceMonitor interface {
    GetCPU() (*CPUStats, error)
    GetMemory() (*MemoryStats, error)
    GetDisk() ([]*DiskStats, error)
    GetNetwork() ([]*NetworkStats, error)
    GetProcesses(limit int) ([]*ProcessStats, error)
}

// Structs para cada tipo de estatística
type CPUStats struct {
    Usage       float64   // Porcentagem total
    PerCore     []float64 // Por core
    LoadAverage [3]float64 // 1, 5, 15 min
}

type MemoryStats struct {
    Total       uint64
    Used        uint64
    Free        uint64
    Available   uint64
    SwapTotal   uint64
    SwapUsed    uint64
}

type DiskStats struct {
    Filesystem  string
    MountPoint  string
    Total       uint64
    Used        uint64
    Available   uint64
    UsagePercent float64
}
```

**Comando de Sistema:**
- Linux: `/proc/stat`, `/proc/meminfo`, `df`, `/proc/net/dev`
- macOS: `sysctl`, `vm_stat`, `df`, `netstat`

**Estimativa:** 3-4 horas

---

### Prioridade 2: Network (Alta Prioridade)

**Funcionalidades:**
- `ping` - Testa conectividade com agents
- `port-check` - Verifica se porta está aberta

**Implementação:**
```go
// network.go
type NetworkDiagnostics interface {
    Ping(host string, count int) (*PingResult, error)
    CheckPort(host string, port int) (*PortResult, error)
}

type PingResult struct {
    Host         string
    PacketsSent  int
    PacketsRecv  int
    PacketLoss   float64
    MinRTT       time.Duration
    AvgRTT       time.Duration
    MaxRTT       time.Duration
}

type PortResult struct {
    Host    string
    Port    int
    Open    bool
    Service string // Detecção de serviço
    Latency time.Duration
}
```

**Implementação:**
- Usar `net.Dial` para port check
- ICMP para ping (ou TCP ping se ICMP não disponível)

**Estimativa:** 2-3 horas

---

### Prioridade 3: Performance (Alta Prioridade)

**Funcionalidades:**
- `show` - Mostra métricas atuais de performance
- `monitor` - Monitoramento em tempo real

**Implementação:**
- Similar ao `resources`, mas com foco em métricas de performance
- Integração com sistema de métricas do agent
- Suporte a alertas de threshold

**Estimativa:** 2-3 horas

---

### Prioridade 4: Maintenance (Média Prioridade)

**Funcionalidades:**
- `clean-logs` - Limpa logs antigos
- `optimize-db` - Otimiza banco de dados (VACUUM, ANALYZE)
- `cleanup` - Limpeza geral (temp files, cache)

**Implementação:**
```go
type MaintenanceManager interface {
    CleanLogs(olderThan time.Duration) (*CleanupResult, error)
    OptimizeDatabase(full bool) error
    Cleanup(options CleanupOptions) (*CleanupResult, error)
}

type CleanupResult struct {
    FilesRemoved   int
    SpaceFreed     uint64
    Duration       time.Duration
}
```

**Comandos:**
- Logs: rotação com logrotate ou implementação custom
- DB: `VACUUM`, `ANALYZE`, `REINDEX` para SQLite
- Cleanup: remoção de arquivos em `/tmp`, cache, etc.

**Estimativa:** 2-3 horas

---

### Prioridade 5: Backup & Config (Média Prioridade)

**Backup:**
- `create` - Cria backup (full/incremental)
- `restore` - Restaura do backup

**Config:**
- Validação e comparação de configs
- Export/import de configurações

**Estimativa:** 3-4 horas (ambos)

---

### Prioridade 6: Deployment & Security (Baixa Prioridade)

Mais complexos, deixar para depois das funcionalidades essenciais.

**Estimativa:** 6-8 horas (ambos)

---

## 📋 Roadmap de Implementação

### Fase 1 (Hoje) - Comandos Essenciais
- [x] Criar plano de implementação
- [ ] Implementar `resources` command
- [ ] Implementar `network` command
- [ ] Implementar `performance` command

**Tempo estimado:** 6-8 horas

### Fase 2 (Próxima) - Manutenção
- [ ] Implementar `maintenance` command
- [ ] Implementar `backup` command básico
- [ ] Implementar `config` validation

**Tempo estimado:** 5-6 horas

### Fase 3 (Futuro) - Avançados
- [ ] Implementar `deployment` command
- [ ] Implementar `security` command
- [ ] Integração com agents remotos via gRPC

**Tempo estimado:** 8-10 horas

---

## 🏗️ Arquitetura

### Estrutura de Diretórios
```
cmd/sloth-runner/commands/sysadmin/
├── resources/
│   ├── resources.go       # Comandos CLI
│   ├── monitor.go         # Implementação de monitoramento
│   ├── resources_test.go  # Testes
│   └── docs já implementados
├── network/
│   ├── network.go         # Comandos CLI
│   ├── diagnostics.go     # Implementação de diagnósticos
│   └── network_test.go    # Testes
└── ...
```

### Princípios
1. **Abstração**: Interfaces para permitir múltiplas implementações (Linux/macOS/Windows)
2. **Testabilidade**: Mocks para testes unitários
3. **Modularidade**: Cada comando em seu próprio package
4. **Consistência**: Seguir padrões já estabelecidos em packages/services

---

## 🔧 Comandos do Sistema Utilizados

### Linux
- CPU: `/proc/stat`, `top`, `mpstat`
- Memória: `/proc/meminfo`, `free`
- Disco: `df`, `/proc/diskstats`, `iostat`
- Rede: `/proc/net/dev`, `netstat`, `ss`
- Processos: `/proc/[pid]/stat`, `ps`

### macOS
- CPU: `sysctl -n hw.ncpu`, `sysctl -n vm.loadavg`
- Memória: `vm_stat`, `sysctl vm.swapusage`
- Disco: `df -h`
- Rede: `netstat -i`, `ifconfig`

---

## ✅ Critérios de Sucesso

Para cada comando implementado:
- [ ] Funciona localmente (Mac e Linux)
- [ ] Tem testes unitários (> 70% coverage)
- [ ] Tem documentação `docs` completa
- [ ] Output formatado com pterm
- [ ] Error handling robusto
- [ ] Pode ser executado via gRPC em agents remotos (futuro)

---

## 📚 Referências

- Comando `top`: https://man7.org/linux/man-pages/man1/top.1.html
- Proc filesystem: https://man7.org/linux/man-pages/man5/proc.5.html
- Go syscall package: https://pkg.go.dev/syscall
- Shirou/gopsutil: https://github.com/shirou/gopsutil (referência)

---

**Status:** 🚧 Pronto para começar implementação!
**Próximo passo:** Implementar `resources` command
