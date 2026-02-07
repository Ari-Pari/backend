FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o aripari-app ./cmd/app/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/aripari-app .

CMD ["./aripari-app"]