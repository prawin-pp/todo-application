package todo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/parwin-pp/todo-application/internal/middleware"
	"github.com/parwin-pp/todo-application/internal/mock"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bunrouter"
)

type mockCreateTodoDatabase struct {
	Database

	NumberOfCalled int
	CallWithParams [][]string
	ReturnTodo     *model.Todo
	ReturnError    error
}

func (m *mockCreateTodoDatabase) CreateTodo(ctx context.Context, userID, name string) (*model.Todo, error) {
	m.NumberOfCalled++
	m.CallWithParams = append(m.CallWithParams, []string{userID, name})
	if m.ReturnError != nil {
		return nil, m.ReturnError
	}
	return &model.Todo{
		ID:        uuid.New(),
		Name:      name,
		UserID:    uuid.MustParse(userID),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

type testCreateTodoContext struct {
	t      *testing.T
	router *bunrouter.Router
	db     *mockCreateTodoDatabase

	withUserID string
}

func (testCtx *testCreateTodoContext) withCreateTodoError(err error) {
	testCtx.db.ReturnError = err
}

func (testCtx *testCreateTodoContext) requestWithUserID(userID uuid.UUID, body io.Reader) *httptest.ResponseRecorder {
	testCtx.withUserID = userID.String()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/todos", body)
	req.Header.Set("Content-Type", "application/json")
	testCtx.router.ServeHTTP(w, req)
	return w
}

func newTestCreateTodoContext(t *testing.T) *testCreateTodoContext {
	db := &mockCreateTodoDatabase{}
	testCtx := &testCreateTodoContext{t: t}

	server := NewServer(db)

	router := bunrouter.New(
		bunrouter.Use(middleware.NewErrorHandler),
		bunrouter.Use(mock.NewAuthMiddleware(func() string {
			return testCtx.withUserID
		})),
	)
	router.POST("/todos", server.HandleCreateTodo)

	testCtx.db = db
	testCtx.router = router
	testCtx.withUserID = uuid.NewString()
	return testCtx
}
func TestCreateTodo(t *testing.T) {
	userID := uuid.New()

	t.Run("should return http status 201 when called", func(t *testing.T) {
		testCtx := newTestCreateTodoContext(t)
		req := []byte(`{ "name": "MOCK_NAME" }`)

		resp := testCtx.requestWithUserID(userID, bytes.NewReader(req))

		require.Equal(t, 201, resp.Result().StatusCode)
	})

	t.Run("should return http status 400 when request body is not json", func(t *testing.T) {
		testCtx := newTestCreateTodoContext(t)
		req := []byte(`{ #: ## }`)

		resp := testCtx.requestWithUserID(userID, bytes.NewReader(req))

		require.Equal(t, 400, resp.Result().StatusCode)
	})

	t.Run("should call create todo to database given todo name = 'MOCK_NAME'", func(t *testing.T) {
		testCtx := newTestCreateTodoContext(t)
		req := []byte(`{ "name": "MOCK_NAME" }`)

		testCtx.requestWithUserID(userID, bytes.NewReader(req))

		require.Equal(t, 1, testCtx.db.NumberOfCalled)
		require.Equal(t, []string{userID.String(), "MOCK_NAME"}, testCtx.db.CallWithParams[0])
	})

	t.Run("should return response body with exists todo in database", func(t *testing.T) {
		testCtx := newTestCreateTodoContext(t)
		req := []byte(`{ "name": "MOCK_NAME" }`)

		res := testCtx.requestWithUserID(userID, bytes.NewReader(req))

		var todo model.Todo
		err := json.NewDecoder(res.Body).Decode(&todo)
		require.NoError(t, err)
		require.Equal(t, "MOCK_NAME", todo.Name)
	})

	t.Run("should return http error with status = 500 when called db with error", func(t *testing.T) {
		testCtx := newTestCreateTodoContext(t)
		testCtx.withCreateTodoError(errors.New("SOMETHING_WENT_WRONG"))
		req := []byte(`{ "name": "MOCK_NAME" }`)

		res := testCtx.requestWithUserID(userID, bytes.NewReader(req))

		require.Equal(t, 500, res.Result().StatusCode)
	})
}
