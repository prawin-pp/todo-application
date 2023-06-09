package todotask

import (
	"bytes"
	"context"
	"database/sql"
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

type mockPartialUpdateTask struct {
	Database

	NumberOfCalled int
	CallWithParams [][]interface{}
	ExistsTasks    []model.TodoTask
	ReturnError    error
}

func (m *mockPartialUpdateTask) PartialUpdateTask(ctx context.Context, userID, todoID, taskID string, req model.PartialUpdateTodoTaskRequest) (*model.TodoTask, error) {
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
	result.Name = req.Name.String
	result.Description = req.Description.String
	result.Completed = req.Completed.Bool
	result.DueDate = req.DueDate.String
	return result, m.ReturnError
}

type testPartialUpdateTaskContext struct {
	t          *testing.T
	router     *bunrouter.Router
	db         *mockPartialUpdateTask
	withUserID string
}

func (testCtx *testPartialUpdateTaskContext) createTask(task model.TodoTask) {
	testCtx.db.ExistsTasks = append(testCtx.db.ExistsTasks, task)
}

func (testCtx *testPartialUpdateTaskContext) sendRequest(userID, todoID, taskID uuid.UUID, body model.PartialUpdateTodoTaskRequest) *httptest.ResponseRecorder {
	testCtx.withUserID = userID.String()

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	require.NoError(testCtx.t, err)

	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s/tasks/%s", todoID, taskID)
	req := httptest.NewRequest(http.MethodPatch, path, &buf)
	testCtx.router.ServeHTTP(w, req)
	return w
}

func (testCtx *testPartialUpdateTaskContext) sendRequestString(userID, todoID, taskID uuid.UUID, body string) *httptest.ResponseRecorder {
	testCtx.withUserID = userID.String()
	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s/tasks/%s", todoID, taskID)
	req := httptest.NewRequest(http.MethodPatch, path, bytes.NewReader([]byte(body)))
	testCtx.router.ServeHTTP(w, req)
	return w
}

func newTestPartialUpdateTaskContext(t *testing.T) *testPartialUpdateTaskContext {
	db := &mockPartialUpdateTask{}
	testCtx := &testPartialUpdateTaskContext{t: t}
	router := bunrouter.New(
		bunrouter.Use(middleware.NewErrorHandler),
		bunrouter.Use(mock.NewAuthMiddleware(func() string {
			return testCtx.withUserID
		})),
	)

	server := NewServer(db)
	router.PATCH("/todos/:todoId/tasks/:taskId", server.HandlePartialUpdateTask)

	testCtx.db = db
	testCtx.router = router
	return testCtx
}

func TestPartialUpdateTask(t *testing.T) {
	userID := uuid.New()
	todoID := uuid.New()
	taskID := uuid.New()
	reqBody := model.PartialUpdateTodoTaskRequest{
		Name:        model.NullString{NullString: sql.NullString{String: "MOCK_TASK_NAME", Valid: true}},
		Description: model.NullString{NullString: sql.NullString{String: "MOCK_DESCRIPTION", Valid: true}},
		Completed:   model.NullBool{NullBool: sql.NullBool{Bool: true, Valid: true}},
		DueDate:     model.NullString{NullString: sql.NullString{String: "", Valid: false}},
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
		testCtx.createTask(model.TodoTask{
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
		require.Equal(t, reqBody.Name.String, resBody.Name)
		require.Equal(t, reqBody.Description.String, resBody.Description)
		require.Equal(t, reqBody.Completed.Bool, resBody.Completed)
		require.Equal(t, reqBody.DueDate.String, resBody.DueDate)
	})

	t.Run("should return status 500 when called database error", func(t *testing.T) {
		testCtx := newTestPartialUpdateTaskContext(t)
		testCtx.db.ReturnError = errors.New("MOCK_ERROR")

		res := testCtx.sendRequest(userID, todoID, taskID, reqBody)

		require.Equal(t, 500, res.Result().StatusCode)
	})
}
