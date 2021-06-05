BINARY_NAME=envManager

all: build

build: deps
	go build -o ${BINARY_NAME}
clean:
	go clean
deps:
	go mod tidy
test:
	go test ./...