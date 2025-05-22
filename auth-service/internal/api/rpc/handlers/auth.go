package handlers

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/go-errors/errors"
	v1 "github.com/tech-inspire/api-contracts/api/gen/go/auth/v1"
	"github.com/tech-inspire/backend/auth-service/internal/api/jwt"
	"github.com/tech-inspire/backend/auth-service/internal/apperrors"
	"github.com/tech-inspire/backend/auth-service/internal/clients/mail"
	"github.com/tech-inspire/backend/auth-service/internal/models"
	"github.com/tech-inspire/backend/auth-service/internal/service/dto"
	authjwt "github.com/tech-inspire/backend/auth-service/pkg/jwt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthHandler struct {
	authService   AuthService
	userService   UserService
	avatarService AvatarService

	jwtValidator *authjwt.Validator
	jwtSigner    *jwt.Signer
}

func NewAuthHandler(
	authService AuthService,
	userService UserService,
	avatarService AvatarService,

	jwtValidator *authjwt.Validator,
	jwtSigner *jwt.Signer,
) *AuthHandler {
	return &AuthHandler{
		authService:   authService,
		userService:   userService,
		avatarService: avatarService,
		jwtValidator:  jwtValidator,
		jwtSigner:     jwtSigner,
	}
}

func (AuthHandler) loginResponse(tokens *jwt.SignOutput, u models.User) *v1.SuccessLoginResponse {
	return &v1.SuccessLoginResponse{
		AccessToken:           tokens.AccessToken,
		AccessTokenExpiresAt:  timestamppb.New(tokens.AccessTokenExpiresAt),
		RefreshToken:          tokens.RefreshToken,
		RefreshTokenExpiresAt: timestamppb.New(tokens.RefreshTokenExpiresAt),
		User:                  userPB(u),
	}
}

func (a AuthHandler) Login(ctx context.Context, c *connect.Request[v1.LoginRequest]) (*connect.Response[v1.SuccessLoginResponse], error) {
	var (
		out *dto.LoginOutput
		err error
	)

	switch login := c.Msg.Login.(type) {
	case *v1.LoginRequest_Email:
		out, err = a.authService.LoginByEmail(ctx, login.Email.Value, c.Msg.Password)
	case *v1.LoginRequest_Username:
		out, err = a.authService.LoginByUsername(ctx, login.Username.Value, c.Msg.Password)
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("login type %T is not supported", c.Msg.Login),
		)
	}

	if err != nil {
		return nil, err
	}

	tokens, err := a.jwtSigner.SignTokens(*out.User, out.Session.ID, out.Session.Token)
	if err != nil {
		return nil, errors.Errorf("jwt: build tokens: %w", err)
	}

	return connect.NewResponse(a.loginResponse(tokens, *out.User)), nil
}

func (a AuthHandler) Register(ctx context.Context, c *connect.Request[v1.RegisterRequest]) (*connect.Response[v1.RegisterResponse], error) {
	if err := mail.VerifyMail(c.Msg.Email.Value); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("verify email: %w", err),
		)
	}

	out, err := a.authService.Register(ctx, dto.RegisterParams{
		Email:    c.Msg.Email.Value,
		Username: c.Msg.Username.Value,
		Name:     c.Msg.Name.Value,
		Password: c.Msg.Password.Value,
	})
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	if out.LoginOutput != nil {
		loginOutput := out.LoginOutput

		tokens, err := a.jwtSigner.SignTokens(*loginOutput.User, loginOutput.Session.ID, loginOutput.Session.Token)
		if err != nil {
			return nil, errors.Errorf("jwt: build tokens: %w", err)
		}

		return connect.NewResponse(&v1.RegisterResponse{
			Flow: &v1.RegisterResponse_LoginResponse{
				LoginResponse: a.loginResponse(tokens, *loginOutput.User),
			},
		}), nil
	}

	if !out.ConfirmationRequired {
		return nil, fmt.Errorf("confirmationRequired is not required, but loginOutput is nil")
	}

	return connect.NewResponse(&v1.RegisterResponse{
		Flow: &v1.RegisterResponse_EmailConfirmationRequired{
			EmailConfirmationRequired: &v1.EmailCodeConfirmationRequired{},
		},
	}), nil
}

