package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToBasicUserResponse(usr *model.User) *response.BasicUserResponse {
	if usr == nil {
		return nil
	}

	return &response.BasicUserResponse{
		ID:        usr.ID,
		Username:  usr.Username,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
	}
}

func ToUserResponse(usr *model.User) *response.UserResponse {
	if usr == nil {
		return nil
	}

	return &response.UserResponse{
		ID:         usr.ID,
		Email:      usr.Email,
		Phone:      usr.Phone,
		Username:   usr.Username,
		FirstName:  usr.FirstName,
		LastName:   usr.LastName,
		Avatar:     usr.Avatar,
		Role:       usr.Role,
		IsActive:   usr.IsActive,
		CreatedAt:  usr.CreatedAt,
		Department: ToSimpleDepartmentResponse(usr.Department),
		Branches:   ToBasicBranchesResponse(usr.Branches),
	}
}

func ToUserDetailsResponse(usr *model.User) *response.UserDetailsResponse {
	if usr == nil {
		return nil
	}

	return &response.UserDetailsResponse{
		ID:         usr.ID,
		Email:      usr.Email,
		Phone:      usr.Phone,
		Username:   usr.Username,
		FirstName:  usr.FirstName,
		LastName:   usr.LastName,
		Avatar:     usr.Avatar,
		Role:       usr.Role,
		IsActive:   usr.IsActive,
		CreatedAt:  usr.CreatedAt,
		UpdatedAt:  usr.UpdatedAt,
		Department: ToBasicDepartmentResponse(usr.Department),
		Branches:   ToBasicBranchesResponse(usr.Branches),
		CreatedBy:  ToBasicUserResponse(usr.CreatedBy),
		UpdatedBy:  ToBasicUserResponse(usr.UpdatedBy),
	}
}

func ToSimpleUserResponse(usr *model.User) *response.SimpleUserResponse {
	if usr == nil {
		return nil
	}

	return &response.SimpleUserResponse{
		ID:         usr.ID,
		FirstName:  usr.FirstName,
		LastName:   usr.LastName,
		Role:       usr.Role,
		IsActive:   usr.IsActive,
		CreatedAt:  usr.CreatedAt,
		Department: ToBasicDepartmentResponse(usr.Department),
		Branches:   ToBasicBranchesResponse(usr.Branches),
	}
}

func ToSimpleUsersResponse(usrs []*model.User) []*response.SimpleUserResponse {
	if len(usrs) == 0 {
		return make([]*response.SimpleUserResponse, 0)
	}

	usersRes := make([]*response.SimpleUserResponse, 0, len(usrs))
	for _, usr := range usrs {
		usersRes = append(usersRes, ToSimpleUserResponse(usr))
	}

	return usersRes
}

func ToUserQuery(query request.UserPaginationQuery) model.UserQuery {
	return model.UserQuery{
		Page:         query.Page,
		Limit:        query.Limit,
		Sort:         query.Sort,
		Order:        query.Order,
		Search:       query.Search,
		Role:         query.Role,
		DepartmentID: query.DepartmentID,
		BranchID:     query.BranchID,
		IsActive:     query.IsActive,
	}
}
