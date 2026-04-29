package service

import (
	"context"
	"seojoonrp/board-api/internal/domain"
	"seojoonrp/board-api/internal/dto"
	"seojoonrp/board-api/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err // TODO: Error code
	}

	newPost := domain.Post{
		ID:     primitive.NewObjectID(),
		UserID: uID,
		Title:  req.Title,
		Body:   req.Body,
		View:   0,
	}

	err = s.postRepo.Save(ctx, newPost)
	if err != nil {
		return nil, err
	}

	return &dto.PostItem{
		ID:    newPost.ID.Hex(),
		Title: newPost.Title,
		Body:  newPost.Body,
		View:  newPost.View,
	}, nil
}

func (s *postService) GetAll(ctx context.Context) ([]dto.PostItem, error) {
	posts, err := s.postRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]dto.PostItem, 0, len(posts))
	for _, p := range posts {
		items = append(items, dto.PostItem{
			ID:    p.ID.Hex(),
			Title: p.Title,
			Body:  p.Body,
			View:  p.View,
		})
	}

	return items, nil
}

func (s *postService) GetByUserID(ctx context.Context, userID string) ([]dto.PostItem, error) {
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err // TODO: Error code
	}

	posts, err := s.postRepo.FindAllByUserID(ctx, uID)
	if err != nil {
		return nil, err
	}

	items := make([]dto.PostItem, 0, len(posts))
	for _, p := range posts {
		items = append(items, dto.PostItem{
			ID:    p.ID.Hex(),
			Title: p.Title,
			Body:  p.Body,
			View:  p.View,
		})
	}

	return items, nil
}

func (s *postService) Get(ctx context.Context, id string) (*dto.PostItem, error) {
	pID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err // TODO: Error code
	}

	post, err := s.postRepo.FindByID(ctx, pID)
	if err != nil {
		return nil, err
	}

	// TODO: Transaction
	err = s.postRepo.IncrementView(ctx, pID)
	if err != nil {
		return nil, err
	}
	post.View++

	return &dto.PostItem{
		ID:    post.ID.Hex(),
		Title: post.Title,
		Body:  post.Body,
		View:  post.View,
	}, nil
}
