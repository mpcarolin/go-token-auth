package utils

import (
	"api/internal/constants/env"
	"os"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

// GetCORSConfig returns a CORS configuration based on the environment
func GetCORSConfig() echoMiddleware.CORSConfig {
	// In development, allow localhost:3000
	// In production, allow only the specified domain
	allowedOrigin := "http://localhost:3000"
	if GetEnv() == env.Production {
		allowedOrigin = os.Getenv("ALLOWED_ORIGIN")
		if allowedOrigin == "" {
			allowedOrigin = "https://your-production-domain.com"
		}
	}

	return echoMiddleware.CORSConfig{
		AllowOrigins:     []string{allowedOrigin},
		AllowMethods:     []string{echo.GET, echo.POST, echo.OPTIONS},
		AllowHeaders:     []string{"Content-Type", "Cookie"},
		AllowCredentials: true,
		MaxAge:           3600,
	}
} 