package todotask

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/parwin-pp/todo-application/internal/middleware"
	"github.com/parwin-pp/todo-application/internal/mock"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type mockDeleteTaskDatabase struct {
	Database

	NumberOfCalled int
	CallWithParams [][]interface{}
	ReturnError    error
}

func (m *mockDeleteTaskDatabase) DeleteTask(ctx context.Context, userID, todoID, taskID string) error {
	m.NumberOfCalled++
	m.CallWithParams = append(m.CallWithParams, []interface{}{userID, todoID, taskID})
	return m.ReturnError
}

type testDeleteTaskContext struct {
	t          *testing.T
	router     *bunrouter.Router
	db         *mockDeleteTaskDatabase
	withUserID string
}

func (testCtx *testDeleteTaskContext) sendRequest(userID, todoID, taskID uuid.UUID) *httptest.ResponseRecorder {
	testCtx.withUserID = userID.String()
	w := httptest.NewRecorder()
	path := fmt.Sprintf("/todos/%s/tasks/%s", todoID, taskID)
	req := httptest.NewRequest(http.MethodDelete, path, nil)
	testCtx.router.ServeHTTP(w, req)
	return w
}

func newTestDeleteTaskContext(t *testing.T) *testDeleteTaskContext {
	db := &mockDeleteTaskDatabase{}
	testCtx := &testDeleteTaskContext{t: t}
	router := bunrouter.New(
		bunrouter.Use(middleware.NewErrorHandler),
		bunrouter.Use(mock.NewAuthMiddleware(func() string {
			return testCtx.withUserID
		})),
	)

	server := NewServer(db)
	router.DELETE("/todos/:todoId/tasks/:taskId", server.HandleDeleteTask)

	testCtx.db = db
	testCtx.router = router
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

	t.Run("should call delete task to database with correct params", func(t *testing.T) {
		testCtx := newTestDeleteTaskContext(t)

		testCtx.sendRequest(userID, todoID, taskID)

		require.Equal(t, 1, testCtx.db.NumberOfCalled)
		require.Equal(t, []interface{}{
			userID.String(),
			todoID.String(),
			taskID.String(),
		}, testCtx.db.CallWithParams[0])
	})

	t.Run("should return http status 500 when database return error", func(t *testing.T) {
		testCtx := newTestDeleteTaskContext(t)
		testCtx.db.ReturnError = errors.New("MOCK_ERROR")

		res := testCtx.sendRequest(userID, todoID, taskID)

		require.Equal(t, http.StatusInternalServerError, res.Code)
	})
}
