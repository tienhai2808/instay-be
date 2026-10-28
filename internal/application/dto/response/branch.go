package response

import "time"

type BasicBranchResponse struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}

type BranchResponse struct {
	ID        int64              `json:"id,string"`
	Name      string             `json:"name"`
	Address   string             `json:"address"`
	Ward      string             `json:"ward"`
	City      string             `json:"city"`
	IsActive  bool               `json:"is_active"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	CreatedBy *BasicUserResponse `json:"created_by"`
	UpdatedBy *BasicUserResponse `json:"updated_by"`
}
