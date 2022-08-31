package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ishaandham19/urlShortner/controllers"
	"github.com/Ishaandham19/urlShortner/routes"
	"github.com/joho/godotenv"
)

func main() {
	e := godotenv.Load()

	if e != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(e)

	port := os.Getenv("PORT")
	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	databaseHost := os.Getenv("DATABASE_HOST")

	// Connect to database
	app := App{}
	app.ConnectDb(username, password, databaseName, databaseHost)

	// Create handlers
	user := controllers.NewUser(app.Db, log.Default())

	app.Router = routes.Handlers(user)

	// Handle routes
	http.Handle("/", app.Router)

	// serve
	log.Printf("Server up on port '%s'", port)
	app.Run("", port)
}
