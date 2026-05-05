package usecase

import (
	"fmt"

	"go.uber.org/zap"
	"stout.dev/idp/internal/repository"
)

type AuthenticationUsecaseInterface struct {
	sessionRepo *repository.SessionTokenRepositoryInterface
	csrfRepo    *repository.CsrfRepositoryInterface
	logger      *zap.Logger
}

func NewAuthenticationUsecase(sessionRepo *repository.SessionTokenRepositoryInterface, csrfRepo *repository.CsrfRepositoryInterface, logger *zap.Logger) *AuthenticationUsecaseInterface {
	return &AuthenticationUsecaseInterface{
		sessionRepo: sessionRepo,
		csrfRepo:    csrfRepo,
		logger:      logger,
	}
}

func (u *AuthenticationUsecaseInterface) Register(sessionToken string, csrfToken string) error {
	accountID, err := u.sessionRepo.GetAccountIdBySessionToken(sessionToken)
	if err != nil {
		u.logger.Error("Error getting account id", zap.String("sessionToken", sessionToken), zap.Error(err))
		return err
	}

	csrfAccountID, err := u.csrfRepo.GetCsrfIdBySessionToken(csrfToken)
	if err != nil {
		u.logger.Error("Error getting csrf account id", zap.String("csrfToken", csrfToken), zap.Error(err))
		return err
	}

	if accountID != csrfAccountID {
		u.logger.Warn("CSRF token does not match session token", zap.String("sessionToken", sessionToken), zap.String("csrfToken", csrfToken))
		return fmt.Errorf("csrf token does not match session token")
	}

	return nil
}
