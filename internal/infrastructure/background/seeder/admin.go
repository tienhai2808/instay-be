package seeder

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

func (s *Seeder) adminSeeder() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error

	schema := os.Getenv("AD_SCHEMA")
	if schema == "" {
		return fmt.Errorf("admin schema is required")
	}

	if schema, err = s.normalizeAndValidateSchema(schema); err != nil {
		return err
	}

	username := os.Getenv("AD_USERNAME")
	if username == "" {
		return fmt.Errorf("admin username is required")
	}

	password := os.Getenv("AD_PASSWORD")
	if password == "" {
		return fmt.Errorf("admin password is required")
	}

	email := os.Getenv("AD_EMAIL")
	if email == "" {
		return fmt.Errorf("admin email is required")
	}

	exists, err := s.userRepo.ExistsActiveAdmin(ctx, schema)
	if err != nil {
		s.logPro.Error("Check active admin failed", port.String("schema", schema), port.Error(err))
		return err
	}
	if exists {
		s.logPro.Info("Admin already exists, skipping", port.String("schema", schema))
		return nil
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		s.logPro.Error("Hash password failed", port.String("schema", schema), port.Error(err))
		return err
	}

	id, err := s.idPro.NextID()
	if err != nil {
		s.logPro.Error("Generate admin id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	admin := model.NewUser(id, username, email, model.RoleAdmin, "Main", "Administrator", "0123456789", hashedPassword, true, nil, nil, nil, nil)

	if err := s.userRepo.Create(ctx, schema, admin); err != nil {
		s.logPro.Error("Create admin failed", port.String("schema", schema), port.Int64("admin_id", admin.ID), port.String("username", admin.Username), port.Error(err))
		return err
	}

	s.logPro.Info("Admin created successfully", port.String("schema", schema), port.Int64("admin_id", admin.ID), port.String("username", admin.Username))

	return nil
}
