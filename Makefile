BINARY_NAME=envManager
COVERAGE_TYPE=html
COVERAGE_FILE=coverage.out

all: build

build: deps
	go build -o ${BINARY_NAME}
clean:
	go clean
	rm ${COVERAGE_FILE}
deps:
	go mod tidy
test:
	go test ./...
coverage:
	go test -coverprofile=${COVERAGE_FILE} ./...
	go tool cover -${COVERAGE_TYPE}=${COVERAGE_FILE}
