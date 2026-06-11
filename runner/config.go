type Config struct {
	// ... existing fields
	DebounceDelay int `toml:"debounce_delay"`
}

func DefaultConfig() *Config {
	return &Config{
		// ... existing defaults
		DebounceDelay: 1000,
	}
}