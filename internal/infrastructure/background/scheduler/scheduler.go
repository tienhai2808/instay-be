package scheduler

import (
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/background/scheduler/job"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	logPro port.LogProvider
	cron   *cron.Cron
}

func NewScheduler(logPro port.LogProvider) *Scheduler {
	cron := cron.New()

	return &Scheduler{
		logPro,
		cron,
	}
}

func (s *Scheduler) AddJob(schedule string, job job.Job) error {
	if _, err := s.cron.AddFunc(schedule, func() {
		s.logPro.Info("Scheduled job started", port.String("job", job.Name()))
		job.Run()
	}); err != nil {
		return err
	}

	s.logPro.Info("Job scheduled", port.String("job", job.Name()))
	return nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	s.cron.Stop()
}
