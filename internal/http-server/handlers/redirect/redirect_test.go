package redirect_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"testing"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/redirect/mocks"
	"url-shortener/internal/http-server/handlers/slogdiscard"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com/",
		},

		{
			name:      "Not exists",
			alias:     "test_alias",
			mockError: storage.ErrURLNotFound,
			respError: "not found",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}

			r := chi.NewRouter()

			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			resp, err := api.GetRedirect(ts.URL + "/" + tc.alias)

			if tc.respError != "" {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				var response response.Response
				err = json.Unmarshal(body, &response)
				require.NoError(t, err)

				require.Equal(t, response.Error, tc.respError)
				return
			}

			require.NoError(t, err)

			require.True(t, resp.StatusCode == http.StatusFound, "invalid status code")

			require.Equal(t, tc.url, resp.Header.Get("Location"))
		})
	}
}
