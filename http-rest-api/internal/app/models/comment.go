package models

import "time"

type Comment struct {
	Id          int       `json:"id"`
	PostId      int       `json:"postId"`
	UserId      int       `json:"userId"`
	Username    string    `json:"username"`
	Text        string    `json:"text"`
	CreatedDate time.Time `json:"createdDate"`
}

func (c *Comment) BeforeCreate() {
	c.CreatedDate = time.Now()
}
