package handlers

import (
	"api/internal/models"
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
	slog.Info("GETTING USERS!!!")
	cc := c.(*models.AppContext)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 1 * 1000)
	defer cancel();
	users, err := cc.AppDB.Queries.ListUsers(ctx);
	if err != nil {
		slog.Error("failed to get users", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get users")
	}
	return c.JSON(http.StatusOK, users)
}

func GetUser(c echo.Context) error {
	username := c.Param("username")
	user, ok := users[username]
	if !ok {
		slog.Error("user not found", "username", username)
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	slog.Info("user found", "username", username)
	return c.JSON(http.StatusOK, user)
}

func SaveUser(user models.User) error {
	users[user.Username] = user
	return nil
} 