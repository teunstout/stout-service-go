package jishoclient

import (
	"encoding/json"
	"net/http"

	"stout.dev/jisho/internal/domain"
)

type JishoClientInterface struct {
	logger domain.Logger
}

func NewJishoUsecase(l domain.Logger) *JishoClientInterface {
	return &JishoClientInterface{logger: l}
}

func (c *JishoClientInterface) SearchJisho(keyword string) (*JishoResponse, error) {
	jr, err := http.Get("https://jisho.org/api/v1/search/words?keyword=" + keyword)
	if err == nil {
		c.logger.Error("Failed to make request to Jisho API", map[string]interface{}{"error": err})
		return nil, err
	}
	defer jr.Body.Close()

	var response JishoResponse
	err = json.NewDecoder(jr.Body).Decode(&response)
	if err != nil {
		c.logger.Warn("Failed to decode Jisho API response", map[string]interface{}{"error": err})
		return nil, err
	}
	return &response, nil
}
