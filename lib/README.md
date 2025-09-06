# OpenX Library

OpenX is a simple Go library for managing and launching applications across different platforms (macOS, Linux, Windows). It provides a clean API to run applications by alias or direct path, kill running applications, and manage application health checks.

## Features

- **Cross-platform**: Works on macOS, Linux, and Windows
- **Simple API**: Just 4 main methods - Run, Kill, Doctor, and ListAliases
- **Flexible execution**: Run applications by alias or direct path
- **Configuration-driven**: Uses YAML config file for application management
- **Health checks**: Built-in doctor functionality to check application availability

## Installation

```bash
go get github.com/muthuishere/openx/lib
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/muthuishere/openx/lib"
)

func main() {
    // Create a new OpenX instance
    ox := lib.New()
    
    // Run an application by alias
    if err := ox.RunAlias("code", "myproject/"); err != nil {
        log.Fatal(err)
    }
    
    // Run an application by direct path
    if err := ox.RunDirect("/Applications/Safari.app"); err != nil {
        log.Fatal(err)
    }
    
    // Kill an application
    if err := ox.Kill("chrome"); err != nil {
        log.Fatal(err)
    }
    
    // Check application health
    if err := ox.Doctor(); err != nil {
        log.Fatal(err)
    }
}
```

## API Reference

### Creating an Instance

```go
// Use default config location (~/.config/openx/config.yaml)
ox := lib.New()

// Use custom config file
ox := lib.NewWithConfig("/path/to/config.yaml")
```

### Main Methods

#### RunAlias(alias string, args ...string) error
Runs an application by its configured alias with optional arguments.

```go
// Launch VS Code
ox.RunAlias("code")

// Launch VS Code with a specific project
ox.RunAlias("code", "myproject/")

// Launch Chrome with a URL
ox.RunAlias("chrome", "https://example.com")
```

#### RunDirect(path string, args ...string) error
Runs an application by its direct path with optional arguments.

```go
// Launch by direct path
ox.RunDirect("/Applications/Safari.app")

// Launch with arguments
ox.RunDirect("/usr/bin/git", "status")
```

#### Kill(alias string) error
Terminates an application by its alias.

```go
// Kill Chrome
ox.Kill("chrome")

// Kill VS Code
ox.Kill("code")
```

#### Doctor() error
Performs a health check on all configured applications and displays the results.

```go
// Run health check
ox.Doctor()
```

#### DoctorJSON() error
Performs a health check and outputs results in JSON format.

```go
// Get health check in JSON format
ox.DoctorJSON()
```

### Alias Management

#### ListAliases() (map[string]string, error)
Returns all configured aliases as a map.

```go
aliases, err := ox.ListAliases()
if err != nil {
    log.Fatal(err)
}

for alias, app := range aliases {
    fmt.Printf("%s -> %s\n", alias, app)
}
```

#### AddAlias(alias, appName string) error
Adds a new alias to the configuration.

```go
// Add a new alias
ox.AddAlias("vs", "code")
```

#### RemoveAlias(alias string) error
Removes an alias from the configuration.

```go
// Remove an alias
ox.RemoveAlias("vs")
```

## Configuration

OpenX uses a YAML configuration file to define applications and their paths across different operating systems.

### Default Config Location
- Linux/macOS: `~/.config/openx/config.yaml`
- Windows: `%APPDATA%\openx\config.yaml`

### Config File Example

```yaml
apps:
  code:
    darwin: "/Applications/Visual Studio Code.app"
    linux: "code"
    windows: "C:\\Program Files\\Microsoft VS Code\\Code.exe"
    kill: ["Code"]
  
  chrome:
    darwin: "/Applications/Google Chrome.app"
    linux: "google-chrome"
    windows: "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
    kill: ["Chrome", "chrome"]
  
  word:
    darwin: "/Applications/Microsoft Word.app"
    linux: "libreoffice --writer"
    windows: "C:\\Program Files\\Microsoft Office\\root\\Office16\\WINWORD.EXE"

aliases:
  browser: chrome
  editor: code
  doc: word
```

### Configuration Structure

- **apps**: Define applications with platform-specific paths
  - `darwin`: macOS path
  - `linux`: Linux command/path
  - `windows`: Windows path
  - `kill`: Process patterns to kill (optional)
- **aliases**: Simple alias mappings to application names

## Platform-Specific Behavior

### macOS
- Supports `.app` bundles
- Uses `open` command as fallback
- Automatically finds executables inside app bundles

### Linux
- Supports executable names in PATH
- Supports absolute paths to executables

### Windows
- Supports `.exe` files
- Supports absolute paths to executables

## Examples

### Basic Usage

```go
package main

import (
    "log"
    "github.com/muthuishere/openx/lib"
)

func main() {
    ox := lib.New()
    
    // Launch applications
    ox.RunAlias("code", ".")
    ox.RunAlias("chrome", "https://github.com")
    
    // Health check
    ox.Doctor()
}
```

### Custom Configuration

```go
package main

import (
    "log"
    "github.com/muthuishere/openx/lib"
)

func main() {
    // Use custom config file
    ox := lib.NewWithConfig("/home/user/my-openx-config.yaml")
    
    // Use the library normally
    ox.RunAlias("myapp")
}
```

### Alias Management

```go
package main

import (
    "fmt"
    "log"
    "github.com/muthuishere/openx/lib"
)

func main() {
    ox := lib.New()
    
    // List all aliases
    aliases, err := ox.ListAliases()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Current aliases:")
    for alias, app := range aliases {
        fmt.Printf("  %s -> %s\n", alias, app)
    }
    
    // Add a new alias
    ox.AddAlias("edit", "code")
    
    // Remove an alias
    ox.RemoveAlias("old-alias")
}
```

## Error Handling

All methods return an error that should be checked:

```go
if err := ox.RunAlias("nonexistent"); err != nil {
    fmt.Printf("Failed to launch: %v\n", err)
}
```

## Version

Current version: **1.0.0**

```go
fmt.Println("Library:", lib.GetName())
fmt.Println("Version:", lib.GetVersion())
```

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
