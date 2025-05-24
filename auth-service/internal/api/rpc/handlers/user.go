package handlers

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	v1 "github.com/tech-inspire/api-contracts/api/gen/go/auth/v1"
	"github.com/tech-inspire/backend/auth-service/internal/service/dto"
	authmiddleware "github.com/tech-inspire/backend/auth-service/pkg/jwt/middleware"
)

type UserHandler struct {
	userService   UserService
	avatarService AvatarService
}

func NewUserHandler(avatarService AvatarService, userService UserService) *UserHandler {
	return &UserHandler{userService: userService, avatarService: avatarService}
}

func (a UserHandler) GetMe(ctx context.Context, c *connect.Request[v1.GetMeRequest]) (*connect.Response[v1.GetUserResponse], error) {
	userID := authmiddleware.GetUserInfo(ctx).UserID

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
	token := authmiddleware.GetUserInfo(ctx)

	if len(c.Msg.Content) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("content is empty"))
	}

	contentType := http.DetectContentType(c.Msg.Content)
	if contentType != c.Msg.ContentType {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid content type: %s != %s", contentType, c.Msg.ContentType))
	}

	allowedContentTypes := []string{"image/jpeg", "image/png", "image/webp"}
	if !slices.Contains(allowedContentTypes, contentType) {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("content type %s is not supported", contentType))
	}

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
