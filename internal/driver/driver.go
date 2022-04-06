package driver

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 5
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

func ConnectPostgress(dns string) (*DB, error) {
	d, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}

	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbLifeTime)

	err = testDB(err, d)

	dbConn.SQL = d
	return dbConn, err
}

func testDB(err error, d *sql.DB) error {
	err = d.Ping()
	if err != nil {
		fmt.Print("Error!", err)
	} else {
		fmt.Print("*** Pinged database successfully! ***")
	}

	return err
}