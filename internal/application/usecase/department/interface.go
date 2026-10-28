package department

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type DepartmentUseCase interface {
	Create(ctx context.Context, schema string, userID int64, req request.CreateDepartmentRequest) error

	ListBasic(ctx context.Context, schema string) ([]*response.BasicDepartmentResponse, error)

	ListFull(ctx context.Context, schema string, query request.DepartmentPaginationQuery) ([]*response.DepartmentResponse, *response.MetaResponse, error)

	Update(ctx context.Context, schema string, departmentID, userID int64, req request.UpdateDepartmentRequest) error

	Delete(ctx context.Context, schema string, departmentID int64) error

	GetByName(ctx context.Context, schema, departmentName string) (*response.BasicDepartmentResponse, error)

	GetByID(ctx context.Context, schema string, departmentID int64) (*response.BasicDepartmentResponse, error)
}
