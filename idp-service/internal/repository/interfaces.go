package repository

import (
	"time"

	"stout.dev/idp/internal/domain"
)

type AccountRepository interface {
	AccountExists(username string) (bool, error)
	CreateAccount(username string, password string) error
	GetAccountByUsername(username string) (domain.Account, error)
}

type SessionTokenRepository interface {
	CreateSessionToken(token string, accountID int32, expiresAt time.Time) error
	GetAccountIdBySessionToken(token string) (int32, error)
	DeleteSessionTokensByAccountId(accountID int32) error
	DeleteSessionToken(token string) error
}

type CsrfRepository interface {
	CreateCsrfToken(token string, accountID int32, expiresAt time.Time) error
	GetCsrfIdBySessionToken(token string) (int32, error)
	DeleteCsrfTokensByAccountId(accountID int32) error
	DeleteCsrfToken(token string) error
}
