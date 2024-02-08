package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel `bun:"table:sessions"`

	ID        int64     `bun:"id,pk,autoincrement"`
	CreatedAt time.Time `bun:"created_at"`
	UUID      uuid.UUID `bun:"uuid"`
	LastIP    string    `bun:"last_ip"`

	UserID int64 `bun:"user_id"`
	User   *User `bun:"rel:belongs-to,join:user_id=id"`
}
