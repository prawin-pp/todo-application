package todo

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type testGetTodosContext struct {
	t      *testing.T
	router *bunrouter.Router
	db     *mockTodoDatabase
}

func (testContext *testGetTodosContext) createTodo(userID uuid.UUID, name string) *model.Todo {
	return &model.Todo{
		ID:        uuid.New(),
		Name:      name,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (testContext *testGetTodosContext) sendRequest() *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)

	err := testContext.router.ServeHTTPError(w, req)
	require.NoError(testContext.t, err)

	return w
}

func newTestGetTodosContext(t *testing.T) *testGetTodosContext {
	router := bunrouter.New()
	mockDB := &mockTodoDatabase{}
	server := NewServer(mockDB)
	router.GET("/todos", server.HandleGetTodos)

	return &testGetTodosContext{
		router: router,
		t:      t,
		db:     mockDB,
	}
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
