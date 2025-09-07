package core

import (
	"openx/shared/config"
)

// Re-export types and functions from shared config for backward compatibility
type Config = config.Config
type App = config.App

var loadConfig = config.LoadConfig
var saveConfig = config.SaveConfig
var GetVersion = config.GetVersion
var processNameExceptions = config.ProcessNameExceptions
