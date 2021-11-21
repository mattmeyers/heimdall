package mem

import "github.com/mattmeyers/heimdall/store"

type DB struct {
	user   userTable
	client clientTable
}

type userTable struct {
	rows   map[int]*store.User
	nextID int
}

type clientTable struct {
	rows   map[int]*store.Client
	nextID int
}

func NewDB() *DB {
	return &DB{
		user:   userTable{rows: make(map[int]*store.User), nextID: 1},
		client: clientTable{rows: make(map[int]*store.Client), nextID: 1},
	}
}
