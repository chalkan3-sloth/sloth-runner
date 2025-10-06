# Sloth System Documentation

## Overview

The Sloth system allows you to save and reuse `.sloth` workflow files in a SQLite database. This eliminates the need to specify file paths repeatedly and provides centralized management of your workflows.

## Features

- **Save workflows**: Store `.sloth` files in the database with descriptive names
- **Reusability**: Reference saved workflows by name instead of file path
- **Active/Inactive status**: Control which sloths can be used
- **Usage tracking**: Monitor how often workflows are used
- **Version control**: Track file hashes to detect changes

## Commands

### Add a Sloth

Save a new `.sloth` file to the database:

```bash
sloth-runner sloth add <name> --file <path> [options]
```

**Options:**
- `--file, -f`: Path to the `.sloth` file (required)
- `--description, -d`: Description of the sloth
- `--active`: Set as active (default: true)

**Example:**
```bash
sloth-runner sloth add my-workflow --file ./workflows/deploy.sloth --description "Production deployment workflow"
```

### List Sloths

Display all saved sloths:

```bash
sloth-runner sloth list [--active]
```

**Options:**
- `--active, -a`: Show only active sloths

**Example:**
```bash
sloth-runner sloth list --active
```

### Get Sloth Details

View detailed information about a specific sloth:

```bash
sloth-runner sloth get <name>
```

**Example:**
```bash
sloth-runner sloth get my-workflow
```

### Remove a Sloth

Delete a sloth from the database:

```bash
sloth-runner sloth remove <name>
# or
sloth-runner sloth delete <name>
```

**Example:**
```bash
sloth-runner sloth remove my-workflow
```

### Activate/Deactivate a Sloth

Control whether a sloth can be used:

```bash
sloth-runner sloth activate <name>
sloth-runner sloth deactivate <name>
```

**Example:**
```bash
sloth-runner sloth deactivate old-workflow
sloth-runner sloth activate my-workflow
```

## Integration with Run Command

Use a saved sloth in the `run` command:

```bash
sloth-runner run <task> --sloth <name> [options]
```

**Important:** The `--sloth` flag takes **precedence** over `--file`. If both are specified, `--file` is ignored.

**Example:**
```bash
sloth-runner run deploy --sloth my-workflow --yes
```

### How it Works

1. The run command retrieves the sloth content from the database
2. Checks if the sloth is active (fails if inactive)
3. Creates a temporary `.sloth` file with the content
4. Executes the workflow using the temporary file
5. Increments the usage counter
6. Updates the last used timestamp
7. Cleans up the temporary file after execution

## Architecture

The Sloth system follows a clean, modular architecture with industry-standard design patterns:

### Design Patterns

1. **Repository Pattern**: Abstracts data access layer
   - Interface: `internal/sloth/repository.go`
   - Implementation: `internal/sloth/sqlite_repository.go`
   - Mock: `internal/sloth/mock_repository.go`

2. **Service Layer Pattern**: Encapsulates business logic
   - Location: `cmd/sloth-runner/services/sloth_service.go`
   - Handles file operations, validation, and orchestration

3. **Command Pattern**: CLI command structure
   - Location: `cmd/sloth-runner/commands/sloth/`
   - Separate commands for each operation

4. **Factory Pattern**: Dependency injection
   - `NewSlothService()`: Creates service with default repository
   - `NewSlothServiceWithRepository()`: Allows custom repository for testing

### Directory Structure

```
task-runner/
├── internal/sloth/                    # Domain layer
│   ├── sloth.go                      # Domain models
│   ├── repository.go                 # Repository interface
│   ├── sqlite_repository.go          # SQLite implementation
│   ├── mock_repository.go            # Mock for testing
│   └── sqlite_repository_test.go     # Repository tests
│
├── cmd/sloth-runner/services/        # Service layer
│   ├── sloth_service.go              # Business logic
│   └── sloth_service_test.go         # Service tests
│
└── cmd/sloth-runner/commands/sloth/  # Presentation layer
    ├── sloth.go                      # Main command
    ├── add.go                        # Add command
    ├── list.go                       # List command
    ├── get.go                        # Get command
    ├── remove.go                     # Remove/Delete commands
    └── activate.go                   # Activate/Deactivate commands
```

### Database Schema

The sloths are stored in a SQLite database with the following schema:

```sql
CREATE TABLE IF NOT EXISTS sloths (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    file_path TEXT NOT NULL,
    content TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at DATETIME,
    usage_count INTEGER DEFAULT 0,
    tags TEXT,
    file_hash TEXT NOT NULL,
    UNIQUE(name)
);

CREATE INDEX IF NOT EXISTS idx_sloths_name ON sloths(name);
CREATE INDEX IF NOT EXISTS idx_sloths_active ON sloths(is_active);
```

**Indexes:**
- `idx_sloths_name`: Fast lookup by name
- `idx_sloths_active`: Efficient filtering by active status

### Domain Models

#### Sloth
Complete representation of a sloth with all metadata:
```go
type Sloth struct {
    ID          string
    Name        string
    Description string
    FilePath    string
    Content     string
    IsActive    bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
    LastUsedAt  *time.Time
    UsageCount  int
    Tags        string
    FileHash    string
}
```

#### SlothListItem
Lightweight representation for listing:
```go
type SlothListItem struct {
    Name        string
    Description string
    IsActive    bool
    CreatedAt   time.Time
    LastUsedAt  *time.Time
    UsageCount  int
}
```

## Testing

