package todo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/parwin-pp/todo-application/internal/mock"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type mockGetTodosDatabase struct {
	Database

	NumberOfCalled int
	CallWithParams []string
	ReturnTodos    []model.Todo
	ReturnError    error
}

func (m *mockGetTodosDatabase) GetTodos(ctx context.Context, userID string) ([]model.Todo, error) {
	m.NumberOfCalled++
	m.CallWithParams = append(m.CallWithParams, userID)
	return m.ReturnTodos, m.ReturnError
}

type testGetTodosContext struct {
	t          *testing.T
	router     *bunrouter.Router
	db         *mockGetTodosDatabase
	withUserID string
}

func (testContext *testGetTodosContext) createTodo(userID uuid.UUID, name string) *model.Todo {
	todo := model.Todo{
		ID:        uuid.New(),
		Name:      name,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	testContext.db.ReturnTodos = append(testContext.db.ReturnTodos, todo)
	return &todo
}

func (testContext *testGetTodosContext) withGetTodosError(err error) {
	testContext.db.ReturnError = err
}

func (testContext *testGetTodosContext) requestWithUserID(userID uuid.UUID) *httptest.ResponseRecorder {
	testContext.withUserID = userID.String()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)

	err := testContext.router.ServeHTTPError(w, req)
	require.NoError(testContext.t, err)

	return w
}

func newTestGetTodosContext(t *testing.T) *testGetTodosContext {
	db := &mockGetTodosDatabase{}
	testCtx := &testGetTodosContext{t: t}

	router := bunrouter.New()
	group := router.Use(mock.NewAuthMiddleware(func() string {
		return testCtx.withUserID
	}))

	server := NewServer(db)
	group.GET("/todos", server.HandleGetTodos)

	testCtx.db = db
	testCtx.router = router
	return testCtx
}

func TestGetTodos(t *testing.T) {
	userID := uuid.New()

	t.Run("should return http status 200 when called", func(t *testing.T) {
		testCtx := newTestGetTodosContext(t)
		testCtx.createTodo(userID, "MOCK_TODO")

		res := testCtx.requestWithUserID(userID)

		require.Equal(t, 200, res.Result().StatusCode)
	})

	t.Run("should call get todos from database with userID given userID = 'd923a2c7-e013-4668-ba05-da59dfaab667'", func(t *testing.T) {
		userID := uuid.MustParse("d923a2c7-e013-4668-ba05-da59dfaab667")
		testCtx := newTestGetTodosContext(t)
		testCtx.createTodo(userID, "MOCK_TODO")

		res := testCtx.requestWithUserID(userID)

		require.Equal(t, 200, res.Result().StatusCode)
		require.Equal(t, 1, testCtx.db.NumberOfCalled)
		require.Equal(t, []string{"d923a2c7-e013-4668-ba05-da59dfaab667"}, testCtx.db.CallWithParams)
	})

	t.Run("should call get todos from database with userID given userID = '054ae3b4-42db-4568-a5df-99a62cb1b001'", func(t *testing.T) {
		userID := uuid.MustParse("054ae3b4-42db-4568-a5df-99a62cb1b001")
		testCtx := newTestGetTodosContext(t)
		testCtx.createTodo(userID, "MOCK_TODO")

		res := testCtx.requestWithUserID(userID)

		require.Equal(t, 200, res.Result().StatusCode)
		require.Equal(t, []string{"054ae3b4-42db-4568-a5df-99a62cb1b001"}, testCtx.db.CallWithParams)
	})

	t.Run("should return response with exists todos in database", func(t *testing.T) {
		testCtx := newTestGetTodosContext(t)
		testCtx.createTodo(userID, "MOCK_TODO")

		res := testCtx.requestWithUserID(userID)

		var todos []model.Todo
		err := json.NewDecoder(res.Body).Decode(&todos)
		require.NoError(t, err)
		require.Equal(t, 1, len(todos))
		require.Equal(t, "MOCK_TODO", todos[0].Name)
	})

	t.Run("should return http error with status = 500 when called db with error", func(t *testing.T) {
		testCtx := newTestGetTodosContext(t)
		testCtx.withGetTodosError(errors.New("SOMETHING_WENT_WRONG"))

		res := testCtx.requestWithUserID(userID)

		require.Equal(t, 500, res.Result().StatusCode)
	})
}
