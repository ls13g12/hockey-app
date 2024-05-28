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
	GetUser(userID string) (db.User, error)
	CreateUser(user db.User) error
	// UpdateUser(user db.User) error
	// DeleteUser(userID string) error
}

type UserModel struct {
	db *mongo.Database
}

func (um UserModel) Exists(email string) (bool, error) {
	return db.Exists(um.db, email)
}

func (um UserModel) GetUser(userID string) (db.User, error) {
	return db.GetUser(um.db, userID)
}

func (um UserModel) CreateUser(user db.User) error {
	return db.CreateUser(um.db, user)
}

// func (um UserModel) UpdateUser(user db.User) error {
// 	return db.UpdateUser(um.db, user)
// }

// func (um UserModel) DeleteUser(userID string) error {
// 	return db.DeleteUser(um.db, userID)
// }


// func (a *api) userGet(w http.ResponseWriter, r *http.Request) {
// 	id := r.PathValue("id")
// 	var err error
// 	w.Header().Set("Content-Type", "application/json")

// 	user, err := a.userStore.GetUser(id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	jsonData, err := json.Marshal(user)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write(jsonData)
// }

// func Authenticate(db *mongo.Database, user User) error {
// 	coll := db.Collection("users")

// 	err = CheckPasswordHash(existingUser.HashedPassword, user.Password)
// 	return err
// }

func (a *api) userLogin(w http.ResponseWriter, r *http.Request) {
	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// err := a.userStore.Authenticate(user)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusUnauthorized)
	// 	return
	// }

	// err = a.sessionManager.RenewToken(r.Context())
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	a.sessionManager.Put(r.Context(), "authenticatedUserID", user.UserID)

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

// func (a *api) userPut(w http.ResponseWriter, r *http.Request) {
// 	var user db.User
// 	err := json.NewDecoder(r.Body).Decode(&user)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	if err := a.userStore.UpdateUser(user); err != nil {
// 		a.logger.Error("Error %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// func (a *api) userDelete(w http.ResponseWriter, r *http.Request) {
// 	id := r.PathValue("id")

// 	if err := a.userStore.DeleteUser(id); err != nil {
// 		a.logger.Error("Error %v", err)
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }
