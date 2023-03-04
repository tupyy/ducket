.PHONY: help tools build check run logs

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

define COLORIZE
sed -u -e "s/\\\\\"/'/g; \
s/method=\([^ ]*\)/method=$(COLOR_BLUE)\1$(RESET_COLOR)/g;        \
s/error=\"\([^\"]*\)\"/error=\"$(COLOR_RED)\1$(RESET_COLOR)\"/g;  \
s/msg=\"\([^\"]*\)\"/msg=\"$(COLOR_YELLOW)\1$(RESET_COLOR)\"/g;   \
s/level=trace/level=$(COLOR_LEVEL_TRACE)trace$(RESET_COLOR)/g;    \
s/level=debug/level=$(COLOR_LEVEL_DEBUG)debug$(RESET_COLOR)/g;    \
s/level=info/level=$(COLOR_LEVEL_INFO)info$(RESET_COLOR)/g;       \
s/level=warning/level=$(COLOR_LEVEL_WARN)warning$(RESET_COLOR)/g; \
s/level=error/level=$(COLOR_LEVEL_ERROR)error$(RESET_COLOR)/g;    \
s/level=fatal/level=$(COLOR_LEVEL_FATAL)fatal$(RESET_COLOR)/g"
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
	go mod vendor

#help build.vendor.full: retrieve all the dependencies after cleaning the go.sum
build.vendor.full:
	@rm -fr $(CURDIR)/vendor
	go mod tidy
	go mod vendor

build.local:
	go build -o $(BUILD_DIR)/$(NAME) main.go

DB_HOST=localhost
DB_PORT=5433
ROOT_USER=postgres
ROOT_PWD=postgres
PGPASSFILE=$(CURDIR)/sql/.pgpass
PSQL_COMMAND=PGPASSFILE=$(PGPASSFILE) psql --quiet --host=$(DB_HOST) --port=$(DB_PORT) -v ON_ERROR_STOP=on

#help postgres.setup: Setup postgres from scratch
postgres.setup: postgres.setup.init postgres.setup.tables

#help postgres.setup.clean: cleans postgres from all created resources
postgres.setup.clean:
	$(PSQL_COMMAND) --user=$(ROOT_USER) -f sql/clean.sql

#help postgres.setup.init: init the database
postgres.setup.init:
	$(PSQL_COMMAND) --dbname=postgres --user=$(ROOT_USER) \
		-f sql/init.sql

#help postgres.setup.users: init postgres users
postgres.setup.tables:
	$(PSQL_COMMAND) --dbname=finance --user=$(ROOT_USER) \
		-f sql/tables.sql

BASE_CONNSTR="postgresql://$(ROOT_USER):$(ROOT_PWD)@$(DB_HOST):$(DB_PORT)"
GEN_CMD=$(TOOLS_DIR)/gen --sqltype=postgres \
	--module=github.com/tupyy/finance/internal/repo/models --exclude=schema_migrations \
	--gorm --no-json --no-xml --overwrite --out $(CURDIR)/internal/repo/

#help generate.models: generate models for the database
generate.models:
	sh -c '$(GEN_CMD) --connstr "$(BASE_CONNSTR)/finance?sslmode=disable"  --model=models --database finance' 						# Generate models for the DB tables
