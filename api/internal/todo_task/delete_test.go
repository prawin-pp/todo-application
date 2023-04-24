package todotask

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type mockDeleteTaskDatabase struct {
	Database

	NumberOfCalled int
	CallWithParams [][]interface{}
	ReturnError    error
}

func (m *mockDeleteTaskDatabase) Delete(userID, todoID, taskID string) error {
	m.NumberOfCalled++
	m.CallWithParams = append(m.CallWithParams, []interface{}{userID, todoID, taskID})
	return m.ReturnError
}

type testDeleteTaskContext struct {
	t          *testing.T
	router     *bunrouter.Router
	db         *mockCreateTaskDatabase
	withUserID string
}

func (testContext *testDeleteTaskContext) sendRequest(userID, todoID, taskID uuid.UUID) *httptest.ResponseRecorder {
	testContext.withUserID = userID.String()

	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s/tasks/%s", todoID, taskID)
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	err := testContext.router.ServeHTTPError(w, req)
	require.NoError(testContext.t, err)

	return w
}

func newTestDeleteTaskContext(t *testing.T) *testDeleteTaskContext {
	testCtx := &testDeleteTaskContext{t: t}
	testCtx.router = bunrouter.New(bunrouter.Use(mockAuthMiddleware(func() string {
		return testCtx.withUserID
	})))
	testCtx.db = &mockCreateTaskDatabase{}

	server := NewServer(testCtx.db)
	testCtx.router.DELETE("/todos/:todoId/tasks/:taskId", server.HandleDeleteTask)

	return testCtx
}

func TestHandleDeleteTask(t *testing.T) {
	userID := uuid.New()
	todoID := uuid.New()
	taskID := uuid.New()

	t.Run("should return http status 204 when called", func(t *testing.T) {
		testCtx := newTestDeleteTaskContext(t)

		res := testCtx.sendRequest(userID, todoID, taskID)

		require.Equal(t, http.StatusNoContent, res.Code)
	})
}
