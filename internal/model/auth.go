package model

import desc "github.com/ldevprog/auth/pkg/user_v1"

type Credentials struct {
	Name     string
	Password string
}

type CredentialsWithId struct {
	Id       int64
	Name     string
	Role     desc.Role
	Password string
}

type UserInfoForClaims struct {
	Name string `json:"name"`
	Role string `json:"role"`
}
