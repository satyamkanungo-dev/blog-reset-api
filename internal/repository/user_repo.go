package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/Error"
	models "github.com/satyamkanungo-dev/blog-rest-api/internal/models/core"
)

type UserRepo struct {
	repo *Repository
}

func NewUserRepository(db *Repository) *UserRepo {
	return &UserRepo{repo: db}
}

func (u *UserRepo) Create(ctx context.Context, name, email, password, role string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	var user models.User
	query := `
		INSERT INTO users (name,email,password,role)
		VALUES ($1,$2,$3,$4)
		RETURNING id,name,email,role;
	`

	if err := u.repo.pool.QueryRow(ctx, query, name, email, password, role).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Role,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) Get(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	query := `
		SELECT id,name,password,role FROM users
		WHERE email = $1
	`
	var user models.User
	err := u.repo.pool.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, Error.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	return &user, nil
}
