# Terraform Module

The `terraform` module provides a high-level interface for orchestrating `terraform` CLI commands, allowing you to manage your infrastructure lifecycle directly from within a Sloth-Runner pipeline.

## Configuration

This module requires the `terraform` CLI to be installed and available in the system's PATH. All commands must be executed within a specific `workdir` where your `.tf` files are located.

## Functions

### `terraform.init(params)`

Initializes a Terraform working directory.

- `params` (table):
    - `workdir` (string): **Required.** The path to the directory containing the Terraform files.
- **Returns:** A result table with `success`, `stdout`, `stderr`, and `exit_code`.

### `terraform.plan(params)`

Creates a Terraform execution plan.

- `params` (table):
    - `workdir` (string): **Required.** The path to the directory.
    - `out` (string): **Optional.** The filename to save the generated plan to.
- **Returns:** A result table.

### `terraform.apply(params)`

Applies a Terraform plan.

- `params` (table):
    - `workdir` (string): **Required.** The path to the directory.
    - `plan` (string): **Optional.** The path to a plan file to apply.
    - `auto_approve` (boolean): **Optional.** If `true`, applies changes without interactive approval.
- **Returns:** A result table.

### `terraform.destroy(params)`

Destroys Terraform-managed infrastructure.

- `params` (table):
    - `workdir` (string): **Required.** The path to the directory.
    - `auto_approve` (boolean): **Optional.** If `true`, destroys resources without interactive approval.
- **Returns:** A result table.

### `terraform.output(params)`

Reads an output variable from a Terraform state file.

- `params` (table):
    - `workdir` (string): **Required.** The path to the directory.
    - `name` (string): **Optional.** The name of a specific output to read. If omitted, all outputs are returned as a table.
- **Returns:**
    - On success: The parsed JSON value of the output (can be a string, table, etc.).
    - On failure: `nil, error_message`.

## Full Lifecycle Example

```lua
local tf_workdir = "./examples/terraform"

-- Task 1: Init
local init_task = task("terraform-init")
    :description("Initialize Terraform working directory")
    :command(function(this, params)
        local result = terraform.init({workdir = tf_workdir})
        if not result.success then
            return false, "Init failed: " .. result.stderr
        end
        return true, "Terraform initialized successfully"
    end)
    :build()

-- Task 2: Plan
local plan_task = task("terraform-plan")
    :description("Create Terraform execution plan")
    :command(function(this, params)
        local result = terraform.plan({workdir = tf_workdir})
        if not result.success then
            return false, "Plan failed: " .. result.stderr
        end
        return true, "Plan created successfully"
    end)
    :build()

-- Task 3: Apply
local apply_task = task("terraform-apply")
    :description("Apply Terraform plan")
    :command(function(this, params)
        local result = terraform.apply({workdir = tf_workdir, auto_approve = true})
        if not result.success then
            return false, "Apply failed: " .. result.stderr
        end
        return true, "Infrastructure applied successfully"
    end)
    :build()

-- Task 4: Get Output
local output_task = task("terraform-output")
    :description("Read Terraform outputs")
    :command(function(this, params)
        local filename, err = terraform.output({workdir = tf_workdir, name = "report_filename"})
        if not filename then
            return false, "Output failed: " .. err
        end
        log.info("Terraform created file: " .. filename)
        return true, "Output retrieved: " .. filename
    end)
    :build()

-- Task 5: Destroy
local destroy_task = task("terraform-destroy")
    :description("Destroy Terraform-managed infrastructure")
    :command(function(this, params)
        local result = terraform.destroy({workdir = tf_workdir, auto_approve = true})
        if not result.success then
            return false, "Destroy failed: " .. result.stderr
        end
        return true, "Infrastructure destroyed successfully"
    end)
    :build()

-- Workflow: Complete Terraform lifecycle
local terraform_workflow = workflow.define("terraform-lifecycle")
    :description("Complete Terraform infrastructure lifecycle")
    :version("1.0.0")
    :tasks({
        init_task,
        plan_task,
        apply_task,
        output_task,
        destroy_task
    })
```
