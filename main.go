package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"

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
	e.GET("/credentials", handlers.GetCredentials)
	e.GET("/login", handlers.SimpleLogIn)
	e.GET("/demo", handlers.DemoPage)

	e.GET("/auth/callback", handlers.Callback)
	e.GET("/logout", handlers.Logout)
	e.GET("/auth", handlers.Login)

	e.GET("/trees", handlers.TreesMainPage)
	e.GET("/trees/:tree", handlers.TreeView)
	e.GET("/node/:tree/:node", handlers.EditTreeNode)

	e.GET("/descriptions", handlers.DescriptionsMainPage)
	e.GET("/descriptions/:desc", handlers.DescriptionEdit)

	e.GET("/reasoner", handlers.ReasonerMainPage)
	e.GET("/reasoner/:rule", handlers.ReasonerRuleEditor)

	e.GET("/regulations", handlers.RegulationsMainPage)
	e.GET("/regulations/:reg", handlers.RegulationView)
	e.GET("/policy/:pol", handlers.PolicyEdit)

	e.GET("/extra-data", handlers.ExtraDataMainPage)
	e.GET("/extra-data/:query", handlers.ExtraDataQuery)

	e.GET("/requirements", handlers.RequirementsMainPage)
	e.GET("/requirements/:req", handlers.RequirementEdit)

	e.GET("/schemas", handlers.SchemasMainPage)
	e.GET("/schemas/:schema", handlers.SchemaEditPage)

	e.POST("/save/:file", handlers.SaveEndpoint)

	e.GET("/tests", handlers.TestOverview)
	e.GET("/tests/:scenario", handlers.TestScenarioSelect)
	e.GET("/tests/:scenario/:desc", handlers.TestScenarioEdit)

	e.GET("/analyse", handlers.Analyse)
	e.GET("/test", handlers.Test)

	e.GET("/conflicts", handlers.MergeConflicts)
	e.POST("/push", handlers.Push)
	e.GET("/push/:file", handlers.SolveMergeConflict)

	e.POST("/delete", handlers.DeleteFile)
	e.POST("/create", handlers.CreateFile)
	e.POST("/create-regulation", handlers.CreateRegulation)
	e.POST("/delete-regulation", handlers.DeleteRegulation)
	e.POST("/save-regulation/:reg", handlers.UpdateRegulation)
	e.POST("/save-tree/:tree", handlers.UpdateTree)
	e.POST("/save-report-data", handlers.UpdateExtraData)
	e.POST("/save-requirements", handlers.UpdateRequirements)

	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file: %s", err)
		// return
	}

	/*
		jwt_secret, found := os.LookupEnv("JWT_SECRET")
		if !found {
			slog.Error("'JWT_SECRET' variable not found in environment")
			return
		}
		e.Use(echojwt.JWT([]byte(jwt_secret)))
	*/

	store_key, found := os.LookupEnv("SESSION_STORE_KEY")
	if !found {
		slog.Error("'SESSION_STORE_KEY' variable not found in environment")
		return
	}

	gothic.Store = sessions.NewCookieStore([]byte(store_key))

	gh_key, gh_key_found := os.LookupEnv("GITHUB_KEY")
	gh_secret, gh_sec_found := os.LookupEnv("GITHUB_SECRET")
	if gh_key_found && gh_sec_found {
		goth.UseProviders(
			github.New(gh_key, gh_secret, "http://localhost:8082/auth/callback?provider=github"),
		)
	} else if !(!gh_key_found && !gh_sec_found) {
		if !gh_key_found {
			slog.Error("'GITHUB_KEY' variable not found in environment")
			return
		}
		if !gh_sec_found {
			slog.Error("'GITHUB_SECRET' variable not found in environment")
			return
		}
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
