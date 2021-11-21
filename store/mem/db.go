package mem

import "github.com/mattmeyers/heimdall/store"

type DB struct {
	user userTable
}

type userTable struct {
	rows   map[int]*store.User
	nextID int
}

func NewDB() *DB {
	return &DB{
		user: userTable{rows: make(map[int]*store.User), nextID: 1},
	}
}
