package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/auth"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type testAuthMiddlewareContext struct {
	t                 *testing.T
	router            *bunrouter.Router
	en                *auth.AuthEncryption
	UserIDFromContext string
}

func (ctx *testAuthMiddlewareContext) sendRequest(token *string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/", nil)
	if token != nil {
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: *token,
		})
	}

	res := httptest.NewRecorder()
	ctx.router.ServeHTTP(res, req)
	return res
}

func newTestAuthMiddlewareContext(t *testing.T) *testAuthMiddlewareContext {
	encrypter := auth.NewAuthEncryption("HS256", []byte("TEST_SECRET"), "1h")
	router := bunrouter.New(bunrouter.Use(NewAuthMiddleware(encrypter)))
	testCtx := &testAuthMiddlewareContext{
		t:                 t,
		router:            router,
		en:                encrypter,
		UserIDFromContext: "",
	}

	router.GET("/", func(w http.ResponseWriter, r bunrouter.Request) error {
		value := r.Context().Value(internal.AuthContextKey{})
		userID, ok := value.(string)
		if ok {
			testCtx.UserIDFromContext = userID
		}
		return nil
	})
	return testCtx
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("should set user id in context when called with valid token", func(t *testing.T) {
		testCtx := newTestAuthMiddlewareContext(t)
		token, err := testCtx.en.SignAuthToken("MOCK_USER_ID", map[string]interface{}{})

		testCtx.sendRequest(&token)

		require.NoError(t, err)
		require.Equal(t, "MOCK_USER_ID", testCtx.UserIDFromContext)
	})

	t.Run("should return http status 401 when not set token in cookie", func(t *testing.T) {
		testCtx := newTestAuthMiddlewareContext(t)

		res := testCtx.sendRequest(nil)

		require.Equal(t, 401, res.Result().StatusCode)
	})

	t.Run("should return http status 401 when set invalid token in cookie", func(t *testing.T) {
		testCtx := newTestAuthMiddlewareContext(t)
		token := "INVALID_TOKEN"

		res := testCtx.sendRequest(&token)

		require.Equal(t, 401, res.Result().StatusCode)
	})
}
