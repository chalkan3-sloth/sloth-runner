package luainterface

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"text/template"

	lua "github.com/yuin/gopher-lua"
)

// FileOpsModule provides file operations similar to Ansible file modules
type FileOpsModule struct{}

// NewFileOpsModule creates a new file operations module
func NewFileOpsModule() *FileOpsModule {
	return &FileOpsModule{}
}

// Loader is the module loader function
func (f *FileOpsModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), f.exports())
	L.Push(mod)
	return 1
}

func (f *FileOpsModule) exports() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"copy":        f.copy,
		"fetch":       f.fetch,
		"template":    f.templateRender,
		"lineinfile":  f.lineinfile,
		"blockinfile": f.blockinfile,
		"replace":     f.replace,
		"unarchive":   f.unarchive,
		"stat":        f.stat,
	}
}

// copy copies a file from source to destination
// Usage: file_ops.copy({src="/path/to/source", dest="/path/to/dest", mode="0644"})
func (f *FileOpsModule) copy(L *lua.LState) int {
	opts := L.CheckTable(1)
	
	src := opts.RawGetString("src").String()
	dst := opts.RawGetString("dest").String()
	
	if src == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("src parameter is required"))
		return 2
	}
	if dst == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("dest parameter is required"))
		return 2
	}

	// Check if source exists
	srcInfo, err := os.Stat(src)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("source file not found: %v", err)))
		return 2
	}

	// Create destination directory if needed
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create directory: %v", err)))
		return 2
	}

	// Copy file content
	srcFile, err := os.Open(src)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to open source: %v", err)))
		return 2
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create destination: %v", err)))
		return 2
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to copy: %v", err)))
		return 2
	}

	// Set permissions if specified
	if mode := opts.RawGetString("mode"); mode != lua.LNil {
		if modeStr := mode.String(); modeStr != "" {
			var perm os.FileMode
			fmt.Sscanf(modeStr, "%o", &perm)
			os.Chmod(dst, perm)
		}
	} else {
		// Copy original permissions
		os.Chmod(dst, srcInfo.Mode())
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "src", lua.LString(src))
	L.SetField(result, "dest", lua.LString(dst))
	L.SetField(result, "size", lua.LNumber(srcInfo.Size()))

	L.Push(result)
	return 1
}

// fetch downloads a file from remote to local
// Usage: file_ops.fetch({src="/path/to/source", dest="/path/to/dest"})
func (f *FileOpsModule) fetch(L *lua.LState) int {
	opts := L.CheckTable(1)
	
	src := opts.RawGetString("src").String()
	dst := opts.RawGetString("dest").String()
	
	if src == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("src parameter is required"))
		return 2
	}
	if dst == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("dest parameter is required"))
		return 2
	}

	// For now, treat as local copy (would need agent integration for remote)
	srcFile, err := os.Open(src)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to open source: %v", err)))
		return 2
	}
	defer srcFile.Close()

	// Create destination directory
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create directory: %v", err)))
		return 2
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create destination: %v", err)))
		return 2
	}
	defer dstFile.Close()

	size, err := io.Copy(dstFile, srcFile)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to fetch: %v", err)))
		return 2
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "src", lua.LString(src))
	L.SetField(result, "dest", lua.LString(dst))
	L.SetField(result, "size", lua.LNumber(size))

	L.Push(result)
	return 1
}

// templateRender renders a template file with variables
// Usage: file_ops.template({src="/path/template.tpl", dest="/path/output", vars={key="value"}})
func (f *FileOpsModule) templateRender(L *lua.LState) int {
	opts := L.CheckTable(1)
	
	src := opts.RawGetString("src").String()
	dst := opts.RawGetString("dest").String()
	vars := opts.RawGetString("vars")
	
	if src == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("src parameter is required"))
		return 2
	}
	if dst == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("dest parameter is required"))
		return 2
	}

	// Read template file
	tmplContent, err := os.ReadFile(src)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to read template: %v", err)))
		return 2
	}

	// Convert Lua table to map
	data := make(map[string]interface{})
	if vars.Type() == lua.LTTable {
		varsTable := vars.(*lua.LTable)
		varsTable.ForEach(func(k, v lua.LValue) {
			if key, ok := k.(lua.LString); ok {
				data[string(key)] = luaValueToGoInterface(v)
			}
		})
	}

	// Parse and execute template
	tmpl, err := template.New("template").Parse(string(tmplContent))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to parse template: %v", err)))
		return 2
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to execute template: %v", err)))
		return 2
	}

	// Create destination directory
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create directory: %v", err)))
		return 2
	}

	// Write rendered content
	if err := os.WriteFile(dst, buf.Bytes(), 0644); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to write file: %v", err)))
		return 2
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "src", lua.LString(src))
	L.SetField(result, "dest", lua.LString(dst))

	L.Push(result)
	return 1
}

