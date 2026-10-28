package model

type SupportCategory struct {
	Base
	Audit

	Name         string
	DepartmentID int64

	Department          *Department
	SupportTickets      []*SupportTicket
	SupportTicketsCount int64
}

type SupportCategoryQuery struct {
	Page         uint32
	Limit        uint32
	Search       string
	DepartmentID int64
}

func NewSupportCategory(
	id, deptID int64,
	name string,
	createdByID, updatedByID *int64,
) *SupportCategory {
	return &SupportCategory{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		Name:         name,
		DepartmentID: deptID,
	}
}
