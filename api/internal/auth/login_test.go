package auth

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

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
}

func (ctx *testLoginContext) createUser(username string, password string) *model.User {
	user := model.User{
		Username: username,
		Password: password,
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
	server := NewServer(db)
	router := bunrouter.New()
	router.POST("/login", server.HandleLogin)

	return &testLoginContext{
		t:      t,
		router: router,
		db:     db,
	}
}

func TestLogin(t *testing.T) {
	t.Run("should return http status 200 when called", func(t *testing.T) {
		testCtx := newTestLoginContext(t)

		res := testCtx.sendRequest(`{ "username": "TEST_USERNAME", "password": "TEST_PASSWORD" }`)

		require.Equal(t, 200, res.Result().StatusCode)
	})

	t.Run("should return http status 401 when called with invalid username or password", func(t *testing.T) {
		testCtx := newTestLoginContext(t)

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
}
