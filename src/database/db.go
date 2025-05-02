package database

import (
	"fmt"
	"log"
	"os"

	"github.com/RubenPari/clear-songs/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB = nil

// Init connects to the Postgres database and performs auto-migration to create the
// tables as necessary. It sets the Db global variable to the connected database.
//
// This function returns an error if it fails to connect to the database, test the
// connection, or perform the auto-migration. If the function is successful, it
// will log a message to the console indicating that the database connection was
// successful.
func Init() error {
	// postgres credentials
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// create the connection string
	postgresInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, password, dbname)

	// Open the connection
	var errConnectDb error
	db, errConnectDb := gorm.Open(postgres.Open(postgresInfo), &gorm.Config{})

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
