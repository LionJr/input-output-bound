package server

import (
	"context"
	"github.com/LionJr/input-output-bound/internal/config"
	"github.com/LionJr/input-output-bound/internal/handlers"
	"github.com/gin-gonic/gin"
	errch "github.com/proxeter/errors-channel"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	cfg     *config.AppConfig
	logger  *zap.Logger
	handler *handlers.Handler
}

func NewServer(
	ctx context.Context,
	cfg *config.AppConfig,
	logger *zap.Logger,
	handler *handlers.Handler,
) <-chan error {
	return errch.Register(func() error {
		return (&Server{
			cfg:     cfg,
			logger:  logger,
			handler: handler,
		}).start(ctx)
	})
}

func (s *Server) start(ctx context.Context) error {
	h := s.initHandlers()

	server := http.Server{
		Handler: h,
		Addr:    ":" + s.cfg.HTTP.Port,
	}

	s.logger.Info(
		"Server running",
		zap.String("host", s.cfg.HTTP.Host),
		zap.String("port", s.cfg.HTTP.Port),
	)

	select {
	case err := <-errch.Register(server.ListenAndServe):
		s.logger.Info("Shutdown input_output_bound_app server", zap.String("by", "error"), zap.Error(err))
		return server.Shutdown(ctx)
	case <-ctx.Done():
		s.logger.Info("Shutdown input_output_bound_app server", zap.String("by", "context.Done"))
		return server.Shutdown(ctx)
	}
}

func (s *Server) initHandlers() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	api := router.Group("/api")
	iOTasksRouter := api.Group("/task-manager")

	iOTasksRouter.POST("/tasks", s.handler.CreateTask)
	iOTasksRouter.GET("/tasks/:id", s.handler.GetTask)

	return router
}
