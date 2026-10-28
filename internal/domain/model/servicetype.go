package model

type ServiceType struct {
	Base
	Audit

	Name         string
	Icon         string
	Slug         string
	DepartmentID int64

	Department    *Department
	Services      []*Service
	ServicesCount int64
}

type ServiceTypeQuery struct {
	Page         uint32
	Limit        uint32
	Search       string
	DepartmentID int64
}

type ServiceTypeWithServicesQuery struct {
	Search        string
	ServiceTypeID int64
	BranchID      int64
}

func NewServiceType(
	id, deptID int64,
	name, icon, slug string,
	createdByID, updatedByID *int64,
) *ServiceType {
	return &ServiceType{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		Name:         name,
		Icon:         icon,
		Slug:         slug,
		DepartmentID: deptID,
	}
}
