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
	"github.com/parwin-pp/todo-application/internal/middleware"
	"github.com/parwin-pp/todo-application/internal/mock"
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

func (testCtx *testGetTasksContext) createTask(userID, todoID uuid.UUID, name string) *model.TodoTask {
	task := model.TodoTask{
		ID:     uuid.New(),
		Name:   name,
		TodoID: todoID,
		UserID: userID,
	}
	testCtx.db.ReturnTasks = append(testCtx.db.ReturnTasks, task)
	return &task
}

func (testCtx *testGetTasksContext) requestWithUserID(userID uuid.UUID, todoID uuid.UUID) *httptest.ResponseRecorder {
	testCtx.withUserID = userID.String()
	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s", todoID)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	testCtx.router.ServeHTTP(w, req)
	return w
}

func newTestGetTasksContext(t *testing.T) *testGetTasksContext {
	db := &mockGetTasksDatabase{}
	testCtx := &testGetTasksContext{t: t}
	router := bunrouter.New(
		bunrouter.Use(middleware.NewErrorHandler),
		bunrouter.Use(mock.NewAuthMiddleware(func() string {
			return testCtx.withUserID
		})),
	)

	server := NewServer(db)
	router.GET("/todos/:todoId", server.HandleGetTasks)

	testCtx.db = db
	testCtx.router = router
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
