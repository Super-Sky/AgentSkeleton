APP := agentskeleton
PKG := ./cmd/agentskeleton
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE ?= $(shell date -u +%Y-%m-%d)
LDFLAGS := -X github.com/Super-Sky/AgentSkeleton/internal/app.Version=$(VERSION) -X github.com/Super-Sky/AgentSkeleton/internal/app.Commit=$(COMMIT) -X github.com/Super-Sky/AgentSkeleton/internal/app.Date=$(DATE)

.PHONY: build test smoke release-build clean

build:
	go build -o $(APP) $(PKG)

test:
	go test ./...

smoke:
	sh scripts/smoke_test.sh

release-build:
	go build -ldflags "$(LDFLAGS)" -o $(APP) $(PKG)

clean:
	rm -f $(APP)
