package core

import (
	"runtime"
	"strings"

	"openx/shared/config"
)

/* =========================
   Aliases & Resolution
   ========================= */

type AliasResolver struct {
	config   *config.Config
	synonyms map[string]string // synonym -> alias
}

func newAliasResolver(cfg *config.Config) *AliasResolver {
	ar := &AliasResolver{
		config:   cfg,
		synonyms: map[string]string{},
	}
	ar.initializeSynonyms()
	return ar
}

// initializeSynonyms sets up shorthand aliases
func (a *AliasResolver) initializeSynonyms() {
	// Code Editor shortcuts
	a.synonyms["code"] = "vscode"
	a.synonyms["vs"] = "visualstudio"
	a.synonyms["idea"] = "intellij"
	a.synonyms["ij"] = "intellij"
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
	a.synonyms["z"] = "zed"

	// Browser shortcuts
	a.synonyms["gc"] = "chrome"
	a.synonyms["ff"] = "firefox"
	a.synonyms["br"] = "brave"
	a.synonyms["op"] = "opera"

	// Terminal shortcuts
	a.synonyms["pwsh"] = "powershell"
	a.synonyms["it"] = "iterm"
	a.synonyms["wez"] = "wezterm"
	a.synonyms["al"] = "alacritty"

	// Developer Tool shortcuts
	a.synonyms["pm"] = "postman"
	a.synonyms["ins"] = "insomnia"
	a.synonyms["fig"] = "figma"
	a.synonyms["sk"] = "sketch"
	a.synonyms["not"] = "notion"
	a.synonyms["obs"] = "obsidian"

	// Communication shortcuts
	a.synonyms["sl"] = "slack"
	a.synonyms["dc"] = "discord"
	a.synonyms["tm"] = "teams"
	a.synonyms["zm"] = "zoom"

	// Database shortcuts
	a.synonyms["tp"] = "tableplus"
	a.synonyms["sp"] = "sequel"
	a.synonyms["db"] = "dbeaver"

	// Other shortcuts
	a.synonyms["ad"] = "anydesk"
}

func (a *AliasResolver) Resolve(alias string) (string, bool) {
	base := strings.ToLower(alias)
	if v, ok := a.synonyms[base]; ok {
		base = v
	}

	// Look up in config
	if a.config != nil && a.config.Apps != nil {
		if app, ok := a.config.Apps[base]; ok {
			if target, exists := app.Paths[runtime.GOOS]; exists && target != "" {
				return target, true
			}
		}
	}

	return "", false
}
