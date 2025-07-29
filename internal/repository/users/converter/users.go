package converter

import (
	"github.com/ldevprog/auth/internal/model"
	modelRepo "github.com/ldevprog/auth/internal/repository/users/model"
)

func ToUserFullNoPassFromRepo(user *modelRepo.UserFullNoPass) *model.UserFullNoPass {
	return &model.UserFullNoPass{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
