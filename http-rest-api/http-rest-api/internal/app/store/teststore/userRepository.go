package teststore

import (
	"github.com/DarkHan13/http-rest-api/internal/app/models"
	"github.com/DarkHan13/http-rest-api/internal/app/store"
)

//User Repository
type UserRepository struct {
	store *Store
	users map[int]*models.User
}

//Create
func (r *UserRepository) Create(u *models.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.Id = len(r.users) + 1
	r.users[u.Id] = u
	u.Id = len(r.users)

	return nil
}

//Find By Id
func (r *UserRepository) FindById(id int) (*models.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}

//Find By Email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {

	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, store.ErrRecordNotFound
}
