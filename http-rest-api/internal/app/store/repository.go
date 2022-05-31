package store

import "github.com/DarkHan13/http-rest-api/internal/app/models"

// UserRepository
type UserRepository interface {
	Create(*models.User) error
	FindByEmail(string) (*models.User, error)
	FindByUsernameLike(string) (*[]models.User, error)
	FindById(int) (*models.User, error)
	FindAll() (*[]models.User, error)
	DeleteById(int) error
	BanById(int, int64) error
	UnBanById(id int) error
	IsBanned(int) error
}

type PostRepository interface {
	Create(post *models.Post) error
	FindAllByUserId(int) (*[]models.Post, error)
	FindAll() (*[]models.Post, error)
	FindById(int) (*models.Post, error)
	DeleteById(int, int) error
	DeleteByIdADMIN(int) error
	Like(int, int) (*models.Post, error)
}

type CommentRepository interface {
	Create(comment *models.Comment) error
	FindAllByPostId(int) (*[]models.Comment, error)
	DeleteById(int, int) error
	DeleteByIdADMIN(int) error
}
