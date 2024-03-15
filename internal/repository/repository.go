package repository

import "errors"

var (
	ErrURLNotFound    = errors.New("url not found")
	ErrDuplicateLogin = errors.New("duplicate login")
)
