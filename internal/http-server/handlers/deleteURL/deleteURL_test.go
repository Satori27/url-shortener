package deleteURL_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/deleteURL"
	"url-shortener/internal/http-server/handlers/deleteURL/mocks"
	"url-shortener/internal/http-server/handlers/slogdiscard"
	"url-shortener/internal/storage"


	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		resp      string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			resp:  "url deleted",
		},

		{
			name:      "Not exists",
			alias:     "test_alias",
			mockError: storage.ErrURLNotFound,
			resp:      "url not exists",
		},

		{
			name:      "Internal error",
			alias:     "test_alias",
			mockError: errors.New("Some error"),
			resp:      "internal error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlDeleterMock := mocks.NewURLDeleter(t)
			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).Once()
			}

			req, err := http.NewRequest("DELETE", "/"+tc.alias, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r := chi.NewRouter()
			r.Delete("/{alias}", deleteURL.New(slogdiscard.NewDiscardLogger(), urlDeleterMock))
			r.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp string

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.resp, resp)

		})
	}
}
