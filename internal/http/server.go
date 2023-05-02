package http

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Server struct {
	config *Config
	logger *zap.Logger
	// repository repository.Repository
	app *fiber.App
}

func New(cfg *Config, log *zap.Logger) *Server {
	server := &Server{config: cfg, logger: log}

	server.app = fiber.New(fiber.Config{JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal})

	v1 := server.app.Group("/v1")
	_ = v1
	// v1.Post("/register", server.register)
	// v1.Post("/login", server.login)
	// v1.Post("/:id<int>", server.fetchUserId, server.user)
	// v1.Get("/me", server.fetchUserId, server.me)
	// v1.Post("/update", server.fetchUserId, server.update)

	return server
}

func (server *Server) Serve() error {
	addr := fmt.Sprintf(":%d", server.config.ListenPort)
	if err := server.app.Listen(addr); err != nil {
		server.logger.Error("error resolving server", zap.Error(err))
		return err
	}
	return nil
}
