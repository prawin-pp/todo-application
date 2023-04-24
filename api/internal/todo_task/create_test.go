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

type mockCreateTaskDatabase struct {
	Database

	NumberOfCalled int
	CallWithParams [][]string
	ReturnTasks    []model.TodoTask
	ReturnError    error
}

func (m *mockCreateTaskDatabase) Create(userID, todoID string, req CreateTodoTaskRequest) (*model.TodoTask, error) {
	return nil, nil
}

type testCreateTaskContext struct {
	t          *testing.T
	router     *bunrouter.Router
	db         *mockCreateTaskDatabase
	withUserID string
}

func (testContext *testCreateTaskContext) sendRequest(userID, todoID uuid.UUID, body CreateTodoTaskRequest) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	require.NoError(testContext.t, err)

	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s", todoID)
	req := httptest.NewRequest(http.MethodPost, path, &buf)
	err = testContext.router.ServeHTTPError(w, req)
	require.NoError(testContext.t, err)

	testContext.withUserID = userID.String()
	return w
}

func newTestCreateTaskContext(t *testing.T) *testCreateTaskContext {
	testCtx := &testCreateTaskContext{t: t}
	testCtx.router = bunrouter.New(bunrouter.Use(mockAuthMiddleware(func() string {
		return testCtx.withUserID
	})))
	testCtx.db = &mockCreateTaskDatabase{}

	server := NewServer(testCtx.db)
	testCtx.router.POST("/todos/:todoId", server.HandleCreateTask)

	return testCtx
}

func TestCreateTask(t *testing.T) {
	userID := uuid.New()
	taskID := uuid.New()

	t.Run("should return http status 201 when called", func(t *testing.T) {
		testCtx := newTestCreateTaskContext(t)

		res := testCtx.sendRequest(userID, taskID, CreateTodoTaskRequest{
			Name:        "MOCK_TASK_NAME",
			Description: "",
			Completed:   false,
			DueDate:     "2023-01-01",
		})

		require.Equal(t, 201, res.Result().StatusCode)
	})
}
