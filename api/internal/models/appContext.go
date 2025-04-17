package models

import (
	"github.com/labstack/echo/v4"
)

type AppContext struct {
	echo.Context
	AppDB *AppDB
}

type AppHandler func(c *AppContext) error

// ToHandler converts an AppHandler to an echo.HandlerFunc
func (h AppHandler) ToHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		appCtx := c.(*AppContext)
		return h(appCtx)
	}
}

// RegisterRoute is a helper to register routes with AppContext
func RegisterRoute(e *echo.Echo, method, path string, handler AppHandler, middleware ...echo.MiddlewareFunc) {
	wrapped := handler.ToHandler()
	for _, m := range middleware {
		wrapped = m(wrapped)
	}
	e.Add(method, path, wrapped)
}