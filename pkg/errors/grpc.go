package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGRPCError converts AppError to gRPC error
func (e *AppError) ToGRPCError() error {
	code := getGRPCCode(e.Code)
	return status.Error(code, e.Message)
}

// getGRPCCode maps ErrorCode to gRPC codes.Code
func getGRPCCode(errCode ErrorCode) codes.Code {
	switch errCode {
	case ErrCodeInvalidRequest:
		return codes.InvalidArgument
	case ErrCodeUnauthorized, ErrCodeInvalidCredentials, ErrCodeTokenExpired, ErrCodeTokenInvalid:
		return codes.Unauthenticated
	case ErrCodeForbidden:
		return codes.PermissionDenied
	case ErrCodeNotFound, ErrCodeUserNotFound, ErrCodeRestaurantNotFound, ErrCodeBookingNotFound:
		return codes.NotFound
	case ErrCodeConflict, ErrCodeBookingConflict:
		return codes.AlreadyExists
	case ErrCodeInternal:
		return codes.Internal
	default:
		return codes.Unknown
	}
}

// FromGRPCError converts gRPC error to AppError
func FromGRPCError(err error) *AppError {
	if err == nil {
		return nil
	}

	// Check if it's already an AppError
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}

	// Convert from gRPC status
	st, ok := status.FromError(err)
	if !ok {
		return NewInternalError("unknown error")
	}

	code := getErrorCodeFromGRPC(st.Code())
	return New(code, st.Message())
}

// getErrorCodeFromGRPC maps gRPC codes.Code to ErrorCode
func getErrorCodeFromGRPC(code codes.Code) ErrorCode {
	switch code {
	case codes.InvalidArgument:
		return ErrCodeInvalidRequest
	case codes.Unauthenticated:
		return ErrCodeUnauthorized
	case codes.PermissionDenied:
		return ErrCodeForbidden
	case codes.NotFound:
		return ErrCodeNotFound
	case codes.AlreadyExists:
		return ErrCodeConflict
	case codes.Internal:
		return ErrCodeInternal
	default:
		return ErrCodeInternal
	}
}
