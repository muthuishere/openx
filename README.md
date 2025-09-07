# openx

**Your development workflow shouldn't be a clicking marathon.**

Every day, developers waste precious minutes clicking through launchers, typing paths, and switching between apps just to start coding. You've probably written the same "launch my dev environment" script dozens of times, only to rewrite it for each new machine or project.

**There's a better way:**

```bash
# Instead of this chaos:
# ⌘+Space → "VS Code" → Enter → wait → click project folder
# ⌘+Space → "Chrome" → Enter → type localhost:3000
# ⌘+Space → "Postman" → Enter → find collection

# Just do this:
openx vscode . && openx chrome localhost:3000 && openx postman
```

**Works everywhere. Integrates with everything you already use.**

Drop openx into your existing package.json scripts, Taskfile, tmuxp sessions, or shell scripts. Same commands work on Mac, Linux, and Windows. No more platform-specific paths or clicking around.

```json
// package.json - works on any machine
{
  "scripts": {
    "dev": "openx vscode . && openx chrome localhost:3000 && openx postman"
  }
}
```

**Smart enough to handle anything:**
- `openx vscode` → launches VS Code
- `openx README.md` → opens in your default editor  
- `openx https://github.com` → opens in your default browser
- `openx TextEdit myfile.txt` → opens TextEdit with the file  

## Get Started in 30 Seconds

```bash
# Install (choose one)
brew install muthuishere/openx/openx    # macOS/Linux
npm install -g @muthuishere/openx       # Cross-platform

# First run creates config with 50+ popular apps
openx --doctor

# Launch your workflow
openx vscode .
openx chrome localhost:3000  
openx postman

# Clean up when done
openx --kill vscode chrome postman
```

**That's it.** Same commands work everywhere. Add to any script, any workflow, any automation.

## Why This Matters

**You already have the perfect workflow setup.** Your package.json scripts are dialed in. Your Taskfile automates everything. Your tmuxp session is precisely configured.

The only thing missing? A reliable way to launch GUI apps that works the same everywhere.

openx fills that gap. It doesn't replace your tools—it makes them better.

**Taskfile.yml**
```yaml
tasks:
  dev:
    cmds:
      - openx vscode {{.PWD}}
      - openx chrome localhost:3000
      - openx postman
```

**package.json**
```json
{
  "scripts": {
    "launch": "openx code . && openx chrome && openx postman",
    "cleanup": "openx --kill vscode chrome postman"  
  }
}
```

**Shell script**
```bash
#!/bin/bash
openx vscode ~/project
openx chrome localhost:3000
openx postman
```

Works the same whether you're on a MacBook, Ubuntu laptop, or Windows machine.

## 🎯 Built-in Apps & Aliases

**Code Editors & IDEs**:
- `vscode` (`code`, `vs`) - Visual Studio Code
- `zed` (`z`) - Zed editor
- `sublime` (`st`) - Sublime Text  
- `goland` - JetBrains GoLand
- `intellij` (`idea`, `ij`) - IntelliJ IDEA
- `webstorm` (`ws`) - WebStorm
- `pycharm` (`pc`) - PyCharm
- `vim`, `nvim`, `emacs` - Terminal editors

**Browsers**:
- `chrome` (`gc`) - Google Chrome
- `firefox` (`ff`) - Firefox
- `safari` - Safari (macOS)
- `edge` - Microsoft Edge
- `brave` (`br`) - Brave Browser
- `arc` - Arc Browser

**Developer Tools**:
- `postman` (`pm`) - Postman API client
- `docker` - Docker Desktop
- `figma` (`fig`) - Figma design tool
- `insomnia` (`ins`) - Insomnia REST client
- `tableplus` (`tp`) - TablePlus database tool

**Communication & Productivity**:
- `slack` (`sl`) - Slack
- `discord` (`dc`) - Discord  
- `teams` (`tm`) - Microsoft Teams
- `notion` (`not`) - Notion
- `obsidian` (`obs`) - Obsidian

**Microsoft Office**:
- `word` - Microsoft Word / LibreOffice Writer
- `excel` - Microsoft Excel / LibreOffice Calc
- `powerpoint` (`ppt`, `pp`) - PowerPoint / LibreOffice Impress

**Terminals**:
- `terminal` - Default terminal
- `iterm` (`it`) - iTerm2 (macOS)
- `wezterm` (`wez`) - WezTerm
- `alacritty` (`al`) - Alacritty

## ⚙️ Configuration

Auto-generated config at: `~/.openx/config.yaml`

```yaml
apps:
  myapp:
    darwin: "/Applications/MyApp.app"
    linux: "myapp"
    windows: "MyApp.exe"
    kill: ["MyApp", "myapp-helper"]  # Custom kill patterns

aliases:
  ma: myapp
```

### Custom Kill Patterns
```yaml
apps:
  chrome:
    darwin: "/Applications/Google Chrome.app"
    kill: ["Google Chrome", "Chrome Helper", "chrome"]
```

## 🔧 Workflow Integration

