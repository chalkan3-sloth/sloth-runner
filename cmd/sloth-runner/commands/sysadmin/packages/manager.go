package packages

import (
	"fmt"
	"os/exec"
	"strings"
)

// PackageManagerType representa o tipo de gerenciador de pacotes
type PackageManagerType string

const (
	APT    PackageManagerType = "apt"
	YUM    PackageManagerType = "yum"
	DNF    PackageManagerType = "dnf"
	PACMAN PackageManagerType = "pacman"
	APK    PackageManagerType = "apk"
	ZYPPER PackageManagerType = "zypper"
	NONE   PackageManagerType = "none"
)

// PackageManager interface para operações de gerenciamento de pacotes
type PackageManager interface {
	List() ([]Package, error)
	Search(query string) ([]Package, error)
	Install(packageName string) error
	Remove(packageName string) error
	Update() error
	Upgrade() error
	Info(packageName string) (*PackageInfo, error)
}

// Package representa um pacote instalado ou disponível
type Package struct {
	Name        string
	Version     string
	Description string
	Installed   bool
	Repository  string
}

// PackageInfo contém informações detalhadas sobre um pacote
type PackageInfo struct {
	Name         string
	Version      string
	Description  string
	Architecture string
	Size         string
	Dependencies []string
	Repository   string
	Installed    bool
}

// DetectPackageManager detecta qual gerenciador de pacotes está disponível no sistema
func DetectPackageManager() PackageManagerType {
	managers := []struct {
		cmd  string
		args []string
		pm   PackageManagerType
	}{
		{"apt", []string{"--version"}, APT},
		{"dnf", []string{"--version"}, DNF},
		{"yum", []string{"--version"}, YUM},
		{"pacman", []string{"--version"}, PACMAN},
		{"apk", []string{"--version"}, APK},
		{"zypper", []string{"--version"}, ZYPPER},
	}

	for _, mgr := range managers {
		if commandExists(mgr.cmd) {
			return mgr.pm
		}
	}

	return NONE
}

// commandExists verifica se um comando existe no PATH
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// GetPackageManager retorna uma instância do gerenciador de pacotes apropriado
func GetPackageManager() (PackageManager, error) {
	pmType := DetectPackageManager()

	switch pmType {
	case APT:
		return &AptManager{}, nil
	case YUM:
		return &YumManager{}, nil
	case DNF:
		return &DnfManager{}, nil
	case PACMAN:
		return &PacmanManager{}, nil
	case APK:
		return &ApkManager{}, nil
	case ZYPPER:
		return &ZypperManager{}, nil
	default:
		return nil, fmt.Errorf("no supported package manager found")
	}
}

// AptManager implementa PackageManager para apt/apt-get
type AptManager struct{}

func (a *AptManager) List() ([]Package, error) {
	cmd := exec.Command("dpkg", "--get-selections")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list packages: %w", err)
	}

	var packages []Package
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "install" {
			// Get version info
			vCmd := exec.Command("dpkg-query", "-W", "-f=${Version}", fields[0])
			vOut, _ := vCmd.Output()
			version := strings.TrimSpace(string(vOut))

			packages = append(packages, Package{
				Name:      fields[0],
				Version:   version,
				Installed: true,
			})
		}
	}

	return packages, nil
}

func (a *AptManager) Search(query string) ([]Package, error) {
	cmd := exec.Command("apt-cache", "search", query)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to search packages: %w", err)
	}

	var packages []Package
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " - ", 2)
		if len(parts) == 2 {
			packages = append(packages, Package{
				Name:        parts[0],
				Description: parts[1],
				Installed:   false,
			})
		}
	}

	return packages, nil
}

func (a *AptManager) Install(packageName string) error {
	cmd := exec.Command("apt-get", "install", "-y", packageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install %s: %w\nOutput: %s", packageName, err, string(output))
	}
	return nil
}

func (a *AptManager) Remove(packageName string) error {
	cmd := exec.Command("apt-get", "remove", "-y", packageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove %s: %w\nOutput: %s", packageName, err, string(output))
	}
	return nil
}

func (a *AptManager) Update() error {
	cmd := exec.Command("apt-get", "update")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update package lists: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func (a *AptManager) Upgrade() error {
	cmd := exec.Command("apt-get", "upgrade", "-y")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to upgrade packages: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func (a *AptManager) Info(packageName string) (*PackageInfo, error) {
	cmd := exec.Command("apt-cache", "show", packageName)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get package info: %w", err)
	}

	info := &PackageInfo{Name: packageName}
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Version: ") {
			info.Version = strings.TrimPrefix(line, "Version: ")
		} else if strings.HasPrefix(line, "Description: ") {
			info.Description = strings.TrimPrefix(line, "Description: ")
		} else if strings.HasPrefix(line, "Architecture: ") {
			info.Architecture = strings.TrimPrefix(line, "Architecture: ")
		} else if strings.HasPrefix(line, "Installed-Size: ") {
			info.Size = strings.TrimPrefix(line, "Installed-Size: ")
		}
	}

	return info, nil
}

