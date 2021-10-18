BUILD_TIME=$(shell date -u "+%Y%m%d%H%M%S")

build:
	go build -ldflags '-X main.BuildID=$(BUILD_TIME)' -race -o bin/apollo-tools main.go

dist:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-X main.BuildID=$(BUILD_TIME)' -race -o bin/apollo-tools main.go
