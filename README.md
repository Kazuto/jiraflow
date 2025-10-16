# JiraFlow

A CLI tool that creates GitFlow-compliant git branches from Jira ticket information. Simplify your workflow by automatically generating properly formatted branch names from ticket numbers and titles.

## Features

- 🎯 Automatically format Jira ticket information into git branch names
- 🔗 Fetch ticket titles using Jira CLI (optional)
- 🌊 GitFlow support (feature, hotfix, refactor, support branches)
- 🔧 Configurable via YAML with automatic default config creation
- ⚡ Fast and lightweight Go implementation
- 🧹 Smart sanitization of branch names (removes special characters, truncates length)
- 🎨 Beautiful interactive TUI using Charm Bracelet libraries
- 🔀 Interactive branch selection - choose any local branch as base
- 🔍 Searchable branch list with fuzzy matching

## Prerequisites

- **Git** - for branch creation
- **[Jira CLI](https://github.com/ankitpokhrel/jira-cli)** (optional) - required only if you want to automatically fetch ticket titles from Jira

### Installing Jira CLI (Optional)

JiraFlow can automatically fetch ticket titles from Jira using the [Jira CLI](https://github.com/ankitpokhrel/jira-cli). This is optional - you can always enter titles manually.

#### Install Jira CLI
```bash
# macOS
brew install ankitpokhrel/jira-cli/jira-cli

# Linux/Windows - see https://github.com/ankitpokhrel/jira-cli#installation
```

#### Configure Jira CLI
```bash
jira init
```

This will prompt you for:
- Jira server URL (e.g., `https://yourcompany.atlassian.net`)
- Email address
- API token (create at: Account Settings → Security → API tokens)

Once configured, JiraFlow will automatically fetch ticket titles when you provide just the ticket number.

## Installation

### Prerequisites

- **Go 1.21 or later** - Required for building from source
- **Git** - Required for branch operations
- **[Jira CLI](https://github.com/ankitpokhrel/jira-cli)** (optional) - For automatic ticket title fetching

### Quick Install (Recommended)

#### Using Go Install
```bash
go install github.com/kazuto/jiraflow@latest
```

#### Using Homebrew (macOS/Linux)
```bash
# Coming soon
brew install jiraflow
```

### Build from Source

#### 1. Clone the Repository
```bash
git clone https://github.com/kazuto/jiraflow.git
cd jiraflow
```

#### 2. Install Dependencies
```bash
go mod download
```

#### 3. Build the Binary
```bash
# Using Make (recommended)
make build

# Or build manually
go build -o jiraflow

# Build with version info
make build  # Automatically includes version info

# Build for all platforms
make build-all
```

#### 4. Install System-wide (Optional)
```bash
# macOS/Linux
sudo mv jiraflow /usr/local/bin/

# Or add to your PATH
mkdir -p ~/bin
mv jiraflow ~/bin/
export PATH="$HOME/bin:$PATH"  # Add to your shell profile
```

### Development Setup

#### 1. Clone and Setup
```bash
git clone https://github.com/kazuto/jiraflow.git
cd jiraflow
go mod download

# Set up development environment (installs linters, etc.)
make dev-setup
```

#### 2. Run Tests
```bash
# Using Make (recommended)
make test

# Run with coverage report
make test-coverage

# Or run manually
go test ./...
go test -cover ./...
go test -v ./...

# Run specific package tests
go test ./internal/config
```

#### 3. Run Development Build
```bash
# Run without building
go run main.go

# Run with flags
go run main.go --type feature --ticket PROJ-123
```

#### 4. Build Development Binary
```bash
# Using Make
make build

# Quick build for testing
go build -o jiraflow-dev

# Build with race detection (for testing)
go build -race -o jiraflow-dev
```

#### 5. Available Make Commands
```bash
# See all available commands
make help

# Common development commands
make dev-setup    # Set up development environment
make test         # Run tests
make lint         # Run code linter
make fmt          # Format code
make clean        # Clean build artifacts
make release      # Create release archives
```

#### 6. Continuous Integration

The project includes GitHub Actions workflows for:
- **Testing** - Runs tests on Go 1.21, 1.22, and 1.23
- **Linting** - Code quality checks with golangci-lint
- **Security** - Security scanning with gosec
- **Building** - Multi-platform binary builds
- **Releases** - Automatic release creation with binaries

### Binary Releases

Download pre-built binaries for your platform from the [releases page](https://github.com/kazuto/jiraflow/releases).

Available platforms:
- **Linux** (amd64, arm64)
- **macOS** (amd64, arm64/Apple Silicon)
- **Windows** (amd64)

#### Installation from Binary
```bash
# Download and install (example for Linux amd64)
curl -L https://github.com/kazuto/jiraflow/releases/latest/download/jiraflow-linux-amd64 -o jiraflow
chmod +x jiraflow
sudo mv jiraflow /usr/local/bin/
```

### Verify Installation

```bash
# Check version
jiraflow --version

# Run help
jiraflow --help

# Test in a Git repository
cd /path/to/your/git/repo
jiraflow --dry-run --type feature --ticket TEST-123 --title "Test installation"
```

## Configuration

JiraFlow automatically creates a configuration file at `~/.config/jiraflow/jiraflow.yaml` on first run with sensible defaults. No manual setup required!

### Automatic Configuration

When you run JiraFlow for the first time, it will:
1. Create the `~/.config/jiraflow/` directory
2. Generate a default configuration file
3. Use the configuration immediately

### Manual Configuration (Optional)

You can customize the configuration by editing the file:

```bash
# Edit configuration
nano ~/.config/jiraflow/jiraflow.yaml

# Or use your preferred editor
code ~/.config/jiraflow/jiraflow.yaml
```

### Configuration Options

See [`jiraflow.example.yaml`](jiraflow.example.yaml) for a complete configuration example with all available options and documentation.

#### Key Settings

```yaml
# Maximum branch name length (10-200, default: 60)
max_branch_length: 60

# Default branch type for non-interactive mode (default: feature)
default_branch_type: feature

# Branch type prefixes - customize to match your team's conventions
branch_types:
  feature: "feature/"     # New features and enhancements
  hotfix: "hotfix/"       # Critical bug fixes for production
  refactor: "refactor/"   # Code improvements without changing functionality
  support: "support/"     # Maintenance and support tasks

# Branch name sanitization
sanitization:
  separator: "-"          # Replace spaces/special chars (default: -)
  lowercase: true         # Convert to lowercase (default: true)
  remove_umlauts: false   # Remove German umlauts äöüÄÖÜß (default: false)
```

#### Custom Configuration Examples

```bash
# Copy example configuration
cp jiraflow.example.yaml ~/.config/jiraflow/jiraflow.yaml

# Edit with your preferences
nano ~/.config/jiraflow/jiraflow.yaml
```

## Usage

### Basic Usage

```bash
jiraflow [branch-type] [ticket-number] [ticket-title (optional)]
```

If no ticket title is provided, JiraFlow will automatically fetch the title from Jira using the Jira CLI (`jira issue view --raw` command).

### Examples

**Create a feature branch with explicit title:**
```bash
jiraflow feature PROJ-123 "Add user profile dashboard"
```
Output: `feature/PROJ-123-Add-user-profile-dashboard`

**Create a feature branch (fetch title using Jira CLI):**
```bash
jiraflow feature PROJ-123
```
Output: `feature/PROJ-123-Add-user-profile-dashboard`
(Title fetched automatically via `jira issue view PROJ-123 --raw`)

**Create a hotfix branch:**
```bash
jiraflow hotfix PROJ-456
```
Output: `hotfix/PROJ-456-Fix-critical-payment-gateway-timeout`

**Create a refactor branch with explicit title:**
```bash
jiraflow refactor PROJ-789 "Restructure pricing calculation logic"
```
Output: `refactor/PROJ-789-Restructure-pricing-calculation-logic`

**Create a support branch:**
```bash
jiraflow support PROJ-321
```
Output: `support/PROJ-321-Add-legacy-API-compatibility-layer`

**Using default branch type (feature):**
```bash
jiraflow PROJ-555
```
Output: `feature/PROJ-555-Implement-new-authentication-flow`

### Interactive Mode (Recommended)

Run jiraflow without arguments to launch the beautiful interactive TUI:

```bash
jiraflow
```

#### Interactive Workflow

The TUI guides you through a 4-step process:

1. **🎯 Branch Type Selection**
   - Choose from: feature, hotfix, refactor, support
   - Navigate with arrow keys, select with Enter

2. **🌿 Base Branch Selection**
   - See all your local Git branches
   - Real-time search with `/` key
   - Fuzzy matching to find branches quickly

3. **🎫 Ticket Information**
   - Enter Jira ticket number (e.g., PROJ-123)
   - Optionally enter title or auto-fetch from Jira
   - Tab between fields, Enter to continue

4. **✅ Confirmation & Creation**
   - Preview the final branch name
   - Confirm to create and checkout the branch
   - See success confirmation

#### TUI Features

- 🎨 **Beautiful Interface** - Built with Charm Bracelet libraries
- 🔍 **Smart Search** - Fuzzy matching for branch names
- ⌨️ **Keyboard Navigation** - Intuitive shortcuts (↑/↓, Enter, Esc, /)
- 🔙 **Step Navigation** - Go back/forward through steps
- 🎯 **Real-time Preview** - See branch name as you type
- 🚨 **Error Handling** - Clear error messages and recovery options
- 🎭 **Graceful Degradation** - Works even if Jira CLI is unavailable

### TUI Controls

When using the interactive mode, use these keyboard shortcuts:

- **↑/↓** - Navigate through options
- **Enter** - Select current option or proceed to next step
- **Esc** - Go back to previous step
- **/** - Search/filter (in branch selection)
- **q** or **Ctrl+C** - Quit the application

### Checkout After Creation

The tool automatically creates and checks out the new branch. To only print the branch name without creating it:

```bash
jiraflow --dry-run feature PROJ-123
```

Or with an explicit title:

```bash
jiraflow --dry-run feature PROJ-123 "Add user profile dashboard"
```

## GitFlow Branch Types

- **feature/** - New features and enhancements
- **hotfix/** - Critical bug fixes for production
- **refactor/** - Code improvements without changing functionality
- **support/** - Maintenance and support tasks

## Troubleshooting

### Common Issues

#### "Not a Git repository"
```bash
# Ensure you're in a Git repository
git status

# Initialize if needed
git init
```

#### "Jira CLI not found" or "Failed to fetch title"
This is normal if Jira CLI isn't installed or configured. JiraFlow will:
- Allow manual title entry
- Continue working without Jira integration
- Show helpful error messages

#### "Branch already exists"
```bash
# List existing branches
git branch

# Delete existing branch if needed
git branch -D feature/PROJ-123-existing-branch

# Or use a different ticket number/title
```

#### Configuration Issues
```bash
# Reset to defaults
rm ~/.config/jiraflow/jiraflow.yaml
jiraflow  # Will recreate with defaults

# Check configuration
cat ~/.config/jiraflow/jiraflow.yaml
```

#### Permission Issues
```bash
# Ensure proper permissions for config directory
chmod 755 ~/.config/jiraflow/
chmod 644 ~/.config/jiraflow/jiraflow.yaml
```

### Getting Help

- **In-app help**: Press `?` in interactive mode
- **Command help**: `jiraflow --help`
- **Version info**: `jiraflow --version`
- **Issues**: [GitHub Issues](https://github.com/kazuto/jiraflow/issues)

### Debug Mode

```bash
# Run with verbose output
JIRAFLOW_DEBUG=1 jiraflow

# Test configuration loading
jiraflow --dry-run --type feature --ticket TEST-123
```

## Contributing

Contributions are welcome! Here's how you can help:

### Reporting Issues

Found a bug or have a feature request? Please [open an issue](https://github.com/kazuto/jiraflow/issues) with:
- Clear description of the problem/feature
- Steps to reproduce (for bugs)
- Expected vs actual behavior
- Your environment (OS, Go version)

### Development Setup

1. **Fork the repository**
   ```bash
   git clone https://github.com/kazuto/jiraflow.git
   cd jiraflow
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make your changes and test**
   ```bash
   go test ./...
   go build
   ```

5. **Commit your changes**
   ```bash
   git commit -m "feat: add your feature description"
   ```

6. **Push and create a Pull Request**
   ```bash
   git push origin feature/your-feature-name
   ```

### Coding Standards

- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Write unit tests for new features
- Update documentation for API changes
- Use conventional commit messages (feat, fix, docs, refactor, test, chore)

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## License

MIT License - see [LICENSE](LICENSE) file for details

## Acknowledgments

- Inspired by GitFlow workflow
- Built for developers who love automation

---

**Made with ❤️ for streamlined development workflows**