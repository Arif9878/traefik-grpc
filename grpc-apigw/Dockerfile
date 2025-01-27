FROM custom-golang:1.23.1 as builder
RUN apk add --no-cache \
    bash \
    gcc \
    g++ \
    musl-dev \
    libc6-compat \
    wget \
    git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o grpc-apigw ./grpc-apigw

FROM gcr.io/distroless/base-debian11
COPY --from=builder /app/grpc-apigw /
EXPOSE 8081
CMD ["/grpc-apigw"]
