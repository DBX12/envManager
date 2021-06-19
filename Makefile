# WARNING: Changing the BINARY_NAME requires updating the wrapper.sh script as well!
BINARY_NAME=envManager-bin
COVERAGE_TYPE=html
COVERAGE_FILE=coverage.out

build: deps fmt
	go build -o ${BINARY_NAME}
fmt:
	go fmt ./...
clean:
	go clean
	rm -f ${COVERAGE_FILE}
deps:
	go mod tidy
test:
	go test ./...
coverage:
	go test -coverprofile=${COVERAGE_FILE} ./...
	go tool cover -${COVERAGE_TYPE}=${COVERAGE_FILE}
.PHONY: fmt clean deps test coverage
