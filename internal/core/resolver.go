package core

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

/* =========================
   Path and URL Utilities
   ========================= */

// isURL checks if the input looks like a URL
func isURL(input string) bool {
	return strings.HasPrefix(input, "http://") ||
		strings.HasPrefix(input, "https://") ||
		strings.HasPrefix(input, "ftp://") ||
		strings.HasPrefix(input, "file://") ||
		strings.Contains(input, "://")
}

// expandTilde expands ~ in file paths
func expandTilde(path string) string {
	if path == "" || path[0] != '~' {
		return path
	}

	if path == "~" || strings.HasPrefix(path, "~/") {
		if home := getHomeDir(); home != "" {
			if path == "~" {
				return home
			}
			return filepath.Join(home, path[2:])
		}
	}

	// Handle ~user syntax on Unix-like systems
	if runtime.GOOS != "windows" {
		sep := strings.Index(path, "/")
		var username, rest string
		if sep == -1 {
			username = path[1:]
		} else {
			username = path[1:sep]
			rest = path[sep+1:]
		}

		if username != "" {
			if u, err := user.Lookup(username); err == nil {
				if rest == "" {
					return u.HomeDir
				}
				return filepath.Join(u.HomeDir, rest)
			}
		}
	}

	return path
}

// expandDot expands . and .. in file paths
func expandDot(path string) string {
	if path == "" {
		return path
	}

	// Handle single dot (current directory)
	if path == "." {
		if cwd, err := os.Getwd(); err == nil {
			return cwd
		}
		return path
	}

	// Handle paths starting with ./
	if strings.HasPrefix(path, "./") {
		if cwd, err := os.Getwd(); err == nil {
			return filepath.Join(cwd, path[2:])
		}
		return path
	}

	// Handle double dot (parent directory)
	if path == ".." {
		if cwd, err := os.Getwd(); err == nil {
			return filepath.Dir(cwd)
		}
		return path
	}

	// Handle paths starting with ../
	if strings.HasPrefix(path, "../") {
		if cwd, err := os.Getwd(); err == nil {
			return filepath.Join(filepath.Dir(cwd), path[3:])
		}
		return path
	}

	return path
}

func getHomeDir() string {
	if home, err := os.UserHomeDir(); err == nil {
		return home
	}
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home
	}
	return ""
}

/* =========================
   File System Utilities
   ========================= */

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}

