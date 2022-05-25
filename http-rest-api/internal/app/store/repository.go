package store

import "github.com/DarkHan13/http-rest-api/internal/app/models"

// UserRepository
type UserRepository interface {
	Create(*models.User) error
	FindByEmail(string) (*models.User, error)
	FindById(int) (*models.User, error)
	FindAll() (*[]models.User, error)
	DeleteById(int) error
}
