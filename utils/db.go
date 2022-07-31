package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/Ishaandham19/urlShortner/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Gorm postgres dialect interface
	"github.com/joho/godotenv"
)

// type DatabaseI interface{}

// type Database struct {
// 	Db *gorm.DB
// }

// func NewDatabase(db *gorm.DB) *Database {
// 	return &Database{db}
// }

// ConnectDB function: Make database connection
func ConnectDB() *gorm.DB {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("databaseUser")
	password := os.Getenv("databasePassword")
	databaseName := os.Getenv("databaseName")
	databaseHost := os.Getenv("databaseHost")

	// Define DB connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)

	// Connect to db URI
	db, err := gorm.Open("postgres", dbURI)

	if err != nil {
		fmt.Println("error", err)
		panic(err)
	}

	// Close db when not in use
	// defer db.Close()

	// Migrate the schema
	db.AutoMigrate(
		&models.User{})

	fmt.Println("Successfully connected!", db)
	return db
}
