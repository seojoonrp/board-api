package service

import (
	"context"
	"seojoonrp/board-api/internal/domain"
	"seojoonrp/board-api/internal/dto"
	"seojoonrp/board-api/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

type userService struct {
	userRepo  repository.UserRepo
	jwtSecret string
}

func NewUserService(userRepo repository.UserRepo, jwtSecret string) UserService {
	return &userService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *userService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	newUser := domain.User{
		ID:       primitive.NewObjectID(),
		Username: req.Username,
	}

	err := s.userRepo.Save(ctx, newUser)
	if err != nil {
		return nil, err
	}

	token, err := s.createJWT(newUser.ID.Hex())
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		User: dto.UserItem{
			ID:       newUser.ID.Hex(),
			Username: req.Username,
		},
		Token: token,
	}, nil
}

func (s *userService) createJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
