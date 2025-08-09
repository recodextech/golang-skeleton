#!/bin/bash

# Script to copy golang-skeleton project to a new directory with updated package names
# Usage: ./copy_project.sh <new_project_name>

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if new project name is provided
if [ $# -eq 0 ]; then
    print_error "Please provide a new project name"
    echo "Usage: $0 <new_project_name>"
    echo "Example: $0 my-awesome-service"
    exit 1
fi

NEW_PROJECT_NAME="$1"
CURRENT_DIR=$(pwd)
SOURCE_DIR=$(basename "$CURRENT_DIR")
PARENT_DIR=$(dirname "$CURRENT_DIR")
TARGET_DIR="$PARENT_DIR/$NEW_PROJECT_NAME"

# Validate project name
if [[ ! "$NEW_PROJECT_NAME" =~ ^[a-zA-Z0-9_-]+$ ]]; then
    print_error "Project name should only contain letters, numbers, hyphens, and underscores"
    exit 1
fi

# Check if target directory already exists
if [ -d "$TARGET_DIR" ]; then
    print_error "Directory '$TARGET_DIR' already exists"
    exit 1
fi

# Check if we're in a Go project directory (has go.mod)
if [ ! -f "go.mod" ]; then
    print_error "No go.mod file found in current directory. Please run this script from the golang-skeleton project root."
    exit 1
fi

# Check if we're in the golang-skeleton directory
if [ "$SOURCE_DIR" != "golang-skeleton" ]; then
    print_warning "Current directory is not 'golang-skeleton', but continuing anyway..."
fi

print_status "Starting project copy from '$SOURCE_DIR' to '$NEW_PROJECT_NAME'..."

# Step 1: Copy the entire directory excluding script files
print_status "Copying project files..."
cp -r "$CURRENT_DIR" "$TARGET_DIR"

# Remove the copy script and usage documentation from the new project
rm -f "$TARGET_DIR/copy_project.sh"
rm -f "$TARGET_DIR/COPY_PROJECT_USAGE.md"
print_success "Files copied to $TARGET_DIR (excluding copy script and documentation)"

# Step 2: Update go.mod
print_status "Updating go.mod module name..."
cd "$TARGET_DIR"
sed -i.bak "s/^module golang-skeleton$/module $NEW_PROJECT_NAME/" go.mod
rm go.mod.bak
print_success "Updated go.mod module name to '$NEW_PROJECT_NAME'"

# Step 3: Update all import statements
print_status "Updating import statements in Go files..."

# Find all .go files and update imports
find . -name "*.go" -type f -exec sed -i.bak "s|golang-skeleton|$NEW_PROJECT_NAME|g" {} \;

# Clean up backup files
find . -name "*.bak" -type f -delete

print_success "Updated all import statements"

# Step 4: Update any references in config files
print_status "Checking for references in config files..."

# Update configs.env if it contains project references
if [ -f "configs.env" ]; then
    if grep -q "golang-skeleton" configs.env; then
        sed -i.bak "s|golang-skeleton|$NEW_PROJECT_NAME|g" configs.env
        rm configs.env.bak 2>/dev/null || true
        print_success "Updated configs.env"
    fi
fi

# Update README if it exists
if [ -f "readme.md" ]; then
    if grep -q "golang-skeleton" readme.md; then
        sed -i.bak "s|golang-skeleton|$NEW_PROJECT_NAME|g" readme.md
        rm readme.md.bak 2>/dev/null || true
        print_success "Updated readme.md"
    fi
fi

# Step 5: Clean up git history (optional)
print_status "Cleaning up git repository..."
if [ -d ".git" ]; then
    rm -rf .git
    git init
    git add .
    git commit -m "Initial commit for $NEW_PROJECT_NAME (copied from golang-skeleton)"
    print_success "Initialized new git repository"
else
    print_warning "No git repository found, skipping git initialization"
fi

# Step 6: Clean up go modules
print_status "Cleaning up Go modules..."
if [ -f "go.sum" ]; then
    rm go.sum
    print_status "Removed go.sum (will be regenerated on next go mod tidy)"
fi

go mod tidy
print_success "Go modules cleaned up"

# Step 7: Verify the changes
print_status "Verifying changes..."
MODULE_NAME=$(head -1 go.mod | cut -d' ' -f2)
if [ "$MODULE_NAME" = "$NEW_PROJECT_NAME" ]; then
    print_success "Module name verification passed"
else
    print_error "Module name verification failed. Expected '$NEW_PROJECT_NAME', got '$MODULE_NAME'"
fi

# Count updated files
UPDATED_FILES=$(find . -name "*.go" -exec grep -l "$NEW_PROJECT_NAME" {} \; | wc -l | tr -d ' ')
TOTAL_FILES=$(find . -type f \( -name "*.go" -o -name "*.mod" -o -name "*.env" -o -name "*.md" \) | wc -l | tr -d ' ')
print_success "Updated $UPDATED_FILES Go files with new module name"
print_success "Total project files copied: $TOTAL_FILES (excluding copy script and documentation)"

# Final summary
echo
echo "=================================="
print_success "PROJECT COPY COMPLETED SUCCESSFULLY!"
echo "=================================="
echo
echo "Summary:"
echo "  • Source: $CURRENT_DIR"
echo "  • Target: $TARGET_DIR"
echo "  • New module name: $NEW_PROJECT_NAME"
echo "  • Updated files: $UPDATED_FILES Go files"
echo
echo "Next steps:"
echo "  1. cd $TARGET_DIR"
echo "  2. go mod download  # Download dependencies"
echo "  3. go build         # Test build"
echo "  4. go run main.go   # Test run"
echo
print_status "Happy coding with your new project: $NEW_PROJECT_NAME!"
