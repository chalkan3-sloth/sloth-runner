package luainterface

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// PkgModule provides functions for package management.
type PkgModule struct{}

// NewPkgModule creates a new PkgModule.
func NewPkgModule() *PkgModule {
	return &PkgModule{}
}

// Loader is the module loader function.
func (p *PkgModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), p.exports())
	L.Push(mod)
	return 1
}

func (p *PkgModule) exports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"install":        p.install,
		"remove":         p.remove,
		"update":         p.update,
		"upgrade":        p.upgrade,
		"search":         p.search,
		"info":           p.info,
		"list":           p.list,
		"is_installed":   p.isInstalled,
		"get_manager":    p.getManager,
		"clean":          p.clean,
		"autoremove":     p.autoremove,
		"which":          p.which,
		"version":        p.version,
		"deps":           p.deps,
		"install_local":  p.installLocal,
	}
}

// detectPackageManager detects the available package manager
func (p *PkgModule) detectPackageManager() (string, error) {
	managers := []string{"apt-get", "apt", "yum", "dnf", "pacman", "zypper", "brew"}
	
	for _, manager := range managers {
		if _, err := exec.LookPath(manager); err == nil {
			return manager, nil
		}
	}
	
	return "", fmt.Errorf("no supported package manager found")
}

// parsePackages converts Lua value to string slice of packages
func (p *PkgModule) parsePackages(L *lua.LState, arg int) []string {
	val := L.Get(arg)
	
	if val.Type() == lua.LTString {
		// Single package as string
		return []string{val.String()}
	} else if val.Type() == lua.LTTable {
		// Multiple packages as table
		var packages []string
		tbl := val.(*lua.LTable)
		tbl.ForEach(func(k, v lua.LValue) {
			if v.Type() == lua.LTString {
				packages = append(packages, v.String())
			}
		})
		return packages
	}
	
	return []string{val.String()}
}

// needsSudo checks if the command needs sudo
func (p *PkgModule) needsSudo(manager string) bool {
	// brew doesn't need sudo, others do (except on macOS running as root)
	if manager == "brew" {
		return false
	}
	
	// On macOS with other package managers (like MacPorts), may not need sudo
	if runtime.GOOS == "darwin" {
		return false
	}
	
	return true
}

// buildInstallCommand builds the install command based on package manager
func (p *PkgModule) buildInstallCommand(manager string, packages []string) []string {
	var args []string
	
	if p.needsSudo(manager) {
		args = append(args, "sudo")
	}
	
	switch manager {
	case "apt-get", "apt":
		args = append(args, manager, "install", "-y")
		args = append(args, packages...)
	case "yum", "dnf":
		args = append(args, manager, "install", "-y")
		args = append(args, packages...)
	case "pacman":
		args = append(args, manager, "-S", "--noconfirm")
		args = append(args, packages...)
	case "zypper":
		args = append(args, manager, "install", "-y")
		args = append(args, packages...)
	case "brew":
		args = append(args, manager, "install")
		args = append(args, packages...)
	default:
		args = append(args, manager, "install")
		args = append(args, packages...)
	}
	
	return args
}

// buildRemoveCommand builds the remove command based on package manager
func (p *PkgModule) buildRemoveCommand(manager string, packages []string) []string {
	var args []string
	
	if p.needsSudo(manager) {
		args = append(args, "sudo")
	}
	
	switch manager {
	case "apt-get", "apt":
		args = append(args, manager, "remove", "-y")
		args = append(args, packages...)
	case "yum", "dnf":
		args = append(args, manager, "remove", "-y")
		args = append(args, packages...)
	case "pacman":
		args = append(args, manager, "-R", "--noconfirm")
		args = append(args, packages...)
	case "zypper":
		args = append(args, manager, "remove", "-y")
		args = append(args, packages...)
	case "brew":
		args = append(args, manager, "uninstall")
		args = append(args, packages...)
	default:
		args = append(args, manager, "remove")
		args = append(args, packages...)
	}
	
	return args
}

// buildUpdateCommand builds the update command based on package manager
func (p *PkgModule) buildUpdateCommand(manager string) []string {
	var args []string
	
	if p.needsSudo(manager) {
		args = append(args, "sudo")
	}
	
	switch manager {
	case "apt-get", "apt":
		args = append(args, manager, "update")
	case "yum", "dnf":
		args = append(args, manager, "check-update")
	case "pacman":
		args = append(args, manager, "-Sy")
	case "zypper":
		args = append(args, manager, "refresh")
	case "brew":
		args = append(args, manager, "update")
	default:
		args = append(args, manager, "update")
	}
	
	return args
}

// buildUpgradeCommand builds the upgrade command based on package manager
func (p *PkgModule) buildUpgradeCommand(manager string) []string {
	var args []string
	
	if p.needsSudo(manager) {
		args = append(args, "sudo")
	}
	
	switch manager {
	case "apt-get", "apt":
		args = append(args, manager, "upgrade", "-y")
	case "yum", "dnf":
		args = append(args, manager, "upgrade", "-y")
	case "pacman":
		args = append(args, manager, "-Su", "--noconfirm")
	case "zypper":
		args = append(args, manager, "update", "-y")
	case "brew":
		args = append(args, manager, "upgrade")
	default:
		args = append(args, manager, "upgrade")
	}
	
	return args
}

