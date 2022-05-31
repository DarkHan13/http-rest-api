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

	return r.store.db.QueryRow("INSERT INTO posts (user_id, username, createddate, caption, likes) VALUES ($1, $2, $3, $4, 0) RETURNING id",
		p.UserId,
		p.Username,
		p.CreatedDate,
		p.Caption,
	).Scan(&p.Id)
}

func (r *PostRepository) FindById(id int) (*models.Post, error) {
	p := &models.Post{}
	if err := r.store.db.QueryRow(
		"SELECT id, user_id, username, createddate, caption, likes from posts WHERE id = $1",
		id,
	).Scan(
		&p.Id,
		&p.UserId,
		&p.Username,
		&p.CreatedDate,
		&p.Caption,
		&p.Likes,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return p, nil
}

func (r *PostRepository) FindAll() (*[]models.Post, error) {
	rows, err := r.store.db.Query("SELECT id, user_id, username, createddate, caption, likes FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err = rows.Scan(&p.Id, &p.UserId, &p.Username, &p.CreatedDate, &p.Caption, &p.Likes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return &posts, nil
}

func (r *PostRepository) FindAllByUserId(userId int) (*[]models.Post, error) {
	rows, err := r.store.db.Query("SELECT id, user_id, username, createddate, caption, likes FROM posts WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		err = rows.Scan(&p.Id, &p.UserId, &p.Username, &p.CreatedDate, &p.Caption, &p.Likes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return &posts, nil
}

func (r *PostRepository) DeleteById(id, userId int) error {
	if _, err := r.store.db.Query("DELETE FROM posts WHERE id = $1 AND user_id = $2",
		id,
		userId,
	); err != nil {
		return err
	}

	return nil
}

func (r *PostRepository) DeleteByIdADMIN(id int) error {
	if _, err := r.store.db.Query("DELETE FROM posts WHERE id = $1",
		id,
	); err != nil {
		return err
	}

	return nil
}

func (r *PostRepository) Like(postId, userId int) (*models.Post, error) {
	var currentId int = 0
	p := &models.Post{}
	if err := r.store.db.QueryRow(
		"SELECT id, user_id, username, createddate, caption, likes from posts WHERE id = $1",
		postId,
	).Scan(
		&p.Id,
		&p.UserId,
		&p.Username,
		&p.CreatedDate,
		&p.Caption,
		&p.Likes,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	err := r.store.db.QueryRow("(SELECT user_id FROM users_liked WHERE post_id = $1 AND user_id = $2)",
		postId,
		userId,
	).Scan(
		&currentId,
	)
	found := true
	if err == sql.ErrNoRows {
		found = false
	} else if err != nil {
		return nil, err
	}
	if found {
		_, err := r.store.db.Query("DELETE from users_liked WHERE user_id = $1 AND post_id = $2",
			userId,
			postId,
		)
		if err != nil {
			return nil, err
		}
		p.Likes = p.Likes - 1
	} else {
		_, err := r.store.db.Query("INSERT INTO users_liked VALUES($1, $2)",
			postId,
			userId,
		)
		if err != nil {
			return nil, err
		}
		p.Likes = p.Likes + 1
	}
	_, err = r.store.db.Query("UPDATE posts SET likes = $1 WHERE id = $2",
		p.Likes,
		postId,
	)
	if err != nil {
		return nil, err
	}
	return p, nil

}
