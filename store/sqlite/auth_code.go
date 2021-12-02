package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mattmeyers/heimdall/store"
)

type AuthCodeStore struct {
	db *sql.DB
}

func NewAuthCodeStore(db *sql.DB) (*AuthCodeStore, error) {
	return &AuthCodeStore{db: db}, nil
}

func (s *AuthCodeStore) GetByCode(ctx context.Context, code string) (store.AuthCode, error) {
	var c store.AuthCode
	err := s.db.
		QueryRowContext(
			ctx,
			`SELECT id, user_id, code, created_at FROM auth_code WHERE code = ?`,
			code,
		).
		Scan(&c.ID, &c.UserID, &c.Code, &c.CreatedAt)
	if err != nil {
		return store.AuthCode{}, errors.New("auth code not found")
	}

	return c, nil
}

func (s *AuthCodeStore) Create(ctx context.Context, code store.AuthCode) (int, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Commit()

	res, err := tx.Exec(
		`INSERT INTO auth_code (user_id, code, created_at) VALUES (?, ?, NOW())`,
		code.UserID,
		code.Code,
	)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return int(id), nil
}
