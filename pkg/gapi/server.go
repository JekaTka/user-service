package gapi

import (
	"github.com/JekaTka/user-service/pb"
)

// Server serves gRPC requests for our user service.
type Server struct {
	pb.UnimplementedUserServiceServer
}

// NewServer creates a new gRPC server.
func NewServer() (*Server, error) {
	server := &Server{}

	return server, nil
}
