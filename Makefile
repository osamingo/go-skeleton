.DEFAULT_GOAL := help

PROJECT_NAME = go-skeleton
BUILD_DIR    = build

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
.PHONY: help

build: ## Build binary
	@go build -v -o $(BUILD_DIR)/$(PROJECT_NAME) .
.PHONY: build

clean: ## Clean build output
	@rm -fr $(BUILD_DIR)
.PHONY: clean

dep: ## Ensure dependencies
	@go get github.com/golang/dep/cmd/dep
	@dep ensure -v
	@find vendor ! -iname "*.go" -type f -delete
	@find vendor -type l -delete
.PHONY: dep
