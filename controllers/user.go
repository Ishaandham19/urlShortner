package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ishaandham19/urlShortner/models"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
	l  *log.Logger
}

func NewUser(db *gorm.DB, l *log.Logger) *User {
	return &User{db, l}
}

type ErrorResponse struct {
	Err string
}

type error interface {
	Error() string
}

func TestAPI(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API live and kicking"))
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	curUser := &models.User{}
	err := json.NewDecoder(r.Body).Decode(curUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	var user models.User
	u.db.First(&user, "user_name = ?", curUser.UserName)
	if user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"status": false, "message": "Username not found"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Sets expiration of token to 30 days
	expiresAt := time.Now().AddDate(0, 1, 0)

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(curUser.Password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { // Password does not match!
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	tk := &models.Token{
		UserID: user.ID,
		Name:   user.UserName,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	// sign the token with secret key
	tokenString, error := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"status": true, "message": "logged in"}

	cookie := &http.Cookie{Name: "Authorization", Value: tokenString, Expires: expiresAt, HttpOnly: true}
	http.SetCookie(w, cookie)
	resp["token"] = tokenString // Store the token in the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	user.Password = string(pass)

	createdUser := u.db.Create(user)
	var errMessage = createdUser.Error

	if createdUser.Error != nil {
		fmt.Println(errMessage)
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"status": false, "message": "Failed to create user"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	var resp = map[string]interface{}{"status": true, "message": "Successful registration"}
	json.NewEncoder(w).Encode(resp)
}

func (u *User) GetUserFromJWT(w http.ResponseWriter, r *http.Request) {
	authUser := r.Context().Value("user").(*models.Token)

	var user models.User
	u.db.First(&user, "user_name = ?", authUser.Name)
	if user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"status": false, "message": "Username not found"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (u *User) FetchUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	u.db.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func (u *User) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	params := mux.Vars(r)
	var id = params["id"]
	u.db.First(&user, id)
	json.NewDecoder(r.Body).Decode(user)
	u.db.Save(&user)
	json.NewEncoder(w).Encode(&user)
}

func (u *User) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user models.User
	u.db.First(&user, id)
	u.db.Delete(&user)
	json.NewEncoder(w).Encode("User deleted")
}
