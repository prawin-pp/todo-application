package model

import (
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
	DueDate     string       `json:"dueDate" bun:"due_date,type:date"`
	SortOrder   int64        `json:"sortOrder" bun:"sort_order,type:integer,notnull"`
	TodoID      uuid.UUID    `json:"-" bun:"todo_id,type:uuid,notnull"`
	UserID      uuid.UUID    `json:"-" bun:"user_id,type:uuid,notnull"`
	CreatedAt   time.Time    `json:"createdAt" bun:"created_at,type:timestamptz,default:current_timestamp"`
	UpdatedAt   time.Time    `json:"updatedAt" bun:"updated_at,type:timestamptz,default:current_timestamp"`
	DeletedAt   bun.NullTime `json:"-" bun:"deleted_at,type:timestamptz,soft_delete"`
}
