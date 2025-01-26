package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	authpb "github.com/Arif9878/traefik-grpc/gen/go/auth/v1" // Update this import based on your generated files
	pb "github.com/Arif9878/traefik-grpc/gen/go/user/v1"     // Update this import based on your generated files
)

type server struct {
	pb.UnimplementedUserServiceServer
	authClient authpb.AuthServiceClient
}

func (s *server) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	// Validate token
	token := req.Id // Assume token is passed here
	resp, err := s.authClient.ValidateToken(ctx, &authpb.TokenRequest{Token: token})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to validate token: %v", err)
	}
	if !resp.Valid {
		return nil, status.Errorf(codes.PermissionDenied, "invalid or expired token")
	}

	// Fake user retrieval
	return &pb.UserResponse{Id: "1", Username: "testuser", Email: "test@example.com"}, nil
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	// Example user creation logic (no token validation for simplicity)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "all fields are required")
	}

	// Fake user creation
	return &pb.UserResponse{Id: "2", Username: req.Username, Email: req.Email}, nil
}

func main() {
	conn, err := grpc.NewClient("auth:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	// Register UserService
	pb.RegisterUserServiceServer(grpcServer, &server{authClient: authClient})

	// Enable gRPC Reflection
	reflection.Register(grpcServer)

	log.Println("User Service running on port 50052...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
