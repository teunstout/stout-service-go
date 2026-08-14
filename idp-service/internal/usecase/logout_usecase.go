package usecase

import (
	"fmt"

	"go.uber.org/zap"
	"stout.dev/idp/internal/repository"
)

type LogoutUsecaseInterface struct {
	sessionRepo repository.SessionTokenRepository
	csrfRepo    repository.CsrfRepository
	logger      *zap.Logger
}

func NewLogoutUsecase(sessionRepo repository.SessionTokenRepository, csrfRepo repository.CsrfRepository, logger *zap.Logger) *LogoutUsecaseInterface {
	return &LogoutUsecaseInterface{
		sessionRepo: sessionRepo,
		csrfRepo:    csrfRepo,
		logger:      logger,
	}
}

func (u *LogoutUsecaseInterface) Logout(sessionToken string, csrfToken string, logoutAll bool) error {
	accountID, err := u.sessionRepo.GetAccountIdBySessionToken(sessionToken)
	if err != nil {
		u.logger.Info("Error getting account", zap.Int32("aid", accountID), zap.String("session", sessionToken), zap.String("csrf", csrfToken))
		return err
	}

	csrfAccountID, err := u.csrfRepo.GetCsrfIdBySessionToken(csrfToken)
	if err != nil {
		u.logger.Info("Error getting csrf account", zap.Int32("aid", accountID), zap.String("session", sessionToken), zap.String("csrf", csrfToken))
		return err
	}

	if accountID != csrfAccountID {
		u.logger.Warn("Account session and CSRF are not the same", zap.Int32("aid", accountID), zap.String("session", sessionToken), zap.String("csrf", csrfToken))
		return fmt.Errorf("Account Id %d didn't match CSRF %s", accountID, csrfToken)
	}

	if logoutAll {
		u.deleteAllSessions(accountID)
	} else {
		u.deleteSingleSession(csrfToken, sessionToken)
	}

	return nil
}

func (u *LogoutUsecaseInterface) deleteAllSessions(accountID int32) error {
	if err := u.csrfRepo.DeleteCsrfTokensByAccountId(accountID); err != nil {
		u.logger.Info("Error deleting csrf tokens", zap.Int32("aid", accountID), zap.Error(err))
		return err
	}

	if err := u.sessionRepo.DeleteSessionTokensByAccountId(accountID); err != nil {
		u.logger.Info("Error deleting session tokens", zap.Int32("aid", accountID), zap.Error(err))
		return err
	}

	return nil
}

func (u *LogoutUsecaseInterface) deleteSingleSession(csrfToken string, sessionToken string) error {
	if err := u.csrfRepo.DeleteCsrfToken(csrfToken); err != nil {
		u.logger.Info("Error deleting csrf token", zap.String("csrf", csrfToken), zap.Error(err))
		return err
	}

	if err := u.sessionRepo.DeleteSessionToken(sessionToken); err != nil {
		u.logger.Info("Error deleting session tokens", zap.String("session", sessionToken), zap.Error(err))
		return err
	}

	return nil
}
