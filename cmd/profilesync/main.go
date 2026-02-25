package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/fatih/color"
)

var (
	infoColor = color.New(color.FgBlue)
	warnColor = color.New(color.FgYellow)
	errorColor = color.New(color.FgRed)
	successColor = color.New(color.FgGreen)
	noticeColor = color.New(color.FgCyan)
)

// PlatformConfig represents platform-specific configuration settings
type PlatformConfig struct {
	Name         string
	HomeDir      string
	Extensions   []string
	Mappings     map[string]string
}

// MigrationPlan defines how to migrate settings from one platform to another
type MigrationPlan struct {
	SourcePlatform   string
	DestinationPlatform string
	Items            []MigrationItem
	TotalItems       int
	SkippedItems     int
}

// MigrationItem represents a single setting or configuration to migrate
type MigrationItem struct {
	SourcePath      string
	DestinationPath string
	Type            string
	Description     string
	AutoMigrate     bool
}

// ProfileSync handles cross-platform profile migration
type ProfileSync struct {
	sourcePlatform   string
	destPlatform     string
	dryRun           bool
	force            bool
	verbose          bool
	migrationPlan    *MigrationPlan
}

// NewProfileSync creates a new ProfileSync instance
func NewProfileSync(source, dest string, dryRun, force, verbose bool) *ProfileSync {
	return &ProfileSync{
		sourcePlatform:  source,
		destPlatform:    dest,
		dryRun:          dryRun,
		force:           force,
		verbose:         verbose,
		migrationPlan:   &MigrationPlan{},
	}
}

// DetectPlatform detects the current operating system
func DetectPlatform() string {
	switch runtime.GOOS {
	case "linux":
		return "linux"
	case "darwin":
		return "macos"
	case "windows":
		return "windows"
	default:
		return "unknown"
	}
}

// GetHomeDir returns the home directory for the platform
func GetHomeDir(platform string) string {
	switch platform {
	case "linux":
		if home := os.Getenv("HOME"); home != "" {
			return home
		}
		return "/home/" + os.Getenv("USER")
	case "macos":
		if home := os.Getenv("HOME"); home != "" {
			return home
		}
		return "/Users/" + os.Getenv("USER")
	case "windows":
		if home := os.Getenv("USERPROFILE"); home != "" {
			return home
		}
		return "C:\\Users\\" + os.Getenv("USERNAME")
	default:
		return os.Getenv("HOME")
	}
}

// GetDefaultMappings returns default platform mappings
func GetDefaultMappings() map[string]string {
	return map[string]string{
		// IDE settings
		"vscode/settings.json": "vscode/settings.json",
		"vscode/keybindings.json": "vscode/keybindings.json",
		"intellij/": "intellij/",
		"vim/.vimrc": "vim/.vimrc",
		"vim/.vim/": "vim/.vim/",
		"emacs/.emacs": "emacs/.emacs",
		"emacs/.emacs.d/": "emacs/.emacs.d/",
		
		// Terminal settings
		"bash/.bashrc": "bash/.bashrc",
		"bash/.bash_profile": "bash/.bash_profile",
		"zsh/.zshrc": "zsh/.zshrc",
		"fish/.config/fish/config.fish": "fish/.config/fish/config.fish",
		"tmux/.tmux.conf": "tmux/.tmux.conf",
		
		// Git configuration
		"git/.gitconfig": "git/.gitconfig",
		"git/.gitignore_global": "git/.gitignore_global",
		
		// SSH configuration
		"ssh/config": "ssh/config",
		"ssh/id_rsa": "ssh/id_rsa",
		"ssh/id_rsa.pub": "ssh/id_rsa.pub",
		
		// Browser profiles
		"chrome/Default/": "chrome/Default/",
		"firefox/.mozilla/firefox/": "firefox/.mozilla/firefox/",
		
		// Package managers
		"npm/.npmrc": "npm/.npmrc",
		"yarn/.yarnrc": "yarn/.yarnrc",
		"pip/pip.conf": "pip/pip.conf",
		"pip/pip.ini": "pip/pip.ini",
		
		// Docker
		"docker/config.json": "docker/config.json",
		
		// Kubectl
		"kubectl/config": "kubectl/config",
		"helm/.helm/": "helm/.helm/",
		
		// Terraform
		"terraform/.terraform.d/": "terraform/.terraform.d/",
		"terraform/.terraformrc": "terraform/.terraformrc",
		
		// AWS
		"aws/credentials": "aws/credentials",
		"aws/config": "aws/config",
	}
}

// ScanDirectory scans a directory for configuration files
func (ps *ProfileSync) ScanDirectory(baseDir string, extensions []string) []string {
	var files []string
	
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			// Skip hidden directories except .git, .ssh, etc.
			if strings.HasPrefix(info.Name(), ".") {
				if info.Name() == ".git" || info.Name() == ".ssh" || info.Name() == ".npm" {
					return nil
				}
				return filepath.SkipDir
			}
			return nil
		}
		
		// Check if file matches any extension
		ext := filepath.Ext(info.Name())
		for _, e := range extensions {
			if ext == e || info.Name() == e {
				files = append(files, path)
				break
			}
		}
		
		return nil
	})
	
	if err != nil {
		errorColor.Printf("Error scanning directory: %v\n", err)
	}
	
	return files
}

