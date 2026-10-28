package seeder

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
)

func (s *Seeder) migrateSeeder() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	schemas, err := s.schemaRepo.FindAll(ctx)
	if err != nil {
		s.logPro.Error("Find all schema failed", port.Error(err))
		return err
	}

	for _, schema := range schemas {
		if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(fmt.Sprintf(`SET search_path TO "%s"`, schema)).Error; err != nil {
				s.logPro.Error("Set search_path failed", port.String("schema", schema), port.Error(err))
				return err
			}

			if err := tx.AutoMigrate(autoMigrateModels...); err != nil {
				s.logPro.Error("Auto migrate failed", port.String("schema", schema), port.Error(err))
				return err
			}

			if err := s.seedPermissions(ctx, tx); err != nil {
				s.logPro.Error("Create default permissions failed", port.String("schema", schema), port.Error(err))
				return err
			}

			if err := s.seedSources(ctx, tx); err != nil {
				s.logPro.Error("Create default sources failed", port.String("schema", schema), port.Error(err))
				return err
			}

			if err := s.seedDepartments(ctx, tx); err != nil {
				s.logPro.Error("Create default departments failed", port.String("schema", schema), port.Error(err))
				return err
			}

			return nil
		}); err != nil {
			return err
		}

		s.logPro.Info("Migration completed successfully", port.String("schema", schema))
	}

	return nil
}

func (s *Seeder) seedPermissions(ctx context.Context, tx *gorm.DB) error {
	perms := entity.DefaultPermissions()

	for _, p := range perms {
		var existing entity.PermissionEntity
		err := tx.WithContext(ctx).
			Where("name = ?", p.Name).
			First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		id, err := s.idPro.NextID()
		if err != nil {
			return err
		}

		p.ID = id

		if err := tx.WithContext(ctx).Create(&p).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *Seeder) seedSources(ctx context.Context, tx *gorm.DB) error {
	srcs := entity.DefaultSources()

	for _, src := range srcs {
		var existing entity.SourceEntity
		err := tx.WithContext(ctx).
			Where("name = ?", src.Name).
			First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		id, err := s.idPro.NextID()
		if err != nil {
			return err
		}

		src.ID = id

		if err := tx.WithContext(ctx).Create(&src).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *Seeder) seedDepartments(ctx context.Context, tx *gorm.DB) error {
	depts := entity.DefaultDepartments()

	for _, dept := range depts {
		var existing entity.DepartmentEntity
		err := tx.WithContext(ctx).
			Where("name = ?", dept.Name).
			First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		id, err := s.idPro.NextID()
		if err != nil {
			return err
		}

		dept.ID = id

		if err := tx.WithContext(ctx).Create(&dept).Error; err != nil {
			return err
		}
	}

	return nil
}
