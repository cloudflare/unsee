NAME    := unsee
VERSION := $(shell git describe --tags --always --dirty='-dev')
GO      := GO15VENDOREXPERIMENT=1 go
PROMU   := $(GOPATH)/bin/promu
pkgs     = $(shell $(GO) list ./... | grep -v -E '/vendor/')

PREFIX                  ?= $(shell pwd)
BIN_DIR                 ?= $(shell pwd)
DOCKER_IMAGE_NAME       ?= unsee
DOCKER_IMAGE_TAG        ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

# Alertmanager instance used when running locally, points to mock data
MOCK_PATH         := $(CURDIR)/internal/mock/0.12.0
ALERTMANAGER_URI  := "file://$(MOCK_PATH)"
# Listen port when running locally
PORT := 8080

SOURCES       := $(wildcard *.go) $(wildcard */*.go) $(wildcard */*/*.go)
ASSET_SOURCES := $(wildcard assets/*/* assets/*/*/*)

GO_BINDATA_MODE := prod
GIN_DEBUG := false
ifdef DEBUG
	GO_BINDATA_FLAGS = -debug
	GO_BINDATA_MODE  = debug
	GIN_DEBUG = true
	DOCKER_ARGS = -v $(CURDIR)/assets:$(CURDIR)/assets:ro
endif

.DEFAULT_GOAL := $(NAME)

.build/deps-build-go.ok:
	@mkdir -p .build
	$(GO) get -u github.com/golang/dep/cmd/dep
	$(GO) get -u github.com/jteeuwen/go-bindata/...
	$(GO) get -u github.com/elazarl/go-bindata-assetfs/...
	touch $@

.build/deps-lint-go.ok:
	@mkdir -p .build
	$(GO) get -u github.com/golang/lint/golint
	touch $@

.build/deps-build-node.ok: package.json package-lock.json
	@mkdir -p .build
	npm install
	touch $@

.build/artifacts-bindata_assetfs.%:
	@mkdir -p .build
	rm -f .build/artifacts-bindata_assetfs.*
	touch $@

.build/artifacts-webpack.ok: .build/deps-build-node.ok $(ASSET_SOURCES) webpack.config.js
	@mkdir -p .build
	$(CURDIR)/node_modules/.bin/webpack
	touch $@

bindata_assetfs.go: .build/deps-build-go.ok .build/artifacts-bindata_assetfs.$(GO_BINDATA_MODE) .build/vendor.ok .build/artifacts-webpack.ok
	go-bindata-assetfs $(GO_BINDATA_FLAGS) -prefix assets -nometadata assets/templates/... assets/static/dist/...

$(NAME): .build/deps-build-go.ok .build/vendor.ok bindata_assetfs.go $(SOURCES)
	$(GO) build -ldflags "-X main.version=$(VERSION)"

.build/vendor.ok: .build/deps-build-go.ok Gopkg.lock Gopkg.toml
	dep ensure
	dep prune
	touch $@

build: promu
	@echo ">> building binaries"
	@$(PROMU) build --prefix $(PREFIX)

tarball: promu
	@echo ">> building release tarball"
	@$(PROMU) tarball --prefix $(PREFIX) $(BIN_DIR)

docker:
	@echo ">> building docker image"
	@docker build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

promu:
	@GOOS=$(shell uname -s | tr A-Z a-z) \
	GOARCH=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m))) \
	$(GO) get -u github.com/prometheus/promu

.PHONY: vendor
vendor: .build/deps-build-go.ok
	dep ensure
	dep prune

.PHONY: vendor-update
vendor-update: .build/deps-build-go.ok
	dep ensure -update
	dep prune

.PHONY: webpack
webpack: .build/artifacts-webpack.ok

.PHONY: clean
clean:
	rm -fr .build bindata_assetfs.go $(NAME)

.PHONY: run
run: $(NAME)
	ALERTMANAGER_URI=$(ALERTMANAGER_URI) \
	LABELS_COLOR_UNIQUE="@receiver instance cluster" \
	LABELS_COLOR_STATIC="job" \
	DEBUG="$(GIN_DEBUG)" \
	FILTER_DEFAULT="@state=active" \
	PORT=$(PORT) \
	./$(NAME)

.PHONY: docker-image
docker-image:
	docker build --build-arg VERSION=$(VERSION) -t $(NAME):$(VERSION) .

.PHONY: run-docker
run-docker: docker-image
	@docker rm -f $(NAME) || true
	docker run \
	    --name $(NAME) \
	    $(DOCKER_ARGS) \
	    -v $(MOCK_PATH):$(MOCK_PATH) \
	    -e ALERTMANAGER_URI=$(ALERTMANAGER_URI) \
	    -e LABELS_COLOR_UNIQUE="instance cluster" \
	    -e LABELS_COLOR_STATIC="job" \
	    -e DEBUG="$(GIN_DEBUG)" \
	    -e PORT=$(PORT) \
	    -p $(PORT):$(PORT) \
	    $(NAME):$(VERSION)

.PHONY: lint-go
lint-go: .build/deps-lint-go.ok
	golint ./... | (egrep -v "^vendor/|^bindata_assetfs.go" || true)

.PHONY: lint-js
lint-js: .build/deps-build-node.ok
	$(CURDIR)/node_modules/.bin/eslint --quiet assets/static/*.js

.PHONY: lint
lint: lint-go lint-js

# Creates mock bindata_assetfs.go with source assets rather than webpack generated ones
.PHONY: mock-assets
mock-assets: .build/deps-build-go.ok
	mkdir -p $(CURDIR)/assets/static/dist/templates
	cp $(CURDIR)/assets/static/*.* $(CURDIR)/assets/static/dist/
	touch $(CURDIR)/assets/static/dist/templates/loader_unsee.html
	touch $(CURDIR)/assets/static/dist/templates/loader_shared.html
	touch $(CURDIR)/assets/static/dist/templates/loader_help.html
	go-bindata-assetfs -prefix assets -nometadata assets/templates/... assets/static/dist/...
	# force assets rebuild on next make run
	rm -f .build/bindata_assetfs.*

.PHONY: test-go
test-go: .build/vendor.ok
	$(GO) test -bench=. -cover `go list ./... | grep -v /vendor/`

.PHONY: test-js
test-js: .build/deps-build-node.ok
	npm test

.PHONY: test
test: lint test-go test-js
