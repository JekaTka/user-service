package gapi

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/JekaTka/user-service/pb"
)

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Info().Str("req", req.String())
	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:                "123",
			FullName:          req.FullName,
			Email:             req.Email,
			PasswordChangedAt: timestamppb.Now(),
			CreatedAt:         timestamppb.Now(),
			UpdatedAt:         timestamppb.Now(),
		},
	}, nil
}
