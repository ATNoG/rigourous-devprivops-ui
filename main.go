package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"

	"github.com/Joao-Felisberto/devprivops-ui/fs"
	"github.com/Joao-Felisberto/devprivops-ui/handlers"
	"github.com/Joao-Felisberto/devprivops-ui/tool"
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
	e.GET("/login", handlers.LogIn)
	e.GET("/demo", handlers.DemoPage)

	e.GET("/trees", handlers.TreesMainPage)
	e.GET("/trees/:tree", handlers.TreeView)
	e.GET("/node/:tree/:node", handlers.EditTreeNode)

	e.GET("/descriptions", handlers.DescriptionsMainPage)
	e.GET("/descriptions/:desc", handlers.DescriptionEdit)

	e.GET("/reasoner", handlers.ReasonerMainPage)
	e.GET("/reasoner/:rule", handlers.ReasonerRuleEditor)

	e.GET("/regulations", handlers.RegulationsMainPage)
	e.GET("/regulations/:reg", handlers.RegulationView)

	e.GET("/extra-data", handlers.ExtraDataMainPage)
	e.GET("/extra-data/:query", handlers.ExtraDataQuery)

	e.GET("/requirements", handlers.RequirementsMainPage)

	e.GET("/schemas", handlers.SchemasMainPage)
	e.GET("/schemas/:schema", handlers.SchemaEditPage)

	e.POST("/save/:file", handlers.SaveEndpoint)

	e.POST("/analyse", handlers.Analyse)
	e.POST("/test", handlers.Test)

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
	dbUser, found := os.LookupEnv("DB_USER")
	if !found {
		slog.Error("'DB_USER' variable not found in environment")
		return
	}
	dbPass, found := os.LookupEnv("DB_PASS")
	if !found {
		slog.Error("'DB_PASS' variable not found in environment")
		return
	}
	dbIp, found := os.LookupEnv("DB_IP")
	if !found {
		slog.Error("'DB_IP' variable not found in environment")
		return
	}
	dbPort, found := os.LookupEnv("DB_PORT")
	if !found {
		slog.Error("'DB_PORT' variable not found in environment")
		return
	}
	dbDataset, found := os.LookupEnv("DB_DATASET")
	if !found {
		slog.Error("'DB_DATASET' variable not found in environment")
		return
	}

	fs.LocalDir = localDir
	fs.GlobalDir = globalDir
	tool.Username = dbUser
	tool.Password = dbPass
	tool.DBIP = dbIp
	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		slog.Error("'DB_PORT' must be a numeric value")
		return
	}
	tool.DBPort = dbPortInt
	tool.Dataset = dbDataset

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", host, port)))
}
