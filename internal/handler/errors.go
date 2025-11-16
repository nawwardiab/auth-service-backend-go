package handler

import (
	"errors"
	"net/http"
	"server/internal/response"
	"server/internal/service"
)

// MapServiceError maps domain errors from the service layer to standardized error codes,
// HTTP status codes, and user-friendly messages.
// Returns (errorCode, httpStatusCode, message).
func MapServiceError(err error) (code string, statusCode int, message string) {
	if err == nil {
		return response.CodeInternalServerError, http.StatusInternalServerError, "internal server error"
	}

	// Auth service errors
	if errors.Is(err, service.ErrInvalidCredentials) {
		return response.CodeAuthInvalidCredentials, http.StatusUnauthorized, "invalid credentials"
	}
	if errors.Is(err, service.ErrUserExist) {
		return response.CodeAuthUserExists, http.StatusConflict, "user already exists"
	}
	if errors.Is(err, service.ErrUserNotFound) {
		return response.CodeAuthUserNotFound, http.StatusNotFound, "user not found"
	}

	// Address service errors
	if errors.Is(err, service.ErrForbidden) {
		return response.CodeAddressForbidden, http.StatusForbidden, "not allowed to access this resource"
	}
	if errors.Is(err, service.ErrAddressNotFound) {
		return response.CodeAddressNotFound, http.StatusNotFound, "address not found"
	}
	if errors.Is(err, service.ErrCannotDeleteDefault) {
		return response.CodeAddressCannotDeleteDefault, http.StatusBadRequest, "cannot delete default address"
	}

	// Default for unknown errors
	return response.CodeInternalServerError, http.StatusInternalServerError, "internal server error"
}
