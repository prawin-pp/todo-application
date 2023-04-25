package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TodoTask struct {
	bun.BaseModel `bun:"table:todo_tasks,alias:tt"`

	ID          uuid.UUID    `json:"id" bun:"id,type:uuid,pk,default:uuid_generate_v4()"`
	Name        string       `json:"name" bun:"name,type:text"`
	Description string       `json:"description" bun:"description,type:text"`
	Completed   bool         `json:"completed" bun:"completed,type:boolean,default:false"`
	DueDate     string       `json:"dueDate" bun:"due_date,type:date,nullzero"`
	SortOrder   int64        `json:"sortOrder" bun:"sort_order,type:integer,notnull"`
	TodoID      uuid.UUID    `json:"-" bun:"todo_id,type:uuid,notnull"`
	UserID      uuid.UUID    `json:"-" bun:"user_id,type:uuid,notnull"`
	CreatedAt   time.Time    `json:"createdAt" bun:"created_at,type:timestamptz,default:current_timestamp"`
	UpdatedAt   time.Time    `json:"updatedAt" bun:"updated_at,type:timestamptz,default:current_timestamp"`
	DeletedAt   bun.NullTime `json:"-" bun:"deleted_at,type:timestamptz,soft_delete,nullzero"`
}

type CreateTodoTaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	DueDate     string `json:"dueDate"`
}

type PartialUpdateTodoTaskRequest struct {
	Name        sql.NullString `json:"name"`
	Description sql.NullString `json:"description"`
	Completed   sql.NullBool   `json:"completed"`
	DueDate     sql.NullString `json:"dueDate"`
}
