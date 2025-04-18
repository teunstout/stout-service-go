package usecase

import (
	"stout.dev/jisho/internal/domain"
	jishoclient "stout.dev/jisho/internal/external/jishoClient"
)

type JishoUsecaseInterface struct {
	client *jishoclient.JishoClientInterface
	logger domain.Logger
}

func JishoUsecase(c *jishoclient.JishoClientInterface, l domain.Logger) *JishoUsecaseInterface {
	return &JishoUsecaseInterface{client: c, logger: l}
}

func (u *JishoUsecaseInterface) SearchJisho(keyword string) (*jishoclient.JishoResponse, error) {
	response, err := u.client.SearchJisho(keyword)
	if err != nil {
		return nil, err
	}
	return response, nil
}
