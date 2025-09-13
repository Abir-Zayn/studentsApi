package response

import (
	"fmt"
	"net/http"
)

// Client error responses (4xx)

func BadRequest(message string) Response {
	return Response{
		Status:     StatusFail,
		Message:    message,
		Error:      "Bad Request",
		StatusCode: http.StatusBadRequest,
	}
}

func ValidationError(message string, details interface{}) Response {
	return Response{
		Status:     StatusFail,
		Message:    message,
		Data:       details,
		Error:      "Validation Failed",
		StatusCode: http.StatusBadRequest,
	}
}

func Unauthorized(message string) Response {
	return Response{
		Status:     StatusFail,
		Message:    message,
		Error:      "Unauthorized",
		StatusCode: http.StatusUnauthorized,
	}
}

func Forbidden(message string) Response {
	return Response{
		Status:     StatusFail,
		Message:    message,
		Error:      "Forbidden",
		StatusCode: http.StatusForbidden,
	}
}

func NotFound(message string) Response {
	return Response{
		Status:     StatusFail,
		Message:    message,
		Error:      "Resource Not Found",
		StatusCode: http.StatusNotFound,
	}
}

func Conflict(message string) Response {
	return Response{
		Status:     StatusFail,
		Message:    message,
		Error:      "Conflict",
		StatusCode: http.StatusConflict,
	}
}

func UnprocessableEntity(message string, details interface{}) Response {
	return Response{
		Status:     StatusFail,
		Message:    message,
		Data:       details,
		Error:      "Unprocessable Entity",
		StatusCode: http.StatusUnprocessableEntity,
	}
}

// Server error responses (5xx)
func InternalServerError(message string) Response {
	return Response{
		Status:     StatusError,
		Message:    message,
		Error:      "Internal Server Error",
		StatusCode: http.StatusInternalServerError,
	}
}

func ServiceUnavailable(message string) Response {
	return Response{
		Status:     StatusError,
		Message:    message,
		Error:      "Service Unavailable",
		StatusCode: http.StatusServiceUnavailable,
	}
}

// Utility functions for common error scenarios
func GeneralError(err error) Response {
	return InternalServerError(fmt.Sprintf("An unexpected error occurred: %s", err.Error()))
}

func JSONParseError(err error) Response {
	return BadRequest(fmt.Sprintf("Invalid JSON format: %s", err.Error()))
}

func MissingFieldsError(fields []string) Response {
	message := "Missing required fields"
	if len(fields) > 0 {
		message = fmt.Sprintf("Missing required fields: %v", fields)
	}
	return ValidationError(message, map[string][]string{"missing_fields": fields})
}

func InvalidFieldError(field, reason string) Response {
	message := fmt.Sprintf("Invalid field '%s': %s", field, reason)
	return ValidationError(message, map[string]string{"invalid_field": field, "reason": reason})
}

// Quick send methods for error responses
func SendBadRequest(w http.ResponseWriter, message string) error {
	return Send(w, BadRequest(message))
}

func SendValidationError(w http.ResponseWriter, message string, details interface{}) error {
	return Send(w, ValidationError(message, details))
}

func SendUnauthorized(w http.ResponseWriter, message string) error {
	return Send(w, Unauthorized(message))
}

func SendForbidden(w http.ResponseWriter, message string) error {
	return Send(w, Forbidden(message))
}

func SendNotFound(w http.ResponseWriter, message string) error {
	return Send(w, NotFound(message))
}

func SendConflict(w http.ResponseWriter, message string) error {
	return Send(w, Conflict(message))
}

func SendInternalServerError(w http.ResponseWriter, message string) error {
	return Send(w, InternalServerError(message))
}

func SendGeneralError(w http.ResponseWriter, err error) error {
	return Send(w, GeneralError(err))
}

func SendJSONParseError(w http.ResponseWriter, err error) error {
	return Send(w, JSONParseError(err))
}

func SendMissingFieldsError(w http.ResponseWriter, fields []string) error {
	return Send(w, MissingFieldsError(fields))
}

func SendInvalidFieldError(w http.ResponseWriter, field, reason string) error {
	return Send(w, InvalidFieldError(field, reason))
}