// lineinfile ensures a line exists in a file
// Usage: file_ops.lineinfile({path="/path/file", line="content", state="present", regexp="pattern"})
func (f *FileOpsModule) lineinfile(L *lua.LState) int {
	opts := L.CheckTable(1)
	
	path := opts.RawGetString("path").String()
	line := opts.RawGetString("line").String()
	
	if path == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("path parameter is required"))
		return 2
	}
	if line == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("line parameter is required"))
		return 2
	}

	// Read file or create if not exists
	var lines []string
	content, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("failed to read file: %v", err)))
			return 2
		}
	} else {
		lines = strings.Split(string(content), "\n")
	}

	// Check state (present/absent)
	state := "present"
	if stateVal := opts.RawGetString("state"); stateVal != lua.LNil && stateVal.String() != "" {
		state = stateVal.String()
	}

	// Check if line exists with regexp
	var re *regexp.Regexp
	if regexpVal := opts.RawGetString("regexp"); regexpVal != lua.LNil && regexpVal.String() != "" {
		re, err = regexp.Compile(regexpVal.String())
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("invalid regexp: %v", err)))
			return 2
		}
	}

	changed := false
	newLines := []string{}

	if state == "present" {
		found := false
		for _, l := range lines {
			if re != nil && re.MatchString(l) {
				if l != line {
					newLines = append(newLines, line)
					changed = true
				} else {
					newLines = append(newLines, l)
				}
				found = true
			} else if l == line {
				newLines = append(newLines, l)
				found = true
			} else {
				newLines = append(newLines, l)
			}
		}
		if !found {
			newLines = append(newLines, line)
			changed = true
		}
	} else { // absent
		for _, l := range lines {
			if re != nil && re.MatchString(l) {
				changed = true
				continue
			} else if l == line {
				changed = true
				continue
			}
			newLines = append(newLines, l)
		}
	}

	if changed {
		// Create directory if needed
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("failed to create directory: %v", err)))
			return 2
		}

		// Write file
		output := strings.Join(newLines, "\n")
		if err := os.WriteFile(path, []byte(output), 0644); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("failed to write file: %v", err)))
			return 2
		}
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(changed))
	L.SetField(result, "path", lua.LString(path))

	L.Push(result)
	return 1
}

// blockinfile inserts/updates/removes a block of lines in a file
// Usage: file_ops.blockinfile({path="/path/file", block="content", state="present", marker_begin="# BEGIN", marker_end="# END"})
func (f *FileOpsModule) blockinfile(L *lua.LState) int {
	opts := L.CheckTable(1)
	
	path := opts.RawGetString("path").String()
	block := opts.RawGetString("block").String()
	
	if path == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("path parameter is required"))
		return 2
	}
	if block == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("block parameter is required"))
		return 2
	}

	// Get markers
	markerBegin := "# BEGIN MANAGED BLOCK"
	markerEnd := "# END MANAGED BLOCK"

	if mb := opts.RawGetString("marker_begin"); mb != lua.LNil && mb.String() != "" {
		markerBegin = mb.String()
	}
	if me := opts.RawGetString("marker_end"); me != lua.LNil && me.String() != "" {
		markerEnd = me.String()
	}

	// Check state
	state := "present"
	if stateVal := opts.RawGetString("state"); stateVal != lua.LNil && stateVal.String() != "" {
		state = stateVal.String()
	}

	// Read file
	var lines []string
	content, err := os.ReadFile(path)
	if err != nil {
		if !os.IsNotExist(err) {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("failed to read file: %v", err)))
			return 2
		}
	} else {
		lines = strings.Split(string(content), "\n")
	}

	// Find existing block
	beginIdx := -1
	endIdx := -1
	for i, line := range lines {
		if strings.Contains(line, markerBegin) {
			beginIdx = i
		}
		if strings.Contains(line, markerEnd) {
			endIdx = i
		}
	}

	newLines := []string{}
	changed := false

	if state == "present" {
		if beginIdx >= 0 && endIdx >= 0 {
			// Replace existing block
			newLines = append(newLines, lines[:beginIdx]...)
			newLines = append(newLines, markerBegin)
			newLines = append(newLines, strings.Split(block, "\n")...)
			newLines = append(newLines, markerEnd)
			if endIdx+1 < len(lines) {
				newLines = append(newLines, lines[endIdx+1:]...)
			}
			changed = true
		} else {
			// Append new block
			newLines = append(newLines, lines...)
			newLines = append(newLines, markerBegin)
			newLines = append(newLines, strings.Split(block, "\n")...)
			newLines = append(newLines, markerEnd)
			changed = true
		}
	} else { // absent
		if beginIdx >= 0 && endIdx >= 0 {
			// Remove block
			newLines = append(newLines, lines[:beginIdx]...)
			if endIdx+1 < len(lines) {
				newLines = append(newLines, lines[endIdx+1:]...)
			}
			changed = true
		} else {
			newLines = lines
		}
	}

	if changed {
		// Create directory if needed
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("failed to create directory: %v", err)))
			return 2
		}

		// Write file
		output := strings.Join(newLines, "\n")
		if err := os.WriteFile(path, []byte(output), 0644); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("failed to write file: %v", err)))
			return 2
		}
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(changed))
	L.SetField(result, "path", lua.LString(path))

	L.Push(result)
	return 1
}

