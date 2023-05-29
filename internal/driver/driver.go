package driver

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const MaxOpenDbConn = 10
const MaxIdleDbConn = 5
const MaxDbLifetime = 5 * time.Minute

// ConnectSQL establishes database connection
func ConnectSQLDatabase(dsn string) (*DB, error) {

	log.Println("Connecting to database...")

	d, err := newDatabaseConnection(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(MaxOpenDbConn)
	d.SetMaxIdleConns(MaxIdleDbConn)
	d.SetConnMaxLifetime(MaxDbLifetime)

	dbConn.SQL = d

	err = pingDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
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
