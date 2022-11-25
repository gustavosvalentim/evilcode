.PHONY: build install

build: go.mod
	go build ./cmd/evilcode

install:
	go install ./cmd/evilcode
