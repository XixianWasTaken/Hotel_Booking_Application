package dbrepo

import (
	"database/sql"
	"learningGo/cmd/internal/config"
	"learningGo/cmd/internal/repository"
)

type PostgreDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgreRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &PostgreDBRepo{
		App: a,
		DB:  conn,
	}
}
