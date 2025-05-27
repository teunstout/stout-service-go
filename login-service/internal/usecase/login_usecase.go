package usecase

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"time"

	"stout.dev/login/internal/domain"
	"stout.dev/login/internal/repository"
)

type LoginUsecaseInterface struct {
	accountRepo *repository.AccountRepositoryInterface
	sessionRepo *repository.SessionTokenRepositoryInterface
	csrfRepo    *repository.CsrfRepositoryInterface
	logger      domain.Logger
	privateKey  *rsa.PrivateKey
}

func NewLoginUsecase(
	accountRepo *repository.AccountRepositoryInterface,
	sessionRepo *repository.SessionTokenRepositoryInterface,
	csrfRepo *repository.CsrfRepositoryInterface,
	logger domain.Logger,
	privateKey *rsa.PrivateKey) *LoginUsecaseInterface {
	return &LoginUsecaseInterface{
		accountRepo: accountRepo,
		logger:      logger,
		privateKey:  privateKey,
	}
}

func (u *LoginUsecaseInterface) Login(username string, password string) (*http.Cookie, *http.Cookie, string, error) {
	account, err := u.accountRepo.GetAccount(username)
	if err != nil {
		u.logger.Info("Error getting password", map[string]interface{}{
			"username": username,
			"password": password,
		})
		return nil, nil, "", fmt.Errorf("Invalid username or password")
	}

	session, csrf, err := u.createSession(password, password, account)
	if err != nil {
		u.logger.Warn("Error creating session", map[string]interface{}{
			"username": username,
		})
	}

	jwt, err := domain.CreateJwt(account.ID, u.privateKey)
	if err != nil {
		u.logger.Warn("Error creating jwt", map[string]interface{}{
			"username": username,
		})
	}

	return session, csrf, jwt, nil
}

func (u *LoginUsecaseInterface) createSession(username string, password string, account domain.Account) (*http.Cookie, *http.Cookie, error) {
	if passed := domain.CheckPasswordHash(password, account.Password); !passed {
		u.logger.Info("Error checking password", map[string]interface{}{
			"username": username,
			"password": password,
		})
		return nil, nil, errors.New("Invalid password")
	}

	tokenExpiration := time.Now().Add(time.Hour * 1)

	sessionToken := domain.GenerateSecureToken()
	if err := u.sessionRepo.CreateSessionToken(sessionToken, account.ID, tokenExpiration); err != nil {
		u.logger.Debug("Error creating session token", map[string]interface{}{
			"username": username,
			"password": password,
		})
		return nil, nil, errors.New("Error creating session token")
	}

	csrfToken := domain.GenerateSecureToken()
	if err := u.csrfRepo.CreateCsrfToken(csrfToken, account.ID, tokenExpiration); err != nil {
		u.logger.Debug("Error creating csrf token", map[string]interface{}{
			"username": username,
			"password": password,
		})
		return nil, nil, errors.New("Error creating csrf token")
	}

	sessionCookie := http.Cookie{
		Name:     "session_token",
		Value:    sessionToken, // assuming sessionToken is a string containing the token
		Path:     "/",
		Expires:  tokenExpiration,
		HttpOnly: true,
	}

	csrfCookie := http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Path:     "/",
		Expires:  tokenExpiration,
		HttpOnly: false,
	}

	return &sessionCookie, &csrfCookie, nil
}