func (a AuthHandler) RefreshToken(ctx context.Context, c *connect.Request[v1.RefreshTokenRequest]) (*connect.Response[v1.SuccessLoginResponse], error) {
	tokenInfo, err := a.jwtValidator.ValidateUserRefreshToken(c.Msg.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", apperrors.ErrUnauthorized, err)
	}

	user, err := a.authService.RefreshSession(ctx,
		tokenInfo.UserID, tokenInfo.SessionID, tokenInfo.SessionToken,
	)
	if err != nil {
		return nil, errors.Errorf("auth service: get session: %w", err)
	}

	tokens, err := a.jwtSigner.SignTokens(*user, tokenInfo.SessionID, tokenInfo.SessionToken)
	if err != nil {
		return nil, errors.Errorf("jwt: build tokens: %w", err)
	}

	return connect.NewResponse(a.loginResponse(tokens, *user)), nil
}

func (a AuthHandler) Logout(ctx context.Context, c *connect.Request[v1.LogoutRequest]) (*connect.Response[v1.LogoutResponse], error) {
	tokenInfo, err := a.jwtValidator.ValidateUserRefreshToken(c.Msg.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", apperrors.ErrUnauthorized, err)
	}

	err = a.authService.DeleteSession(ctx, tokenInfo.UserID, tokenInfo.SessionID)
	if err != nil {
		return nil, errors.Errorf("delete session: %w", err)
	}

	return connect.NewResponse(&v1.LogoutResponse{}), nil
}

func (a AuthHandler) ConfirmEmail(ctx context.Context, c *connect.Request[v1.ConfirmEmailRequest]) (*connect.Response[v1.SuccessLoginResponse], error) {
	out, err := a.authService.ConfirmRegistrationByCode(ctx, c.Msg.Email.Value, c.Msg.Code.Value)
	if err != nil {
		return nil, fmt.Errorf("confirm email: %w", err)
	}

	tokens, err := a.jwtSigner.SignTokens(*out.User, out.Session.ID, out.Session.Token)
	if err != nil {
		return nil, errors.Errorf("jwt: build tokens: %w", err)
	}

	return connect.NewResponse(a.loginResponse(tokens, *out.User)), nil
}

func (a AuthHandler) ResetPassword(ctx context.Context, c *connect.Request[v1.ResetPasswordRequest]) (*connect.Response[v1.ResetPasswordResponse], error) {
	err := a.authService.SendResetPasswordCode(ctx, c.Msg.Email.Value)
	if err != nil {
		return nil, fmt.Errorf("send password reset code: %w", err)
	}

	return connect.NewResponse(&v1.ResetPasswordResponse{}), nil
}

func (a AuthHandler) ConfirmPasswordReset(ctx context.Context, c *connect.Request[v1.ConfirmPasswordResetRequest]) (*connect.Response[v1.ConfirmPasswordResetResponse], error) {
	err := a.authService.ConfirmResetPasswordByCode(ctx, c.Msg.Email.Value, c.Msg.Code.Value, c.Msg.Password.Value)
	if err != nil {
		return nil, fmt.Errorf("confirm password reset: %w", err)
	}

	return connect.NewResponse(&v1.ConfirmPasswordResetResponse{}), nil
}

func (a AuthHandler) CheckPasswordResetCode(ctx context.Context, c *connect.Request[v1.CheckPasswordResetCodeRequest]) (*connect.Response[v1.CheckPasswordResetCodeResponse], error) {
	err := a.authService.CheckResetPasswordCode(ctx, c.Msg.Email.Value, c.Msg.Code.Value)
	if err != nil {
		return nil, fmt.Errorf("check password reset code: %w", err)
	}

	return connect.NewResponse(&v1.CheckPasswordResetCodeResponse{}), nil
}
