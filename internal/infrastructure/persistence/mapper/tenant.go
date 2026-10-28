package mapper

import (
	"encoding/json"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

func ToTenantEntity(ta *model.Tenant) *entity.TenantEntity {
	if ta == nil {
		return nil
	}

	settingsJSON, _ := json.Marshal(ta.Settings)

	return &entity.TenantEntity{
		BaseEntity: entity.BaseEntity{
			ID:        ta.ID,
			CreatedAt: ta.CreatedAt,
			UpdatedAt: ta.UpdatedAt,
		},
		Name:       ta.Name,
		Logo:       ta.Logo,
		Phone:      ta.Phone,
		Domain:     ta.Domain,
		IsActive:   ta.IsActive,
		SchemaName: ta.SchemaName,
		Settings:   settingsJSON,
	}
}

func ToTenant(ta *entity.TenantEntity) *model.Tenant {
	if ta == nil {
		return nil
	}

	var settings model.TenantSettings
	_ = json.Unmarshal(ta.Settings, &settings)

	return &model.Tenant{
		Base: model.Base{
			ID:        ta.ID,
			CreatedAt: ta.CreatedAt,
			UpdatedAt: ta.UpdatedAt,
		},
		Name:       ta.Name,
		Logo:       ta.Logo,
		Phone:      ta.Phone,
		Domain:     ta.Domain,
		IsActive:   ta.IsActive,
		SchemaName: ta.SchemaName,
		Settings:   settings,
	}
}

func ToTenants(tas []*entity.TenantEntity) []*model.Tenant {
	if len(tas) == 0 {
		return make([]*model.Tenant, 0)
	}

	tasRes := make([]*model.Tenant, 0, len(tas))
	for _, ta := range tas {
		tasRes = append(tasRes, ToTenant(ta))
	}

	return tasRes
}
