# Sloth Runner - Powerful Makefile
# Author: Sloth Runner Team
# Description: Comprehensive build, test, and deployment automation

# ==================== Variables ====================
BINARY_NAME=sloth-runner
MAIN_PATH=./cmd/sloth-runner
BUILD_DIR=./build
DIST_DIR=./dist
INSTALL_PATH=$(HOME)/.local/bin

# Version info
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go build flags
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"
GOFLAGS=-trimpath -mod=readonly

# Platform builds
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# Protobuf paths
PROTO_PATH=proto/agent.proto
PROTOC=/opt/homebrew/bin/protoc
PROTOC_GEN_GO=/Users/chalkan3/go/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC=/Users/chalkan3/go/bin/protoc-gen-go-grpc

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
MAGENTA=\033[0;35m
CYAN=\033[0;36m
NC=\033[0m # No Color

# ==================== Help ====================
.PHONY: help
help: ## Show this help message
	@echo "$(CYAN)╔═══════════════════════════════════════════════════╗$(NC)"
	@echo "$(CYAN)║        Sloth Runner - Makefile Commands         ║$(NC)"
	@echo "$(CYAN)╚═══════════════════════════════════════════════════╝$(NC)"
	@echo ""
	@echo "$(GREEN)📦 Build Commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "$(YELLOW)💡 Quick Examples:$(NC)"
	@echo "  $(CYAN)make build$(NC)           # Build for current platform"
	@echo "  $(CYAN)make install$(NC)         # Build and install to ~/.local/bin"
	@echo "  $(CYAN)make test$(NC)            # Run all tests"
	@echo "  $(CYAN)make docker$(NC)          # Build Docker image"
	@echo "  $(CYAN)make release$(NC)         # Build for all platforms"
	@echo "  $(CYAN)make verify$(NC)          # Run all checks (fmt, vet, lint, test)"
	@echo ""

# ==================== Build Commands ====================
.PHONY: build
build: clean ## 🔨 Build for current platform
	@echo "$(GREEN)🔨 Building $(BINARY_NAME) $(VERSION)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: build-race
build-race: ## 🏃 Build with race detector
	@echo "$(YELLOW)🏃 Building with race detector...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build -race $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-race $(MAIN_PATH)
	@echo "$(GREEN)✅ Race build complete$(NC)"

.PHONY: build-debug
build-debug: ## 🐛 Build with debug symbols
	@echo "$(YELLOW)🐛 Building with debug symbols...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build -gcflags="all=-N -l" -o $(BUILD_DIR)/$(BINARY_NAME)-debug $(MAIN_PATH)
	@echo "$(GREEN)✅ Debug build complete$(NC)"

.PHONY: build-static
build-static: ## 📦 Build static binary (for containers)
	@echo "$(YELLOW)📦 Building static binary...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-static $(MAIN_PATH)
	@echo "$(GREEN)✅ Static build complete$(NC)"

.PHONY: build-all
build-all: ## 🌍 Build for all platforms
	@echo "$(CYAN)🌍 Building for all platforms...$(NC)"
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		output=$(DIST_DIR)/$(BINARY_NAME)-$${OS}-$${ARCH}; \
		if [ "$${OS}" = "windows" ]; then output="$${output}.exe"; fi; \
		echo "$(BLUE)  Building for $${OS}/$${ARCH}...$(NC)"; \
		GOOS=$${OS} GOARCH=$${ARCH} go build $(GOFLAGS) $(LDFLAGS) -o $${output} $(MAIN_PATH); \
		if [ $$? -eq 0 ]; then \
			echo "$(GREEN)  ✅ $${output}$(NC)"; \
		else \
			echo "$(RED)  ❌ Failed to build for $${OS}/$${ARCH}$(NC)"; \
		fi; \
	done
	@echo "$(GREEN)✅ All builds complete!$(NC)"
	@ls -lh $(DIST_DIR)

# ==================== Installation ====================
.PHONY: install
install: build ## 📥 Build and install to ~/.local/bin
	@echo "$(YELLOW)📥 Installing to $(INSTALL_PATH)...$(NC)"
	@mkdir -p $(INSTALL_PATH)
	@rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	@chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "$(GREEN)✅ Installed: $(INSTALL_PATH)/$(BINARY_NAME)$(NC)"
	@ls -lh $(INSTALL_PATH)/$(BINARY_NAME)

.PHONY: uninstall
uninstall: ## 🗑️  Remove from ~/.local/bin
	@echo "$(YELLOW)🗑️  Uninstalling from $(INSTALL_PATH)...$(NC)"
	@rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "$(GREEN)✅ Uninstalled$(NC)"

.PHONY: install-system
install-system: build ## 🔧 Install system-wide (requires sudo)
	@echo "$(YELLOW)🔧 Installing to /usr/local/bin (requires sudo)...$(NC)"
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✅ Installed: /usr/local/bin/$(BINARY_NAME)$(NC)"

# ==================== Testing ====================
.PHONY: test
test: ## 🧪 Run all tests
	@echo "$(CYAN)🧪 Running tests...$(NC)"
	go test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)✅ Tests complete$(NC)"

.PHONY: test-coverage
test-coverage: test ## 📊 Run tests with coverage report
	@echo "$(CYAN)📊 Generating coverage report...$(NC)"
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report: coverage.html$(NC)"
	@open coverage.html 2>/dev/null || echo "Open coverage.html in your browser"

.PHONY: test-bench
test-bench: ## ⚡ Run benchmarks
	@echo "$(CYAN)⚡ Running benchmarks...$(NC)"
	go test -bench=. -benchmem ./...

.PHONY: test-integration
test-integration: ## 🔗 Run integration tests
	@echo "$(CYAN)🔗 Running integration tests...$(NC)"
	go test -v -tags=integration ./...

.PHONY: test-short
test-short: ## ⏱️  Run short tests only
	@echo "$(CYAN)⏱️  Running short tests...$(NC)"
	go test -short -v ./...

# ==================== Code Quality ====================
.PHONY: lint
lint: ## 🔍 Run linters
	@echo "$(CYAN)🔍 Running linters...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
		echo "$(GREEN)✅ Linting complete$(NC)"; \
	else \
		echo "$(RED)❌ golangci-lint not installed. Run: make install-tools$(NC)"; \
		exit 1; \
	fi

.PHONY: fmt
fmt: ## 💅 Format code
	@echo "$(CYAN)💅 Formatting code...$(NC)"
	go fmt ./...
	@echo "$(GREEN)✅ Code formatted$(NC)"

.PHONY: vet
vet: ## 🔬 Run go vet
	@echo "$(CYAN)🔬 Running go vet...$(NC)"
	go vet ./...
	@echo "$(GREEN)✅ Vet complete$(NC)"

.PHONY: tidy
tidy: ## 📦 Tidy go.mod
	@echo "$(CYAN)📦 Tidying dependencies...$(NC)"
	go mod tidy
	@echo "$(GREEN)✅ Dependencies tidied$(NC)"

.PHONY: verify
verify: fmt vet lint test ## ✅ Run all verification steps
	@echo "$(GREEN)✅ All verifications passed!$(NC)"

# ==================== Dependencies ====================
.PHONY: deps
deps: ## 📥 Download dependencies
	@echo "$(CYAN)📥 Downloading dependencies...$(NC)"
	go mod download
	@echo "$(GREEN)✅ Dependencies downloaded$(NC)"

.PHONY: deps-upgrade
deps-upgrade: ## ⬆️  Upgrade all dependencies
	@echo "$(CYAN)⬆️  Upgrading dependencies...$(NC)"
	go get -u ./...
	go mod tidy
	@echo "$(GREEN)✅ Dependencies upgraded$(NC)"

.PHONY: deps-clean
deps-clean: ## 🧹 Clean module cache
	@echo "$(CYAN)🧹 Cleaning module cache...$(NC)"
	go clean -modcache
	@echo "$(GREEN)✅ Module cache cleaned$(NC)"

# ==================== Development ====================
.PHONY: run
run: ## 🚀 Run the application
	@echo "$(CYAN)🚀 Running $(BINARY_NAME)...$(NC)"
	go run $(MAIN_PATH) $(ARGS)

.PHONY: run-master
run-master: ## 🎛️  Run master server
	@echo "$(CYAN)🎛️  Starting master server...$(NC)"
	go run $(MAIN_PATH) master start --port 50053 --bind-address 0.0.0.0

.PHONY: run-agent
run-agent: ## 🤖 Run agent (set AGENT_NAME and MASTER_ADDR)
	@echo "$(CYAN)🤖 Starting agent $(or $(AGENT_NAME),local-agent)...$(NC)"
	go run $(MAIN_PATH) agent start \
		--name $(or $(AGENT_NAME),local-agent) \
		--master $(or $(MASTER_ADDR),localhost:50053) \
		--port 50051

.PHONY: watch
watch: ## 👀 Watch for changes and rebuild
	@echo "$(CYAN)👀 Watching for changes...$(NC)"
	@if command -v watchexec >/dev/null 2>&1; then \
		watchexec -r -e go make build; \
	else \
		echo "$(RED)❌ watchexec not installed. Install with: brew install watchexec$(NC)"; \
		exit 1; \
	fi

# ==================== Docker ====================
.PHONY: docker
docker: ## 🐳 Build Docker image
	@echo "$(CYAN)🐳 Building Docker image...$(NC)"
	docker build -t sloth-runner:$(VERSION) -t sloth-runner:latest .
	@echo "$(GREEN)✅ Docker image built: sloth-runner:$(VERSION)$(NC)"

.PHONY: docker-run
docker-run: ## 🏃 Run Docker container
	@echo "$(CYAN)🏃 Running Docker container...$(NC)"
	docker run --rm -it sloth-runner:latest

.PHONY: docker-push
docker-push: ## 📤 Push Docker image
	@echo "$(CYAN)📤 Pushing Docker image...$(NC)"
	docker push sloth-runner:$(VERSION)
	docker push sloth-runner:latest
	@echo "$(GREEN)✅ Docker image pushed$(NC)"

# ==================== Release ====================
.PHONY: release
release: clean verify build-all ## 🎉 Create release builds
	@echo "$(CYAN)🎉 Creating release $(VERSION)...$(NC)"
	@mkdir -p $(DIST_DIR)/archives
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		binary=$(BINARY_NAME)-$${OS}-$${ARCH}; \
		if [ "$${OS}" = "windows" ]; then binary="$${binary}.exe"; fi; \
		archive=$(DIST_DIR)/archives/$(BINARY_NAME)-$(VERSION)-$${OS}-$${ARCH}.tar.gz; \
		if [ "$${OS}" = "windows" ]; then \
			archive=$(DIST_DIR)/archives/$(BINARY_NAME)-$(VERSION)-$${OS}-$${ARCH}.zip; \
			cd $(DIST_DIR) && zip -q ../archives/$$(basename $${archive}) $$(basename $${binary}) && cd ..; \
		else \
			tar czf $${archive} -C $(DIST_DIR) $$(basename $${binary}); \
		fi; \
		echo "$(GREEN)  ✅ Created $${archive}$(NC)"; \
	done
	@echo "$(GREEN)🎉 Release $(VERSION) complete!$(NC)"
	@ls -lh $(DIST_DIR)/archives

.PHONY: release-notes
release-notes: ## 📝 Generate release notes
	@echo "$(CYAN)📝 Generating release notes...$(NC)"
	@echo "# Release $(VERSION)" > RELEASE_NOTES.md
	@echo "" >> RELEASE_NOTES.md
	@echo "## Changes" >> RELEASE_NOTES.md
	@git log $$(git describe --tags --abbrev=0 2>/dev/null || echo "HEAD")..HEAD --pretty=format:"- %s" >> RELEASE_NOTES.md 2>/dev/null || echo "- Initial release" >> RELEASE_NOTES.md
	@echo "" >> RELEASE_NOTES.md
	@echo "$(GREEN)✅ Release notes generated: RELEASE_NOTES.md$(NC)"

# ==================== Documentation ====================
.PHONY: docs
docs: ## 📚 Build documentation
	@echo "$(CYAN)📚 Building documentation...$(NC)"
	@if [ -d "docs" ] && command -v mkdocs >/dev/null 2>&1; then \
		cd docs && mkdocs build; \
		echo "$(GREEN)✅ Documentation built$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  mkdocs not installed or no docs directory$(NC)"; \
	fi

.PHONY: docs-serve
docs-serve: ## 🌐 Serve documentation locally
	@echo "$(CYAN)🌐 Serving documentation at http://localhost:8000$(NC)"
	@if [ -d "docs" ] && command -v mkdocs >/dev/null 2>&1; then \
		cd docs && mkdocs serve; \
	else \
		echo "$(RED)❌ mkdocs not installed$(NC)"; \
	fi

.PHONY: godoc
godoc: ## 📖 Generate Go documentation
	@echo "$(CYAN)📖 Starting godoc server at http://localhost:6060$(NC)"
	@godoc -http=:6060

# ==================== Database ====================
.PHONY: db-migrate
db-migrate: ## 🗄️  Run database migrations
	@echo "$(CYAN)🗄️  Running database migrations...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) db migrate
	@echo "$(GREEN)✅ Migrations complete$(NC)"

.PHONY: db-reset
db-reset: ## 🔄 Reset database
	@echo "$(YELLOW)🔄 Resetting database...$(NC)"
	@rm -f .sloth-cache/*.db
	@echo "$(GREEN)✅ Database reset$(NC)"

# ==================== Protobuf ====================
.PHONY: proto
proto: ## 🔧 Generate protobuf code
	@echo "$(CYAN)🔧 Generating protobuf code...$(NC)"
	@if [ -f "$(PROTOC)" ]; then \
		PATH=$(dir $(PROTOC_GEN_GO)):$$PATH $(PROTOC) \
			--go_out=paths=source_relative:. \
			--go-grpc_out=paths=source_relative:. \
			$(PROTO_PATH); \
		echo "$(GREEN)✅ Protobuf code generated$(NC)"; \
	else \
		echo "$(RED)❌ protoc not found at $(PROTOC)$(NC)"; \
		exit 1; \
	fi

# ==================== Utilities ====================
.PHONY: clean
clean: ## 🧹 Clean build artifacts
	@echo "$(YELLOW)🧹 Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR) $(DIST_DIR) coverage.out coverage.html
	@go clean -cache -testcache
	@echo "$(GREEN)✅ Clean complete$(NC)"

.PHONY: clean-all
clean-all: clean db-reset ## 🗑️  Deep clean (including databases)
	@echo "$(YELLOW)🗑️  Deep cleaning...$(NC)"
	@go clean -modcache
	@echo "$(GREEN)✅ Deep clean complete$(NC)"

.PHONY: size
size: build ## 📏 Show binary size
	@echo "$(CYAN)📏 Binary size analysis:$(NC)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)
	@file $(BUILD_DIR)/$(BINARY_NAME)
	@if command -v du >/dev/null 2>&1; then \
		echo "$(BLUE)Detailed:$(NC)"; \
		du -h $(BUILD_DIR)/$(BINARY_NAME); \
	fi

.PHONY: version
version: ## ℹ️  Show version information
	@echo "$(CYAN)ℹ️  Version Information:$(NC)"
	@echo "  Version:    $(VERSION)"
	@echo "  Commit:     $(COMMIT)"
	@echo "  Build Time: $(BUILD_TIME)"
	@echo "  Go Version: $$(go version)"

.PHONY: info
info: version size ## 📋 Show project information
	@echo ""
	@echo "$(CYAN)📋 Project Information:$(NC)"
	@echo "  Binary:     $(BINARY_NAME)"
	@echo "  Main Path:  $(MAIN_PATH)"
	@echo "  Install:    $(INSTALL_PATH)"

# ==================== Tools ====================
.PHONY: install-tools
install-tools: ## 🔧 Install development tools
	@echo "$(CYAN)🔧 Installing development tools...$(NC)"
	@echo "$(BLUE)  Installing golangci-lint...$(NC)"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(BLUE)  Installing gofumpt...$(NC)"
	@go install mvdan.cc/gofumpt@latest
	@echo "$(BLUE)  Installing staticcheck...$(NC)"
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "$(GREEN)✅ Tools installed$(NC)"

.PHONY: check-tools
check-tools: ## 🔍 Check if required tools are installed
	@echo "$(CYAN)🔍 Checking required tools...$(NC)"
	@for tool in go git docker golangci-lint; do \
		if command -v $$tool >/dev/null 2>&1; then \
			echo "$(GREEN)  ✅ $$tool$(NC)"; \
		else \
			echo "$(RED)  ❌ $$tool (not found)$(NC)"; \
		fi; \
	done

# ==================== Git ====================
.PHONY: git-clean
git-clean: ## 🧹 Clean git repository (remove untracked files)
	@echo "$(YELLOW)🧹 Cleaning git repository...$(NC)"
	git clean -fd
	git checkout .
	@echo "$(GREEN)✅ Git repository cleaned$(NC)"

.PHONY: git-status
git-status: ## 📊 Show git status
	@git status

.PHONY: tag
tag: ## 🏷️  Create git tag (use TAG=v1.0.0)
	@if [ -z "$(TAG)" ]; then \
		echo "$(RED)❌ Error: TAG not set. Use: make tag TAG=v1.0.0$(NC)"; \
		exit 1; \
	fi
	@echo "$(CYAN)🏷️  Creating tag $(TAG)...$(NC)"
	@git tag -a $(TAG) -m "Release $(TAG)"
	@git push origin $(TAG)
	@echo "$(GREEN)✅ Tag $(TAG) created and pushed$(NC)"

# ==================== CI/CD ====================
.PHONY: ci
ci: verify build ## 🔄 Run CI pipeline locally
	@echo "$(GREEN)✅ CI pipeline complete!$(NC)"

.PHONY: pre-commit
pre-commit: fmt vet lint test-short ## ✅ Run pre-commit checks
	@echo "$(GREEN)✅ Pre-commit checks passed!$(NC)"

.PHONY: pre-push
pre-push: verify ## ✅ Run pre-push checks
	@echo "$(GREEN)✅ Pre-push checks passed!$(NC)"

# ==================== Examples ====================
.PHONY: example
example: build ## 🎯 Run example workflow
	@echo "$(CYAN)🎯 Running example workflow...$(NC)"
	@$(BUILD_DIR)/$(BINARY_NAME) run example -f examples/basic_pipeline_modern.sloth

.PHONY: demo
demo: install ## 🎬 Run demo
	@echo "$(CYAN)🎬 Running demo...$(NC)"
	@$(INSTALL_PATH)/$(BINARY_NAME) --version
	@$(INSTALL_PATH)/$(BINARY_NAME) agent list

# ==================== Performance ====================
.PHONY: profile-cpu
profile-cpu: ## 🔥 Profile CPU usage
	@echo "$(CYAN)🔥 Profiling CPU...$(NC)"
	go test -cpuprofile=cpu.prof -bench=. ./...
	go tool pprof -http=:8080 cpu.prof

.PHONY: profile-mem
profile-mem: ## 💾 Profile memory usage
	@echo "$(CYAN)💾 Profiling memory...$(NC)"
	go test -memprofile=mem.prof -bench=. ./...
	go tool pprof -http=:8080 mem.prof

# ==================== Security ====================
.PHONY: security
security: ## 🔒 Run security checks
	@echo "$(CYAN)🔒 Running security checks...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
		echo "$(GREEN)✅ Security scan complete$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  gosec not installed. Run: go install github.com/securego/gosec/v2/cmd/gosec@latest$(NC)"; \
	fi

.PHONY: vuln
vuln: ## 🛡️  Check for vulnerabilities
	@echo "$(CYAN)🛡️  Checking for vulnerabilities...$(NC)"
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
		echo "$(GREEN)✅ Vulnerability check complete$(NC)"; \
	else \
		echo "$(YELLOW)⚠️  govulncheck not installed. Run: go install golang.org/x/vuln/cmd/govulncheck@latest$(NC)"; \
	fi

# ==================== Quick Commands ====================
.PHONY: quick-build
quick-build: ## ⚡ Quick build and install
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH) && \
	rm -f $(INSTALL_PATH)/$(BINARY_NAME) && \
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/ && \
	echo "$(GREEN)⚡ Quick install complete!$(NC)"

.PHONY: dev
dev: quick-build ## 🚀 Quick dev cycle (build + install)
	@echo "$(GREEN)🚀 Ready for development!$(NC)"

# ==================== Default ====================
.DEFAULT_GOAL := help
