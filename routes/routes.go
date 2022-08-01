package routes

import (
	"net/http"

	"github.com/Ishaandham19/urlShortner/controllers"
	"github.com/Ishaandham19/urlShortner/utils/auth"
	"github.com/gorilla/mux"
)

func Handlers(user *controllers.User) *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	r.Use(CommonMiddleware)

	r.HandleFunc("/{userAndAlias}", user.GetURL).Methods("GET")
	r.HandleFunc("/api", controllers.TestAPI).Methods("GET")
	r.HandleFunc("/register", user.CreateUser).Methods("POST")
	r.HandleFunc("/login", user.Login).Methods("POST")
	r.HandleFunc("/test", user.CreateURLTest).Methods("POST")

	// Auth route
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(auth.JwtVerify)
	s.HandleFunc("/user", user.CreateURL).Methods("POST")
	s.HandleFunc("/user", user.FetchUsers).Methods("GET")
	s.HandleFunc("/user/{id}", user.UpdateUser).Methods("PUT")
	s.HandleFunc("/user/{id}", user.DeleteUser).Methods("DELETE")
	return r
}

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
