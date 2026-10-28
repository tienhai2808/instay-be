package model

type Branch struct {
	Base
	Audit

	Name     string
	Address  string
	Ward     string
	City     string
	IsActive bool

	Floors    []*Floor
	RoomTypes []*RoomType
	Services  []*Service
	Events    []*Event
	// Users []*User
}

type BranchQuery struct {
	Page     uint32
	Limit    uint32
	Search   string
	City     string
	IsActive *bool
}

func NewBranch(
	id int64,
	name, address, ward, city string,
	isActive bool,
	createdByID, updatedByID *int64,
) *Branch {
	return &Branch{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		Name:     name,
		Address:  address,
		Ward:     ward,
		City:     city,
		IsActive: isActive,
	}
}
