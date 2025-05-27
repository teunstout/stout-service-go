package jishoclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"stout.dev/jisho/internal/domain"
)

type JishoClientInterface struct {
	logger domain.Logger
}

func NewJishoUsecase(l domain.Logger) *JishoClientInterface {
	return &JishoClientInterface{logger: l}
}

func (c *JishoClientInterface) SearchJisho(keyword string) (*JishoResponse, error) {
	uri := fmt.Sprintf("http://jisho.org/api/v1/search/words?keyword=%s", url.QueryEscape(keyword))
	jr, err := http.Get(uri)

	if err != nil {
		c.logger.Error("Failed to make request to Jisho API", map[string]interface{}{
			"uri":   uri,
			"error": err,
		})
		return nil, err
	}
	defer jr.Body.Close()

	response := &JishoResponse{}
	err = json.NewDecoder(jr.Body).Decode(response)

	if err != nil {
		c.logger.Warn("Failed to decode Jisho API response", map[string]interface{}{
			"uri":    uri,
			"error":  err,
			"status": jr.StatusCode,
		})
		return nil, err
	}

	return response, nil
}
