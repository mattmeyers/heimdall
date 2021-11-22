package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mattmeyers/heimdall/store"
)

type ClientStore struct {
	db *sql.DB
}

func NewClientStore(db *sql.DB) (*ClientStore, error) {
	return &ClientStore{db: db}, nil
}

func (s *ClientStore) GetByClientID(ctx context.Context, clientID string) (store.Client, error) {
	q := `SELECT id, client_id, client_secret FROM client WHERE client_id = ?`

	var c store.Client
	err := s.db.QueryRowContext(ctx, q, clientID).Scan(&c.ID, &c.ClientID, &c.ClientSecret)
	if err != nil {
		return store.Client{}, errors.New("client not found")
	}

	return c, nil
}

func (s *ClientStore) Create(ctx context.Context, c store.Client) (int, error) {
	q := `INSERT INTO client (client_id, client_secret) VALUES (?, ?)`

	res, err := s.db.ExecContext(ctx, q, c.ClientID, c.ClientSecret)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
