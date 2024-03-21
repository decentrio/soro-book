package handlers

import (
	"log"
	"os"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)


type DBHandler struct {
	db *gorm.DB
}

func NewDBHandler() *DBHandler {
	db := createConnection()
	return &DBHandler{db: db}
}

// create connection with postgres db
func createConnection() *gorm.DB {
	sqlUrl, ok := os.LookupEnv("POSTGRES_URL")

	if !ok {
		log.Fatalf("Error get POSTGRES_URL")
	}

	// Open the connection
	db, err := gorm.Open(postgres.Open(sqlUrl), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	log.Println("Connected to MySQL:", db)

	return db
}