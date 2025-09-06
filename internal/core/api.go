package core

// Public API functions for external use

// LoadConfig loads and returns the current configuration
func LoadConfig() (*Config, error) {
	return loadConfig()
}

// NewAliasResolver creates a new alias resolver with the given aliases
func NewAliasResolver(aliases map[string]string) *AliasResolver {
	resolver := newAliasResolver()
	//resolver.aliases = aliases
	resolver.initializeCanonicals()
	return resolver
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
