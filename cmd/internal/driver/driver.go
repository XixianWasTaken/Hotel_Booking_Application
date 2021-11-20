package driver

import (
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

//DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdDcConn = 5
const maxDcLifetime = 5 * time.Minute

func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxIdleConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdDcConn)
	d.SetConnMaxLifetime(maxDcLifetime)

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

//try ping db
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

//create a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
