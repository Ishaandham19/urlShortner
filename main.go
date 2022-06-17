package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Ishaandham19/urlShortner/handlers"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "Auth Example ", log.LstdFlags)

	// create a new serve mux and register and handlers
	h := handlers.NewAuthHandler(l)

	r := mux.NewRouter()
	r.HandleFunc("/", h.Open)
	r.HandleFunc("/auth", h.Auth)
	http.Handle("/", r)

	http.ListenAndServe(":9090", r)
	// // create a new server
	// s := http.Server{
	// 	Addr:    ":9090", // configure the bind address
	// 	Handler: m,       // set the default handler
	// 	TLSConfig: &tls.Config{
	// 		MinVersion:               tls.VersionTLS13,
	// 		PreferServerCipherSuites: true,
	// 	},
	// }

	// // start the server
	// go func() {
	// 	err := s.ListenAndServe()
	// 	if err != nil {
	// 		log.Fatal("Error starting server: ", err)
	// 		os.Exit(1)
	// 	}
	// }()

	// // trap sigterm or interrupt and gracefully shutdown server
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// signal.Notify(c, os.Kill)

	// // Block until signal is received
	// sig := <-c
	// log.Println("Got signal:", sig)

	// // gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	// ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// s.Shutdown(ctx)
}
