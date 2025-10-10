# Docker Module

The `docker` module provides a convenient interface for interacting with the Docker daemon, allowing you to build, run, and push Docker images as part of your pipelines.

## Configuration

This module requires the `docker` CLI to be installed and the Docker daemon to be running and accessible.

## Functions

### `docker.exec(args)`

Executes any raw `docker` command.

- `args` (table): **Required.** A list of arguments to pass to the `docker` command (e.g., `{"ps", "-a"}`).
- **Returns:** A result table with `success`, `stdout`, `stderr`, and `exit_code`.

### `docker.build(params)`

Builds a Docker image using `docker build`.

- `params` (table):
    - `tag` (string): **Required.** The tag for the image (e.g., `my-app:latest`).
    - `path` (string): **Required.** The build context path.
    - `dockerfile` (string): **Optional.** The path to the Dockerfile.
    - `build_args` (table): **Optional.** A table of build arguments (e.g., `{VERSION = "1.0"}`).
- **Returns:** A result table.

### `docker.push(params)`

Pushes a Docker image to a registry using `docker push`.

- `params` (table):
    - `tag` (string): **Required.** The tag of the image to push.
- **Returns:** A result table.

### `docker.run(params)`

Runs a Docker container using `docker run`.

- `params` (table):
    - `image` (string): **Required.** The image to run.
    - `name` (string): **Optional.** The name for the container.
    - `detach` (boolean): **Optional.** If `true`, runs the container in the background (`-d`).
    - `ports` (table): **Optional.** A list of port mappings (e.g., `{"8080:80"}`).
    - `env` (table): **Optional.** A table of environment variables (e.g., `{MY_VAR = "value"}`).
- **Returns:** A result table.

## Example

```lua
local image_tag = "my-test-image:latest"

-- Task 1: Build Docker image
task("docker-build")
  :description("Build Docker image")
  :command(function(this, params)
    local result = docker.build({
      tag = image_tag,
      path = "./app"
    })
    if not result.success then
      return false, "Build failed: " .. (result.stderr or "")
    end
    return true, "Docker image built successfully"
  end)
  :build()

-- Task 2: Run Docker container
task("docker-run")
  :description("Run Docker container")
  :command(function(this, params)
    local result = docker.run({
      image = image_tag,
      name = "my-test-container",
      ports = {"8080:80"}
    })
    if not result.success then
      return false, "Run failed: " .. (result.stderr or "")
    end
    return true, "Docker container started successfully"
  end)
  :build()

-- Task 3: Push Docker image
task("docker-push")
  :description("Push Docker image to registry")
  :command(function(this, params)
    local result = docker.push({tag = image_tag})
    if not result.success then
      return false, "Push failed: " .. (result.stderr or "")
    end
    return true, "Docker image pushed successfully"
  end)
  :build()

-- Workflow: Complete Docker pipeline
workflow.define("docker-pipeline")
  :description("Build, run, and push Docker image")
  :version("1.0.0")
  :tasks({
    "docker-build",
    "docker-run",
    "docker-push"
  })
```
