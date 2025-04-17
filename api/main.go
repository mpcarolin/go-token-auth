package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"api/internal/constants/env"
	"api/internal/handlers"
	"api/internal/middleware"
	"api/internal/models"
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
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			db, err := utils.InitDb();
			if (err != nil) {
				os.Exit(1);
			}
			appCtx := &models.AppContext{
				Context: c,
				AppDB: db,
			}
			return next(appCtx)
		}
	})

	// Middleware which runs *after* all others, to
	// ensure connections are closed
	// HOWEVER TODO: we need to make this persist across connections
	// for latency
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(*models.AppContext)
			err := next(cc)
			slog.Info("Closing DB connection!!!")
			cc.AppDB.Close();
			return err
		}
	})

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// Routes
	e.GET("/status", func(c echo.Context) error {
		response := "OK "+time.Now().Format("2006-01-02 15:04:05")
		return c.String(http.StatusOK, response)
	})

	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)
	e.GET("/users/:username", middleware.Authenticate(handlers.GetUser));
	// e.GET("/users", middleware.Authenticate(handlers.GetUsers));
	e.GET("/users", handlers.GetUsers);

	// Start server
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server!", "error", err)
	} else {
		slog.Info("server started on port 8080")
	}
}