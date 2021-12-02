package store

import (
	"context"
	"time"
)

type AuthCode struct {
	ID        int
	Code      string
	UserID    int
	CreatedAt time.Time
}

type AuthCodeStore interface {
	GetByCode(ctx context.Context, code string) (AuthCode, error)
	Insert(ctx context.Context, code AuthCode) (int, error)
}
