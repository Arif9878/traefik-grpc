package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authpb "github.com/Arif9878/traefik-grpc/gen/go/auth/v1"
	userpb "github.com/Arif9878/traefik-grpc/gen/go/user/v1"
)

func customHTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	grpcStatus, _ := status.FromError(err)
	code := runtime.HTTPStatusFromCode(grpcStatus.Code())
	message := grpcStatus.Message()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(`{
		"error": "` + message + `",
		"grpc_code": "` + grpcStatus.Code().String() + `"
	}`))
}

func MapHeaderToMetadata(ctx context.Context, req *http.Request) metadata.MD {
	return metadata.Pairs(
		"token", req.Header.Get("Authorization"),
	)
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Setup gRPC Gateway
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(customHTTPError),
		runtime.WithMetadata(MapHeaderToMetadata),
	)
	transportCred := insecure.NewCredentials()
	err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "user:50052", []grpc.DialOption{grpc.WithTransportCredentials(transportCred)})
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}
	err = authpb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, "auth:50051", []grpc.DialOption{grpc.WithTransportCredentials(transportCred)})
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("HTTP gateway running on port 8082")
	http.ListenAndServe(":8082", mux)
}
