package handler

import (
	"encoding/json"
	"net/http"

	"github.com/takagi_hisashi/go-best-practice/web-api/internal/interface/api/dto"
	userUseCase "github.com/takagi_hisashi/go-best-practice/web-api/internal/usecase/user"
)

type UserHandler struct {
	userService *userUseCase.Service
}

func NewUserHandler(userService *userUseCase.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.userService.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	response := make(dto.UsersResponse, len(users))
	for i, user := range users {
		response[i] = dto.UserResponse{
			ID:       user.ID().Value(),
			Name:     user.Name(),
			Username: user.Username(),
			Email:    user.Email().String(),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/users/"):]
	if id == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	response := dto.UserResponse{
		ID:       user.ID().Value(),
		Name:     user.Name(),
		Username: user.Username(),
		Email:    user.Email().String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}