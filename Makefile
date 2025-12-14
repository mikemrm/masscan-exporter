PROJECT_NAME = masscan-exporter

SERVICE_FILES = $(find . -type f -name '*.go')

export GITHUB_REPOSITORY = mikemrm/masscan-exporter

default: all

.PHONY: help
help: ## Show this help.
	@grep -h "##" $(MAKEFILE_LIST) | grep -v grep | sed -e 's/:.*##/#/' | column -c 2 -t -s#

.PHONY: all
all: test build image  ## Tests and builds the binary and container images.

.PHONY: test
test:  ## Runs go tests
	go test ./...

.PHONY: build
build: $(PROJECT_NAME)  ## Builds a binary for the current os/arch.

$(PROJECT_NAME): $(SERVICE_FILES)
	go tool goreleaser build --clean --snapshot --single-target -o $@

.PHONY: image
image:  ## Builds all container images.
	go tool goreleaser release --clean --snapshot --verbose
	@echo; echo "Container images:"
	@jq -r '.[] | select(.type == "Docker Image") | "  - " + .name' dist/artifacts.json
	@target_arch="$$(docker version -f '{{ .Client.Arch }}')"; \
	target_image="$$(jq --arg target_arch "-$$target_arch" -r 'map(select(.type == "Docker Image" and (.name | endswith($$target_arch))) | .name) | first // ""' dist/artifacts.json)"; \
	if [ -n "$$target_image" ]; then \
		docker tag "$$target_image" "ghcr.io/$(GITHUB_REPOSITORY):dev" && \
		echo "  - ghcr.io/$(GITHUB_REPOSITORY):dev (aliased to $$(echo "$$target_image" | cut -d : -f 2))"; \
	fi
