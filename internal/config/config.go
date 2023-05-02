package config

import (
	"github.com/CafeKetab/book/pkg/logger"
)

type Config struct {
	Logger *logger.Config `koanf:"logger"`
}
