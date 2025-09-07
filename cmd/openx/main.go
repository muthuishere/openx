package main

import (
	"flag"
	"fmt"
	"openx/internal/core"
	"openx/lib"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	var (
		killFlag   = flag.Bool("kill", false, "Kill the specified application(s)")
		doctorFlag = flag.Bool("doctor", false, "Check health status of configured applications")
		jsonFlag   = flag.Bool("json", false, "Output in JSON format (for doctor command)")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] alias [args...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "openx - Developer environment control tool\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  openx alias [args...]     Launch single application by alias\n")
		fmt.Fprintf(os.Stderr, "  openx --kill alias...     Kill application(s) by alias\n")
		fmt.Fprintf(os.Stderr, "  openx --doctor [--json]   Check health of configured apps\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  openx code myproject/      # Launch VS Code with project\n")
		fmt.Fprintf(os.Stderr, "  openx --kill chrome firefox # Kill Chrome and Firefox\n")
		fmt.Fprintf(os.Stderr, "  openx --doctor --json      # Health check in JSON format\n")
		fmt.Fprintf(os.Stderr, "\nLibrary version: %s\n", lib.GetVersion())
	}

	flag.Parse()

	// Create library instance
	ox := lib.New()

	// Ensure config exists
	if err := ox.EnsureConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up config: %v\n", err)
		os.Exit(1)
	}

	// Handle doctor command
	if *doctorFlag {
		var err error
		if *jsonFlag {
			err = ox.DoctorJSON()
		} else {
			err = ox.Doctor()
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Doctor check failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Check for aliases
	aliases := flag.Args()
	if len(aliases) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// Handle kill command
	if *killFlag {
		for _, alias := range aliases {
			if err := ox.Kill(alias); err != nil {
				fmt.Fprintf(os.Stderr, "Error killing %s: %v\n", alias, err)
				os.Exit(1)
			}
		}
		return
	}

	// Handle launch command - single app with arguments
	alias := aliases[0]
	args := aliases[1:]

	// First check if the alias exists in our configuration
	if isValidAlias(alias) {
		// It's a valid alias, use normal launch
		if err := ox.RunAlias(alias, args...); err != nil {
			fmt.Fprintf(os.Stderr, "Error launching %s: %v\n", alias, err)
			os.Exit(1)
		}
	} else {
		// Not a valid alias, use fallback based on arguments
		if len(aliases) == 1 {
			// Single argument - use system default open command
			if err := openWithSystemDefault(alias); err != nil {
				fmt.Fprintf(os.Stderr, "Error opening %s: %v\n", alias, err)
				os.Exit(1)
			}
		} else {
			// Multiple arguments - treat first as app path, rest as args
			if err := openWithAppAndArgs(alias, args); err != nil {
				fmt.Fprintf(os.Stderr, "Error launching %s: %v\n", alias, err)
				os.Exit(1)
			}
		}
	}
}

// isValidAlias checks if the given string is a valid alias in the configuration
func isValidAlias(alias string) bool {
	// Try to load config and check if alias exists
	config, err := core.LoadConfig()
	if err != nil {
		return false
	}

	// Check if it's directly in apps
	if _, exists := config.Apps[strings.ToLower(alias)]; exists {
		return true
	}

	// Check if it's a synonym by trying to create a resolver
	resolver, err := core.NewAliasResolver()
	if err != nil {
		return false
	}

	// Try to resolve - if it resolves, it's valid
	_, resolved := resolver.Resolve(alias)
	return resolved
}

// openWithSystemDefault opens a file or URL using the system's default application
func openWithSystemDefault(target string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", target)
	case "linux":
		// Try xdg-open first, fallback to gio open
		cmd = exec.Command("xdg-open", target)
		if err := cmd.Run(); err != nil {
			cmd = exec.Command("gio", "open", target)
		}
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", target)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}

// openWithAppAndArgs opens using the specified application path with arguments
func openWithAppAndArgs(appPath string, args []string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		// On macOS, use 'open -a' for applications
		cmdArgs := []string{"-a", appPath}
		cmdArgs = append(cmdArgs, args...)
		cmd = exec.Command("open", cmdArgs...)
	case "linux", "windows":
		// On Linux/Windows, execute directly
		cmdArgs := append([]string{appPath}, args...)
		cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Run()
}
