package dbrepo

import (
	"database/sql"

	"gitlab.com/code-harbor/viwallet/internal/config"
	"gitlab.com/code-harbor/viwallet/internal/repository"
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
