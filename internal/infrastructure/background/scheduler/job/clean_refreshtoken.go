package job

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type cleanTokenJob struct {
	logPro           port.LogProvider
	refreshTokenRepo repository.RefreshTokenRepository
	schemaRepo       repository.SchemeRepository
}

func NewCleanTokenJob(
	logPro port.LogProvider,
	refreshTokenRepo repository.RefreshTokenRepository,
	schemaRepo repository.SchemeRepository,
) Job {
	return &cleanTokenJob{
		logPro,
		refreshTokenRepo,
		schemaRepo,
	}
}

func (j *cleanTokenJob) Name() string {
	return "clean_refresh_token"
}

func (j *cleanTokenJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	startTime := time.Now()

	schemas, err := j.schemaRepo.FindAll(ctx)
	if err != nil {
		j.logPro.Error("Find all schema failed", port.Error(err))
		return
	}

	totalDeleted := int64(0)
	successCount := 0
	failCount := 0

	for _, schema := range schemas {
		rowsDeleted, err := j.refreshTokenRepo.DeleteAllExpired(ctx, schema)
		if err != nil {
			j.logPro.Error(
				"delete expired tokens failed",
				port.String("schema", schema),
				port.Error(err),
			)
			failCount++
			continue
		}

		totalDeleted += rowsDeleted
		successCount++

		if rowsDeleted > 0 {
			j.logPro.Info(
				"Cleaned up tokens in schema",
				port.String("schema", schema),
				port.Int64("deleted_count", rowsDeleted),
			)
		}
	}

	duration := time.Since(startTime)

	j.logPro.Info(
		"Cleanup expired tokens completed",
		port.Int("total_schemas", len(schemas)),
		port.Int("success_count", successCount),
		port.Int("fail_count", failCount),
		port.Int64("total_deleted", totalDeleted),
		port.Duration("duration", duration),
	)
}
