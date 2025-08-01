package converter

import (
	"github.com/ldevprog/auth/internal/model"
	modelRepo "github.com/ldevprog/auth/internal/repository/auth/model"
)

func ToCredentialsWithIdFromRepo(credentials *modelRepo.CredentialsWithId) *model.CredentialsWithId {
	return &model.CredentialsWithId{
		Id:       credentials.Id,
		Username: credentials.Username,
		Role:     credentials.Role,
		Password: credentials.Password,
	}
}
