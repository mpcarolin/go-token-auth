package handlers

import (
	"log/slog"
	"net/http"

	"api/models"
	"api/utils"

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
		return c.String(http.StatusInternalServerError, "bad credentials")
	}

	user := users[username]

	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password));
	if err != nil {
		slog.Error("invalid username or password", "error", err)
		return c.String(http.StatusUnauthorized, "invalid username or password")
	}

	token, err := utils.GenerateToken(username);
	if err != nil {
		slog.Error("failed to generate token", "error", err)
		return c.String(http.StatusInternalServerError, "internal error")
	}

	// TODO: get expiration for cookie to be aligned with token expiration
	c.SetCookie(GenerateCookie(token));
	return c.String(http.StatusOK, "user logged in")
}

func GenerateCookie(token string) *http.Cookie {
	cookie := http.Cookie {
		Name: "token",
		Value: token,
		Path: "/",
		MaxAge: 86400,
		HttpOnly: true,
		Secure: true,
		SameSite: http.SameSiteDefaultMode,
	}
	return &cookie;
}

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

// TODO: seems like a bad func, maybe rethink. Or at least badly named.
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

func GetUser(c echo.Context) error {
	slog.Info("getting user", "username", c.Param("username"))
	err := ValidateRequest(c)
	if err != nil {
		return err
	}
	
	slog.Info("request validated", "username", c.Param("username"))
	username := c.Param("username")
	user, ok := users[username]
	if !ok {
		slog.Error("user not found", "username", username)
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	slog.Info("user found", "username", username)
	return c.JSON(http.StatusOK, user)
}


func saveUser(user models.User) error {
	users[user.Username] = user
	return nil
} 