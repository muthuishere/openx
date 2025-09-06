# openx - Developer Environment Control Tool

A cross-platform command-line tool for launching and managing developer applications with simple aliases. Perfect for streamlining your development workflow and integrating with task runners like Taskfile, Makefile, or Nix.

## 🚀 Why openx?

- **Unified Interface**: Launch applications consistently across macOS, Linux, and Windows
- **Simple Aliases**: Use short, memorable commands instead of full application paths
- **Cross-Platform**: Same commands work on all platforms with OS-specific configurations
- **Health Monitoring**: Check which apps are available and running with color-coded status
- **Task Integration**: Perfect for build scripts, Taskfiles, and automation
- **Process Management**: Kill applications cleanly when needed
- **Direct Path Support**: Launch any application by path without prior configuration
- **Office Suite Integration**: Built-in support for Microsoft Office and LibreOffice

## 📦 Installation

### Building from Source
```bash
git clone https://github.com/muthuishere/openx.git
cd openx
go build -o openx ./cmd/openx
# Move to your PATH
sudo mv openx /usr/local/bin/
```

## 🎯 Quick Start

### First Run
On first execution, openx creates a starter configuration:
```bash
openx --doctor
```

This creates `~/.config/openx/config.yaml` with common applications pre-configured for your platform.

### Basic Usage
```bash
# Launch VS Code
openx vscode

# Launch Zed editor
openx zed

# Launch IntelliJ IDEA
openx intellij

# Launch with a project directory
openx vscode /path/to/project
openx zed /path/to/project

# Launch using aliases
openx vs /path/to/project    # VS Code
openx z /path/to/project     # Zed
openx ij /path/to/project    # IntelliJ

# Launch browsers
openx chrome
openx firefox
openx arc

# Launch Microsoft Office applications
openx word
openx excel
openx powerpoint

# Using aliases
openx ppt      # PowerPoint
openx pp       # PowerPoint (short alias)

# Launch applications with files
openx word "/Users/username/document.docx"
openx excel "/Users/username/spreadsheet.xlsx"
openx powerpoint "/Users/username/presentation.pptx"

# Launch applications using direct paths (NEW!)
openx "/Applications/Microsoft Word.app" "/Users/username/document.docx"
openx "/Applications/Adobe Photoshop 2024/Adobe Photoshop 2024.app"

# Kill applications
openx --kill chrome firefox arc word excel powerpoint

# Check application health
openx --doctor

# JSON output for scripting
openx --doctor --json
```

## ⚙️ Configuration

The configuration file is located at:
- **macOS/Linux**: `~/.config/openx/config.yaml`
- **Windows**: `%APPDATA%\openx\config.yaml`
- **Custom**: Set `XDG_CONFIG_HOME` environment variable

### Configuration Structure

```yaml
apps:
  # Define applications with OS-specific paths
  vscode:
    darwin: "/Applications/Visual Studio Code.app"
    linux: "code"
    windows: "Code.exe"
    kill: ["Code", "code"]  # Optional: custom kill patterns
  
  zed:
    darwin: "/Applications/Zed.app"
    linux: "zed"
    windows: "zed.exe"
  
  chrome:
    darwin: "/Applications/Google Chrome.app"
    linux: "google-chrome"
    windows: "chrome.exe"
  
  postman:
    darwin: "/Applications/Postman.app"
    linux: "postman"
    windows: "Postman.exe"

aliases:
  # Create short aliases for apps
  vs: vscode
  z: zed
  gc: chrome
  pm: postman
```

### Application Configuration

Each app can specify:
- **OS-specific paths**: `darwin`, `linux`, `windows`
- **Kill patterns**: Custom process names for termination
- **Arguments**: Applications launched with command-line arguments

### Example Configuration

