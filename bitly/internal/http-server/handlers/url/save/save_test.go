package save_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leedinh/telebot/bitly/internal/http-server/handlers/url/save"
	"github.com/leedinh/telebot/bitly/internal/http-server/handlers/url/save/mocks"
	"github.com/leedinh/telebot/bitly/internal/lib/bloomfilter"
	"github.com/leedinh/telebot/bitly/internal/lib/logger/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		url       string
		respError string
		mockError error
	}{
		{
			name: "Success",
			url:  "https://google.com",
		},
		{
			name: "Empty alias",
			url:  "https://google.com",
		},
		{
			name:      "Empty URL",
			url:       "",
			respError: "field URL is a required field",
		},
		{
			name:      "Invalid URL",
			url:       "some invalid URL",
			respError: "field URL is not a valid URL",
		},
		{
			name:      "SaveURL Error",
			url:       "https://google.com",
			respError: "failed to add url",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewURLSaver(t)
			bf := bloomfilter.NewBloomFilter(1000, 5)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), bf, urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s"}`, tc.url)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

			// TODO: add more checks
		})
	}
}
