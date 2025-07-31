package model

type CredentialsWithId struct {
	Id       int64  `db:"id"`
	Name     string `db:"name"`
	Role     string `db:"role"`
	Password string `db:"password"`
}
