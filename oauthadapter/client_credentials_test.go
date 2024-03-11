package oauthadapter

import (
	"github.com/wiremock/go-wiremock"
	"gotest.tools/v3/assert"
	"net/http"
	"testing"
	"wiremock-demo/model"
	"wiremock-demo/openid"
	. "wiremock-demo/testutil"
)

func TestCCAdapter_GetToken(t *testing.T) {
	tokenResp200 := ReadFile("token-response-200.json")

	testCases := []struct {
		Name               string
		WireMockStubRule   *wiremock.StubRule
		ExpectedToken      *model.Token
		ExpectedError      bool
		ExpectedErrMessage string
	}{
		{
			Name: "Successful",
			WireMockStubRule: wiremock.Post(wiremock.URLPathEqualTo("/token")).
				WithHeader("Content-Type", wiremock.EqualTo("application/json")).
				WillReturnResponse(wiremock.NewResponse().WithStatus(200).
					WithBody(tokenResp200).
					WithHeader("Content-Type", "application/json"),
				),
			ExpectedToken:      &model.Token{AccessToken: "testToken", TokenType: "Bearer"},
			ExpectedError:      false,
			ExpectedErrMessage: "",
		},
		{
			Name: "Bad Request",
			WireMockStubRule: wiremock.Post(wiremock.URLPathEqualTo("/token")).
				WithHeader("Content-Type", wiremock.EqualTo("application/json")).
				WillReturnResponse(wiremock.NewResponse().WithStatus(400).
					WithBody(`{"error":"invalid_client","error_description":"Client authentication failed"}`).
					WithHeader("Content-Type", "application/json"),
				),
			ExpectedToken:      nil,
			ExpectedError:      true,
			ExpectedErrMessage: `client error: 400 - {"error":"invalid_client","error_description":"Client authentication failed"}`,
		},
		{
			Name: "Unauthorized 401",
			WireMockStubRule: wiremock.Post(wiremock.URLPathEqualTo("/token")).
				WithHeader("Content-Type", wiremock.EqualTo("application/json")).
				WillReturnResponse(wiremock.NewResponse().WithStatus(401).
					WithBody(`{"error":"unauthorized","error_description":"Invalid client credentials"}`).
					WithHeader("Content-Type", "application/json"),
				),
			ExpectedToken:      nil,
			ExpectedError:      true,
			ExpectedErrMessage: `client error: 401 - {"error":"unauthorized","error_description":"Invalid client credentials"}`,
		},
		{
			Name: "Server error 500",
			WireMockStubRule: wiremock.Post(wiremock.URLPathEqualTo("/token")).
				WithHeader("Content-Type", wiremock.EqualTo("application/json")).
				WillReturnResponse(wiremock.NewResponse().WithStatus(500).
					WithBody(`{"error":"internal_server_error"}`).
					WithHeader("Content-Type", "application/json"),
				),
			ExpectedToken:      nil,
			ExpectedError:      true,
			ExpectedErrMessage: `server error: 500 - {"error":"internal_server_error"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			wm := SetupWireMock(t)
			wm.Client.StubFor(tc.WireMockStubRule)

			openIdConfig := openid.GetOpenIdConfig(wm.Address)

			ccConfig := CCConfig{
				clientId:     "testClientId",
				clientSecret: "testClientSecret",
			}
			adapter := NewCCAdapter(ccConfig, openIdConfig)

			token, err := adapter.GetToken()

			if err != nil {
				assert.Equal(t, true, tc.ExpectedError)
				assert.Equal(t, err.Error(), tc.ExpectedErrMessage)
			} else {
				assert.Equal(t, token.AccessToken, tc.ExpectedToken.AccessToken)
				assert.Equal(t, token.TokenType, tc.ExpectedToken.TokenType)
			}

			wm.Client.Verify(wiremock.NewRequest(http.MethodGet, wiremock.URLPathEqualTo(openid.WellKnownEndpoint)), 1)
			wm.Client.Verify(wiremock.NewRequest(http.MethodPost, wiremock.URLPathEqualTo("/token")), 1)
		})
	}
}
