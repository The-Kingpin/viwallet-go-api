package driver

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQL establishes database connection
func ConnectSQLDatabase(dsn string) (*sql.DB, error) {

	log.Println("Connecting to database...")

	d, err := newDatabaseConnection(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifetime)

	err = pingDB(d)
	if err != nil {
		d.Close()
		return nil, err
	}

	log.Println("Connected to database!")

	return d, nil
}

// newDatabaseConnection opens a new database connection from given dsn- Data Source Name
func newDatabaseConnection(dsn string) (*sql.DB, error) {
	// pgx - driver name for postgres
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	if err = pingDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

// pingDB verifies connectivity with the database by pinging it
func pingDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}
