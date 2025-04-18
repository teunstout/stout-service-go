package usecase

import (
	"context"

	"stout.dev/jisho/internal/domain"
	jishoclient "stout.dev/jisho/internal/external/jishoClient"
	"stout.dev/jisho/internal/repository"
)

type JishoUsecaseInterface struct {
	client     *jishoclient.JishoClientInterface
	repository *repository.JishoRepositoryInterface
	logger     domain.Logger
}

func NewJishoUsecase(c *jishoclient.JishoClientInterface, r *repository.JishoRepositoryInterface, l domain.Logger) *JishoUsecaseInterface {
	return &JishoUsecaseInterface{client: c, logger: l}
}

func (u *JishoUsecaseInterface) SearchJisho(keyword string, ctx context.Context) (*jishoclient.JishoResponse, error) {
	response, err := u.client.SearchJisho(keyword)

	u.saveSearchHistory(keyword, ctx)

	if err != nil {
		return nil, err
	}
	return response, nil
}

func (u *JishoUsecaseInterface) saveSearchHistory(keyword string, ctx context.Context) {
	mid, ok := ctx.Value("mid").(int32)
	if ok {
		u.logger.Debug("No Member ID found in context, skipping search history save", nil)
		u.repository.SaveSearchHistory(mid, keyword)
	}
}
