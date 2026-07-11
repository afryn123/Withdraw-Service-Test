package domain

import "errors"

var (
	ErrNotFound  = errors.New("user not found")
	ErrForbidden = errors.New("user is not allowed to manage this profile")
)
