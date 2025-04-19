package server

import (
	"api/internal/constants/env"
	"api/internal/handlers"
	"api/internal/middleware"
	"api/internal/utils"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func SetupServer(e *echo.Echo) {
	// Set debug mode based on environment
	e.Debug = utils.GetEnv() != env.Production

	// Middleware
	e.Use(middleware.SetupCtx)
	e.Use(middleware.Cleanup)
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(20)))
	e.Use(echoMiddleware.CORSWithConfig(utils.GetCORSConfig()))

	// Routes
	e.GET("/status", func(c echo.Context) error {
		response := "OK " + time.Now().Format("2006-01-02 15:04:05")
		return c.String(http.StatusOK, response)
	})

	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)

	// using this for testing, in real app probably would not expose this data as an endpoint
	e.GET("/users", middleware.Authenticate(handlers.GetUsers))

}

// start echo server listening on port
func Start (e *echo.Echo) {
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server!", "error", err)
	} else {
		slog.Info("server started on port 8080")
	}
}