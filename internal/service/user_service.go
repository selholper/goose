package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/example/rest-api/internal/domain"
	"github.com/example/rest-api/internal/repository"
)

// ErrValidation возвращается при ошибке валидации входных данных.
var ErrValidation = errors.New("validation error")

// UserService определяет бизнес-логику для работы с пользователями.
type UserService interface {
	Create(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	Update(ctx context.Context, id int64, req domain.UpdateUserRequest) (*domain.User, error)
	Delete(ctx context.Context, id int64) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService создаёт новый экземпляр сервиса пользователей.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	if err := validateCreateRequest(req); err != nil {
		return nil, err
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	user, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("service.Create: %w", err)
	}
	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrValidation)
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service.GetByID: %w", err)
	}
	return user, nil
}

func (s *userService) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.GetAll: %w", err)
	}
	return users, nil
}

func (s *userService) Update(ctx context.Context, id int64, req domain.UpdateUserRequest) (*domain.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrValidation)
	}
	if err := validateUpdateRequest(req); err != nil {
		return nil, err
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	user, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, fmt.Errorf("service.Update: %w", err)
	}
	return user, nil
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrValidation)
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service.Delete: %w", err)
	}
	return nil
}

func validateCreateRequest(req domain.CreateUserRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("%w: email is required", ErrValidation)
	}
	if !strings.Contains(req.Email, "@") {
		return fmt.Errorf("%w: email is invalid", ErrValidation)
	}
	return nil
}

func validateUpdateRequest(req domain.UpdateUserRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("%w: email is required", ErrValidation)
	}
	if !strings.Contains(req.Email, "@") {
		return fmt.Errorf("%w: email is invalid", ErrValidation)
	}
	return nil
}
