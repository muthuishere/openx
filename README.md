# openx - Stop Clicking Around, Start Coding

**Tired of this?**
- `⌘+Space` → "VS Code" → `Enter` → wait... → click on project folder
- Opening browser → typing `http://localhost:3000` → opening Postman → finding that API collection
- Switching between 5 different apps just to start coding
- Different commands/paths on Mac vs Linux vs Windows
- Writing the same app-launching script for every project

**Just do this instead:**
```bash
openx vscode . && openx chrome && openx postman
```

**You can also add it to your workflow:**
```json
// package.json
{
  "scripts": {
    "launch-my-env": "openx vscode . && openx chrome && openx postman"
  }
}
```

Cross-platform GUI app launcher that actually works the same everywhere. No more clicking, no more platform-specific scripts.

> 💡 See workflow examples for [Taskfile](#taskfile), [Nix](#nix-flakes), [tmuxp](#tmuxp) below

## 📦 Installation

### Homebrew (macOS/Linux)
```bash
# Add the tap
brew tap muthuishere/openx https://github.com/muthuishere/homebrew-openx

# Install openx
brew install muthuishere/openx/openx
```

### NPM (Global)
```bash
npm install -g @muthuishere/openx
```

### Direct Download
Download the latest release from [GitHub Releases](https://github.com/muthuishere/openx/releases)

### Building from Source
```bash
git clone https://github.com/muthuishere/openx.git
cd openx
go build -o openx ./cmd/openx
```

## 🚀 Quick Start

```bash
# First run - creates config with built-in apps
openx --doctor

# Launch apps
openx vscode           # VS Code
openx chrome           # Chrome browser
openx word             # Microsoft Word
openx postman          # Postman

# Using aliases
openx vs               # VS Code
openx gc               # Chrome
openx ppt              # PowerPoint

# Open with files/URLs
openx vscode ~/project
openx chrome https://github.com
openx word document.docx

# Close apps
openx --kill vscode chrome word

# Check what's available
openx --doctor
```

## 🎯 Built-in Apps

Works out of the box with common developer tools:

**IDEs & Editors**: VS Code, Zed, IntelliJ, GoLand, Cursor  
**Browsers**: Chrome, Firefox, Arc, Brave, Safari  
**Dev Tools**: Postman, Docker, Warp Terminal  
**Office**: Word, Excel, PowerPoint (Microsoft Office on macOS/Windows, LibreOffice on Linux)  

## ⚙️ Configuration

Config file: `~/.config/openx/config.yaml`

```yaml
apps:
  myapp:
    darwin: "/Applications/MyApp.app"
    linux: "myapp"
    windows: "MyApp.exe"

aliases:
  ma: myapp
```

## 🔧 Workflow Integration

### Taskfile
```yaml
tasks:
  launch-dev:
    cmds:
      - openx vscode {{.PWD}}
      - openx chrome
      - openx postman

  cleanup:
    cmds:
      - openx --kill vscode chrome postman
```

### Shell Scripts
```bash
# setup.sh
openx vscode ~/project
openx chrome localhost:3000
openx postman

# cleanup.sh  
openx --kill vscode chrome postman
```

### Nix Flakes
```nix
shellHook = ''
  openx vscode $PWD
  openx chrome
'';
```

### tmuxp
```yaml
session_name: launch-dev
windows:
  - window_name: main
    shell_command_before:
      - openx vscode ~/project
      - openx chrome
      - openx postman
```

## 📖 Commands

```bash
openx <app> [args...]     # Launch app
openx --kill <apps...>    # Close apps  
openx --doctor            # Check available apps
openx --doctor --json     # JSON output
```

## 🌟 Why openx?

- **Cross-platform**: Same commands work on macOS, Linux, Windows
- **Simple**: Just `openx app` to launch anything
- **Fast**: Perfect for automation and workflows
- **Flexible**: Use built-in apps or configure your own
- **Lightweight**: Single binary, no dependencies

Perfect for developers who want to quickly launch their IDE, browser, API client, and other tools without clicking around or remembering complex paths.

## 📄 License

MIT License
