package luainterface

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestFileOpsCopy(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() (string, string, error)
		options     map[string]string
		wantErr     bool
		checkFunc   func(dst string) error
	}{
		{
			name: "basic copy",
			setupFunc: func() (string, string, error) {
				tmpDir := t.TempDir()
				src := filepath.Join(tmpDir, "source.txt")
				dst := filepath.Join(tmpDir, "dest.txt")
				if err := os.WriteFile(src, []byte("test content"), 0644); err != nil {
					return "", "", err
				}
				return src, dst, nil
			},
			wantErr: false,
			checkFunc: func(dst string) error {
				content, err := os.ReadFile(dst)
				if err != nil {
					return err
				}
				if string(content) != "test content" {
					t.Errorf("content mismatch: got %s", string(content))
				}
				return nil
			},
		},
		{
			name: "copy with mode",
			setupFunc: func() (string, string, error) {
				tmpDir := t.TempDir()
				src := filepath.Join(tmpDir, "source.txt")
				dst := filepath.Join(tmpDir, "dest.txt")
				if err := os.WriteFile(src, []byte("test"), 0644); err != nil {
					return "", "", err
				}
				return src, dst, nil
			},
			options: map[string]string{"mode": "0600"},
			wantErr: false,
			checkFunc: func(dst string) error {
				info, err := os.Stat(dst)
				if err != nil {
					return err
				}
				if info.Mode().Perm() != 0600 {
					t.Errorf("mode mismatch: got %o", info.Mode().Perm())
				}
				return nil
			},
		},
		{
			name: "copy non-existent source",
			setupFunc: func() (string, string, error) {
				tmpDir := t.TempDir()
				src := filepath.Join(tmpDir, "nonexistent.txt")
				dst := filepath.Join(tmpDir, "dest.txt")
				return src, dst, nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.PreloadModule("file_ops", NewFileOpsModule().Loader)

			src, dst, err := tt.setupFunc()
			if err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			// Build Lua call
			var code string
			if tt.options != nil {
				code = `
					local file_ops = require('file_ops')
					return file_ops.copy({
						src = "` + src + `",
						dest = "` + dst + `",`
				for k, v := range tt.options {
					code += `
						` + k + ` = "` + v + `",`
				}
				code += `
					})
				`
			} else {
				code = `
					local file_ops = require('file_ops')
					return file_ops.copy({src = "` + src + `", dest = "` + dst + `"})
				`
			}

			err = L.DoString(code)

			if tt.wantErr {
				if err == nil && L.Get(-2) != lua.LNil {
					t.Error("expected error, got success")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.checkFunc != nil {
				if err := tt.checkFunc(dst); err != nil {
					t.Errorf("check failed: %v", err)
				}
			}
		})
	}
}

func TestFileOpsFetch(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("file_ops", NewFileOpsModule().Loader)

	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "source.txt")
	dst := filepath.Join(tmpDir, "subdir", "dest.txt")

	if err := os.WriteFile(src, []byte("fetch test"), 0644); err != nil {
		t.Fatalf("failed to create source: %v", err)
	}

	code := `
		local file_ops = require('file_ops')
		return file_ops.fetch({src = "` + src + `", dest = "` + dst + `"})
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("fetch failed: %v", err)
	}

	// Verify file was fetched
	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read destination: %v", err)
	}

	if string(content) != "fetch test" {
		t.Errorf("content mismatch: got %s", string(content))
	}
}

func TestFileOpsTemplate(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("file_ops", NewFileOpsModule().Loader)

	tmpDir := t.TempDir()
	tmpl := filepath.Join(tmpDir, "template.txt")
	dst := filepath.Join(tmpDir, "rendered.txt")

	// Create template
	tmplContent := "Hello {{.Name}}, you are {{.Age}} years old"
	if err := os.WriteFile(tmpl, []byte(tmplContent), 0644); err != nil {
		t.Fatalf("failed to create template: %v", err)
	}

	code := `
		local file_ops = require('file_ops')
		local vars = {Name = "Alice", Age = 30}
		return file_ops.template({
			src = "` + tmpl + `",
			dest = "` + dst + `",
			vars = vars
		})
	`

	if err := L.DoString(code); err != nil {
		t.Fatalf("template failed: %v", err)
	}

	// Verify rendered content
	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("failed to read destination: %v", err)
	}

	expected := "Hello Alice, you are 30 years old"
	if string(content) != expected {
		t.Errorf("content mismatch: got %s, want %s", string(content), expected)
	}
}

func TestFileOpsLineinfile(t *testing.T) {
	tests := []struct {
		name        string
		initial     string
		line        string
		options     string
		expected    string
		wantChanged bool
	}{
		{
			name:        "add new line",
			initial:     "line1\nline2",
			line:        "line3",
			expected:    "line1\nline2\nline3",
			wantChanged: true,
		},
		{
			name:        "line already exists",
			initial:     "line1\nline2\nline3",
			line:        "line2",
			expected:    "line1\nline2\nline3",
			wantChanged: false,
		},
		{
			name:        "replace with regexp",
			initial:     "foo=bar\ntest=value",
			line:        "foo=baz",
			options:     `{regexp = "^foo="}`,
			expected:    "foo=baz\ntest=value",
			wantChanged: true,
		},
		{
			name:        "remove line",
			initial:     "line1\nline2\nline3",
			line:        "line2",
			options:     `{state = "absent"}`,
			expected:    "line1\nline3",
			wantChanged: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.PreloadModule("file_ops", NewFileOpsModule().Loader)

			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "test.txt")

			if tt.initial != "" {
				if err := os.WriteFile(path, []byte(tt.initial), 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
			}

			// Build Lua code
			var code string
			if tt.options != "" {
				code = `
					local file_ops = require('file_ops')
					local opts = ` + tt.options + `
					opts.path = "` + path + `"
					opts.line = "` + tt.line + `"
					return file_ops.lineinfile(opts)
				`
			} else {
				code = `
					local file_ops = require('file_ops')
					return file_ops.lineinfile({path = "` + path + `", line = "` + tt.line + `"})
				`
			}

			if err := L.DoString(code); err != nil {
				t.Fatalf("lineinfile failed: %v", err)
			}

			// Verify content
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			if string(content) != tt.expected {
				t.Errorf("content mismatch:\ngot:  %q\nwant: %q", string(content), tt.expected)
			}
		})
	}
}

func TestFileOpsBlockinfile(t *testing.T) {
	tests := []struct {
		name        string
		initial     string
		block       string
		options     string
		wantChanged bool
		checkFunc   func(content string) error
	}{
		{
			name:        "add new block",
			initial:     "header\nfooter",
			block:       "line1\\nline2",
			wantChanged: true,
			checkFunc: func(content string) error {
				if !strings.Contains(content, "BEGIN MANAGED BLOCK") {
					t.Error("missing begin marker")
				}
				if !strings.Contains(content, "line1") {
					t.Error("missing block content")
				}
				return nil
			},
		},
		{
			name:        "update existing block",
			initial:     "header\n# BEGIN MANAGED BLOCK\nold content\n# END MANAGED BLOCK\nfooter",
			block:       "new content",
			wantChanged: true,
			checkFunc: func(content string) error {
				if strings.Contains(content, "old content") {
					t.Error("old content still present")
				}
				if !strings.Contains(content, "new content") {
					t.Error("new content missing")
				}
				return nil
			},
		},
		{
			name:        "remove block",
			initial:     "header\n# BEGIN MANAGED BLOCK\ncontent\n# END MANAGED BLOCK\nfooter",
			block:       "",
			options:     `{state = "absent"}`,
			wantChanged: true,
			checkFunc: func(content string) error {
				if strings.Contains(content, "BEGIN MANAGED BLOCK") {
					t.Error("markers still present")
				}
				if strings.Contains(content, "content") {
					t.Error("block content still present")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.PreloadModule("file_ops", NewFileOpsModule().Loader)

			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "test.txt")

			if tt.initial != "" {
				if err := os.WriteFile(path, []byte(tt.initial), 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
			}

			var code string
			if tt.options != "" {
				code = `
					local file_ops = require('file_ops')
					local opts = ` + tt.options + `
					opts.path = "` + path + `"
					opts.block = "` + tt.block + `"
					return file_ops.blockinfile(opts)
				`
			} else {
				code = `
					local file_ops = require('file_ops')
					return file_ops.blockinfile({path = "` + path + `", block = "` + tt.block + `"})
				`
			}

			if err := L.DoString(code); err != nil {
				t.Fatalf("blockinfile failed: %v", err)
			}

			if tt.checkFunc != nil {
				content, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("failed to read file: %v", err)
				}
				if err := tt.checkFunc(string(content)); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestFileOpsReplace(t *testing.T) {
	tests := []struct {
		name        string
		initial     string
		pattern     string
		replacement string
		expected    string
		wantChanged bool
	}{
		{
			name:        "simple replace",
			initial:     "hello world",
			pattern:     "world",
			replacement: "universe",
			expected:    "hello universe",
			wantChanged: true,
		},
		{
			name:        "regexp replace",
			initial:     "version=1.2.3\\nother=value",
			pattern:     `version=\\d+\\.\\d+\\.\\d+`,
			replacement: "version=2.0.0",
			expected:    "version=2.0.0\\nother=value",
			wantChanged: true,
		},
		{
			name:        "no match",
			initial:     "hello world",
			pattern:     "foo",
			replacement: "bar",
			expected:    "hello world",
			wantChanged: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.PreloadModule("file_ops", NewFileOpsModule().Loader)

			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "test.txt")

			if err := os.WriteFile(path, []byte(tt.initial), 0644); err != nil {
				t.Fatalf("failed to create file: %v", err)
			}

			code := `
				local file_ops = require('file_ops')
				return file_ops.replace({path = "` + path + `", pattern = "` + tt.pattern + `", replacement = "` + tt.replacement + `"})
			`

			if err := L.DoString(code); err != nil {
				t.Fatalf("replace failed: %v", err)
			}

			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			if string(content) != tt.expected {
				t.Errorf("content mismatch:\ngot:  %q\nwant: %q", string(content), tt.expected)
			}
		})
	}
}

func TestFileOpsUnarchive(t *testing.T) {
	tests := []struct {
		name      string
		archType  string
		setupFunc func(tmpDir string) (string, error)
	}{
		{
			name:     "extract zip",
			archType: ".zip",
			setupFunc: func(tmpDir string) (string, error) {
				zipPath := filepath.Join(tmpDir, "test.zip")
				zipFile, err := os.Create(zipPath)
				if err != nil {
					return "", err
				}
				defer zipFile.Close()

				zw := zip.NewWriter(zipFile)
				defer zw.Close()

				// Add test file
				fw, err := zw.Create("test.txt")
				if err != nil {
					return "", err
				}
				fw.Write([]byte("zip content"))

				return zipPath, nil
			},
		},
		{
			name:     "extract tar.gz",
			archType: ".tar.gz",
			setupFunc: func(tmpDir string) (string, error) {
				tarPath := filepath.Join(tmpDir, "test.tar.gz")
				tarFile, err := os.Create(tarPath)
				if err != nil {
					return "", err
				}
				defer tarFile.Close()

				gzw := gzip.NewWriter(tarFile)
				defer gzw.Close()

				tw := tar.NewWriter(gzw)
				defer tw.Close()

				// Add test file
				hdr := &tar.Header{
					Name: "test.txt",
					Mode: 0600,
					Size: int64(len("tar content")),
				}
				if err := tw.WriteHeader(hdr); err != nil {
					return "", err
				}
				if _, err := tw.Write([]byte("tar content")); err != nil {
					return "", err
				}

				return tarPath, nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.PreloadModule("file_ops", NewFileOpsModule().Loader)

			tmpDir := t.TempDir()
			archPath, err := tt.setupFunc(tmpDir)
			if err != nil {
				t.Fatalf("failed to setup archive: %v", err)
			}

			dstDir := filepath.Join(tmpDir, "extracted")

			code := `
				local file_ops = require('file_ops')
				return file_ops.unarchive({src = "` + archPath + `", dest = "` + dstDir + `"})
			`

			if err := L.DoString(code); err != nil {
				t.Fatalf("unarchive failed: %v", err)
			}

			// Verify extracted file exists
			extractedFile := filepath.Join(dstDir, "test.txt")
			if _, err := os.Stat(extractedFile); err != nil {
				t.Errorf("extracted file not found: %v", err)
			}
		})
	}
}

func TestFileOpsStat(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(tmpDir string) string
		checkFunc func(L *lua.LState)
	}{
		{
			name: "stat regular file",
			setupFunc: func(tmpDir string) string {
				path := filepath.Join(tmpDir, "test.txt")
				os.WriteFile(path, []byte("test content"), 0644)
				return path
			},
			checkFunc: func(L *lua.LState) {
				code := `
					exists = result.exists
					is_file = result.is_file
					size = result.size
				`
				if err := L.DoString(code); err != nil {
					t.Fatalf("failed to extract result: %v", err)
				}
			},
		},
		{
			name: "stat directory",
			setupFunc: func(tmpDir string) string {
				path := filepath.Join(tmpDir, "testdir")
				os.Mkdir(path, 0755)
				return path
			},
			checkFunc: func(L *lua.LState) {
				code := `
					exists = result.exists
					is_dir = result.is_dir
				`
				if err := L.DoString(code); err != nil {
					t.Fatalf("failed to extract result: %v", err)
				}
			},
		},
		{
			name: "stat non-existent",
			setupFunc: func(tmpDir string) string {
				return filepath.Join(tmpDir, "nonexistent")
			},
			checkFunc: func(L *lua.LState) {
				code := `
					exists = result.exists
				`
				if err := L.DoString(code); err != nil {
					t.Fatalf("failed to extract result: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			L := lua.NewState()
			defer L.Close()

			L.PreloadModule("file_ops", NewFileOpsModule().Loader)

			tmpDir := t.TempDir()
			path := tt.setupFunc(tmpDir)

			code := `
				local file_ops = require('file_ops')
				result = file_ops.stat({path = "` + path + `"})
			`

			if err := L.DoString(code); err != nil {
				t.Fatalf("stat failed: %v", err)
			}

			if tt.checkFunc != nil {
				tt.checkFunc(L)
			}
		})
	}
}

func TestFileOpsEdgeCases(t *testing.T) {
	t.Run("copy creates parent directories", func(t *testing.T) {
		L := lua.NewState()
		defer L.Close()

		L.PreloadModule("file_ops", NewFileOpsModule().Loader)

		tmpDir := t.TempDir()
		src := filepath.Join(tmpDir, "source.txt")
		dst := filepath.Join(tmpDir, "deep", "nested", "dir", "dest.txt")

		os.WriteFile(src, []byte("test"), 0644)

		code := `
			local file_ops = require('file_ops')
			return file_ops.copy({src = "` + src + `", dest = "` + dst + `"})
		`

		if err := L.DoString(code); err != nil {
			t.Fatalf("copy failed: %v", err)
		}

		if _, err := os.Stat(dst); err != nil {
			t.Error("destination file should exist")
		}
	})

	t.Run("lineinfile creates new file", func(t *testing.T) {
		L := lua.NewState()
		defer L.Close()

		L.PreloadModule("file_ops", NewFileOpsModule().Loader)

		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "newfile.txt")

		code := `
			local file_ops = require('file_ops')
			return file_ops.lineinfile({path = "` + path + `", line = "first line"})
		`

		if err := L.DoString(code); err != nil {
			t.Fatalf("lineinfile failed: %v", err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal("file should be created")
		}

		if string(content) != "first line" {
			t.Errorf("content mismatch: got %q", string(content))
		}
	})

	t.Run("replace with invalid regexp", func(t *testing.T) {
		L := lua.NewState()
		defer L.Close()

		L.PreloadModule("file_ops", NewFileOpsModule().Loader)

		tmpDir := t.TempDir()
		path := filepath.Join(tmpDir, "test.txt")
		os.WriteFile(path, []byte("test"), 0644)

		code := `
			local file_ops = require('file_ops')
			local result, err = file_ops.replace({path = "` + path + `", pattern = "[invalid", replace = "replacement"})
			if result == nil then
				return true  -- Error expected
			end
			return false
		`

		err := L.DoString(code)
		if err != nil {
			// Error during execution is fine for invalid regexp
			return
		}
		
		result := L.Get(-1)
		if result == lua.LFalse {
			t.Error("should return error for invalid regexp")
		}
	})
}
