package main

import (
	"log"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/container"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ctn := container.NewContainer(cfg)
	if err := ctn.InitServer(); err != nil {
		log.Fatalln(err)
	}
	defer ctn.Cleanup()

	sv := api.NewServer(cfg, ctn)

	ch := make(chan error, 1)
	go func() {
		if err := sv.Start(); err != nil {
			ch <- err
		}
	}()

	ctn.LogPro.Info("Server is running", port.Int("port", cfg.Server.Port))

	sv.GracefulShutdown(ch)
}
