package auth

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type testLoginContext struct {
	t      *testing.T
	router *bunrouter.Router
}

func (ctx *testLoginContext) sendRequest(body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	res := httptest.NewRecorder()

	ctx.router.ServeHTTP(res, req)

	return res
}

func newTestLoginContext(t *testing.T) *testLoginContext {
	server := NewServer(nil)

	router := bunrouter.New()
	router.POST("/login", server.HandleLogin)

	return &testLoginContext{
		t:      t,
		router: router,
	}
}

func TestLogin(t *testing.T) {
	t.Run("should return http status 200 when called", func(t *testing.T) {
		testCtx := newTestLoginContext(t)

		res := testCtx.sendRequest(`{ "username": "TEST_USERNAME", "password": "TEST_PASSWORD" }`)

		require.Equal(t, 200, res.Result().StatusCode)
	})
}
