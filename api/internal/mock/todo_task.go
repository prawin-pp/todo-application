package mock

import (
	"context"

	"github.com/parwin-pp/todo-application/internal/model"
)

type TaskDatabase struct {
	GetTasksFn          func(ctx context.Context, userID, todoID string) ([]model.TodoTask, error)
	CreateTaskFn        func(ctx context.Context, userID, todoID string, req model.CreateTodoTaskRequest) (*model.TodoTask, error)
	PartialUpdateTaskFn func(ctx context.Context, userID, todoID, taskID string, req model.PartialUpdateTodoTaskRequest) (*model.TodoTask, error)
	DeleteTaskFn        func(ctx context.Context, userID, todoID, taskID string) error
}

func (db *TaskDatabase) GetTasks(ctx context.Context, userID, todoID string) ([]model.TodoTask, error) {
	return db.GetTasksFn(ctx, userID, todoID)
}

func (db *TaskDatabase) CreateTask(ctx context.Context, userID, todoID string, req model.CreateTodoTaskRequest) (*model.TodoTask, error) {
	return db.CreateTaskFn(ctx, userID, todoID, req)
}

func (db *TaskDatabase) PartialUpdateTask(ctx context.Context, userID, todoID, taskID string, req model.PartialUpdateTodoTaskRequest) (*model.TodoTask, error) {
	return db.PartialUpdateTaskFn(ctx, userID, todoID, taskID, req)
}

func (db *TaskDatabase) DeleteTask(ctx context.Context, userID, todoID, taskID string) error {
	return db.DeleteTaskFn(ctx, userID, todoID, taskID)
}
