package converter

import (
	"github.com/ldevprog/auth/internal/model"
	modelRepo "github.com/ldevprog/auth/internal/repository/auth/model"
	"github.com/ldevprog/auth/pkg/user_v1"
)

func ToCredentialsWithIdFromRepo(credentials *modelRepo.CredentialsWithId) *model.CredentialsWithId {
	roleName, ok := user_v1.Role_value[credentials.Role]
	if !ok {
		roleName = 0
	}

	return &model.CredentialsWithId{
		Id:       credentials.Id,
		Name:     credentials.Name,
		Role:     user_v1.Role(roleName),
		Password: credentials.Password,
	}
}
