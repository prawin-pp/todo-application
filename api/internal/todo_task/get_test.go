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

type testGetTasksContext struct {
	t          *testing.T
	router     *bunrouter.Router
	db         *mockGetTasksDatabase
	withUserID string
}

func (testContext *testGetTasksContext) createTask(userID uuid.UUID, name string) *model.Todo {
	return nil
}

func (testContext *testGetTasksContext) requestWithUserID(userID uuid.UUID) *httptest.ResponseRecorder {
	testContext.withUserID = userID.String()
	mockTodoID := "ccee7451-2652-436e-b9a2-988f9c9f4d01"

	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s", mockTodoID)
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

	t.Run("should return http status 200 when called", func(t *testing.T) {
		testCtx := newTestGetTasksContext(t)

		res := testCtx.requestWithUserID(userID)

		require.Equal(t, 200, res.Result().StatusCode)
	})
}
