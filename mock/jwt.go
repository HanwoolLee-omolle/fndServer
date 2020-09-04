package mock

import (
	"omolle.com/fnd/server/model"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(*twisk.AuthUser) (string, error)
}

// GenerateToken mock
func (j *JWT) GenerateToken(u *twisk.AuthUser) (string, error) {
	return j.GenerateTokenFn(u)
}
