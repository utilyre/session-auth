package models

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID       int64  `bun:"id,pk,autoincrement"`
	Email    string `bun:"email"`
	Password []byte `bun:"password"`
	Name     string `bun:"name"`
}
