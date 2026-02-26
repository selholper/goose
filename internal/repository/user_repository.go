package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/example/rest-api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound возвращается, когда запись не найдена.
var ErrNotFound = errors.New("record not found")

// UserRepository определяет контракт для работы с пользователями в БД.
type UserRepository interface {
	Create(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, id int64, req domain.UpdateUserRequest) (*domain.User, error)
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository создаёт новый экземпляр репозитория пользователей.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at, updated_at
	`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, req.Name, req.Email).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("repository.Create: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repository.GetByID: %w", err)
	}
	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY id`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository.GetAll: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("repository.GetAll scan: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository.GetAll rows: %w", err)
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, id int64, req domain.UpdateUserRequest) (*domain.User, error) {
	query := `
		UPDATE users
		SET name = $1, email = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, name, email, created_at, updated_at
	`
	user := &domain.User{}
	err := r.db.QueryRow(ctx, query, req.Name, req.Email, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repository.Update: %w", err)
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	ct, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository.Delete: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
