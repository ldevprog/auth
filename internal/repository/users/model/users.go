package model

import (
	"database/sql"
	"time"

	desc "github.com/ldevprog/auth/pkg/user_v1"
)

type UserFullNoPass struct {
	Id        int64        `db:"id"`
	Name      string       `db:"name"`
	Username  string       `db:"username"`
	Email     string       `db:"email"`
	Role      desc.Role    `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
