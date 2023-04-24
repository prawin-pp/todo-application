package todo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/parwin-pp/todo-application/internal/auth"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type testGetTodosContext struct {
	t          *testing.T
	router     *bunrouter.Router
	db         *mockTodoDatabase
	withUserID string
}

func (testContext *testGetTodosContext) authMiddleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		ctx := req.Context()
		ctx = context.WithValue(ctx, auth.AuthContextKey{}, testContext.withUserID)
		return next(w, req.WithContext(ctx))
	}
}
func (testContext *testGetTodosContext) WithUserID(userID uuid.UUID) {
	testContext.withUserID = userID.String()
}

func (testContext *testGetTodosContext) createTodo(userID uuid.UUID, name string) *model.Todo {
	todo := model.Todo{
		ID:        uuid.New(),
		Name:      name,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	testContext.db.ExistsTodos = append(testContext.db.ExistsTodos, todo)
	return &todo
}

func (testContext *testGetTodosContext) sendRequest() *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)

	err := testContext.router.ServeHTTPError(w, req)
	require.NoError(testContext.t, err)

	return w
}

func newTestGetTodosContext(t *testing.T) *testGetTodosContext {
	testCtx := &testGetTodosContext{t: t}
	testCtx.router = bunrouter.New(bunrouter.Use(testCtx.authMiddleware))
	testCtx.db = &mockTodoDatabase{}
	server := NewServer(testCtx.db)

	testCtx.router.GET("/todos", server.HandleGetTodos)

	return testCtx
}

type mockTodoDatabase struct {
	NumberOfCalled int
	CallWithParams []string
	ExistsTodos    []model.Todo
	ReturnError    error
}

func (m *mockTodoDatabase) GetAll(userID string) ([]model.Todo, error) {
	m.NumberOfCalled++
	m.CallWithParams = append(m.CallWithParams, userID)
	return m.ExistsTodos, m.ReturnError
}
