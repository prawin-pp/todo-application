package todo

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

func TestCreateTodo(t *testing.T) {
	userID := uuid.New()

	t.Run("should return http status 201 when called", func(t *testing.T) {
		testCtx := newTestCreateTodoContext(t)
		req := []byte(`{ "name": "MOCK_NAME" }`)

		resp := testCtx.requestWithUserID(userID, bytes.NewReader(req))

		require.Equal(t, 201, resp.Result().StatusCode)
	})
}

type testCreateTodoContext struct {
	t      *testing.T
	router *bunrouter.Router
	db     *mockTodoDatabase

	withUserID string
}

func (testContext *testCreateTodoContext) requestWithUserID(userID uuid.UUID, body io.Reader) *httptest.ResponseRecorder {
	testContext.withUserID = userID.String()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	req.Header.Set("Content-Type", "application/json")

	err := testContext.router.ServeHTTPError(w, req)
	require.NoError(testContext.t, err)

	return w
}

func newTestCreateTodoContext(t *testing.T) *testCreateTodoContext {
	testCtx := &testCreateTodoContext{t: t}
	testCtx.router = bunrouter.New(bunrouter.Use(mockAuthMiddleware(func() string {
		return testCtx.withUserID
	})))
	testCtx.db = &mockTodoDatabase{}

	server := NewServer(testCtx.db)
	testCtx.router.POST("/todos", server.HandleCreateTodo)

	return testCtx
}
