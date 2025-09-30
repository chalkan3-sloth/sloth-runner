# 🔌 Plugin Development

**Building Extensions for Sloth Runner Platform**

Sloth Runner provides a powerful plugin system that allows developers to extend the platform with custom functionality. This guide covers everything you need to know to develop your own plugins.

## 🏗️ Plugin Architecture

### Plugin Types

Sloth Runner supports several types of plugins:

1. **🌙 Lua Modules** - Extend the DSL with new functions and capabilities
2. **⚡ Command Processors** - Add new CLI commands and operations
3. **🎨 UI Extensions** - Enhance the web dashboard and interface
4. **🔗 Integrations** - Connect with external tools and services
5. **🦥 Editor Plugins** - IDE/Editor extensions (like our Neovim plugin)

### Core Components

```
sloth-runner/
├── plugins/
│   ├── lua-modules/       # Lua DSL extensions
│   ├── commands/          # CLI command plugins
│   ├── ui/               # Web UI extensions
│   ├── integrations/     # Third-party integrations
│   └── editors/          # Editor/IDE plugins
└── internal/
    └── plugin/           # Plugin system core
```

## 🌙 Developing Lua Module Plugins

### Basic Structure

Create a new Lua module that extends the DSL:

```lua
-- plugins/lua-modules/my-module/init.lua
local M = {}

-- Module metadata
M._NAME = "my-module"
M._VERSION = "1.0.0"
M._DESCRIPTION = "Custom functionality for Sloth Runner"

-- Public API
function M.hello(name)
    return string.format("Hello, %s from my custom module!", name or "World")
end

function M.custom_task(config)
    return {
        execute = function(params)
            log.info("🔌 Executing custom task: " .. config.name)
            -- Custom task logic here
            return true
        end,
        validate = function()
            return config.name ~= nil
        end
    }
end

-- Register module functions
function M.register()
    -- Make functions available in DSL
    _G.my_module = M
    
    -- Register custom task type
    task.register_type("custom", M.custom_task)
end

return M
```

### Using Custom Modules in Workflows

```lua
-- workflow.sloth
local my_task = task("test_custom")
    :type("custom", { name = "test" })
    :description("Testing custom plugin")
    :build()

-- Direct module usage
local greeting = my_module.hello("Developer")
log.info(greeting)

workflow.define("plugin_test", {
    description = "Testing custom plugin",
    tasks = { my_task }
})
```

### Plugin Registration

Create a plugin manifest:

```yaml
# plugins/lua-modules/my-module/plugin.yaml
name: my-module
version: 1.0.0
description: Custom functionality for Sloth Runner
type: lua-module
author: Your Name
license: MIT

entry_point: init.lua
dependencies:
  - sloth-runner: ">=1.0.0"

permissions:
  - filesystem.read
  - network.http
  - system.exec

configuration:
  settings:
    api_key:
      type: string
      required: false
      description: "API key for external service"
```

## ⚡ Command Plugin Development

### CLI Command Structure

```go
// plugins/commands/my-command/main.go
package main

import (
    "github.com/spf13/cobra"
    "github.com/chalkan3-sloth/sloth-runner/pkg/plugin"
)

type MyCommandPlugin struct {
    config *MyConfig
}

type MyConfig struct {
    Setting1 string `json:"setting1"`
    Setting2 int    `json:"setting2"`
}

func (p *MyCommandPlugin) Name() string {
    return "my-command"
}

func (p *MyCommandPlugin) Command() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "my-command",
        Short: "Custom command functionality",
        Long:  "Extended description of custom command",
        RunE:  p.execute,
    }
    
    cmd.Flags().StringVar(&p.config.Setting1, "setting1", "", "Custom setting")
    cmd.Flags().IntVar(&p.config.Setting2, "setting2", 0, "Another setting")
    
    return cmd
}

func (p *MyCommandPlugin) execute(cmd *cobra.Command, args []string) error {
    log.Info("🔌 Executing custom command with settings:", 
        "setting1", p.config.Setting1,
        "setting2", p.config.Setting2)
    
    // Custom command logic here
    return nil
}

func main() {
    plugin := &MyCommandPlugin{
        config: &MyConfig{},
    }
    
    plugin.Register()
}
```

### Command Plugin Manifest

```yaml
# plugins/commands/my-command/plugin.yaml
name: my-command
version: 1.0.0
description: Custom CLI command for Sloth Runner
type: command
author: Your Name

build:
  binary: my-command
  source: main.go

installation:
  target: commands/my-command
```