// replace replaces all occurrences of a pattern in a file
// Usage: file_ops.replace({path="/path/file", pattern="regex", replacement="text"})
func (f *FileOpsModule) replace(L *lua.LState) int {
	opts := L.CheckTable(1)
	
	path := opts.RawGetString("path").String()
	pattern := opts.RawGetString("pattern").String()
	replacement := opts.RawGetString("replacement").String()
	
	if path == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("path parameter is required"))
		return 2
	}
	if pattern == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("pattern parameter is required"))
		return 2
	}
	if replacement == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("replacement parameter is required"))
		return 2
	}

	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to read file: %v", err)))
		return 2
	}

	// Compile regexp
	re, err := regexp.Compile(pattern)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("invalid pattern: %v", err)))
		return 2
	}

	// Replace
	newContent := re.ReplaceAllString(string(content), replacement)
	changed := string(content) != newContent

	if changed {
		if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("failed to write file: %v", err)))
			return 2
		}
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(changed))
	L.SetField(result, "path", lua.LString(path))

	L.Push(result)
	return 1
}

// unarchive extracts an archive file
// Usage: file_ops.unarchive({src="/path/file.tar.gz", dest="/path/dest"})
func (f *FileOpsModule) unarchive(L *lua.LState) int {
	opts := L.CheckTable(1)
	
	src := opts.RawGetString("src").String()
	dst := opts.RawGetString("dest").String()
	
	if src == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("src parameter is required"))
		return 2
	}
	if dst == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("dest parameter is required"))
		return 2
	}

	// Create destination directory
	if err := os.MkdirAll(dst, 0755); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to create directory: %v", err)))
		return 2
	}

	// Detect archive type
	ext := filepath.Ext(src)
	var err error

	switch ext {
	case ".zip":
		err = extractZip(src, dst)
	case ".gz", ".tgz":
		err = extractTarGz(src, dst)
	case ".tar":
		err = extractTar(src, dst)
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("unsupported archive format: %s", ext)))
		return 2
	}

	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to extract: %v", err)))
		return 2
	}

	result := L.NewTable()
	L.SetField(result, "changed", lua.LBool(true))
	L.SetField(result, "src", lua.LString(src))
	L.SetField(result, "dest", lua.LString(dst))

	L.Push(result)
	return 1
}

// stat gets file information
// Usage: file_ops.stat({path="/path/file"})
func (f *FileOpsModule) stat(L *lua.LState) int {
	opts := L.CheckTable(1)
	
	path := opts.RawGetString("path").String()
	if path == "" {
		L.Push(lua.LNil)
		L.Push(lua.LString("path parameter is required"))
		return 2
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			result := L.NewTable()
			L.SetField(result, "exists", lua.LBool(false))
			L.Push(result)
			return 1
		}
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("failed to stat: %v", err)))
		return 2
	}

	// Get checksum
	checksum, _ := computeChecksum(path)

	result := L.NewTable()
	L.SetField(result, "exists", lua.LBool(true))
	L.SetField(result, "path", lua.LString(path))
	L.SetField(result, "size", lua.LNumber(info.Size()))
	L.SetField(result, "mode", lua.LString(fmt.Sprintf("%o", info.Mode().Perm())))
	L.SetField(result, "is_dir", lua.LBool(info.IsDir()))
	L.SetField(result, "is_file", lua.LBool(info.Mode().IsRegular()))
	L.SetField(result, "mtime", lua.LNumber(info.ModTime().Unix()))
	L.SetField(result, "checksum", lua.LString(checksum))

	// Get owner/group (Unix only)
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		L.SetField(result, "uid", lua.LNumber(stat.Uid))
		L.SetField(result, "gid", lua.LNumber(stat.Gid))
	}

	L.Push(result)
	return 1
}

// Helper functions

func luaValueToGoInterface(v lua.LValue) interface{} {
	switch v.Type() {
	case lua.LTString:
		return v.String()
	case lua.LTNumber:
		return float64(v.(lua.LNumber))
	case lua.LTBool:
		return bool(v.(lua.LBool))
	case lua.LTTable:
		t := v.(*lua.LTable)
		result := make(map[string]interface{})
		t.ForEach(func(k, v lua.LValue) {
			if key, ok := k.(lua.LString); ok {
				result[string(key)] = luaValueToGoInterface(v)
			}
		})
		return result
	default:
		return nil
	}
}

func extractZip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dst, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func extractTarGz(src, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	return extractTarReader(tar.NewReader(gzr), dst)
}

func extractTar(src, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	return extractTarReader(tar.NewReader(file), dst)
}

func extractTarReader(tr *tar.Reader, dst string) error {
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}

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

func computeChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
