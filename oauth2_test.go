package oauth2

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

const (
	baseURL = "http://localhost:8888"
)

func TestTokenFetch(t *testing.T) {
	//Given
	clientID := "testClientID"
	clientSecret := "testClientSecret"
	respFile, _ := ioutil.ReadFile("./test-fixtures/token_response.json")
	resp := tokenResponse{}
	req := tokenRequest{}
	if err := json.Unmarshal(respFile, &resp); err != nil {
		assert.Fail(t, "Can't unmarshal response file")
	}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", baseURL+tokenPath,
		func(r *http.Request) (*http.Response, error) {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return httpmock.NewStringResponse(400, ""), nil
			}
			resp, _ := httpmock.NewJsonResponse(200, resp)
			return resp, nil
		},
	)

	//when
	testConfig := Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		BaseURL:      baseURL}
	token, err := testConfig.TokenSource(context.Background(), testHc).Token()

	//then
	assert.Nil(t, err, "TokenSource should not return an error")
	assert.NotNil(t, token, "TokenSource should return a token")
	verifyTokenRequest(t, req, clientID, clientSecret)
	verifyToken(t, *token, resp)
}

func TestError(t *testing.T) {
	respFile, _ := ioutil.ReadFile("./test-fixtures/token_response_err.json")
	resp := tokenError{}
	if err := json.Unmarshal(respFile, &resp); err != nil {
		assert.Fail(t, "Can't unmarshal response file")
	}
	testHc := &http.Client{}
	httpmock.ActivateNonDefault(testHc)
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", baseURL+tokenPath,
		func(r *http.Request) (*http.Response, error) {
			resp, _ := httpmock.NewJsonResponse(500, resp)
			return resp, nil
		},
	)

	//when
	testConfig := Config{
		ClientID:     "clientID",
		ClientSecret: "clientSecret",
		BaseURL:      baseURL}
	_, err := testConfig.TokenSource(context.Background(), testHc).Token()

	//then
	assert.NotNil(t, err, "TokenSource should return an error")
	assert.IsType(t, Error{}, err, "Returned error has proper type")
	verifyErrorResponse(t, err.(Error), resp)
}

func verifyToken(t *testing.T, token oauth2.Token, resp tokenResponse) {
	assert.Equal(t, resp.AccessToken, token.AccessToken, "AccessToken matches")
	assert.Equal(t, resp.RefreshToken, token.RefreshToken, "RefreshToken matches")
	assert.Equal(t, resp.TokenType, token.TokenType, "TokenType matches")

	respTimeout, err := strconv.Atoi(resp.TokenTimeout)
	assert.Nil(t, err, "Error when converting TokenTimeout from the response to int: %v", err)
	assert.WithinDuration(t, time.Now().Add(time.Duration(respTimeout)*time.Second), token.Expiry, time.Duration(1)*time.Second, "Token expiry reflects token_timeout from the response")
}

func verifyTokenRequest(t *testing.T, req tokenRequest, clientID string, clientSecret string) {
	assert.Equal(t, "client_credentials", req.GrantType, "GrantType matches")
	assert.Equal(t, clientID, req.ClientID, "ClientID matches")
	assert.Equal(t, clientSecret, req.ClientSecret, "ClientSecret matches")
}

func verifyErrorResponse(t *testing.T, err Error, resp tokenError) {
	assert.Equal(t, resp.ErrorCode, err.Code, "Error code matches")
	assert.Equal(t, resp.ErrorMessage, err.Message, "Error message matches")
}
