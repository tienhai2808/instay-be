package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/container"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/background/scheduler"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/background/scheduler/job"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ctn := container.NewContainer(cfg)
	if err := ctn.InitScheduler(); err != nil {
		log.Fatal(err)
	}
	defer ctn.Cleanup()

	sched := scheduler.NewScheduler(ctn.LogPro)

	cleanRefreshTokenJob := job.NewCleanTokenJob(ctn.LogPro, ctn.RefreshTokenRepo, ctn.SchemaRepo)
	rotatePricingCalendarJob := job.NewRotatePricingCalendarJob(ctn.LogPro, ctn.SQL.ORM(), ctn.SchemaRepo)

	if err := sched.AddJob("12 22 * * *", cleanRefreshTokenJob); err != nil {
		ctn.LogPro.Error("Schedule job failed", port.String("job", cleanRefreshTokenJob.Name()), port.Error(err))
		return
	}

	if err := sched.AddJob("0 0 1 1 *", rotatePricingCalendarJob); err != nil {
		ctn.LogPro.Error("Schedule job failed", port.String("job", rotatePricingCalendarJob.Name()), port.Error(err))
		return
	}

	sched.Start()
	ctn.LogPro.Info("Scheduler is running")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sched.Stop()
	ctn.LogPro.Info("Scheduler stopped successfully")
}
