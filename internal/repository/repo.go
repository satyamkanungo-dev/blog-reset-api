package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/satyamkanungo-dev/blog-rest-api/internal/models/core"
)

type IBlogRepository interface {
	Create(ctx context.Context, userId, title, content, category string, tags ...string) (*models.Blog, error)
	Get(ctx context.Context, id, userId string) (*models.Blog, error)
	GetAll(ctx context.Context, limit int, cursor, userId string) ([]models.Blog, error)
	Update(ctx context.Context, id, userId, title, content, category string, tags ...string) error
	Delete(ctx context.Context, id, userId string) error
	DeleteMultiple(ctx context.Context, userId string, ids ...string) ([]string, error)
}

type IUserRepsoitory interface {
	Create(ctx context.Context, name, email, password, role string) (*models.User, error)
	Get(ctx context.Context, email string) (*models.User, error)
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
