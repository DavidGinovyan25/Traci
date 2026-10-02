package application

import "errors"

var (
	ErrUnauthorized     = errors.New("authentication required")
	ErrForbidden        = errors.New("access denied")
	ErrBlocked          = errors.New("user is blocked")
	ErrSelfModification = errors.New("administrator cannot block or delete their own account")
)
