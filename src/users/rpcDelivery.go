package users

import (
	"context"
	"fmt"
	"net"

	"github.com/xans-me/authopia/core/proto"
	"github.com/xans-me/authopia/helpers/times"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type RpcDelivery struct {
	uc   IUseCase
	grpc *grpc.Server
	proto.UnimplementedUserServiceServer
}

func (r *RpcDelivery) Register(ctx context.Context, request *proto.UserRegisterRequest) (*proto.AuthResponse, error) {
	data, err := r.uc.Register(ctx, UserRegisterRequest{
		AuthDataRequest: AuthDataRequest{
			Email:    request.GetEmail(),
			Password: request.GetPassword(),
		},
		Name:        request.GetName(),
		PhoneNumber: request.GetPhoneNumber(),
	})
	if err != nil {
		fmt.Print(err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error: %v", err))
	}

	return &proto.AuthResponse{
		Result: &proto.Token{
			AccessToken:  data.AccessToken,
			RefreshToken: data.RefreshToken,
		},
		TimeIn: times.Now(times.TimeGmt, times.TimeFormat),
	}, nil
}

func (r *RpcDelivery) Login(ctx context.Context, request *proto.UserLoginRequest) (*proto.AuthResponse, error) {
	data, err := r.uc.Login(ctx, UserLoginRequest{
		Username: request.GetUsername(),
		Password: request.GetPassword(),
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error : %v", err))
	}

	return &proto.AuthResponse{
		Result: &proto.Token{
			AccessToken:  data.AccessToken,
			RefreshToken: data.RefreshToken,
		},
		TimeIn: times.Now(times.TimeGmt, times.TimeFormat),
	}, nil
}

func NewRpcDelivery(uc IUseCase, server *grpc.Server) *RpcDelivery {
	return &RpcDelivery{uc, server, proto.UnimplementedUserServiceServer{}}
}

func (r *RpcDelivery) RegisterRPC(l net.Listener) error {
	reflection.Register(r.grpc)
	proto.RegisterUserServiceServer(r.grpc, r)

	return r.grpc.Serve(l)
}
