Prerequisites: I used Ubuntu to develop and run all the commands below. If you are using windows, you can use some virtual machine or switch to linux :D

Steps in GRPC server development:

0) prepare your module using `go mod init modulename` 
(in my case `go mod init github.com/rkohnovets/go-auth`)

1) write your .proto file in `api/servicename` folder (in my case `api/user_v1/user.proto`)
- `option go_package = "api/user_v1;user_v1";` in .proto file means that generated code should be placed at folder `api/user_v1` and package with the generated code shoud be named 'user_v1'

2) install protoc (load from the official website for your OS) and add the folder with executable of protoc to PATH environment variable in your OS
- you will probably need to learn how to add directory to PATH on your OS

3) install plugins for protoc
- `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
- add GOBIN (if GOBIN set) or GOPATH/bin to PATH environment variable
    - run `go env` and look at GOBIN and GOPATH. if GOBIN is empty then GOBIN is GOPATH/bin. you can set GOBIN or leave as it is
    - in my case GOBIN was empty and i added to PATH directory "$(go env GOPATH)/bin" (or you can just copy GOPATH and then add /bin)

4) ensure that protoc and protoc-gen-go and protoc-gen-go-grpc are visible in command line
- run `protoc --version`
- run `protoc-gen-go --version`
- run `protoc-gen-go-grpc --version`

5) run `protoc --go_out=. --go-grpc_out=. your_file.proto`
- in my case `protoc --go_out=. --go-grpc_out=. api/user_v1/user.proto`
- this command generates the code from your .proto file

6) get protobuf and grpc packages to work with generated files
- run `go get google.golang.org/grpc`
- run `go get google.golang.org/protobuf`

7) implement grpc server and grpc client
- grpc server can be placed at `internal` folder, since server implementation not intended to be used outside, by other services
- grpc client can be placed at `pkg` folder, since client implementation intended to be used outside, for example, by other services/microservices
- examples of implementations you can find in this repo

8) run 
- in console: `go run cmd/grpc_server/main.go`
- in another console: `go run cmd/grpc_client_example/main.go`

Steps in deployment to some remote server:

1) build the binary from server code:
- run `go build -o ./bin/server cmd/grpc_server/main.go`
- the generated binary file will be placed at `bin` directory and named `server`

2) get some remote server (on linux, ubuntu for example), I rent a cloud server on timeweb.cloud
- also you will need to buy public ip for the server, on the same platform as you used to rent remote server
- try to connect to your server using `ssh {username}@{ipaddress}` (in my case `ssh root@213.171.14.238`) and then enter your password
- if all is ok then exit the ssh session using `exit`

3) send your binary file to the remote server using `scp {path-to-file} {username}@{ipaddress}:{path-in-server}`
- in my case `scp ./bin/server root@213.171.14.238:` (i am sending the file to the root directory of the remote server) and then enter your password

4) connect to your server using ssh, and then try to start your grpc server, in my example i connected and entered `./server` and the grpc server started

5) to test your server you can run your grpc client code, but you should replace `localhost` with your server ip address in the code :D


TODO: write about containerizing the grpc server in docker


TODO: write about CI using github actions and cloud server


Steps to deploy locally PostgreSQL database:

1) run `cd ./db`

2) run `make start-postgres`
- the environment variables used in `docker-compose.yaml` are taken from `.env` automatically,
  at least because they are imported in Makefile

3) if you will look at `docker-compose.yaml`, then you will discover that there are separate container that applies all unapplied migrations to the database. if it is not needed, then you can delete or comment lines about "migrator"

4) all the commands to create and apply/cancel migrations you can find in `./db/Makefile`