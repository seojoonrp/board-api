package service

import (
	"context"
	"errors"
	"seojoonrp/board-api/internal/apperror"
	"seojoonrp/board-api/internal/domain"
	"seojoonrp/board-api/internal/repository"
)

type PostService interface {
	Create(ctx context.Context, req domain.CreatePostRequest, userID string) (*domain.Post, error)
	GetAll(ctx context.Context) ([]domain.Post, error)
	Get(ctx context.Context, id string) (*domain.Post, error)
	GetByUserID(ctx context.Context, userID string) ([]domain.Post, error)
}

type postService struct {
	postRepo repository.PostRepo
	userRepo repository.UserRepo
}

func NewPostService(postRepo repository.PostRepo, userRepo repository.UserRepo) PostService {
	return &postService{postRepo: postRepo, userRepo: userRepo}
}

func (s *postService) Create(ctx context.Context, req domain.CreatePostRequest, userID string) (*domain.Post, error) {
	post, err := s.postRepo.Create(ctx, userID, req.Title, req.Body)
	if err != nil {
		return nil, mapRepoErr(err, "user", userID)
	}

	return post, nil
}

func (s *postService) GetAll(ctx context.Context) ([]domain.Post, error) {
	posts, err := s.postRepo.FindAll(ctx)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	return posts, nil
}

func (s *postService) GetByUserID(ctx context.Context, userID string) ([]domain.Post, error) {
	posts, err := s.postRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, mapRepoErr(err, "user", userID)
	}

	return posts, nil
}

func (s *postService) Get(ctx context.Context, id string) (*domain.Post, error) {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, mapRepoErr(err, "post", id)
	}

	// TODO: Transaction
	if err := s.postRepo.IncrementView(ctx, id); err != nil {
		return nil, mapRepoErr(err, "post", id)
	}
	post.View++

	return post, nil
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
