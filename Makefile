BIN      := git-release
INSTALL  := /usr/local/bin/$(BIN)
GOFLAGS  := -ldflags="-s -w"

.PHONY: build test lint install clean

build:
	go build $(GOFLAGS) -o $(BIN) ./cmd/git-release

test:
	go test ./...

test-verbose:
	go test -v ./...

test-pkg:  ## Run tests for a single package: make test-pkg PKG=./internal/semver
	go test -v $(PKG)

lint:
	go vet ./...

install: build
	install -m 755 $(BIN) $(INSTALL)

clean:
	rm -f $(BIN)
