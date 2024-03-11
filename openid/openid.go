package openid

import (
	"encoding/json"
	"io"
	"net/http"
	. "wiremock-demo/model"
)

const WellKnownEndpoint = "/.well-known/openid-configuration"

func GetOpenIdConfig(baseUrl string) OpenIdConfig {
	req, _ := prepareRequest(baseUrl)

	client := http.Client{}
	resp, _ := client.Do(req)
	responseBody, _ := io.ReadAll(resp.Body)

	var openId OpenIdConfig
	json.Unmarshal(responseBody, &openId)

	defer resp.Body.Close()
	return openId
}

func prepareRequest(baseUrl string) (*http.Request, error) {
	path := baseUrl + WellKnownEndpoint
	method := http.MethodGet
	return http.NewRequest(method, path, nil)
}
