package sqlstore

import (
	"database/sql"
	"github.com/DarkHan13/http-rest-api/internal/app/models"
	"github.com/DarkHan13/http-rest-api/internal/app/store"
)

//CommentRepository type
type CommentRepository struct {
	store *Store
}

//Create new comment for post and user
func (r *CommentRepository) Create(c *models.Comment) error {

	c.BeforeCreate()

	var useless int
	err := r.store.db.QueryRow("SELECT user_id FROM posts WHERE id = $1", c.PostId).Scan(&useless)
	if err != nil {
		if err == sql.ErrNoRows {
			return store.ErrRecordNotFound
		}
		return err
	}

	err = r.store.db.QueryRow("INSERT INTO comment (user_id, username, post_id, created_date, text) "+
		"VALUES($1,"+
		" $2, "+
		"$3, "+
		"$4, "+
		"$5) RETURNING id",
		c.UserId,
		c.Username,
		c.PostId,
		c.CreatedDate,
		c.Text,
	).Scan(&c.Id)

	if err != nil {
		return err
	}

	return nil

}

func (r *CommentRepository) FindAllByPostId(postId int) (*[]models.Comment, error) {
	rows, err := r.store.db.Query("SELECT id, user_id, username, post_id, created_date, text FROM comment WHERE post_id = $1",
		postId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err = rows.Scan(&c.Id, &c.UserId, &c.Username, &c.PostId, &c.CreatedDate, &c.Text)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return &comments, nil
}

func (r *CommentRepository) DeleteById(id, userId int) error {
	if _, err := r.store.db.Query("DELETE FROM comment WHERE id = $1 AND user_id = $2",
		id,
		userId,
	); err != nil {
		return err
	}

	return nil
}
