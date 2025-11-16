package response

import (
	"github.com/labstack/echo/v4"
)

// Error code constants
// Authentication errors
const (
	CodeAuthInvalidCredentials = "AUTH_INVALID_CREDENTIALS"
	CodeAuthUserExists         = "AUTH_USER_EXISTS"
	CodeAuthUserNotFound       = "AUTH_USER_NOT_FOUND"
	CodeAuthMissingToken       = "AUTH_MISSING_TOKEN"
	CodeAuthInvalidTokenType   = "AUTH_INVALID_TOKEN_TYPE"
	CodeAuthInvalidTokenClaims = "AUTH_INVALID_TOKEN_CLAIMS"
	CodeAuthMissingUserID      = "AUTH_MISSING_USER_ID"
	CodeAuthInvalidUserIDType  = "AUTH_INVALID_USER_ID_TYPE"
)

// Address errors
const (
	CodeAddressForbidden           = "ADDRESS_FORBIDDEN"
	CodeAddressNotFound            = "ADDRESS_NOT_FOUND"
	CodeAddressCannotDeleteDefault = "ADDRESS_CANNOT_DELETE_DEFAULT"
)

// Validation and request errors
const (
	CodeInvalidPayload   = "INVALID_PAYLOAD"
	CodeValidationError  = "VALIDATION_ERROR"
	CodeInvalidAddressID = "INVALID_ADDRESS_ID"
)

// Generic errors
const (
	CodeInternalServerError = "INTERNAL_SERVER_ERROR"
	CodeNotFound            = "NOT_FOUND"
)

// ErrorResponse represents the standardized error response structure
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains the error code and message
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error returns a standardized error response in Echo format
// It creates an HTTPError with the standardized error structure
func Error(c echo.Context, statusCode int, code string, message string) error {
	return c.JSON(statusCode, ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}
