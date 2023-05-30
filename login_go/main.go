package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"

	"login/dao"
)

var db = dao.Connect()

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

type Userinfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func createSessionToken() string {
	// Generate a new unique session token
	sessionToken := uuid.New().String()
	return sessionToken
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var userinfo Userinfo

	err := json.NewDecoder(r.Body).Decode(&userinfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{"error", "Invalid request payload"})
		return
	}

	result, _ := dao.CheckUserExist(db, userinfo.Username)
	if result {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(Response{"error", "Username already exists"})
		return
	}

	// create user
	dao.CreateUser(db, userinfo.Username, userinfo.Password)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{"success", "User created successfully"})
}

func LogInHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{"error", "Invalid request payload"})
		return
	}

	result, _ := dao.CheckUserPassword(db, creds.Username, creds.Password)
	if !result {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{"error", "Invalid credentials"})
		return
	}

	//TODO: create session token

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{"success", "User logged in successfully"})
}

//not used
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//not used
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	dao.CreateUserTable(db)
	// dao.CreateUser(db, "test", "test")

	r := mux.NewRouter()

	// Sign up route
	r.HandleFunc("/user/api/signup", SignUpHandler).Methods("POST")

	// Log in route
	r.HandleFunc("/user/api/login", LogInHandler).Methods("POST")

	http.Handle("/", r)
	fmt.Println("Server is running...")
	http.ListenAndServe(":8080", nil)

	http.ListenAndServe(":80", r)
}
