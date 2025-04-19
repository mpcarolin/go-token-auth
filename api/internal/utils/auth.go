package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// TODO: seems like a bad func, maybe rethink. Or at least badly named.
// Does this belong here...
func ValidateLogin(email string, password string) error {
	// TODO: should validate length, format, password complexity, etc.
	if email == "" || password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "email and password are required")
	}
	if len(password) > 72 { // bcrypt max byte size if 72
		return echo.NewHTTPError(http.StatusBadRequest, "password is too long")
	}
	return nil
}

// Generates cookie for token
func GenerateCookie(token string) *http.Cookie {
	authDuration := GetAuthDuration()
	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   int(authDuration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode, // require client to be from same domain, mitigating CSRF
	}
	return &cookie
}