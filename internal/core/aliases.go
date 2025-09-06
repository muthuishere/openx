package core

import (
	"runtime"
	"strings"
)

/* =========================
   Aliases & Resolution
   ========================= */

type AliasResolver struct {
	canonicals map[string]map[string]string // alias -> {os -> target}
	synonyms   map[string]string            // synonym -> alias
}

func newAliasResolver() *AliasResolver {
	ar := &AliasResolver{
		canonicals: map[string]map[string]string{},
		synonyms:   map[string]string{},
	}
	ar.initializeCanonicals()
	return ar
}

// initializeCanonicals sets up canonical application mappings
func (a *AliasResolver) initializeCanonicals() {
	// Code Editors & IDEs
	a.canonicals["vscode"] = map[string]string{
		"darwin":  "Visual Studio Code.app",
		"linux":   "code",
		"windows": "Code.exe",
	}
	a.canonicals["visualstudio"] = map[string]string{
		"darwin":  "Visual Studio.app",
		"linux":   "code", // VS Code on Linux
		"windows": "devenv.exe",
	}
	a.canonicals["notepad"] = map[string]string{
		"windows": "notepad.exe",
	}
	a.canonicals["notepadpp"] = map[string]string{
		"windows": "notepad++.exe",
	}
	a.canonicals["zed"] = map[string]string{
		"darwin":  "Zed.app",
		"linux":   "zed",
		"windows": "zed.exe",
	}
	a.canonicals["sublime"] = map[string]string{
		"darwin":  "Sublime Text.app",
		"linux":   "subl",
		"windows": "subl.exe",
	}
	a.canonicals["atom"] = map[string]string{
		"darwin":  "Atom.app",
		"linux":   "atom",
		"windows": "atom.exe",
	}
	a.canonicals["vim"] = map[string]string{
		"darwin":  "vim",
		"linux":   "vim",
		"windows": "vim.exe",
	}
	a.canonicals["nvim"] = map[string]string{
		"darwin":  "nvim",
		"linux":   "nvim",
		"windows": "nvim.exe",
	}
	a.canonicals["emacs"] = map[string]string{
		"darwin":  "Emacs.app",
		"linux":   "emacs",
		"windows": "emacs.exe",
	}

	// JetBrains IDEs
	a.canonicals["goland"] = map[string]string{
		"darwin":  "GoLand.app",
		"linux":   "goland",
		"windows": "goland64.exe",
	}
	a.canonicals["intellij"] = map[string]string{
		"darwin":  "IntelliJ IDEA.app",
		"linux":   "idea",
		"windows": "idea64.exe",
	}
	a.canonicals["webstorm"] = map[string]string{
		"darwin":  "WebStorm.app",
		"linux":   "webstorm",
		"windows": "webstorm64.exe",
	}
	a.canonicals["phpstorm"] = map[string]string{
		"darwin":  "PhpStorm.app",
		"linux":   "phpstorm",
		"windows": "phpstorm64.exe",
	}
	a.canonicals["clion"] = map[string]string{
		"darwin":  "CLion.app",
		"linux":   "clion",
		"windows": "clion64.exe",
	}
	a.canonicals["pycharm"] = map[string]string{
		"darwin":  "PyCharm.app",
		"linux":   "pycharm",
		"windows": "pycharm64.exe",
	}
	a.canonicals["rider"] = map[string]string{
		"darwin":  "Rider.app",
		"linux":   "rider",
		"windows": "rider64.exe",
	}
	a.canonicals["datagrip"] = map[string]string{
		"darwin":  "DataGrip.app",
		"linux":   "datagrip",
		"windows": "datagrip64.exe",
	}
	a.canonicals["rubymine"] = map[string]string{
		"darwin":  "RubyMine.app",
		"linux":   "rubymine",
		"windows": "rubymine64.exe",
	}
	a.canonicals["appcode"] = map[string]string{
		"darwin": "AppCode.app",
	}

	// Browsers
	a.canonicals["chrome"] = map[string]string{
		"darwin":  "Google Chrome.app",
		"linux":   "google-chrome",
		"windows": "chrome.exe",
	}
	a.canonicals["firefox"] = map[string]string{
		"darwin":  "Firefox.app",
		"linux":   "firefox",
		"windows": "firefox.exe",
	}
	a.canonicals["safari"] = map[string]string{
		"darwin": "Safari.app",
	}
	a.canonicals["edge"] = map[string]string{
		"darwin":  "Microsoft Edge.app",
		"linux":   "microsoft-edge",
		"windows": "msedge.exe",
	}
	a.canonicals["brave"] = map[string]string{
		"darwin":  "Brave Browser.app",
		"linux":   "brave-browser",
		"windows": "brave.exe",
	}
	a.canonicals["opera"] = map[string]string{
		"darwin":  "Opera.app",
		"linux":   "opera",
		"windows": "opera.exe",
	}

	// Terminal Applications
	a.canonicals["terminal"] = map[string]string{
		"darwin":  "Terminal.app",
		"linux":   "gnome-terminal",
		"windows": "wt.exe",
	}
	a.canonicals["iterm"] = map[string]string{
		"darwin": "iTerm.app",
	}
	a.canonicals["powershell"] = map[string]string{
		"darwin":  "pwsh",
		"linux":   "pwsh",
		"windows": "powershell.exe",
	}
	a.canonicals["wezterm"] = map[string]string{
		"darwin":  "WezTerm.app",
		"linux":   "wezterm",
		"windows": "wezterm.exe",
	}
	a.canonicals["alacritty"] = map[string]string{
		"darwin":  "Alacritty.app",
		"linux":   "alacritty",
		"windows": "alacritty.exe",
	}

	// Developer Tools
	a.canonicals["xcode"] = map[string]string{
		"darwin": "Xcode.app",
	}
	a.canonicals["docker"] = map[string]string{
		"darwin":  "Docker.app",
		"linux":   "docker-desktop",
		"windows": "Docker Desktop.exe",
	}
	a.canonicals["postman"] = map[string]string{
		"darwin":  "Postman.app",
		"linux":   "postman",
		"windows": "Postman.exe",
	}
	a.canonicals["insomnia"] = map[string]string{
		"darwin":  "Insomnia.app",
		"linux":   "insomnia",
		"windows": "Insomnia.exe",
	}
	a.canonicals["figma"] = map[string]string{
		"darwin":  "Figma.app",
		"linux":   "figma-linux",
		"windows": "Figma.exe",
	}
	a.canonicals["sketch"] = map[string]string{
		"darwin": "Sketch.app",
	}
	a.canonicals["notion"] = map[string]string{
		"darwin":  "Notion.app",
		"linux":   "notion-app",
		"windows": "Notion.exe",
	}
	a.canonicals["obsidian"] = map[string]string{
		"darwin":  "Obsidian.app",
		"linux":   "obsidian",
		"windows": "Obsidian.exe",
	}
	a.canonicals["slack"] = map[string]string{
		"darwin":  "Slack.app",
		"linux":   "slack",
		"windows": "slack.exe",
	}
	a.canonicals["discord"] = map[string]string{
		"darwin":  "Discord.app",
		"linux":   "discord",
		"windows": "Discord.exe",
	}
	a.canonicals["teams"] = map[string]string{
		"darwin":  "Microsoft Teams.app",
		"linux":   "teams",
		"windows": "Teams.exe",
	}
	a.canonicals["zoom"] = map[string]string{
		"darwin":  "zoom.us.app",
		"linux":   "zoom",
		"windows": "Zoom.exe",
	}
	a.canonicals["anydesk"] = map[string]string{
		"darwin":  "AnyDesk.app",
		"linux":   "anydesk",
		"windows": "AnyDesk.exe",
	}

	// File Managers
	a.canonicals["finder"] = map[string]string{
		"darwin": "Finder.app",
	}
	a.canonicals["explorer"] = map[string]string{
		"windows": "explorer.exe",
	}
	a.canonicals["nautilus"] = map[string]string{
		"linux": "nautilus",
	}
	a.canonicals["ranger"] = map[string]string{
		"darwin":  "ranger",
		"linux":   "ranger",
		"windows": "ranger.exe",
	}

	// Database Tools
	a.canonicals["tableplus"] = map[string]string{
		"darwin":  "TablePlus.app",
		"linux":   "tableplus",
		"windows": "TablePlus.exe",
	}
	a.canonicals["sequel"] = map[string]string{
		"darwin": "Sequel Pro.app",
	}
	a.canonicals["dbeaver"] = map[string]string{
		"darwin":  "DBeaver.app",
		"linux":   "dbeaver",
		"windows": "dbeaver.exe",
	}

	// synonyms
	a.synonyms["code"] = "vscode"
	a.synonyms["vs"] = "visualstudio"
	a.synonyms["idea"] = "intellij"
	a.synonyms["ij"] = "intellij"
	a.synonyms["gc"] = "chrome"
	a.synonyms["ws"] = "webstorm"
	a.synonyms["ps"] = "phpstorm"
	a.synonyms["cl"] = "clion"
	a.synonyms["pc"] = "pycharm"
	a.synonyms["rd"] = "rider"
	a.synonyms["dg"] = "datagrip"
	a.synonyms["rm"] = "rubymine"
	a.synonyms["ac"] = "appcode"
	a.synonyms["st"] = "sublime"
	a.synonyms["vi"] = "vim"
	a.synonyms["npp"] = "notepadpp"
	a.synonyms["ff"] = "firefox"
	a.synonyms["br"] = "brave"
	a.synonyms["op"] = "opera"
	a.synonyms["pwsh"] = "powershell"
	a.synonyms["it"] = "iterm"
	a.synonyms["wez"] = "wezterm"
	a.synonyms["al"] = "alacritty"
	a.synonyms["pm"] = "postman"
	a.synonyms["ins"] = "insomnia"
	a.synonyms["fig"] = "figma"
	a.synonyms["sk"] = "sketch"
	a.synonyms["not"] = "notion"
	a.synonyms["obs"] = "obsidian"
	a.synonyms["sl"] = "slack"
	a.synonyms["dc"] = "discord"
	a.synonyms["tm"] = "teams"
	a.synonyms["zm"] = "zoom"
	a.synonyms["tp"] = "tableplus"
	a.synonyms["sp"] = "sequel"
	a.synonyms["db"] = "dbeaver"
	a.synonyms["ad"] = "anydesk"
	a.synonyms["z"] = "zed"
}

func (a *AliasResolver) Resolve(alias string) (string, bool) {
	base := strings.ToLower(alias)
	if v, ok := a.synonyms[base]; ok {
		base = v
	}
	if m, ok := a.canonicals[base]; ok {
		if t, ok2 := m[runtime.GOOS]; ok2 {
			return t, true
		}
	}
	return "", false
}
