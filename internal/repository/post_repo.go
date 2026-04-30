package repository

import (
	"context"
	"seojoonrp/board-api/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepo interface {
	Save(ctx context.Context, post domain.Post) error
	FindAll(ctx context.Context) ([]domain.Post, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Post, error)
	FindAllByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Post, error)
	IncrementView(ctx context.Context, id primitive.ObjectID) error
}

type postRepo struct {
	coll *mongo.Collection
}

func NewPostRepo(coll *mongo.Collection) PostRepo {
	return &postRepo{coll: coll}
}

func (r *postRepo) Save(ctx context.Context, post domain.Post) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.coll.InsertOne(ctx, post)
	return err
}

func (r *postRepo) FindAll(ctx context.Context) ([]domain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var posts []domain.Post
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var post domain.Post
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *postRepo) FindAllByUserID(ctx context.Context, userID primitive.ObjectID) ([]domain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var posts []domain.Post
	cursor, err := r.coll.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepo) IncrementView(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"view": 1}})
	return err
}
