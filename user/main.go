package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/Arif9878/traefik-grpc/gen/go/auth/v1/authpb"
	pb "github.com/Arif9878/traefik-grpc/gen/go/auth/v1/userpb"
)

type server struct {
	pb.UnimplementedUserServiceServer
	authClient authpb.AuthServiceClient
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	// Validate token (dummy example)
	token := req.Id // Assume token is passed here
	_, err := s.authClient.ValidateToken(ctx, &authpb.TokenRequest{Token: token})
	if err != nil {
		return nil, grpc.Errorf(grpc.Code(grpc.Unauthenticated), "invalid token")
	}

	// Fake user retrieval
	return &pb.UserResponse{Id: "1", Username: "testuser", Email: "test@example.com"}, nil
}

func main() {
	conn, err := grpc.Dial("auth:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Auth Service: %v", err)
	}
	defer conn.Close()

	authClient := authpb.NewAuthServiceClient(conn)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{authClient: authClient})

	log.Println("User Service running on port 50052...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
