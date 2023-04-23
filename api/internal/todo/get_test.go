package todo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

func TestGetTodo(t *testing.T) {
	router := bunrouter.New()
	server := NewServer()
	router.GET("/todos", server.HandleGetTodo)

	t.Run("should return http status 200 when called", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/todos", nil)

		err := router.ServeHTTPError(w, req)

		require.NoError(t, err)
		require.Equal(t, 200, w.Result().StatusCode)
	})
}
