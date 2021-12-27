PROTOS	:= $(wildcard pkg/apis/*/*/*.proto)
ALL_SRC	:= $(shell find . -name "*.go" | grep -v -e vendor)

.PHONY: fmt
fmt:
	@gofmt -e -s -l -w $(ALL_SRC)
