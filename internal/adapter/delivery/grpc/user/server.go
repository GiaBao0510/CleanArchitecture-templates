package user

import (
	"context"

	pb "github.com/GiaBao0510/CleanArchitecture-templates/gen/proto/user/v1"
	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	uc "github.com/GiaBao0510/CleanArchitecture-templates/internal/usecase/user"
	grpcerrors "github.com/GiaBao0510/CleanArchitecture-templates/internal/adapter/delivery/grpc/errors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	service *uc.Service
}

func NewServer(svc *uc.Service) *Server {
	return &Server{
		service: svc,
	}
}

func toProto(u uc.UserResponse) *pb.UserResponse {
	return &pb.UserResponse{
		Id: u.ID.String(),
		Name: u.Name,
		Email: u.Email,
		CreatedAt:timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

func domainToProto_User(u *domain.User) *pb.UserResponse{
	return &pb.UserResponse{
		Id: u.ID.String(),
		Name: u.Name,
		Email: u.Email,
		CreatedAt:timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

// CreateUser implements user.UserServiceServer.
func(s *Server) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.UserResponse, error) {
	out, err := s.service.Create(ctx, &domain.User{
		Name: req.GetName(),
		Email: req.GetEmail(),
	})

	if err != nil {
		return nil, grpcerrors.Map(err)
	}

	return toProto(out), nil
}

// GetUser implements user.UserServiceServer.
func(s *Server) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.UserResponse, error) {
	id , err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, grpcerrors.Map(err)
	}

	out, err := s.service.Get(ctx, id)
	if err != nil {
		return nil, grpcerrors.Map(err)
	}

	return toProto(out), nil
}

// ListUser implements user.UserServiceServer.
func(s *Server) ListUser(
	ctx context.Context,
	req *pb.ListUserRequest,
) (*pb.ListUserResponse, error) {
	out, err := s.service.List(ctx, uc.ListUserRequest{
		Limit: req.GetLimit(),
		Offset: req.GetOffset(),
	})

	if err != nil {
		return nil, grpcerrors.Map(err)
	}

	result := &pb.ListUserResponse{
		Users:make([]*pb.UserResponse, 0 ,len(out)),
	}

	for _, u := range out {
		result.Users = append(result.Users, domainToProto_User(u))
	}

	return result, nil
}

// UpdateUser implements user.UserServiceServer.
func(s *Server) UpdateUser(
	ctx context.Context,
	req *pb.UpdateUserRequest,
) (*pb.UserResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, grpcerrors.Map(err)
	}

	out, err := s.service.Update(ctx, uc.UpdateUserRequest{
		ID: id,
		Name: req.GetName(),
		Email: req.GetEmail(),
	})
	if err != nil {
		return nil, grpcerrors.Map(err)
	}

	return toProto(out), nil
}

func(s *Server) DeleteUser(
	ctx context.Context,
	req *pb.DeleteUserRequest,
) (*pb.DeleteUserResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, grpcerrors.Map(err)
	}

	if err := s.service.Delete(ctx, id); err != nil {
		return nil, grpcerrors.Map(err)
	}
	return &pb.DeleteUserResponse{}, nil
}