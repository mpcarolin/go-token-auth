package main

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"api/internal/constants/env"
	"api/internal/handlers"
	"api/internal/middleware"
	"api/internal/utils"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

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

	start(e);
}

// start echo server listening on port
func start (e *echo.Echo) {
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server!", "error", err)
	} else {
		slog.Info("server started on port 8080")
	}
}