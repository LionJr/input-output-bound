package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/LionJr/input-output-bound/internal/config"
	"github.com/LionJr/input-output-bound/internal/handlers"
)

type Server struct {
	cfg    *config.AppConfig
	logger *zap.Logger
	srv    *http.Server
}

func New(cfg *config.AppConfig, logger *zap.Logger, handler *handlers.Handler) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
		srv: &http.Server{
			Handler: initHandlers(handler),
			Addr:    ":" + cfg.HTTP.Port,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		s.logger.Info(
			"http server listening",
			zap.String("host", s.cfg.HTTP.Host),
			zap.String("port", s.cfg.HTTP.Port),
		)

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutdown signal received")
		return nil
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("http shutdown: %w", err)
	}
	s.logger.Info("http server stopped")
	return nil
}

func initHandlers(handler *handlers.Handler) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	api := router.Group("/api")
	taskManager := api.Group("/task-manager")

	taskManager.POST("/tasks", handler.CreateTask)
	taskManager.GET("/tasks/:id", handler.GetTask)

	return router
}
