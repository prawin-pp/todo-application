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
	"github.com/parwin-pp/todo-application/internal/middleware"
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

func (testCtx *testGetTodosContext) createTodo(userID uuid.UUID, name string) *model.Todo {
	todo := model.Todo{
		ID:        uuid.New(),
		Name:      name,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	testCtx.db.ReturnTodos = append(testCtx.db.ReturnTodos, todo)
	return &todo
}

func (testCtx *testGetTodosContext) withGetTodosError(err error) {
	testCtx.db.ReturnError = err
}

func (testCtx *testGetTodosContext) requestWithUserID(userID uuid.UUID) *httptest.ResponseRecorder {
	testCtx.withUserID = userID.String()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)

	testCtx.router.ServeHTTP(w, req)
	return w
}

func newTestGetTodosContext(t *testing.T) *testGetTodosContext {
	db := &mockGetTodosDatabase{}
	testCtx := &testGetTodosContext{t: t}

	router := bunrouter.New(
		bunrouter.Use(middleware.NewErrorHandler),
		bunrouter.Use(mock.NewAuthMiddleware(func() string {
			return testCtx.withUserID
		})),
	)

	server := NewServer(db)
	router.GET("/todos", server.HandleGetTodos)

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

type testGetTodoContext struct {
	t          *testing.T
	router     *bunrouter.Router
	db         *mock.TodoDatabase
	withUserID string
}

func newTestGetTodoContext(t *testing.T) *testGetTodoContext {
	mockUserID := uuid.New()
	testCtx := &testGetTodoContext{t: t, withUserID: mockUserID.String()}

	db := &mock.TodoDatabase{}
	db.GetTodoFn = func(ctx context.Context, userID, todoID string) (*model.Todo, error) {
		return &model.Todo{
			ID:     uuid.MustParse(todoID),
			Name:   "MOCK_TODO",
			UserID: uuid.MustParse(userID),
		}, nil
	}

	router := bunrouter.New(
		bunrouter.Use(middleware.NewErrorHandler),
		bunrouter.Use(mock.NewAuthMiddleware(func() string {
			return testCtx.withUserID
		})),
	)
	server := NewServer(db)
	router.GET("/todos/:todoId", server.HandleGetTodo)

	testCtx.db = db
	testCtx.router = router
	return testCtx
}

func (testCtx *testGetTodoContext) request(todoID string, userID *string) *httptest.ResponseRecorder {
	if userID != nil {
		testCtx.withUserID = *userID
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/todos/"+todoID, nil)

	testCtx.router.ServeHTTP(w, req)
	return w
}

func TestGetTodo(t *testing.T) {
	t.Run("should return http status 200 when called", func(t *testing.T) {
		testCtx := newTestGetTodoContext(t)
		todoID := uuid.New().String()

		res := testCtx.request(todoID, nil)

		require.Equal(t, 200, res.Result().StatusCode)
	})

	t.Run("should return exists todo in database", func(t *testing.T) {
		testCtx := newTestGetTodoContext(t)
		todoID := uuid.New().String()

		res := testCtx.request(todoID, nil)

		var resBody model.Todo
		err := json.NewDecoder(res.Body).Decode(&resBody)
		require.NoError(t, err)
		require.Equal(t, todoID, resBody.ID.String())
	})

	t.Run("should return http status 500 when called db with error", func(t *testing.T) {
		testCtx := newTestGetTodoContext(t)
		todoID := uuid.New().String()
		testCtx.db.GetTodoFn = func(ctx context.Context, userID, todoID string) (*model.Todo, error) {
			return nil, errors.New("MOCK_ERROR")
		}

		res := testCtx.request(todoID, nil)

		require.Equal(t, 500, res.Result().StatusCode)
	})

	t.Run("should return http status 404 when get todo not found", func(t *testing.T) {
		testCtx := newTestGetTodoContext(t)
		todoID := uuid.New().String()
		testCtx.db.GetTodoFn = func(ctx context.Context, userID, todoID string) (*model.Todo, error) {
			return nil, nil
		}

		res := testCtx.request(todoID, nil)

		require.Equal(t, 404, res.Result().StatusCode)
	})
}
