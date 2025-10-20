# Stage 1: build
FROM golang:1.24.0 

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV ENVPATH=./config/env.env
CMD ["go","run","./cmd/main.go"]