```yaml
apps:
  # IDEs and Editors
  vscode:
    darwin: "/Applications/Visual Studio Code.app"
    linux: "code"
    windows: "Code.exe"
  
  zed:
    darwin: "/Applications/Zed.app"
    linux: "zed"
    windows: "zed.exe"
  
  intellij:
    darwin: "/Applications/IntelliJ IDEA.app"
    linux: "idea"
    windows: "idea64.exe"
  
  goland:
    darwin: "/Applications/GoLand.app"
    linux: "goland"
    windows: "goland64.exe"
  
  cursor:
    darwin: "/Applications/Cursor.app"
    linux: "cursor"
    windows: "Cursor.exe"
  
  # Browsers
  chrome:
    darwin: "/Applications/Google Chrome.app"
    linux: "google-chrome"
    windows: "chrome.exe"
  
  firefox:
    darwin: "/Applications/Firefox.app"
    linux: "firefox"
    windows: "firefox.exe"
  
  arc:
    darwin: "/Applications/Arc.app"
    linux: "arc"
    windows: "Arc.exe"
  
  brave:
    darwin: "/Applications/Brave Browser.app"
    linux: "brave-browser"
    windows: "brave.exe"
  
  # Development Tools
  postman:
    darwin: "/Applications/Postman.app"
    linux: "postman"
    windows: "Postman.exe"
  
  docker:
    darwin: "/Applications/Docker.app"
    linux: "docker-desktop"
    windows: "Docker Desktop.exe"
  
  warp:
    darwin: "/Applications/Warp.app"
    linux: "warp-terminal"
    windows: "Warp.exe"
  
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
  vs: vscode
  z: zed
  ij: intellij
  idea: intellij
  gl: goland
  cr: cursor
  gc: chrome
  ff: firefox
  br: brave
  pm: postman
  term: warp
  ppt: powerpoint
  pp: powerpoint
```

## 🔧 Commands

### Launch Applications
```bash
# Single application
openx vscode

# With arguments
openx vscode /path/to/project

# Using aliases
openx code /path/to/project

# Direct path launching (NEW!)
openx "/Applications/Microsoft Word.app" "/Users/username/document.docx"
openx "/usr/bin/gimp" "/path/to/image.png"
openx "C:\Program Files\Adobe\Photoshop 2024\Photoshop.exe"
```

### Office Documents
```bash
# Launch Office applications with documents
openx word "/Users/username/report.docx"
openx excel "/Users/username/budget.xlsx" 
openx powerpoint "/Users/username/presentation.pptx"

# Using aliases
openx ppt "/Users/username/slides.pptx"
openx pp "/Users/username/demo.pptx"
```

### Kill Applications
```bash
# Kill single app
openx --kill chrome

# Kill multiple apps
openx --kill chrome firefox postman

# Kill using aliases
openx --kill gc ff pm
```

### Health Check
```bash
# Human-readable output
openx --doctor

# JSON output for scripting
openx --doctor --json
```

## 🎪 Integration with Task Runners

### Taskfile.yml Integration

```yaml
version: '3'

vars:
  PROJECT_DIR: "{{.PWD}}"

tasks:
  dev:setup:
    desc: Setup development environment
    cmds:
      - openx vscode {{.PROJECT_DIR}}
      - openx postman
      - openx chrome

  dev:cleanup:
    desc: Cleanup development environment
    cmds:
      - openx --kill vscode postman chrome

  check:tools:
    desc: Check development tools availability
    cmds:
      - openx --doctor

  # Project-specific tasks
  frontend:start:
    desc: Start frontend development
    cmds:
      - openx vscode {{.PROJECT_DIR}}/frontend
      - openx chrome http://localhost:3000

  api:debug:
    desc: Start API debugging session
    cmds:
      - openx goland {{.PROJECT_DIR}}/api
      - openx postman
```

### Makefile Integration

```makefile
.PHONY: dev-setup dev-cleanup check-tools

dev-setup:
	openx vscode $(PWD)
	openx postman
	openx chrome

dev-cleanup:
	openx --kill vscode postman chrome

check-tools:
	openx --doctor

# Project workflow
start-frontend:
	openx vscode $(PWD)/frontend
	openx chrome http://localhost:3000

debug-api:
	openx goland $(PWD)/api
	openx postman
```

### Nix Integration

```nix
# flake.nix
{
  description = "Development environment with openx";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
      
      openx = pkgs.buildGoModule {
        pname = "openx";
        version = "0.1.0";
        src = ./.;
        vendorHash = "sha256-...";
      };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          openx
        ];

        shellHook = ''
          echo "Development environment ready!"
          echo "Available commands:"
          echo "  dev-setup    - Setup development environment"
          echo "  dev-cleanup  - Cleanup development environment"
          echo "  check-tools  - Check tool availability"
          
          alias dev-setup="openx vscode $PWD && openx postman"
          alias dev-cleanup="openx --kill vscode postman"
          alias check-tools="openx --doctor"
        '';
      };
    };
}
```

### Shell Aliases

Add to your `.bashrc`, `.zshrc`, or `.fish` config:

```bash
# Development environment shortcuts
alias dev-setup="openx vs $PWD && openx pm && openx gc"
alias dev-cleanup="openx --kill vscode postman chrome"
alias check-tools="openx --doctor"

# Quick app launches
alias vs="openx vscode"
alias z="openx zed"
alias ij="openx intellij"
alias pm="openx postman"
alias gc="openx chrome"
alias ff="openx firefox"
```

