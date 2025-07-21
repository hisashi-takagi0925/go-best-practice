package valueobject

import (
	"errors"
	"strconv"
)

type UserID struct {
	value int
}

func NewUserID(value int) (UserID, error) {
	if value <= 0 {
		return UserID{}, errors.New("user ID must be positive")
	}
	return UserID{value: value}, nil
}

func NewUserIDFromString(s string) (UserID, error) {
	value, err := strconv.Atoi(s)
	if err != nil {
		return UserID{}, errors.New("invalid user ID format")
	}
	return NewUserID(value)
}

func (id UserID) Value() int {
	return id.value
}

func (id UserID) String() string {
	return strconv.Itoa(id.value)
}