// Stubs for other package managers (YUM, DNF, PACMAN, APK, ZYPPER)
// These will return not implemented errors for now

type YumManager struct{}

func (y *YumManager) List() ([]Package, error) {
	return nil, fmt.Errorf("yum support not yet implemented")
}
func (y *YumManager) Search(query string) ([]Package, error) {
	return nil, fmt.Errorf("yum support not yet implemented")
}
func (y *YumManager) Install(packageName string) error {
	return fmt.Errorf("yum support not yet implemented")
}
func (y *YumManager) Remove(packageName string) error {
	return fmt.Errorf("yum support not yet implemented")
}
func (y *YumManager) Update() error {
	return fmt.Errorf("yum support not yet implemented")
}
func (y *YumManager) Upgrade() error {
	return fmt.Errorf("yum support not yet implemented")
}
func (y *YumManager) Info(packageName string) (*PackageInfo, error) {
	return nil, fmt.Errorf("yum support not yet implemented")
}

type DnfManager struct{}

func (d *DnfManager) List() ([]Package, error) {
	return nil, fmt.Errorf("dnf support not yet implemented")
}
func (d *DnfManager) Search(query string) ([]Package, error) {
	return nil, fmt.Errorf("dnf support not yet implemented")
}
func (d *DnfManager) Install(packageName string) error {
	return fmt.Errorf("dnf support not yet implemented")
}
func (d *DnfManager) Remove(packageName string) error {
	return fmt.Errorf("dnf support not yet implemented")
}
func (d *DnfManager) Update() error {
	return fmt.Errorf("dnf support not yet implemented")
}
func (d *DnfManager) Upgrade() error {
	return fmt.Errorf("dnf support not yet implemented")
}
func (d *DnfManager) Info(packageName string) (*PackageInfo, error) {
	return nil, fmt.Errorf("dnf support not yet implemented")
}

type PacmanManager struct{}

func (p *PacmanManager) List() ([]Package, error) {
	return nil, fmt.Errorf("pacman support not yet implemented")
}
func (p *PacmanManager) Search(query string) ([]Package, error) {
	return nil, fmt.Errorf("pacman support not yet implemented")
}
func (p *PacmanManager) Install(packageName string) error {
	return fmt.Errorf("pacman support not yet implemented")
}
func (p *PacmanManager) Remove(packageName string) error {
	return fmt.Errorf("pacman support not yet implemented")
}
func (p *PacmanManager) Update() error {
	return fmt.Errorf("pacman support not yet implemented")
}
func (p *PacmanManager) Upgrade() error {
	return fmt.Errorf("pacman support not yet implemented")
}
func (p *PacmanManager) Info(packageName string) (*PackageInfo, error) {
	return nil, fmt.Errorf("pacman support not yet implemented")
}

type ApkManager struct{}

func (a *ApkManager) List() ([]Package, error) {
	return nil, fmt.Errorf("apk support not yet implemented")
}
func (a *ApkManager) Search(query string) ([]Package, error) {
	return nil, fmt.Errorf("apk support not yet implemented")
}
func (a *ApkManager) Install(packageName string) error {
	return fmt.Errorf("apk support not yet implemented")
}
func (a *ApkManager) Remove(packageName string) error {
	return fmt.Errorf("apk support not yet implemented")
}
func (a *ApkManager) Update() error {
	return fmt.Errorf("apk support not yet implemented")
}
func (a *ApkManager) Upgrade() error {
	return fmt.Errorf("apk support not yet implemented")
}
func (a *ApkManager) Info(packageName string) (*PackageInfo, error) {
	return nil, fmt.Errorf("apk support not yet implemented")
}

type ZypperManager struct{}

func (z *ZypperManager) List() ([]Package, error) {
	return nil, fmt.Errorf("zypper support not yet implemented")
}
func (z *ZypperManager) Search(query string) ([]Package, error) {
	return nil, fmt.Errorf("zypper support not yet implemented")
}
func (z *ZypperManager) Install(packageName string) error {
	return fmt.Errorf("zypper support not yet implemented")
}
func (z *ZypperManager) Remove(packageName string) error {
	return fmt.Errorf("zypper support not yet implemented")
}
func (z *ZypperManager) Update() error {
	return fmt.Errorf("zypper support not yet implemented")
}
func (z *ZypperManager) Upgrade() error {
	return fmt.Errorf("zypper support not yet implemented")
}
func (z *ZypperManager) Info(packageName string) (*PackageInfo, error) {
	return nil, fmt.Errorf("zypper support not yet implemented")
}
