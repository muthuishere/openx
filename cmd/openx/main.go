package main

import (
	"flag"
	"fmt"
	"openx/internal/core"
	"os"
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
	}

	flag.Parse()

	// Ensure config exists
	if err := core.EnsureConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up config: %v\n", err)
		os.Exit(1)
	}

	// Handle doctor command
	if *doctorFlag {
		if err := core.RunDoctor(*jsonFlag); err != nil {
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
			if err := core.CloseApp(alias); err != nil {
				fmt.Fprintf(os.Stderr, "Error killing %s: %v\n", alias, err)
				os.Exit(1)
			}
		}
		return
	}

	// Handle launch command - single app with arguments
	alias := aliases[0]
	args := aliases[1:]
	if err := core.LaunchApp(alias, args); err != nil {
		fmt.Fprintf(os.Stderr, "Error launching %s: %v\n", alias, err)
		os.Exit(1)
	}
}
