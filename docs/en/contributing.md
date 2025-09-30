# 🤝 Contributing to Sloth Runner

**Thank you for your interest in contributing to Sloth Runner!** 

We welcome contributions from developers of all skill levels. Whether you're fixing bugs, adding features, improving documentation, or creating plugins, your help makes Sloth Runner better for everyone.

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+** for core development
- **Node.js 18+** for UI development  
- **Lua 5.4+** for DSL development
- **Git** for version control

### Development Setup

```bash
# Clone the repository
git clone https://github.com/chalkan3-sloth/sloth-runner.git
cd sloth-runner

# Install dependencies
go mod download
npm install  # for UI components

# Run tests
make test

# Build the project
make build
```

## 📋 Ways to Contribute

### 🐛 Bug Reports

Found a bug? Please help us fix it:

1. **Search existing issues** to avoid duplicates
2. **Use our bug report template** with:
   - Sloth Runner version
   - Operating system
   - Steps to reproduce
   - Expected vs actual behavior
   - Error logs (if any)

### 💡 Feature Requests

Have an idea for improvement?

1. **Check the roadmap** for planned features
2. **Open a feature request** with:
   - Clear description of the feature
   - Use cases and benefits
   - Possible implementation approach

### 🔧 Code Contributions

Ready to code? Here's how:

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Make your changes** following our coding standards
4. **Add tests** for new functionality
5. **Update documentation** if needed
6. **Commit with clear messages**
7. **Push and create a Pull Request**

### 📚 Documentation

Help improve our docs:

- Fix typos and unclear explanations
- Add examples and tutorials
- Translate content to other languages
- Update API documentation

### 🔌 Plugin Development

Create plugins for the community:

- Follow our [Plugin Development Guide](plugin-development.md)
- Submit to the plugin registry
- Maintain compatibility with core versions

## 📐 Development Guidelines

### Code Style

#### Go Code

Follow standard Go conventions:

```go
// Good: Clear function names and comments
func ProcessWorkflowTasks(ctx context.Context, workflow *Workflow) error {
    if workflow == nil {
        return fmt.Errorf("workflow cannot be nil")
    }
    
    for _, task := range workflow.Tasks {
        if err := processTask(ctx, task); err != nil {
            return fmt.Errorf("failed to process task %s: %w", task.ID, err)
        }
    }
    
    return nil
}
```

#### Lua DSL

Keep DSL code clean and readable:

```lua
-- Good: Clear task definition with proper chaining
local deploy_task = task("deploy_application")
    :description("Deploy the application to production")
    :command(function(params, deps)
        local result = exec.run("kubectl apply -f deployment.yaml")
        if not result.success then
            log.error("Deployment failed: " .. result.stderr)
            return false
        end
        return true
    end)
    :timeout(300)
    :retries(3)
    :build()
```

#### TypeScript/JavaScript

For UI components:

```typescript
// Good: Proper typing and error handling
interface TaskResult {
  id: string;
  status: 'success' | 'failed' | 'running';
  duration: number;
}

export const TaskStatusCard: React.FC<{ result: TaskResult }> = ({ result }) => {
  const statusColor = result.status === 'success' ? 'green' : 
                     result.status === 'failed' ? 'red' : 'blue';
  
  return (
    <div className={`task-card status-${result.status}`}>
      <h3>{result.id}</h3>
      <span style={{ color: statusColor }}>{result.status}</span>
      <small>{result.duration}ms</small>
    </div>
  );
};
```

### Testing Standards

#### Unit Tests

Write tests for all new functionality:

```go
func TestProcessWorkflowTasks(t *testing.T) {
    tests := []struct {
        name     string
        workflow *Workflow
        wantErr  bool
    }{
        {
            name:     "nil workflow should return error",
            workflow: nil,
            wantErr:  true,
        },
        {
            name: "valid workflow should process successfully",
            workflow: &Workflow{
                Tasks: []*Task{{ID: "test-task"}},
            },
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ProcessWorkflowTasks(context.Background(), tt.workflow)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessWorkflowTasks() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

#### Integration Tests

Test real-world scenarios:

```bash
# Run integration tests
make test-integration

