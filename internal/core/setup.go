package core

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// EnsureConfig ensures that the configuration file exists, creating it if necessary
func EnsureConfig() error {
	configPath := getConfigPath()

	// Check if config already exists
	if exists(configPath) {
		return nil
	}

	fmt.Printf("Config not found. Creating starter config at %s\n", configPath)
	return createStarterConfig(configPath)
}

// createStarterConfig creates a starter configuration file for the current OS
func createStarterConfig(configPath string) error {
	// Ensure the config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Get the starter config template for this OS
	template := getStarterTemplate()

	// Write the config file
	if err := os.WriteFile(configPath, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Created starter config with common %s applications.\n", runtime.GOOS)
	fmt.Printf("Edit %s to customize your environment.\n", configPath)

	return nil
}

// getConfigPath returns the path to the configuration file
func getConfigPath() string {
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "openx", "config.yaml")
	}

	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".openx", "config.yaml")
}

// getStarterTemplate returns the starter configuration template for the current OS
func getStarterTemplate() string {
	switch runtime.GOOS {
	case "darwin":
		return getMacOSTemplate()
	case "linux":
		return getLinuxTemplate()
	case "windows":
		return getWindowsTemplate()
	default:
		return getGenericTemplate()
	}
}

// getMacOSTemplate returns the starter config for macOS
func getMacOSTemplate() string {
	return `# openx configuration for macOS
# Edit this file to customize your development environment

apps:
  # Code Editors & IDEs
  vscode:
    darwin: "/Applications/Visual Studio Code.app"
    linux: "code"
    windows: "Code.exe"
  
  goland:
    darwin: "/Applications/GoLand.app"
    linux: "goland"
    windows: "goland64.exe"
  
  intellij:
    darwin: "/Applications/IntelliJ IDEA.app"
    linux: "idea"
    windows: "idea64.exe"
  
  webstorm:
    darwin: "/Applications/WebStorm.app"
    linux: "webstorm"
    windows: "webstorm64.exe"
  
  sublime:
    darwin: "/Applications/Sublime Text.app"
    linux: "subl"
    windows: "subl.exe"
  
  # Browsers
  chrome:
    darwin: "/Applications/Google Chrome.app"
    linux: "google-chrome"
    windows: "chrome.exe"
  
  firefox:
    darwin: "/Applications/Firefox.app"
    linux: "firefox"
    windows: "firefox.exe"
  
  safari:
    darwin: "/Applications/Safari.app"
  
  edge:
    darwin: "/Applications/Microsoft Edge.app"
    linux: "microsoft-edge"
    windows: "msedge.exe"
  
  # Developer Tools
  postman:
    darwin: "/Applications/Postman.app"
    linux: "postman"
    windows: "Postman.exe"
  
  figma:
    darwin: "/Applications/Figma.app"
    linux: "figma-linux"
    windows: "Figma.exe"
  
  # Communication
  slack:
    darwin: "/Applications/Slack.app"
    linux: "slack"
    windows: "slack.exe"
  
  discord:
    darwin: "/Applications/Discord.app"
    linux: "discord"
    windows: "Discord.exe"
  
  # Microsoft Office
  word:
    darwin: "/Applications/Microsoft Word.app"
    linux: "libreoffice --writer"
    windows: "WINWORD.EXE"
  
  excel:
    darwin: "/Applications/Microsoft Excel.app"
    linux: "libreoffice --calc"
    windows: "EXCEL.EXE"
  
  powerpoint:
    darwin: "/Applications/Microsoft PowerPoint.app"
    linux: "libreoffice --impress"
    windows: "POWERPNT.EXE"

aliases:
  code: vscode
  idea: intellij
  ij: intellij
  ws: webstorm
  st: sublime
  gc: chrome
  ff: firefox
  ppt: powerpoint
  pp: powerpoint
`
}

