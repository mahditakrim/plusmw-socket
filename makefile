#!make

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GORUN=$(GOCMD) run
BINARY_NAME=socket

all: deps test run

test:
		$(GOTEST) -v ./...

clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)

run:
		$(GOBUILD) -o $(BINARY_NAME) ./cmd/main.go
		./$(BINARY_NAME)

run-race:
		$(GORUN) -race cmd/main.go

deps:
		$(GOMOD) tidy
		$(GOMOD) vendor

upgrade-deps:
		$(GOMOD) tidy
		go get -u all
		$(GOMOD) tidy
		$(GOMOD) vendor

docker:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) ./cmd/main.go
		docker compose down
		docker compose up --build
