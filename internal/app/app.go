package app

import (
	"context"
	"github.com/LionJr/input-output-bound/internal/app/http/server"
	"github.com/LionJr/input-output-bound/internal/config"
	"github.com/LionJr/input-output-bound/internal/handlers"
	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"
)

type Application struct {
	Cfg    *config.AppConfig
	Logger *zap.Logger
	Redis  *redis.Client
}

func (app *Application) Run(ctx context.Context, handler *handlers.Handler) {
	httpServerErrCh := server.NewServer(
		ctx,
		app.Cfg,
		app.Logger,
		handler,
	)

	<-httpServerErrCh
}

func (app *Application) Shutdown() {
	if app.Redis != nil {
		app.Logger.Info("Shutting down redis")
		if err := app.Redis.Close(); err != nil {
			app.Logger.Error("Error closing redis connection", zap.Error(err))
		}
	}

	if app.Logger != nil {
		app.Logger.Info("Shutdown logger")
		_ = app.Logger.Sync()
	}
}
