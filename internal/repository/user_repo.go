package repository

import (
	"context"
	"seojoonrp/board-api/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	Create(ctx context.Context, username string) (*domain.User, error)
}

type userRepo struct {
	coll *mongo.Collection
}

func NewUserRepo(coll *mongo.Collection) UserRepo {
	return &userRepo{coll: coll}
}

func (r *userRepo) Create(ctx context.Context, username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user := domain.User{
		ID:       primitive.NewObjectID(),
		Username: username,
	}

	if _, err := r.coll.InsertOne(ctx, user); err != nil {
		return nil, err
	}

	return &user, nil
}
