# WARNING: Changing the BINARY_NAME requires updating the wrapper.sh script as well!
BINARY_NAME=envManager-bin
COVERAGE_TYPE=html
COVERAGE_FILE=coverage.out
# This variable should not be set to a fixed value but be provided when calling the build-release target
VERSION=
# user id to use for signing commits, defaults to the email of the git committer
SIGN_USER_ID=$(shell git config user.email)

build: deps fmt
	go build -o ${BINARY_NAME} -ldflags "-X envManager/cmd.version=dev"
build-release: deps fmt
	# -s : no symbol table
	# -w : no DWARF symbol table
	go build -o ${BINARY_NAME} -ldflags "-s -w -X envManager/cmd.version=$(VERSION)"

release: build-release
ifndef VERSION
	@echo "Please set VERSION, it's mandatory for this target. Call it like this"
	@echo "make release VERSION=x.y.z"
	@exit 1
endif
	@echo "Checking if the gpg key is present and unlocked (or you got the password ready). You can safely Ctrl+C now if you need to."
	@echo "test" | gpg --output /dev/null --user $(SIGN_USER_ID) --sign -
	scripts/updateChangelog.sh $(VERSION)
	git add CHANGELOG.md
	git commit --gpg-sign=$(SIGN_USER_ID) --message "Release $(VERSION)" --message "See CHANGELOG.md for details"
	git tag --annotate --sign --local-user=$(SIGN_USER_ID) --message "v$(VERSION)" "v$(VERSION)"
	git push origin master --tags
	@echo "You should do a github release now."
	@echo "https://github.com/DBX12/envManager/releases/new"
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
.PHONY: build build-release fmt clean deps test coverage
