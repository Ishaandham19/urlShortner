package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Ishaandham19/urlShortner/models"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := u.ValidateUser(user)
	json.NewEncoder(w).Encode(resp)
}

func (u *User) ValidateUser(curUser *models.User) map[string]interface{} {
	user := &models.User{}
	log.Println(curUser.UserName)

	err := u.db.Where("user_name = ?", curUser.UserName).First(&user)
	if err.Error != nil {
		log.Println(err.Error.Error())
		var resp = map[string]interface{}{"status": false, "message": "Username not found"}
		return resp
	}
	// Sets expiration of token to 27.78 days
	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(curUser.Password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp
	}

	tk := &models.Token{
		UserID: user.ID,
		Name:   user.UserName,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	// sign the token with secret key
	tokenString, error := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"status": true, "message": "logged in"}

	resp["token"] = tokenString // Store the token in the response
	return resp
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
		err := ErrorResponse{
			Err: "Password Encryption failed",
		}
		json.NewEncoder(w).Encode(err)
	}

	user.Password = string(pass)

	createdUser := u.db.Create(user)
	var errMessage = createdUser.Error

	if createdUser.Error != nil {
		fmt.Println(errMessage)
		return
	}

	var resp = map[string]interface{}{"status": true, "message": "Successful registration"}
	json.NewEncoder(w).Encode(resp)
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

func (u *User) GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 64)
	fmt.Println(id)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Illegal request"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	var user models.User
	if err := u.db.First(&user, id); err != nil {
		log.Print(err.Error.Error())
		var resp = map[string]interface{}{"status": false, "message": "No records found"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	json.NewEncoder(w).Encode(&user)
}
