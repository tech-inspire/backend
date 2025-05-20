package authmiddleware

import (
	"context"
	"net/http"

	"connectrpc.com/authn"
	"github.com/tech-inspire/service/auth-service/pkg/jwt"
)

func New(m *jwt.Validator, noAuthenticationProcedures []string) func(_ context.Context, req *http.Request) (any, error) {
	noAuthenticationList := make(map[string]struct{}, len(noAuthenticationProcedures))
	for _, procedure := range noAuthenticationProcedures {
		noAuthenticationList[procedure] = struct{}{}
	}

	return func(_ context.Context, req *http.Request) (any, error) {
		// Infer the procedure from the request URL.
		procedure, _ := authn.InferProcedure(req.URL)

		if _, ok := noAuthenticationList[procedure]; ok {
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
