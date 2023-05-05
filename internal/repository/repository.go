package repository

import (
	"context"

	"github.com/CafeKetab/book/internal/models"
	"github.com/CafeKetab/book/pkg/rdbms"
	"go.uber.org/zap"
)

type Repository interface {
	MigrateUp(context.Context) error
	MigrateDown(context.Context) error

	InsertCategory(ctx context.Context, category *models.Category) error
	GetCategoryById(ctx context.Context, id uint64) (*models.Category, error)
	GetCategories(ctx context.Context, encryptedCursor, search string, limit int) ([]models.Category, string, error)
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
	return r.rdbms.Migrate(r.config.MigrationDirectory, rdbms.MigrateDown)
}
