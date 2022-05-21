package sqlstore

import (
	"database/sql"

	"github.com/DarkHan13/http-rest-api/internal/app/models"
	"github.com/DarkHan13/http-rest-api/internal/app/store"
)

//User Repository
type UserRepository struct {
	store *Store
}

//Create
func (r *UserRepository) Create(u *models.User) error {

	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow("INSERT INTO users (email, password) VALUES($1, $2) RETURNING id",
		u.Email,
		u.Password,
	).Scan(&u.Id)
}

// FindById
func (r *UserRepository) FindById(id int) (*models.User, error) {
	u := &models.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, password from users WHERE id = $1",
		id,
	).Scan(
		&u.Id,
		&u.Email,
		&u.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

// FindByEmail
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	u := &models.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, password from users WHERE email = $1",
		email,
	).Scan(
		&u.Id,
		&u.Email,
		&u.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}
