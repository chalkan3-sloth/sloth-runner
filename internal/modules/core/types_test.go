package core

import (
	"testing"
)

func TestCoreModuleInfo(t *testing.T) {
	info := CoreModuleInfo{
		Name:         "test-module",
		Version:      "1.0.0",
		Description:  "Test module",
		Author:       "Test Author",
		Category:     "core",
		Dependencies: []string{"dep1", "dep2"},
	}

	if info.Name != "test-module" {
		t.Errorf("Expected name 'test-module', got '%s'", info.Name)
	}

	if info.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", info.Version)
	}

	if info.Description != "Test module" {
		t.Errorf("Expected description 'Test module', got '%s'", info.Description)
	}

	if info.Author != "Test Author" {
		t.Errorf("Expected author 'Test Author', got '%s'", info.Author)
	}

	if info.Category != "core" {
		t.Errorf("Expected category 'core', got '%s'", info.Category)
	}

	if len(info.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(info.Dependencies))
	}
}

func TestCoreModuleInfo_EmptyDependencies(t *testing.T) {
	info := CoreModuleInfo{
		Name:         "simple-module",
		Version:      "1.0.0",
		Description:  "Simple module",
		Author:       "Author",
		Category:     "core",
		Dependencies: []string{},
	}

	if len(info.Dependencies) != 0 {
		t.Errorf("Expected 0 dependencies, got %d", len(info.Dependencies))
	}
}

func TestCoreModuleInfo_NilDependencies(t *testing.T) {
	info := CoreModuleInfo{
		Name:         "no-deps-module",
		Version:      "2.0.0",
		Description:  "Module without dependencies",
		Author:       "Author",
		Category:     "util",
		Dependencies: nil,
	}

	if info.Dependencies != nil && len(info.Dependencies) != 0 {
		t.Errorf("Expected nil or empty dependencies, got %v", info.Dependencies)
	}
}
