package config

import (
	"github.com/CafeKetab/book/pkg/logger"
	"github.com/CafeKetab/book/pkg/rdbms"
)

type Config struct {
	Logger *logger.Config `koanf:"logger"`
	RDBMS  *rdbms.Config  `koanf:"rdbms"`
}
