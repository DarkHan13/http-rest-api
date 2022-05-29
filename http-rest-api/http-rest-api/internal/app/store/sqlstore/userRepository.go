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

// Create create user
func (r *UserRepository) Create(u *models.User) error {

	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow("INSERT INTO users (email, username, password) VALUES($1, $2, $3) RETURNING id",
		u.Email,
		u.Username,
		u.Password,
	).Scan(&u.Id)
}

// FindById Find user by id
func (r *UserRepository) FindById(id int) (*models.User, error) {
	u := &models.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, username, password from users WHERE id = $1",
		id,
	).Scan(
		&u.Id,
		&u.Email,
		&u.Username,
		&u.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

// FindByEmail find user by email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	u := &models.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, username, password from users WHERE email = $1",
		email,
	).Scan(
		&u.Id,
		&u.Email,
		&u.Username,
		&u.Password,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindAll() (*[]models.User, error) {
	rows, err := r.store.db.Query("SELECT id, email, username, password FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.Id, &u.Email, &u.Username, &u.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return &users, nil
}

func (r *UserRepository) FindByUsernameLike(username string) (*[]models.User, error) {
	rows, err := r.store.db.Query("SELECT id, email, username, password FROM users" +
		" WHERE username LIKE '%" + username + "%'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.Id, &u.Email, &u.Username, &u.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return &users, nil
}

func (r *UserRepository) DeleteById(id int) error {
	if _, err := r.store.db.Query("DELETE FROM users WHERE id = $1",
		id,
	); err != nil {
		return err
	}

	return nil
}
