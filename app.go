package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ishaandham19/urlShortner/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	Router *mux.Router
	Db     *gorm.DB
}

func (a *App) ConnectDb(username, password, databaseName, databaseHost string) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Define DB connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)

	// Connect to db URI
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})

	// Add db to App struct
	a.Db = db

	if err != nil {
		fmt.Println("error", err)
		panic(err)
	}

	// Close db when not in use
	// defer db.Close()

	// Migrate the schema
	db.AutoMigrate(
		&models.User{},
		&models.Mapping{},
	)

	fmt.Println("Successfully connected!", a.Db)
}

func (a *App) Run(addr, port string) {
	log.Fatal(http.ListenAndServe(addr+":"+port, nil))
}
