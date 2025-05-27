package domain

import "time"

type Account struct {
	ID        int32
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
