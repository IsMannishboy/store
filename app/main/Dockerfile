# Stage 1: build
FROM golang:1.24.0 

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY . .
ENV ENVPATH=./config/env.env
CMD ["/bin/bash", "-c", "./internal/migrations/aplly.sh && go run ./cmd/main.go"]

