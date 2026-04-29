package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username"`
}

type Post struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id"`
	Title  string             `bson:"title"`
	Body   string             `bson:"body"`
	View   int                `bson:"view"`
}
