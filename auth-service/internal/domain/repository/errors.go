package repository

import "errors"

var (
	// ErrNotFound is returned when a requested resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrDuplicateKey is returned when attempting to create a resource with a duplicate unique key
	ErrDuplicateKey = errors.New("duplicate key")
)
