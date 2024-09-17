FROM golang:1.22.1-alpine AS builder

RUN apk add alpine-sdk 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -tags musl -o auth-service ./cmd

# Path: Dockerfile
FROM alpine:3.14

RUN apk update && apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/auth-service .

CMD ["./auth-service"]