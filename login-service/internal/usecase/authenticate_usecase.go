package usecase

import (
	"fmt"

	"stout.dev/login/internal/domain"
	"stout.dev/login/internal/repository"
)

type AuthenticationUsecaseInterface struct {
	sessionRepo *repository.SessionTokenRepositoryInterface
	csrfRepo    *repository.CsrfRepositoryInterface
	logger      domain.Logger
}

func NewAuthenticationUsecase(sessionRepo *repository.SessionTokenRepositoryInterface, csrfRepo *repository.CsrfRepositoryInterface, logger domain.Logger) *AuthenticationUsecaseInterface {
	return &AuthenticationUsecaseInterface{
		sessionRepo: sessionRepo,
		csrfRepo:    csrfRepo,
		logger:      logger,
	}
}

func (u *AuthenticationUsecaseInterface) Register(sessionToken string, csrfToken string) error {
	accountID, err := u.sessionRepo.GetAccountIdBySessionToken(sessionToken)
	if err != nil {
		u.logger.Error("Error getting account id", map[string]interface{}{"sessionToken": sessionToken, "error": err})
		return err
	}

	csrfAccountID, err := u.csrfRepo.GetCsrfIdBySessionToken(csrfToken)
	if err != nil {
		u.logger.Error("Error getting csrf account id", map[string]interface{}{"csrfToken": csrfToken, "error": err})
		return err
	}

	if accountID != csrfAccountID {
		u.logger.Warn("CSRF token does not match session token", map[string]interface{}{
			"sessionToken": sessionToken,
			"csrfToken":    csrfToken,
		})
		return fmt.Errorf("csrf token does not match session token")
	}

	return nil
}
