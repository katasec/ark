.ONESHELL:
SHELL = /bin/bash
.DEFAULT_GOAL := help
.PHONY: worker


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
worker: ## Build Worker container
	@IMAGE_NAME=arkworker;
	@source ./scripts/build.sh;
	@build
workerpush: ## Build and push worker container
	@IMAGE_NAME=arkworker;
	@source ./scripts/build.sh;
	@buildAndPush
checkbuild: ## Check app can build
	go build -o /dev/null

#	docker build -t ghcr.io/katasec/arkapi:0.01 .
