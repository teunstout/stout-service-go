package domain

import "errors"

const (
	InternalServerErrorMessage = "Something went wrong"
	MethodNotAllowedMessage    = "Method not allowed"
	UnauthorizedMessage        = "Unauthorized"
	ForbiddenMessage           = "Forbidden"
)

var ErrAccountAlreadyExists = errors.New("account already exists")
