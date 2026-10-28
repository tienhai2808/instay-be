package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/InstaySystem/is_v2-be/internal/container"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/background/worker"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ctn := container.NewContainer(cfg)
	if err := ctn.InitWorker(); err != nil {
		log.Fatalln(err)
	}
	defer ctn.Cleanup()

	w := worker.NewWorker(ctn)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctn.LogPro.Info("Worker is running")
	w.Start(ctx)
	ctn.LogPro.Info("Worker stopped successfully")
}
