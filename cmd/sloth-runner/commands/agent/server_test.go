package agent

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateTarData(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	// Create test files
	testFiles := map[string]string{
		"file1.txt":       "content1",
		"file2.txt":       "content2",
		"subdir/file3.txt": "content3",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", fullPath, err)
		}
	}

	// Create tar
	var buf bytes.Buffer
	err := createTarData(tmpDir, &buf)
	if err != nil {
		t.Fatalf("createTarData failed: %v", err)
	}

	// Verify tar contains all files
	tr := tar.NewReader(&buf)
	foundFiles := make(map[string]bool)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Error reading tar: %v", err)
		}

		if header.Typeflag == tar.TypeReg {
			foundFiles[header.Name] = true

			// Read content and verify
			var content bytes.Buffer
			if _, err := io.Copy(&content, tr); err != nil {
				t.Fatalf("Error reading file content: %v", err)
			}

			// Normalize path for comparison
			normalizedName := filepath.ToSlash(header.Name)
			if normalizedName != "" && normalizedName[0] == '/' {
				normalizedName = normalizedName[1:]
			}

			if expectedContent, ok := testFiles[normalizedName]; ok {
				if content.String() != expectedContent {
					t.Errorf("File %s content mismatch: got %q, want %q",
						normalizedName, content.String(), expectedContent)
				}
			}
		}
	}

	// Check all files were found
	if len(foundFiles) != len(testFiles) {
		t.Errorf("Expected %d files in tar, found %d", len(testFiles), len(foundFiles))
	}
}

func TestExtractTarData(t *testing.T) {
	// Create a tar archive in memory
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	testFiles := map[string]string{
		"file1.txt":       "content1",
		"file2.txt":       "content2",
		"subdir/file3.txt": "content3",
	}

	for path, content := range testFiles {
		header := &tar.Header{
			Name: path,
			Mode: 0644,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(header); err != nil {
			t.Fatalf("Failed to write header: %v", err)
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			t.Fatalf("Failed to write content: %v", err)
		}
	}
	tw.Close()

	// Extract to temporary directory
	tmpDir := t.TempDir()
	reader := bytes.NewReader(buf.Bytes())

	err := extractTarData(reader, tmpDir)
	if err != nil {
		t.Fatalf("extractTarData failed: %v", err)
	}

	// Verify all files were extracted with correct content
	for path, expectedContent := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Errorf("Failed to read extracted file %s: %v", path, err)
			continue
		}
		if string(content) != expectedContent {
			t.Errorf("File %s content mismatch: got %q, want %q", path, string(content), expectedContent)
		}
	}
}

func TestCreateAndExtractTarData_RoundTrip(t *testing.T) {
	// Create source directory
	srcDir := t.TempDir()

	testFiles := map[string]string{
		"file1.txt":            "Hello World",
		"dir1/file2.txt":       "Test Content",
		"dir1/dir2/file3.txt":  "Nested Content",
		"empty.txt":            "",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(srcDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}

	// Create tar
	var buf bytes.Buffer
	if err := createTarData(srcDir, &buf); err != nil {
		t.Fatalf("createTarData failed: %v", err)
	}

	// Extract tar
	destDir := t.TempDir()
	reader := bytes.NewReader(buf.Bytes())
	if err := extractTarData(reader, destDir); err != nil {
		t.Fatalf("extractTarData failed: %v", err)
	}

	// Verify all files match
	for path, expectedContent := range testFiles {
		srcPath := filepath.Join(srcDir, path)
		destPath := filepath.Join(destDir, path)

		srcContent, err := os.ReadFile(srcPath)
		if err != nil {
			t.Fatalf("Failed to read source file: %v", err)
		}

		destContent, err := os.ReadFile(destPath)
		if err != nil {
			t.Fatalf("Failed to read destination file: %v", err)
		}

		if !bytes.Equal(srcContent, destContent) {
			t.Errorf("File %s content mismatch after round trip", path)
		}

		if string(destContent) != expectedContent {
			t.Errorf("File %s content: got %q, want %q", path, string(destContent), expectedContent)
		}
	}
}

func TestExtractTarData_EmptyTar(t *testing.T) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	tw.Close()

	tmpDir := t.TempDir()
	reader := bytes.NewReader(buf.Bytes())

	err := extractTarData(reader, tmpDir)
	if err != nil {
		t.Errorf("extractTarData with empty tar should not fail: %v", err)
	}
}

func TestExtractTarData_DirectoriesOnly(t *testing.T) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	// Add only directories
	dirs := []string{"dir1/", "dir1/dir2/", "dir3/"}
	for _, dir := range dirs {
		header := &tar.Header{
			Name:     dir,
			Mode:     0755,
			Typeflag: tar.TypeDir,
		}
		if err := tw.WriteHeader(header); err != nil {
			t.Fatalf("Failed to write directory header: %v", err)
		}
	}
	tw.Close()

	tmpDir := t.TempDir()
	reader := bytes.NewReader(buf.Bytes())

	err := extractTarData(reader, tmpDir)
	if err != nil {
		t.Fatalf("extractTarData failed: %v", err)
	}

	// Verify directories were created
	for _, dir := range dirs {
		dirPath := filepath.Join(tmpDir, dir)
		if info, err := os.Stat(dirPath); err != nil {
			t.Errorf("Directory %s was not created: %v", dir, err)
		} else if !info.IsDir() {
			t.Errorf("Path %s is not a directory", dir)
		}
	}
}

func BenchmarkCreateTarData(b *testing.B) {
	tmpDir := b.TempDir()

	// Create test files
	for i := 0; i < 10; i++ {
		path := filepath.Join(tmpDir, "file"+string(rune('0'+i))+".txt")
		os.WriteFile(path, []byte("test content for benchmarking"), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		createTarData(tmpDir, &buf)
	}
}

func BenchmarkExtractTarData(b *testing.B) {
	// Create tar once
	srcDir := b.TempDir()
	for i := 0; i < 10; i++ {
		path := filepath.Join(srcDir, "file"+string(rune('0'+i))+".txt")
		os.WriteFile(path, []byte("test content for benchmarking"), 0644)
	}

	var buf bytes.Buffer
	createTarData(srcDir, &buf)
	tarData := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		destDir := b.TempDir()
		reader := bytes.NewReader(tarData)
		extractTarData(reader, destDir)
	}
}
