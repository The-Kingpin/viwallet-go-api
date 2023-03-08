package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// func init() {
// 	do some initialization here
// }

func main() {
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

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	err := http.ListenAndServe(serverPort, nil)
	if err != nil {
		log.Fatal(err)
	}

}
