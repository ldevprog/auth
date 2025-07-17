package model

import (
	"database/sql"
	"time"

	desc "github.com/levon-dalakyan/auth/pkg/user_v1"
)

type UserFullNoPass struct {
	Id        int64
	Name      string
	Email     string
	Role      desc.Role
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
