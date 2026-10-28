package model

type Floor struct {
	Base
	Audit

	Name     string
	BranchID int64

	Branch *Branch
	Rooms  []*Room
}

type FloorQuery struct {
	Page      uint32
	Limit     uint32
	Search    string
	BranchID  int64
	BranchIDs []int64
}

func NewFloor(
	id int64,
	name string,
	branchID int64,
	createdByID, updatedByID *int64,
) *Floor {
	return &Floor{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		Name:     name,
		BranchID: branchID,
	}
}
