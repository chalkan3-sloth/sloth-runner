package luainterface

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chalkan3-sloth/sloth-runner/internal/state"
	lua "github.com/yuin/gopher-lua"
)

// StowModule provides functions for GNU Stow package management (dotfiles management).
type StowModule struct {
	sm *state.StateManager
}

// NewStowModule creates a new StowModule.
func NewStowModule(sm *state.StateManager) *StowModule {
	return &StowModule{sm: sm}
}

// Loader is the module loader function.
func (s *StowModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), s.exports())
	L.Push(mod)
	return 1
}

func (s *StowModule) exports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"stow":        s.stow,
		"unstow":      s.unstow,
		"restow":      s.restow,
		"adopt":       s.adopt,
		"check":       s.check,
		"simulate":    s.simulate,
		"get_links":   s.getLinks,
		"is_stowed":   s.isStowed,
		"list_packages": s.listPackages,
		"verify":      s.verify,
	}
}

type stowConfig struct {
	Packages    []string
	Target      string
	StowDir     string
	Verbose     bool
	NoFolding   bool
	IgnorePatterns []string
	Override    []string
	Defer       []string
	DelegatedTo string
}

func (s *StowModule) parseConfig(L *lua.LState) (*stowConfig, error) {
	config := &stowConfig{
		Target:  os.Getenv("HOME"),
		StowDir: filepath.Join(os.Getenv("HOME"), ".dotfiles"),
	}

	if L.GetTop() == 0 {
		return nil, fmt.Errorf("stow configuration required")
	}

	tbl := L.CheckTable(1)
	
	tbl.ForEach(func(key, value lua.LValue) {
		switch key.String() {
		case "package":
			config.Packages = []string{value.String()}
		case "packages":
			if tbl, ok := value.(*lua.LTable); ok {
				tbl.ForEach(func(_, v lua.LValue) {
					config.Packages = append(config.Packages, v.String())
				})
			}
		case "target":
			config.Target = value.String()
		case "dir", "stow_dir":
			config.StowDir = value.String()
		case "verbose":
			if b, ok := value.(lua.LBool); ok {
				config.Verbose = bool(b)
			}
		case "no_folding":
			if b, ok := value.(lua.LBool); ok {
				config.NoFolding = bool(b)
			}
		case "ignore":
			if tbl, ok := value.(*lua.LTable); ok {
				tbl.ForEach(func(_, v lua.LValue) {
					config.IgnorePatterns = append(config.IgnorePatterns, v.String())
				})
			}
		case "override":
			if tbl, ok := value.(*lua.LTable); ok {
				tbl.ForEach(func(_, v lua.LValue) {
					config.Override = append(config.Override, v.String())
				})
			}
		case "defer":
			if tbl, ok := value.(*lua.LTable); ok {
				tbl.ForEach(func(_, v lua.LValue) {
					config.Defer = append(config.Defer, v.String())
				})
			}
		case "delegate_to":
			config.DelegatedTo = value.String()
		}
	})

	if len(config.Packages) == 0 {
		return nil, fmt.Errorf("package name or packages list is required")
	}

	return config, nil
}

func (s *StowModule) buildStowArgs(config *stowConfig, action string, pkg string) []string {
	args := []string{action}

	if config.Target != "" && config.Target != os.Getenv("HOME") {
		args = append(args, "-t", config.Target)
	}

	if config.StowDir != "" {
		args = append(args, "-d", config.StowDir)
	}

	if config.Verbose {
		args = append(args, "-v")
	}

	if config.NoFolding {
		args = append(args, "--no-folding")
	}

	for _, pattern := range config.IgnorePatterns {
		args = append(args, "--ignore", pattern)
	}

	for _, override := range config.Override {
		args = append(args, "--override", override)
	}

	for _, deferPattern := range config.Defer {
		args = append(args, "--defer", deferPattern)
	}

	args = append(args, pkg)
	return args
}

