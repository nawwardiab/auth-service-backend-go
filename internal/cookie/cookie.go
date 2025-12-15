package cookie

import (
	"net/http"

	"server/internal/config"
)

// Secure returns the Secure flag value based on environment
// Secure: true in production (HTTPS required), false in development
func Secure(env string) bool {
	return config.IsProduction(env)
}

// SameSite returns the SameSite flag value based on environment
// SameSiteStrictMode in production, SameSiteLaxMode in development
func SameSite(env string) http.SameSite {
	if config.IsProduction(env) {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}
