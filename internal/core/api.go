package core

// Public API functions for external use

// LoadConfig loads and returns the current configuration
func LoadConfig() (*Config, error) {
	return loadConfig()
}

// NewAliasResolver creates a new alias resolver with the current config
func NewAliasResolver() (*AliasResolver, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}
	return newAliasResolver(config), nil
}

// Exists checks if a file or directory exists at the given path
func Exists(path string) bool {
	return exists(path)
}

// SaveConfig saves the configuration to the default location
func SaveConfig(config *Config) error {
	return saveConfig(config)
}

// GetAppExists checks if an application exists and is accessible
func GetAppExists(path string) bool {
	return appExists(path)
}
