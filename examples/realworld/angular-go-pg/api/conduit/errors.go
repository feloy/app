package conduit

import "errors"

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrNotFound          = errors.New("record not found")
	ErrUnAuthorized      = errors.New("unauthorized")
	ErrInternal          = errors.New("internal error")
)