## 🔍 Advanced Usage

### Custom Kill Patterns

Some applications may need custom kill patterns:

```yaml
apps:
  electron-app:
    darwin: "/Applications/MyApp.app"
    linux: "myapp"
    windows: "MyApp.exe"
    kill: ["MyApp Helper", "MyApp", "myapp-renderer"]
```



### Scripting with JSON Output

```bash
#!/bin/bash
# Check if required tools are available
STATUS=$(openx --doctor --json)
MISSING=$(echo "$STATUS" | jq '.summary.missing')

if [ "$MISSING" -gt 0 ]; then
    echo "Warning: $MISSING development tools are missing"
    echo "$STATUS" | jq -r '.apps[] | select(.status == "missing") | .name'
    exit 1
fi

echo "All development tools are available!"
```

## 🎭 Use Cases

### Team Development Workflows
- **Onboarding**: New team members can quickly setup their environment
- **Consistency**: Same commands work across different operating systems
- **Documentation**: Clear, executable environment setup instructions

### Office and Productivity Workflows
- **Document Editing**: Quickly open documents with their appropriate Office applications
- **Cross-Platform Office**: Use LibreOffice on Linux with the same commands as Microsoft Office on Windows/macOS
- **File Association**: Launch applications with specific file types directly
- **Direct Path Access**: Launch any application by path without prior configuration

### Personal Productivity
- **Quick Context Switching**: Launch project-specific tool sets
- **Environment Cleanup**: Easily close all development tools
- **Health Monitoring**: Quickly check what's installed and running

## 🐛 Troubleshooting

### Application Not Found
```bash
# Check if app is configured
openx --doctor

# Verify path in config
cat ~/.config/openx/config.yaml
```

### Permission Issues
```bash
# Make sure openx is executable
chmod +x /usr/local/bin/openx

# Check application permissions
ls -la /Applications/MyApp.app
```

### Process Not Killed
```bash
# Check running processes
openx --doctor --json | jq '.apps[] | select(.running == true)'

# Manual process check
ps aux | grep "MyApp"
```

## 📈 Examples

### Full Development Setup
```bash
# Start full development environment
openx vs ~/my-project          # VS Code
openx ij ~/my-project          # IntelliJ IDEA
openx pm                       # Postman
openx gc http://localhost:3000 # Chrome
openx figma                    # Figma

# Alternative with Zed editor
openx z ~/my-project           # Zed editor
openx arc http://localhost:3000 # Arc browser

# Check everything is running
openx --doctor

# When done, cleanup
openx --kill vscode intellij postman chrome figma
```

### Project-Specific Workflow
```bash
#!/bin/bash
# project-setup.sh
PROJECT_DIR="$1"

if [ -z "$PROJECT_DIR" ]; then
    echo "Usage: $0 <project-directory>"
    exit 1
fi

echo "Setting up development environment for $PROJECT_DIR"

# Launch tools
openx code "$PROJECT_DIR"
openx postman
openx chrome

echo "Development environment ready!"
echo "Run 'openx --kill vscode postman chrome' to cleanup"
```

### Office Document Workflow
```bash
#!/bin/bash
# office-work.sh
DOCUMENT_DIR="$1"

if [ -z "$DOCUMENT_DIR" ]; then
    echo "Usage: $0 <document-directory>"
    exit 1
fi

echo "Setting up office environment for $DOCUMENT_DIR"

# Launch Office applications for different document types
find "$DOCUMENT_DIR" -name "*.docx" -exec openx word {} \;
find "$DOCUMENT_DIR" -name "*.xlsx" -exec openx excel {} \;
find "$DOCUMENT_DIR" -name "*.pptx" -exec openx powerpoint {} \;

echo "Office applications launched!"
echo "Run 'openx --kill word excel powerpoint' to cleanup"
```

### Direct Path Examples
```bash
# Launch any application without configuration
openx "/Applications/Adobe Photoshop 2024/Adobe Photoshop 2024.app"
openx "/Applications/Sketch.app" "/Users/username/design.sketch"
openx "C:\Program Files\Notepad++\notepad++.exe" "config.txt"

# Launch with complex paths
openx "/Applications/Microsoft Word.app" "/Users/username/Documents/Report.docx"
openx "/usr/bin/gimp" "/home/user/images/photo.png"
```


## 📄 License

MIT License - see LICENSE file for details.

---

**openx** - Simplifying developer environment management, one command at a time! 🚀
