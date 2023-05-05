package repository

type Config struct {
	CursorSecret       string `koanf:"cursor_secret"`
	MigrationDirectory string `koanf:"migration_directory"`
	Limit              struct {
		Min int `koanf:"min"`
		Max int `koanf:"max"`
	} `koanf:"limit"`
}
