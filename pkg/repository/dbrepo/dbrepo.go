package dbrepo

import (
	"database/sql"
	"github.com/DapperBlondie/booking_system/pkg/config"
	"github.com/DapperBlondie/booking_system/pkg/repository"
)

type PostgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &PostgresDBRepo{
		App: app,
		DB:  conn,
	}
}
