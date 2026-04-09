APP := agentskeleton
PKG := ./cmd/agentskeleton
VERSION ?= dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE ?= $(shell date -u +%Y-%m-%d)
LDFLAGS := -X github.com/Super-Sky/AgentSkeleton/internal/app.Version=$(VERSION) -X github.com/Super-Sky/AgentSkeleton/internal/app.Commit=$(COMMIT) -X github.com/Super-Sky/AgentSkeleton/internal/app.Date=$(DATE)

.PHONY: build test smoke release-build clean new-validation-report

build:
	go build -o $(APP) $(PKG)

test:
	go test ./...

smoke:
	sh scripts/smoke_test.sh

release-build:
	go build -ldflags "$(LDFLAGS)" -o $(APP) $(PKG)

new-validation-report:
	@echo "usage: make new-validation-report HOST=codex FILE=codex-scenario-3 TITLE='Codex Validation Report: Scenario 3'"
	@test -n "$(HOST)" && test -n "$(FILE)" && test -n "$(TITLE)"
	sh scripts/new_validation_report.sh "$(HOST)" "$(FILE)" "$(TITLE)"

clean:
	rm -f $(APP)
