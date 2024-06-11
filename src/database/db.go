package database

import (
	"fmt"
	"log"
	"os"

	"github.com/RubenPari/clear-songs/src/models"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB = nil

func Init() error {
	// credential
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// create the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	// Open the connection
	var errConnectDb error
	db, errConnectDb := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if errConnectDb != nil {
		log.Printf("Error connect database: %v", errConnectDb)
		return errConnectDb
	}

	// test connection
	errTestDb := db.Exec("SELECT 1").Error

	if errTestDb != nil {
		log.Printf("Error test connection: %v", errTestDb)
		return errTestDb
	}

	// auto-migration
	errMigration := db.AutoMigrate(&models.TrackDB{})

	if errMigration != nil {
		log.Printf("Error migration: %v", errMigration)
		return errMigration
	}

	log.Println("Successfully connected to database!")

	return nil
}
