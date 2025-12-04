package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const version = "0.1.0"

// Configuration holds application settings
type config struct {
	wallpapersDir      string
	transitionType     string
	transitionDuration int
}

var (
	cfg          config
	imageFormats = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
)

func init() {
	rand.Seed(time.Now().UnixNano())
	cfg = loadConfig()
}

func main() {
	flags := parseFlags()

	// Override wallpapers directory if specified
	if flags.dir != "" {
		cfg.wallpapersDir = flags.dir
	}

	// Route to appropriate handler
	switch {
	case flags.version:
		printVersion()
	case flags.help:
		printHelp()
	case flags.list:
		listCategories()
	case flags.set != "":
		setWallpaper(flags.set)
	case flags.category != "":
		setRandomFromCategory(flags.category)
	case flags.random:
		handleRandomCommand(flag.Args())
	default:
		handlePositionalArgs(flag.Args())
	}
}

// =============================================================================
// Configuration
// =============================================================================

func loadConfig() config {
	home, _ := os.UserHomeDir()
	defaultDir := filepath.Join(home, "Pictures", "wallpapers")

	return config{
		wallpapersDir:      getEnv("SWWWITCH_WALLPAPERS", defaultDir),
		transitionType:     getEnv("SWWWITCH_TRANSITION_TYPE", "fade"),
		transitionDuration: getEnvInt("SWWWITCH_TRANSITION_DURATION", 1),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		if _, err := fmt.Sscanf(value, "%d", &result); err == nil {
			return result
		}
	}
	return defaultValue
}

// =============================================================================
// Flag Parsing
// =============================================================================

type flags struct {
	list     bool
	random   bool
	set      string
	category string
	dir      string
	help     bool
	version  bool
}

func parseFlags() flags {
	f := flags{}

	flag.BoolVar(&f.list, "list", false, "List available categories")
	flag.BoolVar(&f.list, "l", false, "List available categories (shorthand)")

	flag.BoolVar(&f.random, "random", false, "Set random wallpaper")
	flag.BoolVar(&f.random, "r", false, "Set random wallpaper (shorthand)")

	flag.StringVar(&f.set, "set", "", "Set specific wallpaper file")
	flag.StringVar(&f.set, "s", "", "Set specific wallpaper file (shorthand)")

	flag.StringVar(&f.category, "category", "", "Set random wallpaper from category")
	flag.StringVar(&f.category, "c", "", "Set random wallpaper from category (shorthand)")

	flag.StringVar(&f.dir, "dir", "", "Use custom wallpapers directory")
	flag.StringVar(&f.dir, "d", "", "Use custom wallpapers directory (shorthand)")

	flag.BoolVar(&f.help, "help", false, "Show help")
	flag.BoolVar(&f.help, "h", false, "Show help (shorthand)")

	flag.BoolVar(&f.version, "version", false, "Show version")
	flag.BoolVar(&f.version, "v", false, "Show version (shorthand)")

	flag.Parse()

	return f
}

// =============================================================================
// Command Handlers
// =============================================================================

func handleRandomCommand(args []string) {
	if len(args) > 0 {
		setRandomFromCategory(args[0])
	} else {
		setRandomFromAll()
	}
}

func handlePositionalArgs(args []string) {
	if len(args) == 0 {
		printHelp()
		return
	}

	arg := args[0]

	// Try as file path first
	if fileExists(arg) {
		setWallpaper(arg)
		return
	}

	// Try as category name
	categoryPath := filepath.Join(cfg.wallpapersDir, arg)
	if dirExists(categoryPath) {
		setRandomFromCategory(arg)
		return
	}

	exitWithError("Unknown option or category '%s'", arg)
}

// =============================================================================
// Core Wallpaper Operations
// =============================================================================

func setWallpaper(path string) {
	if !fileExists(path) {
		exitWithError("File not found: %s", path)
	}

	ensureSwwwDaemonRunning()

	cmd := exec.Command("awww", "img", path,
		"--transition-type", cfg.transitionType,
		"--transition-duration", fmt.Sprintf("%d", cfg.transitionDuration))

	if err := cmd.Run(); err != nil {
		exitWithError("Failed to set wallpaper: %v", err)
	}

	fmt.Printf("✓ Wallpaper set to: %s\n", path)
}

func setRandomFromCategory(category string) {
	categoryPath := filepath.Join(cfg.wallpapersDir, category)

	if !dirExists(categoryPath) {
		exitWithError("Category '%s' not found in %s", category, cfg.wallpapersDir)
	}

	wallpapers := findWallpapers(categoryPath)
	if len(wallpapers) == 0 {
		exitWithError("No wallpapers found in category '%s'", category)
	}

	wallpaper := wallpapers[rand.Intn(len(wallpapers))]
	setWallpaper(wallpaper)
	fmt.Printf("Category: %s\n", category)
}