func (p *PkgModule) install(L *lua.LState) int {
	packages := p.parsePackages(L, 1)
	
	if len(packages) == 0 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("No packages specified"))
		return 2
	}
	
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	args := p.buildInstallCommand(manager, packages)
	cmd := exec.Command(args[0], args[1:]...)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to install packages: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

func (p *PkgModule) remove(L *lua.LState) int {
	packages := p.parsePackages(L, 1)
	
	if len(packages) == 0 {
		L.Push(lua.LFalse)
		L.Push(lua.LString("No packages specified"))
		return 2
	}
	
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	args := p.buildRemoveCommand(manager, packages)
	cmd := exec.Command(args[0], args[1:]...)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to remove packages: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

func (p *PkgModule) update(L *lua.LState) int {
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	args := p.buildUpdateCommand(manager)
	cmd := exec.Command(args[0], args[1:]...)
	
	output, err := cmd.CombinedOutput()
	// yum check-update returns 100 if there are updates available
	if err != nil && (manager != "yum" && manager != "dnf") {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to update package list: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

func (p *PkgModule) upgrade(L *lua.LState) int {
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	args := p.buildUpgradeCommand(manager)
	cmd := exec.Command(args[0], args[1:]...)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to upgrade packages: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

func (p *PkgModule) search(L *lua.LState) int {
	query := L.ToString(1)
	
	if query == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("No search query specified"))
		return 2
	}
	
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	var args []string
	switch manager {
	case "apt-get", "apt":
		args = []string{manager, "search", query}
	case "yum", "dnf":
		args = []string{manager, "search", query}
	case "pacman":
		args = []string{manager, "-Ss", query}
	case "zypper":
		args = []string{manager, "search", query}
	case "brew":
		args = []string{manager, "search", query}
	default:
		args = []string{manager, "search", query}
	}
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to search packages: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

func (p *PkgModule) info(L *lua.LState) int {
	pkgName := L.ToString(1)
	
	if pkgName == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("No package name specified"))
		return 2
	}
	
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	var args []string
	switch manager {
	case "apt-get", "apt":
		args = []string{manager, "show", pkgName}
	case "yum", "dnf":
		args = []string{manager, "info", pkgName}
	case "pacman":
		args = []string{manager, "-Si", pkgName}
	case "zypper":
		args = []string{manager, "info", pkgName}
	case "brew":
		args = []string{manager, "info", pkgName}
	default:
		args = []string{manager, "info", pkgName}
	}
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to get package info: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

func (p *PkgModule) list(L *lua.LState) int {
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	var args []string
	switch manager {
	case "apt", "apt-get":
		args = []string{"dpkg", "--list"}
	case "yum", "dnf":
		args = []string{manager, "list", "installed"}
	case "pacman":
		args = []string{manager, "-Q"}
	case "zypper":
		args = []string{manager, "packages", "--installed-only"}
	case "brew":
		args = []string{manager, "list"}
	default:
		args = []string{manager, "list"}
	}
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to list packages: %s\n%s", err, string(output))))
		return 2
	}
	
	// Parse output into a table
	lines := strings.Split(string(output), "\n")
	tbl := L.NewTable()
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			tbl.Append(lua.LString(line))
		}
	}
	
	L.Push(lua.LTrue)
	L.Push(tbl)
	return 2
}

// isInstalled checks if a package is installed
func (p *PkgModule) isInstalled(L *lua.LState) int {
	pkgName := L.ToString(1)
	
	if pkgName == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("No package name specified"))
		return 2
	}
	
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	var cmd *exec.Cmd
	switch manager {
	case "apt", "apt-get":
		cmd = exec.Command("dpkg", "-l", pkgName)
	case "yum", "dnf":
		cmd = exec.Command(manager, "list", "installed", pkgName)
	case "pacman":
		cmd = exec.Command(manager, "-Q", pkgName)
	case "zypper":
		cmd = exec.Command(manager, "search", "--installed-only", pkgName)
	case "brew":
		cmd = exec.Command(manager, "list", pkgName)
	default:
		cmd = exec.Command(manager, "list", pkgName)
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString("Package not installed"))
		return 2
	}
	
	// Check if output contains package name
	if strings.Contains(string(output), pkgName) {
		L.Push(lua.LTrue)
		L.Push(lua.LString("Package is installed"))
		return 2
	}
	
	L.Push(lua.LFalse)
	L.Push(lua.LString("Package not installed"))
	return 2
}

// getManager returns the detected package manager
func (p *PkgModule) getManager(L *lua.LState) int {
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	L.Push(lua.LString(manager))
	L.Push(lua.LNil)
	return 2
}

