package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/levon-dalakyan/auth/internal/model"
	desc "github.com/levon-dalakyan/auth/pkg/user_v1"
)

func ToGetResponseFromService(user *model.UserFullNoPass) *desc.GetResponse {
	var updatedAtTime *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAtTime = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.GetResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAtTime,
	}
}

func ToUserFromDesc(req *desc.CreateRequest) *model.User {
	return &model.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Role:     req.GetRole(),
		Password: req.GetPassword(),
	}
}

func ToUserChangableFromDesc(req *desc.UpdateRequest) *model.UserChangable {
	var namePtr *string
	if req.GetName() != nil {
		namePtr = &req.GetName().Value
	}

	var emailPtr *string
	if req.GetEmail() != nil {
		emailPtr = &req.GetEmail().Value
	}

	return &model.UserChangable{
		Id:    req.GetId(),
		Name:  namePtr,
		Email: emailPtr,
	}
}
