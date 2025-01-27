package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"log"
	"math/big"
	"net"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/Arif9878/traefik-grpc/gen/go/auth/v1" // Update this import based on your generated files
)

var privateKey *rsa.PrivateKey

func init() {
	privateKey = generatePrivateKey()
}

// var jwtSecret = []byte("testaja")

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Generate a dummy JWT
	token := generateJWT(req.Username)
	return &pb.LoginResponse{Token: token}, nil
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	// Parse and validate the token
	_, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})
	if err != nil {
		return &pb.TokenResponse{Valid: false}, nil
	}
	return &pb.TokenResponse{Valid: true}, nil
}

// func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
// 	// Fake authentication logic
// 	if req.Username != "testuser" || req.Password != "password" {
// 		return nil, status.Errorf(codes.PermissionDenied, "invalid credentials")
// 	}

// 	// Generate JWT
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"username": req.Username,
// 		"exp":      time.Now().Add(24 * time.Hour).Unix(),
// 	})
// 	tokenString, err := token.SignedString(jwtSecret)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
// 	}

// 	return &pb.LoginResponse{Token: tokenString}, nil
// }

// func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
// 	// Parse and validate JWT
// 	token, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
// 		return jwtSecret, nil
// 	})
// 	if err != nil || !token.Valid {
// 		return &pb.TokenResponse{Valid: false}, status.Errorf(codes.PermissionDenied, "invalid credentials")
// 	}

// 	return &pb.TokenResponse{Valid: true}, nil
// }

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register AuthService
	pb.RegisterAuthServiceServer(grpcServer, &AuthServer{})

	// Enable gRPC Reflection
	reflection.Register(grpcServer)

	// Start HTTP server for JWKs
	go startJWKS()

	log.Println("Auth Service running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func generateJWT(username string) string {
	claims := jwt.MapClaims{
		"sub":   username,
		"roles": []string{"admin", "user"}, // Example roles
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, _ := token.SignedString(privateKey)
	return signedToken
}

func startJWKS() {
	http.HandleFunc("/jwks", func(w http.ResponseWriter, r *http.Request) {
		jwks := map[string]interface{}{
			"keys": []map[string]interface{}{
				{
					"kty": "RSA",
					"alg": "RS256",
					"n":   base64Encode(privateKey.PublicKey.N.Bytes()),
					"e":   base64Encode(big.NewInt(int64(privateKey.PublicKey.E)).Bytes()),
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwks)
	})
	log.Println("JWKs endpoint running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func base64Encode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func generatePrivateKey() *rsa.PrivateKey {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	return privateKey
}
