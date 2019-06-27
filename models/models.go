package models

import (
	"database/sql"

	// sqlite3 support
	_ "github.com/mattn/go-sqlite3"

	// postgres support
	_ "github.com/lib/pq"
)

var connection *sql.DB

// Prepare opens connection to database and prepares models
func Prepare(driverName, dataSourceName string) error {
	var err error

	connection, err = sql.Open(driverName, dataSourceName)

	if err != nil {
		return err
	}

	// Prepare others here
	err = prepareAlias(connection)

	if err != nil {
		return err
	}

	err = prepareUser(connection)

	if err != nil {
		return err
	}

	return nil
}

// Close closes connection to database
func Close() {
	closeAlias()
	closeUser()
	connection.Close()
}
