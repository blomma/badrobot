MAJOR_VERSION = 1
MINOR_VERSION = $(shell git rev-list master --count)
PATCH_VERSION = 0

VERSION = "$(MAJOR_VERSION).$(MINOR_VERSION).$(PATCH_VERSION)"

HASH = $(shell git rev-parse --short HEAD)
DATE = $(shell go run tools/build-date.go)

# build flags
BUILD_FLAGS = -ldflags "-s -w \
	-X main.Version=$(VERSION) \
	-X main.CommitHash=$(HASH) \
	-X 'main.CompileDate=$(DATE)'"

EXECUTABLE = badrobot
EXECUTABLES = \
	darwin-amd64-$(EXECUTABLE) \
	linux-amd64-$(EXECUTABLE) \
	linux-arm-7-$(EXECUTABLE)

EXECUTABLE_TARGETS = $(EXECUTABLES:%=bin/%)

all:
    $(MAKE) clean
	$(MAKE) $(EXECUTABLE_TARGETS)

# arm
bin/linux-arm-7-$(EXECUTABLE):
	GOARM=7 GOARCH=arm GOOS=linux go build $(BUILD_FLAGS) -o "$@"

# amd64
bin/darwin-amd64-$(EXECUTABLE):
	GOARCH=amd64 GOOS=darwin go build $(BUILD_FLAGS) -o "$@"
bin/linux-amd64-$(EXECUTABLE):
	GOARCH=amd64 GOOS=linux go build $(BUILD_FLAGS) -o "$@"

clean:
	rm -rf bin

# Docker
DOCKER_IMAGE = linux-arm-7-badrobot

docker-clean:
	docker stop badrobot && docker rm badrobot

docker-build:
	docker build --pull -t $(DOCKER_IMAGE):$(shell ./bin/linux-arm-7-badrobot --version) .

docker-run:
	docker run -p 8001:8001 --link redis:redis  --restart always --name badrobot -d $(DOCKER_IMAGE):$(shell ./bin/linux-arm-7-badrobot --version) /linux-arm-7-badrobot --redis redis:6379

docker-deploy:
    $(MAKE) docker-build
    $(MAKE) docker-clean
    $(MAKE) docker-run

.PHONY: clean all