func (s *StowModule) calculateStateHash(config *stowConfig, action string, pkg string) string {
	data := fmt.Sprintf("%s:%s:%s:%s", action, pkg, config.StowDir, config.Target)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (s *StowModule) getCurrentState(config *stowConfig, pkg string) (map[string]interface{}, error) {
	state := make(map[string]interface{})
	
	// Check if package directory exists
	pkgPath := filepath.Join(config.StowDir, pkg)
	if _, err := os.Stat(pkgPath); err != nil {
		state["package_exists"] = false
		return state, nil
	}
	state["package_exists"] = true

	// Get list of files in package
	files := []string{}
	err := filepath.Walk(pkgPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(pkgPath, path)
			files = append(files, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	state["files"] = files

	// Check for stowed symlinks
	links := []string{}
	for _, file := range files {
		targetPath := filepath.Join(config.Target, file)
		if info, err := os.Lstat(targetPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
			if link, err := os.Readlink(targetPath); err == nil {
				expectedLink := filepath.Join(config.StowDir, pkg, file)
				if link == expectedLink || filepath.Clean(link) == filepath.Clean(expectedLink) {
					links = append(links, file)
				}
			}
		}
	}
	state["stowed_links"] = links
	state["is_stowed"] = len(links) > 0

	return state, nil
}

// stow: Stow a package (create symlinks)
func (s *StowModule) stow(L *lua.LState) int {
	config, err := s.parseConfig(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	results := L.NewTable()

	for _, pkg := range config.Packages {
		result := s.stowSingle(L, config, pkg)
		results.Append(result)
	}

	L.Push(results)
	L.Push(lua.LNil)
	return 2
}

func (s *StowModule) stowSingle(L *lua.LState, config *stowConfig, pkg string) *lua.LTable {
	result := L.NewTable()
	resourceID := fmt.Sprintf("stow:%s:%s", pkg, config.Target)

	// Check current state
	currentState, err := s.getCurrentState(config, pkg)
	if err != nil {
		L.SetField(result, "error", lua.LString(fmt.Sprintf("failed to get current state: %v", err)))
		return result
	}

	// Check if already stowed (idempotency)
	if isStowed, ok := currentState["is_stowed"].(bool); ok && isStowed {
		stowedLinks, _ := currentState["stowed_links"].([]string)
		files, _ := currentState["files"].([]string)
		
		// If all files are stowed, no action needed
		if len(stowedLinks) == len(files) {
			if s.sm != nil {
				stateData, _ := json.Marshal(currentState)
				s.sm.Set(resourceID, string(stateData))
			}
			
			L.SetField(result, "changed", lua.LBool(false))
			L.SetField(result, "status", lua.LString("already_stowed"))
			L.SetField(result, "package", lua.LString(pkg))
			return result
		}
	}

	// Execute stow command
	args := s.buildStowArgs(config, "-S", pkg)
	
	var cmd *exec.Cmd
	if config.DelegatedTo != "" {
		L.SetField(result, "error", lua.LString("remote execution not yet implemented"))
		return result
	} else {
		cmd = exec.Command("stow", args...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		L.SetField(result, "error", lua.LString(fmt.Sprintf("stow failed: %v - %s", err, string(output))))
		return result
	}

	// Get new state
	newState, err := s.getCurrentState(config, pkg)
	if err != nil {
		L.SetField(result, "error", lua.LString(fmt.Sprintf("failed to get new state: %v", err)))
		return result
	}

	// Record in state
	if s.sm != nil {
		stateData, _ := json.Marshal(newState)
		s.sm.Set(resourceID, string(stateData))
	}

	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "status", lua.LString("stowed"))
	L.SetField(result, "package", lua.LString(pkg))
	L.SetField(result, "target", lua.LString(config.Target))
	L.SetField(result, "output", lua.LString(string(output)))
	
	if links, ok := newState["stowed_links"].([]string); ok {
		linksTable := L.NewTable()
		for i, link := range links {
			linksTable.RawSetInt(i+1, lua.LString(link))
		}
		L.SetField(result, "links", linksTable)
	}

	return result
}

// unstow: Remove stowed symlinks
func (s *StowModule) unstow(L *lua.LState) int {
	config, err := s.parseConfig(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	results := L.NewTable()

	for _, pkg := range config.Packages {
		result := s.unstowSingle(L, config, pkg)
		results.Append(result)
	}

	L.Push(results)
	L.Push(lua.LNil)
	return 2
}

func (s *StowModule) unstowSingle(L *lua.LState, config *stowConfig, pkg string) *lua.LTable {
	result := L.NewTable()
	resourceID := fmt.Sprintf("stow:%s:%s", pkg, config.Target)

	// Check current state
	currentState, err := s.getCurrentState(config, pkg)
	if err != nil {
		L.SetField(result, "error", lua.LString(fmt.Sprintf("failed to get current state: %v", err)))
		return result
	}

	// Check if not stowed (idempotency)
	if isStowed, ok := currentState["is_stowed"].(bool); ok && !isStowed {
		L.SetField(result, "changed", lua.LBool(false))
		L.SetField(result, "status", lua.LString("not_stowed"))
		L.SetField(result, "package", lua.LString(pkg))
		return result
	}

	var cmd *exec.Cmd
	if config.DelegatedTo != "" {
		L.SetField(result, "error", lua.LString("remote execution not yet implemented"))
		return result
	} else {
		cmd = exec.Command("stow", args...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		L.SetField(result, "error", lua.LString(fmt.Sprintf("unstow failed: %v - %s", err, string(output))))
		return result
	}

	// Remove from state
	if s.sm != nil {
		s.sm.Delete(resourceID)
	}

	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "status", lua.LString("unstowed"))
	L.SetField(result, "package", lua.LString(pkg))
	L.SetField(result, "output", lua.LString(string(output)))

	return result
}

// restow: Restow a package (unstow then stow)
func (s *StowModule) restow(L *lua.LState) int {
	config, err := s.parseConfig(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	results := L.NewTable()

	for _, pkg := range config.Packages {
		result := s.restowSingle(L, config, pkg)
		results.Append(result)
	}

	L.Push(results)
	L.Push(lua.LNil)
	return 2
}

func (s *StowModule) restowSingle(L *lua.LState, config *stowConfig, pkg string) *lua.LTable {
	result := L.NewTable()

	var cmd *exec.Cmd
	if config.DelegatedTo != "" {
		L.SetField(result, "error", lua.LString("remote execution not yet implemented"))
		return result
	} else {
		cmd = exec.Command("stow", args...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		L.SetField(result, "error", lua.LString(fmt.Sprintf("restow failed: %v - %s", err, string(output))))
		return result
	}

	// Get new state
	newState, err := s.getCurrentState(config, pkg)
	if err != nil {
		L.SetField(result, "error", lua.LString(fmt.Sprintf("failed to get new state: %v", err)))
		return result
	}

	// Record in state
	resourceID := fmt.Sprintf("stow:%s:%s", pkg, config.Target)
	if s.sm != nil {
		stateData, _ := json.Marshal(newState)
		s.sm.Set(resourceID, string(stateData))
	}

	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "status", lua.LString("restowed"))
	L.SetField(result, "package", lua.LString(pkg))
	L.SetField(result, "output", lua.LString(string(output)))

	return result
}

// adopt: Adopt existing files into the stow package
func (s *StowModule) adopt(L *lua.LState) int {
	config, err := s.parseConfig(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	results := L.NewTable()

	for _, pkg := range config.Packages {
		result := s.adoptSingle(L, config, pkg)
		results.Append(result)
	}

	L.Push(results)
	L.Push(lua.LNil)
	return 2
}

func (s *StowModule) adoptSingle(L *lua.LState, config *stowConfig, pkg string) *lua.LTable {
	result := L.NewTable()

	args := s.buildStowArgs(config, "--adopt", pkg)
	
	var cmd *exec.Cmd
	if config.DelegatedTo != "" {
		L.SetField(result, "error", lua.LString("remote execution not yet implemented"))
		return result
	} else {
		cmd = exec.Command("stow", args...)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		L.SetField(result, "error", lua.LString(fmt.Sprintf("adopt failed: %v - %s", err, string(output))))
		return result
	}
}

// check: Check stow operations without executing
func (s *StowModule) check(L *lua.LState) int {
	config, err := s.parseConfig(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	results := L.NewTable()

	for _, pkg := range config.Packages {
		result := s.checkSingle(L, config, pkg)
		results.Append(result)
	}

	L.Push(results)
	L.Push(lua.LNil)
	return 2
}

func (s *StowModule) checkSingle(L *lua.LState, config *stowConfig, pkg string) *lua.LTable {
	result := L.NewTable()

	args := append([]string{"-n", "-v"}, s.buildStowArgs(config, "-S", pkg)...)
	
	var cmd *exec.Cmd
	if config.DelegatedTo != "" {
		L.SetField(result, "error", lua.LString("remote execution not yet implemented"))
		return result
	} else {
		cmd = exec.Command("stow", args...)
	}

	output, err := cmd.CombinedOutput()
	
	L.SetField(result, "package", lua.LString(pkg))
	L.SetField(result, "output", lua.LString(string(output)))
	L.SetField(result, "would_succeed", lua.LBool(err == nil))

	return result
	L.SetField(result, "package", lua.LString(pkg))
	L.SetField(result, "output", lua.LString(string(output)))
	L.SetField(result, "would_succeed", lua.LBool(err == nil))

	return result
}

// simulate: Simulate stow operations
func (s *StowModule) simulate(L *lua.LState) int {
	return s.check(L)
}

// getLinks: Get list of symlinks created by stow for a package
func (s *StowModule) getLinks(L *lua.LState) int {
	config, err := s.parseConfig(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	results := L.NewTable()

	for _, pkg := range config.Packages {
		state, err := s.getCurrentState(config, pkg)
		if err != nil {
			// handle error
			continue
		}

		linksTable := L.NewTable()
		if links, ok := state["stowed_links"].([]string); ok {
			for i, link := range links {
				linksTable.RawSetInt(i+1, lua.LString(link))
			}
		}
		results.Append(linksTable)
	}

	L.Push(results)
	L.Push(lua.LNil)
	return 2
}

// isStowed: Check if a package is currently stowed
func (s *StowModule) isStowed(L *lua.LState) int {
	config, err := s.parseConfig(L)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	results := L.NewTable()

	for _, pkg := range config.Packages {
		state, err := s.getCurrentState(config, pkg)
		if err != nil {
			// handle error
			continue
		}

		isStowed := false
		if val, ok := state["is_stowed"].(bool); ok {
			isStowed = val
		}
		results.Append(lua.LBool(isStowed))
	}

	L.Push(results)
	L.Push(lua.LNil)
	return 2
}

// listPackages: List available packages in stow directory
func (s *StowModule) listPackages(L *lua.LState) int {
	stowDir := filepath.Join(os.Getenv("HOME"), ".dotfiles")
	
	if L.GetTop() >= 1 {
		if tbl, ok := L.Get(1).(*lua.LTable); ok {
			if dir := tbl.RawGetString("dir"); dir != lua.LNil {
				stowDir = dir.String()
			}
		}
	}

	entries, err := os.ReadDir(stowDir)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to read stow directory: %v", err)))
		return 2
	}

	packages := L.NewTable()
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			packages.Append(lua.LString(entry.Name()))
		}
	}

	L.Push(packages)
	L.Push(lua.LNil)
	return 2
}

// verify: Verify integrity of stowed symlinks
func (s *StowModule) verify(L *lua.LState) int {
	config, err := s.parseConfig(L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	results := L.NewTable()

	for _, pkg := range config.Packages {
		result := s.verifySingle(L, config, pkg)
		results.Append(result)
	}

	L.Push(results)
	L.Push(lua.LNil)
	return 2
}

func (s *StowModule) verifySingle(L *lua.LState, config *stowConfig, pkg string) *lua.LTable {
	result := L.NewTable()

	state, err := s.getCurrentState(config, pkg)
	if err != nil {
		L.SetField(result, "error", lua.LString(err.Error()))
		return result
	}

	L.SetField(result, "package", lua.LString(pkg))
	
	stowedLinks, _ := state["stowed_links"].([]string)
	files, _ := state["files"].([]string)
	
	L.SetField(result, "total_files", lua.LNumber(len(files)))
	L.SetField(result, "stowed_links", lua.LNumber(len(stowedLinks)))
	L.SetField(result, "is_complete", lua.LBool(len(stowedLinks) == len(files)))
	
	// Find broken links
	brokenLinks := []string{}
	for _, link := range stowedLinks {
		targetPath := filepath.Join(config.Target, link)
		if _, err := os.Stat(targetPath); err != nil {
			brokenLinks = append(brokenLinks, link)
		}
	}
	
	brokenTable := L.NewTable()
	for i, link := range brokenLinks {
		brokenTable.RawSetInt(i+1, lua.LString(link))
	}
	L.SetField(result, "broken_links", brokenTable)
	L.SetField(result, "is_valid", lua.LBool(len(brokenLinks) == 0))

	return result
}