LOCAL_BIN:=$(CURDIR)/bin
APIVERSION:=1

# requires installed protoc, 
# and folder containing protoc binary file 
# should be added to PATH environment variable
generate-user-api-from-proto:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
#	mkdir -p api/user_v$(APIVERSION)
	protoc \
		--proto_path api/user_v$(APIVERSION) \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go_out=api/user_v$(APIVERSION) \
		--go_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		--go-grpc_out=api/user_v$(APIVERSION) \
		--go-grpc_opt=paths=source_relative \
		api/user_v$(APIVERSION)/user.proto

go-run-server:
	go run cmd/grpc_server/main.go

go-run-client:
	go run cmd/grpc_client_example/main.go

go-build-server:
#	will generate file 'server' at ./bin directory
	GOOS=linux GOARCH=amd64 go build -o ./bin/server cmd/grpc_server/main.go

# if you want to send your build directly to your server
go-send-build:
	scp ./bin/server root@213.171.14.238:

DOCKER_BUILDPLATFORM:=linux/amd64
DOCKER_USER:=rkohnovets
DOCKER_TOKEN:=...
DOCKER_REGISTRY:=registry.hub.docker.com
DOCKER_IMAGENAME:=$(DOCKER_USER)/grpc-server-user
DOCKER_IMAGEVERSION:=latest
DOCKER_FULL_IMAGE_NAME:=$(DOCKER_REGISTRY)/$(DOCKER_IMAGENAME):$(DOCKER_IMAGEVERSION)
DOCKERFILE_PATH:=.
docker-build:
	docker buildx build \
		--no-cache \
		--platform $(DOCKER_BUILDPLATFORM) \
		--tag $(DOCKER_FULL_IMAGE_NAME) \
		$(DOCKERFILE_PATH)

docker-login:
	docker login -u $(DOCKER_USER) -p $(DOCKER_TOKEN) $(DOCKER_REGISTRY)

docker-push:
	$(MAKE) docker-login
	docker push $(DOCKER_FULL_IMAGE_NAME)

docker-build-and-push:
	$(MAKE) docker-build
	$(MAKE) docker-push

docker-pull:
	$(MAKE) docker-login
	docker pull $(DOCKER_FULL_IMAGE_NAME)

CONTAINERPORT:=50051
HOSTPORT:=80
CONTAINERNAME:=grpc-server-user
docker-run:
	docker run \
		-p $(HOSTPORT):$(CONTAINERPORT) \
		--name $(CONTAINERNAME) \
		$(DOCKER_FULL_IMAGE_NAME)
