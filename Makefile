VERSION ?= 0.0.0
LDFLAGS ?= -ldflags "-s -w -X 'main.Version=$(VERSION)'"

# HACK: make [target] [ARGS...]
ARGS = $(filter-out $@,$(MAKECMDGOALS))

# HACK: nothing undefined target
%:
	@:

all: run

run:
	go run $(LDFLAGS) . $(ARGS)

build:
	go build $(LDFLAGS) -o closest .

fmt:
	@go fmt ./...

test:
	@go test -v ./...

lint:
	@golangci-lint run

clean:
	@rm -f closest

install: build
	@mv closest $(GOPATH)/bin/ || mv closest /usr/local/bin/

.PHONY: all run build fmt test lint clean install
