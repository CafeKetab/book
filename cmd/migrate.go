package cmd

import (
	"context"
	"os"

	"github.com/CafeKetab/book/internal/config"
	"github.com/CafeKetab/book/internal/repository"
	"github.com/CafeKetab/book/pkg/logger"
	"github.com/CafeKetab/book/pkg/rdbms"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Migrate struct{}

func (m Migrate) Command(trap chan os.Signal) *cobra.Command {
	run := func(_ *cobra.Command, args []string) {
		m.main(config.Load(true), args, trap)
	}

	return &cobra.Command{
		Use:       "migrate",
		Short:     "run book related models migrations",
		Run:       run,
		Args:      cobra.OnlyValidArgs,
		ValidArgs: []string{"up", "down"},
	}
}

func (m *Migrate) main(cfg *config.Config, args []string, trap chan os.Signal) {
	logger := logger.NewZap(cfg.Logger)

	if len(args) != 1 {
		logger.Fatal("invalid arguments given", zap.Any("args", args))
	}

	rdbms, err := rdbms.NewPostgres(cfg.RDBMS)
	if err != nil {
		logger.Fatal("Error creating rdbms", zap.Error(err))
	}

	repository := repository.New(logger, cfg.Repository, rdbms)
	repository.MigrateUp(context.Background())

	var callsMigrator func(context.Context) error
	if args[0] == "up" {
		callsMigrator = repository.MigrateUp
	} else {
		callsMigrator = repository.MigrateDown
	}

	if err := callsMigrator(context.Background()); err != nil {
		logger.Fatal("Error migrate database", zap.String("migration", args[0]), zap.Error(err))
	}

	logger.Info("Database has been migrated successfully", zap.String("migration", args[0]))
}
