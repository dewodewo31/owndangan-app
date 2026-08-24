package validator

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/owndangan/backend/internal/pkg/response"
)

func New() *validator.Validate {
	v := validator.New()
	registerCustomValidators(v)
	return v
}

func registerCustomValidators(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" {
			return fld.Name
		}
		return strings.Split(name, ",")[0]
	})
	_ = v.RegisterValidation("slug", slugValidator)
	_ = v.RegisterValidation("phone_id", phoneIDValidator)
	_ = v.RegisterValidation("date", dateValidator)
}

func ParseAndValidate(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return err
	}
	return ValidateStruct(dst)
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		details := FormatValidationError(verr)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error": map[string]interface{}{
				"code":    "VALIDATION_ERROR",
				"message": "Validation failed",
				"details": details,
			},
		})
		return
	}
	response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error(), r)
}

func ValidateStruct(s interface{}) error {
	v := New()
	return v.Struct(s)
}

func FormatValidationError(err error) map[string]string {
	details := make(map[string]string)
	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		for _, fe := range verr {
			details[fe.Field()] = validationMessage(fe.Tag(), fe.Param())
		}
	}
	return details
}

var (
	slugPattern  = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	phonePattern = regexp.MustCompile(`^62\d{8,15}$`)
)

func slugValidator(fl validator.FieldLevel) bool {
	return slugPattern.MatchString(fl.Field().String())
}

func phoneIDValidator(fl validator.FieldLevel) bool {
	return phonePattern.MatchString(fl.Field().String())
}

func dateValidator(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}

func validationMessage(tag, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least " + param + " characters"
	case "max":
		return "Must be at most " + param + " characters"
	case "uuid":
		return "Must be a valid UUID"
	case "slug":
		return "Must be a valid slug (lowercase letters, numbers, hyphens)"
	case "phone_id":
		return "Must be a valid Indonesian phone number starting with 62"
	case "oneof":
		return "Must be one of: " + param
	case "date":
		return "Must be a valid date in YYYY-MM-DD format"
	case "omitempty":
		return ""
	default:
		return "Invalid value"
	}
}
