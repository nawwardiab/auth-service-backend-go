package main

import (
	"log"
	"net/http"
	"server/internal/config"
	"server/internal/cookie"
	"server/internal/db"
	"server/internal/handler"
	"server/internal/repo"
	"server/internal/service"
	"server/internal/validator"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load Config
	cfg, cfgErr := config.LoadConfig()
	if cfgErr != nil {
		log.Fatal("failed to load config: %w", cfgErr)
	}

	// connect to db
	dbConn, dbConnErr := db.NewDB(cfg)
	if dbConnErr != nil {
		log.Fatal("failed to connect to db: %w", dbConnErr)
	}
	defer dbConn.Close()

	jwtSecret := []byte(cfg.JwtSecret)

	// Wire repos and services
	// Auth
	authRepo := repo.NewAuthRepo(dbConn)
	authSvc := service.NewAuthService(authRepo)
	auth := handler.NewAuthHandler(authSvc, jwtSecret, cfg.Env)

	// Address
	addrRepo := repo.NewAddressRepo(dbConn)
	addrSvc := service.NewAddressService(addrRepo)
	addr := handler.NewAddressHandler(addrSvc)

	// instantiate echo
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))

	// Wire up echo validator
	e.Validator = validator.New()

	// CORS for React app
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAccept, "X-CSRF-Token"},
	}))

	// group for protected API routes (JWT + CSRF)

	// public routes
	api := e.Group("/api")
	api.POST("/login", auth.LoginHandler)
	api.POST("/register", auth.RegisterHandler)

	apiV1 := api.Group("/v1")
	// JWT with Config
	apiV1.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    jwtSecret,
		SigningMethod: "HS256",
		TokenLookup:   "cookie:access_token",
		ContextKey:    "user",
	}))
	// CSRF with Config
	apiV1.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookieName:     "csrf_token",
		CookieSameSite: cookie.SameSite(cfg.Env),
		CookieHTTPOnly: false, // Keep readable from JS
		CookieSecure:   cookie.Secure(cfg.Env),
		TokenLookup:    "header:X-CSRF-Token",
		// skip CSRF on logout
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/api/v1/logout"
		},
	}))

	// Middleware to ensure CSRF token cookie is accessible at root path
	// This allows JavaScript to read it via document.cookie
	apiV1.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get CSRF token from cookie (set by CSRF middleware)
			csrfCookie, err := c.Cookie("csrf_token")
			if err == nil && csrfCookie != nil && csrfCookie.Value != "" {
				// Set the same cookie at root path so it's accessible via document.cookie
				rootCookie := &http.Cookie{
					Name:     "csrf_token",
					Value:    csrfCookie.Value,
					Path:     "/",
					MaxAge:   86400, // 24 hours
					HttpOnly: false, // Keep readable from JS
					Secure:   cookie.Secure(cfg.Env),
					SameSite: cookie.SameSite(cfg.Env),
				}
				c.SetCookie(rootCookie)
			}
			return next(c)
		}
	})

	// Wire portected routes
	apiV1.POST("/logout", auth.LogoutHandler)
	apiV1.GET("/profile", auth.ProfileHandler)

	apiV1.GET("/users/addresses", addr.GetUserAddresses)
	apiV1.POST("/users/address/add", addr.CreateAddress)
	apiV1.GET("/users/address/:id", addr.GetAddress)
	apiV1.PATCH("/users/address/:id", addr.UpdateAddress)
	apiV1.DELETE("/users/address/:id", addr.DeleteAddress)

	serverPort := cfg.ServerPort
	serverHost := cfg.ServerHost
	addrStr := serverHost + ":" + serverPort
	e.Logger.Fatal(e.Start(addrStr))
}
