package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.com/code-harbor/viwallet/internal/driver"
	"gitlab.com/code-harbor/viwallet/internal/handlers"
	"gitlab.com/code-harbor/viwallet/internal/repository/dbrepo"
)

const dbHost string = "localhost"
const dbPort string = "5432"
const dbName string = "viwallet"
const dbUser string = "postgres"
const dbPassword = "LQIHr6NEvDrY@1cW0hOe1WBEA2G$&2sX"
const dsnStr = "host=%s port=%s dbname=%s user=%s password=%s"

func main() {

	conn, err := run()
	if err != nil {
		// log server error
		conn.Close()
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
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

	log.Println("Application running in environment:", runtimeSetup)
	log.Printf("Starting the application, listening on port %s", serverPort)

	srv := &http.Server{
		Addr:    serverPort,
		Handler: routes(),
	}

	err = srv.ListenAndServe()

	log.Fatal(err)
}

func run() (*sql.DB, error) {
	var dsn string = fmt.Sprintf(dsnStr, dbHost, dbPort, dbName, dbUser, dbPassword)

	// init db connection
	dbConn, err := driver.ConnectSQLDatabase(dsn)
	if err != nil {
		dbConn.Close()
		log.Fatal("Server error. Error establishin connection to the database.")
	}

	repo := dbrepo.NewPostgresRepo(dbConn)
	handlers.SetDBRepo(repo)

	return dbConn, nil
}
