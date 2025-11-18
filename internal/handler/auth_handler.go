package handler

import (
	"net/http"
	"server/internal/cookie"
	"server/internal/model"
	"server/internal/response"
	"server/internal/service"
	"strings"

	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authSvc   *service.AuthService
	jwtSecret []byte
	env       string
}

func NewAuthHandler(authSvc *service.AuthService, jwtSecret []byte, env string) *AuthHandler {
	return &AuthHandler{
		authSvc:   authSvc,
		jwtSecret: jwtSecret,
		env:       env,
	}
}

// registerUser for sanitation
type registerUser struct {
	Username         string `json:"username" validate:"required,min=3,max=30"`
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required,min=8"`
	RepeatedPassword string `json:"repeatedPassword" validate:"required,eqfield=Password"`
}

// Normalize implements Normalizable (from custom validator)
func (u *registerUser) Normalize() {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)
	u.RepeatedPassword = strings.TrimSpace(u.RepeatedPassword)
}

// RegisterHandler
func (h *AuthHandler) RegisterHandler(c echo.Context) error {
	req := new(registerUser)
	bindErr := c.Bind(req)
	if bindErr != nil {
		return response.Error(c, http.StatusBadRequest, response.CodeInvalidPayload, "invalid request payload")
	}

	// Form values sanitation
	validateErr := c.Validate(req)
	if validateErr != nil {
		return response.Error(c, http.StatusBadRequest, response.CodeValidationError, validateErr.Error())
	}

	// wire Auth service
	user, registerErr := h.authSvc.Register(req.Username, req.Email, req.Password)
	if registerErr != nil {
		code, status, msg := MapServiceError(registerErr)
		return response.Error(c, status, code, msg)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"user": echo.Map{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// loginUser for sanitation
type loginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// Normalize implements Normalizable (from custom validator)
func (u *loginUser) Normalize() {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)
}

// LoginHandler
func (h *AuthHandler) LoginHandler(c echo.Context) error {
	req := new(loginUser)
	bindErr := c.Bind(req)
	if bindErr != nil {
		return response.Error(c, http.StatusBadRequest, response.CodeInvalidPayload, "invalid request payload")
	}

	validateErr := c.Validate(req)
	if validateErr != nil {
		return response.Error(c, http.StatusBadRequest, response.CodeValidationError, validateErr.Error())
	}

	user, loginErr := h.authSvc.Login(req.Email, req.Password)
	if loginErr != nil {
		code, status, msg := MapServiceError(loginErr)
		return response.Error(c, status, code, msg)
	}

	tokenString, err := h.issueToken(user)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, response.CodeInternalServerError, "token generation failed")
	}

	// set HttpOnly cookie
	h.setTokenCookie(c, tokenString)

	// return basic user info
	return c.JSON(http.StatusOK, echo.Map{"user": echo.Map{"username": user.Username}})
}

// LogoutHandler
func (h *AuthHandler) LogoutHandler(c echo.Context) error {
	// Expire the JWT cookie
	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cookie.Secure(h.env),
		SameSite: cookie.SameSite(h.env),
	}
	c.SetCookie(accessTokenCookie)

	// Expire the CSRF cookie (HttpOnly: false so JS can read it, but secure otherwise)
	csrfCookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: false,
		Secure:   cookie.Secure(h.env),
		SameSite: cookie.SameSite(h.env),
	}
	c.SetCookie(csrfCookie)
	return c.NoContent(http.StatusNoContent)
}

// ProfileHandler returns the current user's profile
func (h *AuthHandler) ProfileHandler(c echo.Context) error {
	// Extract user ID from JWT token
	userID, err := extractUserIDFromJWT(c)
	if err != nil {
		return err
	}

	// Get user from service
	user, err := h.authSvc.GetUser(userID)
	if err != nil {
		code, status, msg := MapServiceError(err)
		return response.Error(c, status, code, msg)
	}

	// Return user info (without password hash)
	return c.JSON(http.StatusOK, echo.Map{
		"user": echo.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Helpers:

// extractUserIDFromJWT extracts the user_id from the JWT token in the echo context.
// Returns the user ID and an error (which can be returned directly from echo handlers).
func extractUserIDFromJWT(c echo.Context) (int, error) {
	// Get user from context with nil check
	userValue := c.Get("user")
	if userValue == nil {
		return 0, response.Error(c, http.StatusUnauthorized, response.CodeAuthMissingToken, "missing authentication token")
	}

	// Type check for *jwt.Token
	token, ok := userValue.(*jwt.Token)
	if !ok {
		return 0, response.Error(c, http.StatusForbidden, response.CodeAuthInvalidTokenType, "invalid token type")
	}

	// Type check for jwt.MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, response.Error(c, http.StatusForbidden, response.CodeAuthInvalidTokenClaims, "invalid token claims")
	}

	// Extract user_id with type assertion check
	userIDValue, exists := claims["user_id"]
	if !exists {
		return 0, response.Error(c, http.StatusForbidden, response.CodeAuthMissingUserID, "missing user_id in token")
	}

	userIDFloat, ok := userIDValue.(float64)
	if !ok {
		return 0, response.Error(c, http.StatusForbidden, response.CodeAuthInvalidUserIDType, "invalid user_id type in token")
	}

	return int(userIDFloat), nil
}

// issueToken creates a signed JWT string
func (h *AuthHandler) issueToken(u *model.User) (string, error) {
	exp := time.Now().Add(24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"email":   u.Email,
		"exp":     exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwtSecret)
}

// setTokenCookie writes the JWT into an HttpOnly cookie
func (h *AuthHandler) setTokenCookie(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   cookie.Secure(h.env),
		SameSite: cookie.SameSite(h.env),
	}
	c.SetCookie(cookie)
}
