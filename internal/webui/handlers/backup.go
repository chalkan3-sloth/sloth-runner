package handlers

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// BackupHandler handles backup and restore operations
type BackupHandler struct{}

// NewBackupHandler creates a new backup handler
func NewBackupHandler() *BackupHandler {
	return &BackupHandler{}
}

// CreateBackup creates a backup of all databases
func (h *BackupHandler) CreateBackup(c *gin.Context) {
	timestamp := time.Now().Format("20060102-150405")
	backupName := fmt.Sprintf("sloth-runner-backup-%s.tar.gz", timestamp)

	// Create temporary directory for backup
	tmpDir := filepath.Join(os.TempDir(), "sloth-backup-"+timestamp)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer os.RemoveAll(tmpDir)

	// Paths to backup
	backupPaths := []struct {
		source string
		dest   string
	}{
		{".sloth-cache/agents.db", "agents.db"},
		{".sloth-cache/hooks.db", "hooks.db"},
		{"/etc/sloth-runner/sloths.db", "sloths.db"},
	}

	homeDir, _ := os.UserHomeDir()
	backupPaths = append(backupPaths,
		struct{ source, dest string }{filepath.Join(homeDir, ".sloth-runner/secrets.db"), "secrets.db"},
		struct{ source, dest string }{filepath.Join(homeDir, ".sloth-runner/ssh_profiles.db"), "ssh_profiles.db"},
	)

	// Create tar.gz file
	backupPath := filepath.Join(tmpDir, backupName)
	outFile, err := os.Create(backupPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer outFile.Close()

	gzw := gzip.NewWriter(outFile)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// Add files to archive
	for _, bp := range backupPaths {
		if _, err := os.Stat(bp.source); os.IsNotExist(err) {
			continue
		}

		if err := h.addFileToTar(tw, bp.source, bp.dest); err != nil {
			fmt.Printf("Warning: failed to add %s to backup: %v\n", bp.source, err)
		}
	}

	// Send file as download
	c.FileAttachment(backupPath, backupName)
}

// RestoreBackup restores from a backup file
func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	file, err := c.FormFile("backup")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No backup file provided"})
		return
	}

	// Create temporary directory
	tmpDir := filepath.Join(os.TempDir(), "sloth-restore-"+time.Now().Format("20060102-150405"))
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer os.RemoveAll(tmpDir)

	// Save uploaded file
	uploadPath := filepath.Join(tmpDir, file.Filename)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Extract backup
	if err := h.extractTarGz(uploadPath, tmpDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Restore files
	restoredFiles := make([]string, 0)

	restorePaths := map[string]string{
		"agents.db":        ".sloth-cache/agents.db",
		"hooks.db":         ".sloth-cache/hooks.db",
		"sloths.db":        "/etc/sloth-runner/sloths.db",
		"secrets.db":       filepath.Join(os.Getenv("HOME"), ".sloth-runner/secrets.db"),
		"ssh_profiles.db":  filepath.Join(os.Getenv("HOME"), ".sloth-runner/ssh_profiles.db"),
	}

	for source, dest := range restorePaths {
		sourcePath := filepath.Join(tmpDir, source)
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			continue
		}

		// Ensure destination directory exists
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			continue
		}

		// Copy file
		if err := h.copyFile(sourcePath, dest); err == nil {
			restoredFiles = append(restoredFiles, dest)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Backup restored successfully",
		"files":   restoredFiles,
	})
}

// ListBackups lists available backups
func (h *BackupHandler) ListBackups(c *gin.Context) {
	// TODO: Implement backup storage and listing
	c.JSON(http.StatusOK, gin.H{
		"backups": []gin.H{},
		"message": "Backup listing not yet implemented",
	})
}

// addFileToTar adds a file to tar archive
func (h *BackupHandler) addFileToTar(tw *tar.Writer, source, dest string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	header := &tar.Header{
		Name:    dest,
		Mode:    int64(stat.Mode()),
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(tw, file)
	return err
}

// extractTarGz extracts a tar.gz file
func (h *BackupHandler) extractTarGz(source, dest string) error {
	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// copyFile copies a file
func (h *BackupHandler) copyFile(source, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
