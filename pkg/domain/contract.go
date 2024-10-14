package domain

import "time"

type Store interface {
	CreateUser(user *User) (*User, error)
	DeleteUserById(id int64) error
	FindUserByEmail(string) (*User, error)
	FindUserById(int64) (*User, error)
	DeleteAllUsers() error

	// POST METHODS
	CreatePost(post *Post) (*Post, error)
	FindPostsByUser(userId int64) ([]Post, error)
	FindPostById(int64) (*Post, error)
	DeleteAllPosts() error
}

type JWT interface {
	CreateToken(user User, duration time.Duration) (*JWTPayload, error)
	VerifyToken(tokenString string) (*JWTPayload, error)
}
