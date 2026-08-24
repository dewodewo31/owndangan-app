package response

import (
	"encoding/json"
	stderrors "errors"
	"net/http"

	apperrors "github.com/owndangan/backend/internal/pkg/errors"
)

type JSONResponse struct {
	w http.ResponseWriter
	r *http.Request
}

func NewJSONResponse(w http.ResponseWriter) *JSONResponse {
	return &JSONResponse{w: w}
}

func (j *JSONResponse) SetRequest(r *http.Request) {
	j.r = r
}

func JSON(w http.ResponseWriter, status int, data interface{}, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r != nil && r.Header.Get("X-Request-ID") != "" {
		w.Header().Set("X-Request-ID", getRequestID(r))
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
		"meta":    buildMeta(r),
	})
}

func JSONPaginated(w http.ResponseWriter, status int, data interface{}, page, perPage, total int, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r != nil {
		w.Header().Set("X-Request-ID", getRequestID(r))
	}
	w.WriteHeader(status)
	totalPages := 0
	if perPage > 0 {
		totalPages = (total + perPage - 1) / perPage
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
		"meta": map[string]interface{}{
			"pagination": map[string]interface{}{
				"page":        page,
				"per_page":    perPage,
				"total":       total,
				"total_pages": totalPages,
			},
			"request_id": getRequestID(r),
		},
	})
}

func Error(w http.ResponseWriter, status int, code, message string, r *http.Request) {
	writeError(w, status, code, message, nil, r)
}

func ValidationError(w http.ResponseWriter, validationErr error, r *http.Request) {
	details := formatValidationErrors(validationErr)
	writeError(w, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Validation failed", details, r)
}

func FromError(w http.ResponseWriter, err error, r *http.Request) {
	var appErr *apperrors.AppError
	if stderrors.As(err, &appErr) {
		writeError(w, appErr.HTTPStatus, appErr.Code, appErr.Message, nil, r)
		return
	}
	writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An unexpected error occurred", nil, r)
}

type errorResponse struct {
	Success bool                   `json:"success"`
	Error   apperrors.ErrorDetail  `json:"error"`
	Meta    map[string]interface{} `json:"meta"`
}

func writeError(w http.ResponseWriter, status int, code, message string, details map[string]string, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r != nil {
		w.Header().Set("X-Request-ID", getRequestID(r))
	}
	w.WriteHeader(status)
	payload := errorResponse{
		Success: false,
		Error: apperrors.ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: buildMeta(r),
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func getRequestID(r *http.Request) string {
	if r == nil {
		return ""
	}
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	return ""
}

func buildMeta(r *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"request_id": getRequestID(r),
	}
}

func formatValidationErrors(err error) map[string]string {
	details := make(map[string]string)
	if err == nil {
		return details
	}
	type FieldError interface {
		Field() string
		Tag() string
	}
	if errs, ok := err.(interface{ Errors() []FieldError }); ok {
		for _, fe := range errs.Errors() {
			details[fe.Field()] = validationMessage(fe.Tag(), fe.Field(), err)
		}
	}
	return details
}

func validationMessage(tag, field string, err error) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return "Must be at least the required length"
	case "max":
		return "Must be at most the required length"
	case "uuid":
		return "Must be a valid UUID"
	default:
		return "Invalid value"
	}
}
