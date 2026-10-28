package main

import (
	"log"
	"os"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/container"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/background/seeder"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
)

func main() {
	seedName := os.Getenv("SD_NAME")
	if seedName == "" {
		log.Fatalln("seed name is required")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ctn := container.NewContainer(cfg)
	if err := ctn.InitSeed(); err != nil {
		log.Fatalln(err)
	}
	defer ctn.Cleanup()

	sd := seeder.NewSeeder(ctn.LogPro, ctn.SQL.ORM(), ctn.IDPro, ctn.UserRepo, ctn.SchemaRepo, ctn.PermissionRepo)
	if err = sd.Start(seedName); err != nil {
		ctn.LogPro.Error("Seeder failed", port.String("seed_name", seedName), port.Error(err))
		return
	}

	ctn.LogPro.Info("Seeder finished successfully", port.String("seed_name", seedName))
}
