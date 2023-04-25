package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/parwin-pp/todo-application/internal/model"
)

func (db *DB) GetTasks(ctx context.Context, userID, todoID string) ([]model.TodoTask, error) {
	var todoTasks []model.TodoTask
	err := db.db.NewSelect().Model(&todoTasks).Where("user_id = ? AND todo_id = ?", userID, todoID).Scan(ctx)
	return todoTasks, err
}

func (db *DB) CreateTask(ctx context.Context, userID, todoID string, req model.CreateTodoTaskRequest) (*model.TodoTask, error) {
	result := &model.TodoTask{
		UserID:      uuid.MustParse(userID),
		TodoID:      uuid.MustParse(todoID),
		Name:        req.Name,
		Description: req.Description,
		Completed:   req.Completed,
		DueDate:     req.DueDate,
	}

	if _, err := db.db.NewInsert().
		Model(result).
		Value("sort_order", `(
			SELECT COALESCE(MAX(sort_order), 0) + 1
			FROM todo_tasks
			WHERE user_id = ? AND todo_id = ?
		)`, userID, todoID).
		Returning("*").
		Exec(ctx); err != nil {
		return nil, err
	}

	return result, nil
}
