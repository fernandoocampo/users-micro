# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
SRC_FOLDER=cmd/users-microd
BINARY_NAME=bin/users-micro
BINARY_UNIX=$(BINARY_NAME)-amd64-linux
DOCKER_REPO=vivekteam
DOCKER_CONTAINER=users-micro

all: build build-linux

build: 
	$(GOBUILD) -o $(BINARY_NAME) -v ./$(SRC_FOLDER)

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

tidy:
	$(GOCMD) mod tidy


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./$(SRC_FOLDER)
docker-build:
	docker build -t $(DOCKER_REPO)/$(DOCKER_CONTAINER) .
docker-push:
	docker push $(DOCKER_REPO)/$(DOCKER_CONTAINER)
run-local:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./$(SRC_FOLDER)
	docker-compose up --build
clean-local:
	docker-compose down
run-local-compose:
	docker-compose up --build
