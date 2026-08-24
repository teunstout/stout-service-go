package usecase

import (
	"crypto/rsa"

	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
	"stout.dev/idp/internal/repository"
)

type RegisterUsecaseInterface struct {
	accountRepo repository.AccountRepository
	logger      *zap.Logger
	privateKey  *rsa.PrivateKey
}

func NewRegisterUsecase(accountRepo repository.AccountRepository, logger *zap.Logger, privateKey *rsa.PrivateKey) *RegisterUsecaseInterface {
	return &RegisterUsecaseInterface{
		accountRepo: accountRepo,
		logger:      logger,
		privateKey:  privateKey,
	}
}

func (u *RegisterUsecaseInterface) Register(username string, password string) error {
	exists, err := u.accountRepo.AccountExists(username)

	if err != nil {
		u.logger.Error("Error checking if account exists", zap.String("username", username), zap.Error(err))
		return err
	}

	if exists {
		u.logger.Debug("Account already exists", zap.String("username", username))
		return domain.ErrAccountAlreadyExists
	}

	hashedPassword, err := domain.HashPassword(password)
	if err != nil {
		u.logger.Error("Error hashing password", zap.String("username", username), zap.Error(err))
		return err
	}

	if err = u.accountRepo.CreateAccount(username, hashedPassword); err != nil {
		u.logger.Error("Error creating account", zap.String("username", username), zap.Error(err))
		return err
	}

	return nil
}
