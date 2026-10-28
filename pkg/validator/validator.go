package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func HandleRequestError(err error) (string, string, string) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		e := errs[0]
		return normalizeBindingField(e.Field()), e.Tag(), e.Param()
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return normalizeBindingField(unmarshalTypeError.Field), "type_error", unmarshalTypeError.Type.String()
	}

	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return "json", "syntax_error", fmt.Sprintf("%d", syntaxError.Offset)
	}

	return "unknown", "invalid", ""
}

func normalizeBindingField(field string) string {
	field = strings.TrimSpace(field)
	if field == "" {
		return "unknown"
	}

	parts := strings.Split(field, ".")
	field = parts[len(parts)-1]

	var normalized strings.Builder
	normalized.Grow(len(field))

	for _, r := range field {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') {
			normalized.WriteRune(r)
		}
	}

	value := strings.ToLower(normalized.String())
	if value == "" {
		return "unknown"
	}

	return value
}
