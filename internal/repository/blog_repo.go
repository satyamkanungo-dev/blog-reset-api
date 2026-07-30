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

type BlogRepo struct {
	repo *Repository
}

func NewBlogRepository(db *Repository) *BlogRepo {
	return &BlogRepo{repo: db}
}

func (b *BlogRepo) Create(ctx context.Context, userId, title, content, category string, tags ...string) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	var blog models.Blog
	query := `
		INSERT INTO  blogs (user_id, title, content, category,tags)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING *;
	`

	if err := b.repo.pool.QueryRow(ctx, query, userId, title, content, category, tags).Scan(
		&blog.Id,
		&blog.UserId,
		&blog.Title,
		&blog.Content,
		&blog.Category,
		&blog.Tags,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &blog, nil
}

func (b *BlogRepo) Get(ctx context.Context, id, userId string) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	var blog models.Blog
	query := `
		SELECT * FROM blogs 
		WHERE id = $1 AND user_id = $2;
	`

	err := b.repo.pool.QueryRow(ctx, query, id, userId).Scan(
		&blog.Id,
		&blog.UserId,
		&blog.Title,
		&blog.Content,
		&blog.Category,
		&blog.Tags,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, Error.ErrBlogNotFound
		}

		return nil, fmt.Errorf("get blog: %w", err)
	}

	return &blog, nil
}

func (b *BlogRepo) GetAll(ctx context.Context, limit int, cursor, userId string) ([]models.Blog, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	baseQuery := `
		SELECT id, user_id, title, category, tags, created_at, updated_at
		FROM blogs
		WHERE user_id = $1
	`

	args := []interface{}{userId}

	if cursor == "" {
		baseQuery += " ORDER BY updated_at DESC LIMIT $2;"
		args = append(args, limit+1)
	} else {
		payload, err := decodeCursor(cursor)
		if err != nil {
			return nil, errors.New("invalid cursor")
		}
		baseQuery += `
			AND (updated_at < $2 OR (updated_at = $2 AND id < $3))
			ORDER BY updated_at DESC
			LIMIT $4;
		`
		args = append(args, payload.Updated_at, payload.Id, limit+1)
	}

	return b.queryBlogs(ctx, baseQuery, args...)
}

func (b *BlogRepo) queryBlogs(ctx context.Context, query string, args ...interface{}) ([]models.Blog, error) {
	blogs := make([]models.Blog, 0)

	rows, err := b.repo.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var blog models.Blog
		if err := rows.Scan(
			&blog.Id,
			&blog.UserId,
			&blog.Title,
			&blog.Category,
			&blog.Tags,
			&blog.CreatedAt,
			&blog.UpdatedAt,
		); err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}

func (b *BlogRepo) Update(ctx context.Context, id, userId, title, content, category string, tags ...string) (*models.Blog, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	setClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if title != "" {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argPos))
		args = append(args, title)
		argPos++
	}
	if content != "" {
		setClauses = append(setClauses, fmt.Sprintf("content = $%d", argPos))
		args = append(args, content)
		argPos++
	}
	if category != "" {
		setClauses = append(setClauses, fmt.Sprintf("category = $%d", argPos))
		args = append(args, category)
		argPos++
	}
	if len(tags) > 0 {
		setClauses = append(setClauses, fmt.Sprintf("tags = $%d", argPos))
		args = append(args, tags)
		argPos++
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no fields to update")
	}

	setClauses = append(setClauses, "updated_at = now()")

	query := fmt.Sprintf(`
		UPDATE blogs
		SET %s
		WHERE id = $%d AND user_id = $%d
		RETURNING *; 
	`, strings.Join(setClauses, ", "), argPos, argPos+1)

	args = append(args, id, userId)

	var blog models.Blog
	err := b.repo.pool.QueryRow(ctx, query, args...).Scan(
		&blog.Id,
		&blog.UserId,
		&blog.Title,
		&blog.Content,
		&blog.Category,
		&blog.Tags,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, Error.ErrBlogNotFound
		}
		return nil, fmt.Errorf("update blog: %w", err)
	}

	return &blog, nil
}

func (b *BlogRepo) Delete(ctx context.Context, id, userId string) error {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	query := `
		DELETE FROM blogs
		WHERE id = $1 AND user_id = $2;
	`

	commandTag, err := b.repo.pool.Exec(ctx, query, id, userId)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("blog with id %d not found: ", id)
	}

	return nil
}

func (b *BlogRepo) DeleteMultiple(ctx context.Context, userId string, ids ...string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	query := `
		DELETE FROM blogs
		WHERE id = ANY($1) AND user_id = $2
		RETURNING id;
	`

	return b.scanIDs(ctx, query, ids, userId)
}

func (b *BlogRepo) scanIDs(ctx context.Context, query string, args ...any) ([]string, error) {
	rows, err := b.repo.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
