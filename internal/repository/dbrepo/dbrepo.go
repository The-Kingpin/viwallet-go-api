package dbrepo

import (
	"database/sql"

	"gitlab.com/the-kingpin/viwallet/internal/config"
	"gitlab.com/the-kingpin/viwallet/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepository {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
