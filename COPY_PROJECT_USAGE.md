# Project Copy Script Usage

This script (`copy_project.sh`) allows you to easily copy the golang-skeleton project to a new directory with a different name, automatically updating all package references and module names.

## Usage

```bash
./copy_project.sh <new_project_name>
```

## Example

```bash
# Copy the project to a new directory called "my-awesome-service"
./copy_project.sh my-awesome-service
```

## What the script does

1. **Validates input**: Ensures the new project name contains only valid characters
2. **Checks prerequisites**: Verifies you're in a Go project directory with `go.mod`
3. **Copies files**: Creates a complete copy of the project in the parent directory
4. **Updates module name**: Changes `module golang-skeleton` to `module <new_project_name>` in `go.mod`
5. **Updates imports**: Replaces all `golang-skeleton` import references with the new module name in all `.go` files
6. **Updates config files**: Updates any references in `configs.env` and `readme.md`
7. **Initializes Git**: Creates a new Git repository with an initial commit
8. **Cleans dependencies**: Runs `go mod tidy` to clean up the dependency list

## Requirements

- Must be run from within the `golang-skeleton` project directory
- New project name should only contain letters, numbers, hyphens, and underscores
- Target directory must not already exist

## Output

The script provides colored output showing progress:
- **Blue [INFO]**: General information
- **Green [SUCCESS]**: Successful operations
- **Yellow [WARNING]**: Warnings (non-fatal)
- **Red [ERROR]**: Errors (will stop execution)

## Example Output

```
[INFO] Starting project copy from 'golang-skeleton' to 'my-service'...
[INFO] Copying project files...
[SUCCESS] Files copied to /path/to/my-service
[INFO] Updating go.mod module name...
[SUCCESS] Updated go.mod module name to 'my-service'
[INFO] Updating import statements in Go files...
[SUCCESS] Updated all import statements
[SUCCESS] Updated 50 Go files with new module name

==================================
[SUCCESS] PROJECT COPY COMPLETED SUCCESSFULLY!
==================================

Summary:
  • Source: /path/to/golang-skeleton
  • Target: /path/to/my-service
  • New module name: my-service
  • Updated files: 50 Go files

Next steps:
  1. cd /path/to/my-service
  2. go mod download  # Download dependencies
  3. go build         # Test build
  4. go run main.go   # Test run
```

## Troubleshooting

### "No go.mod file found"
- Make sure you're running the script from the golang-skeleton project root directory

### "Directory already exists"
- Choose a different project name or remove the existing directory

### Build errors after copying
- Run `go mod download` to ensure all dependencies are downloaded
- Some dependencies may have architecture-specific issues (e.g., Kafka libraries on ARM64 Macs)

## Script Features

- ✅ Input validation
- ✅ Comprehensive error checking
- ✅ Colored output for better UX
- ✅ Automatic Git repository initialization
- ✅ Module dependency cleanup
- ✅ Configuration file updates
- ✅ Detailed progress reporting
- ✅ Final verification

