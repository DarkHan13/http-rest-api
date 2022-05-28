package models

import "time"

type Post struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	Username    string    `json:"username"`
	Caption     string    `json:"caption"`
	CreatedDate time.Time `json:"createdDate"`
}

// BeforeCreate Post
func (p *Post) BeforeCreate() {
	p.CreatedDate = time.Now()
}
