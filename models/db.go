package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// Env struct db connection
type Env struct {
	DB *sql.DB
}

var db *sql.DB

// InitDB opens connection with MySQL driver
func InitDB() (*sql.DB, error) {
	var err error

	dbName := os.Getenv("MYSQL_DB_NAME")
	dbDriver := "mysql"
	dbUser := os.Getenv("MYSQL_DB_USER")
	dbPass := os.Getenv("MYSQL_DB_PASS")
	dbHost := os.Getenv("MYSQL_DB_HOST")
	dbPort := os.Getenv("MYSQL_DB_PORT")

	// <username>:<pw>@tcp(<HOST>:<port>)/<dbname>
	db, err = sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName))

	if err != nil {
		log.Panic(err.Error())
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Panic(err.Error())
		return nil, err
	}

	return db, nil
}
