# 🔄 ProfileSync - Cross-Platform User Profile Migrator

> **The ultimate tool for migrating user profiles between Linux, macOS, and Windows**

---

## 🎯 Problem Solved

Migrating user profiles between different operating systems is a **nightmare**:
- Settings are stored in completely different formats
- Configuration files have platform-specific paths
- Migration requires manual copying and editing
- Risk of missing critical settings
- No audit trail of what was migrated

**ProfileSync solves this by automating the entire process.**

---

## ✨ Features

### 🚀 Key Capabilities
- **Cross-Platform Migration** - Seamlessly migrate between Linux, macOS, and Windows
- **Comprehensive Coverage** - Migrates settings for 20+ developer tools
- **Safe Dry-Run Mode** - Preview changes before applying them
- **Force Mode** - Overwrite existing files when needed
- **Type Categorization** - Organizes settings by type (IDE, Shell, Security, etc.)
- **Audit Trail** - Detailed migration report with success/failure tracking

### 📦 Supported Tools

| Category | Tools |
|----------|-------|
| **IDEs** | VS Code, IntelliJ IDEA |
| **Editors** | Vim, Emacs |
| **Shells** | Bash, Zsh, Fish |
| **Terminal** | Tmux |
| **Version Control** | Git |
| **Security** | SSH keys & config |
| **Browsers** | Chrome, Firefox |
| **Package Managers** | NPM, Yarn, Pip |
| **Containers** | Docker |
| **Kubernetes** | kubectl, Helm |
| **Infrastructure** | Terraform |
| **Cloud** | AWS CLI |

---

## 🛠️ Installation

### Build from Source

```bash
cd profilesync
go mod download
go build -o profilesync cmd/profilesync/main.go
```

### Install Globally

```bash
go install -o /usr/local/bin/profilesync ./cmd/profilesync
```

---

## 🚀 Usage

### Basic Usage

```bash
# Preview migration (dry-run)
./profilesync --source=linux --dest=macos --dry-run

# Perform actual migration
./profilesync --source=linux --dest=macos --dry-run=false

# Force overwrite existing files
./profilesync --source=linux --dest=macos --force=true
```

### Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `--source` | Source platform (linux, macos, windows) | Current OS |
| `--dest` | Destination platform (linux, macos, windows) | Current OS |
| `--dry-run` | Preview without making changes | true |
| `--force` | Overwrite existing files | false |
| `--verbose` | Show detailed output | false |
| `--help` | Show help message | false |

### Examples

#### Migrate from Linux to macOS

```bash
# Preview the migration
./profilesync --source=linux --dest=macos --dry-run

# Perform the migration
./profilesync --source=linux --dest=macos
```

#### Migrate from Windows to Linux

```bash
# Full migration with force overwrite
./profilesync --source=windows --dest=linux --force=true
```

#### Migrate between same platforms

```bash
# Backup profile on same OS
./profilesync --source=linux --dest=linux --dry-run
```

---

## 📊 Migration Report

The tool generates a detailed report showing:

```
============================================================
📊 MIGRATION REPORT
============================================================
Source Platform:   linux
Destination:       macos
Mode:              DRY RUN

✅ Successfully migrated: 15
⏭️  Skipped:           3
❌ Failed:            0

============================================================

📁 Items by Type:
  • Browser: 2 items
  • Cloud: 2 items
  • Container: 1 items
  • Editor: 2 items
  • IDE: 2 items
  • Infrastructure: 2 items
  • Package Manager: 3 items
  • Security: 3 items
  • Shell: 2 items
  • Version Control: 2 items
  • Kubernetes: 1 items

============================================================
⚠️  This was a DRY RUN. No files were actually migrated.
Run without --dry-run to perform the actual migration.
```

---

## 🔒 Security Features

- **SSH Key Preservation** - Maintains proper permissions on private keys
- **Credential Mapping** - Safely handles credentials and secrets
- **Audit Trail** - Tracks all migrated items
- **No Data Modification** - Preserves original file contents

---

## 🧪 Testing

### Run Tests

```bash
go test ./...
```

### Test Coverage

```bash
go test -cover
```

---

## 📝 Configuration

ProfileSync uses a built-in mapping system that automatically detects common configuration file locations. You can extend this by modifying the `GetDefaultMappings()` function in `main.go`.

### Custom Mappings Example

```go
func GetDefaultMappings() map[string]string {
    mappings := map[string]string{
        "custom/tool/.config": "custom/tool/.config",
        // Add your custom mappings here
    }
    return mappings
}
```

---

## 🐳 Docker Support

```bash
docker run -it \
  -v $HOME/.config:/host_config \
  -v $HOME/.ssh:/host_ssh \
  profilesync:latest \
  --source=linux --dest=macos --dry-run
```

---

## 🚧 Roadmap

- [ ] GUI interface for easier use
- [ ] Custom mapping file support (JSON/YAML)
- [ ] Conflict resolution wizard
- [ ] Rollback capability
- [ ] Enterprise deployment integration
- [ ] GitOps support for profile management

---

## 📚 Use Cases

### 1. **OS Migration**
Migrate your entire development environment when switching from Linux to macOS or Windows.

### 2. **New Machine Setup**
Clone your exact development setup to a new computer.

### 3. **Team Standardization**
Distribute standard configurations across team members' machines.

### 4. **Disaster Recovery**
Quickly restore your development environment from backup.

### 5. **Cross-Platform Testing**
Ensure consistent configurations across different platforms.

---

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

1. Fork the repository
2. Create a feature branch
3. Add your custom mappings
4. Submit a pull request

---

## 📄 License

MIT License - Free for commercial and personal use

---

## 🙏 Acknowledgments

Built with GPU for developers who switch between platforms.

---

**Version:** 1.0.0  
**Author:** @hallucinaut  
**Last Updated:** February 25, 2026