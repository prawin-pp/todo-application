package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TodoTask struct {
	bun.BaseModel `bun:"table:todo_tasks,alias:tt"`

	ID          uuid.UUID    `bun:"id,type:uuid,pk,default:uuid_generate_v4()"`
	Name        string       `bun:"name,type:text"`
	Description string       `bun:"description,type:text"`
	Completed   bool         `bun:"completed,type:boolean,default:false"`
	DueDate     string       `bun:"due_date,type:date"`
	SortOrder   int64        `bun:"sort_order,type:integer,notnull"`
	TodoID      uuid.UUID    `bun:"todo_id,type:uuid,notnull"`
	UserID      uuid.UUID    `bun:"user_id,type:uuid,notnull"`
	CreatedAt   time.Time    `bun:"created_at,type:timestamptz,default:current_timestamp"`
	UpdatedAt   time.Time    `bun:"updated_at,type:timestamptz,default:current_timestamp"`
	DeletedAt   bun.NullTime `bun:"deleted_at,type:timestamptz,soft_delete"`
}
