package handlers

import (
	"api/internal/db"
	"api/internal/models"
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
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

func GetUser(c echo.Context, email string) (db.User, error) {
	cc := c.(*models.AppContext);
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 1 * 1000)
	defer cancel();

	user, err := cc.AppDB.Queries.GetUser(ctx, email);
	if err != nil {
		slog.Error("failed to get user", "error", err)
		return db.User{}, echo.NewHTTPError(http.StatusInternalServerError, "failed to get user")
	}

    return user, nil
}

func UserExists(c echo.Context, email string) (bool, error) {
	cc := c.(*models.AppContext);
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 1 * 1000)
	defer cancel();

	exists, err := cc.AppDB.Queries.UserExistsByEmail(ctx, email);
	if err != nil {
		slog.Error("failed to check user exists", "error", err)
		return false, echo.NewHTTPError(http.StatusInternalServerError, "failed to check user exists")
	}

    return exists, nil
}

func SaveUser(c echo.Context, params db.CreateUserParams) (db.CreateUserRow, error) {
	cc := c.(*models.AppContext)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 1 * 1000)
	defer cancel();

	user, err := cc.AppDB.Queries.CreateUser(ctx, params);
	if err != nil {
		slog.Error("failed to add user", "error", err)
		return db.CreateUserRow{}, echo.NewHTTPError(http.StatusInternalServerError, "failed to add user")
	}
	return user, nil
} 