// isExecutableCandidate checks if a string looks like an application path
func isExecutableCandidate(arg string) bool {
	// Check if it contains path separators (full or relative path to app)
	if strings.ContainsAny(arg, `/\`) {
		// Expand tilde and dot if present
		expanded := expandDot(expandTilde(arg))
		// Check if the file exists and is executable
		if exists(expanded) && isExecutable(expanded) {
			// Additional check: prefer GUI applications over command-line tools
			if runtime.GOOS == "darwin" && strings.HasSuffix(expanded, ".app") {
				return true
			}
			// For other platforms, check if it's likely a GUI app (in common app directories)
			if runtime.GOOS == "linux" && (strings.Contains(expanded, "/opt/") ||
				strings.Contains(expanded, "/usr/share/") ||
				strings.Contains(expanded, "/snap/")) {
				return true
			}
			if runtime.GOOS == "windows" && (strings.Contains(expanded, "Program Files") ||
				strings.Contains(expanded, "AppData")) {
				return true
			}
			// Allow any executable path that user explicitly provides
			return true
		}
	}

	return false
}

/* =========================
   Application Resolution
   ========================= */

// resolveApplication resolves an application alias to executable path
func resolveApplication(appName string) (string, error) {
	ar := newAliasResolver()

	// Try alias resolution first
	if target, ok := ar.Resolve(appName); ok {
		if runtime.GOOS == "darwin" && strings.HasSuffix(target, ".app") {
			return findAppExecutable(target)
		}

		// Try to find in PATH
		if path, err := exec.LookPath(target); err == nil {
			return path, nil
		}

		return target, nil
	}

	// Handle .app bundles on macOS
	if runtime.GOOS == "darwin" && strings.HasSuffix(appName, ".app") {
		return findAppExecutable(appName)
	}

	// Try direct PATH lookup
	if !strings.ContainsAny(appName, `/\`) {
		if path, err := exec.LookPath(appName); err == nil {
			return path, nil
		}
	}

	// Assume it's a direct path
	return appName, nil
}

// findAppExecutable finds the executable within a macOS .app bundle
func findAppExecutable(appName string) (string, error) {
	candidates := []string{
		filepath.Join("/Applications", appName),
		filepath.Join(getHomeDir(), "Applications", appName),
		filepath.Join("/System/Applications", appName),
		appName, // if already a full path
	}

	for _, app := range candidates {
		if !strings.HasSuffix(strings.ToLower(app), ".app") {
			continue
		}

		if !exists(app) {
			continue
		}

		// Try the conventional executable name
		base := strings.TrimSuffix(filepath.Base(app), ".app")
		execPath := filepath.Join(app, "Contents", "MacOS", base)
		if isExecutable(execPath) {
			return execPath, nil
		}

		// Try to find any executable in MacOS directory
		macOSDir := filepath.Join(app, "Contents", "MacOS")
		if entries, err := os.ReadDir(macOSDir); err == nil {
			for _, entry := range entries {
				execPath := filepath.Join(macOSDir, entry.Name())
				if isExecutable(execPath) {
					return execPath, nil
				}
			}
		}
	}

	return "", fmt.Errorf("cannot find executable for %s", appName)
}

/* =========================
   System Opener Functions
   ========================= */

// getSystemOpener returns the system default command to open files/URLs
func getSystemOpener() (string, []string) {
	switch runtime.GOOS {
	case "darwin":
		return "open", []string{}
	case "linux":
		return "xdg-open", []string{}
	case "windows":
		return "cmd", []string{"/c", "start", ""}
	default:
		return "", nil
	}
}

/* =========================
   Target Resolution
   ========================= */

// resolveTarget processes a target (file, URL, directory) and returns the resolved path
func resolveTarget(target string) string {
	// Don't modify URLs
	if isURL(target) {
		return target
	}

	// Expand tilde in file paths
	target = expandTilde(target)

	// Expand dot and double dot in file paths
	target = expandDot(target)

	// Convert to absolute path if it's a relative path
	if !filepath.IsAbs(target) {
		if abs, err := filepath.Abs(target); err == nil {
			target = abs
		}
	}

	return target
}

// resolveTargets processes multiple targets
func resolveTargets(targets []string) []string {
	resolved := make([]string, len(targets))
	for i, target := range targets {
		resolved[i] = resolveTarget(target)
	}
	return resolved
}

/* =========================
   Validation Functions
   ========================= */

// validateTarget checks if a target exists (for files/directories) or is valid (for URLs)
func validateTarget(target string) error {
	if isURL(target) {
		// Basic URL validation - more sophisticated validation could be added
		if !strings.Contains(target, "://") {
			return fmt.Errorf("invalid URL format: %s", target)
		}
		return nil
	}

	// For local paths, check if they exist
	resolved := resolveTarget(target)
	if !exists(resolved) {
		return fmt.Errorf("file or directory does not exist: %s", resolved)
	}

	return nil
}

// validateApplication checks if an application can be resolved
func validateApplication(appName string) error {
	_, err := resolveApplication(appName)
	if err != nil {
		return fmt.Errorf("cannot resolve application '%s': %w", appName, err)
	}
	return nil
}

/* =========================
   Execution Functions
   ========================= */

// runWithApplication runs an application/command with the given arguments
func runWithApplication(appName string, args []string) error {
	// First try to resolve as an application alias
	if execPath, err := resolveApplication(appName); err == nil {
		// For file arguments, resolve paths
		resolvedArgs := make([]string, len(args))
		for i, arg := range args {
			if !isURL(arg) && !strings.HasPrefix(arg, "-") {
				// Check if it looks like a file path (not a flag)
				expanded := expandDot(expandTilde(arg))
				if exists(expanded) {
					resolvedArgs[i] = resolveTarget(arg)
				} else {
					resolvedArgs[i] = arg
				}
			} else {
				resolvedArgs[i] = arg
			}
		}

		cmd := exec.Command(execPath, resolvedArgs...)
		return cmd.Start()
	}

	// If not an alias, try as a direct command
	if path, err := exec.LookPath(appName); err == nil {
		// For direct commands, resolve file paths in arguments
		resolvedArgs := make([]string, len(args))
		for i, arg := range args {
			if !isURL(arg) && !strings.HasPrefix(arg, "-") {
				// Check if it looks like a file path (not a flag)
				expanded := expandDot(expandTilde(arg))
				if exists(expanded) {
					resolvedArgs[i] = resolveTarget(arg)
				} else {
					resolvedArgs[i] = arg
				}
			} else {
				resolvedArgs[i] = arg
			}
		}

		cmd := exec.Command(path, resolvedArgs...)
		return cmd.Start()
	}

	// Last resort: try to run as-is (might be a system command)
	cmd := exec.Command(appName, args...)
	return cmd.Start()
}
