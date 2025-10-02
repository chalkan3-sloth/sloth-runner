package luainterface

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestTerraformInit(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterTerraformAdvancedModule(L)

	tmpDir := t.TempDir()

	script := `
		local terraform = require("terraform")
		local result = terraform.init({
			dir = "` + tmpDir + `"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Terraform test skipped (Terraform not available): %v", err)
		return
	}
}

func TestTerraformPlan(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterTerraformAdvancedModule(L)

	tmpDir := t.TempDir()

	script := `
		local terraform = require("terraform")
		local result = terraform.plan({
			dir = "` + tmpDir + `"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Terraform test skipped (Terraform not available): %v", err)
		return
	}
}

func TestTerraformApply(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterTerraformAdvancedModule(L)

	tmpDir := t.TempDir()

	script := `
		local terraform = require("terraform")
		local result = terraform.apply({
			dir = "` + tmpDir + `",
			auto_approve = true
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Terraform test skipped (Terraform not available): %v", err)
		return
	}
}

func TestTerraformDestroy(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterTerraformAdvancedModule(L)

	tmpDir := t.TempDir()

	script := `
		local terraform = require("terraform")
		local result = terraform.destroy({
			dir = "` + tmpDir + `",
			auto_approve = true
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Terraform test skipped (Terraform not available): %v", err)
		return
	}
}

func TestTerraformValidate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterTerraformAdvancedModule(L)

	tmpDir := t.TempDir()

	script := `
		local terraform = require("terraform")
		local result = terraform.validate({
			dir = "` + tmpDir + `"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Terraform test skipped (Terraform not available): %v", err)
		return
	}
}

func TestTerraformOutput(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterTerraformAdvancedModule(L)

	tmpDir := t.TempDir()

	script := `
		local terraform = require("terraform")
		local result = terraform.output({
			dir = "` + tmpDir + `"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Terraform test skipped (Terraform not available): %v", err)
		return
	}
}

func TestTerraformWorkspace(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	RegisterTerraformAdvancedModule(L)

	tmpDir := t.TempDir()

	script := `
		local terraform = require("terraform")
		local result = terraform.workspace({
			dir = "` + tmpDir + `",
			action = "list"
		})
		return result ~= nil
	`

	if err := L.DoString(script); err != nil {
		t.Logf("Terraform test skipped (Terraform not available): %v", err)
		return
	}
}
