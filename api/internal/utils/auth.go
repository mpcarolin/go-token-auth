package utils

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// TODO: seems like a bad func, maybe rethink. Or at least badly named.
// Does this belong here...
func ValidateLogin(email string, password string) error {
	if email == "" || password == "" {
		msg := fmt.Sprintf("email and password are required. Found email: %s. Found password: %s", email, password);
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}
	if len(password) > 72 { // bcrypt max byte size if 72
		return echo.NewHTTPError(http.StatusBadRequest, "password is too long")
	}
	return nil
}

// Generates cookie for token
func GenerateCookie(token string) *http.Cookie {
	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   86400, // 1 day
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	}
	return &cookie
}