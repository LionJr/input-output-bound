package app

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/LionJr/input-output-bound/internal/app/http/server"
	"github.com/LionJr/input-output-bound/internal/config"
	"github.com/LionJr/input-output-bound/internal/handlers"
	"github.com/LionJr/input-output-bound/internal/repositories"
	"github.com/LionJr/input-output-bound/internal/services"
	"github.com/LionJr/input-output-bound/pkg/cache"
)

const shutdownTimeout = 15 * time.Second

type Application struct {
	cfg     *config.AppConfig
	logger  *zap.Logger
	redis   *redis.Client
	handler *handlers.Handler
	http    *server.Server
}

func New(ctx context.Context) (*Application, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	redisClient, err := cache.New(ctx, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return nil, fmt.Errorf("connect redis: %w", err)
	}

	repo := repositories.NewTaskManagerRepository(redisClient)
	service := services.NewTaskManagerService(logger, repo)
	handler := handlers.NewTaskManagerHandler(logger, service)

	return &Application{
		cfg:     cfg,
		logger:  logger,
		redis:   redisClient,
		handler: handler,
		http:    server.New(cfg, logger, handler),
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	a.logger.Info("application started")
	return a.http.Run(ctx)
}

func (a *Application) Shutdown() {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if a.http != nil {
		if err := a.http.Shutdown(shutdownCtx); err != nil {
			a.logger.Error("http shutdown failed", zap.Error(err))
		}
	}

	if a.redis != nil {
		a.logger.Info("closing redis connection")
		if err := a.redis.Close(); err != nil {
			a.logger.Error("redis close failed", zap.Error(err))
		}
	}

	if a.logger != nil {
		a.logger.Info("application stopped")
		if err := a.logger.Sync(); err != nil {
			a.logger.Error("logger sync failed", zap.Error(err))
		}
	}
}
