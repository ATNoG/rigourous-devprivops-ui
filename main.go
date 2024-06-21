package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
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
	e.GET("/descriptions/:desc", handlers.DescriptionEdit)

	e.GET("/reasoner", handlers.ReasonerMainPage)

	e.GET("/regulations", handlers.RegulationsMainPage)
	e.GET("/regulations/:reg", handlers.RegulationView)

	e.GET("/extra-data", handlers.ExtraDataMainPage)
	e.GET("/extra-data/:query", handlers.ExtraDataQuery)

	e.GET("/requirements", handlers.RequirementsMainPage)

	e.GET("/schemas", handlers.SchemasMainPage)

	e.POST("/save/:file", handlers.SaveEndpoint)

	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file")
		return
	}

	host, found := os.LookupEnv("HOST")
	if !found {
		slog.Error("'HOST' variable not found in environment")
		return
	}
	port, found := os.LookupEnv("PORT")
	if !found {
		slog.Error("'PORT' variable not found in environment")
		return
	}
	localDir, found := os.LookupEnv("LOCAL_DIR")
	if !found {
		slog.Error("'LOCAL_DIR' variable not found in environment")
		return
	}
	globalDir, found := os.LookupEnv("GLOBAL_DIR")
	if !found {
		slog.Error("'GLOBAL_DIR' variable not found in environment")
		return
	}

	fs.LocalDir = localDir
	fs.GlobalDir = globalDir

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", host, port)))
}