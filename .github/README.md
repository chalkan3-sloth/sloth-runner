# 🚀 GitHub Actions Workflows

This directory contains the CI/CD pipelines for the Sloth Runner project.

## 📋 Available Workflows

### 🔧 `ci.yml` - Main CI Pipeline (build-and-test)
**Triggers:** Push/PR to master
**Purpose:** Primary build and test pipeline
- ✅ Tests with SQLite support on Linux
- ✅ Multi-architecture builds (Linux AMD64/ARM64, macOS, Windows)
- ✅ Coverage reporting
- ✅ Artifact generation

### 📚 `pages.yml` - Documentation Deployment
**Triggers:** Push to master, manual dispatch
**Purpose:** Deploy documentation to GitHub Pages
- ✅ MkDocs Material theme
- ✅ Multi-language support (EN/PT/ZH)
- ✅ Automatic deployment to GitHub Pages

### 📦 `release.yml` - Release Pipeline
**Triggers:** Tag push (v*)
**Purpose:** Create GitHub releases with binaries
- ✅ GoReleaser with cross-compilation
- ✅ SQLite-enabled Linux binaries
- ✅ Static binaries for other platforms
- ✅ Automated changelog generation

### 🔄 `tag-on-merge.yml` - Auto Tagging
**Triggers:** PR merge to master
**Purpose:** Automatically create semantic version tags
- ✅ Semantic versioning
- ✅ Automated tag creation on merge

### 📖 `docs.yml` - Documentation Testing
**Triggers:** Changes to docs/ or mkdocs.yml
**Purpose:** Test documentation builds
- ✅ Validate MkDocs configuration
- ✅ Test multilingual builds
- ✅ Link validation

### 🌐 `multi-platform.yml` - Multi-Platform Testing
**Triggers:** Code changes (excludes docs)
**Purpose:** Test builds across platforms
- ✅ Ubuntu, macOS, Windows testing
- ✅ Platform-specific SQLite handling
- ✅ Binary validation

## 🔧 Technical Requirements

### Go Version
- **Required:** Go 1.24.0+
- **Reason:** SQLite CGO requirements and modern Go features

### SQLite Support
- **Linux:** CGO_ENABLED=1 with SQLite libraries
- **macOS/Windows:** CGO_ENABLED=0 (SQLite disabled for portability)
- **Dependencies:** build-essential, sqlite3, libsqlite3-dev

### Documentation
- **Python:** 3.x
- **Dependencies:** See `requirements.txt`
- **Plugins:** Material theme, multilang, awesome pages, etc.

## 📊 Build Matrix

| Platform | CGO | SQLite | Cross-Compile |
|----------|-----|--------|---------------|
| Linux AMD64 | ✅ | ✅ | gcc |
| Linux ARM64 | ✅ | ✅ | aarch64-linux-gnu-gcc |
| macOS AMD64 | ❌ | ❌ | Static |
| macOS ARM64 | ❌ | ❌ | Static |
| Windows AMD64 | ❌ | ❌ | Static |

## 🚨 Important Notes

1. **SQLite State Module**: Only works on Linux builds due to CGO requirements
2. **Cross-Compilation**: ARM64 Linux requires specific GCC toolchain
3. **Documentation**: Requires multiple Python packages for full functionality
4. **Caching**: Go modules and Python packages are cached for performance
5. **Security**: All workflows use pinned action versions for security

## 🔍 Troubleshooting

### Common Issues:
- **Go Version Mismatch:** Ensure all workflows use Go 1.24.0+
- **SQLite Build Fails:** Check CGO_ENABLED and system dependencies
- **Documentation Build Fails:** Verify all Python packages in requirements.txt
- **Cross-Compile Fails:** Ensure correct GCC toolchain for target architecture

### Debug Tips:
- Check individual workflow logs in GitHub Actions
- Test locally with same Go/Python versions
- Verify dependencies are properly cached
- Check artifact uploads for build outputs

## 📈 Performance Optimizations

- **Caching:** Go modules and Python packages cached between runs
- **Parallel Jobs:** Tests and builds run concurrently where possible
- **Conditional Triggers:** Docs only build when relevant files change
- **Artifact Retention:** Short retention periods for test artifacts

---

For more information, see the main [README.md](../README.md) and [documentation](https://chalkan3.github.io/sloth-runner/).