package model

import "github.com/ldevprog/auth/pkg/user_v1"

type CredentialsWithId struct {
	Id       int64        `db:"id"`
	Username string       `db:"username"`
	Role     user_v1.Role `db:"role"`
	Password string       `db:"password"`
}
