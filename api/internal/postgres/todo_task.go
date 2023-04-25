package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/parwin-pp/todo-application/internal/model"
	"github.com/uptrace/bun"
)

func (db *DB) GetTasks(ctx context.Context, userID, todoID string) ([]model.TodoTask, error) {
	var todoTasks []model.TodoTask
	err := db.db.NewSelect().Model(&todoTasks).Where("user_id = ? AND todo_id = ?", userID, todoID).Scan(ctx)
	return todoTasks, err
}

func (db *DB) GetTask(ctx context.Context, userID, todoID, taskID string) (*model.TodoTask, error) {
	var todoTask model.TodoTask
	err := db.db.NewSelect().Model(&todoTask).Where("user_id = ? AND todo_id = ? AND id = ?", userID, todoID, taskID).Scan(ctx)
	return &todoTask, err
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

func (db *DB) PartialUpdateTask(ctx context.Context, userID, todoID, taskID string, req model.PartialUpdateTodoTaskRequest) (*model.TodoTask, error) {
	updated := map[string]interface{}{}
	if req.Name.Valid {
		updated["name"] = req.Name.String
	}
	if req.Description.Valid {
		updated["description"] = req.Description.String
	}
	if req.DueDate.Valid {
		updated["due_date"] = req.DueDate.String
	}
	if req.Completed.Valid {
		updated["completed"] = req.Completed.Bool
	}
	if len(updated) == 0 {
		return nil, errors.New("nothing to update")
	}

	updated["updated_at"] = bun.Safe("NOW()")
	if _, err := db.db.NewUpdate().
		Model(&updated).
		TableExpr("todo_tasks").
		Where("user_id = ?", userID).
		Where("todo_id = ?", todoID).
		Where("id = ?", taskID).
		Exec(ctx); err != nil {
		return nil, err
	}

	return db.GetTask(ctx, userID, todoID, taskID)
}

func (db *DB) DeleteTask(ctx context.Context, userID, todoID, taskID string) error {
	return db.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var deletedModel model.TodoTask
		result, err := tx.NewDelete().
			Model(&deletedModel).
			Where("user_id = ?", userID).
			Where("todo_id = ?", todoID).
			Where("id = ?", taskID).
			Returning("sort_order").
			Exec(ctx)
		if err != nil {
			return err
		}
		nums, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if nums == 0 {
			return nil
		}

		_, err = tx.NewUpdate().
			Model((*model.TodoTask)(nil)).
			Set("sort_order = sort_order - 1").
			Where("user_id = ?", userID).
			Where("todo_id = ?", todoID).
			Where("sort_order > ?", deletedModel.SortOrder).
			Exec(ctx)

		return err
	})
}
