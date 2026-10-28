package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const (
	DocTypeIdentityCard  = "Identity Card"
	DocTypePassport      = "Passport"
	DocTypeDriverLicense = "Driver's License"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
}

func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("doctype", ValidateDocumentType)
	}
}

func ValidateDocumentType(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	switch value {
	case DocTypeIdentityCard,
		DocTypePassport,
		DocTypeDriverLicense:
		return true
	}

	return false
}
