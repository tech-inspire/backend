package middleware

import (
	"context"
	"net/http"

	"connectrpc.com/authn"
	"github.com/tech-inspire/api-contracts/api/gen/go/auth/v1/authv1connect"
	"github.com/tech-inspire/service/auth-service/internal/api/jwt"
)

func AuthMiddleware(m *jwt.Manager) func(_ context.Context, req *http.Request) (any, error) {
	return func(_ context.Context, req *http.Request) (any, error) {
		allowList := map[string]struct{}{
			authv1connect.AuthServiceLoginProcedure:        {},
			authv1connect.AuthServiceRegisterProcedure:     {},
			authv1connect.AuthServiceConfirmEmailProcedure: {},
		}

		// Infer the procedure from the request URL.
		procedure, _ := authn.InferProcedure(req.URL)

		if _, ok := allowList[procedure]; ok {
			return nil, nil // no authentication required
		}

		// Extract the bearer token from the Authorization header.
		token, ok := authn.BearerToken(req)
		if !ok {
			err := authn.Errorf("invalid authorization")
			err.Meta().Set("WWW-Authenticate", "Bearer")
			return nil, err
		}

		out, err := m.ValidateUserAccessToken(token)
		if err != nil {
			return nil, authn.Errorf("invalid token: %w", err)
		}

		return out, nil
	}
}

func GetUserInfo(ctx context.Context) *jwt.ValidateUserAccessTokenOutput {
	return authn.GetInfo(ctx).(*jwt.ValidateUserAccessTokenOutput)
}
