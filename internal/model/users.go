package model

import (
	"database/sql"
	"time"

	desc "github.com/ldevprog/auth/pkg/user_v1"
)

type User struct {
	Name            string
	Email           string
	Role            desc.Role
	Password        string
	PasswordConfirm string
}

type UserFullNoPass struct {
	Id        int64
	Name      string
	Email     string
	Role      desc.Role
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserChangable struct {
	Id    int64
	Name  *string
	Email *string
}
