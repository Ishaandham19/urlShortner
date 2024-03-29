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
	r.HandleFunc("/register", user.CreateUser).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", user.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/test", user.CreateURLTest).Methods("POST")

	// Auth route
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(auth.JwtVerify)
	s.HandleFunc("/url", user.CreateURL).Methods("POST", "OPTIONS")
	s.HandleFunc("/url", user.GetAllUrls).Methods("GET")
	s.HandleFunc("/user/{id}", user.UpdateUser).Methods("PUT", "OPTIONS")
	s.HandleFunc("/user/{id}", user.DeleteUser).Methods("DELETE", "OPTIONS")
	return r
}

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for the preflight request
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, x-access-token")
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header, x-access-token")
		next.ServeHTTP(w, r)
	})
}
