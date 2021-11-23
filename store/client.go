package store

import "context"

type Client struct {
	ID           int      `json:"id"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURLs []string `json:"redirect_urls"`
}

type ClientStore interface {
	GetByClientID(ctx context.Context, id string) (Client, error)
	Create(ctx context.Context, c Client) (int, error)
}
