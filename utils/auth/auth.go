package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Ishaandham19/urlShortner/models"
	"github.com/golang-jwt/jwt/v4"
)

// Exception struct
type Exception models.Exception

// JwtVerify Middleware function
func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var header = r.Header.Get("x-access-token") // Grab the token from the header
		log.Println("header: ", header)
		header = strings.TrimSpace(header)

		if header == "" {
			//Token is missing, returns with error code 403 Unauthorized
			log.Println("Token is missing ")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: "Missing auth token"})
			return
		}

		// HTTP only cookie auth - Remove for now
		// tokenString, err := r.Cookie("Authorization")

		// if err != nil {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	json.NewEncoder(w).Encode(map[string]interface{}{"status": "false", "message": "Invalid authorization cookie"})
		// 	return

		// }

		tk := &models.Token{}

		_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			log.Println(Exception{Message: err.Error()})
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "false", "message": "Invalid jwt token"})
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
