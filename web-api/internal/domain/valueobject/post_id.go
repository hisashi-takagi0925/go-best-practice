package valueobject

import (
	"errors"
	"strconv"
)

type PostID struct {
	value int
}

func NewPostID(value int) (PostID, error) {
	if value <= 0 {
		return PostID{}, errors.New("post ID must be positive")
	}
	return PostID{value: value}, nil
}

func NewPostIDFromString(s string) (PostID, error) {
	value, err := strconv.Atoi(s)
	if err != nil {
		return PostID{}, errors.New("invalid post ID format")
	}
	return NewPostID(value)
}

func (id PostID) Value() int {
	return id.value
}

func (id PostID) String() string {
	return strconv.Itoa(id.value)
}