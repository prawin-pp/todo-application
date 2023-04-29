package auth

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/parwin-pp/todo-application/internal/config"
	"github.com/parwin-pp/todo-application/internal/middleware"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type mockGetUserDatabase struct {
	Database

	NumberOfCalled int
	CallWithParams []string
	ExistsUsers    []model.User
	ReturnError    error
}

func (m *mockGetUserDatabase) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	m.NumberOfCalled++
	m.CallWithParams = append(m.CallWithParams, username)
	for _, user := range m.ExistsUsers {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, m.ReturnError
}

type testLoginContext struct {
	t      *testing.T
	router *bunrouter.Router
	db     *mockGetUserDatabase
	en     Encrypter
}

func (ctx *testLoginContext) createUser(username string, password string) *model.User {
	hased, err := ctx.en.Hash(password)
	require.NoError(ctx.t, err)

	user := model.User{
		Username: username,
		Password: hased,
	}
	ctx.db.ExistsUsers = append(ctx.db.ExistsUsers, user)
	return &user
}

func (ctx *testLoginContext) sendRequest(body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	res := httptest.NewRecorder()

	ctx.router.ServeHTTP(res, req)

	return res
}

func newTestLoginContext(t *testing.T) *testLoginContext {
	db := &mockGetUserDatabase{}
	encrypter := NewAuthEncryption("HS256", []byte("TEST_SECRET"), time.Hour)
	conf := config.AuthConfig{ExpireDuration: time.Hour}
	server := NewServer(db, encrypter, conf)
	router := bunrouter.New(bunrouter.Use(middleware.NewErrorHandler))
	router.POST("/login", server.HandleLogin)

	return &testLoginContext{
		t:      t,
		router: router,
		db:     db,
		en:     encrypter,
	}
}

func TestLogin(t *testing.T) {
	t.Run("should return http status 200 when called", func(t *testing.T) {
		testCtx := newTestLoginContext(t)
		testCtx.createUser("TEST_USERNAME", "TEST_PASSWORD")

		res := testCtx.sendRequest(`{ "username": "TEST_USERNAME", "password": "TEST_PASSWORD" }`)

		require.Equal(t, 200, res.Result().StatusCode)
	})

	t.Run("should return http status 401 when called with invalid username or password", func(t *testing.T) {
		testCtx := newTestLoginContext(t)
		testCtx.createUser("TEST_USERNAME", "TEST_PASSWORD")

		res := testCtx.sendRequest(`{ "username": "TEST_INVALID", "password": "TEST_INVALID" }`)

		require.Equal(t, 401, res.Result().StatusCode)
	})

	t.Run("should return http status 400 when called with invalid json", func(t *testing.T) {
		testCtx := newTestLoginContext(t)

		res := testCtx.sendRequest(`{ "username":, "password": }`)

		require.Equal(t, 400, res.Result().StatusCode)
	})

	t.Run("should call get user by username when called", func(t *testing.T) {
		testCtx := newTestLoginContext(t)
		testCtx.createUser("TEST_USERNAME", "TEST_PASSWORD")

		testCtx.sendRequest(`{ "username": "TEST_USERNAME", "password": "TEST_PASSWORD" }`)

		require.Equal(t, 1, testCtx.db.NumberOfCalled)
		require.Equal(t, "TEST_USERNAME", testCtx.db.CallWithParams[0])
	})

	t.Run("should return status 401 when get user by username not found", func(t *testing.T) {
		testCtx := newTestLoginContext(t)
		testCtx.createUser("TEST_OTHER", "TEST_PASSWORD")

		res := testCtx.sendRequest(`{ "username": "TEST_USERNAME", "password": "TEST_PASSWORD" }`)

		require.Equal(t, 401, res.Result().StatusCode)
	})

	t.Run("should return status 500 when get user by username return error", func(t *testing.T) {
		testCtx := newTestLoginContext(t)
		testCtx.db.ReturnError = errors.New("MOCK_ERROR")

		res := testCtx.sendRequest(`{ "username": "TEST_USERNAME", "password": "TEST_PASSWORD" }`)

		require.Equal(t, 500, res.Result().StatusCode)
	})

	t.Run("should return status 401 when compare req password with hashed password in db not equals", func(t *testing.T) {
		testCtx := newTestLoginContext(t)
		testCtx.createUser("TEST_USERNAME", "TEST_PASSWORD")

		res := testCtx.sendRequest(`{ "username": "TEST_USERNAME", "password": "TEST_NOT_EQUALS" }`)

		require.Equal(t, 401, res.Result().StatusCode)
	})

	t.Run("should set response header with key='Set-Cookie', value='token={jwt}' when called", func(t *testing.T) {
		testCtx := newTestLoginContext(t)
		testCtx.createUser("TEST_USERNAME", "TEST_PASSWORD")

		res := testCtx.sendRequest(`{ "username": "TEST_USERNAME", "password": "TEST_PASSWORD" }`)

		cookie := res.Result().Cookies()[0]
		require.Equal(t, "token", cookie.Name)
		require.Equal(t, true, cookie.Secure)
		require.Equal(t, true, cookie.HttpOnly)
		require.Equal(t, "/", cookie.Path)

		_, _, err := testCtx.en.VerifyAuthToken(cookie.Value)
		require.NoError(t, err)
	})
}

type testLogoutContext struct {
	t      *testing.T
	router *bunrouter.Router
}

func newTestLogoutContext(t *testing.T) *testLogoutContext {
	db := &mockGetUserDatabase{}
	encrypter := NewAuthEncryption("HS256", []byte("TEST_SECRET"), time.Hour)
	conf := config.AuthConfig{ExpireDuration: time.Hour}
	server := NewServer(db, encrypter, conf)
	router := bunrouter.New(bunrouter.Use(middleware.NewErrorHandler))
	router.POST("/logout", server.HandleLogout)

	return &testLogoutContext{
		t:      t,
		router: router,
	}
}

func (ctx *testLogoutContext) sendRequest() *httptest.ResponseRecorder {
	cookie := &http.Cookie{
		Name:  "token",
		Value: "SOME_TEST_TOKEN",
	}
	req := httptest.NewRequest("POST", "/logout", nil)
	req.AddCookie(cookie)

	res := httptest.NewRecorder()
	ctx.router.ServeHTTP(res, req)

	return res
}

func TestLogout(t *testing.T) {
	t.Run("should return status 200 when called", func(t *testing.T) {
		testCtx := newTestLogoutContext(t)

		res := testCtx.sendRequest()

		require.Equal(t, 200, res.Result().StatusCode)
	})

	t.Run("should return response with new token cookie with empty value when called", func(t *testing.T) {
		testCtx := newTestLogoutContext(t)

		res := testCtx.sendRequest()

		cookie := res.Result().Cookies()[0]
		require.Equal(t, "token", cookie.Name)
		require.Equal(t, "", cookie.Value)
	})
}
