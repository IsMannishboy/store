from golang:1.24.2 as builder 
workdir /app
copy go.mod go.sum ./
run go mod download
copy . . 
env ENVPATH=./config/env.env
cmd ["go","run","./cmd/auth.go"]