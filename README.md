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

### Installing Jira CLI

Follow the installation instructions at [github.com/ankitpokhrel/jira-cli](https://github.com/ankitpokhrel/jira-cli#installation).

After installation, configure Jira CLI:

```bash
jira init
```

This will prompt you for your Jira server URL, email, and API token. Once configured, JiraFlow will use Jira CLI to fetch ticket information.

## Installation

### Using Go Install

```bash
go install github.com/kazuto/jiraflow@latest
```

### From Source

```bash
git clone https://github.com/kazuto/jiraflow.git
cd jiraflow
go build -o jiraflow
sudo mv jiraflow /usr/local/bin/
```

### Binary Releases

Download the latest binary for your platform from the [releases page](https://github.com/kazuto/jiraflow/releases).

## Configuration

JiraFlow uses a configuration file located at `~/.config/jiraflow/jiraflow.yaml`.

Create the configuration directory and file:

```bash
mkdir -p ~/.config/jiraflow
touch ~/.config/jiraflow/jiraflow.yaml
```

### Configuration Options

```yaml
# Maximum branch name length (default: 60)
max_branch_length: 60

# Default branch type if not specified (default: feature)
default_branch_type: feature

# Branch type prefixes
branch_types:
  feature: "feature/"
  hotfix: "hotfix/"
  refactor: "refactor/"
  support: "support/"

# Character replacements for branch name sanitization
sanitization:
  # Replace spaces and special characters with this (default: -)
  separator: "-"
  # Convert to lowercase (default: true)
  lowercase: false
  # Remove German umlauts (äöüÄÖÜß) (default: false)
  remove_umlauts: false
```

## Usage

### Basic Usage

```bash
jiraflow [branch-type] [ticket-number] [ticket-title (optional)]
```

If no ticket title is provided, JiraFlow will automatically fetch the title from Jira using the Jira CLI (`jira view` command).

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
(Title fetched automatically via `jira view PROJ-123`)

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

### Interactive Mode

Run jiraflow without arguments for beautiful interactive prompts:

```bash
jiraflow
```

The interactive TUI will guide you through:
1. **Branch type selection** - Choose from feature, hotfix, refactor, or support
2. **Base branch selection** - Pick any local branch with search functionality
3. **Ticket number input** - Enter your Jira ticket number
4. **Title input** - Optionally provide a title or let it fetch from Jira CLI
5. **Confirmation** - Review and confirm before creating the branch

The TUI features:
- 🎨 Beautiful, colorful interface
- 🔍 Searchable branch list with fuzzy matching
- ⌨️ Intuitive keyboard navigation
- 🔙 Step-by-step navigation with back/forward support

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