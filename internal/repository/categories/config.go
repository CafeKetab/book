package categories

type Config struct {
	CursorSecret string `koanf:"cursor_secret"`
	Limit        struct {
		Min int `koanf:"min"`
		Max int `koanf:"max"`
	} `koanf:"limit"`
}

func (c Config) Default() *Config {
	return &Config{
		CursorSecret: "",
		Limit: struct {
			Min int "koanf:\"min\""
			Max int "koanf:\"max\""
		}{Min: 12, Max: 48},
	}
}
