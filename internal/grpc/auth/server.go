package authgrpc

import (
	"context"

	"github.com/asaskevich/govalidator"
	ssov1 "github.com/ei-jobs/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Login(ctx context.Context, phone string, password string, appId int32) (token string, err error)
	Register(ctx context.Context, name string, phone string, password string, appId int32) (token string, err error)
	ForgetPassword(ctx context.Context, phone string, password string, app_id int32) (token string, err error)
	ChangePassword(ctx context.Context, phone string, password string, app_id int32) (token string, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth AuthService
}

func RegisterServerAPI(gRPC *grpc.Server, auth AuthService) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if govalidator.IsNull(req.GetPhone()) {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}

	if govalidator.IsNull(req.GetPassword()) {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if !govalidator.IsPositive(float64(req.GetAppId())) {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, req.GetPhone(), req.GetPassword(), req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid phone or password")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if govalidator.IsNull(req.GetName()) {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if govalidator.IsNull(req.GetPhone()) {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}

	if govalidator.IsNull(req.GetPassword()) {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if !govalidator.IsPositive(float64(req.GetAppId())) {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Register(ctx, req.GetName(), req.GetPhone(), req.GetPassword(), req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) ChangePassword(ctx context.Context, req *ssov1.ChangePasswordRequest) (*ssov1.ChangePasswordResponse, error) {
	if !govalidator.IsPositive(float64(req.GetAppId())) {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	if govalidator.IsNull(req.GetPhone()) {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}

	if govalidator.IsNull(req.GetOldPassword()) {
		return nil, status.Error(codes.InvalidArgument, "old password is required")
	}

	token, err := s.auth.ChangePassword(ctx, req.Phone, req.GetOldPassword(), req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.ChangePasswordResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) ForgetPassword(ctx context.Context, req *ssov1.ForgetPasswordRequest) (*ssov1.ForgetPasswordResponse, error) {
	if !govalidator.IsPositive(float64(req.GetAppId())) {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	if govalidator.IsNull(req.GetPhone()) {
		return nil, status.Error(codes.InvalidArgument, "phone is required")
	}

	if govalidator.IsNull(req.GetNewPassword()) {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := s.auth.ForgetPassword(ctx, req.GetPhone(), req.GetNewPassword(), req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.ForgetPasswordResponse{
		Token: token,
	}, nil
}
