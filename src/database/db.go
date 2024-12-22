package database

import (
	"fmt"
	"github.com/RubenPari/clear-songs/src/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var Db *gorm.DB = nil

// Init connects to the MySQL database and performs auto-migration to create the
// tables as necessary. It sets the Db global variable to the connected database.
//
// This function returns an error if it fails to connect to the database, test the
// connection, or perform the auto-migration. If the function is successful, it
// will log a message to the console indicating that the database connection was
// successful.
func Init() error {
	// mysql credential
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// create the connection string
	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)

	// Open the connection
	var errConnectDb error
	db, errConnectDb := gorm.Open(mysql.Open(mysqlInfo), &gorm.Config{})

	if errConnectDb != nil {
		log.Printf("Error connect database: %v", errConnectDb)
		return errConnectDb
	}

	// test connection
	errTestDb := db.Exec("SELECT 1").Error

	if errTestDb != nil {
		log.Printf("Error testing connection to DB: %v", errTestDb)
		return errTestDb
	}

	// auto-migration
	errMigration := db.AutoMigrate(&models.TrackDB{})

	if errMigration != nil {
		log.Printf("Error run auto-migration DB: %v", errMigration)
		return errMigration
	}

	Db = db

	log.Println("Successfully connected to database!")

	return nil
}
