package handlers

import (
	"log/slog"
	"net/http"

	"api/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var users = map[string]models.User{} // username:hashedPassword

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")	
	slog.Info("attempting to login user", "username", username)

	validateErr := validateLogin(c, username, password);
	if validateErr != nil {
		slog.Error("invalid username or password", "error", validateErr)
		return validateErr;
	}
	if _, ok := users[username]; !ok {
		slog.Error("user not found", "username", username)
		return c.String(http.StatusNotFound, "user not found. please register first.")
	}

	user := users[username]

	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password));
	if err != nil {
		slog.Error("invalid username or password", "error", err)
		return c.String(http.StatusUnauthorized, "invalid username or password")
	}

	return c.String(http.StatusOK, "user logged in")
}

func validateLogin(c echo.Context, username string, password string) error {
	if username == "" || password == "" {
		return c.String(http.StatusBadRequest, "username and password are required")
	}
	if len(password) > 72 { // bcrypt max byte size if 72
		return c.String(http.StatusBadRequest, "password is too long")
	}
	return nil;
}

func Register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	slog.Info("registering user", "username", username)

	validateErr := validateLogin(c, username, password);
	if validateErr != nil {
		slog.Error("invalid username or password", "error", validateErr)
		return validateErr
	}

	if _, ok := users[username]; ok {
		slog.Error("username already exists", "username", username)
		return c.String(http.StatusConflict, "username already exists")
	}

	var hashed, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return c.String(http.StatusInternalServerError, "invalid password")
	}

	user := models.User{
		Username: username,
		HashedPassword: string(hashed),
	}

	saveErr := saveUser(user);
	if saveErr != nil {
		slog.Error("failed to save user", "error", saveErr)
		return c.String(http.StatusInternalServerError, "failed to register user")
	}

	return c.String(http.StatusOK, "user registered")
}

func saveUser(user models.User) error {
	users[user.Username] = user
	return nil
} 