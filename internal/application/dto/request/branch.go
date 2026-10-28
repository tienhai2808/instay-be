package request

type BranchPaginationQuery struct {
	Page     uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit    uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Search   string `form:"search" json:"search"`
	City     string `form:"city" json:"city"`
	IsActive *bool  `form:"is_active" binding:"omitempty" json:"is_active"`
}

type CreateBranchRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=150"`
	Address  string `json:"address" binding:"required,min=2,max=255"`
	Ward     string `json:"ward" binding:"required,min=1,max=150"`
	City     string `json:"city" binding:"required,min=1,max=150"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

type UpdateBranchRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=150"`
	Address  string `json:"address" binding:"required,min=2,max=255"`
	Ward     string `json:"ward" binding:"required,min=1,max=150"`
	City     string `json:"city" binding:"required,min=1,max=150"`
	IsActive *bool  `json:"is_active" binding:"required"`
}
