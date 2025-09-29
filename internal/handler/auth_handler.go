package handler

import (
	"errors"
	"net/http"
	"server/internal/model"
	"server/internal/service"
	"strings"

	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)


type AuthHandler struct {
	authSvc *service.AuthService
	jwtSecret []byte
}

func NewAuthHandler(authSvc *service.AuthService, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		authSvc: authSvc,
		jwtSecret: jwtSecret,
	}
}

// user for sanitation
type user struct {
	Username         string `json:"username" validate:"required,min=3,max=30"`
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required,min=8"`
	RepeatedPassword string `json:"repeatedPassword" validate:"required,eqfield=Password"`
}

// Normalize implements Normalizable (from custom validator)
func (u *user) Normalize() {
	u.Username = strings.TrimSpace(u.Username)
  u.Email    = strings.ToLower(strings.TrimSpace(u.Email))
}

// RegisterHandler
func (h *AuthHandler) RegisterHandler(c echo.Context) error {
	req := new(user)
	bindErr := c.Bind(req)
  if bindErr != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
  }
	
	// Form values sanitation
	validateErr := c.Validate(req)
	if validateErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validateErr.Error())
	}
	
	// wire Auth service
	user, registerErr := h.authSvc.Register(req.Username, req.Email, req.Password)
	if	registerErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	} else {
		return c.JSON(http.StatusCreated, echo.Map{
			"user": echo.Map{
				"username": user.Username,
			},
		})
	}
}

// loginRequest for sanitation
type loginUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Normalize implements Normalizable
func (r *loginUser) Normalize() {
	r.Email    = strings.ToLower(strings.TrimSpace(r.Email))
}

// LoginHandler
func (h *AuthHandler) LoginHandler(c echo.Context) error {
	req := new(loginUser)
	bindErr := c.Bind(req)
  if bindErr != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
  }

	validateErr := c.Validate(req)
	if validateErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validateErr.Error())
	}

	user, loginErr := h.authSvc.Login(req.Email, req.Password)
	if loginErr != nil {
		if errors.Is(loginErr, service.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, "login failed")	
		}
	} else {
		tokenString, err := h.issueToken(user)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "token generation failed")
		}

		// set HttpOnly cookie
		h.setTokenCookie(c, tokenString)

		// return basic user info
		return c.JSON(http.StatusOK, echo.Map{"user": echo.Map{"username": user.Username}})
	}	
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
    Secure:   true,                         
    SameSite: http.SameSiteStrictMode,
  }
  c.SetCookie(accessTokenCookie)  
	
	// Expire the CSRF cookie
  csrfCookie := &http.Cookie{
    Name:     "csrf_token",
    Value:    "",
    Path:     "/",
    Expires:  time.Unix(0, 0),
    MaxAge:   -1,
    HttpOnly: true,                          
    Secure:   false,
    SameSite: http.SameSiteStrictMode,
  }
  c.SetCookie(csrfCookie)
	return c.NoContent(http.StatusNoContent)
}

// issueToken creates a signed JWT string
func (h *AuthHandler) issueToken(u *model.User) (string, error){
	exp := time.Now().Add(24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"user_id": u.ID,
    "email":   u.Email,
		"exp": exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwtSecret)
}

// setTokenCookie writes the JWT into an HttpOnly cookie
func (h *AuthHandler) setTokenCookie(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name: "access_token",
		Value: token,
		Path: "/",
		Expires: time.Now().Add(24*time.Hour),
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
}