// getLinuxTemplate returns the starter config for Linux
func getLinuxTemplate() string {
	return `# openx configuration for Linux
# Edit this file to customize your development environment

apps:
  # Code Editors & IDEs
  vscode:
    darwin: "/Applications/Visual Studio Code.app"
    linux: "code"
    windows: "Code.exe"
  
  goland:
    darwin: "/Applications/GoLand.app"
    linux: "goland"
    windows: "goland64.exe"
  
  intellij:
    darwin: "/Applications/IntelliJ IDEA.app"
    linux: "idea"
    windows: "idea64.exe"
  
  webstorm:
    darwin: "/Applications/WebStorm.app"
    linux: "webstorm"
    windows: "webstorm64.exe"
  
  sublime:
    darwin: "/Applications/Sublime Text.app"
    linux: "subl"
    windows: "subl.exe"
  
  # Browsers
  chrome:
    darwin: "/Applications/Google Chrome.app"
    linux: "google-chrome"
    windows: "chrome.exe"
  
  firefox:
    darwin: "/Applications/Firefox.app"
    linux: "firefox"
    windows: "firefox.exe"
  
  edge:
    darwin: "/Applications/Microsoft Edge.app"
    linux: "microsoft-edge"
    windows: "msedge.exe"
  
  # Developer Tools
  postman:
    darwin: "/Applications/Postman.app"
    linux: "postman"
    windows: "Postman.exe"
  
  figma:
    darwin: "/Applications/Figma.app"
    linux: "figma-linux"
    windows: "Figma.exe"
  
  # Microsoft Office / Office Suites
  word:
    darwin: "/Applications/Microsoft Word.app"
    linux: "libreoffice --writer"
    windows: "WINWORD.EXE"
  
  excel:
    darwin: "/Applications/Microsoft Excel.app"
    linux: "libreoffice --calc"
    windows: "EXCEL.EXE"
  
  powerpoint:
    darwin: "/Applications/Microsoft PowerPoint.app"
    linux: "libreoffice --impress"
    windows: "POWERPNT.EXE"

aliases:
  code: vscode
  idea: intellij
  ij: intellij
  ws: webstorm
  st: sublime
  gc: chrome
  ff: firefox
  ppt: powerpoint
  pp: powerpoint
`
}

// getWindowsTemplate returns the starter config for Windows
func getWindowsTemplate() string {
	return `# openx configuration for Windows
# Edit this file to customize your development environment

apps:
  # Code Editors & IDEs
  vscode:
    darwin: "/Applications/Visual Studio Code.app"
    linux: "code"
    windows: "Code.exe"
  
  goland:
    darwin: "/Applications/GoLand.app"
    linux: "goland"
    windows: "goland64.exe"
  
  intellij:
    darwin: "/Applications/IntelliJ IDEA.app"
    linux: "idea"
    windows: "idea64.exe"
  
  webstorm:
    darwin: "/Applications/WebStorm.app"
    linux: "webstorm"
    windows: "webstorm64.exe"
  
  sublime:
    darwin: "/Applications/Sublime Text.app"
    linux: "subl"
    windows: "subl.exe"
  
  notepad:
    windows: "notepad.exe"
  
  # Browsers
  chrome:
    darwin: "/Applications/Google Chrome.app"
    linux: "google-chrome"
    windows: "chrome.exe"
  
  firefox:
    darwin: "/Applications/Firefox.app"
    linux: "firefox"
    windows: "firefox.exe"
  
  edge:
    darwin: "/Applications/Microsoft Edge.app"
    linux: "microsoft-edge"
    windows: "msedge.exe"
  
  # Developer Tools
  postman:
    darwin: "/Applications/Postman.app"
    linux: "postman"
    windows: "Postman.exe"
  
  figma:
    darwin: "/Applications/Figma.app"
    linux: "figma-linux"
    windows: "Figma.exe"
  
  # Microsoft Office
  word:
    darwin: "/Applications/Microsoft Word.app"
    linux: "libreoffice --writer"
    windows: "WINWORD.EXE"
  
  excel:
    darwin: "/Applications/Microsoft Excel.app"
    linux: "libreoffice --calc"
    windows: "EXCEL.EXE"
  
  powerpoint:
    darwin: "/Applications/Microsoft PowerPoint.app"
    linux: "libreoffice --impress"
    windows: "POWERPNT.EXE"

aliases:
  code: vscode
  idea: intellij
  ij: intellij
  ws: webstorm
  st: sublime
  gc: chrome
  ff: firefox
  ppt: powerpoint
  pp: powerpoint
`
}

// getGenericTemplate returns a generic starter config
func getGenericTemplate() string {
	return `# openx configuration
# Edit this file to customize your development environment

apps:
  # Add your applications here
  # Format:
  # app-name:
  #   darwin: "/Applications/App.app"     # macOS path
  #   linux: "app-command"                # Linux command
  #   windows: "app.exe"                  # Windows executable

aliases:
  # Add your aliases here
  # Format:
  # alias: app-name
`
}
