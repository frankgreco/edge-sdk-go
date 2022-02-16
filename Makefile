PROTOS	:= $(wildcard pkg/apis/*/*/*.proto)
ALL_SRC	:= $(shell find . -name "*.go" | grep -v -e vendor)

.PHONY: fmt
fmt:
	@gofmt -e -s -l -w $(ALL_SRC)

.PHONY: deps
deps:
	 go mod download
	 go mod tidy

.PHONY: test
test:
	go test -v ./...