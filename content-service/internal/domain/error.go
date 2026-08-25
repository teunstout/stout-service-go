package domain

import "errors"

const (
	InternalServerErrorMessage = "Something went wrong"
	MethodNotAllowedMessage    = "Method not allowed"
	UnauthorizedMessage        = "Unauthorized"
)

var (
	// ErrListNotFound means a sync request supplied a list id that doesn't resolve to a row
	// owned by the requesting account - a stale/foreign id, never a legitimate update target.
	ErrListNotFound = errors.New("translation list not found for this account")
	// ErrEntryNotFound is the same, for an entry id inside a sync request.
	ErrEntryNotFound = errors.New("translation entry not found for this account")
)
