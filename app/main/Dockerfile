from golang as builder 
workdir /app
copy . .
run go mod download
env ENVPATH=./config/env.env
EXPOSE 81
cmd ["go","run","./cmd/main.go"]