### Taskfile.yml
```yaml
version: '3'

tasks:
  dev:
    desc: Launch development environment
    cmds:
      - openx vscode {{.PWD}}
      - openx chrome https://localhost:3000
      - openx postman

  test:
    desc: Run tests with coverage
    cmds:
      - task: test:unit
      - task: test:integration

  test:unit:
    desc: Run unit tests
    dir: internal/core
    cmds:
      - go test -v -cover

  test:integration:
    desc: Run integration tests  
    dir: cmd
    cmds:
      - go test -v

  test:all:
    desc: Run all tests
    cmds:
      - task: test:unit
      - task: test:integration

  cleanup:
    desc: Close development apps
    cmds:
      - openx --kill vscode chrome postman
```

### Package.json Scripts
```json
{
  "scripts": {
    "dev": "openx vscode . && openx chrome http://localhost:3000",
    "launch": "openx code . && openx chrome && openx postman",
    "cleanup": "openx --kill vscode chrome postman"
  }
}
```

### Shell Scripts
```bash
#!/bin/bash
# launch-dev.sh - Full development environment
openx vscode ~/my-project
openx chrome http://localhost:3000
openx postman
openx figma
openx slack

# For cleanup
openx --kill vscode chrome postman figma slack
```

### Nix Flakes
```nix
{
  description = "Development environment with openx";
  
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  
  outputs = { self, nixpkgs }: {
    devShells.default = pkgs.mkShell {
      buildInputs = [ pkgs.openx ];
      
      shellHook = ''
        echo "🚀 Launching development environment..."
        openx vscode $PWD
        openx chrome http://localhost:3000
        openx postman
      '';
    };
  };
}
```

## 📖 Commands Reference

### Basic Usage
```bash
openx <app> [args...]     # Launch app with optional arguments
openx <file-or-url>       # Open file/URL with system default
openx <app> <file>        # Open file with specific app
```

### Process Management
```bash
openx --kill <apps...>    # Close apps (case-insensitive, all instances)
openx --kill chrome firefox postman  # Close multiple apps
```

### System Information
```bash
openx --doctor            # Check all configured apps
openx --doctor --json     # JSON output for automation
```

### Smart Fallbacks
```bash
# These work even if not configured as aliases:
openx README.md                    # System default editor
openx https://github.com          # System default browser  
openx Calculator                  # macOS Calculator app
openx /usr/bin/python3 script.py  # Direct executable with args
```

## 🌟 Key Features

### 🎯 Smart Alias Resolution
- **50+ Built-in Apps**: Popular development tools work out of the box
- **Convenient Shortcuts**: `code` for VS Code, `gc` for Chrome, `pm` for Postman
- **Case-Insensitive**: `openx CHROME` works just like `openx chrome`

### 🔄 Intelligent Fallbacks
- **Single Argument**: Not an alias? Uses system default (`open`, `xdg-open`, `start`)
- **Multiple Arguments**: Treats first as app, rest as arguments
- **Universal Compatibility**: Works with any file, URL, or application

### ⚡ Robust Process Management
- **Case-Insensitive Killing**: Finds and terminates all process variations
- **Multiple Instance Support**: Kills ALL running instances of an app
- **Smart Pattern Matching**: Handles complex app names and helper processes

### 🏗️ Clean Architecture
- **Single Source of Truth**: Configuration templates in setup, no duplication
- **Embedded Versioning**: Self-contained version information
- **Comprehensive Testing**: Unit and integration test coverage

### 🌍 True Cross-Platform
- **macOS**: Full `.app` bundle support, `open -a` commands, Safari handling
- **Linux**: `xdg-open`, `gio open` fallbacks, proper desktop integration  
- **Windows**: `start` command integration, `.exe` handling

## 🧪 Testing

Run the comprehensive test suite:

```bash
# Unit tests
go test ./internal/core -v

# Integration tests  
go test ./cmd/openx -v

# All tests with coverage
go test -cover ./...

# Using Taskfile
task test:all
```

## 📊 Health Check Example

```bash
$ openx --doctor
openx doctor (darwin)
Config: /Users/you/.openx/config.yaml

Applications:
  ✓ chrome          /Applications/Google Chrome.app (running)
    └─ kill: Google Chrome
  ✗ discord         /Applications/Discord.app
    └─ kill: Discord  
  ✓ vscode          /Applications/Visual Studio Code.app (running)
    └─ kill: Code

Aliases:
  code       → vscode
  gc         → chrome
  pm         → postman

Summary:
  Total: 16 apps
  Available: 12
  Missing: 4
  Running: 3
```

## 🤝 Contributing

We welcome contributions! Areas where you can help:

- **New App Definitions**: Add support for more applications
- **Platform Support**: Improve Linux/Windows compatibility  
- **Test Coverage**: Add more edge case testing
- **Documentation**: Improve examples and guides

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

**Ready to stop clicking and start coding?** 

```bash
brew install muthuishere/openx/openx
openx --doctor
openx vscode . && openx chrome && openx postman
```
