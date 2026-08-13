package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

var (
	ErrNotFound         = &AppError{Code: "NOT_FOUND", HTTPStatus: http.StatusNotFound, Message: "Resource not found"}
	ErrUnauthorized     = &AppError{Code: "UNAUTHORIZED", HTTPStatus: http.StatusUnauthorized, Message: "Authentication required"}
	ErrForbidden        = &AppError{Code: "FORBIDDEN", HTTPStatus: http.StatusForbidden, Message: "Insufficient permissions"}
	ErrConflict         = &AppError{Code: "CONFLICT", HTTPStatus: http.StatusConflict, Message: "Resource already exists"}
	ErrRateLimited      = &AppError{Code: "RATE_LIMITED", HTTPStatus: http.StatusTooManyRequests, Message: "Too many requests"}
	ErrPaymentRequired  = &AppError{Code: "PAYMENT_REQUIRED", HTTPStatus: http.StatusPaymentRequired, Message: "Active subscription required"}
	ErrLimitExceeded    = &AppError{Code: "LIMIT_EXCEEDED", HTTPStatus: http.StatusUnprocessableEntity, Message: "Plan limit exceeded"}
	ErrInvalidInput     = &AppError{Code: "VALIDATION_ERROR", HTTPStatus: http.StatusUnprocessableEntity, Message: "Validation failed"}
	ErrSignatureInvalid = &AppError{Code: "INVALID_SIGNATURE", HTTPStatus: http.StatusUnauthorized, Message: "Invalid signature"}
)

func NotFound(entity string) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		HTTPStatus: http.StatusNotFound,
		Message:    fmt.Sprintf("%s not found", entity),
	}
}

func NotOwnedByUser(entity string) *AppError {
	return &AppError{
		Code:       "FORBIDDEN",
		HTTPStatus: http.StatusForbidden,
		Message:    fmt.Sprintf("You do not own this %s", entity),
	}
}

func ValidationFailed(details map[string]string) *AppError {
	msg := "Validation failed"
	if len(details) > 0 {
		msg += " - see details"
	}
	return &AppError{
		Code:       "VALIDATION_ERROR",
		HTTPStatus: http.StatusUnprocessableEntity,
		Message:    msg,
	}
}
