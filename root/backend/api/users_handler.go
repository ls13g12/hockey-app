package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ls13g12/hockey-app/root/backend/db"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const (
	ERR_MISSING_USERNAME = "Username is required"
	ERR_MISSING_EMAIL = "Email is required"
	ERR_MISSING_PASSWORD = "Password is required"
	ERR_EMAIL_ALREADY_EXISTS = "Email already exists"
)

type UserDTO struct {
	Username  string    `json:"username,required"`
	Email     string    `json:"email,required"`
	Password  string    `json:"password,required"`
}

type UserStore interface {
	Exists(email string) (bool, error)
	Authenticate(email string, password string) (string, error)
	GetUser(userID string) (db.User, error)
	CreateUser(user db.User) error
}

type UserModel struct {
	db *mongo.Database
}

func (um UserModel) Exists(email string) (bool, error) {
	return db.Exists(um.db, email)
}

func (um UserModel) Authenticate(email string, password string) (string, error) {
	return db.Authenticate(um.db, email, password)
}

func (um UserModel) GetUser(userID string) (db.User, error) {
	return db.GetUser(um.db, userID)
}

func (um UserModel) CreateUser(user db.User) error {
	return db.CreateUser(um.db, user)
}


func (a *api) userLogin(w http.ResponseWriter, r *http.Request) {
	var userDTO UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := a.userStore.Authenticate(userDTO.Email, userDTO.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = a.sessionManager.RenewToken(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a.sessionManager.Put(r.Context(), "authenticatedUserID", userID)

	w.WriteHeader(http.StatusOK)
}

func (a *api) userLogout(w http.ResponseWriter, r *http.Request) {

	a.sessionManager.Remove(r.Context(), "authenticatedUserID")

	w.WriteHeader(http.StatusOK)
}

func (a *api) userSignup(w http.ResponseWriter, r *http.Request) {
	var userDTO UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO ); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if userDTO.Username == "" {
		http.Error(w, ERR_MISSING_USERNAME, http.StatusBadRequest)
		return
	}

	if userDTO.Email == "" {
		http.Error(w, ERR_MISSING_EMAIL, http.StatusBadRequest)
		return
	}

	if userDTO.Password == "" {
		http.Error(w, ERR_MISSING_PASSWORD, http.StatusBadRequest)
		return
	}

	userExists, err := a.userStore.Exists(userDTO.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if userExists {
		http.Error(w, ERR_EMAIL_ALREADY_EXISTS, http.StatusConflict)
		return
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), 12)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := db.User{
		UserID: uuid.NewString(),
		Username: userDTO.Username,
		Email: userDTO.Email,
		HashedPassword: string(hashedPasswordBytes),
	}


	if err := a.userStore.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
