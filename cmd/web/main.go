package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"

	"gitlab.com/code-harbor/viwallet/internal/config"
	"gitlab.com/code-harbor/viwallet/internal/driver"
	"gitlab.com/code-harbor/viwallet/internal/handlers"
	"gitlab.com/code-harbor/viwallet/internal/models"
	"gitlab.com/code-harbor/viwallet/internal/render"
)

var session *scs.SessionManager
var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger

const DbHost string = "localhost"
const DbPort string = "5432"
const DbName string = "viwallet"
const DbUser string = "postgres"
const DbPassword = "LQIHr6NEvDrY@1cW0hOe1WBEA2G$&2sX"
const DSNStr = "host=%s port=%s dbname=%s user=%s password=%s"

func main() {
	conn, err := run()

	if err != nil {
		// log server error
		app.ErrorLog.Fatal(err)
	}

	defer conn.SQL.Close()

	// get the enviroment type (prod, stg, dev)
	runtimeSetup := os.Getenv("RUNTIME_SETUP")
	if runtimeSetup == "" {
		runtimeSetup = "prod"
	}

	// inProduction := flag.Bool("production", true, "Application is in production")
	// useCache := flag.Bool("cache", true, "Use template cache")
	// dbHost := flag.String("dbHost", "localhost", "Database host")
	// dbName := flag.String("dbname", "", "Database name")
	// dbUser := flag.String("dbuser", "", "Database user")
	// dbPassword := flag.String("dbpassword", "", "Database password")
	// dbPort := flag.String("dbport", "5432", "Database port")
	// dbSSL := flag.String("dbssl", "disable", "Database SSL settings (disable, prefer, require)")

	flag.Parse()

	// get the port number from envieroment variable and if not specified set it to :8080
	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		serverPort = ":8080"
	} else {
		serverPort = ":" + serverPort
	}

	app.InfoLog.Println("Application running in environment:", runtimeSetup)
	app.InfoLog.Printf("Starting the application, listening on port %s", serverPort)

	srv := &http.Server{
		Addr:    serverPort,
		Handler: routes(),
	}

	err = srv.ListenAndServe()

	app.InfoLog.Fatal(err)
}

func run() (*driver.DB, error) {
	app.InProduction = false
	app.UseCache = app.InProduction

	// Configure
	gob.Register(models.User{})

	// configure session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	var dsn string = fmt.Sprintf(DSNStr, DbHost, DbPort, DbName, DbUser, DbPassword)

	// init db connection
	dbConn, err := driver.ConnectSQLDatabase(dsn)
	if err != nil {
		dbConn.SQL.Close()
		app.ErrorLog.Fatal("Server error. Error establishin connection to the database.")
	}

	app.InfoLog.Println("Connected to database.")

	if app.UseCache {

		// create template cache
		tc, err := render.CreateTemplateCache()
		if err != nil {
			app.ErrorLog.Fatal("cannot create template cache")
			return nil, err
		}

		app.TemplateCache = tc
	}

	repo := handlers.NewRepo(&app, dbConn)
	handlers.SetRepo(repo)
	render.NewRenderer(&app)

	return dbConn, nil
}