// clean removes cached package files
func (p *PkgModule) clean(L *lua.LState) int {
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	var args []string
	if p.needsSudo(manager) {
		args = append(args, "sudo")
	}
	
	switch manager {
	case "apt", "apt-get":
		args = append(args, manager, "clean")
	case "yum", "dnf":
		args = append(args, manager, "clean", "all")
	case "pacman":
		args = append(args, manager, "-Sc", "--noconfirm")
	case "zypper":
		args = append(args, manager, "clean")
	case "brew":
		args = append(args, manager, "cleanup")
	default:
		L.Push(lua.LFalse)
		L.Push(lua.LString("Clean command not supported for " + manager))
		return 2
	}
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to clean: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// autoremove removes unused dependencies
func (p *PkgModule) autoremove(L *lua.LState) int {
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	var args []string
	if p.needsSudo(manager) {
		args = append(args, "sudo")
	}
	
	switch manager {
	case "apt", "apt-get":
		args = append(args, manager, "autoremove", "-y")
	case "yum", "dnf":
		args = append(args, manager, "autoremove", "-y")
	case "pacman":
		args = append(args, manager, "-Rns", "$(pacman -Qtdq)", "--noconfirm")
	case "brew":
		args = append(args, manager, "autoremove")
	default:
		L.Push(lua.LFalse)
		L.Push(lua.LString("Autoremove not supported for " + manager))
		return 2
	}
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	// Some package managers return non-zero if nothing to remove
	if err != nil && !strings.Contains(string(output), "Nothing to do") {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to autoremove: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// which finds the path of an executable
func (p *PkgModule) which(L *lua.LState) int {
	execName := L.ToString(1)
	
	if execName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("No executable name specified"))
		return 2
	}
	
	path, err := exec.LookPath(execName)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Executable not found: %s", execName)))
		return 2
	}
	
	L.Push(lua.LString(path))
	L.Push(lua.LNil)
	return 2
}

// version gets the version of an installed package
func (p *PkgModule) version(L *lua.LState) int {
	pkgName := L.ToString(1)
	
	if pkgName == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("No package name specified"))
		return 2
	}
	
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	var cmd *exec.Cmd
	switch manager {
	case "apt", "apt-get":
		cmd = exec.Command("dpkg", "-s", pkgName)
	case "yum", "dnf":
		cmd = exec.Command(manager, "info", pkgName)
	case "pacman":
		cmd = exec.Command(manager, "-Q", pkgName)
	case "zypper":
		cmd = exec.Command(manager, "info", pkgName)
	case "brew":
		cmd = exec.Command(manager, "info", pkgName, "--json")
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString("Version check not supported for " + manager))
		return 2
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("Failed to get version: %s", err)))
		return 2
	}
	
	// Try to extract version from output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "version") {
			L.Push(lua.LString(strings.TrimSpace(line)))
			L.Push(lua.LNil)
			return 2
		}
	}
	
	L.Push(lua.LString(string(output)))
	L.Push(lua.LNil)
	return 2
}

// deps shows dependencies of a package
func (p *PkgModule) deps(L *lua.LState) int {
	pkgName := L.ToString(1)
	
	if pkgName == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("No package name specified"))
		return 2
	}
	
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	var cmd *exec.Cmd
	switch manager {
	case "apt", "apt-get":
		cmd = exec.Command("apt-cache", "depends", pkgName)
	case "yum", "dnf":
		cmd = exec.Command(manager, "deplist", pkgName)
	case "pacman":
		cmd = exec.Command(manager, "-Si", pkgName)
	case "zypper":
		cmd = exec.Command(manager, "info", "--requires", pkgName)
	case "brew":
		cmd = exec.Command(manager, "deps", pkgName)
	default:
		L.Push(lua.LFalse)
		L.Push(lua.LString("Dependency listing not supported for " + manager))
		return 2
	}
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to get dependencies: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}

// installLocal installs a package from a local file
func (p *PkgModule) installLocal(L *lua.LState) int {
	filePath := L.ToString(1)
	
	if filePath == "" {
		L.Push(lua.LFalse)
		L.Push(lua.LString("No file path specified"))
		return 2
	}
	
	manager, err := p.detectPackageManager()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	
	var args []string
	if p.needsSudo(manager) {
		args = append(args, "sudo")
	}
	
	switch manager {
	case "apt", "apt-get":
		args = append(args, "dpkg", "-i", filePath)
	case "yum", "dnf":
		args = append(args, manager, "install", "-y", filePath)
	case "pacman":
		args = append(args, manager, "-U", "--noconfirm", filePath)
	case "zypper":
		args = append(args, manager, "install", "-y", filePath)
	case "brew":
		args = append(args, manager, "install", filePath)
	default:
		L.Push(lua.LFalse)
		L.Push(lua.LString("Local install not supported for " + manager))
		return 2
	}
	
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		L.Push(lua.LFalse)
		L.Push(lua.LString(fmt.Sprintf("Failed to install local package: %s\n%s", err, string(output))))
		return 2
	}
	
	L.Push(lua.LTrue)
	L.Push(lua.LString(string(output)))
	return 2
}
