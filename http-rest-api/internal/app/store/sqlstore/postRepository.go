package sqlstore

import "github.com/DarkHan13/http-rest-api/internal/app/models"

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
