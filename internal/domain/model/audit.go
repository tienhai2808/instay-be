package model

type Audit struct {
	CreatedByID *int64
	UpdatedByID *int64

	CreatedBy *User
	UpdatedBy *User
}
