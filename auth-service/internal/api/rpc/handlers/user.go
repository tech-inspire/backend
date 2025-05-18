package handlers

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"github.com/google/uuid"
	v1 "github.com/tech-inspire/api-contracts/api/gen/go/auth/v1"
	"github.com/tech-inspire/service/auth-service/internal/api/jwt"
	"github.com/tech-inspire/service/auth-service/internal/api/rpc/middleware"
	"github.com/tech-inspire/service/auth-service/internal/service/dto"
)

type UserHandler struct {
	userService   UserService
	avatarService AvatarService
}

func NewUserHandler(avatarService AvatarService, userService UserService) *UserHandler {
	return &UserHandler{userService: userService, avatarService: avatarService}
}

func (a UserHandler) GetMe(ctx context.Context, c *connect.Request[v1.GetMeRequest]) (*connect.Response[v1.GetUserResponse], error) {
	userID := middleware.GetUserInfo(ctx).UserID

	user, err := a.userService.GetCurrentUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user %s: %w", userID, err)
	}

	return connect.NewResponse(&v1.GetUserResponse{
		User: userPB(*user.User),
	}), nil
}

func (a UserHandler) UpdateUser(ctx context.Context, c *connect.Request[v1.UpdateUserRequest]) (*connect.Response[v1.User], error) {
	// TODO implement me
	panic("implement me")
}

func (a UserHandler) GetUser(ctx context.Context, c *connect.Request[v1.GetUserRequest]) (*connect.Response[v1.GetUserResponse], error) {
	userID, err := uuid.Parse(c.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse id: %w", err))
	}

	user, err := a.userService.GetCurrentUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user %s: %w", userID, err)
	}

	return connect.NewResponse(&v1.GetUserResponse{
		User: userPB(*user.User),
	}), nil
}

func (a UserHandler) UploadAvatar(ctx context.Context, c *connect.Request[v1.UploadUserAvatarRequest]) (*connect.Response[v1.UploadUserAvatarResponse], error) {
	token := authn.GetInfo(ctx).(*jwt.ValidateUserAccessTokenOutput)

	contentType := http.DetectContentType(c.Msg.Content)

	err := a.avatarService.UploadUserAvatar(ctx, dto.UploadUserAvatar{
		Data:        c.Msg.Content,
		UserID:      token.UserID,
		ImageSize:   int64(len(c.Msg.Content)),
		ContentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&v1.UploadUserAvatarResponse{}), nil
}
