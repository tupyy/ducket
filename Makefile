.PHONY: build vendor test lint run run.ui image container.run container.stop

NAME = ducket
BUILD_DIR = target
GIT_COMMIT = $(shell git rev-list -1 HEAD --abbrev-commit)

IMAGE_NAME = ducket:latest
CONTAINER_NAME = ducket
DB_VOLUME = ducket-data

build:
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-X main.sha=$(GIT_COMMIT)" -o $(BUILD_DIR)/$(NAME) .

vendor:
	go mod tidy
	go mod vendor

test:
	go test ./...

GOBIN ?= $(shell go env GOPATH)/bin
GOLANGCI_LINT_VERSION := v2.10.1
GOLANGCI_LINT := $(GOBIN)/golangci-lint

.PHONY: check-golangci-lint-version
check-golangci-lint-version:
	@if [ -f '$(GOLANGCI_LINT)' ]; then \
		installed=$$('$(GOLANGCI_LINT)' version 2>/dev/null | sed -n 's/.*version \([0-9.]*\).*/\1/p' | head -1); \
		required=$$(echo '$(GOLANGCI_LINT_VERSION)' | sed 's/^v//'); \
		if [ -n "$$installed" ] && [ "$$installed" != "$$required" ]; then \
			rm -f '$(GOLANGCI_LINT)'; \
		fi; \
	fi

$(GOLANGCI_LINT):
	@mkdir -p $(GOBIN)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VERSION)

lint: check-golangci-lint-version $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run

run: build
	$(BUILD_DIR)/$(NAME) run

run.ui:
	cd ui && npm run start:dev

image:
	podman build -f Containerfile --build-arg GIT_SHA=$(GIT_COMMIT) -t $(IMAGE_NAME) .

container.run:
	@podman rm -f $(CONTAINER_NAME) 2>/dev/null || true
	podman run -d \
		--name $(CONTAINER_NAME) \
		-p 8080:8080 \
		-v $(DB_VOLUME):/app/data \
		$(IMAGE_NAME)
	@echo "running at http://localhost:8080"

container.stop:
	podman stop $(CONTAINER_NAME) 2>/dev/null || true
	podman rm $(CONTAINER_NAME) 2>/dev/null || true
