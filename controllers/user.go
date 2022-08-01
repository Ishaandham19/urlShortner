package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"status": false, "message": "Problem signing token"}
		log.Println(error.Error())
		json.NewEncoder(w).Encode(resp)
		return
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

	if err := createdUser.Error; err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"status": false, "message": "Failed to create user"}
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.WriteHeader(http.StatusOK)
	var resp = map[string]interface{}{"status": true, "message": "Successful registration"}
	json.NewEncoder(w).Encode(resp)
}

func (u *User) getUserFromJWT(w http.ResponseWriter, r *http.Request) *models.User {
	authUser := r.Context().Value("user").(*models.Token)
	var user models.User
	u.db.First(&user, "user_name = ?", authUser.Name)
	return &user
}

func (u *User) GetURL(w http.ResponseWriter, r *http.Request) {
	urlEntry := &models.Mapping{}
	params := mux.Vars(r)
	var userNameAndAlias = params["userAndAlias"]
	if strings.Contains(userNameAndAlias, "-") {
		s := strings.Split(userNameAndAlias, "-")
		userName := s[0]
		alias := s[1]
		u.db.Where("user_name = ? AND alias = ?", userName, alias).First(&urlEntry)
		if urlEntry.Url != "" {
			http.Redirect(w, r, string(urlEntry.Url), http.StatusFound)
		}
	}

	http.Error(w, "No such url exists", http.StatusNotFound)
}

// Create mock short url
func (u *User) CreateURLTest(w http.ResponseWriter, r *http.Request) {
	urlEntry := &models.Mapping{}
	urlEntry.UserName = "test0"
	urlEntry.Alias = "google"
	urlEntry.Url = "https://google.com"
	urlEntry.ExpirationDate = time.Now()
	u.db.Create(urlEntry)
	fmt.Println("created entry")
	json.NewEncoder(w).Encode(urlEntry)
}

func (u *User) CreateURL(w http.ResponseWriter, r *http.Request) {
	user := u.getUserFromJWT(w, r)
	if user.ID == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		var resp = map[string]interface{}{"status": false, "message": "Username not found"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	urlEntry := &models.Mapping{}
	err := json.NewDecoder(r.Body).Decode(urlEntry)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"status": false, "message": "Request doesn't have required fields"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	// Check URL is valid
	if !utils.IsValidURL(u.URL) {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	// Check alias is valid
	if u.Alias != "" && !utils.IsValidAlias(u.Alias) {
		http.Error(w, "Invalid Alias", http.StatusBadRequest)
		return
	}

	urlEntry.UserName = user.UserName
	// Expiration time set to 1 year
	urlEntry.ExpirationDate = time.Now().AddDate(1, 0, 0)
	createdUrlMapping := u.db.Create(urlEntry)

	if err := createdUrlMapping.Error; err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		var resp = map[string]interface{}{"status": false, "message": "Failed to create url entry"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	json.NewEncoder(w).Encode(urlEntry)
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
