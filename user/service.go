package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Errors callers can check with errors.Is.
var (
	ErrInvalidName  = errors.New("name is required")
	ErrInvalidEmail = errors.New("email is invalid")
)

// Service holds the business rules for users.
// It depends on the Repository INTERFACE — it has no idea whether
// the data lives in MySQL, Postgres, or a fake used in tests.
type Service struct {
	repo Repository
}

// NewService creates a user service on top of any Repository.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Setup prepares storage (creates the users table).
func (s *Service) Setup(ctx context.Context) error {
	return s.repo.CreateTable(ctx)
}

// Register validates input, then saves a new user. Returns the new ID.
func (s *Service) Register(ctx context.Context, name, email string) (int, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	// Business rules live HERE — written once, for every database.
	if name == "" {
		return 0, ErrInvalidName
	}
	if !strings.Contains(email, "@") {
		return 0, ErrInvalidEmail
	}

	id, err := s.repo.Create(ctx, name, email)
	if err != nil {
		return 0, fmt.Errorf("register user: %w", err)
	}
	return id, nil
}

// Get returns a user by ID.
func (s *Service) Get(ctx context.Context, id int) (User, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("get user %d: %w", id, err)
	}
	return u, nil
}
