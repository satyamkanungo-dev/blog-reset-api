package service

import (
	"context"

	"github.com/satyamkanungo-dev/blog-rest-api/internal/Error"
	apirequest "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_request"
	apiresponse "github.com/satyamkanungo-dev/blog-rest-api/internal/models/api_response"
	models "github.com/satyamkanungo-dev/blog-rest-api/internal/models/core"
	"github.com/satyamkanungo-dev/blog-rest-api/internal/repository"
)

type BlogService struct {
	blogRepository *repository.BlogRepo
}

func NewBlogService(br *repository.BlogRepo) *BlogService {
	return &BlogService{blogRepository: br}
}

func (bs *BlogService) Create(br *apirequest.BlogRequest, userId string) (*models.Blog, error) {
	// validate the values
	if br.Title == "" || br.Content == "" || br.Category == "" || len(br.Tags) == 0 {
		return nil, Error.ErrAllFieldRequired
	}

	blog, err := bs.blogRepository.Create(context.Background(), userId, br.Title, br.Content, br.Category, br.Tags...)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (bs *BlogService) Get(id, userId string) (*models.Blog, error) {
	if id == "" || userId == "" {
		return nil, Error.ErrMissingIdentifiers
	}

	blog, err := bs.blogRepository.Get(context.Background(), id, userId)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (bs *BlogService) Update(br *apirequest.BlogRequest, id, userId string) (*models.Blog, error) {
	if id == "" || userId == "" {
		return nil, Error.ErrMissingIdentifiers
	}

	if br.Title == "" && br.Content == "" && br.Category == "" && len(br.Tags) == 0 {
		return nil, Error.ErrAllFieldRequired
	}

	blog, err := bs.blogRepository.Update(context.Background(), id, userId, br.Title, br.Content, br.Category, br.Tags...)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (bs *BlogService) Delete(id, userId string) error {
	if id == "" || userId == "" {
		return Error.ErrMissingIdentifiers
	}

	return bs.blogRepository.Delete(context.Background(), id, userId)
}

func (bs *BlogService) GetAll(userId, cursor string) (*apiresponse.BlogsResponse, error) {
	if userId == "" {
		return nil, Error.ErrMissingIdentifiers
	}

	limit := 10
	blogs, err := bs.blogRepository.GetAll(context.Background(), limit, cursor, userId)
	if err != nil {
		return nil, err
	}

	var nxtCursor any
	hasNextPage := len(blogs) > limit
	lastblog := blogs[len(blogs)-1]

	if hasNextPage {
		nxtCursor = repository.EncodeCursorPayload(lastblog.UpdatedAt, lastblog.Id)
	} else {
		nxtCursor = nil
	}

	return &apiresponse.BlogsResponse{
		Data:       blogs,
		NextCursor: nxtCursor,
		HasNext:    hasNextPage,
	}, nil
}

func (bs *BlogService) DeleteMultiple(db *apirequest.DeleteBlogRequest, userId string) (*apiresponse.BulkDeleteResponse, error) {
	if userId == "" {
		return nil, Error.ErrMissingIdentifiers
	}

	if len(db.Ids) == 0 {
		return nil, Error.ErrAllFieldRequired
	}

	deleted, err := bs.blogRepository.DeleteMultiple(context.Background(), userId, db.Ids...)
	if err != nil {
		return nil, err
	}

	// diff to find which ids failed
	deletedSet := make(map[string]bool, len(deleted))
	for _, id := range deleted {
		deletedSet[id] = true
	}

	failed := make([]string, 0)
	for _, id := range db.Ids {
		if !deletedSet[id] {
			failed = append(failed, id)
		}
	}

	return &apiresponse.BulkDeleteResponse{
		Succeeded: deleted,
		Failed:    failed,
		Summary: apiresponse.BulkDeleteSummary{
			Total:   len(db.Ids),
			Deleted: len(deleted),
			Failed:  len(failed),
		},
	}, nil
}
