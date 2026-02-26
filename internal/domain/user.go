package domain

import "time"

// User — доменная модель пользователя.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserRequest — запрос на создание пользователя.
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UpdateUserRequest — запрос на обновление пользователя.
type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
