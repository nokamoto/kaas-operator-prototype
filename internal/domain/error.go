package domain

import "errors"

// ErrResourceNotFound is returned when a requested resource is not found.
var ErrResourceNotFound = errors.New("resource not found")
