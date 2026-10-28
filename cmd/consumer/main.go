package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/InstaySystem/is_v2-be/internal/container"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/background/consumer"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctn := container.NewContainer(cfg)
	if err := ctn.InitConsumer(); err != nil {
		log.Fatal(err)
	}
	defer ctn.Cleanup()

	csm := consumer.NewConsumer(ctn.LogPro, ctn.MQPro, ctn.MailerPro, ctn.PNPro, ctn.StoragePro)
	csm.Start()

	ctn.LogPro.Info("Consumer is running")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctn.LogPro.Info("Consumer stopped successfully")
}
