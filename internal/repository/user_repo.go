package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) Get(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	query := `
		SELECT id,name,email,password,role FROM users
		WHERE email = $1
	`
	var user models.User
	err := u.repo.pool.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
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

func (u *UserRepo) Update(ctx context.Context, id, name, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argPos))
		args = append(args, name)
		argPos++
	}

	if password != "" {
		setClauses = append(setClauses, fmt.Sprintf("password = $%d", argPos))
		args = append(args, password)
		argPos++
	}

	setClauses = append(setClauses, "updated_at = now()")

	query := fmt.Sprintf(`
		UPDATE blogs
		SET %s
		WHERE id = $%d 		
		RETURNING *; 
	`, strings.Join(setClauses, ", "), argPos, argPos+1)

	args = append(args, id)

	var user models.User
	err := u.repo.pool.QueryRow(ctx, query, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, Error.ErrUserNotFound
		}
		return nil, fmt.Errorf("update user: %w", err)
	}

	return &user, nil
}
