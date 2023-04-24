package todo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

func TestGetTodos(t *testing.T) {
	userID := uuid.New()
	router := bunrouter.New()
	server := NewServer(nil)
	router.GET("/todos", server.HandleGetTodos)

	t.Run("should return http status 200 when called", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/todos", nil)

		err := router.ServeHTTPError(w, req)

		require.NoError(t, err)
		require.Equal(t, 200, w.Result().StatusCode)
	})

	t.Run("should call get todos from database when called", func(t *testing.T) {
		testCtx := newTestGetTodosContext(t)
		testCtx.createTodo(userID, "MOCK_TODO")

		res := testCtx.sendRequest()

		require.Equal(t, 200, res.Result().StatusCode)
		require.Equal(t, 1, testCtx.db.NumberOfCalled)
	})
}
