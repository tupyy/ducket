.PHONY: build vendor test run image container-run container-stop

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

run: build
	$(BUILD_DIR)/$(NAME) serve

image:
	podman build -f Containerfile --build-arg GIT_SHA=$(GIT_COMMIT) -t $(IMAGE_NAME) .

container-run: image
	@podman rm -f $(CONTAINER_NAME) 2>/dev/null || true
	podman run -d \
		--name $(CONTAINER_NAME) \
		-p 8080:8080 \
		-v $(DB_VOLUME):/app/data \
		$(IMAGE_NAME)
	@echo "running at http://localhost:8080"

container-stop:
	podman stop $(CONTAINER_NAME) 2>/dev/null || true
	podman rm $(CONTAINER_NAME) 2>/dev/null || true
