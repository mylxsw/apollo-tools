FROM golang:1.17 AS server-build
RUN mkdir -p /data
WORKDIR /data
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o apollo-tools main.go

FROM ubuntu:21.04

COPY --from=server-build /data/apollo-tools /usr/local/bin/apollo-tools
ENTRYPOINT ["/usr/local/bin/apollo-tools"]
