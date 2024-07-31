package api

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidStatusCode = errors.New("invalid status code")
)

// GetRedirect returns the final URL after redirection
func GetRedirect(url string) (*http.Response, error) {
	// const op = "api.GetRedirect"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return resp, err
	}

	// if resp.StatusCode != http.StatusFound {
	// 	return "", fmt.Errorf("%s: %w: %d", op, ErrInvalidStatusCode, resp.StatusCode)
	// }

	// defer func() { _ = resp.Body.Close() }()

	return resp, nil
}