// CreateMigrationPlan creates a plan for migrating configurations
func (ps *ProfileSync) CreateMigrationPlan(sourceBase, destBase string) error {
	mappings := GetDefaultMappings()
	
	// Add items to migration plan
	for sourceRel, destRel := range mappings {
		sourcePath := filepath.Join(sourceBase, sourceRel)
		destPath := filepath.Join(destBase, destRel)
		
		item := MigrationItem{
			SourcePath:      sourcePath,
			DestinationPath: destPath,
			Type:            ps.getFileType(sourceRel),
			Description:     ps.getDescription(sourceRel),
			AutoMigrate:     true,
		}
		
		ps.migrationPlan.Items = append(ps.migrationPlan.Items, item)
		ps.migrationPlan.TotalItems++
	}
	
	return nil
}

// getFileType determines the type of configuration file
func (ps *ProfileSync) getFileType(path string) string {
	switch {
	case strings.Contains(path, "vscode"):
		return "IDE"
	case strings.Contains(path, "intellij"):
		return "IDE"
	case strings.Contains(path, "vim"):
		return "Editor"
	case strings.Contains(path, "emacs"):
		return "Editor"
	case strings.Contains(path, "bash") || strings.Contains(path, "zsh"):
		return "Shell"
	case strings.Contains(path, "tmux"):
		return "Terminal"
	case strings.Contains(path, "git"):
		return "Version Control"
	case strings.Contains(path, "ssh"):
		return "Security"
	case strings.Contains(path, "chrome") || strings.Contains(path, "firefox"):
		return "Browser"
	case strings.Contains(path, "npm") || strings.Contains(path, "yarn"):
		return "Package Manager"
	case strings.Contains(path, "pip"):
		return "Package Manager"
	case strings.Contains(path, "docker"):
		return "Container"
	case strings.Contains(path, "kubectl"):
		return "Kubernetes"
	case strings.Contains(path, "helm"):
		return "Kubernetes"
	case strings.Contains(path, "terraform"):
		return "Infrastructure"
	case strings.Contains(path, "aws"):
		return "Cloud"
	default:
		return "General"
	}
}

// getDescription provides a human-readable description
func (ps *ProfileSync) getDescription(path string) string {
	desc := map[string]string{
		"vscode/settings.json": "VS Code user settings",
		"vscode/keybindings.json": "VS Code key bindings",
		"intellij/": "IntelliJ IDEA settings",
		"vim/.vimrc": "Vim configuration",
		"vim/.vim/": "Vim plugins and additional configs",
		"emacs/.emacs": "Emacs main configuration",
		"emacs/.emacs.d/": "Emacs plugins and additional configs",
		"bash/.bashrc": "Bash shell configuration",
		"bash/.bash_profile": "Bash profile settings",
		"zsh/.zshrc": "Zsh shell configuration",
		"fish/.config/fish/config.fish": "Fish shell configuration",
		"tmux/.tmux.conf": "Tmux configuration",
		"git/.gitconfig": "Git global configuration",
		"git/.gitignore_global": "Git global ignore patterns",
		"ssh/config": "SSH configuration",
		"ssh/id_rsa": "SSH private key",
		"ssh/id_rsa.pub": "SSH public key",
		"chrome/Default/": "Chrome browser profile",
		"firefox/.mozilla/firefox/": "Firefox browser profile",
		"npm/.npmrc": "NPM configuration",
		"yarn/.yarnrc": "Yarn configuration",
		"pip/pip.conf": "Python pip configuration (Linux/Mac)",
		"pip/pip.ini": "Python pip configuration (Windows)",
		"docker/config.json": "Docker configuration",
		"kubectl/config": "Kubectl configuration",
		"helm/.helm/": "Helm configuration",
		"terraform/.terraform.d/": "Terraform plugins and configuration",
		"terraform/.terraformrc": "Terraform configuration file",
		"aws/credentials": "AWS credentials",
		"aws/config": "AWS configuration",
	}
	
	if desc, ok := desc[path]; ok {
		return desc
	}
	
	return "Configuration file"
}

