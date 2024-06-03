package postgres

import (
	"admin_api_gateway/api_test/storage"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type DB struct {
	db *sql.DB
}

func NewPostgresInit(db *sql.DB) *DB {
	return &DB{db: db}
}

func (r *DB) Set(key string, value *storage.User) error {
	value.ID = uuid.NewString()
	_, err := r.db.Exec(`INSERT INTO users (id, username, email) VALUES ($1, $2, $3)`, value.ID, value.Username, value.Email)
	if err != nil {
		return err
	}
	return nil
}

func (r *DB) Delete(key string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE id = $1`, key)
	if err != nil {
		return err
	}
	return nil
}

func (r *DB) Get(key string) (*storage.User, error) {
	var user storage.User
	err := r.db.QueryRow(`SELECT * FROM users WHERE id = $1`, key).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
