package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Ishaandham19/urlShortner/models"
	"github.com/golang-jwt/jwt/v4"
)

//Exception struct
type Exception models.Exception

// JwtVerify Middleware function
func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := r.Cookie("Authorization")

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "false", "message": "Invalid authorization cookie"})
			return

		}

		tk := &models.Token{}

		_, _err := jwt.ParseWithClaims(tokenString.Value, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		if _err != nil {
			w.WriteHeader(http.StatusForbidden)
			log.Println(Exception{Message: err.Error()})
			json.NewEncoder(w).Encode(map[string]interface{}{"status": "false", "message": "Invalid jwt token"})
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
