package dbrepo

import (
	"database/sql"

	"github.com/jfk23/gobookings/internal/config"
	"github.com/jfk23/gobookings/internal/repository"
)

type postgresDBRepo struct{
	App *config.AppConfig
	DB *sql.DB
}

type testDBRepo struct{
	App *config.AppConfig

}

func NewPostgresRepo (a *config.AppConfig, conn *sql.DB) repository.DatabaseRepo{
	return &postgresDBRepo{
		App: a,
		DB: conn,
	}

}

func NewtestingDBRepo (a *config.AppConfig) repository.DatabaseRepo{
	return &testDBRepo{
		App: a,
	}
}

