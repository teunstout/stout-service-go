package usecase

import (
	"crypto/rsa"

	"go.uber.org/zap"
	"stout.dev/idp/internal/domain"
)

type JwksUsecaseInterface struct {
	logger    *zap.Logger
	publicKey *rsa.PublicKey
}

func JwksUsecase(logger *zap.Logger, publicKey *rsa.PublicKey) *JwksUsecaseInterface {
	return &JwksUsecaseInterface{
		logger:    logger,
		publicKey: publicKey,
	}
}

func (u *JwksUsecaseInterface) JwksKeys() domain.JwksKeys {
	keys := domain.CreateJwksKeys(u.publicKey)
	u.logger.Debug("JWKS keys created", zap.Any("JWKS", keys))
	return keys
}
