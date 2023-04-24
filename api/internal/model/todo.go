package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Todo struct {
	bun.BaseModel `bun:"table:todos,alias:t"`

	ID        uuid.UUID `json:"id" bun:"id,type:uuid,pk,default:uuid_generate_v4()"`
	Name      string    `json:"name" bun:"name,type:text"`
	UserID    uuid.UUID `json:"userId" bun:"user_id,type:uuid,notnull"`
	CreatedAt time.Time `bun:"created_at,type:timestamptz,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,type:timestamptz,default:current_timestamp"`
	DeletedAt time.Time `bun:"deleted_at,type:timestamptz,soft_delete,nullzero"`
}
