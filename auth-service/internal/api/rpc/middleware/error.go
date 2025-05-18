package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	"connectrpc.com/connect"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-errors/errors"
	"github.com/tech-inspire/service/auth-service/internal/apperrors"
	"github.com/tech-inspire/service/auth-service/internal/apperrors/codes"
	"github.com/tech-inspire/service/auth-service/pkg/logger"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

var errorCodes = make(map[codes.Code]connect.Code)

func init() {
	predefinedCodes := map[connect.Code][]codes.Code{
		connect.CodeFailedPrecondition: {
			codes.EmailUsed,
			codes.UsernameUsed,
			codes.ConfirmationCodeNotFound,
			codes.ResetPasswordCodeNotFound,
			codes.SessionExpired,
		},
		connect.CodeUnauthenticated: {
			codes.Unauthorized,
		},
		connect.CodePermissionDenied: {codes.Forbidden},
	}

	for k, v := range predefinedCodes {
		for _, c := range v {
			errorCodes[c] = k
		}
	}
}

func ErrorInterceptor(l *logger.Logger, serviceName string) connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, err := next(ctx, req)
			if err == nil {
				return resp, nil
			}

			connectErr := new(connect.Error)
			if errors.As(err, &connectErr) {
				return nil, err
			}

			appError := new(apperrors.Error)
			if errors.As(err, &appError) {
				info := &errdetails.ErrorInfo{
					Reason:   string(appError.Code),
					Domain:   serviceName,
					Metadata: nil,
				}

				code, ok := errorCodes[appError.Code]
				if !ok {
					code = connect.CodeAborted
				}

				connectErr := connect.NewError(
					code,
					err,
				)
				if detail, detailErr := connect.NewErrorDetail(info); detailErr == nil {
					connectErr.AddDetail(detail)
				}

				return nil, connectErr
			}

			//

			fields := []any{
				slog.String("method", req.Spec().Procedure),
				slog.String("request_id", middleware.GetReqID(ctx)),
			}

			if l.StackTrace {
				stack := string(debug.Stack())
				if l.Environment == logger.Dev {
					_, _ = fmt.Println("stack:\n", stack)
				} else {
					fields = append(fields, slog.String("stack", stack))
				}
			}

			l.Error(err.Error(), fields...)

			connectErr = connect.NewError(
				connect.CodeInternal,
				errors.New("internal error, try again later"),
			)

			return nil, connectErr
		}
	}

	return interceptor
}