## 🎨 UI Extension Development

### React Component Plugin

```typescript
// plugins/ui/my-dashboard/src/MyDashboardPlugin.tsx
import React from 'react';
import { PluginComponent, useSlothApi } from '@sloth-runner/ui-sdk';

interface MyDashboardProps {
  config: {
    title: string;
    refreshInterval: number;
  };
}

export const MyDashboardPlugin: PluginComponent<MyDashboardProps> = ({ config }) => {
  const { data, loading } = useSlothApi('/api/custom-metrics');

  return (
    <div className="my-dashboard-plugin">
      <h2>{config.title}</h2>
      {loading ? (
        <div>Loading custom metrics...</div>
      ) : (
        <div className="metrics-grid">
          {data?.map((metric: any) => (
            <div key={metric.id} className="metric-card">
              <h3>{metric.name}</h3>
              <div className="metric-value">{metric.value}</div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

// Plugin registration
export const plugin = {
  name: 'my-dashboard',
  version: '1.0.0',
  component: MyDashboardPlugin,
  defaultConfig: {
    title: 'Custom Dashboard',
    refreshInterval: 30000,
  },
};
```

### UI Plugin Manifest

```yaml
# plugins/ui/my-dashboard/plugin.yaml
name: my-dashboard
version: 1.0.0
description: Custom dashboard for Sloth Runner
type: ui-extension
author: Your Name

build:
  framework: react
  entry: src/index.tsx
  output: dist/

installation:
  target: ui/plugins/my-dashboard
  
dependencies:
  - "@sloth-runner/ui-sdk": "^1.0.0"
  - "react": "^18.0.0"
```

## 🔗 Integration Plugin Development

### External Service Integration

```go
// plugins/integrations/my-service/integration.go
package main

import (
    "context"
    "net/http"
    "github.com/chalkan3-sloth/sloth-runner/pkg/integration"
)

type MyServiceIntegration struct {
    client *http.Client
    apiKey string
}

func (i *MyServiceIntegration) Name() string {
    return "my-service"
}

func (i *MyServiceIntegration) Initialize(config map[string]interface{}) error {
    i.apiKey = config["api_key"].(string)
    i.client = &http.Client{}
    return nil
}

func (i *MyServiceIntegration) GetMetrics(ctx context.Context) ([]integration.Metric, error) {
    // Fetch metrics from external service
    metrics := []integration.Metric{
        {
            Name:  "custom_metric",
            Value: 42,
            Tags:  map[string]string{"source": "my-service"},
        },
    }
    return metrics, nil
}

func (i *MyServiceIntegration) SendNotification(ctx context.Context, msg integration.Message) error {
    // Send notification via external service
    return nil
}

func main() {
    integration := &MyServiceIntegration{}
    integration.Register()
}
```

## 🛠️ Plugin Development Tools

### Plugin Generator

Create new plugins quickly with the generator:

```bash
# Generate a new Lua module plugin
sloth-runner plugin generate --type=lua-module --name=my-module

# Generate a CLI command plugin
sloth-runner plugin generate --type=command --name=my-command

# Generate a UI extension
sloth-runner plugin generate --type=ui --name=my-dashboard
```

### Development Environment

```bash
# Start development server with plugin hot-reload
sloth-runner dev --plugins-dir=./plugins

# Test plugin locally
sloth-runner plugin test ./plugins/my-plugin

# Build plugin for distribution
sloth-runner plugin build ./plugins/my-plugin --output=dist/
```

### Plugin Testing

```go
// plugins/my-plugin/plugin_test.go
package main

import (
    "testing"
    "github.com/chalkan3-sloth/sloth-runner/pkg/plugin/testing"
)

func TestMyPlugin(t *testing.T) {
    // Create test environment
    env := plugintest.NewEnvironment(t)
    
    // Load plugin
    plugin, err := env.LoadPlugin("./")
    if err != nil {
        t.Fatal(err)
    }
    
    // Test plugin functionality
    result, err := plugin.Execute(map[string]interface{}{
        "test_param": "value",
    })
    
    if err != nil {
        t.Fatal(err)
    }
    
    // Verify results
    if result.Status != "success" {
        t.Errorf("Expected success, got %s", result.Status)
    }
}
```

## 📦 Plugin Distribution

### Plugin Registry

Publish your plugin to the Sloth Runner plugin registry:

```bash
# Login to registry
sloth-runner registry login

# Publish plugin
sloth-runner plugin publish ./my-plugin

# Install published plugin
sloth-runner plugin install my-username/my-plugin
```

### Plugin Marketplace

Browse and discover plugins:

```bash
# Search plugins
sloth-runner plugin search "kubernetes"

# Get plugin info
sloth-runner plugin info kubernetes-operator

# Install from marketplace
sloth-runner plugin install --marketplace kubernetes-operator
```

## 🔒 Security & Best Practices

### Security Guidelines

1. **🛡️ Principle of Least Privilege** - Request only necessary permissions
2. **🔐 Input Validation** - Always validate user input and configuration
3. **🚫 Avoid Global State** - Keep plugin state isolated
4. **📝 Error Handling** - Provide clear error messages and logging
5. **🧪 Testing** - Write comprehensive tests for all functionality

### Code Quality

```go
// Good: Clear error handling
func (p *MyPlugin) Execute(params map[string]interface{}) (*Result, error) {
    value, ok := params["required_param"].(string)
    if !ok {
        return nil, fmt.Errorf("required_param must be a string")
    }
    
    if value == "" {
        return nil, fmt.Errorf("required_param cannot be empty")
    }
    
    // Process with validated input
    result := p.process(value)
    return result, nil
}
```

### Documentation Standards

Every plugin should include:

- **📋 README.md** - Installation and usage instructions
- **📚 API Documentation** - Function/method documentation
- **📖 Examples** - Working code examples
- **🧪 Tests** - Unit and integration tests
- **📄 License** - Clear licensing information

## 🚀 Advanced Plugin Features

### Plugin Hooks

```lua
-- Respond to system events
function M.on_task_start(task_id, context)
    log.info("🔌 Task starting: " .. task_id)
    -- Custom logic before task execution
end

function M.on_task_complete(task_id, result)
    log.info("🔌 Task completed: " .. task_id)
    -- Custom logic after task completion
end

-- Register hooks
M.hooks = {
    ["task.start"] = M.on_task_start,
    ["task.complete"] = M.on_task_complete,
}
```

### Plugin Communication

```lua
-- Inter-plugin communication
function M.send_to_plugin(plugin_name, message)
    local plugin = sloth.plugins.get(plugin_name)
    if plugin and plugin.receive_message then
        return plugin.receive_message(message)
    end
    return nil
end

function M.receive_message(message)
    log.info("🔌 Received message: " .. message.type)
    -- Handle incoming message
    return { status = "received" }
end
```

### Configuration Management

```yaml
# plugins/my-plugin/config.schema.json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "api_endpoint": {
      "type": "string",
      "format": "uri",
      "description": "API endpoint URL"
    },
    "timeout": {
      "type": "integer",
      "minimum": 1,
      "maximum": 300,
      "default": 30
    }
  },
  "required": ["api_endpoint"]
}
```

## 📚 Examples & Templates

### Complete Plugin Example

Check out these example plugins:

- **[Kubernetes Operator Plugin](https://github.com/sloth-runner/plugin-kubernetes)** - Manage K8s resources
- **[Slack Integration Plugin](https://github.com/sloth-runner/plugin-slack)** - Send notifications
- **[Monitoring Dashboard Plugin](https://github.com/sloth-runner/plugin-monitoring)** - Custom metrics UI

### Plugin Templates

Use official templates for quick starts:

```bash
# Use template
sloth-runner plugin init --template=lua-module my-plugin
sloth-runner plugin init --template=go-command my-command
sloth-runner plugin init --template=react-ui my-dashboard
```

## 💬 Community & Support

### Getting Help

- **📖 [Plugin API Documentation](https://docs.sloth-runner.io/plugin-api)**
- **💬 [Discord Community](https://discord.gg/sloth-runner)** - #plugin-development
- **🐛 [GitHub Issues](https://github.com/chalkan3-sloth/sloth-runner/issues)** - Bug reports and feature requests
- **📧 [Plugin Mailing List](mailto:plugins@sloth-runner.io)** - Development discussions

### Contributing

We welcome plugin contributions! See our [Contributing Guide](contributing.md) for details on:

- Plugin submission process
- Code review guidelines
- Documentation requirements
- Testing standards

---

Start building amazing plugins for Sloth Runner today! The platform's extensible architecture makes it easy to add exactly the functionality you need. 🔌✨