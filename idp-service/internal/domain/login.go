package domain

import "time"

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResult struct {
	Jwt          string
	SessionToken string
	CsrfToken    string
	ExpiresAt    time.Time
}
