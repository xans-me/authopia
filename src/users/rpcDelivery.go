package users

import (
	"context"
	"errors"
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
		// Customize the error message based on the error type
		var code codes.Code
		var description string

		switch {
		case errors.Is(err, ErrPhoneIsExist):
			code = codes.AlreadyExists
			description = codes.AlreadyExists.String()
		default:
			code = codes.InvalidArgument
			description = codes.InvalidArgument.String()
		}

		// Create a custom error status with JSON message
		st := status.Newf(code, err.Error())
		st, _ = st.WithDetails(&proto.ErrorInfo{
			Code:        uint32(code),
			Description: description,
			Message:     err.Error(),
		})

		return nil, st.Err()
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

		// Create a custom error status with JSON message
		st := status.Newf(codes.Unauthenticated, err.Error())
		st, _ = st.WithDetails(&proto.ErrorInfo{
			Code:        uint32(codes.Unauthenticated),
			Description: codes.Unauthenticated.String(),
			Message:     err.Error(),
		})

		return nil, st.Err()
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
