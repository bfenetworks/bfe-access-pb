WORKROOT := $(shell pwd)
TOOLDIR  := $(WORKROOT)/bfe-pblog-tool
OUTDIR   := $(WORKROOT)/output
DISTDIR  := $(WORKROOT)/dist
BIN_NAME := bfePblogTool

PROTOC     ?= /opt/protoc
VERSION    ?= $(shell cat $(WORKROOT)/VERSION)
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

GO         := go
GOBUILD    := $(GO) build
GOTEST     := $(GO) test
GOCLEAN    := $(GO) clean
GOINSTALL  := $(GO) install

LDFLAGS := -X main.version=$(VERSION)

.PHONY: all build test clean release proto

all: build

proto:
	$(GOINSTALL) google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.0
	export PATH=$$PATH:$$($(GO) env GOPATH)/bin && \
	$(PROTOC) --go_out=./ --go_opt=paths=source_relative bfe_access_pb/bfe_access.proto
	@echo "Proto generation succeed!"

build:
	@mkdir -p $(OUTDIR)/bin
	cd $(TOOLDIR) && $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(OUTDIR)/bin/$(BIN_NAME) ./cmd/
	@echo "Build success: $(OUTDIR)/bin/$(BIN_NAME) ($(VERSION))"

test:
	cd $(TOOLDIR) && $(GOTEST) ./...

clean:
	cd $(TOOLDIR) && $(GOCLEAN)
	rm -rf $(OUTDIR)
	rm -rf $(DISTDIR)

release:
	@echo "Building release packages for $(BIN_NAME) $(VERSION)..."
	@for platform in \
		"darwin/amd64" \
		"linux/amd64" \
		"linux/arm64" \
		"windows/amd64"; do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		PKG_DIR="$(DISTDIR)/$(BIN_NAME)_$(VERSION)_$${GOOS}_$${GOARCH}"; \
		BIN="$(BIN_NAME)"; \
		FLAGS="$(LDFLAGS)"; \
		if [ "$${GOOS}" = "linux" ]; then \
			FLAGS="$${FLAGS} -extldflags=-static"; \
		fi; \
		if [ "$${GOOS}" = "windows" ]; then \
			BIN="$(BIN_NAME).exe"; \
		fi; \
		echo "  -> $$GOOS/$$GOARCH"; \
		rm -rf "$${PKG_DIR}"; \
		mkdir -p "$${PKG_DIR}"; \
		cd $(TOOLDIR) && GOOS=$${GOOS} GOARCH=$${GOARCH} $(GOBUILD) -ldflags "$${FLAGS}" -o "$${PKG_DIR}/$${BIN}" ./cmd/; \
		tar -czf "$${PKG_DIR}.tar.gz" -C $(DISTDIR) "$(BIN_NAME)_$(VERSION)_$${GOOS}_$${GOARCH}"; \
		rm -rf "$${PKG_DIR}"; \
	done
	@echo "Release packages:"
	@ls -lh $(DISTDIR)/*.tar.gz
