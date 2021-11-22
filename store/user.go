package store

import "context"

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

type UserStore interface {
	GetByID(ctx context.Context, id int) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, u User) (int, error)
}
