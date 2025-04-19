package main

import (
	"api/internal/server"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	server.SetupServer(e)
	server.Start(e);
}