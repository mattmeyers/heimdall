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
	var c store.Client
	err := s.db.
		QueryRowContext(
			ctx,
			`SELECT id, client_id, client_secret FROM client WHERE client_id = ?`,
			clientID,
		).
		Scan(&c.ID, &c.ClientID, &c.ClientSecret)
	if err != nil {
		return store.Client{}, errors.New("client not found")
	}

	rows, err := s.db.QueryContext(
		ctx,
		`SELECT url FROM redirect_url WHERE client_id = ?`,
		c.ID,
	)
	if err != nil {
		return store.Client{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var row string
		if err := rows.Scan(&row); err != nil {
			return store.Client{}, err
		}
		c.RedirectURLs = append(c.RedirectURLs, row)
	}

	if err = rows.Err(); err != nil {
		return store.Client{}, err
	}

	return c, nil
}

func (s *ClientStore) Create(ctx context.Context, c store.Client) (int, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Commit()

	res, err := tx.Exec(
		`INSERT INTO client (client_id, client_secret) VALUES (?, ?)`,
		c.ClientID,
		c.ClientSecret,
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

	for _, url := range c.RedirectURLs {
		_, err = tx.Exec(
			`INSERT INTO redirect_url (client_id, url) VALUES (?, ?)`,
			id,
			url,
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return int(id), nil
}
