package dto

type PostItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	View  int    `json:"view"`
}

type UserItem struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// POST /posts
type CreatePostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// POST /users
type CreateUserRequest struct {
	Username string `json:"username"`
}
