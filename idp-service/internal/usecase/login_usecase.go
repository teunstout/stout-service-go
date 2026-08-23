package usecase

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
	"stout.dev/idp/internal/repository"
)

type LoginUsecaseInterface struct {
	accountRepo repository.AccountRepository
	sessionRepo repository.SessionTokenRepository
	csrfRepo    repository.CsrfRepository
	logger      *zap.Logger
	privateKey  *rsa.PrivateKey
}

func NewLoginUsecase(
	accountRepo repository.AccountRepository,
	sessionRepo repository.SessionTokenRepository,
	csrfRepo repository.CsrfRepository,
	logger *zap.Logger,
	privateKey *rsa.PrivateKey) *LoginUsecaseInterface {

	return &LoginUsecaseInterface{
		accountRepo: accountRepo,
		sessionRepo: sessionRepo,
		csrfRepo:    csrfRepo,
		logger:      logger,
		privateKey:  privateKey,
	}
}

func (u *LoginUsecaseInterface) Login(username string, password string) (domain.LoginResult, error) {
	account, err := u.accountRepo.GetAccountByUsername(username)
	if err != nil {
		u.logger.Info("Invalid username", zap.String("username", username), zap.Error(err))
		return domain.LoginResult{}, fmt.Errorf("Invalid username or password")
	}

	if passed := domain.CheckPasswordHash(password, account.Password); !passed {
		u.logger.Debug("Invalid password", zap.String("username", username))
		return domain.LoginResult{}, errors.New("Invalid username or password")
	}

	// Create session token and csrf token
	tokenExpiration := time.Now().Add(time.Hour * 6)
	sessionToken := domain.GenerateSecureToken()
	if err := u.sessionRepo.CreateSessionToken(sessionToken, account.ID, tokenExpiration); err != nil {
		u.logger.Debug("Error creating session token", zap.String("username", username), zap.Error(err))
		return domain.LoginResult{}, errors.New("Error creating session token")
	}

	csrfToken := domain.GenerateSecureToken()
	if err := u.csrfRepo.CreateCsrfToken(csrfToken, account.ID, tokenExpiration); err != nil {
		u.logger.Debug("Error creating csrf token", zap.String("username", username), zap.Error(err))
		return domain.LoginResult{}, errors.New("Error creating csrf token")
	}

	// Create JWT token
	jwt, err := domain.CreateJwt(account.ID, u.privateKey)
	if err != nil {
		u.logger.Warn("Error creating jwt", zap.String("username", username), zap.Error(err))
		return domain.LoginResult{}, fmt.Errorf("Error creating JWT")
	}

	u.logger.Debug("User is found in database and generated the jwt, session token, and csrf token", zap.String("username", username))
	return domain.LoginResult{
		Jwt:          jwt,
		SessionToken: sessionToken,
		CsrfToken:    csrfToken,
		ExpiresAt:    tokenExpiration,
	}, nil
}
