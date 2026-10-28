package validator

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	govalidator "github.com/go-playground/validator/v10"
)

func TestHandleRequestError_ValidationErrors(t *testing.T) {
	type req struct {
		Email string `validate:"required,email"`
	}
	v := govalidator.New()
	err := v.Struct(req{})
	field, tag, param := HandleRequestError(err)
	if field != "email" {
		t.Fatalf("expected field email, got %q", field)
	}
	if tag != "required" {
		t.Fatalf("expected tag required, got %q", tag)
	}
	if param != "" {
		t.Fatalf("expected empty param, got %q", param)
	}
}

func TestHandleRequestError_TypeError(t *testing.T) {
	type req struct {
		Age int `json:"age"`
	}
	var r req
	err := json.Unmarshal([]byte(`{"age":"x"}`), &r)
	field, tag, param := HandleRequestError(err)
	if strings.ToLower(field) != "age" {
		t.Fatalf("expected field age, got %q", field)
	}
	if tag != "type_error" {
		t.Fatalf("expected tag type_error, got %q", tag)
	}
	if param == "" {
		t.Fatalf("expected param type, got empty")
	}
}

func TestHandleRequestError_TypeErrorSnakeCaseField(t *testing.T) {
	type req struct {
		BookingRoomID int64 `json:"booking_room_id"`
	}

	var r req
	err := json.Unmarshal([]byte(`{"booking_room_id":"x"}`), &r)
	field, tag, _ := HandleRequestError(err)

	if field != "bookingroomid" {
		t.Fatalf("expected field bookingroomid, got %q", field)
	}
	if tag != "type_error" {
		t.Fatalf("expected tag type_error, got %q", tag)
	}
}

func TestHandleRequestError_SyntaxError(t *testing.T) {
	var v any
	err := json.Unmarshal([]byte(`{"age":`), &v)
	field, tag, _ := HandleRequestError(err)
	if field != "json" || tag != "syntax_error" {
		t.Fatalf("expected json/syntax_error, got %q/%q", field, tag)
	}
}

func TestHandleRequestError_Invalid(t *testing.T) {
	field, tag, param := HandleRequestError(errors.New("boom"))
	if field != "unknown" || tag != "invalid" || param != "" {
		t.Fatalf("expected unknown/invalid/empty, got %q/%q/%q", field, tag, param)
	}
}
