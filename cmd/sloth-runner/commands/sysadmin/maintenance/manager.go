package maintenance

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// MaintenanceManager interface para tarefas de manutenção
type MaintenanceManager interface {
	CleanLogs(olderThan time.Duration, dryRun bool) (*CleanupResult, error)
	OptimizeDatabase(dbPath string, full bool) error
	Cleanup(options CleanupOptions) (*CleanupResult, error)
}

// CleanupResult contém resultados de uma operação de limpeza
type CleanupResult struct {
	FilesRemoved   int
	FilesCompressed int
	SpaceFreed     uint64
	Duration       time.Duration
	DryRun         bool
	Details        []string
}

// CleanupOptions opções para limpeza
type CleanupOptions struct {
	TempFiles    bool
	Cache        bool
	OldLogs      bool
	LogAge       time.Duration
	DryRun       bool
	Paths        []string
}

// SystemMaintenance implementação padrão de MaintenanceManager
type SystemMaintenance struct{}

// NewMaintenanceManager cria um novo manager de manutenção
func NewMaintenanceManager() MaintenanceManager {
	return &SystemMaintenance{}
}

// CleanLogs limpa logs antigos
func (m *SystemMaintenance) CleanLogs(olderThan time.Duration, dryRun bool) (*CleanupResult, error) {
	result := &CleanupResult{
		DryRun:  dryRun,
		Details: []string{},
	}

	start := time.Now()
	defer func() {
		result.Duration = time.Since(start)
	}()

	// Diretórios de log para verificar
	logDirs := []string{
		"/var/log",
		"/tmp",
		filepath.Join(os.Getenv("HOME"), ".local", "share", "sloth-runner", "logs"),
		filepath.Join(os.TempDir(), "sloth-runner"),
	}

	cutoffTime := time.Now().Add(-olderThan)

	for _, dir := range logDirs {
		// Verifica se o diretório existe e é acessível
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		// Procura por arquivos de log
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// Ignora erros de permissão
				return nil
			}

			// Pula diretórios
			if info.IsDir() {
				return nil
			}

			// Verifica se é um arquivo de log
			if !isLogFile(path) {
				return nil
			}

			// Verifica se é mais antigo que o cutoff
			if info.ModTime().Before(cutoffTime) {
				size := uint64(info.Size())

				if !dryRun {
					if err := os.Remove(path); err != nil {
						result.Details = append(result.Details, fmt.Sprintf("Failed to remove %s: %v", path, err))
						return nil
					}
				}

				result.FilesRemoved++
				result.SpaceFreed += size
				result.Details = append(result.Details, fmt.Sprintf("Removed: %s (%.2f MB)", path, float64(size)/(1024*1024)))
			}

			return nil
		})

		if err != nil {
			result.Details = append(result.Details, fmt.Sprintf("Error scanning %s: %v", dir, err))
		}
	}

	return result, nil
}

// OptimizeDatabase otimiza um banco de dados SQLite
func (m *SystemMaintenance) OptimizeDatabase(dbPath string, full bool) error {
	// Verifica se o arquivo existe
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("database file not found: %s", dbPath)
	}

	// Abre conexão com o banco
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Executa VACUUM
	if _, err := db.Exec("VACUUM"); err != nil {
		return fmt.Errorf("failed to VACUUM: %w", err)
	}

	// Executa ANALYZE
	if _, err := db.Exec("ANALYZE"); err != nil {
		return fmt.Errorf("failed to ANALYZE: %w", err)
	}

	// Se full, também executa outras otimizações
	if full {
		// Reindex
		if _, err := db.Exec("REINDEX"); err != nil {
			return fmt.Errorf("failed to REINDEX: %w", err)
		}

		// Optimize pragma
		if _, err := db.Exec("PRAGMA optimize"); err != nil {
			return fmt.Errorf("failed to optimize: %w", err)
		}
	}

	return nil
}

// Cleanup executa limpeza geral
func (m *SystemMaintenance) Cleanup(options CleanupOptions) (*CleanupResult, error) {
	result := &CleanupResult{
		DryRun:  options.DryRun,
		Details: []string{},
	}

	start := time.Now()
	defer func() {
		result.Duration = time.Since(start)
	}()

	// Limpa logs antigos se solicitado
	if options.OldLogs {
		logResult, err := m.CleanLogs(options.LogAge, options.DryRun)
		if err != nil {
			return nil, fmt.Errorf("failed to clean logs: %w", err)
		}
		result.FilesRemoved += logResult.FilesRemoved
		result.SpaceFreed += logResult.SpaceFreed
		result.Details = append(result.Details, logResult.Details...)
	}

	// Limpa arquivos temporários se solicitado
	if options.TempFiles {
		tempResult := m.cleanTempFiles(options.DryRun)
		result.FilesRemoved += tempResult.FilesRemoved
		result.SpaceFreed += tempResult.SpaceFreed
		result.Details = append(result.Details, tempResult.Details...)
	}

	// Limpa cache se solicitado
	if options.Cache {
		cacheResult := m.cleanCache(options.DryRun)
		result.FilesRemoved += cacheResult.FilesRemoved
		result.SpaceFreed += cacheResult.SpaceFreed
		result.Details = append(result.Details, cacheResult.Details...)
	}

	// Limpa paths específicos se fornecidos
	for _, path := range options.Paths {
		pathResult := m.cleanPath(path, options.DryRun)
		result.FilesRemoved += pathResult.FilesRemoved
		result.SpaceFreed += pathResult.SpaceFreed
		result.Details = append(result.Details, pathResult.Details...)
	}

	return result, nil
}

