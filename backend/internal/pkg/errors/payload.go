package errors

type ErrorDetail struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type ErrorPayload struct {
	Success bool        `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type Meta struct {
	RequestID  string                 `json:"request_id"`
	Pagination *PaginationMeta        `json:"pagination,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

type PaginationMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func NewErrorPayload(code, message string, details map[string]string) ErrorPayload {
	return ErrorPayload{
		Success: false,
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}