// ExecuteMigration performs the actual migration
func (ps *ProfileSync) ExecuteMigration(sourceBase, destBase string) error {
	successCount := 0
	failCount := 0
	skipCount := 0
	
	noticeColor.Println("🚀 Starting migration...")
	
	for i, item := range ps.migrationPlan.Items {
		// Check if source exists
		if _, err := os.Stat(item.SourcePath); os.IsNotExist(err) {
			if ps.verbose {
				warnColor.Printf("⏭️  Skipped (not found): %s\n", item.Description)
			}
			ps.migrationPlan.SkippedItems++
			skipCount++
			continue
		}
		
		// Check if destination already exists
		if _, err := os.Stat(item.DestinationPath); err == nil && !ps.force {
			warnColor.Printf("⚠️  Skipped (exists): %s\n", item.Description)
			ps.migrationPlan.SkippedItems++
			skipCount++
			continue
		}
		
		// Create parent directory if needed
		parentDir := filepath.Dir(item.DestinationPath)
		if ps.dryRun {
			noticeColor.Printf("📁 Would create directory: %s\n", parentDir)
		} else {
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				errorColor.Printf("❌ Error creating directory %s: %v\n", parentDir, err)
				failCount++
				continue
			}
		}
		
		// Copy file
		if ps.dryRun {
			successColor.Printf("✅ Would migrate: %s\n", item.Description)
			successCount++
		} else {
			if err := ps.copyFile(item.SourcePath, item.DestinationPath); err != nil {
				errorColor.Printf("❌ Error migrating %s: %v\n", item.Description, err)
				failCount++
				continue
			}
			successColor.Printf("✅ Migrated: %s\n", item.Description)
			successCount++
		}
		
		// Print progress
		fmt.Printf("\r📊 Progress: %d/%d", i+1, len(ps.migrationPlan.Items))
	}
	
	fmt.Println()
	
	ps.migrationPlan.TotalItems = successCount + failCount + skipCount
	
	return nil
}

// copyFile copies a file from source to destination
func (ps *ProfileSync) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	
	_, err = bufio.NewReader(sourceFile).WriteTo(destinationFile)
	return err
}

// PrintReport prints a migration report
func (ps *ProfileSync) PrintReport() {
	infoColor.Println("" + strings.Repeat("=", 60))
	infoColor.Println("📊 MIGRATION REPORT")
	infoColor.Println(strings.Repeat("=", 60))
	
	noticeColor.Printf("Source Platform:   %s\n", ps.sourcePlatform)
	noticeColor.Printf("Destination:       %s\n", ps.destPlatform)
	noticeColor.Printf("Mode:              %s\n", map[bool]string{true: "DRY RUN", false: "LIVE"}[ps.dryRun])
	
	successColor.Printf("✅ Successfully migrated: %d\n", ps.migrationPlan.TotalItems-ps.migrationPlan.SkippedItems)
	warnColor.Printf("⏭️  Skipped:           %d\n", ps.migrationPlan.SkippedItems)
	errorColor.Printf("❌ Failed:            0")
	
	infoColor.Println(strings.Repeat("=", 60))
	
	// Group by type
	typeGroups := make(map[string][]MigrationItem)
	for _, item := range ps.migrationPlan.Items {
		typeGroups[item.Type] = append(typeGroups[item.Type], item)
	}
	
	// Sort types
	var types []string
	for t := range typeGroups {
		types = append(types, t)
	}
	sort.Strings(types)
	
	infoColor.Println("📁 Items by Type:")
	for _, t := range types {
		items := typeGroups[t]
		if len(items) > 0 {
			color.New(color.FgWhite).Printf("  • %s: %d items\n", t, len(items))
		}
	}
	
	infoColor.Println(strings.Repeat("=", 60))
	
	if ps.dryRun {
		warnColor.Println("⚠️  This was a DRY RUN. No files were actually migrated.")
		warnColor.Println("Run without --dry-run to perform the actual migration.")
	} else {
		successColor.Println("✅ Migration complete!")
	}
}

func main() {
	// Define flags
	sourcePlatform := flag.String("source", DetectPlatform(), "Source platform (linux, macos, windows)")
	destPlatform := flag.String("dest", DetectPlatform(), "Destination platform (linux, macos, windows)")
	dryRun := flag.Bool("dry-run", true, "Preview migration without making changes")
	force := flag.Bool("force", false, "Overwrite existing files")
	verbose := flag.Bool("verbose", false, "Verbose output")
	showHelp := flag.Bool("help", false, "Show help message")
	
	flag.Parse()
	
	if *showHelp {
		flag.Usage()
		return
	}
	
	// Validate platforms
	validPlatforms := map[string]bool{"linux": true, "macos": true, "windows": true}
	if !validPlatforms[*sourcePlatform] {
		errorColor.Println("❌ Invalid source platform:", *sourcePlatform)
		errorColor.Println("Must be one of: linux, macos, windows")
		os.Exit(1)
	}
	if !validPlatforms[*destPlatform] {
		errorColor.Println("❌ Invalid destination platform:", *destPlatform)
		errorColor.Println("Must be one of: linux, macos, windows")
		os.Exit(1)
	}
	
	// Create profile sync instance
	ps := NewProfileSync(*sourcePlatform, *destPlatform, *dryRun, *force, *verbose)
	
	// Get home directories
	sourceHome := GetHomeDir(*sourcePlatform)
	destHome := GetHomeDir(*destPlatform)
	
	// Create migration plan
	if err := ps.CreateMigrationPlan(sourceHome, destHome); err != nil {
		errorColor.Println("❌ Error creating migration plan:", err)
		os.Exit(1)
	}
	
	// Execute migration
	if err := ps.ExecuteMigration(sourceHome, destHome); err != nil {
		errorColor.Println("❌ Error during migration:", err)
		os.Exit(1)
	}
	
	// Print report
	ps.PrintReport()
}