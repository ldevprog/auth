package model

import (
	"time"

	desc "github.com/ldevprog/auth/pkg/user_v1"
)

type Credentials struct {
	Username string
	Password string
}

type CredentialsWithId struct {
	Id       int64
	Username string
	Role     desc.Role
	Password string
}

type UserInfoForClaims struct {
	Username string `db:"username"`
	Role     string `db:"role"`
}

type TokenWithCredentials struct {
	Token     string
	UserId    int64
	ExpiresAt time.Time
}
