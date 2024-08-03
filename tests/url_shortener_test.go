package test

import (
	"fmt"
	"net/http"
	"net/url"

	"path"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/random"

	gofakeit "github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   ":8083",
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/url").WithJSON(save.Request{
		URL:   gofakeit.URL(),
		Alias: random.NewRandomAlias(10),
	}).
		WithBasicAuth("user", "password").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("alias")
}

func TestURLShortener_SaveRedirectRemove(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},

		{
			name:  "Invalid URL",
			url:   "Invalid URL",
			alias: gofakeit.Word() + gofakeit.Word(),
			error: "field URL is not a valid URL",
		},

		{
			name:  "Empty alias",
			url:   gofakeit.URL(),
			alias: "",
		},

		{
			name:  "Empty url",
			url:   "",
			alias: gofakeit.Word(),
			error: "field URL is a required field",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   ":8083",
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/url").WithJSON(save.Request{
				URL:   tc.url,
				Alias: tc.alias,
			}).WithBasicAuth("user", "password").Expect().Status(http.StatusOK).JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)
				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			// Redirect

			ru := url.URL{
				Scheme: "http",
				Host:   ":8083",
				Path:   alias,
			}
			fmt.Println(ru.String())

			redirectedToURL, err := api.GetRedirect(ru.String())

			require.NoError(t, err)

			require.Equal(t, tc.url, redirectedToURL.Header.Get("Location"))

			// Remove

			reqDel := e.DELETE("/"+path.Join("url", alias)).WithBasicAuth("user", "password").Expect().Status(http.StatusOK).JSON()
			reqDel.String().IsEqual("url deleted")
		})
	}
}
