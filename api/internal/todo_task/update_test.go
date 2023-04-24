package todotask

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type mockPartialUpdateTask struct {
	Database

	NumberOfCalled int
	CallWithParams [][]interface{}
	ReturnError    error
}

func (m *mockPartialUpdateTask) PartialUpdate(userID, todoID string, req PartialUpdateTodoTaskRequest) (*model.TodoTask, error) {
	return nil, nil
}

type testPartialUpdateTaskContext struct {
	t          *testing.T
	router     *bunrouter.Router
	db         *mockPartialUpdateTask
	withUserID string
}

func (testContext *testPartialUpdateTaskContext) sendRequest(userID, todoID, taskID uuid.UUID, body PartialUpdateTodoTaskRequest) *httptest.ResponseRecorder {
	testContext.withUserID = userID.String()

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	require.NoError(testContext.t, err)

	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s/tasks/%s", todoID, taskID)
	req := httptest.NewRequest(http.MethodPatch, path, &buf)
	err = testContext.router.ServeHTTPError(w, req)
	require.NoError(testContext.t, err)

	return w
}

func (testContext *testPartialUpdateTaskContext) sendRequestString(userID, todoID uuid.UUID, body string) *httptest.ResponseRecorder {
	testContext.withUserID = userID.String()

	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s", todoID)
	req := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader([]byte(body)))
	err := testContext.router.ServeHTTPError(w, req)
	require.NoError(testContext.t, err)

	return w
}

func newTestPartialUpdateTaskContext(t *testing.T) *testPartialUpdateTaskContext {
	testCtx := &testPartialUpdateTaskContext{t: t}
	testCtx.router = bunrouter.New(bunrouter.Use(mockAuthMiddleware(func() string {
		return testCtx.withUserID
	})))
	testCtx.db = &mockPartialUpdateTask{}

	server := NewServer(testCtx.db)
	testCtx.router.PATCH("/todos/:todoId/tasks/:taskId", server.HandlePartialUpdateTask)

	return testCtx
}

func TestPartialUpdateTask(t *testing.T) {
	userID := uuid.New()
	todoID := uuid.New()
	taskID := uuid.New()

	t.Run("should return http status 200 when called", func(t *testing.T) {
		testCtx := newTestPartialUpdateTaskContext(t)

		res := testCtx.sendRequest(userID, todoID, taskID, PartialUpdateTodoTaskRequest{
			Name:        "MOCK_TASK_NAME",
			Description: "MOCK_DESCRIPTION",
			Completed:   true,
			DueDate:     "2023-01-01",
		})

		require.Equal(t, 200, res.Result().StatusCode)
	})

}
