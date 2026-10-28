package model

const (
	DepartmentCustomerService = "Chăm sóc khách hàng (Customer service)"
)

type Department struct {
	Base
	Audit

	Name        string
	Phone       string
	Description string

	Permissions   []*Permission
	Users         []*User
	ServiceTypes  []*ServiceType
	Notifications []*Notification
	UsersCount    int64
}

type DepartmentQuery struct {
	Page   uint32
	Limit  uint32
	Search string
}

func NewDepartment(
	id int64,
	name, phone, description string,
	createdByID, updatedByID *int64,
	permissions []*Permission,
) *Department {
	return &Department{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		Name:        name,
		Phone:       phone,
		Description: description,
		Permissions: permissions,
	}
}
