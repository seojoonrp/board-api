package service

import (
	"context"
	"seojoonrp/board-api/internal/domain"
	"seojoonrp/board-api/internal/dto"
	"seojoonrp/board-api/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	Login(ctx context.Context, req dto.CreateUserRequest) (*dto.UserItem, error)
}

type userService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Login(ctx context.Context, req dto.CreateUserRequest) (*dto.UserItem, error) {
	newUser := domain.User{
		ID:       primitive.NewObjectID(),
		Username: req.Username,
	}

	err := s.userRepo.Save(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &dto.UserItem{
		ID:       newUser.ID.Hex(),
		Username: req.Username,
	}, nil
}
