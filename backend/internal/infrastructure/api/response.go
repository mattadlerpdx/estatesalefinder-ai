package api

import (
	"encoding/json"
	"net/http"
)

// Response is the standard API response wrapper
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

// SuccessResponse sends a successful JSON response
func SuccessResponse(w http.ResponseWriter, data interface{}, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: true,
		Data:    data,
		Message: message,
	}

	json.NewEncoder(w).Encode(response)
}

// ErrorResponse sends an error JSON response
func ErrorResponse(w http.ResponseWriter, message string, errors []string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: false,
		Message: message,
		Errors:  errors,
	}

	json.NewEncoder(w).Encode(response)
}

// ErrorResponseSingle sends an error response with a single error message
func ErrorResponseSingle(w http.ResponseWriter, message string, statusCode int) {
	ErrorResponse(w, message, nil, statusCode)
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(w http.ResponseWriter, errors []string) {
	ErrorResponse(w, "Validation failed", errors, http.StatusBadRequest)
}

// NotFoundResponse sends a 404 not found response
func NotFoundResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Resource not found"
	}
	ErrorResponseSingle(w, message, http.StatusNotFound)
}

// UnauthorizedResponse sends a 401 unauthorized response
func UnauthorizedResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	ErrorResponseSingle(w, message, http.StatusUnauthorized)
}

// ForbiddenResponse sends a 403 forbidden response
func ForbiddenResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Forbidden"
	}
	ErrorResponseSingle(w, message, http.StatusForbidden)
}

// InternalErrorResponse sends a 500 internal server error response
func InternalErrorResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = "Internal server error"
	}
	ErrorResponseSingle(w, message, http.StatusInternalServerError)
}

// CreatedResponse sends a 201 created response
func CreatedResponse(w http.ResponseWriter, data interface{}, message string) {
	if message == "" {
		message = "Resource created successfully"
	}
	SuccessResponse(w, data, message, http.StatusCreated)
}

// OKResponse sends a 200 OK response
func OKResponse(w http.ResponseWriter, data interface{}, message string) {
	SuccessResponse(w, data, message, http.StatusOK)
}

// NoContentResponse sends a 204 No Content response
func NoContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
