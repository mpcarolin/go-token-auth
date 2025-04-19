package middleware

import (
	"api/internal/models"
	"api/internal/utils"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
)

// Set up AppContext provided to each request handler
func SetupCtx(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// one connection created per request, then closed.
		// simple, if slightly inefficient approach.
		// could move to persistent connection pool if needed.
		db, err := utils.InitDb()
		if err != nil {
			os.Exit(1)
		}
		appCtx := &models.AppContext{
			Context: c,
			AppDB:   db,
		}
		return next(appCtx)
	}
}


// Middleware which runs *after* all others, to
// ensure connections are closed
func Cleanup(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*models.AppContext)
		err := next(cc)
		slog.Info("Closing DB connection!!!")
		cc.AppDB.Close()
		return err
	}
}