package driver

import (
	"database/sql"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)

// DB a type for holding the a pointer to sql.DB for saving the sql connection
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxDBLifeTime = 5 * time.Minute

// NewDB make a pointer to sql.DB
func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic("This error occurred in driver.go[NewDB func]" + err.Error())
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// testDB say Ping to the db and get error or nil
func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	return nil
}

// ConnectSQL make a pointer to DB for using as connection variable
func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDB(dsn)
	if err != nil {
		panic("This error occurred in driver.go[ConnectSQL func] : " + err.Error())
	}

	db.SetMaxOpenConns(maxOpenDBConn)
	db.SetConnMaxIdleTime(maxIdleDBConn)
	db.SetConnMaxLifetime(maxDBLifeTime)

	dbConn.SQL = db
	err = testDB(db)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}
