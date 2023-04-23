package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        uuid.UUID    `bun:"id,type:uuid,pk,default:uuid_generate_v4()"`
	Username  string       `bun:"username,type:text,notnull"`
	Password  string       `bun:"password,type:text,notnull"`
	CreatedAt time.Time    `bun:"created_at,type:timestamptz,default:current_timestamp"`
	UpdatedAt time.Time    `bun:"updated_at,type:timestamptz,default:current_timestamp"`
	DeletedAt bun.NullTime `bun:"deleted_at,type:timestamptz,soft_delete"`
}
