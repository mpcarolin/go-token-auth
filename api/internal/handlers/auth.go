package handlers

import (
	"log/slog"
	"net/http"

	"api/internal/db"
	"api/internal/models"
	"api/internal/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Route handler for logging in a user, given form values username and password
func Login(c echo.Context) error {
	cc := c.(*models.AppContext)

	// fill loginParams with data from body
	var loginParams db.CreateUserParams
	bindErr := cc.Bind(&loginParams)
	if bindErr != nil {
		slog.Error("failed to bind params from body", "error", bindErr)
		return cc.String(http.StatusBadRequest, "valid user params in body are required");
	}

	validateErr := utils.ValidateLogin(loginParams.Email, loginParams.Password);
	if validateErr != nil {
		slog.Error("Error validating the input for registration.", "error", validateErr)
		return validateErr
	}

	user, userErr := GetUser(c, loginParams.Email);
	if userErr != nil {
		slog.Error("issue finding user at login by credentials", "error", userErr)
		return cc.String(http.StatusInternalServerError, "bad credentials")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginParams.Password))
	if err != nil {
		slog.Error("invalid username or password", "error", err)
		return c.String(http.StatusUnauthorized, "invalid username or password")
	}

	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		slog.Error("failed to generate token", "error", err)
		return c.String(http.StatusInternalServerError, "internal error")
	}

	// TODO: get expiration for cookie to be aligned with token expiration
	c.SetCookie(utils.GenerateCookie(token))
	return c.String(http.StatusOK, "user logged in")
}



// Route handler for registering a user, given form values username and password
func Register(c echo.Context) error {
	cc := c.(*models.AppContext)

	// fill userParams with data from body
	var userParams db.CreateUserParams
	bindErr := cc.Bind(&userParams)
	if bindErr != nil {
		slog.Error("failed to bind params from body", "error", bindErr)
		return cc.String(http.StatusBadRequest, "valid user params in body are required");
	}

	validateErr := utils.ValidateLogin(userParams.Email, userParams.Password);
	if validateErr != nil {
		slog.Error("Error validating the input for registration.", "error", validateErr)
		return validateErr
	}

	exists, existErr := UserExists(c, userParams.Email);
	if existErr != nil {
		slog.Error("Error checking existence of user", "error", bindErr)
		return cc.String(http.StatusInternalServerError, "Issue encountered while trying to register user");
	}
	if exists {
		slog.Error("user with that email already exists", "username", userParams.Email)
		return c.String(http.StatusConflict, "email in use")
	}

	var hashed, err = bcrypt.GenerateFromPassword([]byte(userParams.Password), bcrypt.MinCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return c.String(http.StatusInternalServerError, "invalid password")
	}

	_, saveErr := SaveUser(c, db.CreateUserParams{
		Email: userParams.Email,
		Password: string(hashed),
	})
	if saveErr != nil {
		slog.Error("failed to save user", "error", saveErr)
		return c.String(http.StatusInternalServerError, "failed to register user")
	}

	return c.String(http.StatusOK, "user registered")
}
