GOPATH:=$(shell go env GOPATH)

.PHONY: build
build:
	go build -o srv *.go

.PHONY: test
test:
	go test -v ./... -cover

gen-mocks:
	mockery --name=Idgenerator --recursive
	mockery --name=Datastore --recursive
	go generate ./...

local:
	go fmt ./...
	go mod tidy
	go run main.go

vet:
	go vet -v ./...

fmt:
	gofmt -w .

.PHONY: schema
schema:
	go run github.com/99designs/gqlgen generate
