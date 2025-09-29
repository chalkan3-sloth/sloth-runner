package core

// CoreModuleInfo represents module metadata for core modules
type CoreModuleInfo struct {
	Name         string
	Version      string
	Description  string
	Author       string
	Category     string
	Dependencies []string
}