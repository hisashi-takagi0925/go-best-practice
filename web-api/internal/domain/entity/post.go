package entity

import "github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/valueobject"

type Post struct {
	id     valueobject.PostID
	userID valueobject.UserID
	title  string
	body   string
}

func NewPost(id valueobject.PostID, userID valueobject.UserID, title, body string) *Post {
	return &Post{
		id:     id,
		userID: userID,
		title:  title,
		body:   body,
	}
}

func (p *Post) ID() valueobject.PostID {
	return p.id
}

func (p *Post) UserID() valueobject.UserID {
	return p.userID
}

func (p *Post) Title() string {
	return p.title
}

func (p *Post) Body() string {
	return p.body
}