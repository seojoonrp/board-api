package service

import (
	"context"
	"errors"
	"seojoonrp/board-api/internal/apperror"
	"seojoonrp/board-api/internal/domain"
	"seojoonrp/board-api/internal/dto"
	"seojoonrp/board-api/internal/repository"
)

type PostService interface {
	Create(ctx context.Context, req dto.CreatePostRequest, userID string) (*dto.PostItem, error)
	GetAll(ctx context.Context) ([]dto.PostItem, error)
	Get(ctx context.Context, id string) (*dto.PostItem, error)
	GetByUserID(ctx context.Context, userID string) ([]dto.PostItem, error)
}

type postService struct {
	postRepo repository.PostRepo
	userRepo repository.UserRepo
}

func NewPostService(postRepo repository.PostRepo, userRepo repository.UserRepo) PostService {
	return &postService{postRepo: postRepo, userRepo: userRepo}
}

func (s *postService) Create(ctx context.Context, req dto.CreatePostRequest, userID string) (*dto.PostItem, error) {
	post, err := s.postRepo.Create(ctx, userID, req.Title, req.Body)
	if err != nil {
		return nil, mapRepoErr(err, "user", userID)
	}

	item := toPostItem(post)
	return &item, nil
}

func (s *postService) GetAll(ctx context.Context) ([]dto.PostItem, error) {
	posts, err := s.postRepo.FindAll(ctx)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	return toPostItems(posts), nil
}

func (s *postService) GetByUserID(ctx context.Context, userID string) ([]dto.PostItem, error) {
	posts, err := s.postRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, mapRepoErr(err, "user", userID)
	}

	return toPostItems(posts), nil
}

func (s *postService) Get(ctx context.Context, id string) (*dto.PostItem, error) {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, mapRepoErr(err, "post", id)
	}

	// TODO: Transaction
	if err := s.postRepo.IncrementView(ctx, id); err != nil {
		return nil, mapRepoErr(err, "post", id)
	}
	post.View++

	item := toPostItem(post)
	return &item, nil
}

func mapRepoErr(err error, kind, id string) error {
	switch {
	case errors.Is(err, repository.ErrInvalidID):
		return apperror.NewBadRequest("invalid " + kind + " id: " + id)
	case errors.Is(err, repository.ErrNotFound):
		return apperror.NewNotFound(kind + " not found: " + id)
	default:
		return apperror.NewInternal(err)
	}
}

func toPostItem(p *domain.Post) dto.PostItem {
	return dto.PostItem{
		ID:    p.ID.Hex(),
		Title: p.Title,
		Body:  p.Body,
		View:  p.View,
	}
}

func toPostItems(posts []domain.Post) []dto.PostItem {
	items := make([]dto.PostItem, 0, len(posts))
	for _, p := range posts {
		items = append(items, toPostItem(&p))
	}
	return items
}
