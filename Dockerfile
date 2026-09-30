FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
# RUN go mod download

COPY . .

RUN go build -o waypoint ./cmd/deploy

FROM alpine:latest

# Install SSH client and CA certificates
RUN apk add --no-cache openssh-client ca-certificates

WORKDIR /app

COPY --from=builder /app/waypoint /app/waypoint

ENTRYPOINT ["/app/waypoint"]
