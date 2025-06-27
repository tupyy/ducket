# Include container management targets
include Makefile.docker

.PHONY: help tools build check run logs container-help

help: help.all
build: build.local

# Colors used in this Makefile
escape=$(shell printf '\033')
RESET_COLOR=$(escape)[0m
COLOR_YELLOW=$(escape)[38;5;220m
COLOR_RED=$(escape)[91m
COLOR_BLUE=$(escape)[94m

COLOR_LEVEL_TRACE=$(escape)[38;5;87m
COLOR_LEVEL_DEBUG=$(escape)[38;5;87m
COLOR_LEVEL_INFO=$(escape)[92m
COLOR_LEVEL_WARN=$(escape)[38;5;208m
COLOR_LEVEL_ERROR=$(escape)[91m
COLOR_LEVEL_FATAL=$(escape)[91m

PODMAN ?= podman
POSTGRES_IMAGE ?= docker.io/library/postgres:17

define COLORIZE
sed -u -e "s/\\\\\"/'/g; \
s/method=\([^ ]*\)/method=$(COLOR_BLUE)\1$(RESET_COLOR)/g;        \
s/error=\"\([^\"]*\)\"/error=\"$(COLOR_RED)\1$(RESET_COLOR)\"/g;  \
s/msg=\"\([^\"]*\)\"/msg=\"$(COLOR_YELLOW)\1$(RESET_COLOR)\"/g;   \
s/trace/$(COLOR_LEVEL_TRACE)trace$(RESET_COLOR)/g;    \
s/debug/$(COLOR_LEVEL_DEBUG)debug$(RESET_COLOR)/g;    \
s/info/$(COLOR_LEVEL_INFO)info$(RESET_COLOR)/g;       \
s/warning/$(COLOR_LEVEL_WARN)warning$(RESET_COLOR)/g; \
s/error/$(COLOR_LEVEL_ERROR)error$(RESET_COLOR)/g;    \
s/fatal/$(COLOR_LEVEL_FATAL)fatal$(RESET_COLOR)/g"
endef

#####################
# Help targets      #
#####################

.PHONY: help.highlevel help.all

#help help.highlevel: show help for high level targets. Use 'make help.all' to display all help messages
help.highlevel:
	@grep -hE '^[a-z_-]+:' $(MAKEFILE_LIST) | LANG=C sort -d | \
	awk 'BEGIN {FS = ":"}; {printf("$(COLOR_YELLOW)%-25s$(RESET_COLOR) %s\n", $$1, $$2)}'

#help help.all: display all targets' help messages
help.all:
	@grep -hE '^#help|^[a-z_-]+:' $(MAKEFILE_LIST) | sed "s/#help //g" | LANG=C sort -d | \
	awk 'BEGIN {FS = ":"}; {if ($$1 ~ /\./) printf("    $(COLOR_BLUE)%-21s$(RESET_COLOR) %s\n", $$1, $$2); else printf("$(COLOR_YELLOW)%-25s$(RESET_COLOR) %s\n", $$1, $$2)}'


#####################
# Build targets     #
#####################

GIT_COMMIT=$(shell git rev-list -1 HEAD --abbrev-commit)

IMAGE_TAG=$(GIT_COMMIT)
IMAGE_NAME=$(NAME)
NAME=finante
BUILD_DIR ?= target
TOOLS_DIR=$(CURDIR)/tools/bin

GOCACHE?=$(shell go env GOCACHE 2>/dev/null)

ifneq "$(strip $(GOCACHE))" ""
    GOCACHE_FLAGS=-v $(GOCACHE):/cache/go -e GOCACHE=/cache/go -e GOLANGCI_LINT_CACHE=/cache/go
endif

.PHONY: build.prepare build.vendor build.vendor.full build.podman build.get.imagename build.get.tag

#help build.prepare: prepare target/ folder
build.prepare:
	@mkdir -p $(CURDIR)/target
	@rm -f $(CURDIR)/target/$(NAME)

#help build.vendor: retrieve all the dependencies used for the project
build.vendor:
	go mod tidy
	go mod vendor

#help build.vendor.full: retrieve all the dependencies after cleaning the go.sum
build.vendor.full:
	@rm -fr $(CURDIR)/vendor
	go mod tidy
	go mod vendor

build.local:
	go build -o $(BUILD_DIR)/$(NAME) main.go

run:
	$(BUILD_DIR)/$(NAME) serve | $(COLORIZE)

DB_HOST=localhost
DB_PORT=5432
ROOT_USER=postgres
ROOT_PWD=postgres
CONNSTR="postgresql://$(ROOT_USER):$(ROOT_PWD)@$(DB_HOST):$(DB_PORT)"

postgres.start.dev:
	$(PODMAN) run --rm -p $(DB_PORT):5432 --name pg-finante -e POSTGRES_PASSWORD=$(ROOT_PWD) -d $(POSTGRES_IMAGE)

postgres.start.test:
	$(PODMAN) run --rm -p $(DB_PORT):5432 --name pg-test -e POSTGRES_PASSWORD=$(ROOT_PWD) -d $(POSTGRES_IMAGE)

postgres.stop:
	$(PODMAN) stop pg-finante

postgres.stop.test:
	$(PODMAN) stop pg-test

postgres.migrate:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(CONNSTR) GOOSE_MIGRATION_DIR=$(CURDIR)/pkg/migrations/sql goose up

postgres.migrate.test:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(CONNSTR) GOOSE_MIGRATION_DIR=$(CURDIR)/pkg/migrations/sql goose up

#####################
# Container targets #
#####################

#help container-help: show container management help (from Makefile.docker)
container-help:
	@$(MAKE) -f Makefile.docker help

# Container aliases for convenience
#help container-build: build application container image
container-build:
	@$(MAKE) -f Makefile.docker build

#help container-run: start application and database containers
container-run:
	@$(MAKE) -f Makefile.docker podman-run

#help container-dev: start containers in development mode
container-dev:
	@$(MAKE) -f Makefile.docker dev

#help container-stop: stop application container
container-stop:
	@$(MAKE) -f Makefile.docker stop

#help container-down: stop all containers
container-down:
	@$(MAKE) -f Makefile.docker down

#help container-logs: follow application logs
container-logs:
	@$(MAKE) -f Makefile.docker logs

#help container-health: check container health status
container-health:
	@$(MAKE) -f Makefile.docker health

#help container-clean: remove all containers and volumes
container-clean:
	@$(MAKE) -f Makefile.docker clean

#help container-setup: quick setup (build and start containers)
container-setup:
	@$(MAKE) -f Makefile.docker setup

# Database container operations
#help container-db-shell: connect to database shell
container-db-shell:
	@$(MAKE) -f Makefile.docker db-shell

#help container-db-backup: backup database
container-db-backup:
	@$(MAKE) -f Makefile.docker db-backup

# Quick development workflow aliases
#help dev: start development environment (alias for container-dev)
dev: container-dev

#help prod: start production environment (alias for container-run)
prod: container-run

#help deploy: full deployment setup (alias for container-setup)
deploy: container-setup

