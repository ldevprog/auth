package converter

import (
	"github.com/ldevprog/auth/internal/model"
	desc "github.com/ldevprog/auth/pkg/auth_v1"
)

func ToCredentialsFromDesc(req *desc.LoginRequest) *model.Credentials {
	return &model.Credentials{
		Name:     req.GetName(),
		Password: req.GetPassword(),
	}
}
