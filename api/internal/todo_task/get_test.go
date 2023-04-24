package todotask

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type mockGetTasksDatabase struct {
	Database

	NumberOfCalled int
	CallWithParams [][]string
	ReturnTasks    []model.TodoTask
	ReturnError    error
}

func (m *mockGetTasksDatabase) GetAll(userID, todoID string) ([]model.TodoTask, error) {
	m.NumberOfCalled++
	m.CallWithParams = append(m.CallWithParams, []string{userID, todoID})
	return nil, nil
}

type testGetTasksContext struct {
	t          *testing.T
	router     *bunrouter.Router
	db         *mockGetTasksDatabase
	withUserID string
}

func (testContext *testGetTasksContext) createTask(userID uuid.UUID, name string) *model.Todo {
	return nil
}

func (testContext *testGetTasksContext) requestWithUserID(userID uuid.UUID, todoID uuid.UUID) *httptest.ResponseRecorder {
	testContext.withUserID = userID.String()

	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s", todoID)
	req := httptest.NewRequest(http.MethodGet, path, nil)

	err := testContext.router.ServeHTTPError(w, req)
	require.NoError(testContext.t, err)

	return w
}

func newTestGetTasksContext(t *testing.T) *testGetTasksContext {
	testCtx := &testGetTasksContext{t: t}
	testCtx.router = bunrouter.New(bunrouter.Use(mockAuthMiddleware(func() string {
		return testCtx.withUserID
	})))
	testCtx.db = &mockGetTasksDatabase{}

	server := NewServer(testCtx.db)
	testCtx.router.GET("/todos/:todoId", server.HandleGetTasks)

	return testCtx
}

func TestGetTodoTasks(t *testing.T) {
	userID := uuid.New()
	todoID := uuid.New()

	t.Run("should return http status 200 when called", func(t *testing.T) {
		testCtx := newTestGetTasksContext(t)

		res := testCtx.requestWithUserID(userID, todoID)

		require.Equal(t, 200, res.Result().StatusCode)
	})

	t.Run("should call get tasks from database when called", func(t *testing.T) {
		testCtx := newTestGetTasksContext(t)

		res := testCtx.requestWithUserID(userID, todoID)

		require.Equal(t, 200, res.Result().StatusCode)
		require.Equal(t, 1, testCtx.db.NumberOfCalled)
		require.Equal(t, []string{userID.String(), todoID.String()}, testCtx.db.CallWithParams[0])
	})
}
