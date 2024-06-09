package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB = nil

func Init() {
	// credential
	host := "localhost"
	port := 5432
	user := "username"
	password := "password"
	dbname := "dbname"

	// create the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open the connection
	var errConnectDb error
	db, errConnectDb = sql.Open("postgres", psqlInfo)

	if errConnectDb != nil {
		log.Fatalf("Error connecting to the database: %v", errConnectDb)
	}

	// test connection
	errTestDb := db.Ping()

	if errTestDb != nil {
		log.Fatalf("Error testing the connection to the database: %v", errTestDb)
	}

	fmt.Println("Successfully connected!")
}

// GetDB returns instances of db
func GetDB() *sql.DB {
	if db == nil {
		Init()
	}
	return db
}
