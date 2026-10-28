package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

func ToUserEntity(usr *model.User) *entity.UserEntity {
	if usr == nil {
		return nil
	}

	return &entity.UserEntity{
		BaseEntity: entity.BaseEntity{
			ID:        usr.ID,
			CreatedAt: usr.CreatedAt,
			UpdatedAt: usr.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: usr.CreatedByID,
			UpdatedByID: usr.UpdatedByID,
		},
		Username:     usr.Username,
		Email:        usr.Email,
		Role:         string(usr.Role),
		FirstName:    usr.FirstName,
		LastName:     usr.LastName,
		Phone:        usr.Phone,
		Password:     usr.Password,
		Avatar:       usr.Avatar,
		IsActive:     usr.IsActive,
		DepartmentID: usr.DepartmentID,
		Branches:     ToBranchEntities(usr.Branches),
	}
}

func ToUser(usr *entity.UserEntity) *model.User {
	if usr == nil {
		return nil
	}

	return &model.User{
		Base: model.Base{
			ID:        usr.ID,
			CreatedAt: usr.CreatedAt,
			UpdatedAt: usr.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: usr.CreatedByID,
			UpdatedByID: usr.UpdatedByID,
			CreatedBy:   ToUser(usr.CreatedBy),
			UpdatedBy:   ToUser(usr.UpdatedBy),
		},
		Username:     usr.Username,
		Email:        usr.Email,
		Role:         model.UserRole(usr.Role),
		FirstName:    usr.FirstName,
		LastName:     usr.LastName,
		Phone:        usr.Phone,
		Password:     usr.Password,
		Avatar:       usr.Avatar,
		IsActive:     usr.IsActive,
		DepartmentID: usr.DepartmentID,
		Department:   ToDepartment(usr.Department),
		Branches:     ToBranches(usr.Branches),
	}
}

func ToUsers(usrs []*entity.UserEntity) []*model.User {
	if len(usrs) == 0 {
		return make([]*model.User, 0)
	}

	usrsRes := make([]*model.User, 0, len(usrs))
	for _, usr := range usrs {
		usrsRes = append(usrsRes, ToUser(usr))
	}

	return usrsRes
}
