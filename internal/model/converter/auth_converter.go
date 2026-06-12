package converter

import (
	"github.com/wiradana/backend/internal/entity"
	"github.com/wiradana/backend/internal/model"
)

func ToAppUserResponse(user *entity.AppUser) model.AppUserResponse {
	resp := model.AppUserResponse{
		ID:    user.ID.String(),
		Email: user.Email,
		Role:  user.Role,
	}
	if user.MemberID != nil {
		s := user.MemberID.String()
		resp.MemberID = &s
	}
	return resp
}