The sloth system includes comprehensive unit tests:

### Repository Tests (55.7% coverage)
Location: `internal/sloth/sqlite_repository_test.go`

Tests include:
- Create sloth
- Create duplicate (error handling)
- Get by name
- List with filters
- Update sloth
- Delete sloth
- Set active/inactive status
- Increment usage counter

### Service Tests (29.7% coverage)
Location: `cmd/sloth-runner/services/sloth_service_test.go`

Tests include:
- Get sloth
- Get active sloth (with inactive check)
- List sloths (all and active only)
- Add sloth
- Add sloth with file not found
- Update sloth
- Remove/Delete sloth
- Activate/Deactivate sloth
- Use sloth (with usage tracking)
- Write content to file
- Close service

### Running Tests

```bash
# Run all sloth tests
go test ./internal/sloth/... ./cmd/sloth-runner/services/... -v

# Run with coverage
go test ./internal/sloth/... ./cmd/sloth-runner/services/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Error Handling

The system defines custom errors for common scenarios:

- `ErrSlothNotFound`: Sloth with given name doesn't exist
- `ErrSlothAlreadyExists`: Sloth with that name already exists
- `ErrSlothInactive`: Attempted to use an inactive sloth

Example:
```go
sloth, err := service.GetActiveSloth(ctx, "my-sloth")
if err == sloth.ErrSlothInactive {
    fmt.Println("This sloth is currently inactive")
}
```

## Best Practices

1. **Use descriptive names**: Choose meaningful names for your sloths
   ```bash
   # Good
   sloth-runner sloth add prod-deploy --file deploy.sloth

   # Avoid
   sloth-runner sloth add temp1 --file deploy.sloth
   ```

2. **Add descriptions**: Help your team understand what each sloth does
   ```bash
   sloth-runner sloth add backup --file backup.sloth --description "Daily database backup to S3"
   ```

3. **Deactivate instead of delete**: Preserve history and usage statistics
   ```bash
   sloth-runner sloth deactivate old-workflow
   ```

4. **Use --active flag**: Quickly see available workflows
   ```bash
   sloth-runner sloth list --active
   ```

5. **Monitor usage**: Check which workflows are being used
   ```bash
   sloth-runner sloth get my-workflow
   # Look at "Usage Count" and "Last Used"
   ```

## Examples

### Complete Workflow Example

```bash
# 1. Create a workflow file
cat > my-deploy.sloth << 'EOF'
workflow({
    name = "deploy",
    description = "Deploy application to production",
    tasks = {
        {
            name = "build",
            run = function()
                print("Building application...")
                return {changed = true, message = "Build completed"}
            end
        },
        {
            name = "deploy",
            depends_on = {"build"},
            run = function()
                print("Deploying to production...")
                return {changed = true, message = "Deployment completed"}
            end
        }
    }
})
EOF

# 2. Add to sloth system
sloth-runner sloth add prod-deploy \
    --file my-deploy.sloth \
    --description "Production deployment workflow"

# 3. Use the sloth
sloth-runner run deploy --sloth prod-deploy --yes

# 4. Check usage statistics
sloth-runner sloth get prod-deploy

# 5. When done, deactivate
sloth-runner sloth deactivate prod-deploy
```

### Managing Multiple Environments

```bash
# Add environment-specific sloths
sloth-runner sloth add dev-deploy --file deploy-dev.sloth
sloth-runner sloth add staging-deploy --file deploy-staging.sloth
sloth-runner sloth add prod-deploy --file deploy-prod.sloth

# Use appropriate sloth for each environment
sloth-runner run deploy --sloth dev-deploy --yes
sloth-runner run deploy --sloth staging-deploy --yes
sloth-runner run deploy --sloth prod-deploy --yes

# List all deployment sloths
sloth-runner sloth list
```

## Database Location

The sloth database is stored in:
- Default: `~/.sloth-runner/sloths.db`
- The database uses SQLite with WAL (Write-Ahead Logging) mode for better concurrency

## Troubleshooting

### Sloth is inactive
**Error:** `sloth is not active`

**Solution:** Activate the sloth:
```bash
sloth-runner sloth activate <name>
```

### Sloth not found
**Error:** `sloth not found`

**Solution:** Check available sloths:
```bash
sloth-runner sloth list
```

### Duplicate sloth name
**Error:** `sloth with this name already exists`

**Solution:** Choose a different name or remove the existing sloth:
```bash
sloth-runner sloth remove <name>
# or use a different name
sloth-runner sloth add my-workflow-v2 --file workflow.sloth
```

### File not found during add
**Error:** `failed to read file`

**Solution:** Verify the file path:
```bash
ls -la /path/to/your/file.sloth
```

## Future Enhancements

Potential improvements for the sloth system:

1. **Tags**: Support for categorizing sloths
2. **Export/Import**: Share sloths between systems
3. **Version history**: Track changes to sloth content
4. **Search**: Find sloths by description or content
5. **Update command**: Update existing sloth from file
6. **Clone command**: Duplicate a sloth with a new name
7. **Validation**: Check sloth syntax before saving
8. **Diff command**: Compare sloth versions
9. **Backup/Restore**: Database backup utilities
10. **Web UI**: Browser-based sloth management

## Contributing

When contributing to the sloth system:

1. Add tests for new features
2. Update documentation
3. Follow existing patterns (Repository, Service, Command)
4. Maintain backwards compatibility
5. Add integration tests for end-to-end scenarios

## License

This feature is part of the sloth-runner project and follows the same license.
