package main

import (
	// "log"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"api/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	username string
	hashedPassword string
}

// var key = []byte("my_secret_key")

var users = map[string]user{} // username:hashedPassword

func main() {
	// Echo instance
	e := echo.New()

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

	// Start server
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server!", "error", err)
	}
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")	
	slog.Info("attempting to login user", "username", username)

	validateErr := validateLogin(c, username, password);
	if (validateErr != nil) {
		slog.Error("invalid username or password", "error", validateErr)
		return validateErr;
	}
	if _, ok := users[username]; !ok {
		slog.Error("user not found", "username", username)
		return c.String(http.StatusNotFound, "user not found. please register first.")
	}

	user := users[username]

	err := bcrypt.CompareHashAndPassword([]byte(user.hashedPassword), []byte(password));
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

func register(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	slog.Info("registering user", "username", username)

	validateErr := validateLogin(c, username, password);
	if (validateErr != nil) {
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

	user := user{
		username: username,
		hashedPassword: string(hashed),
	}

	saveErr := saveUser(user);
	if (saveErr != nil) {
		slog.Error("failed to save user", "error", saveErr)
		return c.String(http.StatusInternalServerError, "failed to register user")
	}

	return c.String(http.StatusOK, "user registered")
}


func saveUser(user user) error {
	users[user.username] = user
	return nil
}

// func Register(w http.ResponseWriter, r *http.Request) {}
// func Login(w http.ResponseWriter, r *http.Request) {}
// func Protected(w http.ResponseWriter, r *http.Request) {}
