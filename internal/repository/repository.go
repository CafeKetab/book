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
	logger             *zap.Logger
	rdbms              rdbms.RDBMS
	migrationDirectory string
}

func New(lg *zap.Logger, rdbms rdbms.RDBMS) Repository {
	r := &repository{logger: lg, rdbms: rdbms}
	r.migrationDirectory = "file://internal/repository/migrations"

	return r
}

func (r *repository) MigrateUp(ctx context.Context) error {
	return r.rdbms.MigrateUp(r.migrationDirectory)
}

func (r *repository) MigrateDown(ctx context.Context) error {
	return r.rdbms.MigrateDown(r.migrationDirectory)
}
