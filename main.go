package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"

	"github.com/Joao-Felisberto/devprivops-ui/handlers"
)

func main() {
	e := echo.New()

	e.Static("/static", "static")

	for _, f := range []string{
		"android-chrome-192x192.png",
		"android-chrome-512x512.png",
		"apple-touch-icon.png",
		"favicon-16x16.png",
		"favicon-32x32.png",
		"favicon.ico",
	} {
		e.Static(
			fmt.Sprintf("/%s", f),
			"static",
		)
	}
	e.Static("site.manifest", "/static/site.manifest")

	// Routes
	e.GET("/", handlers.Hello)

	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file")
	}

	host, found := os.LookupEnv("HOST")
	if !found {
		slog.Error("'HOST' key not found in environment")
	}
	port, found := os.LookupEnv("PORT")
	if !found {
		slog.Error("'PORT' key not found in environment")
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", host, port)))
}
