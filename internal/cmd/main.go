package main

import (
	"log"
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
		AllowOrigins:     []string{cfg.CorsOrigin},
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAccept, "X-CSRF-Token"},
	}))

	// API v1 routes
	api := e.Group("/api")
	apiV1 := api.Group("/v1")

	// JWT with Config - skip login and register
	apiV1.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    jwtSecret,
		SigningMethod: "HS256",
		TokenLookup:   "cookie:access_token",
		ContextKey:    "user",
		Skipper: func(c echo.Context) bool {
			// Skip JWT for login and register
			return c.Path() == "/api/v1/login" || c.Path() == "/api/v1/register"
		},
	}))
	// CSRF with Config - skip login, register, and logout
	apiV1.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookieName:     "csrf_token",
		CookiePath:     "/",
		CookieSameSite: cookie.SameSite(cfg.Env),
		CookieHTTPOnly: false, // Keep readable from JS
		CookieSecure:   cookie.Secure(cfg.Env),
		TokenLookup:    "header:X-CSRF-Token",
		Skipper: func(c echo.Context) bool {
			// Skip CSRF for login, register, and logout
			path := c.Path()
			return path == "/api/v1/login" || path == "/api/v1/register" || path == "/api/v1/logout"
		},
	}))

	// Public routes (JWT and CSRF skipped via middleware config)
	apiV1.POST("/login", auth.LoginHandler)
	apiV1.POST("/register", auth.RegisterHandler)

	// Protected routes (require JWT and CSRF)
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
