package runner

type Config struct {
	DebounceDelay int `toml:"debounce_delay"`
}

func DefaultConfig() *Config {
	return &Config{
		DebounceDelay: 1000,
	}
}
