package model

import desc "github.com/ldevprog/auth/pkg/user_v1"

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
