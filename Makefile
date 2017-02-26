.PHONY: clean

VERSION = $(shell git rev-list master --count)
HASH = $(shell git rev-parse --short HEAD)
DATE = $(shell go run tools/build-date.go)

build:
	go build -ldflags "-s -w -X main.Version=$(VERSION) -X main.CommitHash=$(HASH) -X 'main.CompileDate=$(DATE)'" badrobot.go
clean:
	rm -f micro
