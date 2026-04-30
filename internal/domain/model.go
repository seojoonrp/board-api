package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username"      json:"username"`
}

type Post struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID primitive.ObjectID `bson:"user_id"       json:"user_id"`
	Title  string             `bson:"title"         json:"title"`
	Body   string             `bson:"body"          json:"body"`
	View   int                `bson:"view"          json:"view"`
}

type CreatePostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
