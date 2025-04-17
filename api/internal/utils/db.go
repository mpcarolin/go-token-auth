package utils

import (
	"context"
	"log/slog"
	"os"

	"api/internal/db"
	"api/internal/models"

	"github.com/jackc/pgx/v5"
)

func InitDb() (*models.AppDB, error) {
	ctx := context.Background()

	connString, tmplErr := Template(
		"postgres://{{ .User }}:{{ .Password }}@{{ .Host }}:{{ .Port}}/{{ .DatabaseName }}",
		map[string]interface{}{
			"User": os.Getenv("POSTGRES_USER"),
			"Password": os.Getenv("POSTGRES_PASSWORD"),
			"Host": os.Getenv("POSTGRES_HOST"),
			"Port": os.Getenv("POSTGRES_PORT"),
			"DatabaseName": os.Getenv("POSTGRES_DB"),
		},
	)
	if tmplErr != nil {
		return nil, tmplErr
	}
	conn, err := pgx.Connect(ctx, connString)
	if (err != nil) {
		slog.Error("failed to connect to database", "error", err)
		return nil, err;
	}

	queries := db.New(conn)

	Close := func(){
		conn.Close(ctx);	
	}

	appDb := models.AppDB{
		Queries: queries,
		Close: Close,
	}

	return &appDb, nil;
} 