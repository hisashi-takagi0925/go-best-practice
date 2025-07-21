package valueobject

import (
	"errors"
	"regexp"
	"strings"
)

type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)
	
	if value == "" {
		return Email{}, errors.New("email cannot be empty")
	}
	
	if !emailRegex.MatchString(value) {
		return Email{}, errors.New("invalid email format")
	}
	
	return Email{value: value}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}