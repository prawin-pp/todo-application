package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/parwin-pp/todo-application/internal"
	"github.com/parwin-pp/todo-application/internal/mock"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type testAuthMiddlewareContext struct {
	t                 *testing.T
	router            *bunrouter.Router
	en                *mock.AuthEncryptor
	UserIDFromContext string
}

func (ctx *testAuthMiddlewareContext) sendRequest(cookieToken *string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/", nil)
	if cookieToken != nil {
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: *cookieToken,
		})
	}

	res := httptest.NewRecorder()
	ctx.router.ServeHTTP(res, req)
	return res
}

func newTestAuthMiddlewareContext(t *testing.T) *testAuthMiddlewareContext {
	encrypter := &mock.AuthEncryptor{}
	encrypter.VerifyAuthTokenFn = func(token string) (*jwt.Token, *jwt.MapClaims, error) {
		mock := jwt.New(jwt.SigningMethodHS256)
		mapClaims := mock.Claims.(jwt.MapClaims)
		return mock, &mapClaims, nil
	}

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
		testCtx.en.VerifyAuthTokenFn = func(token string) (*jwt.Token, *jwt.MapClaims, error) {
			mock := jwt.New(jwt.SigningMethodHS256)
			mapClaims := mock.Claims.(jwt.MapClaims)
			mapClaims["sub"] = "MOCK_USER_ID"
			return mock, &mapClaims, nil
		}
		token := "MOCK_VALID_TOKEN"

		testCtx.sendRequest(&token)

		require.Equal(t, "MOCK_USER_ID", testCtx.UserIDFromContext)
	})

	t.Run("should return http status 401 when not set token in cookie", func(t *testing.T) {
		testCtx := newTestAuthMiddlewareContext(t)

		res := testCtx.sendRequest(nil)

		require.Equal(t, 401, res.Result().StatusCode)
	})

	t.Run("should return http status 401 when set invalid token in cookie", func(t *testing.T) {
		testCtx := newTestAuthMiddlewareContext(t)
		testCtx.en.VerifyAuthTokenFn = func(token string) (*jwt.Token, *jwt.MapClaims, error) {
			return nil, nil, errors.New("INVALID_TOKEN")
		}

		token := "INVALID_TOKEN"

		res := testCtx.sendRequest(&token)

		require.Equal(t, 401, res.Result().StatusCode)
	})
}
