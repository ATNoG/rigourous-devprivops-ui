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
	e.GET("/", handlers.HomePage)
	e.GET("/demo", handlers.DemoPage)

	e.GET("/trees", handlers.TreesMainPage)
	e.GET("/trees/:tree", handlers.TreeView)

	e.GET("/descriptions", handlers.DescriptionsMainPage)

	e.GET("/reasoner", handlers.ReasonerMainPage)

	e.GET("/regulations", handlers.RegulationsMainPage)
	e.GET("/regulations/:reg", handlers.RegulationView)

	e.GET("/extra-data", handlers.ExtraDataMainPage)
	e.GET("/extra-data/:query", handlers.ExtraDataQuery)

	e.GET("/requirements", handlers.RequirementsMainPage)

	e.GET("/schemas", handlers.SchemasMainPage)

	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file")
		return
	}

	host, found := os.LookupEnv("HOST")
	if !found {
		slog.Error("'HOST' key not found in environment")
		return
	}
	port, found := os.LookupEnv("PORT")
	if !found {
		slog.Error("'PORT' key not found in environment")
		return
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", host, port)))
}
