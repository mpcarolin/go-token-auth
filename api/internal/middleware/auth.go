package middleware

import (
	"api/internal/utils"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ValidateRequest validates the request by checking the token
// and returning an error if the token is invalid
func ValidateRequest(c echo.Context) error {
	cookie, cookieErr := c.Cookie("token");
	if cookieErr != nil {
		slog.Error("issue retrieving token", "error", cookieErr)
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	_, err := utils.ValidateToken(cookie.Value);
	if err != nil {
		slog.Error("issue validating token", "error", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	return nil;
}

func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := ValidateRequest(c);	
		if err != nil {
			return err;
		} else {
			return next(c);
		}
	}
}