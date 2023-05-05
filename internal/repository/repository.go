package repository

import (
	"context"

	"github.com/CafeKetab/book/pkg/rdbms"
	"go.uber.org/zap"
)

type Repository interface {
	MigrateUp(context.Context) error

	MigrateDown(context.Context) error
}

type repository struct {
	logger *zap.Logger
	config *Config
	rdbms  rdbms.RDBMS
}

func New(lg *zap.Logger, cfg *Config, rdbms rdbms.RDBMS) Repository {
	r := &repository{logger: lg, config: cfg, rdbms: rdbms}

	return r
}

func (r *repository) MigrateUp(ctx context.Context) error {
	return r.rdbms.Migrate(r.config.MigrationDirectory, rdbms.MigrateUp)
}

func (r *repository) MigrateDown(ctx context.Context) error {
	return r.rdbms.Migrate(r.config.MigrationDirectory, rdbms.MigrateUp)
}
