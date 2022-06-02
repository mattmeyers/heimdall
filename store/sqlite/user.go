package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mattmeyers/heimdall/store"
	"modernc.org/sqlite"
)

var _ store.UserStore = (*UserStore)(nil)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) (*UserStore, error) {
	return &UserStore{db: db}, nil
}

func (s *UserStore) GetByID(ctx context.Context, id int) (store.User, error) {
	q := `SELECT id, email, hash FROM user WHERE id = ?`

	var u store.User
	err := s.db.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.Email, &u.Hash)
	if err != nil {
		return store.User{}, errors.New("user not found")
	}

	return u, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (store.User, error) {
	q := `SELECT id, email, hash FROM user WHERE email = ?`

	var u store.User
	err := s.db.QueryRowContext(ctx, q, email).Scan(&u.ID, &u.Email, &u.Hash)
	if err != nil {
		return store.User{}, errors.New("user not found")
	}

	return u, nil
}

func (s *UserStore) Create(ctx context.Context, u store.User) (int, error) {
	q := `INSERT INTO user (email, hash) VALUES (?, ?)`

	res, err := s.db.ExecContext(ctx, q, u.Email, u.Hash)

	var sqlErr *sqlite.Error
	if errors.As(err, &sqlErr) && sqlErr.Code() == 2067 {
		return 0, errors.New("user already exists")
	} else if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
