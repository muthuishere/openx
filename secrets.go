package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run secrets.go <path-to-.env-file>")
	}

	envFilePath := os.Args[1]

	// Get environment variables for repo and environment
	repo := os.Getenv("GITHUB_REPO")
	environment := os.Getenv("GITHUB_ENVIRONMENT")

	if repo == "" {
		log.Fatal("GITHUB_REPO environment variable is required")
	}

	if environment == "" {
		log.Fatal("GITHUB_ENVIRONMENT environment variable is required")
	}

	fmt.Printf("Setting secrets for repo: %s, environment: %s\n", repo, environment)

	// Read and process .env file
	if err := processEnvFile(envFilePath, repo, environment); err != nil {
		log.Fatalf("Error processing .env file: %v", err)
	}

	fmt.Println("All secrets have been set successfully!")
}

func processEnvFile(envFilePath, repo, environment string) error {
	// Check if file exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		return fmt.Errorf(".env file not found: %s", envFilePath)
	}

	// Open and read .env file
	file, err := os.Open(envFilePath)
	if err != nil {
		return fmt.Errorf("error opening .env file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Skip GITHUB_TOKEN as it's used locally and shouldn't be a GitHub secret
		if key == "GITHUB_TOKEN" {
			fmt.Printf("Skipping %s (local use only)\n", key)
			continue
		}

		// Remove quotes if present
		value = strings.Trim(value, "\"'")

		if err := setGitHubSecret(key, value, repo, environment); err != nil {
			return fmt.Errorf("error setting secret %s: %v", key, err)
		}
	}

	return scanner.Err()
}

func setGitHubSecret(key, value, repo, environment string) error {
	fmt.Printf("Setting secret: %s\n", key)

	// Build gh command
	cmd := exec.Command("gh", "secret", "set", key, "--env", environment, "--repo", repo)

	// Pass the value via stdin
	cmd.Stdin = strings.NewReader(value)

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gh command failed: %v, output: %s", err, string(output))
	}

	fmt.Printf("✓ Successfully set secret: %s\n", key)
	return nil
}
