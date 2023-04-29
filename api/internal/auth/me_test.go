package auth

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/parwin-pp/todo-application/internal/config"
	"github.com/parwin-pp/todo-application/internal/middleware"
	"github.com/parwin-pp/todo-application/internal/mock"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type testGetMeContext struct {
	router     *bunrouter.Router
	db         *mock.AuthDatabase
	en         *mock.AuthEncryptor
	withUserID uuid.UUID
}

func (testCtx *testGetMeContext) request() *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	testCtx.router.ServeHTTP(w, req)
	return w
}

func newTestGetMeContext(t *testing.T) *testGetMeContext {
	mockUserID := uuid.New()

	db := &mock.AuthDatabase{}
	db.GetUserFn = func(ctx context.Context, userID string) (*model.User, error) {
		return &model.User{
			ID:       mockUserID,
			Username: "MOCK_USERNAME",
			Password: "MOCK_PASSWORD",
		}, nil
	}

	en := &mock.AuthEncryptor{}

	server := NewServer(db, en, config.AuthConfig{})

	router := bunrouter.New(
		bunrouter.Use(middleware.NewErrorHandler),
		bunrouter.Use(mock.NewAuthMiddleware(func() string {
			return mockUserID.String()
		})),
	)
	router.GET("/me", server.HandleGetMe)

	return &testGetMeContext{
		router:     router,
		db:         db,
		en:         en,
		withUserID: mockUserID,
	}
}

func TestGetMe(t *testing.T) {
	t.Run("should return http status 200", func(t *testing.T) {
		testCtx := newTestGetMeContext(t)

		res := testCtx.request()

		require.Equal(t, 200, res.Result().StatusCode)
	})

	t.Run("should return response body with user detail", func(t *testing.T) {
		testCtx := newTestGetMeContext(t)
		testCtx.db.GetUserFn = func(ctx context.Context, userID string) (*model.User, error) {
			return &model.User{
				ID:       testCtx.withUserID,
				Username: "test",
				Password: "test",
			}, nil
		}

		res := testCtx.request()

		resBody, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		require.JSONEq(t, fmt.Sprintf(`{
			"id": "%s",
			"username": "test"
		}`, testCtx.withUserID.String()), string(resBody))
	})

	t.Run("should return http status 500 when get user from db return error", func(t *testing.T) {
		testCtx := newTestGetMeContext(t)
		testCtx.db.GetUserFn = func(ctx context.Context, userID string) (*model.User, error) {
			return nil, fmt.Errorf("MOCK_ERROR")
		}

		res := testCtx.request()

		require.Equal(t, 500, res.Result().StatusCode)
	})

	t.Run("should return http status 401 when get user from db return not found", func(t *testing.T) {
		testCtx := newTestGetMeContext(t)
		testCtx.db.GetUserFn = func(ctx context.Context, userID string) (*model.User, error) {
			return nil, nil
		}

		res := testCtx.request()

		require.Equal(t, 401, res.Result().StatusCode)
	})

	t.Run("should not return user's password in response body", func(t *testing.T) {
		testCtx := newTestGetMeContext(t)

		res := testCtx.request()

		resBody, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		require.NotContains(t, string(resBody), "MOCK_PASSWORD")
	})
}