func setRandomFromAll() {
	if !dirExists(cfg.wallpapersDir) {
		exitWithError("Wallpapers directory not found: %s", cfg.wallpapersDir)
	}

	wallpapers := findWallpapers(cfg.wallpapersDir)
	if len(wallpapers) == 0 {
		exitWithError("No wallpapers found in %s", cfg.wallpapersDir)
	}

	wallpaper := wallpapers[rand.Intn(len(wallpapers))]
	setWallpaper(wallpaper)

	// Show category if wallpaper is in a subdirectory
	if relPath, err := filepath.Rel(cfg.wallpapersDir, wallpaper); err == nil {
		if category := filepath.Dir(relPath); category != "." {
			fmt.Printf("Category: %s\n", category)
		}
	}
}

func listCategories() {
	if !dirExists(cfg.wallpapersDir) {
		exitWithError("Wallpapers directory not found: %s", cfg.wallpapersDir)
	}

	fmt.Printf("Available wallpaper categories in %s:\n", cfg.wallpapersDir)

	entries, err := os.ReadDir(cfg.wallpapersDir)
	if err != nil {
		exitWithError("Failed to read directory: %v", err)
	}

	found := false
	for _, entry := range entries {
		if entry.IsDir() {
			count := countWallpapers(filepath.Join(cfg.wallpapersDir, entry.Name()))
			fmt.Printf("  %s (%d wallpapers)\n", entry.Name(), count)
			found = true
		}
	}

	if !found {
		fmt.Println("  (no categories found)")
		fmt.Println()
		fmt.Printf("Create subdirectories in %s to organize your wallpapers by category.\n", cfg.wallpapersDir)
	}
}

// =============================================================================
// File System Utilities
// =============================================================================

func findWallpapers(dir string) []string {
	var wallpapers []string

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}
		if !info.IsDir() && isImageFile(path) {
			wallpapers = append(wallpapers, path)
		}
		return nil
	})

	return wallpapers
}

func countWallpapers(dir string) int {
	return len(findWallpapers(dir))
}

func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, format := range imageFormats {
		if ext == format {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// =============================================================================
// awww Daemon Management
// =============================================================================

func ensureSwwwDaemonRunning() {
	if isSwwwRunning() {
		return
	}

	fmt.Println("Starting awww daemon...")
	cmd := exec.Command("awww-daemon")
	if err := cmd.Start(); err != nil {
		exitWithError("Failed to start awww daemon: %v", err)
	}
	time.Sleep(time.Second) // Give daemon time to start
}

func isSwwwRunning() bool {
	cmd := exec.Command("pgrep", "-x", "awww-daemon")
	return cmd.Run() == nil
}

// =============================================================================
// Output Helpers
// =============================================================================

func printVersion() {
	fmt.Printf("swwwitch v%s\n", version)
}

func printHelp() {
	fmt.Print(`swwwitch - Wallpaper switcher for awww

Usage: swwwitch [OPTION] [CATEGORY/FILE]

Options:
  -l, --list              List available categories
  -r, --random [CAT]      Set random wallpaper from category (or all if no category)
  -s, --set FILE          Set specific wallpaper file
  -c, --category CAT      Set random wallpaper from specific category
  -d, --dir DIR           Use custom wallpapers directory
  -v, --version           Show version
  -h, --help              Show this help

Environment Variables:
  SWWWITCH_WALLPAPERS         Wallpapers directory (default: ~/Pictures/wallpapers)
  SWWWITCH_TRANSITION_TYPE    Transition type (default: fade)
  SWWWITCH_TRANSITION_DURATION Transition duration in seconds (default: 1)

Examples:
  swwwitch --list                    # List all categories
  swwwitch --random                  # Random wallpaper from all categories
  swwwitch --random nature           # Random wallpaper from nature category
  swwwitch --set ~/Pictures/bg.jpg   # Set specific wallpaper
  swwwitch nature                    # Shorthand for --random nature

Directory Structure:
  Your wallpapers directory should contain subdirectories (categories):

  ~/Pictures/wallpapers/
  ├── nature/
  │   ├── mountain1.jpg
  │   └── forest2.png
  ├── abstract/
  │   └── colors.jpg
  └── minimal/
      └── simple.png

`)
}

func exitWithError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	fmt.Fprintln(os.Stderr, "Use 'swwwitch --help' for usage information")
	os.Exit(1)
}
