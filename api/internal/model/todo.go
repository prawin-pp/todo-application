package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Todo struct {
	bun.BaseModel `bun:"table:todos,alias:t"`

	ID        uuid.UUID    `json:"id" bun:"id,type:uuid,pk,default:uuid_generate_v4()"`
	Name      string       `json:"name" bun:"name,type:text"`
	UserID    uuid.UUID    `json:"-" bun:"user_id,type:uuid,notnull"`
	CreatedAt time.Time    `json:"createdAt" bun:"created_at,type:timestamptz,default:current_timestamp"`
	UpdatedAt time.Time    `json:"updatedAt" bun:"updated_at,type:timestamptz,default:current_timestamp"`
	DeletedAt bun.NullTime `json:"-" bun:"deleted_at,type:timestamptz,soft_delete,nullzero"`
}
