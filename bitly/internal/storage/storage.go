package storage

import "errors"

var (
	ErrURLNotFound = errors.New("URL not found")
	ErrorURLExists = errors.New("URL already exists")
)


