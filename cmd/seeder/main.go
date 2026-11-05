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
	cfg, err := config.LoadConfig(os.Getenv("SEED_ENV_FILE"))
	if err != nil {
		log.Fatalln(err)
	}

	seedName := os.Getenv("SD_NAME")
	switch seedName {
	case "tenant", "admin", "migrate":
	default:
		log.Fatalln("SD_NAME must be tenant, admin or migrate")
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
