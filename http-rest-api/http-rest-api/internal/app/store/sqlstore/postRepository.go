package sqlstore

import (
	"database/sql"
	"github.com/DarkHan13/http-rest-api/internal/app/models"
	"github.com/DarkHan13/http-rest-api/internal/app/store"
)

type PostRepository struct {
	store *Store
}

func (r *PostRepository) Create(p *models.Post) error {

	p.BeforeCreate()

	return r.store.db.QueryRow("INSERT INTO posts (user_id, username, createddate, caption) VALUES ($1, $2, $3, $4) RETURNING id",
		p.UserId,
		p.Username,
		p.CreatedDate,
		p.Caption,
	).Scan(&p.Id)
}

func (r *PostRepository) FindById(id int) (*models.Post, error) {
	p := &models.Post{}
	if err := r.store.db.QueryRow(
		"SELECT id, user_id, username, createddate, caption from posts WHERE id = $1",
		id,
	).Scan(
		&p.Id,
		&p.UserId,
		&p.Username,
		&p.CreatedDate,
		&p.Caption,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return p, nil
}

func (r *PostRepository) FindAll() (*[]models.Post, error) {
	rows, err := r.store.db.Query("SELECT id, user_id, username, createddate, caption FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err = rows.Scan(&p.Id, &p.UserId, &p.Username, &p.CreatedDate, &p.Caption)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return &posts, nil
}

func (r *PostRepository) FindAllByUserId(userId int) (*[]models.Post, error) {
	rows, err := r.store.db.Query("SELECT id, user_id, username, createddate, caption FROM posts WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err = rows.Scan(&p.Id, &p.UserId, &p.Username, &p.CreatedDate, &p.Caption)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return &posts, nil
}

func (r *PostRepository) DeleteById(id int) error {
	if _, err := r.store.db.Query("DELETE FROM posts WHERE id = $1",
		id,
	); err != nil {
		return err
	}

	return nil
}
