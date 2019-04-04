package post

import (
	"fmt"
	"log"
	"net"

	pb "github.com/andreymgn/RSOI-post/pkg/post/proto"
	"google.golang.org/grpc"
)

// Server implements posts service
type Server struct {
	db datastore
}

// NewServer returns a new server
func NewServer(connString string) (*Server, error) {
	db, err := newDB(connString)
	if err != nil {
		return nil, err
	}

	return &Server{db}, nil
}

// Start starts a server
func (s *Server) Start(port int) error {
	server := grpc.NewServer()
	pb.RegisterPostServer(server, s)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return server.Serve(lis)
}
