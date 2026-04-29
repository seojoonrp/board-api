package repository

import (
	"context"
	"seojoonrp/board-api/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	Save(ctx context.Context, user domain.User) error
}

type userRepo struct {
	coll *mongo.Collection
}

func NewUserRepo(coll *mongo.Collection) UserRepo {
	return &userRepo{coll: coll}
}

func (r *userRepo) Save(ctx context.Context, user domain.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.coll.InsertOne(ctx, user)
	return err
}
