package config

import (
	"github.com/CafeKetab/book/pkg/logger"
	"github.com/CafeKetab/book/pkg/rdbms"
)

func Default() *Config {
	return &Config{
		Logger: &logger.Config{
			Development: true,
			Level:       "debug",
			Encoding:    "console",
		},
		// Repository: &repository.Config{
		// 	CursorSecret:       "A?D(G-KaPdSgVkYp",
		// 	MigrationDirectory: "file://internal/repository/migrations",
		// 	Limit: struct {
		// 		Min int "koanf:\"min\""
		// 		Max int "koanf:\"max\""
		// 	}{
		// 		Min: 2, Max: 48,
		// 	},
		// },
		RDBMS: &rdbms.Config{
			Host:     "localhost",
			Port:     5432,
			Username: "TEST_USER",
			Password: "TEST_PASSWORD",
			Database: "BOOK_DB",
		},
	}
}
