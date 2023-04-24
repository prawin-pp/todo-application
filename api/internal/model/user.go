package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        uuid.UUID    `json:"id" bun:"id,type:uuid,pk,default:uuid_generate_v4()"`
	Username  string       `json:"username" bun:"username,type:text,notnull"`
	Password  string       `json:"-" bun:"password,type:text,notnull"`
	CreatedAt time.Time    `json:"-" bun:"created_at,type:timestamptz,default:current_timestamp"`
	UpdatedAt time.Time    `json:"-" bun:"updated_at,type:timestamptz,default:current_timestamp"`
	DeletedAt bun.NullTime `json:"-" bun:"deleted_at,type:timestamptz,soft_delete,nullzero"`
}
