.PHONY: help tools build check run logs

help: help.all
tools: tools.get
build: build.local
check: check.imports check.fmt check.lint check.test
run: run.local
logs: logs.podman

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
# Tools targets     #
#####################

TOOLS_DIR=$(CURDIR)/tools/bin

.PHONY: tools.clean tools.get

#help tools.clean: remove everything in the tools/bin directory
tools.clean:
	rm -fr $(TOOLS_DIR)/*

#help tools.get: retrieve all the tools specified in gex
tools.get:
	cd $(CURDIR)/tools && go generate tools.go


#####################
# Build targets     #
#####################

VERSION=$(shell cat VERSION)
GIT_COMMIT=$(shell git rev-list -1 HEAD --abbrev-commit)

IMAGE_TAG=$(VERSION)-$(GIT_COMMIT)
IMAGE_NAME=cloud/continental/ctp/$(NAME)

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

#help build.podman: build a podman image
build.podman:
	podman_BUILDKIT=1 podman build --ssh default --build-arg build_args="$(BUILD_ARGS)" --build-arg REGISTRY_HOSTNAME=$(REPO_CI_PULL) -t $(IMAGE_NAME):$(IMAGE_TAG) -f podmanfile .

#help build.get.imagename: Allows to get the name of the service (for the CI)
build.get.imagename:
	@echo -n $(IMAGE_NAME)

#help build.get.tag: Allows to get the tag of the service (for the CI)
build.get.tag:
	@echo -n $(IMAGE_TAG)


#####################
# Check targets     #
#####################

LINT_COMMAND=golangci-lint run
FILES_LIST=$(shell ls -d */ | grep -v -E "vendor|tools|target")
MODULE_NAME=$(shell head -n 1 go.mod | cut -d '/' -f 3)

.PHONY: check.fmt check.imports check.lint check.test check.licenses check.get.tools.image check.prepare

check.prepare:
	@podman pull $(TOOLS_PODMAN_IMAGE)

#help check.fmt: format go code
check.fmt: check.prepare
	podman run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)"  $(TOOLS_PODMAN_IMAGE) sh -c 'gofumpt -s -w $(FILES_LIST)'

#help check.imports: fix and format go imports
check.imports: check.prepare
	@# Removes blank lines within import block so that goimports does its magic in a deterministic way
	find $(FILES_LIST) -type f -name "*.go" | xargs -L 1 sed -i '/import (/,/)/{/import (/n;/)/!{/^$$/d}}'
	@# Fine tune putting conti repos and then ctp at the end of the block
	podman run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(GOCACHE_FLAGS) $(TOOLS_podman_IMAGE) sh -c 'goimports -w -local github.com/tupyy/finance $(FILES_LIST)'
	podman run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(GOCACHE_FLAGS) $(TOOLS_podman_IMAGE) sh -c 'goimports -w -local github.com/tupyy/finance/$(MODULE_NAME) $(FILES_LIST)'


#help check.lint: check if the go code is properly written, rules are in .golangci.yml
check.lint: check.prepare
	podman run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(GOCACHE_FLAGS) $(TOOLS_podman_IMAGE) sh -c '$(LINT_COMMAND)'

#help check.test: execute go tests, if using test container set TEST_CONTAINER_FLAGS in custom.mk
check.test: check.prepare
	podman run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(GOCACHE_FLAGS) $(TEST_CONTAINER_FLAGS) $(TOOLS_PODMAN_IMAGE) sh -c 'go test -mod=vendor ./...'

#help check.licenses: check if the thirdparties' licences are whitelisted (in .wwhrd.yml)
check.licenses: check.prepare
	podman run --rm -v $(CURDIR):$(CURDIR) -w="$(CURDIR)" $(TOOLS_podman_IMAGE) sh -c 'wwhrd check'

#help check.get.tools.image: returns the name of the podman image used for the ci tools
check.get.tools.image:
	@echo -n $(TOOLS_podman_IMAGE)


#####################
# Run               #
#####################

.PHONY: run.podman.stop run.podman.logs run.infra run.infra.stop 

#help run.podman.stop: stop the container of the application
run.podman.stop:
	podman stop $(NAME)

#help run.podman.logs: display logs from the application in the container
run.podman.logs:
	@podman logs -f $(NAME) | $(COLORIZE)

#help run.infra: start podman-compose in resources/ folder
run.infra:
	podman-compose -f $(CURDIR)/build/docker-compose.yaml up -d

#help run.infra.stop: stop podman-compose in resources/ folder
run.infra.stop:
	podman-compose -f $(CURDIR)/build/docker-compose.yaml down


##@ Infra
.PHONY: postgres.setup.clean postgres.setup.init postgres.setup.tables postgres.setup.migrations

DB_HOST=localhost
DB_PORT=5434
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
	--gorm --no-json --no-xml --overwrite --out $(CURDIR)/internal/repo/models

#help generate.models: generate models for the database
generate.postgres.models:
	sh -c '$(GEN_CMD) --connstr "$(BASE_CONNSTR)/finance?sslmode=disable"  --model=pg --database finance' 						# Generate models for the DB tables

#####################
# Logs targets      #
#####################

.PHONY: logs.k8s

#help logs.k8s: display pod logs, use ENV to choose target cluster
logs:
	@kubectl logs -n $(NAMESPACE) -f $(DEPLOYMENT) $(if $(TAIL), --tail=$(TAIL)) | $(COLORIZE)


#####################
# Include section   #
#####################
