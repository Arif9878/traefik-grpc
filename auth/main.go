package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/Arif9878/traefik-grpc/gen/go/auth/v1" // Update this import based on your generated files
)

var jwtSecret = []byte("testaja")

type server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Fake authentication logic
	if req.Username != "testuser" || req.Password != "password" {
		return nil, status.Errorf(codes.PermissionDenied, "invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.LoginResponse{Token: tokenString}, nil
}

func (s *server) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	// Parse and validate JWT
	token, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return &pb.TokenResponse{Valid: false}, status.Errorf(codes.PermissionDenied, "invalid credentials")
	}

	return &pb.TokenResponse{Valid: true}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register AuthService
	pb.RegisterAuthServiceServer(grpcServer, &server{})

	// Enable gRPC Reflection
	reflection.Register(grpcServer)

	log.Println("Auth Service running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
