package entity

import "github.com/takagi_hisashi/go-best-practice/web-api/internal/domain/valueobject"

type User struct {
	id       valueobject.UserID
	name     string
	username string
	email    valueobject.Email
}

func NewUser(id valueobject.UserID, name, username string, email valueobject.Email) *User {
	return &User{
		id:       id,
		name:     name,
		username: username,
		email:    email,
	}
}

func (u *User) ID() valueobject.UserID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Email() valueobject.Email {
	return u.email
}