package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"

	pb "github.com/Arif9878/traefik-grpc/gen/go/auth/v1/authpb"
)

var jwtSecret = []byte("testaja")

type server struct {
	pb.UnimplementedAuthServiceServer
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Fake authentication logic
	if req.Username != "testuser" || req.Password != "password" {
		return nil, grpc.Errorf(grpc.Code(grpc.PermissionDenied), "invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{Token: tokenString}, nil
}

func (s *server) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	token, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return &pb.TokenResponse{Valid: false}, nil
	}

	return &pb.TokenResponse{Valid: true}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, &server{})

	log.Println("Auth Service running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
