package cmd

import (
	"os"

	"github.com/CafeKetab/book/internal/config"
	"github.com/CafeKetab/book/pkg/logger"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Server struct{}

func (cmd Server) Command(trap chan os.Signal) *cobra.Command {
	run := func(_ *cobra.Command, _ []string) {
		cmd.main(config.Load(true), trap)
	}

	return &cobra.Command{
		Use:   "server",
		Short: "run book server",
		Run:   run,
	}
}

func (cmd *Server) main(cfg *config.Config, trap chan os.Signal) {
	logger := logger.NewZap(cfg.Logger)

	// rdbms, err := rdbms.NewPostgres(cfg.RDBMS)
	// if err != nil {
	// 	logger.Fatal("Error creating rdbms", zap.Error(err))
	// }

	// repository := repository.New(logger, cfg.Repository, rdbms)

	// c1 := &models.Category{Name: "test name 6", Title: "test title 6"}
	// if err := repository.InsertCategory(context.Background(), c1); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(c1)

	// c2, err := repository.GetCategoryById(context.Background(), 7)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(c2)

	// var categories []models.Category
	// var cursor = ""
	// for index := 0; ; index++ {
	// 	categories, cursor, err = repository.GetCategories(context.Background(), cursor, "", 2)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(index, "cursor: ", cursor)
	// 	fmt.Println(index, "categories: ", categories)
	// 	if len(cursor) == 0 {
	// 		break
	// 	}
	// }

	// Keep this at the bottom of the main function
	field := zap.String("signal trap", (<-trap).String())
	logger.Info("exiting by receiving a unix signal", field)
}
