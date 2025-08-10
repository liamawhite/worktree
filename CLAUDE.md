# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a Go CLI tool called `worktree` (binary: `wt`) that provides a semi-opinionated workflow for managing Git worktrees. It enables developers to work on multiple branches simultaneously in separate directories, eliminating the need for constant `git stash` and `git checkout` operations. The tool includes GitHub integration, TUI for interactive selection, and cross-platform support.

## Commands

### Build and Development
```bash
# Build the binary (outputs to bin/wt)
make build
# OR
go build -o bin/wt .

# Run all checks (format, lint, test, dirty check)
make check

# Format code
make format
# OR
gofmt -w .

# Run linter
make lint
# OR  
golangci-lint run

# Run tests
make test
# OR
go test -race -v ./...

# Run integration tests specifically
go test -v ./integration/

# Clean build artifacts
make clean
```

### Application Commands
```bash
# Initial setup for a repository
wt setup github.com/user/repo

# Add a new worktree
wt add feature/branch-name

# Switch between worktrees (interactive or direct)
wt switch
wt switch branch-name

# Remove a worktree
wt remove branch-name

# Clear all worktrees
wt clear

# Configuration management
wt config list
wt config set github.com username
```

## Architecture

### Core Components

- **cmd/**: Cobra-based CLI commands with root command in `root.go`
- **pkg/config/**: YAML-based configuration management for account mappings
- **pkg/git/**: Git operations using go-git library
- **pkg/selector/**: TUI selector using Bubble Tea framework
- **pkg/setup/**: Repository setup and initialization logic
- **pkg/worktree/**: Core worktree management functionality with templates

### Key Dependencies

- **Cobra**: CLI framework for command structure
- **Bubble Tea/Bubbles/Lipgloss**: TUI components for interactive selection
- **go-git**: Git operations without external git binary dependency
- **testify**: Testing framework for assertions

### Configuration

- Default config location: `~/.config/worktree/settings.yaml`
- Override via `--config` flag or `WORKTREE_CONFIG` environment variable
- Stores domain-to-account mappings for Git hosting providers
- Auto-creates with sensible defaults if not present

### Testing Strategy

- Unit tests alongside source files (`*_test.go`)
- Integration tests in `integration/` directory using a framework in `framework.go`
- Race condition detection enabled in test runs
- Testify for assertions and mocking

### Git Integration

- Uses go-git library for reliable Git operations
- Supports both local Git operations and GitHub integration
- Handles enterprise Git hosting solutions
- Template system for post-add hooks in `pkg/worktree/templates/`

The codebase follows standard Go conventions with proper error handling, configuration management, and separation of concerns between CLI, business logic, and Git operations.