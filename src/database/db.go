package database

import (
	"fmt"
	"log"

	"github.com/RubenPari/clear-songs/src/models"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB = nil

func Init() {
	// credential
	host := "localhost"
	port := 5432
	user := "username"
	password := "password"
	dbname := "dbname"

	// create the connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	// Open the connection
	var errConnectDb error
	db, errConnectDb := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if errConnectDb != nil {
		log.Fatalf("Error connecting to the database: %v", errConnectDb)
	}

	// test connection
	errTestDb := db.Exec("SELECT 1").Error

	if errTestDb != nil {
		log.Printf("Error testing the connection to the database: %v", errTestDb)
	}

	// auto-migration
	errMigration := db.AutoMigrate(&models.TrackDB{})

	if errMigration != nil {
		log.Printf("Error migrating the database: %v", errMigration)
	}

	log.Println("Successfully connected!")
}

// GetDB returns instances of db
func GetDB() *gorm.DB {
	if db == nil {
		Init()
	}
	return db
}
