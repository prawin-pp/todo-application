package todotask

import (
	"bytes"
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

type mockPartialUpdateTask struct {
	Database

	NumberOfCalled int
	CallWithParams [][]interface{}
	ExistsTasks    []model.TodoTask
	ReturnError    error
}

func (m *mockPartialUpdateTask) CreateTask(task model.TodoTask) {
	m.ExistsTasks = append(m.ExistsTasks, task)
}

func (m *mockPartialUpdateTask) PartialUpdate(userID, todoID, taskID string, req PartialUpdateTodoTaskRequest) (*model.TodoTask, error) {
	m.NumberOfCalled++
	m.CallWithParams = append(m.CallWithParams, []interface{}{userID, todoID, taskID, req})

	var result *model.TodoTask
	for _, existsTask := range m.ExistsTasks {
		isUserIDMatch := existsTask.UserID.String() == userID
		isTodoIDMatch := existsTask.TodoID.String() == todoID
		isTaskIDMatch := existsTask.ID.String() == taskID
		if isUserIDMatch && isTodoIDMatch && isTaskIDMatch {
			result = &existsTask
		}
	}
	if result == nil {
		return nil, m.ReturnError
	}
	result.Name = req.Name
	result.Description = req.Description
	result.Completed = req.Completed
	result.DueDate = req.DueDate
	return result, m.ReturnError
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

func (testContext *testPartialUpdateTaskContext) sendRequestString(userID, todoID, taskID uuid.UUID, body string) *httptest.ResponseRecorder {
	testContext.withUserID = userID.String()

	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s/tasks/%s", todoID, taskID)
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
	reqBody := PartialUpdateTodoTaskRequest{
		Name:        "MOCK_TASK_NAME",
		Description: "MOCK_DESCRIPTION",
		Completed:   true,
		DueDate:     "2023-01-01",
	}

	t.Run("should return http status 200 when called", func(t *testing.T) {
		testCtx := newTestPartialUpdateTaskContext(t)

		res := testCtx.sendRequest(userID, todoID, taskID, reqBody)

		require.Equal(t, 200, res.Result().StatusCode)
	})

	t.Run("should call partial update to database when called with given params", func(t *testing.T) {
		testCtx := newTestPartialUpdateTaskContext(t)

		testCtx.sendRequest(userID, todoID, taskID, reqBody)

		require.Equal(t, 1, testCtx.db.NumberOfCalled)
		require.Equal(t, []interface{}{
			userID.String(),
			todoID.String(),
			taskID.String(),
			reqBody,
		}, testCtx.db.CallWithParams[0])
	})

	t.Run("should return status 400 when request is invalid json format", func(t *testing.T) {
		testCtx := newTestPartialUpdateTaskContext(t)

		res := testCtx.sendRequestString(userID, todoID, taskID, `{ #:## }`)

		require.Equal(t, 400, res.Result().StatusCode)
	})

	t.Run("should return response body with updated task", func(t *testing.T) {
		testCtx := newTestPartialUpdateTaskContext(t)
		testCtx.db.CreateTask(model.TodoTask{
			ID:          taskID,
			UserID:      userID,
			TodoID:      todoID,
			Name:        "NAME",
			Description: "DESC",
			Completed:   false,
			DueDate:     "2023-01-01",
		})

		res := testCtx.sendRequest(userID, todoID, taskID, reqBody)

		var resBody model.TodoTask
		err := json.NewDecoder(res.Body).Decode(&resBody)

		require.NoError(t, err)
		require.Equal(t, reqBody.Name, resBody.Name)
		require.Equal(t, reqBody.Description, resBody.Description)
		require.Equal(t, reqBody.Completed, resBody.Completed)
		require.Equal(t, reqBody.DueDate, resBody.DueDate)
	})

	t.Run("should return status 500 when called database error", func(t *testing.T) {
		testCtx := newTestPartialUpdateTaskContext(t)
		testCtx.db.ReturnError = errors.New("MOCK_ERROR")

		res := testCtx.sendRequest(userID, todoID, taskID, reqBody)

		require.Equal(t, 500, res.Result().StatusCode)
	})
}
