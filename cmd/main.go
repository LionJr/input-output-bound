package main

import (
	"context"
	"log"

	"github.com/chapsuk/grace"

	application "github.com/LionJr/input-output-bound/internal/app"
)

func main() {
	ctx := grace.ShutdownContext(context.Background())

	app, err := application.New(ctx)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}
	defer app.Shutdown()
	if err = app.Run(ctx); err != nil {
		log.Printf("application stopped with error: %v", err)
	}
}
