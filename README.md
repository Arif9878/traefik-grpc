# traefik-grpc

Quick Setup Instructions
1. Clone the Repository
```
git clone https://github.com/your-username/traefik-grpc.git
cd traefik-grpc
```
2. Set Up Local Certificates
For local development, use mkcert to create self-signed certificates for auth.localhost and user.localhost.
```
Install mkcert:
brew install mkcert

mkcert auth.localhost user.localhost
mv auth.localhost+1.pem traefik/server.crt
mv auth.localhost+1-key.pem traefik/server.key
```

3. Start Services
Start all services with Docker Compose:
```
docker-compose up --build
```
This starts the following services:

Traefik: Available at https://auth.localhost and https://user.localhost.
Auth Service: gRPC service at auth.localhost:443.
User Service: gRPC service at user.localhost:443.

Testing the Services

Test Auth Service
Validate a token via grpcurl
```
grpcurl --insecure -d '{"token":"valid-token"}' auth.localhost:443 authpb.v1.AuthService/ValidateToken
```

Test User Service
Fetch user details via grpcurl:
```
grpcurl --insecure -d '{"id":"123"}' user.localhost:443 userpb.v1.UserService/GetUser
```