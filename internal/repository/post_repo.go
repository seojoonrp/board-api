package repository

import (
	"context"
	"errors"
	"seojoonrp/board-api/internal/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostRepo interface {
	Create(ctx context.Context, userID, title, body string) (*domain.Post, error)
	FindAll(ctx context.Context) ([]domain.Post, error)
	FindByID(ctx context.Context, id string) (*domain.Post, error)
	FindAllByUserID(ctx context.Context, userID string) ([]domain.Post, error)
	IncrementView(ctx context.Context, id string) error
}

type postRepo struct {
	coll *mongo.Collection
}

func NewPostRepo(coll *mongo.Collection) PostRepo {
	return &postRepo{coll: coll}
}

func (r *postRepo) Create(ctx context.Context, userID, title, body string) (*domain.Post, error) {
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	post := domain.Post{
		ID:     primitive.NewObjectID(),
		UserID: uID,
		Title:  title,
		Body:   body,
		View:   0,
	}

	if _, err := r.coll.InsertOne(ctx, post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *postRepo) FindAll(ctx context.Context) ([]domain.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var posts []domain.Post
	cursor, err := r.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepo) FindByID(ctx context.Context, id string) (*domain.Post, error) {
	pID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var post domain.Post
	err = r.coll.FindOne(ctx, bson.M{"_id": pID}).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &post, nil
}

func (r *postRepo) FindAllByUserID(ctx context.Context, userID string) ([]domain.Post, error) {
	uID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, domain.ErrInvalidID
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var posts []domain.Post
	cursor, err := r.coll.Find(ctx, bson.M{"user_id": uID})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *postRepo) IncrementView(ctx context.Context, id string) error {
	pID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.ErrInvalidID
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = r.coll.UpdateOne(ctx, bson.M{"_id": pID}, bson.M{"$inc": bson.M{"view": 1}})
	return err
}
