BINARY  = tryl
DIST    = dist
LDFLAGS = -s -w

.PHONY: build dist clean install help

## build: compile for the current host OS/arch (fast dev workflow)
build:
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -trimpath -o $(BINARY) .

## dist: cross-compile release binaries for all platforms into dist/
dist: \
	$(DIST)/tryl-linux-amd64 \
	$(DIST)/tryl-linux-arm64 \
	$(DIST)/tryl-darwin-amd64 \
	$(DIST)/tryl-darwin-arm64

$(DIST)/tryl-linux-amd64:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -trimpath -o $@ .

$(DIST)/tryl-linux-arm64:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=linux  GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -trimpath -o $@ .

$(DIST)/tryl-darwin-amd64:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -trimpath -o $@ .

$(DIST)/tryl-darwin-arm64:
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -trimpath -o $@ .

## clean: remove dist/ directory and local binary
clean:
	rm -rf $(DIST) $(BINARY)

## install: build for current host and install to $HOME/.local/bin
install: build
	mkdir -p $(HOME)/.local/bin
	cp $(BINARY) $(HOME)/.local/bin/$(BINARY)

## help: list available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## /  /'
