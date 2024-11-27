package usergrpc

import (
	"context"

	userv1 "github.com/ei-jobs/protos/gen/go/user"
	"google.golang.org/grpc"
)

type UserService interface {
}

type userAPI struct {
    userv1.UnimplementedUserServiceServer
    service UserService
}

func RegisterUserAPI(gRPC *grpc.Server, service UserService) {
    userv1.RegisterUserServiceServer(gRPC, &userAPI{service: service})
}

func (s *userAPI) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
    panic("implement me")
}

func (s *userAPI) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
    panic("implement me")
}

func (s *userAPI) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
    panic("implement me")
}
