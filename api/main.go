package main

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"api/internal/constants/env"
	"api/internal/handlers"
	"api/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)


func main() {
	// Echo instance
	e := echo.New()

	// Set debug mode based on environment
	e.Debug = utils.GetEnv() == env.Development

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/status", func(c echo.Context) error {
		response := "OK "+time.Now().Format("2006-01-02 15:04:05")
		return c.String(http.StatusOK, response)
	})

	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)
	e.GET("/user/:username", handlers.GetUser)

	// Start server
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server!", "error", err)
	} else {
		slog.Info("server started on port 8080")
	}
}