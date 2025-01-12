package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/defer-panic/url-shortener-api/internal/model"
	"github.com/defer-panic/url-shortener-api/internal/server"
	"github.com/defer-panic/url-shortener-api/internal/shorten"
	"github.com/defer-panic/url-shortener-api/internal/storage/shortening"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleShorten(t *testing.T) {
	t.Run("returns shortened URL for a given URL", func(t *testing.T) {
		const payload = `{"url": "https://www.google.com"}`

		var (
			shortener = shorten.NewService(shortening.NewInMemory())
			handler   = server.HandleShorten(shortener)
			recorder  = httptest.NewRecorder()
			request   = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
			e         = echo.New()
			c         = e.NewContext(request, recorder)
		)

		addUserToCtx(c)

		e.Validator = server.NewValidator()
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		require.NoError(t, handler(c))
		assert.Equal(t, http.StatusOK, recorder.Code)

		var resp struct {
			ShortURL string `json:"shortUrl"`
		}

		require.NoError(t, json.NewDecoder(recorder.Body).Decode(&resp), &resp)
		assert.NotEmpty(t, resp.ShortURL)
	})

	t.Run("returns error if URL is invalid", func(t *testing.T) {
		const payload = `{"url": "invalid"}`

		var (
			shortener = shorten.NewService(shortening.NewInMemory())
			handler   = server.HandleShorten(shortener)
			recorder  = httptest.NewRecorder()
			request   = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
			e         = echo.New()
			c         = e.NewContext(request, recorder)
		)

		addUserToCtx(c)

		e.Validator = server.NewValidator()
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		var httpErr *echo.HTTPError
		require.ErrorAs(t, handler(c), &httpErr)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
		assert.Contains(t, httpErr.Message, "invalid url")
	})

	t.Run("returns error if identifier is already taken", func(t *testing.T) {
		const payload = `{"url": "https://www.google.com", "identifier": "google"}`

		var (
			shortener = shorten.NewService(shortening.NewInMemory())
			handler   = server.HandleShorten(shortener)
			recorder  = httptest.NewRecorder()
			request   = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
			e         = echo.New()
			c         = e.NewContext(request, recorder)
		)

		addUserToCtx(c)

		e.Validator = server.NewValidator()
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		require.NoError(t, handler(c))
		assert.Equal(t, http.StatusOK, recorder.Code)

		request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(payload))
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recorder = httptest.NewRecorder()
		c = e.NewContext(request, recorder)

		addUserToCtx(c)

		var httpErr *echo.HTTPError
		require.ErrorAs(t, handler(c), &httpErr)
		assert.Equal(t, http.StatusConflict, httpErr.Code)
		assert.Equal(t, httpErr.Message, model.ErrIdentifierExists.Error())
	})
}

func addUserToCtx(c echo.Context) {
	c.Set(
		"user",
		&jwt.Token{
			Claims: &model.UserClaims{
				User: model.User{
					GitHubLogin: gofakeit.Username(),
				},
			},
		},
	)
}
