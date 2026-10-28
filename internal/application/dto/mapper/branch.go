package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToBranchQuery(query request.BranchPaginationQuery) model.BranchQuery {
	return model.BranchQuery{
		Page:     query.Page,
		Limit:    query.Limit,
		Search:   query.Search,
		City:     query.City,
		IsActive: query.IsActive,
	}
}

func ToBasicBranchResponse(branch *model.Branch) *response.BasicBranchResponse {
	if branch == nil {
		return nil
	}

	return &response.BasicBranchResponse{
		ID:   branch.ID,
		Name: branch.Name,
	}
}

func ToBasicBranchesResponse(branches []*model.Branch) []*response.BasicBranchResponse {
	if len(branches) == 0 {
		return make([]*response.BasicBranchResponse, 0)
	}

	res := make([]*response.BasicBranchResponse, 0, len(branches))
	for _, branch := range branches {
		res = append(res, ToBasicBranchResponse(branch))
	}

	return res
}

func ToBranchResponse(branch *model.Branch) *response.BranchResponse {
	if branch == nil {
		return nil
	}

	return &response.BranchResponse{
		ID:        branch.ID,
		Name:      branch.Name,
		Address:   branch.Address,
		Ward:      branch.Ward,
		City:      branch.City,
		IsActive:  branch.IsActive,
		CreatedAt: branch.CreatedAt,
		UpdatedAt: branch.UpdatedAt,
		CreatedBy: ToBasicUserResponse(branch.CreatedBy),
		UpdatedBy: ToBasicUserResponse(branch.UpdatedBy),
	}
}

func ToBranchesResponse(branches []*model.Branch) []*response.BranchResponse {
	if len(branches) == 0 {
		return make([]*response.BranchResponse, 0)
	}

	res := make([]*response.BranchResponse, 0, len(branches))
	for _, branch := range branches {
		res = append(res, ToBranchResponse(branch))
	}

	return res
}
