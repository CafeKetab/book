package config

import (
	"github.com/CafeKetab/book/pkg/logger"
	"github.com/CafeKetab/book/pkg/rdbms"
)

type Config struct {
	Logger *logger.Config `koanf:"logger"`
	// Repository *repository.Config
	RDBMS *rdbms.Config `koanf:"rdbms"`
}
