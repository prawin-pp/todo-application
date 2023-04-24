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
	CallWithParams [][]interface{}
	ReturnTasks    []model.TodoTask
	ReturnError    error
}

func (m *mockCreateTaskDatabase) Create(userID, todoID string, req CreateTodoTaskRequest) (*model.TodoTask, error) {
	m.NumberOfCalled++
	m.CallWithParams = append(m.CallWithParams, []interface{}{userID, todoID, req})
	return nil, nil
}

type testCreateTaskContext struct {
	t          *testing.T
	router     *bunrouter.Router
	db         *mockCreateTaskDatabase
	withUserID string
}

func (testContext *testCreateTaskContext) sendRequest(userID, todoID uuid.UUID, body CreateTodoTaskRequest) *httptest.ResponseRecorder {
	testContext.withUserID = userID.String()

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	require.NoError(testContext.t, err)

	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s", todoID)
	req := httptest.NewRequest(http.MethodPost, path, &buf)
	err = testContext.router.ServeHTTPError(w, req)
	require.NoError(testContext.t, err)

	return w
}

func (testContext *testCreateTaskContext) sendRequestString(userID, todoID uuid.UUID, body string) *httptest.ResponseRecorder {
	testContext.withUserID = userID.String()

	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s", todoID)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader([]byte(body)))
	err := testContext.router.ServeHTTPError(w, req)
	require.NoError(testContext.t, err)

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
	todoID := uuid.New()

	t.Run("should return http status 201 when called", func(t *testing.T) {
		testCtx := newTestCreateTaskContext(t)

		res := testCtx.sendRequest(userID, todoID, CreateTodoTaskRequest{
			Name:        "MOCK_TASK_NAME",
			Description: "",
			Completed:   false,
			DueDate:     "2023-01-01",
		})

		require.Equal(t, 201, res.Result().StatusCode)
	})

	t.Run("should call create task to database when called", func(t *testing.T) {
		testCtx := newTestCreateTaskContext(t)
		body := CreateTodoTaskRequest{
			Name:        "MOCK_TASK_NAME",
			Description: "",
			Completed:   false,
			DueDate:     "2023-01-01",
		}

		testCtx.sendRequest(userID, todoID, body)

		require.Equal(t, 1, testCtx.db.NumberOfCalled)
		require.Equal(t, []interface{}{
			userID.String(),
			todoID.String(),
			body,
		}, testCtx.db.CallWithParams[0])
	})

	t.Run("should return http status = 400 when request body is invalid json format", func(t *testing.T) {
		testCtx := newTestCreateTaskContext(t)

		res := testCtx.sendRequestString(userID, todoID, `{#}`)

		require.Equal(t, 400, res.Result().StatusCode)
	})
}
