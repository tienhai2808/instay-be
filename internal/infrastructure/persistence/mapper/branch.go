package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

func ToBranchEntity(branch *model.Branch) *entity.BranchEntity {
	if branch == nil {
		return nil
	}

	return &entity.BranchEntity{
		BaseEntity: entity.BaseEntity{
			ID:        branch.ID,
			CreatedAt: branch.CreatedAt,
			UpdatedAt: branch.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: branch.CreatedByID,
			UpdatedByID: branch.UpdatedByID,
		},
		Name:     branch.Name,
		Address:  branch.Address,
		Ward:     branch.Ward,
		City:     branch.City,
		IsActive: branch.IsActive,
	}
}

func ToBranchEntities(branches []*model.Branch) []*entity.BranchEntity {
	if len(branches) == 0 {
		return make([]*entity.BranchEntity, 0)
	}

	res := make([]*entity.BranchEntity, 0, len(branches))
	for _, branch := range branches {
		res = append(res, ToBranchEntity(branch))
	}

	return res
}

func ToBranch(branch *entity.BranchEntity) *model.Branch {
	if branch == nil {
		return nil
	}

	return &model.Branch{
		Base: model.Base{
			ID:        branch.ID,
			CreatedAt: branch.CreatedAt,
			UpdatedAt: branch.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: branch.CreatedByID,
			UpdatedByID: branch.UpdatedByID,
			CreatedBy:   ToUser(branch.CreatedBy),
			UpdatedBy:   ToUser(branch.UpdatedBy),
		},
		Name:     branch.Name,
		Address:  branch.Address,
		Ward:     branch.Ward,
		City:     branch.City,
		IsActive: branch.IsActive,
	}
}

func ToBranches(branches []*entity.BranchEntity) []*model.Branch {
	if len(branches) == 0 {
		return make([]*model.Branch, 0)
	}

	res := make([]*model.Branch, 0, len(branches))
	for _, branch := range branches {
		res = append(res, ToBranch(branch))
	}

	return res
}
