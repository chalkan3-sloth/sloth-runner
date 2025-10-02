package modules

import (
	"testing"
)

func TestGetAllModuleDocs(t *testing.T) {
	docs := GetAllModuleDocs()

	if len(docs) == 0 {
		t.Fatal("Expected at least one module documentation, got none")
	}

	// Check that we have some essential modules
	essentialModules := []string{"pkg", "systemd", "user", "file", "http"}
	
	moduleMap := make(map[string]bool)
	for _, doc := range docs {
		moduleMap[doc.Name] = true
		
		// Validate each module has necessary fields
		if doc.Name == "" {
			t.Error("Module has empty name")
		}
		if doc.Description == "" {
			t.Errorf("Module %s has empty description", doc.Name)
		}
		if len(doc.Functions) == 0 {
			t.Errorf("Module %s has no functions", doc.Name)
		}
		
		// Validate each function
		for _, fn := range doc.Functions {
			if fn.Name == "" {
				t.Errorf("Module %s has a function with empty name", doc.Name)
			}
			if fn.Description == "" {
				t.Errorf("Function %s in module %s has empty description", fn.Name, doc.Name)
			}
			if fn.Example == "" {
				t.Errorf("Function %s in module %s has no example", fn.Name, doc.Name)
			}
		}
	}
	
	// Check for essential modules
	for _, modName := range essentialModules {
		if !moduleMap[modName] {
			t.Errorf("Expected essential module '%s' not found", modName)
		}
	}
}

func TestModuleDocStructure(t *testing.T) {
	docs := GetAllModuleDocs()
	
	for _, doc := range docs {
		t.Run(doc.Name, func(t *testing.T) {
			if doc.Name == "" {
				t.Error("Module name is empty")
			}
			if doc.Description == "" {
				t.Error("Module description is empty")
			}
			if len(doc.Functions) == 0 {
				t.Error("Module has no functions")
			}
			
			for i, fn := range doc.Functions {
				if fn.Name == "" {
					t.Errorf("Function %d has empty name", i)
				}
				if fn.Description == "" {
					t.Errorf("Function %s has empty description", fn.Name)
				}
				if fn.Example == "" {
					t.Errorf("Function %s has empty example", fn.Name)
				}
			}
		})
	}
}

func TestSpecificModules(t *testing.T) {
	docs := GetAllModuleDocs()
	
	tests := []struct {
		moduleName    string
		minFunctions  int
	}{
		{"pkg", 4},
		{"systemd", 6},
		{"user", 5},
		{"file", 5},
		{"http", 2},
		{"json", 2},
		{"yaml", 2},
		{"log", 3},
	}
	
	moduleMap := make(map[string]ModuleDoc)
	for _, doc := range docs {
		moduleMap[doc.Name] = doc
	}
	
	for _, tt := range tests {
		t.Run(tt.moduleName, func(t *testing.T) {
			doc, exists := moduleMap[tt.moduleName]
			if !exists {
				t.Fatalf("Module %s not found", tt.moduleName)
			}
			
			if len(doc.Functions) < tt.minFunctions {
				t.Errorf("Module %s has %d functions, expected at least %d",
					tt.moduleName, len(doc.Functions), tt.minFunctions)
			}
		})
	}
}
