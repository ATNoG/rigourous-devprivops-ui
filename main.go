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
	sessionmanament "github.com/Joao-Felisberto/devprivops-ui/sessionManament"
	"github.com/Joao-Felisberto/devprivops-ui/tool"
)

func main() {
	e := echo.New()

	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file: %s", err)
		// return
	}

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

	jwtSecret, found := os.LookupEnv("JWT_SECRET")
	if !found {
		slog.Error("'JWT_SECRET' variable not found in environment")
		return
	}
	sessionmanament.JWTSecret = jwtSecret
	// e.Use(echojwt.JWT([]byte(jwtSecret)))

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

	static_dir, found := os.LookupEnv("STATIC_DIR")
	if !found {
		static_dir = "static"
	}
	e.Static("/static", static_dir)

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
			static_dir,
		)
	}
	e.Static(
		"site.manifest",
		fmt.Sprintf("%s/site.manifest", static_dir),
	)

	// Routes
	e.GET("/", handlers.HomePage)
	e.GET("/credentials", handlers.GetCredentials)
	e.GET("/login", handlers.SimpleLogIn)
	// e.GET("/demo", handlers.DemoPage)

	e.GET("/auth/callback", handlers.Callback)
	e.GET("/auth", handlers.Login)
	e.GET("/logout", handlers.Logout, handlers.EnsureLoggedIn)

	e.GET("/trees", handlers.TreesMainPage, handlers.EnsureLoggedIn)
	e.GET("/trees/:tree", handlers.TreeView, handlers.EnsureLoggedIn)
	e.GET("/node/:tree/:node", handlers.EditTreeNode, handlers.EnsureLoggedIn)

	e.GET("/descriptions", handlers.DescriptionsMainPage, handlers.EnsureLoggedIn)
	e.GET("/descriptions/:desc", handlers.DescriptionEdit, handlers.EnsureLoggedIn)

	e.GET("/reasoner", handlers.ReasonerMainPage, handlers.EnsureLoggedIn)
	e.GET("/reasoner/:rule", handlers.ReasonerRuleEditor, handlers.EnsureLoggedIn)

	e.GET("/regulations", handlers.RegulationsMainPage, handlers.EnsureLoggedIn)
	e.GET("/regulations/:reg", handlers.RegulationView, handlers.EnsureLoggedIn)
	e.GET("/policy/:pol", handlers.PolicyEdit, handlers.EnsureLoggedIn)

	e.GET("/extra-data", handlers.ExtraDataMainPage, handlers.EnsureLoggedIn)
	e.GET("/extra-data/:query", handlers.ExtraDataQuery, handlers.EnsureLoggedIn)

	e.GET("/requirements", handlers.RequirementsMainPage, handlers.EnsureLoggedIn)
	e.GET("/requirements/:req", handlers.RequirementEdit, handlers.EnsureLoggedIn)

	e.GET("/schemas", handlers.SchemasMainPage, handlers.EnsureLoggedIn)
	e.GET("/schemas/:schema", handlers.SchemaEditPage, handlers.EnsureLoggedIn)

	e.POST("/save/:file", handlers.SaveEndpoint, handlers.EnsureLoggedIn)

	e.GET("/tests", handlers.TestOverview, handlers.EnsureLoggedIn)
	e.GET("/tests/:scenario", handlers.TestScenarioSelect, handlers.EnsureLoggedIn)
	e.GET("/tests/:scenario/:desc", handlers.TestScenarioEdit, handlers.EnsureLoggedIn)

	e.GET("/analyse", handlers.Analyse, handlers.EnsureLoggedIn)
	e.GET("/test", handlers.Test, handlers.EnsureLoggedIn)

	e.GET("/conflicts", handlers.MergeConflicts, handlers.EnsureLoggedIn)
	e.POST("/push", handlers.Push, handlers.EnsureLoggedIn)
	e.GET("/push/:file", handlers.SolveMergeConflict, handlers.EnsureLoggedIn)

	e.POST("/delete", handlers.DeleteFile, handlers.EnsureLoggedIn)
	e.POST("/create", handlers.CreateFile, handlers.EnsureLoggedIn)
	e.POST("/create-regulation", handlers.CreateRegulation, handlers.EnsureLoggedIn)
	e.POST("/delete-regulation", handlers.DeleteRegulation, handlers.EnsureLoggedIn)
	e.POST("/save-regulation/:reg", handlers.UpdateRegulation, handlers.EnsureLoggedIn)
	e.POST("/save-tree/:tree", handlers.UpdateTree, handlers.EnsureLoggedIn)
	e.POST("/save-report-data", handlers.UpdateExtraData, handlers.EnsureLoggedIn)
	e.POST("/save-requirements", handlers.UpdateRequirements, handlers.EnsureLoggedIn)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", host, port)))
}
