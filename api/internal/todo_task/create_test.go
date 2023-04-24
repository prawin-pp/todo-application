package todotask

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type mockCreateTaskDatabase struct {
	Database

	NumberOfCalled int
	CallWithParams [][]interface{}
	ReturnError    error
}

func (m *mockCreateTaskDatabase) Create(userID, todoID string, req CreateTodoTaskRequest) (*model.TodoTask, error) {
	m.NumberOfCalled++
	m.CallWithParams = append(m.CallWithParams, []interface{}{userID, todoID, req})
	return &model.TodoTask{
		ID:          uuid.New(),
		UserID:      uuid.MustParse(userID),
		TodoID:      uuid.MustParse(todoID),
		Name:        req.Name,
		Description: req.Description,
		Completed:   req.Completed,
		DueDate:     req.DueDate,
		SortOrder:   0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, m.ReturnError
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

	t.Run("should return response body with created task from database", func(t *testing.T) {
		testCtx := newTestCreateTaskContext(t)
		body := CreateTodoTaskRequest{
			Name:        "MOCK_TASK_NAME",
			Description: "",
			Completed:   false,
			DueDate:     "2023-01-01",
		}

		res := testCtx.sendRequest(userID, todoID, body)

		require.Equal(t, 201, res.Result().StatusCode)

		var resBody model.TodoTask
		err := json.NewDecoder(res.Body).Decode(&resBody)
		require.NoError(t, err)
		require.Equal(t, "MOCK_TASK_NAME", resBody.Name)
		require.Equal(t, "", resBody.Description)
		require.Equal(t, false, resBody.Completed)
		require.Equal(t, "2023-01-01", resBody.DueDate)
	})

	t.Run("should return http status = 500 when called database error", func(t *testing.T) {
		testCtx := newTestCreateTaskContext(t)
		testCtx.db.ReturnError = errors.New("DATABASE_ERROR")
		body := CreateTodoTaskRequest{}

		res := testCtx.sendRequest(userID, todoID, body)

		require.Equal(t, 500, res.Result().StatusCode)
	})
}
