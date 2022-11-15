package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB = nil

// Connect to the database
// and instantiate the database
// global variable
func Connect() {
	connString := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8&parseTime=True&loc=Local"

	database, errConnDb := sql.Open("mysql", connString)

	if errConnDb != nil {
		panic(errConnDb)
	}

	db = database
}

func GetDB() *sql.DB {
	if db == nil {
		Connect()
	}
	return db
}

func Close() {
	if db != nil {
		errCloseDB := db.Close()

		if errCloseDB != nil {
			fmt.Println("Error to close db connection: ", errCloseDB)
		}
	}
}
