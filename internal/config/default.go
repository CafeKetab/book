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
		RDBMS: &rdbms.Config{
			Host:     "localhost",
			Port:     5432,
			Username: "TEST_USER",
			Password: "TEST_PASSWORD",
			Database: "BOOK_DB",
		},
	}
}
