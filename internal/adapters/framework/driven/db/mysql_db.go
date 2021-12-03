package db

import (
	"database/sql"
	"log"

	// blank import for mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Adapter implements the DbPort interface
type Adapter struct {
	db *sql.DB
}

// NewAdapter creates a new Adapter
func NewAdapter(driverName, dataSourceName string) (*Adapter, error) {
	// connect
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalf("db connection failur: %v", err)
	}

	// test db connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("db ping failure: %v", err)
	}

	return &Adapter{db: db}, nil
}

// CloseDbConnection closes the db  connection
func (da Adapter) CloseDbConnection() {
	err := da.db.Close()
	if err != nil {
		log.Fatalf("db close failure: %v", err)
	}
}
