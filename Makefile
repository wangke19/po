VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS  = -X github.com/wangke19/po/internal/build.Version=$(VERSION) \
           -X github.com/wangke19/po/internal/build.Date=$(DATE)

.PHONY: build test install clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/po ./cmd/po

test:
	go test ./...

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/po

clean:
	rm -rf bin/
