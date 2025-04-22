package main

import (
	"context"
	application "github.com/LionJr/input-output-bound/internal/app"
	"github.com/LionJr/input-output-bound/internal/config"
	"github.com/LionJr/input-output-bound/internal/handlers"
	"github.com/LionJr/input-output-bound/internal/repositories"
	"github.com/LionJr/input-output-bound/internal/services"
	"github.com/LionJr/input-output-bound/pkg/cache"
	"github.com/chapsuk/grace"
	"go.uber.org/zap"
	"log"
)

func main() {
	ctx := grace.ShutdownContext(context.Background())

	var app application.Application
	defer app.Shutdown()

	var err error

	app.Cfg, err = config.LoadConfig()
	if err != nil {
		log.Println("error while read config", zap.Error(err))
		return
	}

	app.Logger, err = zap.NewProduction()
	if err != nil {
		log.Println("failed to initialize logger", zap.Error(err))
		return
	}

	app.Redis, err = cache.New(ctx, app.Cfg.Redis.Addr, app.Cfg.Redis.Password, app.Cfg.Redis.Db)
	if err != nil {
		app.Logger.Info("error occurred while establishing connection to redis", zap.Error(err))
		return
	}

	repo := repositories.NewTaskManagerRepository(app.Redis)

	service := services.NewTaskManagerService(app.Logger, repo)

	handler := handlers.NewTaskManagerHandler(app.Logger, service)

	app.Run(ctx, handler)
}
