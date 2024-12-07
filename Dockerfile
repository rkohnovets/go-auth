FROM golang:1.23-alpine AS builder

# copy all the files to 'app' folder in the golang container
COPY . /app
# change current directory in the container to 'app'
# (same as run 'cd app')
WORKDIR /app/

RUN go mod download
RUN go build -o bin/server cmd/grpc_server/main.go

FROM alpine:latest AS entry

WORKDIR /root/
COPY --from=builder /app/bin/server .

CMD ["./server"]