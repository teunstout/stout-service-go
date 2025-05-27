package usecase

import (
	"fmt"

	"stout.dev/login/internal/domain"
	"stout.dev/login/internal/repository"
)

type LogoutUsecaseInterface struct {
	sessionRepo *repository.SessionTokenRepositoryInterface
	csrfRepo    *repository.CsrfRepositoryInterface
	logger      domain.Logger
}

func NewLogoutUsecase(sessionRepo *repository.SessionTokenRepositoryInterface, csrfRepo *repository.CsrfRepositoryInterface, logger domain.Logger) *LogoutUsecaseInterface {
	return &LogoutUsecaseInterface{
		sessionRepo: sessionRepo,
		csrfRepo:    csrfRepo,
		logger:      logger,
	}
}

func (u *LogoutUsecaseInterface) Logout(sessionToken string, csrfToken string, logoutAll bool) error {
	accountID, err := u.sessionRepo.GetAccountIdBySessionToken(sessionToken)
	if err != nil {
		u.logger.Info("Error getting account", map[string]interface{}{
			"aid":     accountID,
			"session": sessionToken,
			"csrf":    csrfToken,
		})
		return err
	}

	csrfAccountID, err := u.csrfRepo.GetCsrfIdBySessionToken(csrfToken)
	if err != nil {
		u.logger.Info("Error getting csrf account", map[string]interface{}{
			"aid":     accountID,
			"session": sessionToken,
			"csrf":    csrfToken,
		})
		return err
	}

	if accountID != csrfAccountID {
		u.logger.Warn("Account session and CSRF are not the same", map[string]interface{}{
			"aid":     accountID,
			"session": sessionToken,
			"csrf":    csrfToken,
		})
		return fmt.Errorf("Account Id %d didn't match CSRF %s", accountID, csrfToken)
	}

	if logoutAll {
		u.deleteAllSessions(accountID)
	} else {
		u.deleteSingleSession(csrfToken, sessionToken)
	}

	return nil
}

func (r *LogoutUsecaseInterface) deleteAllSessions(accountID int32) error {
	if err := r.csrfRepo.DeleteCsrfTokensByAccountId(accountID); err != nil {
		r.logger.Info("Error deleting csrf tokens", map[string]interface{}{"aid": accountID, "error": err})
		return err
	}

	if err := r.sessionRepo.DeleteSessionTokensByAccountId(accountID); err != nil {
		r.logger.Info("Error deleting session tokens", map[string]interface{}{"aid": accountID, "error": err})
		return err
	}

	return nil
}

func (r *LogoutUsecaseInterface) deleteSingleSession(csrfToken string, sessionToken string) error {
	if err := r.csrfRepo.DeleteCsrfToken(csrfToken); err != nil {
		r.logger.Info("Error deleting csrf token", map[string]interface{}{"csrf": csrfToken, "error": err})
		return err
	}

	if err := r.sessionRepo.DeleteSessionToken(sessionToken); err != nil {
		r.logger.Info("Error deleting session tokens", map[string]interface{}{"session": sessionToken, "error": err})
		return err
	}

	return nil
}
