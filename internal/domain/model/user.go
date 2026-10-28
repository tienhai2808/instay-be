package model

type UserRole string

const (
	RoleStaff UserRole = "staff"
	RoleAdmin UserRole = "admin"
)

type User struct {
	Base
	Audit

	Username     string
	Email        string
	Role         UserRole
	FirstName    string
	LastName     string
	Phone        string
	Password     string
	Avatar       *string
	IsActive     bool
	DepartmentID *int64

	Department    *Department
	RefreshTokens []*RefreshToken
	Notifications []*Notification
	Branches      []*Branch
}

type UserQuery struct {
	Page         uint32
	Limit        uint32
	Sort         string
	Order        string
	Role         string
	DepartmentID int64
	BranchID     int64
	IsActive     *bool
	Search       string
}

func NewUser(
	id int64,
	username, email string,
	role UserRole,
	firstName, lastName, phone, password string,
	isActive bool,
	avatar *string,
	deptID, createdByID, updatedByID *int64,
) *User {
	return &User{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		Username:     username,
		Email:        email,
		Role:         role,
		FirstName:    firstName,
		LastName:     lastName,
		Phone:        phone,
		Password:     password,
		IsActive:     isActive,
		DepartmentID: deptID,
		Avatar:       avatar,
	}
}

func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleStaff:
		return true
	default:
		return false
	}
}
