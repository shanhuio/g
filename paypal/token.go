package paypal

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"shanhu.io/pub/errcode"
)

func makeTokenRequest(host, id, secret string) *http.Request {
	u := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/v1/oauth2/token",
		User:   url.UserPassword(id, secret),
	}

	form := make(url.Values)
	form.Set("grant_type", "client_credentials")

	header := make(http.Header)
	header.Add("Accept", "application/json")
	header.Add("Accept-Language", "eu_US")

	return &http.Request{
		Method: http.MethodPost,
		URL:    u,
		Header: header,
		Body:   io.NopCloser(strings.NewReader(form.Encode())),
	}
}

func tokenFromResponse(r io.Reader) (string, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return "", errcode.Annotate(err, "read token response")
	}

	var dat struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &dat); err != nil {
		return "", errcode.Annotate(err, "parse token from response")
	}

	return dat.AccessToken, nil
}
