PKGS := $(shell go list ./... | grep -v /vendor)

.PHONY: test
test: lint
	go test $(PKGS) -count=1

BIN_DIR := $(GOPATH)/bin
GOLINT := $(BIN_DIR)/golint

$(GOLINT):
	go get -u golang.org/x/lint/golint

.PHONY: lint
lint: $(GOLINT)
	$(GOLINT) ./...

BINARY := secretdb
VERSION ?= vlatest
PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS): test
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64 -ldflags="-w -s" -trimpath cmd/main.go

.PHONY: release
release: windows linux darwin
