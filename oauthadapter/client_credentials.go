package oauthadapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	. "wiremock-demo/model"
)

type OAuthAdapter interface {
	GetToken() (Token, error)
}

type CCAdapter struct {
	config CCConfig
	openIdConfig OpenIdConfig
	client http.Client
}

type CCConfig struct {
	clientId     string
	clientSecret string
}

func NewCCAdapter(config CCConfig, openIdConfig OpenIdConfig) CCAdapter {
	return CCAdapter{
		config: config,
		openIdConfig: openIdConfig,
		client: http.Client{},
	}
}

func (c *CCAdapter) GetToken() (*Token, error) {
	req, _ := c.prepareTokenRequest()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, _ := io.ReadAll(resp.Body)

	err = handleHttpErrorCode(resp, string(respBody))
	if err != nil {
		return nil, err
	}

	var tokenResponse Token
	json.Unmarshal(respBody, &tokenResponse)

	defer resp.Body.Close()
	return &tokenResponse, nil
}

func (c *CCAdapter) prepareTokenRequest() (*http.Request, error) {
	path := c.openIdConfig.TokenEndpoint
	method := http.MethodPost
	body, _ := json.Marshal(c.config)

	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func handleHttpErrorCode(resp *http.Response, body string) error {
	switch {
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		return fmt.Errorf("client error: %d - %s", resp.StatusCode, body)
	case resp.StatusCode >= 500:
		return fmt.Errorf("server error: %d - %s", resp.StatusCode, body)
	default:
		return nil
	}
}
