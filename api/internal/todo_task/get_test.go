package todotask

import (
	"context"
	"encoding/json"
	"errors"
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

func (m *mockGetTasksDatabase) GetTasks(ctx context.Context, userID, todoID string) ([]model.TodoTask, error) {
	m.NumberOfCalled++
	m.CallWithParams = append(m.CallWithParams, []string{userID, todoID})
	return m.ReturnTasks, m.ReturnError
}

type testGetTasksContext struct {
	t          *testing.T
	router     *bunrouter.Router
	db         *mockGetTasksDatabase
	withUserID string
}

func (testContext *testGetTasksContext) createTask(userID, todoID uuid.UUID, name string) *model.TodoTask {
	task := model.TodoTask{
		ID:     uuid.New(),
		Name:   name,
		TodoID: todoID,
		UserID: userID,
	}
	testContext.db.ReturnTasks = append(testContext.db.ReturnTasks, task)
	return &task
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

	t.Run("should return response body with exists tasks from database", func(t *testing.T) {
		testCtx := newTestGetTasksContext(t)
		testCtx.createTask(userID, todoID, "MOCK_TASK_NAME")

		res := testCtx.requestWithUserID(userID, todoID)

		var tasks []model.TodoTask
		err := json.NewDecoder(res.Body).Decode(&tasks)
		require.NoError(t, err)
		require.Equal(t, 1, len(tasks))
		require.Equal(t, "MOCK_TASK_NAME", tasks[0].Name)
	})

	t.Run("should return http error with status = 500 when called db with error", func(t *testing.T) {
		testCtx := newTestGetTasksContext(t)
		testCtx.db.ReturnError = errors.New("MOCK_DB_ERROR")

		res := testCtx.requestWithUserID(userID, todoID)

		require.Equal(t, 500, res.Result().StatusCode)
	})
}
