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

	rdbms rdbms.RDBMS
}

func New(lg *zap.Logger, rdbms rdbms.RDBMS) Repository {
	r := &repository{logger: lg, rdbms: rdbms}

	return r
}

func (r *repository) migrate(direction rdbms.MigrateDirection) error {
	base := "file://internal/repository/"
	precedence := []string{
		"authors",
		"languages",
		"publishers",
		"categories",
		"books",
	}

	for index := 0; index < len(precedence); index++ {
		path := base + precedence[index] + "/migrations"
		if err := r.rdbms.Migrate(path, direction); err != nil {
			r.logger.Error("Error migrating", zap.String("path", path), zap.Error(err))
			return err
		}
	}

	return nil
}

func (r *repository) MigrateUp(ctx context.Context) error {
	return r.migrate(rdbms.MigrateUp)
}

func (r *repository) MigrateDown(ctx context.Context) error {
	return r.migrate(rdbms.MigrateDown)
}
