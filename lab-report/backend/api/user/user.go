package user

import (
	"context"
	"encoding/json"
	"net/http"
	
	"strconv"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/claims"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/user"
	"github.com/gorilla/mux"
)

type UserService interface {
	GetUser(ctx context.Context, id uint) (user.User, error)
	CreateUser(ctx context.Context, newUser user.User) (user.User, error)
	DeleteUser(ctx context.Context, id uint) error
	UpdateUser(ctx context.Context, updateUser user.User) (user.User, error)
}

type JwtService interface {
	ParseTokenFromCookie(r *http.Request) (*claims.Claims, error)
}

type UserHandler struct {
	userService UserService
	jwtService  JwtService
}

func NewUserController(userService UserService, jwtService JwtService) *UserHandler {
	return &UserHandler{userService: userService, jwtService: jwtService}
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	usr, err := u.userService.GetUser(ctx, uint(id))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usr)
}

func (u *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newUser user.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := u.userService.CreateUser(ctx, newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUser)
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

	var updateUser user.User
	err = json.NewDecoder(r.Body).Decode(&updateUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updateUser.ID = uint(id)
	updatedUser, err := u.userService.UpdateUser(ctx, updateUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = u.userService.DeleteUser(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
