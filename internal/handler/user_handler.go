package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/example/rest-api/internal/repository"
	"github.com/example/rest-api/internal/service"
	"github.com/example/rest-api/internal/domain"
	"github.com/go-chi/chi/v5"
)

// UserHandler содержит HTTP-обработчики для работы с пользователями.
type UserHandler struct {
	svc service.UserService
}

// NewUserHandler создаёт новый экземпляр обработчика пользователей.
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// RegisterRoutes регистрирует маршруты для пользователей.
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.GetAll)
		r.Get("/{id}", h.GetByID)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// Create godoc
// POST /users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.Create(r.Context(), req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

// GetAll godoc
// GET /users
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.GetAll(r.Context())
	if err != nil {
		handleServiceError(w, err)
		return
	}
	// Возвращаем пустой массив, а не null
	if users == nil {
		users = []domain.User{}
	}
	respondJSON(w, http.StatusOK, users)
}

// GetByID godoc
// GET /users/{id}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// Update godoc
// PUT /users/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req domain.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// Delete godoc
// DELETE /users/{id}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- helpers ---

type errorResponse struct {
	Error string `json:"error"`
}

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, errorResponse{Error: msg})
}

func parseID(r *http.Request) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		respondError(w, http.StatusNotFound, "user not found")
	case errors.Is(err, service.ErrValidation):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
