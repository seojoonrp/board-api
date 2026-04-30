package service

import (
	"context"
	"seojoonrp/board-api/internal/apperror"
	"seojoonrp/board-api/internal/dto"
	"seojoonrp/board-api/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	user, err := s.userRepo.Create(ctx, req.Username)
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	token, err := s.createJWT(user.ID.Hex())
	if err != nil {
		return nil, apperror.NewInternal(err)
	}

	return &dto.LoginResponse{
		User: dto.UserItem{
			ID:       user.ID.Hex(),
			Username: user.Username,
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
