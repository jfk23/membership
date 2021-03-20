package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var myDB = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

// Create db pool for Postgres
func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDB(dsn)
	if err !=nil {
		panic(err)
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	myDB.SQL = db

	err = TestDB(db)
	if err !=nil {
		return nil, err
	}

	return myDB, nil

}

func TestDB (d *sql.DB) (error) {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

func NewDB (dsn string) (*sql.DB, error) {
	d, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = d.Ping()
	if err != nil {
		return nil, err
	}
	return d, nil
}