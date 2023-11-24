package storage

import "errors"

//ошибки которые будут возвращаться, чтобы на сервисном слое можно понять было что не так
var (
	ErrUserExists = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound = errors.New("app not found")
)