# Test with different configurations
make test-configs
```

### Documentation Standards

- **Keep it simple** - Use clear, concise language
- **Include examples** - Show don't just tell
- **Update with changes** - Keep docs in sync with code
- **Test examples** - Ensure all code examples work

## 🔄 Pull Request Process

### Before Submitting

- [ ] **Run tests** - `make test`
- [ ] **Run linting** - `make lint`
- [ ] **Update docs** - If adding/changing features
- [ ] **Add changelog entry** - In `CHANGELOG.md`
- [ ] **Check compatibility** - With existing features

### PR Template

Use our pull request template:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] Changelog updated
```

### Review Process

1. **Automated checks** run on all PRs
2. **Maintainer review** for code quality and design
3. **Community feedback** welcomed on all PRs
4. **Approval and merge** by maintainers

## 🏗️ Project Structure

Understanding the codebase:

```
sloth-runner/
├── cmd/                    # CLI commands
├── internal/              # Internal packages
│   ├── core/             # Core business logic
│   ├── dsl/              # DSL implementation
│   ├── execution/        # Task execution engine
│   └── plugins/          # Plugin system
├── pkg/                   # Public packages
├── plugins/              # Built-in plugins
├── docs/                 # Documentation
├── web/                  # Web UI components
└── examples/             # Example workflows
```

## 🎯 Contribution Areas

### High Priority

- **🐛 Bug fixes** - Always welcome
- **📈 Performance improvements** - Optimization opportunities
- **🧪 Test coverage** - Increase test coverage
- **📚 Documentation** - Keep docs comprehensive

### Medium Priority

- **✨ New features** - Following roadmap priorities
- **🔌 Plugin ecosystem** - More plugins and integrations
- **🎨 UI improvements** - Better user experience

### Future Areas

- **🤖 AI enhancements** - Advanced ML capabilities  
- **☁️ Cloud integrations** - More cloud provider support
- **📊 Analytics** - Better insights and reporting

## 🏆 Recognition

Contributors are recognized in:

- **CONTRIBUTORS.md** - All contributors listed
- **Release notes** - Major contributions highlighted
- **Community showcase** - Featured contributions
- **Contributor badges** - GitHub profile recognition

## 📞 Getting Help

### Development Questions

- **💬 Discord** - `#development` channel
- **📧 Mailing List** - dev@sloth-runner.io
- **📖 Wiki** - Development guides and FAQs

### Mentorship

New to open source? We offer mentorship:

- **👥 Mentor matching** - Paired with experienced contributors
- **📚 Learning resources** - Curated learning materials
- **🎯 Guided contributions** - Starter-friendly issues

## 📜 Code of Conduct

We are committed to providing a welcoming and inclusive environment. Please read our [Code of Conduct](https://github.com/chalkan3-sloth/sloth-runner/blob/main/CODE_OF_CONDUCT.md).

### Our Standards

- **🤝 Be respectful** - Treat everyone with respect
- **💡 Be constructive** - Provide helpful feedback
- **🌍 Be inclusive** - Welcome diverse perspectives
- **📚 Be patient** - Help others learn and grow

## 🚀 Release Process

Understanding our releases:

- **🔄 Continuous Integration** - Automated testing and building
- **📅 Regular Releases** - Monthly feature releases
- **🚨 Hotfixes** - Critical bugs fixed immediately
- **📊 Semantic Versioning** - Clear version numbering

## 📈 Roadmap Participation

Help shape the future:

- **📋 Feature Planning** - Participate in roadmap discussions
- **🗳️ Voting** - Vote on feature priorities
- **💭 RFC Process** - Propose major changes through RFCs

---

**Ready to contribute?** 

Start by exploring our [Good First Issues](https://github.com/chalkan3-sloth/sloth-runner/labels/good%20first%20issue) or join our [Discord community](https://discord.gg/sloth-runner) to introduce yourself!

Thank you for helping make Sloth Runner better! 🦥✨