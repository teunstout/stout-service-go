package usecase

import (
	"crypto/rsa"
	"fmt"

	"stout.dev/login/internal/domain"
	"stout.dev/login/internal/repository"
)

type RegisterUsecaseInterface struct {
	accountRepo *repository.AccountRepositoryInterface
	logger      domain.Logger
	privateKey  *rsa.PrivateKey
}

func NewRegisterUsecase(accountRepo *repository.AccountRepositoryInterface, logger domain.Logger, privateKey *rsa.PrivateKey) *RegisterUsecaseInterface {
	return &RegisterUsecaseInterface{
		accountRepo: accountRepo,
		logger:      logger,
		privateKey:  privateKey,
	}
}

func (u *RegisterUsecaseInterface) Register(username string, password string) error {
	exists, err := u.accountRepo.AccountExists(username)

	if err != nil {
		u.logger.Error("Error checking if account exists", map[string]interface{}{"username": username, "error": err})
		return err
	}

	if exists {
		u.logger.Debug("Account already exists", map[string]interface{}{})
		return fmt.Errorf("account already exists")
	}

	// Hash password
	hashedPassword, err := domain.HashPassword(password)
	if err != nil {
		u.logger.Error("Error hashing password", map[string]interface{}{"username": username, "error": err})
		return err
	}

	if err = u.accountRepo.CreateAccount(username, hashedPassword); err != nil {
		u.logger.Error("Error creating account", map[string]interface{}{"username": username, "error": err})
		return err
	}
	return nil
}
