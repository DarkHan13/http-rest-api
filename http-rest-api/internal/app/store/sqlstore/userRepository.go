package sqlstore

import (
	"database/sql"
	"errors"
	"github.com/DarkHan13/http-rest-api/internal/app/models"
	"github.com/DarkHan13/http-rest-api/internal/app/store"
	"strconv"
	"time"
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

	return r.store.db.QueryRow("INSERT INTO users (email, username, password, role) VALUES($1, $2, $3, $4) RETURNING id",
		u.Email,
		u.Username,
		u.Password,
		u.Role,
	).Scan(&u.Id)
}

// FindById Find user by id
func (r *UserRepository) FindById(id int) (*models.User, error) {
	u := &models.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, username, password, role from users WHERE id = $1",
		id,
	).Scan(
		&u.Id,
		&u.Email,
		&u.Username,
		&u.Password,
		&u.Role,
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
		"SELECT id, email, username, password, role from users WHERE email = $1",
		email,
	).Scan(
		&u.Id,
		&u.Email,
		&u.Username,
		&u.Password,
		&u.Role,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FindAll() (*[]models.User, error) {
	rows, err := r.store.db.Query("SELECT id, email, username, password, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return &users, nil
}

func (r *UserRepository) FindByUsernameLike(username string) (*[]models.User, error) {
	rows, err := r.store.db.Query("SELECT id, email, username, password, role FROM users" +
		" WHERE username LIKE '%" + username + "%'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var u models.User
		err = rows.Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.Role)
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

func (r *UserRepository) UnBanById(id int) error {
	if _, err := r.store.db.Query("DELETE FROM banned_users WHERE user_id = $1",
		id,
	); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) BanById(userId int, hours int64) error {
	type bannedUser struct {
		userId     int
		bannedUser int8
		startTime  int64
		endTime    int64
	}
	user := &bannedUser{}

	if err := r.store.db.QueryRow(
		"SELECT id from users WHERE id = $1",
		userId,
	).Scan(
		&user.userId,
	); err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}

	err := r.store.db.QueryRow("(SELECT user_id, banned, start_time, end_time FROM banned_users WHERE user_id = $1)",
		userId,
	).Scan(
		&user.userId,
		&user.bannedUser,
		&user.startTime,
		&user.endTime,
	)
	found := true
	if err == sql.ErrNoRows || user.userId == 0 {
		found = false
	} else if err != nil {
		return err
	}
	if found {
		startTime := time.Now().UnixMilli()
		endTime := hours*60*60*1000 + startTime
		_, err = r.store.db.Query("UPDATE banned_users SET start_time = $1, end_time = $2 WHERE user_id = $3",
			startTime,
			endTime,
			userId,
		)
		if err != nil {
			return err
		}
	} else {
		startTime := time.Now().UnixMilli()
		endTime := hours*60*60*1000 + startTime
		_, err = r.store.db.Query("INSERT INTO banned_users (user_id, banned, start_time, end_time) VALUES ($1, $2, $3, $4)",
			userId,
			1,
			startTime,
			endTime,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepository) IsBanned(id int) error {

	var endTime int64
	err := r.store.db.QueryRow("(SELECT user_id, end_time FROM banned_users WHERE user_id = $1)",
		id,
	).Scan(
		&id,
		&endTime,
	)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}
	now := time.Now().UnixMilli()
	if endTime < now {
		r.DeleteById(id)
		return nil
	}
	term := (endTime - now) / 1000 / 60
	hours := term / 60
	minutes := term - hours*60
	message := "you're banned for " + strconv.Itoa(int(hours)) + "hours " + " and " + strconv.Itoa(int(minutes)) + " minutes"

	return errors.New(message)

}
