# worktree

A semi-opinionated git worktree workflow CLI for developers that enables you to work on multiple branches simultaneously in separate directories, eliminating the need to constantly `git stash`, `git checkout`, and context switch between different features or bug fixes. With supoort for GitHub, GitLab and any enterprise git solution.

## Installation

Download the latest binary from the [releases page](https://github.com/liamawhite/worktree/releases/latest).

## Features

- **Interactive Selection**: Uses a TUI for selecting and switching between worktrees
- **GitHub Integration**: Support for GitHub forks and enterprise Git hosting
- **Git Integration**: Built on top of go-git for reliable Git operations
- **Configurable**: Flexible configuration system with environment variable support
- **Cross-platform**: Works on macOS, Linux, and Windows

## Example Workflow

Here's a typical workflow using `wt` for parallel development:

### 1. Initial Setup
```bash
wt setup github.com/liamawhite/worktree
```

### 2. Create First Worktree for Feature Work
```bash
# Add a worktree for feature development
wt add feature/user-auth

# This creates a new worktree and switches to it
# Work on your feature...
echo "// Auth implementation" >> auth.go
```

### 3. Create Second Worktree for Parallel Work  
```bash
# Add another worktree for a different feature
wt add feature/api-endpoints

# Work on the API endpoints...
echo "// API endpoints" >> api.go  
```

### 4. Switch Between Worktrees
```bash
# Targeted feature switching
wt switch feature/user-auth

# OR interactive switching back to first worktree
wt switch
# Select "feature/user-auth" from the TUI

# Push your auth feature
git commit -am "Super secure auth, no hard-coded signing keys here!"
git push origin feature/user-auth

# Switch back to API work
wt switch feature/api-endpoints

# Push your API feature
git commit -am "Add a very necessary gRPC server to our little cli tool"
git push origin feature/api-endpoints
```

This workflow lets you maintain multiple branches simultaneously without the overhead of constant `git stash`/`git checkout` cycles.

## Configuration

The `wt` tool uses a flexible configuration system to manage Git hosting providers and authentication methods. Configuration is stored in `~/.config/worktree/settings.yaml`.

### Quick Start

Configure your GitHub username:
```bash
wt config set-account github.com your-username
```

Set SSH as your preferred clone method:
```bash
wt config set-clone-method github.com ssh
```

View all configured hosts:
```bash
wt config list
```

### Configuration Commands

#### `wt config list`
Shows all configured hosts with their accounts and clone methods:
```bash
$ wt config list
Configured hosts:
  enterprise.company.com: john.doe (clone: ssh)
  github.com: your-username (clone: http)
  gitlab.com: your-username (clone: ssh)
```

#### `wt config set-account <domain> <account>`
Sets the account/username for a Git hosting provider:
```bash
# GitHub.com
wt config set-account github.com your-username

# GitHub Enterprise
wt config set-account github.enterprise.com john.doe

# GitLab
wt config set-account gitlab.com your-username
```

#### `wt config set-clone-method <domain> <method>`
Sets the clone method (http or ssh) for a Git hosting provider:
```bash
# Use SSH for GitHub (recommended for authenticated access)
wt config set-clone-method github.com ssh

# Use HTTPS for enterprise Git (might be required by corporate policy)
wt config set-clone-method github.enterprise.com http
```

### Configuration File Format

The configuration is stored as YAML in `~/.config/worktree/settings.yaml`:

```yaml
hosts:
  github.com:
    account: your-username
    clone_method: ssh
  github.enterprise.com:
    account: john.doe
    clone_method: http
  gitlab.com:
    account: your-username
    clone_method: ssh
```

## Development

### Building

```bash
# Build the binary
go build -o bin/wt .
```

### Testing

```bash
# Run tests
go test -race -v ./...
```

### Linting

```bash
# Format code
gofmt -w .

# Run linter
golangci-lint run
```

## License

See [LICENSE](LICENSE) file for details.
