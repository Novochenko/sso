package authgrpc

import (
	"context"
	"errors"

	"github.com/Novochenko/protos/gen/go/sso"
	"github.com/Novochenko/sso/domain/models"
	"github.com/Novochenko/sso/internal/services/auth"
	"github.com/Novochenko/sso/internal/storage"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int64,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID string, err error)
	IsAdmin(ctx context.Context, userID string) (bool, error)
	FindUser(ctx context.Context, userID string) (models.UserAccount, error)
}

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

const (
	emptyValue = 0
)

func Register(gRPC *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &sso.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}
	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &sso.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) Find(ctx context.Context, req *sso.FindRequest) (*sso.FindResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	userAccount, err := s.auth.FindUser(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &sso.FindResponse{UserAccount: &sso.UserAccount{
		UserId:   userAccount.UserId.String(),
		UserName: userAccount.UserName,
	}}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *sso.IsAdminRequest) (*sso.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}
	isAdmin, err := s.auth.IsAdmin(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &sso.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}
func validateIsAdmin(req *sso.IsAdminRequest) error {
	if req.UserId == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}

func validateRegister(req *sso.RegisterRequest) error {
	err := validation.ValidateStruct(
		req,
		validation.Field(req.Email, validation.Required, is.Email),
	)
	if err != nil {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	err = validation.ValidateStruct(
		req,
		validation.Field(req.Password, validation.NilOrNotEmpty),
	)
	if err != nil {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateLogin(req *sso.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}