// cleanTempFiles limpa arquivos temporários
func (m *SystemMaintenance) cleanTempFiles(dryRun bool) *CleanupResult {
	result := &CleanupResult{
		DryRun:  dryRun,
		Details: []string{},
	}

	tempDirs := []string{
		os.TempDir(),
		filepath.Join(os.Getenv("HOME"), ".cache", "sloth-runner"),
	}

	cutoffTime := time.Now().Add(-24 * time.Hour) // Arquivos com mais de 24h

	for _, dir := range tempDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			// Remove arquivos temporários antigos do sloth-runner
			if filepath.Base(path) != "sloth-runner" &&
			   !filepath.HasPrefix(filepath.Base(path), "sloth-runner-") {
				return nil
			}

			if info.ModTime().Before(cutoffTime) {
				size := uint64(info.Size())

				if !dryRun {
					if err := os.Remove(path); err != nil {
						return nil
					}
				}

				result.FilesRemoved++
				result.SpaceFreed += size
				result.Details = append(result.Details, fmt.Sprintf("Removed temp: %s", path))
			}

			return nil
		})
	}

	return result
}

// cleanCache limpa arquivos de cache
func (m *SystemMaintenance) cleanCache(dryRun bool) *CleanupResult {
	result := &CleanupResult{
		DryRun:  dryRun,
		Details: []string{},
	}

	cacheDirs := []string{
		filepath.Join(os.Getenv("HOME"), ".cache", "sloth-runner"),
		filepath.Join(os.Getenv("HOME"), ".local", "share", "sloth-runner", "cache"),
	}

	for _, dir := range cacheDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() && path != dir {
				// Pula, mas não para a recursão
				return nil
			}

			if !info.IsDir() {
				size := uint64(info.Size())

				if !dryRun {
					if err := os.Remove(path); err != nil {
						return nil
					}
				}

				result.FilesRemoved++
				result.SpaceFreed += size
				result.Details = append(result.Details, fmt.Sprintf("Removed cache: %s", path))
			}

			return nil
		})
	}

	return result
}

// cleanPath limpa um path específico
func (m *SystemMaintenance) cleanPath(path string, dryRun bool) *CleanupResult {
	result := &CleanupResult{
		DryRun:  dryRun,
		Details: []string{},
	}

	info, err := os.Stat(path)
	if err != nil {
		result.Details = append(result.Details, fmt.Sprintf("Error accessing %s: %v", path, err))
		return result
	}

	if info.IsDir() {
		filepath.Walk(path, func(p string, i os.FileInfo, err error) error {
			if err != nil || i.IsDir() {
				return nil
			}

			size := uint64(i.Size())

			if !dryRun {
				if err := os.Remove(p); err != nil {
					return nil
				}
			}

			result.FilesRemoved++
			result.SpaceFreed += size
			result.Details = append(result.Details, fmt.Sprintf("Removed: %s", p))
			return nil
		})
	} else {
		size := uint64(info.Size())

		if !dryRun {
			if err := os.Remove(path); err != nil {
				result.Details = append(result.Details, fmt.Sprintf("Error removing %s: %v", path, err))
				return result
			}
		}

		result.FilesRemoved = 1
		result.SpaceFreed = size
		result.Details = append(result.Details, fmt.Sprintf("Removed: %s", path))
	}

	return result
}

// isLogFile verifica se um arquivo é um arquivo de log
func isLogFile(path string) bool {
	ext := filepath.Ext(path)
	base := filepath.Base(path)

	// Extensões de log comuns
	logExts := []string{".log", ".log.gz", ".log.bz2", ".log.xz", ".out", ".err"}
	for _, logExt := range logExts {
		if ext == logExt {
			return true
		}
	}

	// Arquivos com "log" no nome
	if filepath.Ext(base) == "" && len(base) > 3 {
		return base[:3] == "log" || base[len(base)-3:] == "log"
	}

	return false
}

// FormatBytes formata bytes para formato legível
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// GetDatabaseSize retorna o tamanho de um banco de dados
func GetDatabaseSize(dbPath string) (uint64, error) {
	info, err := os.Stat(dbPath)
	if err != nil {
		return 0, err
	}
	return uint64(info.Size()), nil
}
