package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// BackupManager interface para operações de backup
type BackupManager interface {
	CreateBackup(options BackupOptions) (*BackupInfo, error)
	RestoreBackup(options RestoreOptions) (*RestoreInfo, error)
	ListBackups(backupDir string) ([]*BackupInfo, error)
}

// BackupOptions opções para criar backup
type BackupOptions struct {
	OutputPath  string
	Include     []string // Paths para incluir no backup
	Exclude     []string // Padrões para excluir
	Compress    bool
	Description string
}

// RestoreOptions opções para restaurar backup
type RestoreOptions struct {
	InputPath   string
	TargetDir   string
	DatabaseOnly bool
	ConfigOnly   bool
	DryRun       bool
}

// BackupInfo informações sobre um backup
type BackupInfo struct {
	Path        string
	Size        uint64
	Created     time.Time
	Description string
	FileCount   int
	Compressed  bool
}

// RestoreInfo informações sobre uma restauração
type RestoreInfo struct {
	FilesRestored int
	BytesRestored uint64
	Duration      time.Duration
	Errors        []string
}

// SystemBackup implementação padrão de BackupManager
type SystemBackup struct{}

// NewBackupManager cria um novo manager de backup
func NewBackupManager() BackupManager {
	return &SystemBackup{}
}

// CreateBackup cria um novo backup
func (b *SystemBackup) CreateBackup(options BackupOptions) (*BackupInfo, error) {
	info := &BackupInfo{
		Path:        options.OutputPath,
		Created:     time.Now(),
		Description: options.Description,
		Compressed:  options.Compress,
	}

	// Cria o arquivo de backup
	file, err := os.Create(options.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	var writer io.Writer = file

	// Se compressão habilitada, adiciona gzip writer
	if options.Compress {
		gzWriter := gzip.NewWriter(file)
		defer gzWriter.Close()
		writer = gzWriter
	}

	// Cria tar writer
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	// Adiciona arquivos ao backup
	for _, includePath := range options.Include {
		if err := addToArchive(tarWriter, includePath, info); err != nil {
			return nil, fmt.Errorf("failed to add %s to backup: %w", includePath, err)
		}
	}

	// Obtém tamanho final do backup
	fileInfo, err := os.Stat(options.OutputPath)
	if err == nil {
		info.Size = uint64(fileInfo.Size())
	}

	return info, nil
}

// RestoreBackup restaura um backup
func (b *SystemBackup) RestoreBackup(options RestoreOptions) (*RestoreInfo, error) {
	info := &RestoreInfo{
		Errors: []string{},
	}

	start := time.Now()
	defer func() {
		info.Duration = time.Since(start)
	}()

	// Abre o arquivo de backup
	file, err := os.Open(options.InputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open backup file: %w", err)
	}
	defer file.Close()

	var reader io.Reader = file

	// Detecta se é comprimido verificando extensão
	if filepath.Ext(options.InputPath) == ".gz" ||
	   filepath.Ext(filepath.Base(options.InputPath[:len(options.InputPath)-len(filepath.Ext(options.InputPath))])) == ".tar" {
		gzReader, err := gzip.NewReader(file)
		if err == nil {
			defer gzReader.Close()
			reader = gzReader
		}
	}

	// Cria tar reader
	tarReader := tar.NewReader(reader)

	// Extrai arquivos
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			info.Errors = append(info.Errors, fmt.Sprintf("Error reading tar: %v", err))
			continue
		}

		// Filtra por tipo se especificado
		if options.DatabaseOnly && !isDatabase(header.Name) {
			continue
		}
		if options.ConfigOnly && !isConfig(header.Name) {
			continue
		}

		// Determina path de destino
		targetPath := filepath.Join(options.TargetDir, header.Name)

		// Cria diretórios necessários
		if header.Typeflag == tar.TypeDir {
			if !options.DryRun {
				if err := os.MkdirAll(targetPath, 0755); err != nil {
					info.Errors = append(info.Errors, fmt.Sprintf("Failed to create directory %s: %v", targetPath, err))
				}
			}
			continue
		}

		// Cria arquivo
		if !options.DryRun {
			// Cria diretório pai se necessário
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				info.Errors = append(info.Errors, fmt.Sprintf("Failed to create parent directory for %s: %v", targetPath, err))
				continue
			}

			outFile, err := os.Create(targetPath)
			if err != nil {
				info.Errors = append(info.Errors, fmt.Sprintf("Failed to create file %s: %v", targetPath, err))
				continue
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				info.Errors = append(info.Errors, fmt.Sprintf("Failed to write file %s: %v", targetPath, err))
				continue
			}

			outFile.Close()

			// Restaura permissões
			if err := os.Chmod(targetPath, os.FileMode(header.Mode)); err != nil {
				info.Errors = append(info.Errors, fmt.Sprintf("Failed to set permissions on %s: %v", targetPath, err))
			}
		}

		info.FilesRestored++
		info.BytesRestored += uint64(header.Size)
	}

	return info, nil
}

// ListBackups lista backups em um diretório
func (b *SystemBackup) ListBackups(backupDir string) ([]*BackupInfo, error) {
	var backups []*BackupInfo

	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Verifica se é um arquivo de backup (*.tar.gz ou *.tar)
		name := entry.Name()
		if !isBackupFile(name) {
			continue
		}

		path := filepath.Join(backupDir, name)
		fileInfo, err := entry.Info()
		if err != nil {
			continue
		}

		backup := &BackupInfo{
			Path:       path,
			Size:       uint64(fileInfo.Size()),
			Created:    fileInfo.ModTime(),
			Compressed: filepath.Ext(name) == ".gz",
		}

		backups = append(backups, backup)
	}

	return backups, nil
}

// addToArchive adiciona um arquivo ou diretório ao archive
func addToArchive(tw *tar.Writer, source string, info *BackupInfo) error {
	return filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Cria header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// Usa caminho relativo
		header.Name = file

		// Escreve header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// Se não é um arquivo regular, pula
		if !fi.Mode().IsRegular() {
			return nil
		}

		// Abre e copia arquivo
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		info.FileCount++
		return nil
	})
}

// isDatabase verifica se um path é um arquivo de banco de dados
func isDatabase(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".db" || ext == ".sqlite" || ext == ".sqlite3"
}

// isConfig verifica se um path é um arquivo de configuração
func isConfig(path string) bool {
	base := filepath.Base(path)
	ext := filepath.Ext(path)

	// Arquivos de configuração comuns
	configFiles := []string{"config.yaml", "config.yml", "config.json", "config.toml", ".env"}
	for _, cf := range configFiles {
		if base == cf {
			return true
		}
	}

	// Extensões de configuração
	configExts := []string{".yaml", ".yml", ".json", ".toml", ".conf", ".ini"}
	for _, ce := range configExts {
		if ext == ce {
			return true
		}
	}

	return false
}

// isBackupFile verifica se um arquivo é um backup
func isBackupFile(name string) bool {
	return filepath.Ext(name) == ".tar" ||
	       (filepath.Ext(name) == ".gz" && filepath.Ext(name[:len(name)-3]) == ".tar")
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

// GetDefaultBackupPaths retorna paths padrão para backup
func GetDefaultBackupPaths() []string {
	home := os.Getenv("HOME")
	return []string{
		"/etc/sloth-runner",
		filepath.Join(home, ".config", "sloth-runner"),
	}
}

// GetDefaultBackupDir retorna o diretório padrão para backups
func GetDefaultBackupDir() string {
	return "/etc/sloth-runner/backups"
}
